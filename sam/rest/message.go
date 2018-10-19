package rest

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
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
	return ctrl.wrap(ctrl.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) CreateReply(ctx context.Context, r *request.MessageCreateReply) (interface{}, error) {
	return ctrl.wrap(ctrl.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: r.ChannelID,
		ReplyTo:   r.MessageID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) GetReplies(ctx context.Context, r *request.MessageGetReplies) (interface{}, error) {
	return ctrl.wrapSet(ctrl.svc.msg.With(ctx).Find(&types.MessageFilter{
		ChannelID: r.ChannelID,
		RepliesTo: r.MessageID,
	}))
}

func (ctrl *Message) History(ctx context.Context, r *request.MessageHistory) (interface{}, error) {
	return ctrl.wrapSet(ctrl.svc.msg.With(ctx).Find(&types.MessageFilter{
		ChannelID: r.ChannelID,
		FirstID:   r.LastMessageID,
	}))
}

func (ctrl *Message) Edit(ctx context.Context, r *request.MessageEdit) (interface{}, error) {
	return ctrl.wrap(ctrl.svc.msg.With(ctx).Update(&types.Message{
		ID:        r.MessageID,
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) Delete(ctx context.Context, r *request.MessageDelete) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Delete(r.MessageID)
}

func (ctrl *Message) Search(ctx context.Context, r *request.MessageSearch) (interface{}, error) {
	return ctrl.wrapSet(ctrl.svc.msg.With(ctx).Find(&types.MessageFilter{
		ChannelID: r.ChannelID,
		Query:     r.Query,
	}))
}

func (ctrl *Message) Pin(ctx context.Context, r *request.MessagePin) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Pin(r.MessageID)
}

func (ctrl *Message) Unpin(ctx context.Context, r *request.MessageUnpin) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Unpin(r.MessageID)
}

func (ctrl *Message) Flag(ctx context.Context, r *request.MessageFlag) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Flag(r.MessageID)
}

func (ctrl *Message) Unflag(ctx context.Context, r *request.MessageUnflag) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Unflag(r.MessageID)
}

func (ctrl *Message) React(ctx context.Context, r *request.MessageReact) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).React(r.MessageID, r.Reaction)
}

func (ctrl *Message) Unreact(ctx context.Context, r *request.MessageUnreact) (interface{}, error) {
	return nil, ctrl.svc.msg.With(ctx).Unreact(r.MessageID, r.Reaction)
}
func (ctrl *Message) wrap(m *types.Message, err error) (*outgoing.Message, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Message(m), nil
	}
}

func (ctrl *Message) wrapSet(mm types.MessageSet, err error) (*outgoing.MessageSet, error) {
	if err != nil {
		return nil, err
	} else {
		return payload.Messages(mm), nil
	}
}
