package tgcalls

import "github.com/gotd/td/tg"

func (calls *TGCalls) Stream(channel *tg.InputChannel, file string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch(
		"stream",
		map[string]interface{}{
			"chatId":     channel.ChannelID,
			"accessHash": channel.AccessHash,
			"isChat":     false,
			"file":       file,
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
			"chatId": chatId,
			"isChat": true,
			"file":   file,
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
		map[string]interface{}{"chatId": chatId},
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
		map[string]interface{}{"chatId": chatId},
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
		map[string]interface{}{"chatId": chatId},
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
		map[string]interface{}{"chatId": chatId},
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
		map[string]interface{}{"chatId": chatId},
	)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
