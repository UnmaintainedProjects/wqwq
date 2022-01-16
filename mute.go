package tgcalls

func (calls *TGCalls) Mute() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("mute", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
