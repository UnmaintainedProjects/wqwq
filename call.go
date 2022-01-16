package tgcalls

import (
	"fmt"

	"github.com/gotd/td/tg"
)

func (calls *TGCalls) Stream(channel *tg.InputChannel, file string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch(
		"stream",
		map[string]interface{}{
			"id":   fmt.Sprint(channel.ChannelID),
			"file": file,
			"joinCallParams": map[string]interface{}{
				"accessHash": fmt.Sprint(channel.AccessHash),
				"isChannel":  true,
			},
		},
	)
	return err
}

func (calls *TGCalls) StreamChat(chatId int, file string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch(
		"stream",
		map[string]interface{}{
			"id":   fmt.Sprint(chatId),
			"file": file,
			"joinCallParams": map[string]interface{}{
				"isChannel": false,
			},
		},
	)
	return err
}

func (calls *TGCalls) Mute(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch(
		"mute",
		map[string]interface{}{"id": fmt.Sprint(chatId)},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}

func (calls *TGCalls) Unmute(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch(
		"unmute",
		map[string]interface{}{"id": fmt.Sprint(chatId)},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}

func (calls *TGCalls) Pause(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch(
		"pause",
		map[string]interface{}{"id": fmt.Sprint(chatId)},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}

func (calls *TGCalls) Resume(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch(
		"resume",
		map[string]interface{}{"id": fmt.Sprint(chatId)},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}

func (calls *TGCalls) Stop(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch(
		"stop",
		map[string]interface{}{"id": fmt.Sprint(chatId)},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
