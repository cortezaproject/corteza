package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	channel struct {
		db  db
		ctx context.Context

		event EventService
		ac    applicationAccessController

		actionlog actionlog.Recorder

		channel repository.ChannelRepository
		cmember repository.ChannelMemberRepository
		unread  repository.UnreadRepository
		message repository.MessageRepository

		sysmsgs types.MessageSet
	}

	applicationAccessController interface {
		CanCreatePublicChannel(context.Context) bool
		CanCreatePrivateChannel(context.Context) bool
		CanCreateGroupChannel(context.Context) bool
		CanUpdateChannel(context.Context, *types.Channel) bool
		CanReadChannel(context.Context, *types.Channel) bool
		CanJoinChannel(context.Context, *types.Channel) bool
		CanLeaveChannel(context.Context, *types.Channel) bool
		CanDeleteChannel(context.Context, *types.Channel) bool
		CanUndeleteChannel(context.Context, *types.Channel) bool
		CanArchiveChannel(context.Context, *types.Channel) bool
		CanUnarchiveChannel(context.Context, *types.Channel) bool
		CanManageChannelMembers(context.Context, *types.Channel) bool
		CanChangeChannelMembershipPolicy(context.Context, *types.Channel) bool
		CanSendMessage(context.Context, *types.Channel) bool
		CanUpdateOwnMessages(context.Context, *types.Channel) bool
		CanUpdateMessages(context.Context, *types.Channel) bool
		CanDeleteOwnMessages(context.Context, *types.Channel) bool
		CanDeleteMessages(context.Context, *types.Channel) bool
	}

	ChannelService interface {
		With(ctx context.Context) ChannelService

		FindByID(channelID uint64) (*types.Channel, error)
		Find(types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error)

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
	}
)

const (
	settingsChannelNameLength  = 40
	settingsChannelTopicLength = 200
)

func Channel(ctx context.Context) ChannelService {
	return (&channel{}).With(ctx)
}

func (svc *channel) With(ctx context.Context) ChannelService {
	db := repository.DB(ctx)
	return &channel{
		db:  db,
		ctx: ctx,

		event: Event(ctx),
		ac:    DefaultAccessControl,

		actionlog: DefaultActionlog,

		channel: repository.Channel(ctx, db),
		cmember: repository.ChannelMember(ctx, db),
		unread:  repository.Unread(ctx, db),
		message: repository.Message(ctx, db),

		// System messages should be flushed at the end of each session
		sysmsgs: types.MessageSet{},
	}
}

func (svc *channel) FindByID(ID uint64) (ch *types.Channel, err error) {
	if ch, err = svc.findByID(ID); err != nil {
		return
	}

	if !svc.ac.CanReadChannel(svc.ctx, ch) {
		return nil, ErrNoPermissions.withStack()
	}

	return
}

func (svc *channel) findByID(ID uint64) (ch *types.Channel, err error) {
	ch, err = svc.channel.FindByID(ID)
	if err != nil {
		if repository.ErrChannelNotFound.Eq(err) {
			return nil, ChannelErrNotFound()
		}

		return nil, err
	}

	if err = svc.preloadExtras(types.ChannelSet{ch}); err != nil {
		return nil, err
	}

	return
}

func (svc *channel) Find(filter types.ChannelFilter) (set types.ChannelSet, f types.ChannelFilter, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	set, f, err = svc.channel.Find(filter)
	if err == nil {
		err = svc.preloadExtras(set)
	}

	set, err = set.Filter(func(c *types.Channel) (b bool, e error) {
		return svc.ac.CanReadChannel(svc.ctx, c), nil
	})

	return
}

// preloadExtras pre-loads channel's members, views
func (svc *channel) preloadExtras(cc types.ChannelSet) (err error) {
	if err = svc.preloadMembers(cc); err != nil {
		return
	}

	if err = cc.Walk(svc.setPermissionFlags); err != nil {
		return err
	}

	if err = svc.preloadUnreads(cc); err != nil {
		return
	}

	return
}

