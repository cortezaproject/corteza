package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"
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

	switch r.Kind {
	case "":
		return nil, errors.New("missing value for activity kind")
	case "connected", "disconnected":
		return nil, errors.New("can not use reserved values for activity kind")
	}

	return true, ctrl.event.With(ctx).Activity(&types.Activity{
		UserID:    auth.GetIdentityFromContext(ctx).Identity(),
		ChannelID: r.ChannelID,
		MessageID: r.MessageID,
		Kind:      r.Kind,
	})
}
