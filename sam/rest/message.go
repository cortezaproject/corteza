package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Message struct {
		svc struct {
			msg service.MessageService
		}
	}
)

func (Message) New() *Message {
	ctrl := &Message{}
	ctrl.svc.msg = service.DefaultMessage
	return ctrl
}

func (ctrl *Message) Create(ctx context.Context, r *request.MessageCreate) (interface{}, error) {
	return ctrl.svc.msg.Create(ctx, &types.Message{
		ChannelID: r.ChannelID,
		Message:   r.Message,
	})
}

func (ctrl *Message) History(ctx context.Context, r *request.MessageHistory) (interface{}, error) {
	return ctrl.svc.msg.Find(ctx, &types.MessageFilter{
		ChannelID:     r.ChannelID,
		FromMessageID: r.LastMessageID,
	})
}

func (ctrl *Message) Edit(ctx context.Context, r *request.MessageEdit) (interface{}, error) {
	return ctrl.svc.msg.Update(ctx, &types.Message{
		ID:        r.MessageID,
		ChannelID: r.ChannelID,
		Message:   r.Message,
	})
}

func (ctrl *Message) Delete(ctx context.Context, r *request.MessageDelete) (interface{}, error) {
	return nil, ctrl.svc.msg.Delete(ctx, r.MessageID)
}

func (ctrl *Message) Search(ctx context.Context, r *request.MessageSearch) (interface{}, error) {
	return ctrl.svc.msg.Find(ctx, &types.MessageFilter{
		ChannelID: r.ChannelID,
		Query:     r.Query,
	})
}

func (ctrl *Message) Pin(ctx context.Context, r *request.MessagePin) (interface{}, error) {
	return nil, ctrl.svc.msg.Pin(ctx, r.MessageID)
}

func (ctrl *Message) Unpin(ctx context.Context, r *request.MessageUnpin) (interface{}, error) {
	return nil, ctrl.svc.msg.Unpin(ctx, r.MessageID)
}

func (ctrl *Message) Flag(ctx context.Context, r *request.MessageFlag) (interface{}, error) {
	return nil, ctrl.svc.msg.Flag(ctx, r.MessageID)
}

func (ctrl *Message) Unflag(ctx context.Context, r *request.MessageUnflag) (interface{}, error) {
	return nil, ctrl.svc.msg.Unflag(ctx, r.MessageID)
}

func (ctrl *Message) React(ctx context.Context, r *request.MessageReact) (interface{}, error) {
	return nil, ctrl.svc.msg.React(ctx, r.MessageID, r.Reaction)
}

func (ctrl *Message) Unreact(ctx context.Context, r *request.MessageUnreact) (interface{}, error) {
	return nil, ctrl.svc.msg.Unreact(ctx, r.MessageID, r.Reaction)
}
