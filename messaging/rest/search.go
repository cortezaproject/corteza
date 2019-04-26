package rest

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Search struct {
	svc struct {
		msg service.MessageService
	}
}

func (Search) New() *Search {
	ctrl := &Search{}
	ctrl.svc.msg = service.DefaultMessage
	return ctrl

}

func (ctrl *Search) Messages(ctx context.Context, r *request.SearchMessages) (interface{}, error) {
	return ctrl.wrapSet(ctx)(ctrl.svc.msg.With(ctx).Find(&types.MessageFilter{
		ChannelID:      r.ChannelID,
		AfterID:        r.AfterMessageID,
		BeforeID:       r.BeforeMessageID,
		FromID:         r.FromMessageID,
		ToID:           r.ToMessageID,
		ThreadID:       r.ThreadID,
		UserID:         r.UserID,
		Type:           r.Type,
		PinnedOnly:     r.PinnedOnly,
		BookmarkedOnly: r.BookmarkedOnly,
		Limit:          r.Limit,

		Query: r.Query,
	}))
}

func (ctrl *Search) wrapSet(ctx context.Context) func(mm types.MessageSet, err error) (*outgoing.MessageSet, error) {
	return func(mm types.MessageSet, err error) (*outgoing.MessageSet, error) {
		if err != nil {
			return nil, err
		} else {
			return payload.Messages(ctx, mm), nil
		}
	}
}
