package tgcalls

func (calls *TGCalls) Stream(chatId int64, file string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch("stream", map[string]interface{}{"chatId": chatId, "file": file})
	return err
}

func (calls *TGCalls) Mute(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("mute", map[string]interface{}{"chatId": chatId})
	if err != nil {
		return Err, err
	}
	return int(result.(int64)), nil
}

func (calls *TGCalls) Unmute(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("unmute", map[string]interface{}{"chatId": chatId})
	if err != nil {
		return Err, err
	}
	return int(result.(int64)), nil
}

func (calls *TGCalls) Pause(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("pause", map[string]interface{}{"chatId": chatId})
	if err != nil {
		return Err, err
	}
	return int(result.(int64)), nil
}

func (calls *TGCalls) Resume(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("resume", map[string]interface{}{"chatId": chatId})
	if err != nil {
		return Err, err
	}
	return int(result.(int64)), nil
}

func (calls *TGCalls) Stop(chatId int64) (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("stop", map[string]interface{}{"chatId": chatId})
	if err != nil {
		return Err, err
	}
	return int(result.(int64)), nil
}
