package websocket

import (
	"context"
	"encoding/json"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/pkg/errors"
)

func (s *Session) dispatch(raw []byte) (err error) {
	var p = &incoming.Payload{}
	if err = json.Unmarshal(raw, p); err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	ctx := s.Context()
	switch {
	case p.MessageCreate != nil:
		return s.messageCreate(ctx, p)
	case p.ChannelOpen != nil:
		return s.channelOpen(ctx, p)
	}

	return nil
}

func (s *Session) messageCreate(ctx context.Context, payload *incoming.Payload) (err error) {
	var (
		request = payload.MessageCreate
		msg     = &types.Message{Message: request.Message}
	)

	msg.ChannelID = parseUInt64(request.ChannelID)
	if msg, err = service.Message().Create(ctx, msg); err != nil {
		return
	} else {
		// @todo move this to outgoing.FromMessage(*types.Message) *outgoing.WsMessage
		store.MessageFanout(payloadFromMessage(msg))
	}

	return
}

func (s *Session) channelOpen(ctx context.Context, payload *incoming.Payload) (err error) {
	var (
		request = payload.ChannelOpen
		filter  = &types.MessageFilter{}
	)

	filter.ChannelID = parseUInt64(request.ChannelID)
	filter.FromMessageID = parseUInt64(request.Since)

	if messages, err := service.Message().Find(ctx, filter); err != nil {
		return err
	} else {
		for _, msg := range messages {
			s.send <- payloadFromMessage(msg)
		}
	}

	return
}
