package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/organization"
	"github.com/crusttech/crust/messaging/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	channel struct {
		db  db
		ctx context.Context

		usr systemService.UserService

		evl EventService
		prm PermissionsService

		channel repository.ChannelRepository
		cmember repository.ChannelMemberRepository
		unread  repository.UnreadRepository
		message repository.MessageRepository

		sysmsgs types.MessageSet
	}

	ChannelService interface {
		With(ctx context.Context) ChannelService

		FindByID(channelID uint64) (*types.Channel, error)
		Find(filter *types.ChannelFilter) (types.ChannelSet, error)

		Create(channel *types.Channel) (*types.Channel, error)
		Update(channel *types.Channel) (*types.Channel, error)

		FindMembers(channelID uint64) (types.ChannelMemberSet, error)

		InviteUser(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		DeleteMember(channelID uint64, memberIDs ...uint64) (err error)

		SetFlag(ID uint64, flag types.ChannelMembershipFlag) (*types.Channel, error)

		Archive(ID uint64) (*types.Channel, error)
		Unarchive(ID uint64) (*types.Channel, error)
		Delete(ID uint64) (*types.Channel, error)
		Undelete(ID uint64) (*types.Channel, error)
		RecordView(userID, channelID, lastMessageID uint64) error
	}
)

var (
	ErrUnknownChannelType = errors.New("Unknown ChannelType")
	ErrNoPermission       = errors.New("No permissions")
)

const (
	settingsChannelNameLength  = 40
	settingsChannelTopicLength = 200
)

func Channel() ChannelService {
	return (&channel{
		usr: systemService.DefaultUser,
		evl: DefaultEvent,
		prm: DefaultPermissions,
	}).With(context.Background())
}

func (svc *channel) With(ctx context.Context) ChannelService {
	db := repository.DB(ctx)
	return &channel{
		db:  db,
		ctx: ctx,

		usr: svc.usr.With(ctx),
		evl: svc.evl.With(ctx),
		prm: svc.prm.With(ctx),

		channel: repository.Channel(ctx, db),
		cmember: repository.ChannelMember(ctx, db),
		unread:  repository.Unread(ctx, db),
		message: repository.Message(ctx, db),

		// System messages should be flushed at the end of each session
		sysmsgs: types.MessageSet{},
	}
}

func (svc *channel) FindByID(ID uint64) (ch *types.Channel, err error) {
	ch, err = svc.channel.FindChannelByID(ID)
	if err != nil {
		return
	} else if err = svc.preloadExtras(types.ChannelSet{ch}); err != nil {
		return
	}

	if !ch.CanObserve {
		return nil, errors.New("Not allowed to access channel")
	}

	return
}

func (svc *channel) Find(filter *types.ChannelFilter) (cc types.ChannelSet, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return cc, svc.db.Transaction(func() (err error) {
		if cc, err = svc.channel.FindChannels(filter); err != nil {
			return
		} else if err = svc.preloadExtras(cc); err != nil {
			return
		}

		cc, err = cc.Filter(func(c *types.Channel) (b bool, e error) {
			return c.CanObserve, nil
		})

		return
	})
}

// preloadExtras pre-loads channel's members, views
func (svc *channel) preloadExtras(cc types.ChannelSet) (err error) {
	if err = svc.preloadMembers(cc); err != nil {
		return
	}

	if err = cc.Walk(svc.setPermissionFlags); err != nil {
		return err
	}

	if err = svc.preloadViews(cc); err != nil {
		return
	}

	return
}

func (svc *channel) preloadMembers(cc types.ChannelSet) (err error) {
	var (
		userID = auth.GetIdentityFromContext(svc.ctx).Identity()
		mm     types.ChannelMemberSet
	)

	if mm, err = svc.cmember.Find(&types.ChannelMemberFilter{ComembersOf: userID}); err != nil {
		return
	} else {
		err = cc.Walk(func(ch *types.Channel) error {
			ch.Members = mm.MembersOf(ch.ID)
			ch.Member = mm.FindByChannelID(ch.ID).FindByUserID(userID)
			return nil
		})
	}

	return
}

