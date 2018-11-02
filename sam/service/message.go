package service

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	authService "github.com/crusttech/crust/auth/service"
	authTypes "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	message struct {
		db  db
		ctx context.Context

		attachment repository.AttachmentRepository
		channel    repository.ChannelRepository
		cmember    repository.ChannelMemberRepository
		cview      repository.ChannelViewRepository
		message    repository.MessageRepository
		mflag      repository.MessageFlagRepository

		usr authService.UserService
		evl EventService
	}

	MessageService interface {
		With(ctx context.Context) MessageService

		Find(filter *types.MessageFilter) (types.MessageSet, error)
		FindThreads(filter *types.MessageFilter) (types.MessageSet, error)

		Create(messages *types.Message) (*types.Message, error)
		Update(messages *types.Message) (*types.Message, error)

		React(messageID uint64, reaction string) error
		RemoveReaction(messageID uint64, reaction string) error

		Pin(messageID uint64) error
		RemovePin(messageID uint64) error

		Bookmark(messageID uint64) error
		RemoveBookmark(messageID uint64) error

		Delete(ID uint64) error
	}
)

const (
	settingsMessageBodyLength = 0
)

func Message() MessageService {
	return &message{
		usr: authService.DefaultUser,
		evl: DefaultEvent,
	}
}

func (svc *message) With(ctx context.Context) MessageService {
	db := repository.DB(ctx)
	return &message{
		db:  db,
		ctx: ctx,

		usr: svc.usr.With(ctx),
		evl: svc.evl.With(ctx),

		attachment: repository.Attachment(ctx, db),
		channel:    repository.Channel(ctx, db),
		cmember:    repository.ChannelMember(ctx, db),
		cview:      repository.ChannelView(ctx, db),
		message:    repository.Message(ctx, db),
		mflag:      repository.MessageFlag(ctx, db),
	}
}

func (svc *message) Find(filter *types.MessageFilter) (mm types.MessageSet, err error) {
	// @todo get user from context
	filter.CurrentUserID = repository.Identity(svc.ctx)

	// @todo verify if current user can access & read from this channel
	_ = filter.ChannelID

	mm, err = svc.message.FindMessages(filter)
	if err != nil {
		return nil, err
	}

	return mm, svc.preload(mm)
}

func (svc *message) FindThreads(filter *types.MessageFilter) (mm types.MessageSet, err error) {
	// @todo get user from context
	filter.CurrentUserID = repository.Identity(svc.ctx)

	// @todo verify if current user can access & read from this channel
	_ = filter.ChannelID

	mm, err = svc.message.FindThreads(filter)
	if err != nil {
		return nil, err
	}

	return mm, svc.preload(mm)
}

func (svc *message) Create(in *types.Message) (message *types.Message, err error) {
	if in == nil {
		in = &types.Message{}
	}

	in.Message = strings.TrimSpace(in.Message)
	var mlen = len(in.Message)

	if mlen == 0 {
		return nil, errors.Errorf("Refusing to create message without contents")
	} else if settingsMessageBodyLength > 0 && mlen > settingsMessageBodyLength {
		return nil, errors.Errorf("Message length (%d characters) too long (max: %d)", mlen, settingsMessageBodyLength)
	}

	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	in.UserID = currentUserID

	return message, svc.db.Transaction(func() (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}

		if in.ReplyTo > 0 {
			var original *types.Message
			var replyTo = in.ReplyTo

			for replyTo > 0 {
				// Find original message
				original, err = svc.message.FindMessageByID(in.ReplyTo)
				if err != nil {
					return
				}

				replyTo = original.ReplyTo
			}

			if !original.Type.IsRepliable() {
				return errors.Errorf("Unable to reply on this message (type = %s)", original.Type)
			}

			// We do not want to have multi-level threads
			// Take original's reply-to and use it
			in.ReplyTo = original.ID

			in.ChannelID = original.ChannelID

			// Increment counter, on struct and in repostiry.
			original.Replies++
			if err = svc.message.IncReplyCount(original.ID); err != nil {
				return
			}

			// Broadcast updated original
			bq = append(bq, original)
		}

		if in.ChannelID == 0 {
			return errors.New("ChannelID missing")
		}

		// @todo [SECURITY] verify if current user can access & write to this channel

		if message, err = svc.message.CreateMessage(in); err != nil {
			return
		}

		if err = svc.cview.Inc(message.ChannelID, message.UserID); err != nil {
			return
		}

		return svc.sendEvent(append(bq, message)...)
	})
}

func (svc *message) Update(in *types.Message) (message *types.Message, err error) {
	if in == nil {
		in = &types.Message{}
	}

	in.Message = strings.TrimSpace(in.Message)
	var mlen = len(in.Message)

	if mlen == 0 {
		return nil, errors.Errorf("Refusing to update message without contents")
	} else if settingsMessageBodyLength > 0 && mlen > settingsMessageBodyLength {
		return nil, errors.Errorf("Message length (%d characters) too long (max: %d)", mlen, settingsMessageBodyLength)
	}

	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return message, svc.db.Transaction(func() (err error) {
		original, err := svc.message.FindMessageByID(in.ID)
		if err != nil {
			return err
		}

		if original.Message == in.Message {
			// Nothing changed
			return nil
		}

		if original.UserID != currentUserID {
			return errors.New("Not an owner")
		}

		// Allow message content to be changed, ignore everything else
		original.Message = in.Message

		if message, err = svc.message.UpdateMessage(original); err != nil {
			return err
		}

		return svc.sendEvent(original)
	})
}

