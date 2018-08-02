package service

import (
	"context"
	"fmt"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

type (
	channel struct {
		rpo channelRepository
		//
		//sec struct {
		//	ch channelSecurity
		//}
	}

	channelRepository interface {
		repository.Transactionable
		repository.Channel
	}

	//channelSecurity interface {
	//	CanRead(ctx context.Context, ch *types.Channel) bool
	//}
)

func Channel() *channel {
	var svc = &channel{}

	svc.rpo = repository.New()
	//svc.sec.ch = ChannelSecurity(svc.rpo)

	return svc
}

func (svc channel) FindByID(ctx context.Context, id uint64) (ch *types.Channel, err error) {
	ch, err = svc.rpo.FindChannelByID(id)
	if err != nil {
		return
	}

	//if !svc.sec.ch.CanRead(ch) {
	//	return nil, errors.New("Not allowed to access channel")
	//}

	return
}

func (svc channel) Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error) {
	// @todo: permission check to return only channels that channel has access to
	// @todo: actual searching not just a full select
	return svc.rpo.FindChannels(filter)
}

func (svc channel) Create(ctx context.Context, in *types.Channel) (out *types.Channel, err error) {
	// @todo: [SECURITY] permission check if user can add channel

	return out, svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var msg *types.Message

		if in.Type == types.ChannelTypeDirect {
			return errors.New("Not allowed to create direct channels")
		}

		// @todo get organisation from somewhere
		var organisationID uint64 = 0

		var chCreatorID = auth.GetIdentityFromContext(ctx).Identity()

		// @todo [SECURITY] check if channel topic can be set
		if in.Topic != "" && false {
			return errors.New("Not allowed to set channel topic")
		}

		// @todo [SECURITY] check if user can create public channels
		if in.Type == types.ChannelTypePublic && false {
			return errors.New("Not allowed to create public channels")
		}

		// @todo [SECURITY] check if user can create private channels
		if in.Type == types.ChannelTypePrivate && false {
			return errors.New("Not allowed to create public channels")
		}

		// This is a fresh channel, just copy values
		out = &types.Channel{
			Name:           in.Name,
			Topic:          in.Topic,
			Type:           in.Type,
			OrganisationID: organisationID,
			CreatorID:      chCreatorID,
		}

		// Save the channel
		if out, err = r.CreateChannel(out); err != nil {
			return
		}

		// Join current user as an member & owner
		_, err = r.AddChannelMember(&types.ChannelMember{
			ChannelID: out.ID,
			UserID:    chCreatorID,
			Type:      types.ChannelMembershipTypeOwner,
		})

		if err != nil {
			// Could not add member
			return
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		msg, err = r.CreateMessage(svc.makeSystemMessage(
			out,
			"%s created new %s channel, topic is: %s",
			"<USERNAME>",
			"<PRIVATE-OR-PUBLIC>",
			"<TOPIC>"))

		if err != nil {
			// Message creation failed
			return
		}

		// @todo send channel creation to the event-loop
		// @todo send msg to the event-loop
		_ = msg

		return nil
	})
}