func (svc *channel) preloadViews(cc types.ChannelSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if vv, err := svc.unread.Find(&types.UnreadFilter{UserID: userID}); err != nil {
		return err
	} else {
		cc.Walk(func(ch *types.Channel) error {
			ch.Unread = vv.FindByChannelId(ch.ID)
			return nil
		})
	}

	return nil
}

// FindMembers loads all members (and full users) for a specific channel
func (svc *channel) FindMembers(channelID uint64) (out types.ChannelMemberSet, err error) {
	return out, svc.db.Transaction(func() (err error) {
		if _, err = svc.FindByID(channelID); err != nil {
			return
		}

		out, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID})
		if err != nil {
			return err
		}

		if uu, err := svc.usr.Find(nil); err != nil {
			return err
		} else {
			return out.Walk(func(member *types.ChannelMember) error {
				member.User = uu.FindByID(member.UserID)
				return nil
			})
		}
	})
}

func (svc *channel) Create(in *types.Channel) (out *types.Channel, err error) {
	return out, svc.db.Transaction(func() (err error) {
		var msg *types.Message

		var organisationID = organization.Crust().ID

		var chCreatorID = repository.Identity(svc.ctx)

		mm := svc.buildMemberSet(chCreatorID, in.Members...)

		if in.Type == types.ChannelTypeGroup {
			if out, err = svc.checkGroupExistance(mm); err != nil {
				return err
			} else if out != nil && out.CanObserve {
				// Group already exists so let's just return it
				return nil
			} else if out != nil && !out.CanObserve {
				return errors.New("Not allowed to create this channel due to permission settings")
			}
		}

		if in.Topic != "" && false {
			return errors.New("Not allowed to set channel topic")
		}

		if in.Type == types.ChannelTypePublic && !svc.prm.CanCreatePublicChannel() {
			return errors.New("Not allowed to create public channels")
		}

		if in.Type == types.ChannelTypePrivate && !svc.prm.CanCreatePrivateChannel() {
			return errors.New("Not allowed to create private channels")
		}

		if in.Type == types.ChannelTypeGroup && !svc.prm.CanCreateDirectChannel() {
			return errors.New("Not allowed to create group channels")
		}

		if len(in.Name) == 0 && in.Type != types.ChannelTypeGroup {
			return errors.New("Channel name not provided")
		}

		if settingsChannelNameLength > 0 && len(in.Name) > settingsChannelNameLength {
			return errors.Errorf("Channel name (%d characters) too long (max: %d)", len(in.Name), settingsChannelNameLength)
		}

		if len(in.Topic) > 0 && settingsChannelTopicLength > 0 && len(in.Topic) > settingsChannelTopicLength {
			return errors.Errorf("Channel topic (%d characters) too long (max: %d)", len(in.Topic), settingsChannelTopicLength)
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

		err = mm.Walk(func(m *types.ChannelMember) (err error) {
			// Assign channel ID to membership
			m.ChannelID = out.ID

			// Create member
			if m, err = svc.cmember.Create(m); err != nil {
				return err
			}

			// Subscribe all members
			return svc.evl.Join(m.UserID, out.ID)
		})

		if err != nil {
			// Could not add member
			return
		}

		// Copy all member IDs to channel's member slice
		out.Members = mm.AllMemberIDs()

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if len(out.Name) == 0 {
			svc.scheduleSystemMessage(out, `<@%d> created %s channel`, chCreatorID, out.Type)
		} else if len(out.Topic) == 0 {
			svc.scheduleSystemMessage(out, `<@%d> created %s channel **%s**`, chCreatorID, out.Type, out.Name)
		} else {
			svc.scheduleSystemMessage(out, `<@%d> created %s channel **%s**, topic: %s`, chCreatorID, out.Type, out.Name, out.Topic)
		}

		_ = msg
		if err != nil {
			// Message creation failed
			return
		}

		svc.flushSystemMessages()

		return svc.sendChannelEvent(out)
	})
}

func (svc *channel) buildMemberSet(owner uint64, members ...uint64) (mm types.ChannelMemberSet) {
	// Join current user as an member & owner
	mm = types.ChannelMemberSet{&types.ChannelMember{
		UserID: owner,
		Type:   types.ChannelMembershipTypeOwner,
	}}

	// Add all required members and make sure that list is unique
	for _, m := range members {
		if mm.FindByUserID(m) == nil {
			mm = append(mm, &types.ChannelMember{
				UserID: m,
				Type:   types.ChannelMembershipTypeMember,
			})
		}
	}

	return mm
}

func (svc *channel) checkGroupExistance(mm types.ChannelMemberSet) (out *types.Channel, err error) {
	if out, err = svc.channel.FindChannelByMemberSet(mm.AllMemberIDs()...); err == repository.ErrChannelNotFound {
		return nil, nil
	} else if out != nil && err == nil {
		err = svc.preloadExtras(types.ChannelSet{out})
	}

	return
}

func (svc *channel) Update(in *types.Channel) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var changed bool

		if ch, err = svc.FindByID(in.ID); err != nil {
			return
		}

		if !ch.CanUpdate {
			return errors.New("Not allowed to update this channel")
		}

		if in.Type.IsValid() && ch.Type != in.Type {
			// @todo [SECURITY] check if user can update channel type to public
			if in.Type == types.ChannelTypePublic && false {
				return errors.New("Not allowed to change type of this channel to **public**")
			}

			// @todo [SECURITY] check if user can create update channel type to private
			if in.Type == types.ChannelTypePrivate && false {
				return errors.New("Not allowed to change type of this channel to **private**")
			}

			// @todo [SECURITY] check if user can update channel type to group
			if in.Type == types.ChannelTypeGroup && false {
				return errors.New("Not allowed to change type of this channel to **group**")
			}

			changed = true
		}

		var chUpdatorId = repository.Identity(svc.ctx)

		if len(in.Name) > 0 && ch.Name != in.Name {
			// @todo [SECURITY] can we change channel's name?
			if false {
				return errors.New("Not allowed to rename channel")
			} else if settingsChannelNameLength > 0 && len(in.Name) > settingsChannelNameLength {
				return errors.Errorf("Channel name (%d characters) too long (max: %d)", len(in.Name), settingsChannelNameLength)
			} else if ch.Name != "" {
				svc.scheduleSystemMessage(in, "<@%d> renamed channel **%s** (was: %s)", chUpdatorId, in.Name, ch.Name)
			} else {
				svc.scheduleSystemMessage(in, "<@%d> set channel name to **%s**", chUpdatorId, in.Name)
			}

			ch.Name = in.Name
			changed = true
		}

		if len(in.Topic) > 0 && ch.Topic != in.Topic {
			// @todo [SECURITY] can we change channel's topic?
			if false {
				return errors.New("Not allowed to change channel topic")
			} else if settingsChannelTopicLength > 0 && len(in.Topic) > settingsChannelTopicLength {
				return errors.Errorf("Channel topic (%d characters) too long (max: %d)", len(in.Topic), settingsChannelTopicLength)
			} else if ch.Topic != "" {
				svc.scheduleSystemMessage(in, "<@%d> changed channel topic: %s (was: %s)", chUpdatorId, in.Topic, ch.Topic)
			} else {
				svc.scheduleSystemMessage(in, "<@%d> set channel topic to %s", chUpdatorId, in.Topic)
			}

			ch.Topic = in.Topic
			changed = true
		}

		if !changed {
			return nil
		}
		// Save the updated channel
		if ch, err = svc.channel.UpdateChannel(in); err != nil {
			return
		}

		_ = svc.flushSystemMessages()

		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Delete(ID uint64) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !ch.CanDelete {
			return errors.New("Not allowed to delete this channel")
		}

		if ch.DeletedAt != nil {
			return errors.New("Channel already deleted")
		} else {
			now := time.Now()
			ch.DeletedAt = &now
		}

		svc.scheduleSystemMessage(ch, "<@%d> deleted this channel", userID)

		if err = svc.channel.DeleteChannelByID(ID); err != nil {
			return
		} else {
			// Set deletedAt timestamp so that our clients can react properly...
			ch.DeletedAt = timeNowPtr()
		}

		_ = svc.sendChannelEvent(ch)
		_ = svc.flushSystemMessages()
		return nil
	})
}

