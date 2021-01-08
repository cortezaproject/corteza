package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"

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
	mm, _, err := ctrl.svc.msg.Find(ctx, types.MessageFilter{
		ChannelID:      payload.ParseUint64s(r.ChannelID),
		AfterID:        r.AfterMessageID,
		BeforeID:       r.BeforeMessageID,
		FromID:         r.FromMessageID,
		ToID:           r.ToMessageID,
		ThreadID:       payload.ParseUint64s(r.ThreadID),
		UserID:         payload.ParseUint64s(r.UserID),
		Type:           r.Type,
		PinnedOnly:     r.PinnedOnly,
		BookmarkedOnly: r.BookmarkedOnly,
		Limit:          r.Limit,

		Query: r.Query,
	})

	return ctrl.wrapSet(ctx, mm, err)
}

func (ctrl *Search) Threads(ctx context.Context, r *request.SearchThreads) (interface{}, error) {
	mm, _, err := ctrl.svc.msg.FindThreads(ctx, types.MessageFilter{
		ChannelID: payload.ParseUint64s(r.ChannelID),
		Limit:     r.Limit,

		Query: r.Query,
	})

	return ctrl.wrapSet(ctx, mm, err)

}

func (ctrl *Search) wrapSet(ctx context.Context, mm types.MessageSet, err error) (*outgoing.MessageSet, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Messages(ctx, mm), nil
	}
}
