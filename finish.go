package tgcalls

func (calls *TGCalls) Finish() (int, error) {
	if !calls.running {
		return Err, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("finish", nil)
	if err != nil {
		return Err, err
	}
	return int(result.(float64)), nil
}
