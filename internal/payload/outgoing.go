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
		ID:        msg.ID,
		ChannelID: Uint64toa(msg.ChannelID),
		Message:   msg.Message,
		Type:      string(msg.Type),
		ReplyTo:   msg.ReplyTo,
		Replies:   msg.Replies,

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
		View:          ChannelView(ch.View),

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

func ChannelMember(m *sam.ChannelMember) *outgoing.ChannelMember {
	return &outgoing.ChannelMember{
		User:      User(m.User),
		Type:      string(m.Type),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ChannelMembers(members sam.ChannelMemberSet) *outgoing.ChannelMemberSet {
	mm := make([]*outgoing.ChannelMember, len(members))
	for k, c := range members {
		mm[k] = ChannelMember(c)
	}
	retval := outgoing.ChannelMemberSet(mm)
	return &retval
}

func ChannelView(v *sam.ChannelView) *outgoing.ChannelView {
	if v == nil {
		return nil
	}

	return &outgoing.ChannelView{
		LastMessageID:    Uint64toa(v.LastMessageID),
		NewMessagesCount: v.NewMessagesCount,
	}
}

func ChannelJoin(channelID, userID uint64) *outgoing.ChannelJoin {
	return &outgoing.ChannelJoin{
		ID:     Uint64toa(channelID),
		UserID: Uint64toa(userID),
	}
}

func ChannelPart(channelID, userID uint64) *outgoing.ChannelPart {
	return &outgoing.ChannelPart{
		ID:     Uint64toa(channelID),
		UserID: Uint64toa(userID),
	}
}

func User(user *auth.User) *outgoing.User {
	if user == nil {
		return nil
	}

	return &outgoing.User{
		ID:       Uint64toa(user.ID),
		Name:     user.Name,
		Handle:   user.Handle,
		Username: user.Username,
		Email:    user.Email,
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
		var ext = in.Meta.Preview.Extension
		if ext == "" {
			ext = "jpg"
		}

		preview = fmt.Sprintf(attachmentPreviewURL, in.ID, ext)
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

func Command(cmd *sam.Command) *outgoing.Command {
	if cmd == nil {
		return nil
	}

	return &outgoing.Command{
		Name:        cmd.Name,
		Description: cmd.Description,
	}
}

func Commands(cc sam.CommandSet) *outgoing.CommandSet {
	out := make([]*outgoing.Command, len(cc))
	for k, m := range cc {
		out[k] = Command(m)
	}
	retval := outgoing.CommandSet(out)
	return &retval
}
