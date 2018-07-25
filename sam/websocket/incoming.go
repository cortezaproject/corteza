package websocket

import (
	"context"
	"encoding/json"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"github.com/pkg/errors"
	"strconv"
)

func (s *Session) dispatch(raw []byte) (err error) {
	var payload = &incoming.Message{}
	if err = json.Unmarshal(raw, payload); err != nil {
		return errors.Wrap(err, "Session.incoming: payload malformed")
	}

	if p := payload.MessageCreate; p != nil {
		return s.dispatchMessageCreate(p)
	}

	if p := payload.ChannelOpen; p != nil {
		return s.dispatchChannelOpen(p)
	}

	return nil
}

func (s *Session) dispatchMessageCreate(payload *incoming.MessageCreate) (err error) {
	var (
		msg = &types.Message{Message: payload.Message}
	)

	if msg.ChannelId, err = strconv.ParseUint(payload.ChannelId, 10, 64); err != nil {
		return
	}

	if msg, err = service.Message().Create(context.TODO(), msg); err != nil {
		return
	} else {
		// @todo move this to outgoing.FromMessage(*types.Message) *outgoing.WsMessage
		store.MessageFanout(outgoing.FromMessage(msg))
	}

	return
}

func (s *Session) dispatchChannelOpen(payload *incoming.ChannelOpen) (err error) {
	var filter = &types.MessageFilter{}

	if filter.ChannelId, err = strconv.ParseUint(payload.ChannelId, 10, 64); err != nil {
		return
	}

	if filter.FromMessageId, err = strconv.ParseUint(payload.Since, 10, 64); err != nil {
		return
	}

	if messages, err := service.Message().Find(context.TODO(), filter); err != nil {
		return err
	} else {
		for _, msg := range messages {
			s.send <- outgoing.FromMessage(msg)
		}
	}

	return
}
