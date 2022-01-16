package connection

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sync"
)

var ErrNotStarted = errors.New("connection not started")

type Connection struct {
	input    io.Reader
	output   io.Writer
	handlers map[string]Handler
	channels sync.Map
	running  bool
}

func New(input io.Reader, output io.Writer) *Connection {
	return &Connection{
		input:    input,
		output:   output,
		handlers: map[string]Handler{},
		channels: sync.Map{},
	}
}

func (c *Connection) Start() {
	if !c.running {
		go c.worker()
		c.running = true
	}
}

func (c *Connection) Stop() {
	c.running = false
}

func (c *Connection) Dispatch(method string, params Params) (interface{}, error) {
	r := Response{}
	if !c.running {
		return r, ErrNotStarted
	}
	if params == nil {
		params = map[string]interface{}{}
	}
	id, err := getRandomId()
	if err != nil {
		return r, err
	}
	data, err := json.Marshal(Request{Id: id, Method: method, Params: params})
	if err != nil {
		return r, err
	}
	data = append(data, '\n')
	c.channels.Store(id, make(chan Response))
	_, err = c.output.Write(data)
	if err != nil {
		return r, err
	}
	channel, _ := c.channels.Load(id)
	r = <-channel.(chan Response)
	if !r.Ok {
		return nil, errors.New(r.Result.(string))
	}
	return r.Result, nil
}

func (c *Connection) Respond(id string, ok bool, result interface{}) error {
	if !c.running {
		return ErrNotStarted
	}
	data, err := json.Marshal(Response{Id: id, Ok: ok, Result: result})
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = c.output.Write(data)
	return err
}

func (c *Connection) Handle(method string, handler Handler) {
	c.handlers[method] = handler
}

func (c *Connection) worker() {
	reader := bufio.NewReader(c.input)
	for {
		if !c.running {
			break
		}
		line, err := reader.ReadBytes('\n')
		go func() {
			if err != nil {
				c.running = false
				return
			}
			if len(line) < 3 {
				return
			}
			line = line[:len(line)-1]
			if d, ok := isRequest(line); ok {
				if handler, ok := c.handlers[d.Method]; ok {
					result, err := handler(d)
					if err != nil {
						c.Respond(d.Id, false, err.Error())
					} else if result != nil {
						c.Respond(d.Id, true, result)
					}
				}
			} else if r, ok := isResponse(line); ok {
				if channel, ok := c.channels.Load(r.Id); ok {
					channel.(chan Response) <- r
					c.channels.Delete(r.Id)
				}
			}
		}()
	}
}
