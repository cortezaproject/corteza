package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/store"
	"io"
	"regexp"
	"strings"
)

type (
	message struct {
		ctx     context.Context
		ac      messageAccessController
		channel ChannelService
		store   store.Storer
		event   EventService
	}

	messageAccessController interface {
		CanReadChannel(context.Context, *types.Channel) bool
		CanReplyMessage(context.Context, *types.Channel) bool
		CanSendMessage(context.Context, *types.Channel) bool
		CanUpdateMessages(context.Context, *types.Channel) bool
		CanUpdateOwnMessages(context.Context, *types.Channel) bool
		CanReactMessage(context.Context, *types.Channel) bool
	}

	MessageService interface {
		With(ctx context.Context) MessageService

		Find(types.MessageFilter) (types.MessageSet, types.MessageFilter, error)
		FindThreads(types.MessageFilter) (types.MessageSet, types.MessageFilter, error)

		Create(messages *types.Message) (*types.Message, error)
		Update(messages *types.Message) (*types.Message, error)

		CreateWithAvatar(message *types.Message, avatar io.Reader) (*types.Message, error)

		React(messageID uint64, reaction string) error
		RemoveReaction(messageID uint64, reaction string) error

		MarkAsRead(channelID, threadID, lastReadMessageID uint64) (uint64, uint32, uint32, error)

		Pin(messageID uint64) error
		RemovePin(messageID uint64) error

		Bookmark(messageID uint64) error
		RemoveBookmark(messageID uint64) error

		Delete(messageID uint64) error
	}
)

const (
	settingsMessageBodyLength = 0
	mentionRE                 = `<([@#])(\d+)((?:\s)([^>]+))?>`
)

var (
	mentionsFinder = regexp.MustCompile(mentionRE)
)

func Message(ctx context.Context) MessageService {
	return (&message{
		ac:      DefaultAccessControl,
		channel: DefaultChannel,
		store:   DefaultStore,
	}).With(ctx)
}

func (svc message) With(ctx context.Context) MessageService {
	return &message{
		ctx:     ctx,
		ac:      svc.ac,
		channel: svc.channel,
		store:   svc.store,
		event:   Event(ctx),
	}
}

func (svc message) Find(filter types.MessageFilter) (mm types.MessageSet, f types.MessageFilter, err error) {
	f = filter
	f.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	if f.ChannelID, err = svc.readableChannels(f); err != nil {
		return
	}

	mm, f, err = store.SearchMessagingMessages(svc.ctx, svc.store, f)
	if err != nil {
		return
	}

	return mm, f, svc.preload(svc.ctx, svc.store, mm)
}

func (svc message) FindThreads(filter types.MessageFilter) (mm types.MessageSet, f types.MessageFilter, err error) {
	f = filter
	f.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	if f.ChannelID, err = svc.readableChannels(f); err != nil {
		return
	}

	mm, f, err = store.SearchMessagingThreads(svc.ctx, svc.store, f)
	if err != nil {
		return
	}

	return mm, f, svc.preload(svc.ctx, svc.store, mm)
}

func (svc message) CreateWithAvatar(in *types.Message, avatar io.Reader) (*types.Message, error) {
	// @todo: avatar
	return svc.Create(in)
}

// Returns list of readable channels
//
// Either all (when len(f.ChannelID) == 0) or subset of channel IDs (from f.ChannelID)
func (svc message) readableChannels(f types.MessageFilter) ([]uint64, error) {
	cc, _, err := svc.channel.With(svc.ctx).Find(types.ChannelFilter{
		CurrentUserID:  f.CurrentUserID,
		ChannelID:      f.ChannelID,
		IncludeDeleted: true,
	})

	if err != nil {
		return nil, err
	}

	if len(cc) == 0 {
		// None of the channels requested were returned as accessible
		return nil, ErrNoPermissions.withStack()
	}

	return cc.IDs(), nil
}

