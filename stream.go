package tgcalls

func (calls *TGCalls) Stream(audio string, video string) error {
	if !calls.running {
		return ErrNotRunning
	}
	_, err := calls.conn.Dispatch(
		"stream",
		map[string]interface{}{"audio": audio, "video": video},
	)
	return err
}
