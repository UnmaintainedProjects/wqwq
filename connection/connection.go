package connection

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"sync"
)

var ErrNotStarted = errors.New("connection not started")

type Params = map[string]interface{}

type Handler = func(data Data) (interface{}, error)

type Data struct {
	Id     string `json:"id"`
	Event  string `json:"event"`
	Params Params `json:"params"`
}

func data(line []byte) (Data, bool) {
	r := Data{}
	err := json.Unmarshal(line, &r)
	if err != nil {
		return r, false
	}
	if r.Event == "" {
		return r, false
	}
	return r, true
}

type Response struct {
	Id     string      `json:"id"`
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

func response(line []byte) (Response, bool) {
	r := Response{}
	err := json.Unmarshal(line, &r)
	if err != nil {
		return r, false
	}
	if r.Result == nil {
		return r, false
	}
	return r, true
}

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

func id() (string, error) {
	s := ""
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return s, err
	}
	for _, i := range b {
		s += strconv.Itoa(int(i))
	}
	return s, nil
}

func (c *Connection) Dispatch(event string, params Params) (interface{}, error) {
	r := Response{}
	if !c.running {
		return r, ErrNotStarted
	}
	id, err := id()
	if err != nil {
		return r, err
	}
	data, err := json.Marshal(Data{Id: id, Event: event, Params: params})
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

func (c *Connection) Handle(event string, handler Handler) {
	c.handlers[event] = handler
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
			if d, ok := data(line); ok {
				if handler, ok := c.handlers[d.Event]; ok {
					result, err := handler(d)
					if err != nil {
						c.Respond(d.Id, false, err.Error())
					} else if result != nil {
						c.Respond(d.Id, true, result)
					}
				}
			} else if r, ok := response(line); ok {
				if channel, ok := c.channels.Load(r.Id); ok {
					channel.(chan Response) <- r
					c.channels.Delete(r.Id)
				}
			}
		}()
	}
}
