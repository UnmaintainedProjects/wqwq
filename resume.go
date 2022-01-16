package tgcalls

func (calls *TGCalls) Resume() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("resume", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
