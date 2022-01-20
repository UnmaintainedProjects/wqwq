package tgcalls

import (
	"encoding/json"

	"github.com/gotd/td/tg"
)

func (calls *TGCalls) joinCall(params map[string]interface{}) (string, error) {
	payload, ok := params["payload"].(map[string]interface{})
	if !ok {
		return "", ErrUnexpectedType
	}
	inviteHash := ""
	var joinAs tg.InputPeerClass = &tg.InputPeerSelf{}
	if calls.opts != nil {
		inviteHash = calls.opts.InviteHash
		if calls.opts.JoinAs != nil {
			joinAs = calls.opts.JoinAs
		}
	}
	fullChannel, err := calls.api.ChannelsGetFullChannel(calls.ctx, calls.chat)
	if err != nil {
		return "", err
	}
	call, ok := fullChannel.FullChat.GetCall()
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
		"ssrc-groups": []interface{}{
			map[string]interface{}{
				"semantics": "FID",
				"groups":    payload["sourceGroups"],
			},
		},
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	updates, err := calls.api.PhoneJoinGroupCall(
		calls.ctx,
		&tg.PhoneJoinGroupCallRequest{
			Call:         call,
			Muted:        false,
			VideoStopped: false,
			JoinAs:       joinAs,
			InviteHash:   inviteHash,
			Params: tg.DataJSON{
				Data: string(data),
			},
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
