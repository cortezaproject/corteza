package websocket

import (
	"context"
	"encoding/json"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/pkg/errors"
	"strconv"
)

func (s *Session) dispatch(raw []byte) (err error) {
	var payload = &incoming.Message{}
	if err = json.Unmarshal(raw, payload); err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	ctx := s.Context()
	if p := payload.MessageCreate; p != nil {
		return s.dispatchMessageCreate(ctx, p)
	}

	if p := payload.ChannelOpen; p != nil {
		return s.dispatchChannelOpen(ctx, p)
	}

	return nil
}

func (s *Session) dispatchMessageCreate(ctx context.Context, payload *incoming.MessageCreate) (err error) {
	var (
		msg = &types.Message{Message: payload.Message}
	)

	if msg.ChannelID, err = strconv.ParseUint(payload.ChannelID, 10, 64); err != nil {
		return
	}

	if msg, err = service.Message().Create(ctx, msg); err != nil {
		return
	} else {
		// @todo move this to outgoing.FromMessage(*types.Message) *outgoing.WsMessage
		store.MessageFanout(payloadFromMessage(msg))
	}

	return
}

func (s *Session) dispatchChannelOpen(ctx context.Context, payload *incoming.ChannelOpen) (err error) {
	var filter = &types.MessageFilter{}

	if filter.ChannelID, err = strconv.ParseUint(payload.ChannelID, 10, 64); err != nil {
		return
	}

	if filter.FromMessageID, err = strconv.ParseUint(payload.Since, 10, 64); err != nil {
		return
	}

	if messages, err := service.Message().Find(ctx, filter); err != nil {
		return err
	} else {
		for _, msg := range messages {
			s.send <- payloadFromMessage(msg)
		}
	}

	return
}
