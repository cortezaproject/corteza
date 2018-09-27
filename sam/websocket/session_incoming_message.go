package websocket

import (
	"context"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/sam/types"
)

func (s *Session) messageCreate(ctx context.Context, p *incoming.MessageCreate) error {
	_, err := s.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: payload.ParseUInt64(p.ChannelID),
		Message:   p.Message,
	})

	return err
}

func (s *Session) messageUpdate(ctx context.Context, p *incoming.MessageUpdate) error {
	_, err := s.svc.msg.With(ctx).Update(&types.Message{
		ID:      payload.ParseUInt64(p.ID),
		Message: p.Message,
	})

	return err
}

func (s *Session) messageDelete(ctx context.Context, p *incoming.MessageDelete) error {
	return s.svc.msg.With(ctx).Delete(payload.ParseUInt64(p.ID))
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
