package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

var _ = errors.Wrap

type (
	Message struct {
		svc struct {
			msg     service.MessageService
			command service.CommandService
		}
	}
)

func (Message) New() *Message {
	ctrl := &Message{}
	ctrl.svc.msg = service.DefaultMessage
	ctrl.svc.command = service.DefaultCommand
	return ctrl
}

func (ctrl *Message) Create(ctx context.Context, r *request.MessageCreate) (interface{}, error) {
	return ctrl.wrap(ctx)(ctrl.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) ReplyCreate(ctx context.Context, r *request.MessageReplyCreate) (interface{}, error) {
	return ctrl.wrap(ctx)(ctrl.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: r.ChannelID,
		ReplyTo:   r.MessageID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) Edit(ctx context.Context, r *request.MessageEdit) (interface{}, error) {
	return ctrl.wrap(ctx)(ctrl.svc.msg.With(ctx).Update(&types.Message{
		ID:        r.MessageID,
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl Message) ExecuteCommand(ctx context.Context, r *request.MessageExecuteCommand) (interface{}, error) {
	return ctrl.svc.command.With(ctx).Do(r.ChannelID, r.Command, r.Input)
}

func (ctrl *Message) Delete(ctx context.Context, r *request.MessageDelete) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).Delete(r.MessageID)
}

func (ctrl *Message) MarkAsRead(ctx context.Context, r *request.MessageMarkAsRead) (interface{}, error) {
	return ctrl.svc.msg.With(ctx).MarkAsRead(r.ChannelID, r.ThreadID, r.LastReadMessageID)
}

func (ctrl *Message) PinCreate(ctx context.Context, r *request.MessagePinCreate) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).Pin(r.MessageID)
}

func (ctrl *Message) PinRemove(ctx context.Context, r *request.MessagePinRemove) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).RemovePin(r.MessageID)
}

func (ctrl *Message) BookmarkCreate(ctx context.Context, r *request.MessageBookmarkCreate) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).Bookmark(r.MessageID)
}

func (ctrl *Message) BookmarkRemove(ctx context.Context, r *request.MessageBookmarkRemove) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).RemoveBookmark(r.MessageID)
}

func (ctrl *Message) ReactionCreate(ctx context.Context, r *request.MessageReactionCreate) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).React(r.MessageID, r.Reaction)
}

func (ctrl *Message) ReactionRemove(ctx context.Context, r *request.MessageReactionRemove) (interface{}, error) {
	return resputil.OK(), ctrl.svc.msg.With(ctx).RemoveReaction(r.MessageID, r.Reaction)
}

func (ctrl *Message) wrap(ctx context.Context) func(m *types.Message, err error) (*outgoing.Message, error) {
	return func(m *types.Message, err error) (*outgoing.Message, error) {
		if err != nil || m == nil {
			return nil, err
		} else {
			return payload.Message(ctx, m), nil
		}
	}
}