func (svc *channel) preloadMembers(cc types.ChannelSet) (err error) {
	var (
		userID = auth.GetIdentityFromContext(svc.ctx).Identity()
		mm     types.ChannelMemberSet
	)

	// Load membership info of all channels
	if mm, err = svc.cmember.Find(types.ChannelMemberFilterChannels(cc.IDs()...)); err != nil {
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

// preload channel unread info for a single user
func (svc *channel) preloadUnreads(cc types.ChannelSet) error {
	var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

	if uu, err := svc.unread.Count(userID, 0); err != nil {
		return err
	} else {
		_ = cc.Walk(func(ch *types.Channel) error {
			ch.Unread = uu.FindByChannelId(ch.ID)
			return nil
		})
	}

	if uu, err := svc.unread.CountThreads(userID, 0); err != nil {
		return err
	} else {
		_ = cc.Walk(func(ch *types.Channel) error {
			var u = uu.FindByChannelId(ch.ID)

			if u == nil {
				return nil
			}

			if ch.Unread == nil {
				ch.Unread = &types.Unread{}
			}

			ch.Unread.ThreadCount = u.ThreadCount
			ch.Unread.ThreadTotal = u.ThreadTotal

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

		out, err = svc.cmember.Find(types.ChannelMemberFilterChannels(channelID))
		if err != nil {
			return err
		}

		return
	})
}

func (svc *channel) Create(new *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: new}
	)

	err = svc.db.Transaction(func() (err error) {
		if !new.Type.IsValid() {
			return ChannelErrInvalidType()
		}

		if len(new.Name) == 0 && new.Type != types.ChannelTypeGroup {
			return ChannelErrNameEmpty()
		}

		if settingsChannelNameLength > 0 && len(new.Name) > settingsChannelNameLength {
			return ChannelErrNameLength()
		}

		if len(new.Topic) > 0 && settingsChannelTopicLength > 0 && len(new.Topic) > settingsChannelTopicLength {
			return ChannelErrTopicLength()
		}

		var chCreatorID = auth.GetIdentityFromContext(svc.ctx).Identity()

		mm := svc.buildMemberSet(chCreatorID, new.Members...)

		if new.Type == types.ChannelTypeGroup {
			if ch, err = svc.checkGroupExistance(mm); err != nil {
				return err
			} else if ch != nil && ch.CanObserve {
				// Group already exists so let's just return it
				return nil
			} else if ch != nil && !ch.CanObserve {
				return ChannelErrNotAllowedToRead()
			}
		}

		if new.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(svc.ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if new.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(svc.ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if new.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(svc.ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if !new.MembershipPolicy.IsValid() {
			// Reset invalid membership flag to default
			new.MembershipPolicy = types.ChannelMembershipPolicyDefault
		}

		if new.MembershipPolicy != types.ChannelMembershipPolicyDefault && !svc.ac.CanChangeChannelMembershipPolicy(svc.ctx, new) {
			return ChannelErrNotAllowedToCreate()
		}

		// This is a fresh channel, just copy values
		ch = &types.Channel{
			Name:             new.Name,
			Topic:            new.Topic,
			Type:             new.Type,
			MembershipPolicy: new.MembershipPolicy,
			CreatorID:        chCreatorID,
		}

		// Save the channel
		if ch, err = svc.channel.Create(ch); err != nil {
			return
		}

		err = mm.Walk(func(m *types.ChannelMember) (err error) {
			// Assign channel ID to membership
			m.ChannelID = ch.ID

			// Create member
			if m, err = svc.createMember(m); err != nil {
				return err
			}

			// Subscribe all members
			return svc.event.Join(m.UserID, ch.ID)
		})

		if err != nil {
			// Could not add member
			return
		}

		// Copy all member IDs to channel's member slice
		ch.Members = mm.AllMemberIDs()

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if len(ch.Name) == 0 {
			svc.scheduleSystemMessage(ch, `<@%d> created %s channel`, chCreatorID, ch.Type)
		} else if len(ch.Topic) == 0 {
			svc.scheduleSystemMessage(ch, `<@%d> created %s channel **%s**`, chCreatorID, ch.Type, ch.Name)
		} else {
			svc.scheduleSystemMessage(ch, `<@%d> created %s channel **%s**, topic: %s`, chCreatorID, ch.Type, ch.Name, ch.Topic)
		}

		_ = svc.flushSystemMessages()

		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionCreate, err)
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
	if out, err = svc.channel.FindByMemberSet(mm.AllMemberIDs()...); err == repository.ErrChannelNotFound {
		return nil, nil
	} else if out != nil && err == nil {
		err = svc.preloadExtras(types.ChannelSet{out})
	}

	return
}

func (svc *channel) Update(upd *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: upd}
	)

	err = svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return ChannelErrInvalidID()
		}

		if !upd.Type.IsValid() {
			return ChannelErrInvalidType()
		}

		if len(upd.Name) == 0 && upd.Type != types.ChannelTypeGroup {
			return ChannelErrNameEmpty()
		}

		if settingsChannelNameLength > 0 && len(upd.Name) > settingsChannelNameLength {
			return ChannelErrNameLength()
		}

		if len(upd.Topic) > 0 && settingsChannelTopicLength > 0 && len(upd.Topic) > settingsChannelTopicLength {
			return ChannelErrTopicLength()
		}
		var changed bool

		if ch, err = svc.FindByID(upd.ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUpdateChannel(svc.ctx, ch) {
			return ChannelErrNotAllowedToUpdate()
		}

		if upd.Type.IsValid() && ch.Type != upd.Type {
			if upd.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(svc.ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			if upd.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(svc.ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			if upd.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(svc.ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			changed = true
		}

		var chUpdatorId = auth.GetIdentityFromContext(svc.ctx).Identity()

		if len(upd.Name) > 0 && ch.Name != upd.Name {
			if settingsChannelNameLength > 0 && len(upd.Name) > settingsChannelNameLength {
				return errors.Errorf("channel name (%d characters) too long (max: %d)", len(upd.Name), settingsChannelNameLength)
			} else if ch.Name != "" {
				svc.scheduleSystemMessage(upd, "<@%d> renamed channel **%s** (was: %s)", chUpdatorId, upd.Name, ch.Name)
			} else {
				svc.scheduleSystemMessage(upd, "<@%d> set channel name to **%s**", chUpdatorId, upd.Name)
			}

			ch.Name = upd.Name
			changed = true
		}

		if len(upd.Topic) > 0 && ch.Topic != upd.Topic {
			if settingsChannelTopicLength > 0 && len(upd.Topic) > settingsChannelTopicLength {
				return errors.Errorf("channel topic (%d characters) too long (max: %d)", len(upd.Topic), settingsChannelTopicLength)
			} else if ch.Topic != "" {
				svc.scheduleSystemMessage(upd, "<@%d> changed channel topic: %s (was: %s)", chUpdatorId, upd.Topic, ch.Topic)
			} else {
				svc.scheduleSystemMessage(upd, "<@%d> set channel topic to %s", chUpdatorId, upd.Topic)
			}

			ch.Topic = upd.Topic
			changed = true
		}

		if ch.MembershipPolicy != upd.MembershipPolicy && !svc.ac.CanChangeChannelMembershipPolicy(svc.ctx, ch) {
			return ChannelErrNotAllowedToUpdate()
		} else {
			ch.MembershipPolicy = upd.MembershipPolicy
			changed = true
		}

		if !changed {
			return nil
		}
		// Save the updated channel
		if ch, err = svc.channel.Update(upd); err != nil {
			return
		}

		_ = svc.flushSystemMessages()

		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionUpdate, err)

}

func (svc *channel) Delete(ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}

		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.findByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanDeleteChannel(svc.ctx, ch) {
			return ChannelErrNotAllowedToDelete()
		}

		if ch.DeletedAt != nil {
			return ChannelErrAlreadyDeleted()
		} else {
			now := time.Now()
			ch.DeletedAt = &now
		}

		svc.scheduleSystemMessage(ch, "<@%d> deleted this channel", userID)

		if err = svc.channel.DeleteByID(ID); err != nil {
			return
		} else {
			// Set deletedAt timestamp so that our clients can react properly...
			ch.DeletedAt = timeNowPtr()
		}

		_ = svc.sendChannelEvent(ch)
		_ = svc.flushSystemMessages()
		return nil
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionDelete, err)
}

func (svc *channel) Undelete(ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}

		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.findByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUndeleteChannel(svc.ctx, ch) {
			return ChannelErrNotAllowedToUndelete()
		}

		if ch.DeletedAt == nil {
			return ChannelErrNotDeleted()
		}

		svc.scheduleSystemMessage(ch, "<@%d> undeleted this channel", userID)

		if err = svc.channel.UndeleteByID(ID); err != nil {
			return
		} else {
			// Remove deletedAt timestamp so that our clients can react properly...
			ch.DeletedAt = nil
		}

		svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionUndelete, err)
}

func (svc *channel) SetFlag(ID uint64, flag types.ChannelMembershipFlag) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{flag: string(flag)}
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var membership *types.ChannelMember
		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if members, err := svc.cmember.Find(types.ChannelMemberFilter{ChannelID: []uint64{ch.ID}, MemberID: []uint64{userID}}); err != nil {
			return err
		} else if len(members) == 1 {
			membership = members[0]
			membership.Flag = flag
		}

		if membership == nil {
			return ChannelErrNotMember()
		}

		if ch.Member, err = svc.cmember.Update(membership); err != nil {
			return
		}

		return svc.flushSystemMessages()

		// Setting a flag on a channel is a private thing,
		// no need to send channel event back to everyone
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionSetFlag, err)

}