func (svc message) Create(msg *types.Message) (*types.Message, error) {
	msg.ID = nextID()
	msg.CreatedAt = *now()
	msg.Message = strings.TrimSpace(msg.Message)

	if l := len(msg.Message); l == 0 {
		return nil, fmt.Errorf("refusing to create message without contents")
	} else if settingsMessageBodyLength > 0 && l > settingsMessageBodyLength {
		return nil, fmt.Errorf("message length (%d characters) too long (max: %d)", l, settingsMessageBodyLength)
	}

	// keep pre-existing user id set
	if msg.UserID == 0 {
		msg.UserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	}

	err := store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}
		var ch *types.Channel

		if msg.ReplyTo > 0 {
			var original *types.Message
			var replyTo = msg.ReplyTo

			for replyTo > 0 {
				// Find original message
				original, err = store.LookupMessagingMessageByID(ctx, s, msg.ReplyTo)
				if err != nil {
					return
				}

				replyTo = original.ReplyTo
			}

			if !original.Type.IsRepliable() {
				return fmt.Errorf("unable to reply on this message (type = %s)", original.Type)
			}

			// We do not want to have multi-level threads
			// Take original's reply-to and use it
			msg.ReplyTo = original.ID

			msg.ChannelID = original.ChannelID

			if original.Replies == 0 {
				// First reply,
				//
				// reset unreads for all members
				var mm types.ChannelMemberSet
				mm, _, err = store.SearchMessagingChannelMembers(ctx, s, types.ChannelMemberFilterChannels(original.ChannelID))
				if err != nil {
					return err
				}

				if err = store.PresetMessagingUnread(ctx, s, original.ChannelID, original.ID, mm.AllMemberIDs()...); err != nil {
					return err
				}
			}

			// Increment counter, on struct and in store.
			original.Replies++
			if err = store.UpdateMessagingMessageReplyCount(ctx, s, original.ID, original.Replies); err != nil {
				return
			}

			// Broadcast updated original
			bq = append(bq, original)
		}

		if msg.ChannelID == 0 {
			return fmt.Errorf("channelID missing")
		} else if ch, err = svc.findChannelByID(msg.ChannelID); err != nil {
			return
		}

		if msg.ReplyTo > 0 && !svc.ac.CanReplyMessage(ctx, ch) {
			return ErrNoPermissions.withStack()
		}
		if !svc.ac.CanSendMessage(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if err = store.CreateMessagingMessage(ctx, s, msg); err != nil {
			return
		}

		mentions := svc.extractMentions(msg)
		if err = svc.updateMentions(ctx, s, msg.ID, mentions); err != nil {
			return
		}

		svc.sendNotifications(msg, mentions)

		// Count unreads in the background and send updates to all users
		if err = svc.countUnreads(ctx, s, ch, msg, 0); err != nil {
			return
		}

		return svc.sendEvent(append(bq, msg)...)
	})

	if err != nil {
		return nil, err
	}

	return msg, err
}

func (svc message) Update(upd *types.Message) (msg *types.Message, err error) {
	if upd.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	upd.Message = strings.TrimSpace(upd.Message)

	if l := len(upd.Message); l == 0 {
		return nil, fmt.Errorf("refusing to update message without contents")
	} else if settingsMessageBodyLength > 0 && l > settingsMessageBodyLength {
		return nil, fmt.Errorf("message length (%d characters) too long (max: %d)", l, settingsMessageBodyLength)
	}

	var currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var ch *types.Channel

		if ch, err = svc.findChannelByID(upd.ChannelID); err != nil {
			return err
		} else if !svc.ac.CanReadChannel(ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		msg, err = store.LookupMessagingMessageByID(ctx, s, upd.ID)

		if err != nil {
			return fmt.Errorf("could not load message for editing: %w", err)
		}

		if msg.Message == upd.Message {
			// Nothing changed
			return nil
		}

		if msg.UserID == currentUserID && !svc.ac.CanUpdateOwnMessages(ctx, ch) {
			return ErrNoPermissions.withStack()
		} else if msg.UserID != currentUserID && !svc.ac.CanUpdateMessages(ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		// Update message
		msg.Message = upd.Message
		msg.UpdatedAt = now()

		if err = store.UpdateMessagingMessage(ctx, s, msg); err != nil {
			return err
		}

		if err = svc.updateMentions(ctx, s, msg.ID, svc.extractMentions(msg)); err != nil {
			return
		}

		return svc.sendEvent(msg)
	})

	if err != nil {
		return nil, err
	}

	return msg, err
}

