package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) messageCreate(ctx context.Context, p incoming.MessageCreate) error {
	var (
		msg = &types.Message{
			ChannelID: parseUInt64(p.ChannelID),
			Message:   p.Message,
		}
	)

	msg, err := service.Message().Create(ctx, msg)
	if err != nil {
		return err
	}

	return s.broadcast(payloadFromMessage(msg), &p.ChannelID)
}

func (s *Session) messageUpdate(ctx context.Context, p incoming.MessageUpdate) error {
	var (
		msg = &types.Message{
			ID:      parseUInt64(p.ID),
			Message: p.Message,
		}
	)
	msg, err := service.Message().Update(ctx, msg)
	if err != nil {
		return err
	}

	return s.broadcast(&outgoing.MessageUpdate{ID: p.ID, Message: msg.Message}, &p.ID)
}

func (s *Session) messageDelete(ctx context.Context, p incoming.MessageDelete) error {
	var (
		id = parseUInt64(p.ID)
	)

	if err := service.Message().Delete(ctx, id); err != nil {
		return err
	}

	return s.broadcast(&outgoing.MessageDelete{ID: p.ID}, &p.ChannelID)
}
