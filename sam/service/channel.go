package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"

	authService "github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	channel struct {
		db  *factory.DB
		ctx context.Context

		usr authService.UserService
		evl EventService

		channel repository.ChannelRepository
		cmember repository.ChannelMemberRepository
		message repository.MessageRepository
	}

	ChannelService interface {
		With(ctx context.Context) ChannelService

		FindByID(channelID uint64) (*types.Channel, error)
		Find(filter *types.ChannelFilter) (types.ChannelSet, error)
		FindByMembership() (rval []*types.Channel, err error)
		FindMembers(channelID uint64) (types.ChannelMemberSet, error)

		Create(channel *types.Channel) (*types.Channel, error)
		Update(channel *types.Channel) (*types.Channel, error)

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
	}

	// channelSecurity interface {
	// 	CanRead(ch *types.Channel) bool
	// }
)

func Channel() ChannelService {
	return (&channel{
		usr: authService.DefaultUser,
		evl: DefaultEvent,
	}).With(context.Background())
}

func (svc *channel) With(ctx context.Context) ChannelService {
	db := repository.DB(ctx)
	return &channel{
		db:  db,
		ctx: ctx,

		usr: svc.usr.With(ctx),
		evl: svc.evl.With(ctx),

		channel: repository.Channel(ctx, db),
		cmember: repository.ChannelMember(ctx, db),
		message: repository.Message(ctx, db),
	}
}

func (svc *channel) FindByID(id uint64) (ch *types.Channel, err error) {
	ch, err = svc.channel.FindChannelByID(id)
	if err != nil {
		return
	}

	// if !svc.sec.ch.CanRead(ch) {
	// 	return nil, errors.New("Not allowed to access channel")
	// }

	return
}

func (svc *channel) Find(filter *types.ChannelFilter) (types.ChannelSet, error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if cc, err := svc.channel.FindChannels(filter); err != nil {
		return nil, err
	} else {
		return cc, svc.preloadMembers(cc)
	}
}

func (svc *channel) preloadMembers(cc types.ChannelSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if mm, err := svc.cmember.Find(&types.ChannelMemberFilter{ComembersOf: userID}); err != nil {
		return err
	} else {
		cc.Walk(func(ch *types.Channel) error {
			ch.Members = mm.MembersOf(ch.ID)
			return nil
		})
	}

	return nil
}

// FindMembers loads all members (and full users) for a specific channel
func (svc *channel) FindMembers(channelID uint64) (out types.ChannelMemberSet, err error) {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	// @todo [SECURITY] check if we can return members on this channel
	_ = channelID
	_ = userID

	return out, svc.db.Transaction(func() (err error) {
		out, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID})
		if err != nil {
			return err
		}

		if uu, err := svc.usr.Find(nil); err != nil {
			return err
		} else {
			return out.Walk(func(member *types.ChannelMember) error {
				member.User = uu.FindById(member.UserID)
				return nil
			})
		}
	})
}

// Returns all channels with membership info
func (svc *channel) FindByMembership() (rval []*types.Channel, err error) {
	return rval, svc.db.Transaction(func() error {
		var chMemberId = repository.Identity(svc.ctx)

		var mm []*types.ChannelMember

		if mm, err = svc.cmember.Find(&types.ChannelMemberFilter{MemberID: chMemberId}); err != nil {
			return err
		}

		if rval, err = svc.channel.FindChannels(nil); err != nil {
			return err
		}

		for _, m := range mm {
			for _, c := range rval {
				if c.ID == m.ChannelID {
					c.Member = m
				}
			}
		}

		return nil
	})
}

