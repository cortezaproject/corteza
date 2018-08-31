package service

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

type (
	message struct {
		rpo messageRepository
	}

	MessageService interface {
		Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error)

		Create(ctx context.Context, messages *types.Message) (*types.Message, error)
		Update(ctx context.Context, messages *types.Message) (*types.Message, error)

		React(ctx context.Context, messageID uint64, reaction string) error
		Unreact(ctx context.Context, messageID uint64, reaction string) error

		Pin(ctx context.Context, messageID uint64) error
		Unpin(ctx context.Context, messageID uint64) error

		Flag(ctx context.Context, messageID uint64) error
		Unflag(ctx context.Context, messageID uint64) error

		Attach(ctx context.Context) (*types.Attachment, error)
		Detach(ctx context.Context, messageID uint64) error

		Direct(ctx context.Context, recipientID uint64, in *types.Message) (out *types.Message, err error)

		deleter
	}

	messageRepository interface {
		repository.Transactionable
		repository.Message
		repository.Reaction
		repository.Attachment
		repository.Channel
	}
)

func Message() *message {
	m := &message{rpo: repository.New()}
	return m
}

func (svc message) Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & read from this channel
	_ = currentUserID
	_ = filter.ChannelID

	return svc.rpo.FindMessages(filter)
}

func (svc message) Direct(ctx context.Context, recipientID uint64, in *types.Message) (out *types.Message, err error) {
	return out, svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var currentUserID = auth.GetIdentityFromContext(ctx).Identity()

		// @todo [SECURITY] verify if current user can send direct messages to anyone?
		if false {
			return errors.New("Not allowed to send direct messages")
		}

		// @todo [SECURITY] verify if current user can send direct messages to this user
		if false {
			return errors.New("Not allowed to send direct messages to this user")
		}

		dch, err := r.FindDirectChannelByUserID(currentUserID, recipientID)
		if err == repository.ErrChannelNotFound {
			dch, err = r.CreateChannel(&types.Channel{
				Type: types.ChannelTypeGroup,
			})

			if err != nil {
				return
			}

			membership := &types.ChannelMember{ChannelID: dch.ID, Type: types.ChannelMembershipTypeOwner}

			membership.UserID = currentUserID
			if _, err = r.AddChannelMember(membership); err != nil {
				return
			}

			membership.UserID = recipientID
			if _, err = r.AddChannelMember(membership); err != nil {
				return
			}

		} else if err != nil {
			return errors.Wrap(err, "Could not send direct message")
		}

		// Make sure our message is sent to the right channel
		in.ChannelID = dch.ID
		in.UserID = currentUserID

		// @todo send new msg to the event-loop
		out, err = r.CreateMessage(in)
		return
	})
}

func (svc message) Create(ctx context.Context, mod *types.Message) (*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel

	mod.UserID = currentUserID

	return svc.rpo.CreateMessage(mod)
}

func (svc message) Update(ctx context.Context, mod *types.Message) (*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message

	// @todo verify ownership

	return svc.rpo.UpdateMessage(mod)
}

func (svc message) Delete(ctx context.Context, id uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message

	// @todo verify ownership

	return svc.rpo.DeleteMessageByID(id)
}

func (svc message) React(ctx context.Context, messageID uint64, reaction string) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	var m *types.Message

	// @todo validate reaction

	r := &types.Reaction{
		UserID:    currentUserID,
		MessageID: messageID,
		ChannelID: m.ChannelID,
		Reaction:  reaction,
	}

	if _, err := svc.rpo.CreateReaction(r); err != nil {
		return err
	}

	return nil
}

func (svc message) Unreact(ctx context.Context, messageID uint64, reaction string) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load reaction and verify ownership
	var r *types.Reaction

	return svc.rpo.DeleteReactionByID(r.ID)
}

func (svc message) Pin(ctx context.Context, messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc message) Unpin(ctx context.Context, messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc message) Flag(ctx context.Context, messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc message) Unflag(ctx context.Context, messageID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil
}

func (svc message) Attach(ctx context.Context) (*types.Attachment, error) {
	// @todo define func signature

	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return nil, nil
}

func (svc message) Detach(ctx context.Context, attachmentID uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo verify if current user can remove this attachment

	return svc.rpo.DeleteAttachmentByID(attachmentID)
}

var _ MessageService = &message{}
