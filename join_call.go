package tgcalls

import (
	"encoding/json"

	"github.com/gotd/td/tg"
)

func (calls *TGCalls) joinCall(params map[string]interface{}) (string, error) {
	chatId, ok := params["chatId"].(float64)
	if !ok {
		return "", ErrUnexpectedType
	}
	payload, ok := params["payload"].(map[string]interface{})
	if !ok {
		return "", ErrUnexpectedType
	}
	fullChannel, err := calls.api.ChannelsGetFullChannel(
		calls.ctx,
		&tg.InputChannel{
			ChannelID:  int64(chatId),
			AccessHash: calls.GetAccessHash(int64(chatId)),
		},
	)
	if err != nil {
		return "", err
	}
	call, ok := fullChannel.FullChat.GetCall()
	if !ok {
		return "", ErrNoCall
	}
	params = map[string]interface{}{
		"ufrag": payload["ufrag"],
		"pwd":   payload["pwd"],
		"fingerprints": []map[string]interface{}{{
			"hash":        payload["hash"],
			"setup":       payload["setup"],
			"fingerprint": payload["fingerprint"],
		}},
		"ssrc": payload["source"],
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	updates, err := calls.api.PhoneJoinGroupCall(
		calls.ctx,
		&tg.PhoneJoinGroupCallRequest{
			Call:  call,
			Muted: false,
			Params: tg.DataJSON{
				Data: string(data),
			},
			JoinAs: &tg.InputPeerSelf{},
		},
	)
	if err != nil {
		return "", err
	}
	if updates, ok := updates.(*tg.Updates); ok {
		for _, update := range updates.Updates {
			if update, ok := update.(*tg.UpdateGroupCallConnection); ok {
				return update.Params.Data, nil
			}
		}
	}
	return "", ErrUnexpectedType
}
