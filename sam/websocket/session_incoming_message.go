package websocket

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/sam/types"
)

func (s *Session) messageCreate(ctx context.Context, p *incoming.MessageCreate) error {
	var (
		msg = &types.Message{
			ChannelID: payload.ParseUInt64(p.ChannelID),
			Message:   p.Message,
		}
	)

	msg, err := s.svc.msg.With(ctx).Create(msg)
	if err != nil {
		return err
	}

	return s.sendToAllSubscribers(payload.Message(msg), p.ChannelID)
}

func (s *Session) messageUpdate(ctx context.Context, p *incoming.MessageUpdate) error {
	var (
		msg = &types.Message{
			ID:      payload.ParseUInt64(p.ID),
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
		id = payload.ParseUInt64(p.ID)
	)

	if err := s.svc.msg.With(ctx).Delete(id); err != nil {
		return err
	}

	return s.sendToAllSubscribers(&outgoing.MessageDelete{ID: p.ID}, p.ChannelID)
}

func (s *Session) messageHistory(ctx context.Context, p *incoming.Messages) error {
	var (
		filter = &types.MessageFilter{
			ChannelID:      payload.ParseUInt64(p.ChannelID),
			FromMessageID:  payload.ParseUInt64(p.FromID),
			UntilMessageID: payload.ParseUInt64(p.UntilID),

			// Max no. of messages we will return
			Limit: 50,
		}
	)

	messages, err := s.svc.msg.With(ctx).Find(filter)
	if err != nil {
		return err
	}

	return s.sendReply(payload.Messages(messages))
}
