package service

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	message struct {
		repository struct {
			message    messageRepository
			reaction   messageReactionRepository
			attachment messageAttachmentRepository
		}
	}

	messageRepository interface {
		FindByID(ctx context.Context, messageID uint64) (*types.Message, error)
		Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error)

		Create(ctx context.Context, message *types.Message) (*types.Message, error)
		Update(ctx context.Context, message *types.Message) (*types.Message, error)

		deleter
	}

	messageReactionRepository interface {
		FindByID(ctx context.Context, reactionID uint64) (*types.Reaction, error)
		Create(ctx context.Context, reaction *types.Reaction) (*types.Reaction, error)
		Delete(ctx context.Context, reactionID uint64) error
	}

	messageAttachmentRepository interface {
		FindByID(ctx context.Context, attachmentID uint64) (*types.Attachment, error)
		Create(ctx context.Context, attachment *types.Attachment) (*types.Attachment, error)
		Delete(ctx context.Context, attachmentID uint64) error
	}
)

func Message() *message {
	m := &message{}
	m.repository.message = repository.Message()
	m.repository.reaction = repository.Reaction()
	m.repository.attachment = repository.Attachment()

	return m
}

func (svc message) Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & read from this channel
	_ = currentUserID
	_ = filter.ChannelID

	return svc.repository.message.Find(ctx, filter)
}

func (svc message) Create(ctx context.Context, mod *types.Message) (*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	return svc.repository.message.Create(ctx, mod)
}

func (svc message) Update(ctx context.Context, mod *types.Message) (*types.Message, error) {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message

	// @todo verify ownership

	return svc.repository.message.Update(ctx, mod)
}

func (svc message) Delete(ctx context.Context, id uint64) error {
	// @todo get user from context
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access & write to this channel
	_ = currentUserID

	// @todo load current message

	// @todo verify ownership

	return svc.repository.message.Delete(ctx, id)
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

	if _, err := svc.repository.reaction.Create(ctx, r); err != nil {
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

	return svc.repository.reaction.Delete(ctx, r.ID)
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

	return svc.repository.attachment.Delete(ctx, attachmentID)
}
