package tgcalls

import (
	"encoding/json"

	"github.com/gotd/td/tg"
	"github.com/mitchellh/mapstructure"
)

type joinCallParams struct {
	chatId     int64
	accessHash int64
	isChat     bool
	payload    map[string]interface{}
}

func (calls *TGCalls) joinCall(params_ map[string]interface{}) (string, error) {
	var params joinCallParams
	err := mapstructure.Decode(params_, &params)
	if err != nil {
		return "", err
	}
	var fullChat tg.ChatFullClass
	if params.isChat {
		full, err := calls.api.MessagesGetFullChat(
			calls.ctx,
			params.chatId,
		)
		if err != nil {
			return "", err
		}
		fullChat = full.FullChat
	} else {
		full, err := calls.api.ChannelsGetFullChannel(
			calls.ctx,
			&tg.InputChannel{
				ChannelID:  params.chatId,
				AccessHash: params.accessHash,
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
	result := map[string]interface{}{
		"ufrag": params.payload["ufrag"],
		"pwd":   params.payload["pwd"],
		"fingerprints": []map[string]interface{}{{
			"hash":        params.payload["hash"],
			"setup":       params.payload["setup"],
			"fingerprint": params.payload["fingerprint"],
		}},
		"ssrc": params.payload["source"],
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