func (svc message) Delete(messageID uint64) (err error) {
	var currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	_ = currentUserID

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}
		var deletedMsg, original *types.Message
		var ch *types.Channel

		deletedMsg, err = store.LookupMessagingMessageByID(ctx, s, messageID)
		if err != nil {
			return err
		}

		deletedMsg.DeletedAt = now()

		if ch, err = svc.findChannelByID(deletedMsg.ChannelID); err != nil {
			return err
		} else if !svc.ac.CanReadChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if deletedMsg.UserID == currentUserID && !svc.ac.CanUpdateOwnMessages(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		} else if deletedMsg.UserID != currentUserID && !svc.ac.CanUpdateMessages(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if deletedMsg.ReplyTo > 0 {
			original, err = store.LookupMessagingMessageByID(ctx, s, deletedMsg.ReplyTo)
			if err != nil {
				return err
			}

			// This is a reply to another message, decrease reply counter on the original, on struct and in the
			// repository
			if original.Replies > 0 {
				original.Replies--
			}

			if err = store.UpdateMessagingMessageReplyCount(ctx, s, original.ID, original.Replies); err != nil {
				return
			}

			// Broadcast updated original
			bq = append(bq, original)
		}

		if err = store.UpdateMessagingMessage(ctx, s, deletedMsg); err != nil {
			return
		}

		// Set deletedAt timestamp so that our clients can react properly...
		deletedMsg.DeletedAt = now()

		if err = svc.updateMentions(ctx, s, messageID, nil); err != nil {
			return
		}

		// Count unreads in the background and send updates to all users
		if err = svc.countUnreads(ctx, s, ch, deletedMsg, 0); err != nil {
			return
		}

		return svc.sendEvent(append(bq, deletedMsg)...)
	})

	return
}

// MarkAsRead marks channel/thread as read
//
// If lastReadMessageID is set, it uses that message as last read message
func (svc message) MarkAsRead(channelID, threadID, lastReadMessageID uint64) (uint64, uint32, uint32, error) {
	var (
		currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

		count       uint32
		threadCount uint32
		err         error
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var ch *types.Channel
		var thread *types.Message
		var lastMessage *types.Message

		if ch, err = svc.findChannelByID(channelID); err != nil {
			return err
		} else if !svc.ac.CanReadChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		} else if !ch.IsValid() {
			return fmt.Errorf("invalid channel")
		}

		if threadID > 0 {
			// Validate thread
			if thread, err = store.LookupMessagingMessageByID(ctx, s, threadID); err != nil {
				return fmt.Errorf("unable to verify thread: %w", err)
			} else if !thread.IsValid() {
				return fmt.Errorf("invalid thread")
			}
		} else {
			// This is request for channel,
			// count all thread unreads
			var uu types.UnreadSet
			uu, err = store.CountMessagingUnreadThreads(ctx, s, currentUserID, channelID)
			if err != nil {
				return err
			}

			if u := uu.FindByChannelId(channelID); u != nil {
				threadCount = u.ThreadCount
			}
		}

		if lastReadMessageID > 0 {
			// Validate messageID/threadID/channelID combo
			if lastMessage, err = store.LookupMessagingMessageByID(ctx, s, lastReadMessageID); err != nil {
				return fmt.Errorf("unable to verify last message: %w", err)
			} else if !lastMessage.IsValid() {
				return fmt.Errorf("invalid message")
			} else if lastMessage.ChannelID != channelID {
				return fmt.Errorf("last read message not in the same channel")
			} else if threadID > 0 && lastMessage.ReplyTo != threadID {
				return fmt.Errorf("last read message not in the same thread")
			}

			count, err = store.CountMessagingMessagesFromID(ctx, s, channelID, threadID, lastReadMessageID)
			if err != nil {
				return fmt.Errorf("unable to count unread messages: %w", err)
			}

		} else {
			// use last message ID
			if lastReadMessageID, err = store.LastMessagingMessageID(ctx, s, channelID, threadID); err != nil {
				return fmt.Errorf("unable to find last message: %w", err)
			}

			// no need to count
			count = 0
		}

		u := &types.Unread{
			ChannelID:     channelID,
			UserID:        currentUserID,
			ReplyTo:       threadID,
			LastMessageID: lastReadMessageID,
			Count:         count,
		}

		if err = store.UpsertMessagingUnread(ctx, s, u); err != nil {
			return fmt.Errorf("unable to record unread messages: %w", err)
		}

		// Remove unread counts from all threads when doing mark-channel-as-read
		if threadID == 0 {
			if err = store.ResetMessagingUnreadThreads(ctx, s, channelID, currentUserID); err != nil {
				return fmt.Errorf("unable to clear channel threads: %w", err)
			}
		}

		// Re-count unreads and send updates to this user
		if err = svc.countUnreads(ctx, s, ch, nil, currentUserID); err != nil {
			return
		}

		return nil
	})

	return lastReadMessageID, count, threadCount, err
}