func (svc *channel) Archive(ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.findByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanArchiveChannel(svc.ctx, ch) {
			return ChannelErrNotAllowedToUndelete()
		}

		if ch.ArchivedAt != nil {
			return ChannelErrAlreadyArchived()
		}

		svc.scheduleSystemMessage(ch, "<@%d> archived this channel", userID)

		if err = svc.channel.ArchiveByID(ID); err != nil {
			return
		} else {
			// Set archivedAt timestamp so that our clients can react properly...
			ch.ArchivedAt = timeNowPtr()
		}

		_ = svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionArchive, err)
}

func (svc *channel) Unarchive(ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.findByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUnarchiveChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.ArchivedAt == nil {
			return ChannelErrNotArchived()
		}

		if err = svc.channel.UnarchiveByID(ID); err != nil {
			return
		} else {
			// Unset archivedAt timestamp so that our clients can react properly...
			ch.ArchivedAt = nil
		}

		svc.scheduleSystemMessage(ch, "<@%d> unarchived this channel", userID)

		_ = svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionUnarchive, err)
}

func (svc *channel) InviteUser(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if channelID == 0 {
			return ChannelErrInvalidID()
		}

		for _, memberID := range memberIDs {
			if memberID == 0 {
				return ChannelErrInvalidID()
			}
		}

		var (
			userID   = auth.GetIdentityFromContext(svc.ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		out = types.ChannelMemberSet{}

		if ch, err = svc.FindByID(channelID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if ch.Type == types.ChannelTypeGroup {
			return ChannelErrUnableToManageGroupMembers()
		}

		if !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
			return ChannelErrNotAllowedToManageMembers()
		}

		if existing, err = svc.cmember.Find(types.ChannelMemberFilterChannels(channelID)); err != nil {
			return
		}

		for _, memberID := range memberIDs {
			if e := existing.FindByUserID(memberID); e != nil {
				// Already a member/invited
				out = append(out, e)
				continue
			}

			svc.scheduleSystemMessage(ch, "<@%d> invited <@%d> to the channel", userID, memberID)

			member := &types.ChannelMember{
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeInvitee,
			}

			if member, err = svc.createMember(member); err != nil {
				return err
			}

			out = append(out, member)
		}

		return svc.flushSystemMessages()
	})

	return out, svc.recordAction(svc.ctx, aProps, ChannelActionInviteMember, err)
}