func (svc *message) Delete(ID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message
	// @todo verify ownership

	return svc.db.Transaction(func() (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}
		var deletedMsg, original *types.Message

		deletedMsg, err = svc.message.FindMessageByID(ID)
		if err != nil {
			return err
		}

		if deletedMsg.ReplyTo > 0 {
			original, err = svc.message.FindMessageByID(deletedMsg.ReplyTo)
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

		if err = svc.message.DeleteMessageByID(ID); err != nil {
			return
		}

		if err = svc.cview.Dec(deletedMsg.ChannelID, deletedMsg.UserID); err != nil {
			return err
		} else {
			// Set deletedAt timestamp so that our clients can react properly...
			deletedMsg.DeletedAt = timeNowPtr()
		}

		return svc.sendEvent(append(bq, deletedMsg)...)
	})
}

// React on a message with an emoji
func (svc *message) React(messageID uint64, reaction string) error {
	return svc.flag(messageID, reaction, false)
}

// Remove reaction on a message
func (svc *message) RemoveReaction(messageID uint64, reaction string) error {
	return svc.flag(messageID, reaction, true)
}

// Pin message to the channel
func (svc *message) Pin(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagPinnedToChannel, false)
}

// Remove pin from message
func (svc *message) RemovePin(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagPinnedToChannel, true)
}

// Bookmark message (private)
func (svc *message) Bookmark(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagBookmarkedMessage, false)
}

// Remove bookmark message (private)
func (svc *message) RemoveBookmark(messageID uint64) error {
	return svc.flag(messageID, types.MessageFlagBookmarkedMessage, true)
}

// React on a message with an emoji
func (svc *message) flag(messageID uint64, flag string, remove bool) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	if strings.TrimSpace(flag) == "" {
		// Sanitize
		flag = types.MessageFlagPinnedToChannel
	}

	// @todo validate flags beyond empty string

	err := svc.db.Transaction(func() (err error) {
		var flagOwnerId = currentUserID
		var f *types.MessageFlag

		// @todo [SECURITY] verify if current user can access & write to this channel

		if flag == types.MessageFlagPinnedToChannel {
			// It does not matter how is the owner of the pin,
			flagOwnerId = 0
		}

		f, err = svc.mflag.FindByFlag(messageID, flagOwnerId, flag)
		if f.ID == 0 && remove {
			// Skip removing, flag does not exists
			return nil
		} else if f.ID > 0 && !remove {
			// Skip adding, flag already exists
			return nil
		} else if err != nil && err != repository.ErrMessageFlagNotFound {
			// Other errors, exit
			return
		}

		// Check message
		var msg *types.Message
		msg, err = svc.message.FindMessageByID(messageID)
		if err != nil {
			return
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
			return err
		}

		svc.sendFlagEvent(f)

		return
	})

	return errors.Wrap(err, "Can not flag/un-flag message")
}

func (svc *message) preload(mm types.MessageSet) (err error) {
	if err = svc.preloadUsers(mm); err != nil {
		return
	}

	if err = svc.preloadAttachments(mm); err != nil {
		return
	}

	if err = svc.preloadFlags(mm); err != nil {
		return
	}

	return
}

// Preload for all messages
func (svc *message) preloadUsers(mm types.MessageSet) (err error) {
	var uu authTypes.UserSet

	for _, msg := range mm {
		if msg.User != nil || msg.UserID == 0 {
			continue
		}

		if msg.User = uu.FindById(msg.UserID); msg.User != nil {
			continue
		}

		if msg.User, _ = svc.usr.FindByID(msg.UserID); msg.User != nil {
			// @todo fix this handler errors (ignore user-not-found, return others)
			uu = append(uu, msg.User)
		}
	}

	return
}

// Preload for all messages
func (svc *message) preloadFlags(mm types.MessageSet) (err error) {
	var ff types.MessageFlagSet

	ff, err = svc.mflag.FindByMessageIDs(mm.IDs()...)
	if err != nil {
		return
	}

	return ff.Walk(func(flag *types.MessageFlag) error {
		mm.FindById(flag.MessageID).Flags = append(mm.FindById(flag.MessageID).Flags, flag)
		return nil
	})
}

func (svc *message) preloadAttachments(mm types.MessageSet) (err error) {
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
				if m := mm.FindById(a.MessageID); m != nil {
					m.Attachment = &a.Attachment
				}
			}

			return nil
		})
	}
}

// Sends message to event loop
func (svc *message) sendEvent(mm ...*types.Message) (err error) {
	if err = svc.preload(mm); err != nil {
		return
	}

	for _, msg := range mm {
		if msg.User == nil {
			// @todo fix this handler errors (ignore user-not-found, return others)
			msg.User, _ = svc.usr.FindByID(msg.UserID)
		}

		if err = svc.evl.Message(msg); err != nil {
			return
		}
	}

	return
}

// Sends message to event loop
func (svc *message) sendFlagEvent(ff ...*types.MessageFlag) (err error) {
	for _, f := range ff {
		if err = svc.evl.MessageFlag(f); err != nil {
			return
		}
	}

	return
}

var _ MessageService = &message{}
