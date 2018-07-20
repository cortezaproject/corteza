package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Message struct {
		service messageService
	}

	messageService interface {
		Find(context.Context, *types.MessageFilter) ([]*types.Message, error)

		Create(context.Context, *types.Message) (*types.Message, error)
		Update(context.Context, *types.Message) (*types.Message, error)

		React(context.Context, uint64, string) error
		Unreact(context.Context, uint64, string) error

		Pin(context.Context, uint64) error
		Unpin(context.Context, uint64) error

		Flag(context.Context, uint64) error
		Unflag(context.Context, uint64) error

		Attach(context.Context) (*types.Attachment, error)
		Detach(context.Context, uint64) error

		deleter
	}
)

func (Message) New() *Message {
	return &Message{
		service: service.Message(),
	}
}

func (ctrl *Message) Create(ctx context.Context, r *server.MessageCreateRequest) (interface{}, error) {
	spew.Dump(r)
	return ctrl.service.Create(ctx, (&types.Message{}).
		SetChannelId(r.ChannelId).
		SetMessage(r.Message))
}

func (ctrl *Message) History(ctx context.Context, r *server.MessageHistoryRequest) (interface{}, error) {
	return ctrl.service.Find(ctx, &types.MessageFilter{
		ChannelId:     r.ChannelId,
		LastMessageId: r.LastMessageId,
	})
}

func (ctrl *Message) Edit(ctx context.Context, r *server.MessageEditRequest) (interface{}, error) {
	return ctrl.service.Update(ctx, (&types.Message{}).
		SetID(r.MessageId).
		SetChannelId(r.ChannelId).
		SetMessage(r.Message))
}

func (ctrl *Message) Delete(ctx context.Context, r *server.MessageDeleteRequest) (interface{}, error) {
	return nil, ctrl.service.Delete(ctx, r.MessageId)
}

func (ctrl *Message) Attach(ctx context.Context, r *server.MessageAttachRequest) (interface{}, error) {
	return ctrl.service.Attach(ctx)
}

func (ctrl *Message) Search(ctx context.Context, r *server.MessageSearchRequest) (interface{}, error) {
	return ctrl.service.Find(ctx, &types.MessageFilter{
		ChannelId: r.ChannelId,
		Query:     r.Query,
	})
}

func (ctrl *Message) Pin(ctx context.Context, r *server.MessagePinRequest) (interface{}, error) {
	return nil, ctrl.service.Pin(ctx, r.MessageId)
}

func (ctrl *Message) Unpin(ctx context.Context, r *server.MessageUnpinRequest) (interface{}, error) {
	return nil, ctrl.service.Unpin(ctx, r.MessageId)
}

func (ctrl *Message) Flag(ctx context.Context, r *server.MessageFlagRequest) (interface{}, error) {
	return nil, ctrl.service.Flag(ctx, r.MessageId)
}

func (ctrl *Message) Unflag(ctx context.Context, r *server.MessageUnflagRequest) (interface{}, error) {
	return nil, ctrl.service.Unflag(ctx, r.MessageId)
}

func (ctrl *Message) React(ctx context.Context, r *server.MessageReactRequest) (interface{}, error) {
	return nil, ctrl.service.React(ctx, r.MessageId, r.Reaction)
}

func (ctrl *Message) Unreact(ctx context.Context, r *server.MessageUnreactRequest) (interface{}, error) {
	return nil, ctrl.service.Unreact(ctx, r.MessageId, r.Reaction)
}
