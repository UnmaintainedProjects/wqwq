package tgcalls

func (calls *TGCalls) Stream(file string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch(
		"stream",
		map[string]interface{}{"file": file},
	)
	return err
}
