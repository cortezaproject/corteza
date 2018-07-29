package websocket

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/outgoing"
)

func payloadFromMessage(msg *types.Message) *outgoing.Message {
	return &outgoing.Message{
		Message:   msg.Message,
		ID:        uint64toa(msg.ID),
		ChannelID: uint64toa(msg.ChannelID),
		Type:      msg.Type,
		UserID:    uint64toa(msg.UserID),
		ReplyTo:   uint64toa(msg.ReplyTo),

		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
}

func payloadFromMessages(msg []*types.Message) *outgoing.Messages {
	msgs := make([]*outgoing.Message, len(msg))
	for k, m := range msg {
		msgs[k] = payloadFromMessage(m)
	}
	retval := outgoing.Messages(msgs)
	return &retval
}

func payloadFromChannel(ch *types.Channel) *outgoing.Channel {
	return &outgoing.Channel{
		ID:            uint64toa(ch.ID),
		Name:          ch.Name,
		LastMessageID: uint64toa(ch.LastMessageID),
		Topic:         ch.Topic,
	}
}

func payloadFromChannels(channels []*types.Channel) *outgoing.Channels {
	cc := make([]*outgoing.Channel, len(channels))
	for k, c := range channels {
		cc[k] = payloadFromChannel(c)
	}
	retval := outgoing.Channels(cc)
	return &retval
}
