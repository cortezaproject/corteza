package service

import (
	"context"

	authService "github.com/crusttech/crust/auth/service"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

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

		Direct(recipientID uint64, in *types.Message) (out *types.Message, err error)

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

	mm.Walk(func(i *types.Message) (err error) {
		i.User, err = svc.usr.FindByID(i.UserID)
		return
	})

	return mm, svc.loadAttachments(mm)
}

func (svc *message) loadAttachments(mm types.MessageSet) (err error) {
	var ids []uint64
	mm.Walk(func(m *types.Message) error {
		if m.Type == types.MessageTypeAttachment || m.Type == types.MessageTypeInlineImage {
			ids = append(ids, m.ID)
		}
		return nil
	})

	if set, err := svc.attachment.FindAttachmentByMessageID(ids...); err != nil {
		return err
	} else {
		return set.Walk(func(a *types.MessageAttachment) error {
			if a.MessageID > 0 {
				if m := mm.FindById(a.MessageID); m != nil {
					m.Attachment = &a.Attachment
				}
			}

			return nil
		})
	}
}

func (svc *message) Direct(recipientID uint64, in *types.Message) (out *types.Message, err error) {
	return out, svc.db.Transaction(func() (err error) {
		var currentUserID = repository.Identity(svc.ctx)

		// @todo [SECURITY] verify if current user can send direct messages to anyone?
		if false {
			return errors.New("Not allowed to send direct messages")
		}

		// @todo [SECURITY] verify if current user can send direct messages to this user
		if false {
			return errors.New("Not allowed to send direct messages to this user")
		}

		dch, err := svc.channel.FindDirectChannelByUserID(currentUserID, recipientID)
		if err == repository.ErrChannelNotFound {
			dch, err = svc.channel.CreateChannel(&types.Channel{
				Type: types.ChannelTypeGroup,
			})

			if err != nil {
				return
			}

			membership := &types.ChannelMember{ChannelID: dch.ID, Type: types.ChannelMembershipTypeOwner}

			membership.UserID = currentUserID
			if _, err = svc.cmember.Create(membership); err != nil {
				return
			}

			membership.UserID = recipientID
			if _, err = svc.cmember.Create(membership); err != nil {
				return
			}

		} else if err != nil {
			return errors.Wrap(err, "Could not send direct message")
		}

		// Make sure our message is sent to the right channel
		in.ChannelID = dch.ID
		in.UserID = currentUserID
		in.Type = types.MessageTypeSimpleMessage

		// @todo send new msg to the event-loop
		out, err = svc.message.CreateMessage(in)
		if err != nil {
			return
		}

		if err = svc.cview.Inc(in.ChannelID, in.UserID); err != nil {
			return err
		}

		return svc.sendEvent(out)
	})
}

func (svc *message) Create(mod *types.Message) (message *types.Message, err error) {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel

	mod.UserID = currentUserID

	return message, svc.db.Transaction(func() (err error) {
		if message, err = svc.message.CreateMessage(mod); err != nil {
			return err
		}

		if err = svc.cview.Inc(message.ChannelID, message.UserID); err != nil {
			return err
		}

		return svc.sendEvent(message)
	})
}

func (svc *message) Update(mod *types.Message) (*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message

	// @todo verify ownership

	message, err := svc.message.UpdateMessage(mod)

	if err == nil {
		return nil, err
	}

	return message, svc.sendEvent(message)
}

func (svc *message) Delete(id uint64) error {
	// @todo get user from context
	var currentUserID uint64 = repository.Identity(svc.ctx)

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message
	// @todo verify ownership

	return svc.db.Transaction(func() (err error) {
		if err = svc.message.DeleteMessageByID(id); err != nil {
			return
		}

		// if err == nil {
		// 	err = svc.evl.Message(message)
		// }

		// if err = svc.cview.Dec(message.ChannelID, message.UserID); err != nil {
		// 	return err
		// }

		return nil
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
func (svc *message) sendEvent(msg *types.Message) (err error) {
	if msg.User == nil {
		// @todo pull user from cache
		if msg.User, err = svc.usr.FindByID(msg.UserID); err != nil {
			return
		}
	}

	return svc.evl.Message(msg)
}

var _ MessageService = &message{}
