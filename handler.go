package tgcalls

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func (calls *TGCalls) Handler() telegram.UpdateHandlerFunc {
	return func(ctx context.Context, u tg.UpdatesClass) error {
		var err error
		switch u := u.(type) {
		case *tg.Updates:
			for _, u := range u.Updates {
				switch u := u.(type) {
				case *tg.UpdateGroupCall:
					{
						if _, ok := u.Call.(*tg.GroupCallDiscarded); ok {
							_, err = calls.Stop(u.ChatID)
						}
					}
				}
			}
		}
		return err
	}
}
