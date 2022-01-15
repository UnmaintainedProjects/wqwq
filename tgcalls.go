package tgcalls

import (
	"context"
	"errors"
	"io"
	"os/exec"

	"github.com/gotd/td/tg"
	"github.com/gotgcalls/tgcalls/connection"
)

var (
	ErrNotRunning     = errors.New("tgcalls not running")
	ErrUnexpectedType = errors.New("got an unexpected type")
	ErrNoCall         = errors.New("no active call in the provided chat")
	ErrNoAccessHash   = errors.New("no access hash for the provided chat")
)

const (
	Ok           = 0
	NotMuted     = 1
	AlreadyMuted = 1
	NotPaused    = 1
	NotStreaming = 1
	NotInCall    = 2
	Err          = 3
)

type TGCalls struct {
	GetAccessHash func(chatId int64) int64
	api           *tg.Client
	ctx           context.Context
	cmd           *exec.Cmd
	out           io.ReadCloser
	in            io.WriteCloser
	conn          *connection.Connection
	running       bool
}

func New(api *tg.Client, ctx context.Context) *TGCalls {
	return &TGCalls{api: api, ctx: ctx}
}

func Start(calls *TGCalls) error {
	if calls.running {
		return nil
	}
	cmd := exec.Command("npx", "gotgcalls-server")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	calls.cmd = cmd
	calls.out = out
	calls.in = in
	calls.conn = connection.New(out, in)
	calls.conn.Handle("joinCall", func(data connection.Data) (interface{}, error) {
		return calls.joinCall(data.Params)
	})
	calls.conn.Start()
	calls.running = true
	return nil
}

func Stop(calls *TGCalls) error {
	if !calls.running {
		return nil
	}
	calls.conn.Stop()
	err := calls.cmd.Process.Kill()
	if err != nil {
		return err
	}
	calls.cmd = nil
	calls.out = nil
	calls.in = nil
	calls.conn = nil
	calls.running = false
	return nil
}
