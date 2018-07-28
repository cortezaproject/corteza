package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
)

func (s *Session) channelJoin(ctx context.Context, p incoming.ChannelJoin) error {
	var ()
	// @todo: check access to channel
	s.subs.Add(p.ChannelID, &Subscription{})
	return nil
}

func (s *Session) channelPart(ctx context.Context, p incoming.ChannelPart) error {
	var ()
	// @todo: check access to channel
	s.subs.Delete(p.ChannelID)
	return nil
}

func (s *Session) channelPartAll(ctx context.Context, p incoming.ChannelPartAll) error {
	if p.Leave {
		s.subs.DeleteAll()
	}
	return nil
}

func (s *Session) channelOpen(ctx context.Context, p incoming.ChannelOpen) error {
	var (
		filter = &types.MessageFilter{
			ChannelID:     parseUInt64(p.ChannelID),
			FromMessageID: parseUInt64(p.Since),
		}
	)

	messages, err := service.Message().Find(ctx, filter)
	if err != nil {
		return err
	}
	return s.sendMessage(payloadFromMessages(messages))
}
