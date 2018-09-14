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
		Type:      string(msg.Type),
		UserID:    uint64toa(msg.UserID),
		ReplyTo:   uint64toa(msg.ReplyTo),

		Attachment: payloadFromAttachment(msg.Attachment),

		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}
}

func payloadFromMessages(msg types.MessageSet) *outgoing.Messages {
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
		Type:          string(ch.Type),
		Members:       payloadFromUsers(ch.Members),
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

func payloadFromUser(user *types.User) *outgoing.User {
	return &outgoing.User{
		ID:       uint64toa(user.ID),
		Username: user.Username,
	}
}

func payloadFromUsers(users []*types.User) *outgoing.Users {
	uu := make([]*outgoing.User, len(users))
	for k, u := range users {
		uu[k] = payloadFromUser(u)
		uu[k].Connections = 0

		// @todo this is current instance only, need to sync this across all instances
		store.Walk(func(session *Session) {
			if session.user.ID == u.ID {
				uu[k].Connections++
			}
		})
	}

	retval := outgoing.Users(uu)

	return &retval
}

func payloadFromAttachment(in *types.Attachment) *outgoing.Attachment {
	if in == nil {
		return nil
	}

	return &outgoing.Attachment{
		ID:         uint64toa(in.ID),
		UserID:     uint64toa(in.UserID),
		Url:        in.Url,
		PreviewUrl: in.PreviewUrl,
		Size:       in.Size,
		Mimetype:   in.Mimetype,
		Name:       in.Name,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}
}
