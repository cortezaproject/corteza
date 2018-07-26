package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Message struct {
		svc messageService
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

func (Message) New(messageSvc messageService) *Message {
	var ctrl = &Message{}
	ctrl.svc = messageSvc
	return ctrl
}

func (ctrl *Message) Create(ctx context.Context, r *server.MessageCreateRequest) (interface{}, error) {
	return ctrl.svc.Create(ctx, (&types.Message{}).
		SetChannelID(r.ChannelID).
		SetMessage(r.Message))
}

func (ctrl *Message) History(ctx context.Context, r *server.MessageHistoryRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.MessageFilter{
		ChannelID:     r.ChannelID,
		FromMessageID: r.LastMessageID,
	})
}

func (ctrl *Message) Edit(ctx context.Context, r *server.MessageEditRequest) (interface{}, error) {
	return ctrl.svc.Update(ctx, (&types.Message{}).
		SetID(r.MessageID).
		SetChannelID(r.ChannelID).
		SetMessage(r.Message))
}

func (ctrl *Message) Delete(ctx context.Context, r *server.MessageDeleteRequest) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.MessageID)
}

func (ctrl *Message) Attach(ctx context.Context, r *server.MessageAttachRequest) (interface{}, error) {
	return ctrl.svc.Attach(ctx)
}

func (ctrl *Message) Search(ctx context.Context, r *server.MessageSearchRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.MessageFilter{
		ChannelID: r.ChannelID,
		Query:     r.Query,
	})
}

func (ctrl *Message) Pin(ctx context.Context, r *server.MessagePinRequest) (interface{}, error) {
	return nil, ctrl.svc.Pin(ctx, r.MessageID)
}

func (ctrl *Message) Unpin(ctx context.Context, r *server.MessageUnpinRequest) (interface{}, error) {
	return nil, ctrl.svc.Unpin(ctx, r.MessageID)
}

func (ctrl *Message) Flag(ctx context.Context, r *server.MessageFlagRequest) (interface{}, error) {
	return nil, ctrl.svc.Flag(ctx, r.MessageID)
}

func (ctrl *Message) Unflag(ctx context.Context, r *server.MessageUnflagRequest) (interface{}, error) {
	return nil, ctrl.svc.Unflag(ctx, r.MessageID)
}

func (ctrl *Message) React(ctx context.Context, r *server.MessageReactRequest) (interface{}, error) {
	return nil, ctrl.svc.React(ctx, r.MessageID, r.Reaction)
}

func (ctrl *Message) Unreact(ctx context.Context, r *server.MessageUnreactRequest) (interface{}, error) {
	return nil, ctrl.svc.Unreact(ctx, r.MessageID, r.Reaction)
}
