package tgcalls

import (
	"encoding/json"
	"strconv"

	"github.com/gotd/td/tg"
)

func (calls *TGCalls) joinCall(params map[string]interface{}) (string, error) {
	idString, ok := params["id"].(string)
	if !ok {
		return "", ErrUnexpectedType
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return "", err
	}
	joinCallParams, ok := params["joinCallParams"].(map[string]interface{})
	if !ok {
		return "", ErrUnexpectedType
	}
	isChannel, ok := joinCallParams["isChannel"].(bool)
	if !ok {
		return "", ErrUnexpectedType
	}
	payload, ok := params["payload"].(map[string]interface{})
	if !ok {
		return "", ErrUnexpectedType
	}
	var fullChat tg.ChatFullClass
	if isChannel {
		accessHashString, ok := joinCallParams["accessHash"].(string)
		if !ok {
			return "", ErrUnexpectedType
		}
		accessHash, err := strconv.ParseInt(accessHashString, 10, 64)
		if err != nil {
			return "", err
		}
		full, err := calls.api.ChannelsGetFullChannel(
			calls.ctx,
			&tg.InputChannel{
				ChannelID:  id,
				AccessHash: accessHash,
			},
		)
		if err != nil {
			return "", err
		}
		fullChat = full.FullChat
	} else {
		full, err := calls.api.MessagesGetFullChat(
			calls.ctx,
			id,
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
	result := map[string]interface{}{
		"ufrag": payload["ufrag"],
		"pwd":   payload["pwd"],
		"fingerprints": []map[string]interface{}{{
			"hash":        payload["hash"],
			"setup":       payload["setup"],
			"fingerprint": payload["fingerprint"],
		}},
		"ssrc": payload["source"],
	}
	data, err := json.Marshal(result)
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
