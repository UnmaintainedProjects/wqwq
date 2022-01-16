package tgcalls

func (calls *TGCalls) Pause() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("pause", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
