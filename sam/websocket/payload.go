package websocket

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"strconv"
)

func payloadFromMessage(msg *types.Message) *outgoing.Message {
	return &outgoing.Message{
		Message:   msg.Message,
		ID:        strconv.FormatUint(msg.ID, 10),
		ChannelID: strconv.FormatUint(msg.ChannelID, 10),
		Type:      msg.Type,
		UserID:    strconv.FormatUint(msg.UserID, 10),
		ReplyTo:   strconv.FormatUint(msg.ReplyTo, 10),
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
