package payload

import (
	"context"
	"fmt"
	"net/url"

	authTypes "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload/outgoing"
	samTypes "github.com/crusttech/crust/sam/types"
)

const (
	attachmentURL        = "/attachment/%d/original/%s"
	attachmentPreviewURL = "/attachment/%d/preview.%s"
)

func Message(ctx context.Context, msg *samTypes.Message) *outgoing.Message {
	var currentUserID = auth.GetIdentityFromContext(ctx).Identity()
	var canEdit = msg.Type.IsEditable() && msg.UserID == currentUserID
	var canReply = msg.Type.IsRepliable() && msg.ReplyTo == 0

	return &outgoing.Message{
		ID:        msg.ID,
		ChannelID: Uint64toa(msg.ChannelID),
		Message:   msg.Message,
		Type:      string(msg.Type),
		ReplyTo:   msg.ReplyTo,
		Replies:   msg.Replies,

		User:         User(msg.User),
		Attachment:   Attachment(msg.Attachment),
		Reactions:    MessageReactions(msg.Flags),
		IsPinned:     msg.Flags.IsPinned(),
		IsBookmarked: msg.Flags.IsBookmarked(currentUserID),

		CanReply:  canReply,
		CanEdit:   canEdit,
		CanDelete: canEdit,

		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		DeletedAt: msg.DeletedAt,
	}
}

func Messages(ctx context.Context, msg samTypes.MessageSet) *outgoing.MessageSet {
	msgs := make([]*outgoing.Message, len(msg))
	for k, m := range msg {
		msgs[k] = Message(ctx, m)
	}
	retval := outgoing.MessageSet(msgs)
	return &retval
}

func MessageReactions(flags samTypes.MessageFlagSet) outgoing.ReactionSet {
	var (
		rr     = make([]*outgoing.Reaction, 0)
		rIndex = map[string]int{}
		has    bool
		i      int
	)

	_ = flags.Walk(func(flag *samTypes.MessageFlag) error {
		if flag.IsReaction() {
			r := &outgoing.Reaction{Reaction: flag.Flag, UserIDs: []string{}, Count: 0}

			if i, has = rIndex[flag.Flag]; !has {
				i, rIndex[flag.Flag] = len(rr), len(rr)
				rr = append(rr, r)
			}

			rr[i].UserIDs = append(rr[i].UserIDs, Uint64toa(flag.UserID))
			rr[i].Count++
		}

		return nil
	})

	return rr
}

func Channel(ch *samTypes.Channel) *outgoing.Channel {
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

func Channels(channels samTypes.ChannelSet) *outgoing.ChannelSet {
	cc := make([]*outgoing.Channel, len(channels))
	for k, c := range channels {
		cc[k] = Channel(c)
	}
	retval := outgoing.ChannelSet(cc)
	return &retval
}

func ChannelMember(m *samTypes.ChannelMember) *outgoing.ChannelMember {
	return &outgoing.ChannelMember{
		User:      User(m.User),
		Type:      string(m.Type),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ChannelMembers(members samTypes.ChannelMemberSet) *outgoing.ChannelMemberSet {
	mm := make([]*outgoing.ChannelMember, len(members))
	for k, c := range members {
		mm[k] = ChannelMember(c)
	}
	retval := outgoing.ChannelMemberSet(mm)
	return &retval
}

func ChannelView(v *samTypes.ChannelView) *outgoing.ChannelView {
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

func User(user *authTypes.User) *outgoing.User {
	if user == nil {
		return nil
	}

	return &outgoing.User{
		ID:       user.ID,
		Name:     user.Name,
		Handle:   user.Handle,
		Username: user.Username,
		Email:    user.Email,
	}
}

func Users(users []*authTypes.User) *outgoing.UserSet {
	uu := make([]*outgoing.User, len(users))
	for k, u := range users {
		uu[k] = User(u)
	}

	retval := outgoing.UserSet(uu)

	return &retval
}

func Attachment(in *samTypes.Attachment) *outgoing.Attachment {
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

func Command(cmd *samTypes.Command) *outgoing.Command {
	if cmd == nil {
		return nil
	}

	return &outgoing.Command{
		Name:        cmd.Name,
		Description: cmd.Description,
	}
}

func Commands(cc samTypes.CommandSet) *outgoing.CommandSet {
	out := make([]*outgoing.Command, len(cc))
	for k, m := range cc {
		out[k] = Command(m)
	}
	retval := outgoing.CommandSet(out)
	return &retval
}