func (svc *channel) Undelete(ID uint64) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !ch.CanDelete {
			return errors.New("Not allowed to undelete this channel")
		}

		if ch.DeletedAt == nil {
			return errors.New("Channel not deleted")
		}

		svc.scheduleSystemMessage(ch, "<@%d> undeleted this channel", userID)

		if err = svc.channel.UndeleteChannelByID(ID); err != nil {
			return
		} else {
			// Remove deletedAt timestamp so that our clients can react properly...
			ch.DeletedAt = nil
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) SetFlag(ID uint64, flag types.ChannelMembershipFlag) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var membership *types.ChannelMember
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if members, err := svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: ch.ID, MemberID: userID}); err != nil {
			return err
		} else if len(members) == 1 {
			membership = members[0]
			membership.Flag = flag
		}

		if membership == nil {
			return errors.New("not a member")
		}

		if ch.Member, err = svc.cmember.Update(membership); err != nil {
			return
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Archive(ID uint64) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !ch.CanArchive {
			return errors.New("Not allowed to archive this channel")
		}

		if ch.ArchivedAt != nil {
			return errors.New("Channel already archived")
		}

		svc.scheduleSystemMessage(ch, "<@%d> archived this channel", userID)

		if err = svc.channel.ArchiveChannelByID(ID); err != nil {
			return
		} else {
			// Set archivedAt timestamp so that our clients can react properly...
			ch.ArchivedAt = timeNowPtr()
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Unarchive(ID uint64) (ch *types.Channel, err error) {
	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !ch.CanArchive {
			return errors.New("Not allowed to unarchive this channel")
		}

		if ch.ArchivedAt == nil {
			return errors.New("Channel not archived")
		}

		if err = svc.channel.UnarchiveChannelByID(ID); err != nil {
			return
		} else {
			// Unset archivedAt timestamp so that our clients can react properly...
			ch.ArchivedAt = nil
		}

		svc.scheduleSystemMessage(ch, "<@%d> unarchived this channel", userID)

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) InviteUser(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	out = types.ChannelMemberSet{}

	if ch, err = svc.FindByID(channelID); err != nil {
		return
	}

	if ch.Type == types.ChannelTypeGroup {
		return nil, errors.New("Adding members to a group is not currently supported")
	}

	if !ch.CanChangeMembers {
		return nil, errors.New("Not allowed to invite members")
	}

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		users, err := svc.usr.Find(nil)
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			user := users.FindByID(memberID)
			if user == nil {
				return errors.New("Unexisting user")
			}

			if e := existing.FindByUserID(memberID); e != nil {
				// Already a member/invited
				e.User = user
				out = append(out, e)
				continue
			}

			svc.scheduleSystemMessage(ch, "<@%d> invited <@%d> to the channel", userID, memberID)

			member := &types.ChannelMember{
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeInvitee,
			}

			if member, err = svc.cmember.Create(member); err != nil {
				return err
			}

			out = append(out, member)
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	out = types.ChannelMemberSet{}

	if ch, err = svc.FindByID(channelID); err != nil {
		return
	}

	if ch.Type == types.ChannelTypeGroup {
		return nil, errors.New("Adding members to a group is not currently supported")
	}

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		users, err := svc.usr.Find(nil)
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			var exists bool

			user := users.FindByID(memberID)
			if user == nil {
				return errors.New("Unexisting user")
			}

			if e := existing.FindByUserID(memberID); e != nil {
				if e.Type != types.ChannelMembershipTypeInvitee {
					e.User = user
					out = append(out, e)
					continue
				} else {
					exists = true
				}
			}

			// @todo [SECURITY] implement proper checking
			if !(ch.CanChangeMembers || memberID == userID && ch.Type == types.ChannelTypePublic) {
				return errors.New("Not allowed to add members")
			}

			if !exists {
				if userID == memberID {
					svc.scheduleSystemMessage(ch, "<@%d> joined", memberID)
				} else {
					svc.scheduleSystemMessage(ch, "<@%d> added <@%d> to the channel", userID, memberID)
				}
			}

			member := &types.ChannelMember{
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeOwner,
				User:      user,
			}

			if exists {
				member, err = svc.cmember.Update(member)
			} else {
				member, err = svc.cmember.Create(member)
			}

			svc.evl.Join(memberID, channelID)

			if err != nil {
				return err
			}

			out = append(out, member)
		}

		// Push channel to all members
		if err = svc.sendChannelEvent(ch); err != nil {
			return
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) DeleteMember(channelID uint64, memberIDs ...uint64) (err error) {
	var (
		userID   = repository.Identity(svc.ctx)
		ch       *types.Channel
		existing types.ChannelMemberSet
	)

	if ch, err = svc.FindByID(channelID); err != nil {
		return
	}

	if ch.Type == types.ChannelTypeGroup {
		return errors.New("removing members from a group is currently not supported")
	}

	return svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
			return
		}

		for _, memberID := range memberIDs {
			if existing.FindByUserID(memberID) == nil {
				// Not really a member...
				continue
			}

			// @todo [SECURITY] implement proper checking
			if !(ch.CanChangeMembers || memberID == userID && ch.Type == types.ChannelTypePublic) {
				return errors.New("Not allowed to add members")
			}

			if userID == memberID {
				svc.scheduleSystemMessage(ch, "<@%d> left the channel", memberID)
			} else {
				svc.scheduleSystemMessage(ch, "<@%d> removed from the channel", memberID)
			}

			if err = svc.cmember.Delete(channelID, memberID); err != nil {
				return err
			}

			svc.evl.Part(memberID, channelID)
		}

		return svc.flushSystemMessages()
	})
}

