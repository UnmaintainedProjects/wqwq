package tgcalls

import (
	"context"
	"errors"
	"io"
	"os/exec"

	"github.com/gotd/td/tg"
	"github.com/gotgcalls/tgcalls/connection"
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

var (
	DefaultName = "npx"
	DefaultArgs = []string{"gotgcalls-server"}
)

var (
	ErrNotRunning     = errors.New("tgcalls not running")
	ErrUnexpectedType = errors.New("got an unexpected type")
	ErrNoCall         = errors.New("no active call in the provided chat")
	ErrNoAccessHash   = errors.New("no access hash for the provided chat")
)

type TGCalls struct {
	ctx  context.Context
	chat *tg.InputChannel
	api  *tg.Client
	opts *TGCallsOpts

	cmd  *exec.Cmd
	conn *connection.Connection
	in   io.WriteCloser
	out  io.ReadCloser

	running bool
}

type TGCallsOpts struct {
	Cmd        *TGCallsCmdOpts
	JoinAs     tg.InputPeerClass
	InviteHash string
}

type TGCallsCmdOpts struct {
	Name string
	Args []string
}

func New(
	ctx context.Context,
	chat *tg.InputChannel,
	api *tg.Client,
	opts *TGCallsOpts,
) *TGCalls {
	return &TGCalls{
		ctx:  ctx,
		chat: chat,
		api:  api,
		opts: opts,
	}
}

func Start(calls *TGCalls) error {
	if calls.running {
		return nil
	}
	name := DefaultName
	args := DefaultArgs
	if calls.opts != nil {
		if calls.opts.Cmd != nil {
			name = calls.opts.Cmd.Name
			args = calls.opts.Cmd.Args
		}
	}
	cmd := exec.Command(name, args...)
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
	calls.conn.Handle(
		"joinCall",
		func(request connection.Request) (interface{}, error) {
			return calls.joinCall(request.Params)
		},
	)
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
