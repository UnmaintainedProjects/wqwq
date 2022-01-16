package tgcalls

func (calls *TGCalls) Stop() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("stop", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
