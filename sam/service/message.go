package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	authService "github.com/crusttech/crust/auth/service"

	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	message struct {
		db  *factory.DB
		ctx context.Context

		attachment repository.AttachmentRepository
		channel    repository.ChannelRepository
		cmember    repository.ChannelMemberRepository
		cview      repository.ChannelViewRepository
		message    repository.MessageRepository
		reaction   repository.ReactionRepository

		usr authService.UserService
		evl EventService
	}

	MessageService interface {
		With(ctx context.Context) MessageService

		Find(filter *types.MessageFilter) (types.MessageSet, error)

		Create(messages *types.Message) (*types.Message, error)
		Update(messages *types.Message) (*types.Message, error)

		React(messageID uint64, reaction string) error
		Unreact(messageID uint64, reaction string) error

		Pin(messageID uint64) error
		Unpin(messageID uint64) error

		Flag(messageID uint64) error
		Unflag(messageID uint64) error

		Delete(ID uint64) error
	}
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
		reaction:   repository.Reaction(ctx, db),
	}
}

func (svc *message) Find(filter *types.MessageFilter) (mm types.MessageSet, err error) {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & read from this channel
	_ = currentUserID
	_ = filter.ChannelID

	mm, err = svc.message.FindMessages(filter)
	if err != nil {
		return nil, err
	}

	_ = mm.Walk(func(i *types.Message) (err error) {
		// @todo fix this handler errors (ignore user-not-found, return others)
		i.User, err = svc.usr.FindByID(i.UserID)
		return
	})

	return mm, svc.loadAttachments(mm)
}

func (svc *message) loadAttachments(mm types.MessageSet) (err error) {
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

func (svc *message) Create(mod *types.Message) (message *types.Message, err error) {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	mod.UserID = currentUserID

	return message, svc.db.Transaction(func() (err error) {
		// Broadcast queue
		var bq = types.MessageSet{}

		if mod.ReplyTo > 0 {
			var original *types.Message
			original, err = svc.message.FindMessageByID(mod.ReplyTo)
			if err != nil {
				return
			}

			if original.ReplyTo > 0 {
				// We do not want to have multi-level threads
				// Take original's reply-to and use it
				mod.ReplyTo = original.ReplyTo
			}

			mod.ChannelID = original.ChannelID

			// Increment counter, on struct and in repostiry.
			original.Replies++
			if err = svc.message.IncReplyCount(original.ID); err != nil {
				return
			}

			// Broadcast updated original
			svc.sendEvent(original)
			bq = append(bq, original)
		}

		if mod.ChannelID == 0 {
			return errors.New("ChannelID missing")
		}

		// @todo [SECURITY] verify if current user can access & write to this channel

		if message, err = svc.message.CreateMessage(mod); err != nil {
			return
		}

		if err = svc.cview.Inc(message.ChannelID, message.UserID); err != nil {
			return
		}

		return svc.sendEvent(append(bq, message)...)
	})
}

func (svc *message) Update(mod *types.Message) (message *types.Message, err error) {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return message, svc.db.Transaction(func() (err error) {
		original, err := svc.message.FindMessageByID(mod.ID)
		if err != nil {
			return err
		}

		if original.UserID != currentUserID {
			return errors.New("Not an owner")
		}

		// Allow message content to be changed, ignore everything else
		original.Message = mod.Message

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

		// Set deletedAt timestamp so that our clients can react properly...
		now := time.Now()
		deletedMsg.DeletedAt = &now

		if err = svc.cview.Dec(deletedMsg.ChannelID, deletedMsg.UserID); err != nil {
			return err
		}

		return svc.sendEvent(append(bq, deletedMsg)...)
	})
}

func (svc *message) React(messageID uint64, reaction string) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	var m *types.Message

	// @todo validate reaction

	r := &types.Reaction{
		UserID:    currentUserID,
		MessageID: messageID,
		ChannelID: m.ChannelID,
		Reaction:  reaction,
	}

	if _, err := svc.reaction.CreateReaction(r); err != nil {
		return err
	}

	return nil
}

func (svc *message) Unreact(messageID uint64, reaction string) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load reaction and verify ownership
	var r *types.Reaction

	return svc.reaction.DeleteReactionByID(r.ID)
}

func (svc *message) Pin(messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc *message) Unpin(messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc *message) Flag(messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc *message) Unflag(messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

// Sends message to event loop
//
// It also preloads user
func (svc *message) sendEvent(mm ...*types.Message) (err error) {
	svc.loadAttachments(mm)

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

var _ MessageService = &message{}
