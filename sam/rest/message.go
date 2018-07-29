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
		Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error)

		Create(ctx context.Context, messages *types.Message) (*types.Message, error)
		Update(ctx context.Context, messages *types.Message) (*types.Message, error)

		React(ctx context.Context, messageID uint64, reaction string) error
		Unreact(ctx context.Context, messageID uint64, reaction string) error

		Pin(ctx context.Context, messageID uint64) error
		Unpin(ctx context.Context, messageID uint64) error

		Flag(ctx context.Context, messageID uint64) error
		Unflag(ctx context.Context, messageID uint64) error

		Attach(ctx context.Context) (*types.Attachment, error)
		Detach(ctx context.Context, messageID uint64) error

		deleter
	}
)

func (Message) New(messageSvc messageService) *Message {
	var ctrl = &Message{}
	ctrl.svc = messageSvc
	return ctrl
}

func (ctrl *Message) Create(ctx context.Context, r *server.MessageCreateRequest) (interface{}, error) {
	return ctrl.svc.Create(ctx, &types.Message{
		ChannelID: r.ChannelID,
		Message:   r.Message,
	})
}

func (ctrl *Message) History(ctx context.Context, r *server.MessageHistoryRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.MessageFilter{
		ChannelID:     r.ChannelID,
		FromMessageID: r.LastMessageID,
	})
}

func (ctrl *Message) Edit(ctx context.Context, r *server.MessageEditRequest) (interface{}, error) {
	return ctrl.svc.Update(ctx, &types.Message{
		ID:        r.MessageID,
		ChannelID: r.ChannelID,
		Message:   r.Message,
	})
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
