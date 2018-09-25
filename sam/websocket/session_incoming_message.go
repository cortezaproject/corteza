package websocket

import (
	"context"

	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func (s *Session) messageCreate(ctx context.Context, p *incoming.MessageCreate) error {
	var (
		msg = &types.Message{
			ChannelID: parseUInt64(p.ChannelID),
			Message:   p.Message,
		}
	)

	msg, err := s.svc.msg.With(ctx).Create(msg)
	if err != nil {
		return err
	}

	return s.sendToAllSubscribers(payloadFromMessage(msg), p.ChannelID)
}

func (s *Session) messageUpdate(ctx context.Context, p *incoming.MessageUpdate) error {
	var (
		msg = &types.Message{
			ID:      parseUInt64(p.ID),
			Message: p.Message,
		}
	)
	msg, err := s.svc.msg.With(ctx).Update(msg)
	if err != nil {
		return err
	}

	omsg := &outgoing.MessageUpdate{
		ID:        p.ID,
		Message:   msg.Message,
		UpdatedAt: *msg.UpdatedAt,
	}

	return s.sendToAllSubscribers(omsg, p.ID)
}

func (s *Session) messageDelete(ctx context.Context, p *incoming.MessageDelete) error {
	var (
		id = parseUInt64(p.ID)
	)

	if err := s.svc.msg.With(ctx).Delete(id); err != nil {
		return err
	}

	return s.sendToAllSubscribers(&outgoing.MessageDelete{ID: p.ID}, p.ChannelID)
}

func (s *Session) messageHistory(ctx context.Context, p *incoming.Messages) error {
	var (
		filter = &types.MessageFilter{
			ChannelID:      parseUInt64(p.ChannelID),
			FromMessageID:  parseUInt64(p.FromID),
			UntilMessageID: parseUInt64(p.UntilID),

			// Max no. of messages we will return
			Limit: 50,
		}
	)

	messages, err := s.svc.msg.With(ctx).Find(filter)
	if err != nil {
		return err
	}

	return s.sendReply(payloadFromMessages(messages))
}