func (svc *channel) AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		if channelID == 0 {
			return ChannelErrInvalidID()
		}

		for _, memberID := range memberIDs {
			if memberID == 0 {
				return ChannelErrInvalidID()
			}
		}

		var (
			userID   = auth.GetIdentityFromContext(svc.ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		out = types.ChannelMemberSet{}

		if ch, err = svc.FindByID(channelID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if ch.Type == types.ChannelTypeGroup {
			return ChannelErrUnableToManageGroupMembers()
		}

		if existing, err = svc.cmember.Find(types.ChannelMemberFilterChannels(channelID)); err != nil {
			return
		}

		for _, memberID := range memberIDs {
			var exists bool

			if e := existing.FindByUserID(memberID); e != nil {
				if e.Type != types.ChannelMembershipTypeInvitee {
					out = append(out, e)
					continue
				} else {
					exists = true
				}
			}

			if memberID == userID && !svc.ac.CanJoinChannel(svc.ctx, ch) {
				return ChannelErrNotAllowedToJoin()
			} else if memberID != userID && !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
				return ChannelErrNotAllowedToManageMembers()
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
				Type:      types.ChannelMembershipTypeMember,
			}

			if exists {
				member, err = svc.cmember.Update(member)
			} else {
				member, err = svc.createMember(member)
			}

			if err != nil {
				return err
			}

			svc.event.Join(memberID, channelID)

			out = append(out, member)
		}

		// Push channel to all members
		if err = svc.sendChannelEvent(ch); err != nil {
			return
		}

		return svc.flushSystemMessages()
	})

	return out, svc.recordAction(svc.ctx, aProps, ChannelActionAddMember, err)
}

