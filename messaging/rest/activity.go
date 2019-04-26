package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/types"
)

var _ = errors.Wrap

type Activity struct {
	event service.EventService
}

func (Activity) New() *Activity {
	ctrl := &Activity{}
	ctrl.event = service.DefaultEvent
	return ctrl
}

// SendActivity Forwards channel activity to event service
func (ctrl *Activity) Send(ctx context.Context, r *request.ActivitySend) (interface{}, error) {
	if r.ChannelID == 0 && r.MessageID > 0 {
		return nil, errors.New("can not send activity on message without channel ID")
	}

	return true, ctrl.event.With(ctx).Activity(&types.Activity{
		UserID:    auth.GetIdentityFromContext(ctx).Identity(),
		ChannelID: r.ChannelID,
		MessageID: r.MessageID,
		Kind:      r.Kind,
	})
}
