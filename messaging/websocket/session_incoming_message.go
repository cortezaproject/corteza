package websocket

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/incoming"
)

func (s *Session) messageCreate(ctx context.Context, p *incoming.MessageCreate) error {
	_, err := s.svc.msg.With(ctx).Create(&types.Message{
		ChannelID: payload.ParseUint64(p.ChannelID),
		ReplyTo:   p.ReplyTo,
		Message:   p.Message,
	})

	return err
}

func (s *Session) messageUpdate(ctx context.Context, p *incoming.MessageUpdate) error {
	_, err := s.svc.msg.With(ctx).Update(&types.Message{
		ID:      payload.ParseUint64(p.ID),
		Message: p.Message,
	})

	return err
}

func (s *Session) messageDelete(ctx context.Context, p *incoming.MessageDelete) error {
	return s.svc.msg.With(ctx).Delete(payload.ParseUint64(p.ID))
}
