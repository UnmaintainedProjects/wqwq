package tgcalls

func (calls *TGCalls) Stop() (bool, error) {
	if !calls.running {
		return false, ErrNotRunning
	}
	result, err := calls.conn.Dispatch("stop", nil)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}
