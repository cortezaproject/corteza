package payload

import (
	"context"
	"fmt"
	"net/url"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload/outgoing"
	samTypes "github.com/crusttech/crust/messaging/types"
	systemTypes "github.com/crusttech/crust/system/types"
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

		ReplyTo:     msg.ReplyTo,
		Replies:     msg.Replies,
		RepliesFrom: Uint64stoa(msg.RepliesFrom),

		User:         User(msg.User),
		Attachment:   Attachment(msg.Attachment),
		Mentions:     messageMentionSet(msg.Mentions),
		Reactions:    messageReactionSumSet(msg.Flags),
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

func messageReactionSumSet(flags samTypes.MessageFlagSet) outgoing.MessageReactionSumSet {
	var (
		rr     = make([]*outgoing.MessageReactionSum, 0)
		rIndex = map[string]int{}
		has    bool
		i      int
	)

	_ = flags.Walk(func(flag *samTypes.MessageFlag) error {
		if flag.IsReaction() {
			r := &outgoing.MessageReactionSum{Reaction: flag.Flag, UserIDs: []string{}, Count: 0}

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

// Converts slice of mentions into slice of strings containing all user IDs
// These are IDs of users mentioned in the message
func messageMentionSet(mm samTypes.MentionSet) outgoing.MessageMentionSet {
	return Uint64stoa(mm.UserIDs())
}

func MessageReaction(f *samTypes.MessageFlag) *outgoing.MessageReaction {
	return &outgoing.MessageReaction{
		UserID:    f.UserID,
		MessageID: f.MessageID,
		Reaction:  f.Flag,
	}
}

func MessageReactionRemoved(f *samTypes.MessageFlag) *outgoing.MessageReactionRemoved {
	return &outgoing.MessageReactionRemoved{
		UserID:    f.UserID,
		MessageID: f.MessageID,
		Reaction:  f.Flag,
	}
}

func MessagePin(f *samTypes.MessageFlag) *outgoing.MessagePin {
	return &outgoing.MessagePin{
		UserID:    f.UserID,
		MessageID: f.MessageID,
	}
}

func MessagePinRemoved(f *samTypes.MessageFlag) *outgoing.MessagePinRemoved {
	return &outgoing.MessagePinRemoved{
		UserID:    f.UserID,
		MessageID: f.MessageID,
	}
}

func Channel(ch *samTypes.Channel) *outgoing.Channel {
	var flag = samTypes.ChannelMembershipFlagNone

	if ch.Member != nil {
		flag = ch.Member.Flag
	}

	return &outgoing.Channel{
		ID:             Uint64toa(ch.ID),
		Name:           ch.Name,
		LastMessageID:  Uint64toa(ch.LastMessageID),
		Topic:          ch.Topic,
		Type:           string(ch.Type),
		MembershipFlag: string(flag),
		Members:        Uint64stoa(ch.Members),
		Unread:         Unread(ch.Unread),

		CanJoin:           ch.CanJoin,
		CanPart:           ch.CanPart,
		CanObserve:        ch.CanObserve,
		CanSendMessages:   ch.CanSendMessages,
		CanDeleteMessages: ch.CanDeleteMessages,
		CanChangeMembers:  ch.CanChangeMembers,
		CanUpdate:         ch.CanUpdate,
		CanArchive:        ch.CanArchive,
		CanDelete:         ch.CanDelete,

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

func Unread(v *samTypes.Unread) *outgoing.Unread {
	if v == nil {
		return nil
	}

	return &outgoing.Unread{
		LastMessageID: v.LastMessageID,
		Count:         v.Count,
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

func User(user *systemTypes.User) *outgoing.User {
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

func Users(users []*systemTypes.User) *outgoing.UserSet {
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