func (svc channel) Update(ctx context.Context, in *types.Channel) (out *types.Channel, err error) {
	return out, svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var msgs []*types.Message

		// @todo [SECURITY] can user access this channel?
		if out, err = r.FindChannelByID(in.ID); err != nil {
			return
		}

		if out.Type != types.ChannelTypeDirect {
			return errors.New("Not allowed to change direct channels")
		}

		if out.ArchivedAt != nil {
			return errors.New("Not allowed to edit archived channels")
		} else if out.DeletedAt != nil {
			return errors.New("Not allowed to edit deleted channels")
		}

		// var chCreatorID = auth.GetIdentityFromContext(ctx).Identity()

		// Copy values
		if out.Name != in.Name {
			// @todo [SECURITY] can we change channel's name?
			if false {
				return errors.New("Not allowed to rename channel")
			} else {
				msgs = append(msgs, svc.makeSystemMessage(
					out, "%s renamed channel %s (was: %s)", "<USERNAME>", out.Name, in.Name))
			}
			out.Name = in.Name
		}

		if out.Topic != in.Topic && true {
			// @todo [SECURITY] can we change channel's topic?
			if false {
				return errors.New("Not allowed to change channel topic")
			} else {
				msgs = append(msgs, svc.makeSystemMessage(
					out, "%s changed channel topic: %s (was: %s)", "<USERNAME>", out.Topic, in.Topic))
			}

			out.Topic = in.Topic
		}

		if out.Type != in.Type {
			// @todo [SECURITY] check if user can create public channels
			if in.Type == types.ChannelTypePublic && false {
				return errors.New("Not allowed to make this channel public")
			}

			// @todo [SECURITY] check if user can create private channels
			if in.Type == types.ChannelTypePrivate && false {
				return errors.New("Not allowed to make this channel private")
			}
		}

		// Save the updated channel
		if out, err = r.UpdateChannel(in); err != nil {
			return
		}

		// @todo send channel creation to the event-loop

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		for _, msg := range msgs {
			if msg, err = r.CreateMessage(msg); err != nil {
				// @todo send new msg to the event-loop
				return err
			}
		}

		if err != nil {
			// Message creation failed
			return
		}

		return nil
	})
}

func (svc channel) Delete(ctx context.Context, ID uint64) error {
	return svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var ch *types.Channel

		// @todo [SECURITY] can user delete this channel?

		if ch.DeletedAt != nil {
			return errors.New("Channel already deleted")
		}

		_, err = r.CreateMessage(svc.makeSystemMessage(ch,
			"%s deleted this channel"))

		return r.DeleteChannelByID(ID)
	})
}

func (svc channel) Recover(ctx context.Context, ID uint64) error {
	return svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var ch *types.Channel

		// @todo [SECURITY] can user recover this channel?

		if ch.DeletedAt == nil {
			return errors.New("Channel not deleted")
		}

		_, err = r.CreateMessage(svc.makeSystemMessage(ch,
			"%s recovered this channel"))

		return r.DeleteChannelByID(ID)
	})
}

func (svc channel) Archive(ctx context.Context, ID uint64) error {
	return svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var ch *types.Channel

		// @todo [SECURITY] can user archive this channel?

		if ch.ArchivedAt != nil {
			return errors.New("Channel already archived")
		}

		_, err = r.CreateMessage(svc.makeSystemMessage(ch,
			"%s archived this channel"))

		return r.ArchiveChannelByID(ID)
	})
}

func (svc channel) Unarchive(ctx context.Context, ID uint64) error {
	return svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {
		var ch *types.Channel

		// @todo [SECURITY] can user unarchive this channel?

		if ch.ArchivedAt == nil {
			return errors.New("Channel not archived")
		}

		_, err = r.CreateMessage(svc.makeSystemMessage(ch,
			"%s unarchived this channel"))

		return r.ArchiveChannelByID(ID)
	})

}

func (svc channel) makeSystemMessage(ch *types.Channel, format string, a ...interface{}) *types.Message {
	return &types.Message{
		ChannelID: ch.ID,
		Message:   fmt.Sprintf(format, a...),
	}
}

//// @todo temp location, move this somewhere else
//type (
//	nativeChannelSec struct {
//		rpo struct {
//			ch nativeChannelSecChRepo
//		}
//	}
//
//	nativeChannelSecChRepo interface {
//		FindMember(ctx context.Context, channelId uint64, userId uint64) (*types.User, error)
//	}
//)
//
//func ChannelSecurity(chRpo nativeChannelSecChRepo) channelSecurity {
//	var sec = &nativeChannelSec{}
//
//	sec.rpo.ch = chRpo
//
//	return sec
//}
//
//// Current user can read the channel if he is a member
//func (sec nativeChannelSec) CanRead(ctx context.Context, ch *types.Channel) bool {
//	// @todo check if channel is public?
//
//	var currentUserID = auth.GetIdentityFromContext(ctx).Identity()
//
//	user, err := sec.rpo.FindMember(ch.ID, currentUserID)
//
//	return err != nil && user.Valid()
//}
