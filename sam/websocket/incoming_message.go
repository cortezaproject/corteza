package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) messageCreate(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.MessageCreate
		msg     = &types.Message{
			ChannelID: parseUInt64(request.ChannelID),
			Message:   request.Message,
		}
	)

	msg, err := service.Message().Create(ctx, msg)
	if err != nil {
		return err
	}
	return s.sendMessageChannel(uint64toa(msg.ChannelID), payloadFromMessage(msg))
}

func (s *Session) messageUpdate(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.MessageUpdate
		msg     = &types.Message{
			ID:      parseUInt64(request.ID),
			Message: request.Message,
		}
	)
	msg, err := service.Message().Update(ctx, msg)
	if err != nil {
		return err
	}
	return s.sendMessageChannel(uint64toa(msg.ChannelID), &outgoing.MessageUpdate{request.ID, msg.Message})
}

func (s *Session) messageDelete(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.MessageDelete
		id      = parseUInt64(request.ID)
	)

	if err := service.Message().Delete(ctx, id); err != nil {
		return err
	}
	// @todo: delete broadcast to channel
	return s.sendMessageChannel("TODO", &outgoing.MessageDelete{request.ID})
}
