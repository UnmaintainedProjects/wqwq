package tgcalls

func (calls *TGCalls) Unmute() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("unmute", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
