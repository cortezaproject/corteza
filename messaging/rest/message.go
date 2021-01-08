package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/pkg/errors"
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
	return ctrl.wrap(ctx)(ctrl.svc.msg.Create(ctx, &types.Message{
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) ReplyCreate(ctx context.Context, r *request.MessageReplyCreate) (interface{}, error) {
	return ctrl.wrap(ctx)(ctrl.svc.msg.Create(ctx, &types.Message{
		ChannelID: r.ChannelID,
		ReplyTo:   r.MessageID,
		Message:   r.Message,
	}))
}

func (ctrl *Message) Edit(ctx context.Context, r *request.MessageEdit) (interface{}, error) {
	return ctrl.wrap(ctx)(ctrl.svc.msg.Update(ctx, &types.Message{
		ID:        r.MessageID,
		ChannelID: r.ChannelID,
		Message:   r.Message,
	}))
}

func (ctrl Message) ExecuteCommand(ctx context.Context, r *request.MessageExecuteCommand) (interface{}, error) {
	return ctrl.svc.command.Do(ctx, r.ChannelID, r.Command, r.Input)
}

func (ctrl *Message) Delete(ctx context.Context, r *request.MessageDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.Delete(ctx, r.MessageID)
}

func (ctrl *Message) MarkAsRead(ctx context.Context, r *request.MessageMarkAsRead) (interface{}, error) {
	var messageID, count, tcount, err = ctrl.svc.msg.MarkAsRead(ctx, r.ChannelID, r.ThreadID, r.LastReadMessageID)

	return outgoing.Unread{
		LastMessageID: messageID,
		Count:         count,
		ThreadCount:   tcount,
	}, err
}

func (ctrl *Message) PinCreate(ctx context.Context, r *request.MessagePinCreate) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.Pin(ctx, r.MessageID)
}

func (ctrl *Message) PinRemove(ctx context.Context, r *request.MessagePinRemove) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.RemovePin(ctx, r.MessageID)
}

func (ctrl *Message) BookmarkCreate(ctx context.Context, r *request.MessageBookmarkCreate) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.Bookmark(ctx, r.MessageID)
}

func (ctrl *Message) BookmarkRemove(ctx context.Context, r *request.MessageBookmarkRemove) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.RemoveBookmark(ctx, r.MessageID)
}

func (ctrl *Message) ReactionCreate(ctx context.Context, r *request.MessageReactionCreate) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.React(ctx, r.MessageID, r.Reaction)
}

func (ctrl *Message) ReactionRemove(ctx context.Context, r *request.MessageReactionRemove) (interface{}, error) {
	return api.OK(), ctrl.svc.msg.RemoveReaction(ctx, r.MessageID, r.Reaction)
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