// RecordView
// @deprecated
func (svc *channel) RecordView(userID, channelID, lastMessageID uint64) error {
	return svc.db.Transaction(func() (err error) {
		return svc.unread.Record(userID, channelID, 0, lastMessageID, 0)
	})
}

func (svc *channel) scheduleSystemMessage(ch *types.Channel, format string, a ...interface{}) {
	svc.sysmsgs = append(svc.sysmsgs, &types.Message{
		ChannelID: ch.ID,
		Message:   fmt.Sprintf(format, a...),
		Type:      types.MessageTypeChannelEvent,
	})
}

// Flushes sys message stack, stores them into repo & pushes them into event loop
func (svc *channel) flushSystemMessages() (err error) {
	defer func() {
		svc.sysmsgs = types.MessageSet{}
	}()

	return svc.sysmsgs.Walk(func(msg *types.Message) error {
		if msg, err = svc.message.CreateMessage(msg); err != nil {
			return err
		} else {
			return svc.evl.Message(msg)
		}
	})
}

// Sends channel event
func (svc *channel) sendChannelEvent(ch *types.Channel) (err error) {
	if ch.DeletedAt == nil && ch.ArchivedAt == nil {
		// Looks like a valid channel

		// Preload members, if needed
		if len(ch.Members) == 0 || ch.Member == nil {
			if mm, err := svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: ch.ID}); err != nil {
				return err
			} else {
				ch.Members = mm.AllMemberIDs()
				ch.Member = mm.FindByUserID(repository.Identity(svc.ctx))
			}
		}
	}

	if err = svc.setPermissionFlags(ch); err != nil {
		return
	}

	if err = svc.evl.Channel(ch); err != nil {
		return
	}

	return nil
}

func (svc *channel) setPermissionFlags(ch *types.Channel) (err error) {
	var userID = repository.Identity(svc.ctx)

	var (
		isMember  = ch.Member != nil
		isCreator = ch.CreatorID == userID
		isOwner   = isCreator || (isMember && ch.Member.Type == types.ChannelMembershipTypeOwner)
		isPublic  = ch.Type == types.ChannelTypePublic
	)

	ch.CanJoin = (ch.IsValid() && isPublic) || isOwner
	ch.CanPart = isMember && ch.Type != types.ChannelTypeGroup
	ch.CanObserve = (ch.IsValid() && isPublic) || isMember
	ch.CanSendMessages = ch.CanObserve && isMember
	ch.CanDeleteMessages = isOwner
	ch.CanChangeMembers = isOwner
	ch.CanUpdate = isOwner
	ch.CanArchive = isOwner
	ch.CanDelete = isOwner

	return nil
}

var _ ChannelService = &channel{}
