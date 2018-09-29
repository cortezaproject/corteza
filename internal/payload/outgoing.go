package payload

import (
	"fmt"
	"net/url"

	auth "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/internal/payload/outgoing"
	sam "github.com/crusttech/crust/sam/types"
)

const (
	attachmentURL        = "/attachment/%d/original/%s"
	attachmentPreviewURL = "/attachment/%d/preview.%s"
)

func Message(msg *sam.Message) *outgoing.Message {
	return &outgoing.Message{
		ID:        Uint64toa(msg.ID),
		ChannelID: Uint64toa(msg.ChannelID),
		Message:   msg.Message,
		Type:      string(msg.Type),
		ReplyTo:   Uint64toa(msg.ReplyTo),

		User:       User(msg.User),
		Attachment: Attachment(msg.Attachment),

		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		DeletedAt: msg.DeletedAt,
	}
}

func Messages(msg sam.MessageSet) *outgoing.MessageSet {
	msgs := make([]*outgoing.Message, len(msg))
	for k, m := range msg {
		msgs[k] = Message(m)
	}
	retval := outgoing.MessageSet(msgs)
	return &retval
}

func Channel(ch *sam.Channel) *outgoing.Channel {
	return &outgoing.Channel{
		ID:            Uint64toa(ch.ID),
		Name:          ch.Name,
		LastMessageID: Uint64toa(ch.LastMessageID),
		Topic:         ch.Topic,
		Type:          string(ch.Type),
		Members:       Uint64stoa(ch.Members),

		CreatedAt:  ch.CreatedAt,
		UpdatedAt:  ch.UpdatedAt,
		ArchivedAt: ch.ArchivedAt,
		DeletedAt:  ch.DeletedAt,
	}
}

func Channels(channels sam.ChannelSet) *outgoing.ChannelSet {
	cc := make([]*outgoing.Channel, len(channels))
	for k, c := range channels {
		cc[k] = Channel(c)
	}
	retval := outgoing.ChannelSet(cc)
	return &retval
}

func User(user *auth.User) *outgoing.User {
	if user == nil {
		return nil
	}

	return &outgoing.User{
		ID:       Uint64toa(user.ID),
		Username: user.Username,
	}
}

func Users(users []*auth.User) *outgoing.UserSet {
	uu := make([]*outgoing.User, len(users))
	for k, u := range users {
		uu[k] = User(u)
		uu[k].Connections = 0
	}

	retval := outgoing.UserSet(uu)

	return &retval
}

func Attachment(in *sam.Attachment) *outgoing.Attachment {
	if in == nil {
		return nil
	}

	var preview string

	if in.Meta.Preview != nil {
		preview = fmt.Sprintf(attachmentPreviewURL, in.ID, in.Meta.Preview.Extension)
	}

	return &outgoing.Attachment{
		ID:         Uint64toa(in.ID),
		UserID:     Uint64toa(in.UserID),
		Url:        fmt.Sprintf(attachmentURL, in.ID, url.PathEscape(in.Name)),
		PreviewUrl: preview,
		Meta:       in.Meta,
		Name:       in.Name,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}
}