// createMember orchestrates member creation
func (svc channel) createMember(member *types.ChannelMember) (m *types.ChannelMember, err error) {
	if m, err = svc.cmember.Create(member); err != nil {
		return
	}

	// Create zero-count unread record
	if err = svc.unread.Preset(m.ChannelID, 0, m.UserID); err != nil {
		return
	}

	return
}

func (svc *channel) DeleteMember(channelID uint64, memberIDs ...uint64) (err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = svc.db.Transaction(func() (err error) {
		var (
			userID   = auth.GetIdentityFromContext(svc.ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		if ch, err = svc.FindByID(channelID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if ch.Type == types.ChannelTypeGroup {
			return ChannelErrUnableToManageGroupMembers()
		}

		if existing, err = svc.cmember.Find(types.ChannelMemberFilterChannels(channelID)); err != nil {
			return
		}

		for _, memberID := range memberIDs {
			if existing.FindByUserID(memberID) == nil {
				// Not really a member...
				continue
			}

			if memberID == userID && !svc.ac.CanLeaveChannel(svc.ctx, ch) {
				return ChannelErrNotAllowedToPart()
			} else if memberID != userID && !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
				return ChannelErrNotAllowedToManageMembers()
			}

			if userID == memberID {
				svc.scheduleSystemMessage(ch, "<@%d> left the channel", memberID)
			} else {
				svc.scheduleSystemMessage(ch, "<@%d> removed from the channel", memberID)
			}

			if err = svc.cmember.Delete(channelID, memberID); err != nil {
				return err
			}

			_ = svc.event.Part(memberID, channelID)
		}

		return svc.flushSystemMessages()
	})

	return svc.recordAction(svc.ctx, aProps, ChannelActionRemoveMember, err)

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
		if msg, err = svc.message.Create(msg); err != nil {
			return err
		} else {
			return svc.event.Message(msg)
		}
	})
}

// Sends channel event
func (svc *channel) sendChannelEvent(ch *types.Channel) (err error) {
	if ch.DeletedAt == nil && ch.ArchivedAt == nil {
		// Looks like a valid channel

		// Preload members, if needed
		if len(ch.Members) == 0 || ch.Member == nil {
			if mm, err := svc.cmember.Find(types.ChannelMemberFilterChannels(ch.ID)); err != nil {
				return err
			} else {
				ch.Members = mm.AllMemberIDs()
				ch.Member = mm.FindByUserID(auth.GetIdentityFromContext(svc.ctx).Identity())
			}
		}
	}

	if err = svc.setPermissionFlags(ch); err != nil {
		return
	}

	if err = svc.event.Channel(ch); err != nil {
		return
	}

	return nil
}

func (svc *channel) setPermissionFlags(ch *types.Channel) (err error) {
	ch.CanJoin = svc.ac.CanJoinChannel(svc.ctx, ch)
	ch.CanPart = svc.ac.CanLeaveChannel(svc.ctx, ch)
	ch.CanObserve = svc.ac.CanReadChannel(svc.ctx, ch)
	ch.CanSendMessages = svc.ac.CanSendMessage(svc.ctx, ch)

	ch.CanDeleteMessages = svc.ac.CanDeleteMessages(svc.ctx, ch)
	ch.CanDeleteOwnMessages = svc.ac.CanDeleteOwnMessages(svc.ctx, ch)
	ch.CanUpdateMessages = svc.ac.CanUpdateMessages(svc.ctx, ch)
	ch.CanUpdateOwnMessages = svc.ac.CanUpdateOwnMessages(svc.ctx, ch)
	ch.CanChangeMembers = svc.ac.CanManageChannelMembers(svc.ctx, ch)
	// @todo migrate to proper change-membership-policy action check
	ch.CanChangeMembershipPolicy = svc.ac.CanChangeChannelMembershipPolicy(svc.ctx, ch)

	ch.CanUpdate = svc.ac.CanUpdateChannel(svc.ctx, ch)
	ch.CanArchive = svc.ac.CanArchiveChannel(svc.ctx, ch)
	ch.CanUnarchive = svc.ac.CanUnarchiveChannel(svc.ctx, ch)
	ch.CanDelete = svc.ac.CanDeleteChannel(svc.ctx, ch)
	ch.CanUndelete = svc.ac.CanUndeleteChannel(svc.ctx, ch)

	return nil
}

var _ ChannelService = &channel{}
