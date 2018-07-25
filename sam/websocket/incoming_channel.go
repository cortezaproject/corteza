package websocket

import (
	"context"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
)

func (s *Session) channelJoin(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.ChannelJoin
	)
	// @todo: check access to channel
	s.subs.Add(request.ChannelID, &Subscription{})
	return nil
}

func (s *Session) channelPart(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.ChannelJoin
	)
	// @todo: check access to channel
	s.subs.Delete(request.ChannelID)
	return nil
}

func (s *Session) channelPartAll(ctx context.Context, payload *incoming.Payload) error {
	if payload.ChannelPartAll.Leave {
		s.subs.DeleteAll()
	}
	return nil
}

func (s *Session) channelOpen(ctx context.Context, payload *incoming.Payload) error {
	var (
		request = payload.ChannelOpen
		filter  = &types.MessageFilter{
			ChannelID:     parseUInt64(request.ChannelID),
			FromMessageID: parseUInt64(request.Since),
		}
	)

	messages, err := service.Message().Find(ctx, filter)
	if err != nil {
		return err
	}
	return s.sendMessage(payloadFromMessages(messages))
}