func (svc *channel) Create(in *types.Channel) (out *types.Channel, err error) {
	// @todo: [SECURITY] permission check if user can add channel

	return out, svc.db.Transaction(func() (err error) {
		var msg *types.Message

		// @todo get organisation from somewhere
		var organisationID uint64 = 0

		var chCreatorID = repository.Identity(svc.ctx)

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

		// @todo [SECURITY] check if user can create private channels
		if in.Type == types.ChannelTypeGroup && false {
			return errors.New("Not allowed to create group channels")
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
		if out, err = svc.channel.CreateChannel(out); err != nil {
			return
		}

		// Join current user as an member & owner
		_, err = svc.cmember.Create(&types.ChannelMember{
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
		msg, err = svc.message.CreateMessage(svc.makeSystemMessage(
			out,
			"@%d created new %s channel, topic is: %s",
			chCreatorID,
			"<PRIVATE-OR-PUBLIC>",
			"<TOPIC>"))

		_ = msg
		if err != nil {
			// Message creation failed
			return
		}

		svc.sendMessageEvent(msg)

		return svc.evl.Channel(out)
	})
}

func (svc *channel) Update(in *types.Channel) (out *types.Channel, err error) {
	return out, svc.db.Transaction(func() (err error) {
		var msgs types.MessageSet

		// @todo [SECURITY] can user access this channel?
		if out, err = svc.channel.FindChannelByID(in.ID); err != nil {
			return
		}

		if out.ArchivedAt != nil {
			return errors.New("Not allowed to edit archived channels")
		} else if out.DeletedAt != nil {
			return errors.New("Not allowed to edit deleted channels")
		}

		if out.Type != in.Type {
			// @todo [SECURITY] check if user can create public channels
			if in.Type == types.ChannelTypePublic && false {
				return errors.New("Not allowed to change type of this channel to public")
			}

			// @todo [SECURITY] check if user can create private channels
			if in.Type == types.ChannelTypePrivate && false {
				return errors.New("Not allowed to change type of this channel to private")
			}

			// @todo [SECURITY] check if user can create group channels
			if in.Type == types.ChannelTypeGroup && false {
				return errors.New("Not allowed to change type of this channel to group")
			}
		}

		var chUpdatorId = repository.Identity(svc.ctx)

		// Copy values
		if out.Name != in.Name {
			// @todo [SECURITY] can we change channel's name?
			if false {
				return errors.New("Not allowed to rename channel")
			} else {
				msgs = append(msgs, svc.makeSystemMessage(
					out, "@%d renamed channel %s (was: %s)", chUpdatorId, out.Name, in.Name))
			}
			out.Name = in.Name
		}

		if out.Topic != in.Topic && true {
			// @todo [SECURITY] can we change channel's topic?
			if false {
				return errors.New("Not allowed to change channel topic")
			} else {
				msgs = append(msgs, svc.makeSystemMessage(
					out, "@%d changed channel topic: %s (was: %s)", chUpdatorId, out.Topic, in.Topic))
			}

			out.Topic = in.Topic
		}

		// Save the updated channel
		if out, err = svc.channel.UpdateChannel(in); err != nil {
			return
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		for _, msg := range msgs {
			if msg, err = svc.message.CreateMessage(msg); err != nil {
				return err
			}

			svc.sendMessageEvent(msg)
		}

		if err != nil {
			// Message creation failed
			return
		}

		return svc.evl.Channel(out)
	})
}

func (svc *channel) Delete(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user delete this channel?

		if ch.DeletedAt != nil {
			return errors.New("Channel already deleted")
		}

		msg, err := svc.message.CreateMessage(svc.makeSystemMessage(ch, "@%d deleted this channel", userID))
		if err != nil {
			return
		}
		svc.sendMessageEvent(msg)

		if err = svc.channel.DeleteChannelByID(id); err != nil {
			return
		}

		return svc.evl.Channel(ch)
	})
}

func (svc *channel) Recover(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user recover this channel?

		if ch.DeletedAt == nil {
			return errors.New("Channel not deleted")
		}

		msg, err := svc.message.CreateMessage(svc.makeSystemMessage(ch, "@%d recovered this channel", userID))
		if err != nil {
			return
		}
		svc.sendMessageEvent(msg)

		err = svc.channel.UnarchiveChannelByID(id)
		if err != nil {
			return
		}

		return svc.evl.Channel(ch)

	})
}

func (svc *channel) Archive(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user archive this channel?

		if ch.ArchivedAt != nil {
			return errors.New("Channel already archived")
		}

		msg, err := svc.message.CreateMessage(svc.makeSystemMessage(ch, "@%d archived this channel", userID))
		if err != nil {
			return
		}
		svc.sendMessageEvent(msg)

		err = svc.channel.ArchiveChannelByID(id)
		if err != nil {
			return
		}

		return svc.evl.Channel(ch)
	})
}

func (svc *channel) Unarchive(id uint64) error {
	return svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)
		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(id); err != nil {
			return
		}

		// @todo [SECURITY] can user unarchive this channel?

		if ch.ArchivedAt == nil {
			return errors.New("Channel not archived")
		}

		msg, err := svc.message.CreateMessage(svc.makeSystemMessage(ch, "@%d unarchived this channel", userID))
		if err != nil {
			return
		}
		svc.sendMessageEvent(msg)

		err = svc.channel.ArchiveChannelByID(id)
		if err != nil {
			return
		}

		return svc.evl.Channel(ch)
	})
}

func (svc *channel) AddMember(m *types.ChannelMember) (out *types.ChannelMember, err error) {
	return out, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		var ch *types.Channel

		// @todo [SECURITY] can user access this channel?
		if ch, err = svc.channel.FindChannelByID(m.ChannelID); err != nil {
			return
		}

		// @todo [SECURITY] can user add members to this channel?

		msg, err := svc.message.CreateMessage(svc.makeSystemMessage(ch, "@%d added a new member to this channel: @%d", userID, m.UserID))
		if err != nil {
			return
		}
		svc.sendMessageEvent(msg)

		return err
	})
}

func (svc *channel) makeSystemMessage(ch *types.Channel, format string, a ...interface{}) *types.Message {
	return &types.Message{
		ChannelID: ch.ID,
		Message:   fmt.Sprintf(format, a...),
		Type:      types.MessageTypeChannelEvent,
	}
}

// Sends message to event loop
//
// It also preloads user
func (svc *channel) sendMessageEvent(msg *types.Message) (err error) {
	if msg.User == nil {
		// @todo pull user from cache
		if msg.User, err = svc.usr.FindByID(msg.UserID); err != nil {
			return
		}
	}

	return svc.evl.Message(msg)
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
//		FindMember(channelId uint64, userId uint64) (*types.User, error)
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
//func (sec nativeChannelSec) CanRead(ch *types.Channel) bool {
//	// @todo check if channel is public?
//
//	var currentUserID = repository.Identity(svc.ctx)
//
//	user, err := sec.rpo.FindMember(ch.ID, currentUserID)
//
//	return err != nil && user.Valid()
//}

var _ ChannelService = &channel{}
