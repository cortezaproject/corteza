package service

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/payload"
)

type (
	message struct {
		db     db
		ctx    context.Context
		logger *zap.Logger
		ac     messageAccessController

		channel ChannelService

		attachment repository.AttachmentRepository
		cmember    repository.ChannelMemberRepository
		unread     repository.UnreadRepository
		message    repository.MessageRepository
		mflag      repository.MessageFlagRepository
		mentions   repository.MentionRepository

		event EventService
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
		logger: DefaultLogger.Named("message"),

		ac:      DefaultAccessControl,
		channel: DefaultChannel,
	}).With(ctx)
}

func (svc message) With(ctx context.Context) MessageService {
	db := repository.DB(ctx)
	return &message{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac:      svc.ac,
		channel: svc.channel,

		event: Event(ctx),

		attachment: repository.Attachment(ctx, db),
		cmember:    repository.ChannelMember(ctx, db),
		unread:     repository.Unread(ctx, db),
		message:    repository.Message(ctx, db),
		mflag:      repository.MessageFlag(ctx, db),
		mentions:   repository.Mention(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc message) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc message) Find(filter types.MessageFilter) (mm types.MessageSet, f types.MessageFilter, err error) {
	f = filter
	f.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	if f.ChannelID, err = svc.readableChannels(f); err != nil {
		return
	}

	mm, f, err = svc.message.Find(f)
	if err != nil {
		return
	}

	return mm, f, svc.preload(mm)
}

func (svc message) FindThreads(filter types.MessageFilter) (mm types.MessageSet, f types.MessageFilter, err error) {
	f = filter
	f.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	if f.ChannelID, err = svc.readableChannels(f); err != nil {
		return
	}

	mm, f, err = svc.message.FindThreads(f)
	if err != nil {
		return
	}

	return mm, f, svc.preload(mm)
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

func (svc message) Create(in *types.Message) (m *types.Message, err error) {
	if in == nil {
		in = &types.Message{}
	}

	in.Message = strings.TrimSpace(in.Message)

	var mlen = len(in.Message)

	if mlen == 0 {
		return nil, errors.Errorf("refusing to create message without contents")
	}

	if settingsMessageBodyLength > 0 && mlen > settingsMessageBodyLength {
		return nil, errors.Errorf("message length (%d characters) too long (max: %d)", mlen, settingsMessageBodyLength)
	}

	// keep pre-existing user id set
	if in.UserID == 0 {
		in.UserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	}

	return m, svc.db.Transaction(func() (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}
		var ch *types.Channel

		if in.ReplyTo > 0 {
			var original *types.Message
			var replyTo = in.ReplyTo

			for replyTo > 0 {
				// Find original message
				original, err = svc.message.FindByID(in.ReplyTo)
				if err != nil {
					return
				}

				replyTo = original.ReplyTo
			}

			if !original.Type.IsRepliable() {
				return errors.Errorf("unable to reply on this message (type = %s)", original.Type)
			}

			// We do not want to have multi-level threads
			// Take original's reply-to and use it
			in.ReplyTo = original.ID

			in.ChannelID = original.ChannelID

			if original.Replies == 0 {
				// First reply,
				//
				// reset unreads for all members
				var mm types.ChannelMemberSet
				mm, err = svc.cmember.Find(types.ChannelMemberFilterChannels(original.ChannelID))
				if err != nil {
					return err
				}

				err = svc.unread.Preset(original.ChannelID, original.ID, mm.AllMemberIDs()...)
				if err != nil {
					return err
				}
			}

			// Increment counter, on struct and in repository.
			original.Replies++
			if err = svc.message.IncReplyCount(original.ID); err != nil {
				return
			}

			// Broadcast updated original
			bq = append(bq, original)
		}

		if in.ChannelID == 0 {
			return errors.New("channelID missing")
		} else if ch, err = svc.findChannelByID(in.ChannelID); err != nil {
			return
		}

		if in.ReplyTo > 0 && !svc.ac.CanReplyMessage(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}
		if !svc.ac.CanSendMessage(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if m, err = svc.message.Create(in); err != nil {
			return
		}

		mentions := svc.extractMentions(m)
		if err = svc.updateMentions(m.ID, mentions); err != nil {
			return
		}

		svc.sendNotifications(m, mentions)

		// Count unreads in the background and send updates to all users
		svc.countUnreads(ch, m, 0)

		return svc.sendEvent(append(bq, m)...)
	})
}

func (svc message) Update(in *types.Message) (message *types.Message, err error) {
	if in.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if in == nil {
		in = &types.Message{}
	}

	in.Message = strings.TrimSpace(in.Message)
	var mlen = len(in.Message)

	if mlen == 0 {
		return nil, errors.Errorf("refusing to update message without contents")
	}
	if settingsMessageBodyLength > 0 && mlen > settingsMessageBodyLength {
		return nil, errors.Errorf("message length (%d characters) too long (max: %d)", mlen, settingsMessageBodyLength)
	}

	var currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return message, svc.db.Transaction(func() (err error) {
		var ch *types.Channel

		if ch, err = svc.findChannelByID(in.ChannelID); err != nil {
			return err
		} else if !svc.ac.CanReadChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		message, err = svc.message.FindByID(in.ID)

		if err != nil {
			return errors.Wrap(err, "could not load message for editing")
		}

		if message.Message == in.Message {
			// Nothing changed
			return nil
		}

		if message.UserID == currentUserID && !svc.ac.CanUpdateOwnMessages(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		} else if message.UserID != currentUserID && !svc.ac.CanUpdateMessages(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		// Allow message content to be changed
		message.Message = in.Message

		if message, err = svc.message.Update(message); err != nil {
			return err
		}

		if err = svc.updateMentions(message.ID, svc.extractMentions(message)); err != nil {
			return
		}

		return svc.sendEvent(message)
	})
}

func (svc message) Delete(messageID uint64) error {
	var currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	_ = currentUserID

	return svc.db.Transaction(func() (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}
		var deletedMsg, original *types.Message
		var ch *types.Channel

		deletedMsg, err = svc.message.FindByID(messageID)
		if err != nil {
			return err
		}

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
			original, err = svc.message.FindByID(deletedMsg.ReplyTo)
			if err != nil {
				return err
			}

			// This is a reply to another message, decrease reply counter on the original, on struct and in the
			// repository
			if original.Replies > 0 {
				original.Replies--
			}

			if err = svc.message.DecReplyCount(original.ID); err != nil {
				return err
			}

			// Broadcast updated original
			bq = append(bq, original)
		}

		if err = svc.message.DeleteByID(messageID); err != nil {
			return
		}

		// Set deletedAt timestamp so that our clients can react properly...
		deletedMsg.DeletedAt = timeNowPtr()

		if err = svc.updateMentions(messageID, nil); err != nil {
			return
		}

		// Count unreads in the background and send updates to all users
		svc.countUnreads(ch, deletedMsg, 0)

		return svc.sendEvent(append(bq, deletedMsg)...)
	})
}

// MarkAsRead marks channel/thread as read
//
// If lastReadMessageID is set, it uses that message as last read message
func (svc message) MarkAsRead(channelID, threadID, lastReadMessageID uint64) (uint64, uint32, uint32, error) {
	var (
		currentUserID uint64 = auth.GetIdentityFromContext(svc.ctx).Identity()
		count         uint32
		threadCount   uint32
		err           error
	)

	err = svc.db.Transaction(func() (err error) {
		var ch *types.Channel
		var thread *types.Message
		var lastMessage *types.Message

		if ch, err = svc.findChannelByID(channelID); err != nil {
			return err
		} else if !svc.ac.CanReadChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		} else if !ch.IsValid() {
			return errors.New("invalid channel")
		}

		if threadID > 0 {
			// Validate thread
			if thread, err = svc.message.FindByID(threadID); err != nil {
				return errors.Wrap(err, "unable to verify thread")
			} else if !thread.IsValid() {
				return errors.New("invalid thread")
			}
		} else {
			// This is request for channel,
			// count all thread unreads
			var uu types.UnreadSet
			uu, err = svc.unread.CountThreads(currentUserID, channelID)
			if err != nil {
				return err
			}

			if u := uu.FindByChannelId(channelID); u != nil {
				threadCount = u.ThreadCount
			}
		}

		if lastReadMessageID > 0 {
			// Validate messageID/threadID/channelID combo
			if lastMessage, err = svc.message.FindByID(lastReadMessageID); err != nil {
				return errors.Wrap(err, "unable to verify last message")
			} else if !lastMessage.IsValid() {
				return errors.New("invalid message")
			} else if lastMessage.ChannelID != channelID {
				return errors.New("last read message not in the same channel")
			} else if threadID > 0 && lastMessage.ReplyTo != threadID {
				return errors.New("last read message not in the same thread")
			}

			count, err = svc.message.CountFromMessageID(channelID, threadID, lastReadMessageID)
			if err != nil {
				return errors.Wrap(err, "unable to count unread messages")
			}

		} else {
			// use last message ID
			if lastReadMessageID, err = svc.message.LastMessageID(channelID, threadID); err != nil {
				return errors.Wrap(err, "unable to find last message")
			}

			// no need to count
			count = 0
		}

		err = svc.unread.Record(currentUserID, channelID, threadID, lastReadMessageID, count)
		if err != nil {
			return errors.Wrap(err, "unable to record unread messages")
		}

		// Remove unread counts from all threads when doing mark-channel-as-read
		if threadID == 0 {
			err = svc.unread.ClearThreads(channelID, currentUserID)
			if err != nil {
				return errors.Wrap(err, "unable to clear channel threads")
			}
		}

		// Re-count unreads and send updates to this user
		svc.countUnreads(ch, nil, currentUserID)

		return nil
	})

	return lastReadMessageID, count, threadCount, errors.Wrap(err, "unable to mark as read")
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
func (svc message) flag(messageID uint64, flag string, remove bool) error {
	var currentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	_ = currentUserID

	if strings.TrimSpace(flag) == "" {
		// Sanitize
		flag = types.MessageFlagPinnedToChannel
	}

	// @todo validate flags beyond empty string

	err := svc.db.Transaction(func() (err error) {
		var flagOwnerId = currentUserID
		var f *types.MessageFlag
		var msg *types.Message
		var ch *types.Channel

		if flag == types.MessageFlagPinnedToChannel {
			// It does not matter how is the owner of the pin,
			flagOwnerId = 0
		}

		f, err = svc.mflag.FindByFlag(messageID, flagOwnerId, flag)
		if f == nil && remove {
			// Skip removing, flag does not exists
			return nil
		}

		if f != nil && f.ID > 0 && !remove {
			// Skip adding, flag already exists
			return nil
		}

		if err != nil && err != repository.ErrMessageFlagNotFound {
			// Other errors, exit
			return
		}

		if msg, err = svc.message.FindByID(messageID); err != nil {
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
			err = svc.mflag.DeleteByID(f.ID)
			f.DeletedAt = timeNowPtr()
		} else {
			f, err = svc.mflag.Create(&types.MessageFlag{
				UserID:    currentUserID,
				ChannelID: msg.ChannelID,
				MessageID: msg.ID,
				Flag:      flag,
			})
		}

		if err != nil {
			return
		}

		_ = svc.sendFlagEvent(f)
		return
	})

	return errors.Wrap(err, "can not flag/un-flag message")
}

func (svc message) preload(mm types.MessageSet) (err error) {
	if err = svc.preloadAttachments(mm); err != nil {
		return
	}

	if err = svc.preloadFlags(mm); err != nil {
		return
	}

	if err = svc.preloadMentions(mm); err != nil {
		return
	}

	if err = svc.preloadUnreads(mm); err != nil {
		return
	}

	if err = svc.message.PrefillThreadParticipants(mm); err != nil {
		return
	}

	return
}

// Preload for all messages
func (svc message) preloadFlags(mm types.MessageSet) (err error) {
	var ff types.MessageFlagSet

	ff, err = svc.mflag.FindByMessageIDs(mm.IDs()...)
	if err != nil {
		return
	}

	return ff.Walk(func(flag *types.MessageFlag) error {
		mm.FindByID(flag.MessageID).Flags = append(mm.FindByID(flag.MessageID).Flags, flag)
		return nil
	})
}

// Preload for all messages
func (svc message) preloadMentions(mm types.MessageSet) (err error) {
	var mentions types.MentionSet

	mentions, err = svc.mentions.FindByMessageIDs(mm.IDs()...)
	if err != nil {
		return
	}

	return mm.Walk(func(m *types.Message) error {
		m.Mentions = mentions.FindByMessageID(m.ID)
		return nil
	})
}

func (svc message) preloadAttachments(mm types.MessageSet) (err error) {
	var (
		ids []uint64
		aa  types.MessageAttachmentSet
	)

	err = mm.Walk(func(m *types.Message) error {
		if m.Type == types.MessageTypeAttachment || m.Type == types.MessageTypeInlineImage {
			ids = append(ids, m.ID)
		}
		return nil
	})

	if err != nil {
		return
	}

	if aa, err = svc.attachment.FindAttachmentByMessageID(ids...); err != nil {
		return
	} else {
		return aa.Walk(func(a *types.MessageAttachment) error {
			if a.MessageID > 0 {
				if m := mm.FindByID(a.MessageID); m != nil {
					m.Attachment = &a.Attachment
				}
			}

			return nil
		})
	}
}

func (svc message) preloadUnreads(mm types.MessageSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	// Filter out only relevant messages -- ones with replies
	mm, _ = mm.Filter(func(m *types.Message) (b bool, e error) {
		return m.Replies > 0, nil
	})

	if len(mm) == 0 {
		return nil
	}

	if vv, err := svc.unread.Count(userID, 0, mm.IDs()...); err != nil {
		return err
	} else {
		return mm.Walk(func(m *types.Message) error {
			m.Unread = vv.FindByThreadId(m.ID)
			return nil
		})
	}
}

// Sends message to event loop
func (svc message) sendEvent(mm ...*types.Message) (err error) {
	if err = svc.preload(mm); err != nil {
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
func (svc message) countUnreads(ch *types.Channel, m *types.Message, userID uint64) {
	var (
		err                           error
		uuBase, uuThreads, uuChannels types.UnreadSet
		// mm  types.ChannelMemberSet
		threadIDs []uint64
	)

	if m != nil {
		if m.DeletedAt != nil {
			// When deleting message, all existing counters are decreased!
			if err = svc.unread.Dec(m.ChannelID, m.ReplyTo, m.UserID); err != nil {
				svc.logger.With(zap.Error(err)).Info("could not decrement unread counter")
				return
			}
		} else if m.UpdatedAt == nil {
			// Reset user's counter and set current message ID as last read.
			err = svc.unread.Record(
				m.UserID,
				m.ChannelID,
				m.ReplyTo,
				m.ID,
				0,
			)

			// When new message is created, update all existing counters
			if err = svc.unread.Inc(m.ChannelID, m.ReplyTo, m.UserID); err != nil {
				svc.logger.With(zap.Error(err)).Info("could not increment unread counter")
				return
			}
		}

		if m.ReplyTo > 0 {
			threadIDs = []uint64{m.ReplyTo}
		}
	}

	uuBase, err = svc.unread.Count(userID, ch.ID, threadIDs...)
	if err != nil {
		svc.logger.With(zap.Error(err)).Info("could not count unread messages")
		return
	}

	if len(threadIDs) > 0 {
		// If base count was done for a thread,
		// Do another count for channel
		uuChannels, err = svc.unread.Count(userID, ch.ID)
		if err != nil {
			svc.logger.With(zap.Error(err)).Info("could not count unread messages")
			return
		}

		uuBase = uuBase.Merge(uuChannels)

		// Now recount all threads for this channel
		uuThreads, err = svc.unread.CountThreads(userID, ch.ID)
		if err != nil {
			svc.logger.With(zap.Error(err)).Info("could not count unread messages")
			return
		}

		uuBase = uuBase.Merge(uuThreads)
	}

	// This is a reply, make sure we fetch the new stats about unread replies and push them to users
	err = svc.event.UnreadCounters(uuBase)
	if err != nil {
		svc.logger.With(zap.Error(err)).Info("could not send unread count event")
		return
	}
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
		ChannelID:     m.ChannelID,
		MessageID:     m.ID,
		MentionedByID: m.UserID,
	}

	for m := 0; m < len(match); m++ {
		uid := payload.ParseUInt64(match[m][reSubID])
		if len(mm.FindByUserID(uid)) == 0 {
			// Copy template & assign user id
			mnt := tpl
			mnt.UserID = uid
			mm = append(mm, &mnt)
		}
	}

	return
}

func (svc message) updateMentions(messageID uint64, mm types.MentionSet) error {
	if existing, err := svc.mentions.FindByMessageIDs(messageID); err != nil {
		return errors.Wrap(err, "could not update mentions")
	} else if len(mm) > 0 {
		add, _, del := existing.Diff(mm)

		err = add.Walk(func(m *types.Mention) error {
			_, err = svc.mentions.Create(m)
			return err
		})

		if err != nil {
			return errors.Wrap(err, "could not create mentions")
		}

		err = del.Walk(func(m *types.Mention) error {
			return svc.mentions.DeleteByID(m.ID)
		})

		if err != nil {
			return errors.Wrap(err, "could not delete mentions")
		}
	} else {
		return svc.mentions.DeleteByMessageID(messageID)
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
	} else if mm, err = svc.cmember.Find(types.ChannelMemberFilterChannels(ch.ID)); err != nil {
		return nil, err
	} else {
		ch.Members = mm.AllMemberIDs()
		ch.Member = mm.FindByUserID(currentUserID)
	}

	return
}

var _ MessageService = &message{}