// React on a message with an emoji
func (svc message) React(messageID uint64, reaction string) error {
	return svc.flag(messageID, reaction, false)
}

// Remove reaction on a message
func (svc message) RemoveReaction(messageID uint64, reaction string) error {
	return svc.flag(messageID, reaction, true)
}

// Pin message to the channel
func (svc message) Pin(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagPinnedToChannel, false)
}

// Remove pin from message
func (svc message) RemovePin(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagPinnedToChannel, true)
}

// Bookmark message (private)
func (svc message) Bookmark(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagBookmarkedMessage, false)
}

// Remove bookmark message (private)
func (svc message) RemoveBookmark(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagBookmarkedMessage, true)
}

// React on a message with an emoji
func (svc message) flag(messageID uint64, flag string, remove bool) (err error) {
	var (
		currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
		flagOwnerID   = currentUserID
	)

	if currentUserID == 0 {
		return fmt.Errorf("unable to flag with unknown user")
	}

	if strings.TrimSpace(flag) == "" {
		// Sanitize
		flag = types.MessageFlagPinnedToChannel
	}
	// @todo validate flags more strictly

	if flag == types.MessageFlagPinnedToChannel {
		// It does not matter how is the owner of the pin,
		flagOwnerID = 0
	}

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			ff  types.MessageFlagSet
			f   *types.MessageFlag
			msg *types.Message
			ch  *types.Channel

			filter = types.MessageFlagFilter{
				Flag:      flag,
				MessageID: []uint64{messageID},
				OwnerID:   flagOwnerID,
			}
		)

		if ff, _, err = store.SearchMessagingFlags(ctx, s, filter); err != nil {
			return
		}

		if len(ff) == 0 && remove {
			// Skip removing, flag does not exists
			return nil
		}

		if len(ff) > 0 && !remove {
			// Skip adding, flag already exists
			return nil
		}

		if msg, err = store.LookupMessagingMessageByID(ctx, s, messageID); err != nil {
			return
		}

		if ch, err = svc.findChannelByID(msg.ChannelID); err != nil {
			return
		}

		if !svc.ac.CanReadChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if f != nil && f.IsReaction() && !svc.ac.CanReactMessage(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if remove {
			// setting deleted to inform clients about removed flag
			_ = ff.Walk(func(f *types.MessageFlag) error {
				f.DeletedAt = now()
				return nil
			})
			err = store.DeleteMessagingFlag(ctx, s, ff...)
		} else {
			ff = []*types.MessageFlag{
				{
					ID:        nextID(),
					UserID:    currentUserID,
					CreatedAt: *now(),
					ChannelID: msg.ChannelID,
					MessageID: msg.ID,
					Flag:      flag,
				},
			}
			err = store.CreateMessagingFlag(ctx, s, ff...)
		}

		if err != nil {
			return
		}

		_ = svc.sendFlagEvent(ff...)
		return
	})

	if err != nil {
		return fmt.Errorf("can not flag/un-flag message: %w", err)
	}

	return nil
}

func (svc message) preload(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	if err = svc.preloadAttachments(ctx, s, mm); err != nil {
		return
	}

	if err = svc.preloadFlags(ctx, s, mm); err != nil {
		return
	}

	if err = svc.preloadMentions(ctx, s, mm); err != nil {
		return
	}

	if err = svc.preloadUnreads(ctx, s, mm); err != nil {
		return
	}

	if err = svc.preloadThreadParticipants(ctx, s, mm); err != nil {
		return
	}

	return
}

