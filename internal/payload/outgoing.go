package payload

import (
	"context"
	"fmt"
	"net/url"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	messagingTypes "github.com/cortezaproject/corteza-server/messaging/types"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

const (
	attachmentURL        = "/attachment/%d/original/%s"
	attachmentPreviewURL = "/attachment/%d/preview.%s"
)

func Activity(a *messagingTypes.Activity) *outgoing.Activity {
	return &outgoing.Activity{
		MessageID: a.MessageID,
		ChannelID: a.ChannelID,
		Kind:      a.Kind,
		UserID:    a.UserID,
	}
}

func Message(ctx context.Context, msg *messagingTypes.Message) *outgoing.Message {
	var currentUserID = auth.GetIdentityFromContext(ctx).Identity()
	var canEdit = msg.Type.IsEditable() && msg.UserID == currentUserID
	var canReply = msg.Type.IsRepliable() && msg.ReplyTo == 0

	return &outgoing.Message{
		ID:        msg.ID,
		Type:      string(msg.Type),
		ChannelID: msg.ChannelID,
		Message:   msg.Message,
		UserID:    msg.UserID,

		ReplyTo:     msg.ReplyTo,
		Replies:     msg.Replies,
		RepliesFrom: Uint64stoa(msg.RepliesFrom),
		Unread:      Unread(msg.Unread),

		Attachment:   Attachment(msg.Attachment, currentUserID),
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

func Messages(ctx context.Context, msg messagingTypes.MessageSet) *outgoing.MessageSet {
	msgs := make([]*outgoing.Message, len(msg))
	for k, m := range msg {
		msgs[k] = Message(ctx, m)
	}
	retval := outgoing.MessageSet(msgs)
	return &retval
}

func messageReactionSumSet(flags messagingTypes.MessageFlagSet) outgoing.MessageReactionSumSet {
	var (
		rr     = make([]*outgoing.MessageReactionSum, 0)
		rIndex = map[string]int{}
		has    bool
		i      int
	)

	_ = flags.Walk(func(flag *messagingTypes.MessageFlag) error {
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
func messageMentionSet(mm messagingTypes.MentionSet) outgoing.MessageMentionSet {
	return Uint64stoa(mm.UserIDs())
}

func MessageReaction(f *messagingTypes.MessageFlag) *outgoing.MessageReaction {
	return &outgoing.MessageReaction{
		UserID:    f.UserID,
		MessageID: f.MessageID,
		Reaction:  f.Flag,
	}
}

func MessageReactionRemoved(f *messagingTypes.MessageFlag) *outgoing.MessageReactionRemoved {
	return &outgoing.MessageReactionRemoved{
		UserID:    f.UserID,
		MessageID: f.MessageID,
		Reaction:  f.Flag,
	}
}

func MessagePin(f *messagingTypes.MessageFlag) *outgoing.MessagePin {
	return &outgoing.MessagePin{
		UserID:    f.UserID,
		MessageID: f.MessageID,
	}
}

func MessagePinRemoved(f *messagingTypes.MessageFlag) *outgoing.MessagePinRemoved {
	return &outgoing.MessagePinRemoved{
		UserID:    f.UserID,
		MessageID: f.MessageID,
	}
}

func Channel(ch *messagingTypes.Channel) *outgoing.Channel {
	var flag = messagingTypes.ChannelMembershipFlagNone

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

func Channels(channels messagingTypes.ChannelSet) *outgoing.ChannelSet {
	cc := make([]*outgoing.Channel, len(channels))
	for k, c := range channels {
		cc[k] = Channel(c)
	}
	retval := outgoing.ChannelSet(cc)
	return &retval
}

func ChannelMember(m *messagingTypes.ChannelMember) *outgoing.ChannelMember {
	return &outgoing.ChannelMember{
		UserID:    m.UserID,
		Type:      string(m.Type),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ChannelMembers(members messagingTypes.ChannelMemberSet) *outgoing.ChannelMemberSet {
	mm := make([]*outgoing.ChannelMember, len(members))
	for k, c := range members {
		mm[k] = ChannelMember(c)
	}
	retval := outgoing.ChannelMemberSet(mm)
	return &retval
}

func Unread(v *messagingTypes.Unread) *outgoing.Unread {
	if v == nil {
		return nil
	}

	return &outgoing.Unread{
		LastMessageID: v.LastMessageID,
		Count:         v.Count,
		InThreadCount: v.InThreadCount,
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

func Attachment(in *messagingTypes.Attachment, userID uint64) *outgoing.Attachment {
	if in == nil {
		return nil
	}

	var (
		signParams = fmt.Sprintf("?sign=%s&userID=%d", auth.DefaultSigner.Sign(userID, in.ID), userID)
		preview    string
	)

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
		Url:        fmt.Sprintf(attachmentURL, in.ID, url.PathEscape(in.Name)) + signParams,
		PreviewUrl: preview + signParams,
		Meta:       in.Meta,
		Name:       in.Name,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}
}

func Command(cmd *messagingTypes.Command) *outgoing.Command {
	if cmd == nil {
		return nil
	}

	return &outgoing.Command{
		Name:        cmd.Name,
		Description: cmd.Description,
	}
}

func Commands(cc messagingTypes.CommandSet) *outgoing.CommandSet {
	out := make([]*outgoing.Command, len(cc))
	for k, m := range cc {
		out[k] = Command(m)
	}
	retval := outgoing.CommandSet(out)
	return &retval
}
