package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"

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
		ChannelID:      payload.ParseUInt64s(r.ChannelID),
		AfterID:        r.AfterMessageID,
		BeforeID:       r.BeforeMessageID,
		FromID:         r.FromMessageID,
		ToID:           r.ToMessageID,
		ThreadID:       payload.ParseUInt64s(r.ThreadID),
		UserID:         payload.ParseUInt64s(r.UserID),
		Type:           r.Type,
		PinnedOnly:     r.PinnedOnly,
		BookmarkedOnly: r.BookmarkedOnly,
		Limit:          r.Limit,

		Query: r.Query,
	}))
}

func (ctrl *Search) Threads(ctx context.Context, r *request.SearchThreads) (interface{}, error) {
	return ctrl.wrapSet(ctx)(ctrl.svc.msg.With(ctx).FindThreads(&types.MessageFilter{
		ChannelID: payload.ParseUInt64s(r.ChannelID),
		Limit:     r.Limit,

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