// Preload for all messages
func (message) preloadFlags(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	if len(mm) == 0 {
		return
	}

	var ff types.MessageFlagSet
	ff, _, err = store.SearchMessagingFlags(ctx, s, types.MessageFlagFilter{MessageID: mm.IDs()})
	if err != nil {
		return
	}

	return ff.Walk(func(flag *types.MessageFlag) error {
		mm.FindByID(flag.MessageID).Flags = append(mm.FindByID(flag.MessageID).Flags, flag)
		return nil
	})
}

// Preload for all messages
func (message) preloadMentions(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	if len(mm) == 0 {
		return
	}

	var mentions types.MentionSet
	mentions, _, err = store.SearchMessagingMentions(ctx, s, types.MentionFilter{MessageID: mm.IDs()})
	if err != nil {
		return
	}

	return mm.Walk(func(m *types.Message) error {
		m.Mentions = mentions.FindByMessageID(m.ID)
		return nil
	})
}

func (message) preloadAttachments(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	var (
		aa types.AttachmentSet
	)

	if err != nil || len(mm) == 0 {
		return
	}

	if aa, _, err = store.SearchMessagingAttachments(ctx, s, types.AttachmentFilter{MessageID: mm.IDs()}); err != nil {
		return
	} else {
		_ = aa
		//return aa.Walk(func(a *types.Attachment) error {
		//	if a.MessageID > 0 {
		//		if m := mm.FindByID(a.MessageID); m != nil {
		//			m.Attachment = &a.Attachment
		//		}
		//	}
		//
		//	return nil
		//})
		return
	}
}

func (message) preloadUnreads(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	if len(mm) == 0 {
		return nil
	}

	var (
		unread types.UnreadSet
		userID = auth.GetIdentityFromContext(ctx).Identity()
	)

	// Filter out only relevant messages -- ones with replies
	mm, _ = mm.Filter(func(m *types.Message) (b bool, e error) {
		return m.Replies > 0, nil
	})

	unread, err = store.CountMessagingUnread(ctx, s, userID, 0, mm.IDs()...)
	if err != nil {
		return
	}

	return mm.Walk(func(m *types.Message) error {
		m.Unread = unread.FindByThreadId(m.ID)
		return nil
	})
}

func (message) preloadThreadParticipants(ctx context.Context, s store.Storer, mm types.MessageSet) (err error) {
	if len(mm) == 0 {
		return nil
	}

	var (
		unread types.UnreadSet
		userID = auth.GetIdentityFromContext(ctx).Identity()
	)

	// Filter out only relevant messages -- ones with replies
	mm, _ = mm.Filter(func(m *types.Message) (b bool, e error) {
		return m.Replies > 0, nil
	})

	unread, err = store.CountMessagingUnread(ctx, s, userID, 0, mm.IDs()...)
	if err != nil {
		return
	}

	return mm.Walk(func(m *types.Message) error {
		m.Unread = unread.FindByThreadId(m.ID)
		return nil
	})
}

// Sends message to event loop
func (svc message) sendEvent(mm ...*types.Message) (err error) {
	if err = svc.preload(svc.ctx, svc.store, mm); err != nil {
		return
	}

	for _, msg := range mm {
		if err = svc.event.Message(msg); err != nil {
			return
		}
	}

	return
}

// Generates and sends notifications from the new message
//
//
func (svc message) sendNotifications(message *types.Message, mentions types.MentionSet) {
	// @todo implementation
}

