package tgcalls

import (
	"encoding/json"

	"github.com/gotd/td/tg"
)

func (calls *TGCalls) joinCall(params map[string]interface{}) (string, error) {
	chatIdFloat, ok := params["chatId"].(float64)
	if !ok {
		return "", ErrUnexpectedType
	}
	chatId := int64(chatIdFloat)
	isChat, ok := params["isChat"].(bool)
	if !ok {
		return "", ErrUnexpectedType
	}
	payload, ok := params["payload"].(map[string]interface{})
	if !ok {
		return "", ErrUnexpectedType
	}
	var fullChat tg.ChatFullClass
	if isChat {
		full, err := calls.api.MessagesGetFullChat(
			calls.ctx,
			chatId,
		)
		if err != nil {
			return "", err
		}
		fullChat = full.FullChat
	} else {
		accessHashFloat, ok := params["accessHash"].(float64)
		if !ok {
			return "", ErrUnexpectedType
		}
		accessHash := int64(accessHashFloat)
		full, err := calls.api.ChannelsGetFullChannel(
			calls.ctx,
			&tg.InputChannel{
				ChannelID:  chatId,
				AccessHash: accessHash,
			},
		)
		if err != nil {
			return "", err
		}
		fullChat = full.FullChat
	}
	call, ok := fullChat.GetCall()
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
