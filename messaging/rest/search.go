package rest

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/service"
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
		Query:     r.Query,
		ChannelID: r.InChannel,
		UserID:    r.FromUser,
		FirstID:   r.FirstID,
		LastID:    r.LastID,
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