// countUnreads orchestrates unread-related operations (inc/dec, (re)counting & sending events)
//
// 1. increases/decreases unread counters for channel or thread
// 2. collects all counters for channel or thread
// 3. sends unread events to subscribers
func (svc message) countUnreads(ctx context.Context, s store.Storer, ch *types.Channel, m *types.Message, userID uint64) (err error) {
	var (
		uuBase, uuThreads, uuChannels types.UnreadSet
		// mm  types.ChannelMemberSet
		threadIDs []uint64
	)

	if m != nil {
		if m.DeletedAt != nil {
			// When deleting message, all existing counters are decreased!
			if err = store.DecMessagingUnreadCount(ctx, s, m.ChannelID, m.ReplyTo, m.UserID); err != nil {
				return
			}
		} else if m.UpdatedAt == nil {
			// Reset user's counter and set current message ID as last read.
			u := &types.Unread{
				ChannelID:     m.ChannelID,
				UserID:        m.UserID,
				ReplyTo:       m.ReplyTo,
				LastMessageID: m.ID,
				Count:         0,
			}

			if err = store.UpsertMessagingUnread(ctx, s, u); err != nil {
				return err
			}

			// When new message is created, update all existing counters
			if err = store.IncMessagingUnreadCount(ctx, s, m.ChannelID, m.ReplyTo, m.UserID); err != nil {
				return
			}
		}

		if m.ReplyTo > 0 {
			threadIDs = []uint64{m.ReplyTo}
		}
	}

	uuBase, err = store.CountMessagingUnread(ctx, s, userID, ch.ID, threadIDs...)
	if err != nil {
		return
	}

	if len(threadIDs) > 0 {
		// If base count was done for a thread,
		// Do another count for channel
		uuChannels, err = store.CountMessagingUnread(ctx, s, userID, ch.ID)
		if err != nil {
			return
		}

		uuBase = uuBase.Merge(uuChannels)

		// Now recount all threads for this channel
		uuThreads, err = store.CountMessagingUnreadThreads(ctx, s, userID, ch.ID)
		if err != nil {
			return
		}

		uuBase = uuBase.Merge(uuThreads)
	}

	// This is a reply, make sure we fetch the new stats about unread replies and push them to users
	return svc.event.UnreadCounters(uuBase)
}

// Sends message to event loop
func (svc message) sendFlagEvent(ff ...*types.MessageFlag) (err error) {
	for _, f := range ff {
		if err = svc.event.MessageFlag(f); err != nil {
			return
		}
	}

	return
}

func (svc message) extractMentions(m *types.Message) (mm types.MentionSet) {
	const reSubID = 2
	mm = types.MentionSet{}

	match := mentionsFinder.FindAllStringSubmatch(m.Message, -1)

	// Prepopulated with all we know from message
	tpl := types.Mention{
		ID:            nextID(),
		ChannelID:     m.ChannelID,
		MessageID:     m.ID,
		MentionedByID: m.UserID,
		CreatedAt:     *now(),
	}

	for m := 0; m < len(match); m++ {
		uid := payload.ParseUint64(match[m][reSubID])
		if len(mm.FindByUserID(uid)) == 0 {
			// Copy template & assign user id
			mnt := tpl
			mnt.UserID = uid
			mm = append(mm, &mnt)
		}
	}

	return
}

func (svc message) updateMentions(ctx context.Context, s store.Storer, messageID uint64, mm types.MentionSet) error {
	if existing, _, err := store.SearchMessagingMentions(ctx, s, types.MentionFilter{MessageID: []uint64{messageID}}); err != nil {
		return fmt.Errorf("could not update mentions: %w", err)
	} else if len(mm) > 0 {
		add, _, del := existing.Diff(mm)

		if err = store.CreateMessagingMention(ctx, s, add...); err != nil {
			return fmt.Errorf("could not create mentions: %w", err)
		}

		if err = store.DeleteMessagingMention(ctx, s, del...); err != nil {
			return fmt.Errorf("could not delete mentions: %w", err)
		}
	} else {
		return store.DeleteMessagingMentionByID(ctx, s, messageID)
	}

	return nil
}

// findChannelByID loads channel and it's members
func (svc message) findChannelByID(channelID uint64) (ch *types.Channel, err error) {
	var (
		currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
		mm            types.ChannelMemberSet
	)

	if ch, err = svc.channel.With(svc.ctx).FindByID(channelID); err != nil {
		return nil, err
	} else if mm, _, err = store.SearchMessagingChannelMembers(svc.ctx, svc.store, types.ChannelMemberFilterChannels(ch.ID)); err != nil {
		return nil, err
	} else {
		ch.Members = mm.AllMemberIDs()
		ch.Member = mm.FindByUserID(currentUserID)
	}

	return
}

var _ MessageService = &message{}
