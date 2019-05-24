package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/organization"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	channel struct {
		db  db
		ctx context.Context

		event EventService
		ac    applicationAccessController

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
		CanSendMessage(context.Context, *types.Channel) bool
		CanUpdateOwnMessages(context.Context, *types.Channel) bool
		CanUpdateMessages(context.Context, *types.Channel) bool
		CanDeleteOwnMessages(context.Context, *types.Channel) bool
		CanDeleteMessages(context.Context, *types.Channel) bool
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

		channel: repository.Channel(ctx, db),
		cmember: repository.ChannelMember(ctx, db),
		unread:  repository.Unread(ctx, db),
		message: repository.Message(ctx, db),

		// System messages should be flushed at the end of each session
		sysmsgs: types.MessageSet{},
	}
}

func (svc *channel) FindByID(ID uint64) (ch *types.Channel, err error) {
	ch, err = svc.channel.FindByID(ID)
	if err != nil {
		return
	} else if err = svc.preloadExtras(types.ChannelSet{ch}); err != nil {
		return
	}

	if !svc.ac.CanReadChannel(svc.ctx, ch) {
		return nil, ErrNoPermissions.withStack()
	}

	return
}

func (svc *channel) Find(filter *types.ChannelFilter) (cc types.ChannelSet, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return cc, svc.db.Transaction(func() (err error) {
		if cc, err = svc.channel.Find(filter); err != nil {
			return
		} else if err = svc.preloadExtras(cc); err != nil {
			return
		}

		cc, err = cc.Filter(func(c *types.Channel) (b bool, e error) {
			return svc.ac.CanReadChannel(svc.ctx, c), nil
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

		return
	})
}

func (svc *channel) Create(in *types.Channel) (out *types.Channel, err error) {
	if len(in.Name) == 0 && in.Type != types.ChannelTypeGroup {
		return nil, errors.New("channel name not provided")
	}

	if settingsChannelNameLength > 0 && len(in.Name) > settingsChannelNameLength {
		return nil, errors.Errorf("channel name (%d characters) too long (max: %d)", len(in.Name), settingsChannelNameLength)
	}

	if len(in.Topic) > 0 && settingsChannelTopicLength > 0 && len(in.Topic) > settingsChannelTopicLength {
		return nil, errors.Errorf("channel topic (%d characters) too long (max: %d)", len(in.Topic), settingsChannelTopicLength)
	}

	return out, svc.db.Transaction(func() (err error) {
		var msg *types.Message

		var organisationID = organization.Corteza().ID

		var chCreatorID = repository.Identity(svc.ctx)

		mm := svc.buildMemberSet(chCreatorID, in.Members...)

		if in.Type == types.ChannelTypeGroup {
			if out, err = svc.checkGroupExistance(mm); err != nil {
				return err
			} else if out != nil && out.CanObserve {
				// Group already exists so let's just return it
				return nil
			} else if out != nil && !out.CanObserve {
				return ErrNoPermissions.withStack()
			}
		}

		if in.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(svc.ctx) {
			return ErrNoPermissions.withStack()
		}

		if in.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(svc.ctx) {
			return ErrNoPermissions.withStack()
		}

		if in.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(svc.ctx) {
			return ErrNoPermissions.withStack()
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
		if out, err = svc.channel.Create(out); err != nil {
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
			return svc.event.Join(m.UserID, out.ID)
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
	if out, err = svc.channel.FindByMemberSet(mm.AllMemberIDs()...); err == repository.ErrChannelNotFound {
		return nil, nil
	} else if out != nil && err == nil {
		err = svc.preloadExtras(types.ChannelSet{out})
	}

	return
}

func (svc *channel) Update(in *types.Channel) (ch *types.Channel, err error) {
	if in.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if len(in.Name) == 0 && in.Type != types.ChannelTypeGroup {
		return nil, errors.New("channel name not provided")
	}

	if settingsChannelNameLength > 0 && len(in.Name) > settingsChannelNameLength {
		return nil, errors.Errorf("channel name (%d characters) too long (max: %d)", len(in.Name), settingsChannelNameLength)
	}

	if len(in.Topic) > 0 && settingsChannelTopicLength > 0 && len(in.Topic) > settingsChannelTopicLength {
		return nil, errors.Errorf("channel topic (%d characters) too long (max: %d)", len(in.Topic), settingsChannelTopicLength)
	}

	return ch, svc.db.Transaction(func() (err error) {
		var changed bool

		if ch, err = svc.FindByID(in.ID); err != nil {
			return
		}

		if !svc.ac.CanUpdateChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if in.Type.IsValid() && ch.Type != in.Type {
			if in.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(svc.ctx) {
				return ErrNoPermissions.withStack()
			}

			if in.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(svc.ctx) {
				return ErrNoPermissions.withStack()
			}

			if in.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(svc.ctx) {
				return ErrNoPermissions.withStack()
			}

			changed = true
		}

		var chUpdatorId = repository.Identity(svc.ctx)

		if len(in.Name) > 0 && ch.Name != in.Name {
			if settingsChannelNameLength > 0 && len(in.Name) > settingsChannelNameLength {
				return errors.Errorf("channel name (%d characters) too long (max: %d)", len(in.Name), settingsChannelNameLength)
			} else if ch.Name != "" {
				svc.scheduleSystemMessage(in, "<@%d> renamed channel **%s** (was: %s)", chUpdatorId, in.Name, ch.Name)
			} else {
				svc.scheduleSystemMessage(in, "<@%d> set channel name to **%s**", chUpdatorId, in.Name)
			}

			ch.Name = in.Name
			changed = true
		}

		if len(in.Topic) > 0 && ch.Topic != in.Topic {
			if settingsChannelTopicLength > 0 && len(in.Topic) > settingsChannelTopicLength {
				return errors.Errorf("channel topic (%d characters) too long (max: %d)", len(in.Topic), settingsChannelTopicLength)
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
		if ch, err = svc.channel.Update(in); err != nil {
			return
		}

		_ = svc.flushSystemMessages()

		return svc.sendChannelEvent(ch)
	})
}

func (svc *channel) Delete(ID uint64) (ch *types.Channel, err error) {
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.DeletedAt != nil {
			return errors.New("channel already deleted")
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
}

func (svc *channel) Undelete(ID uint64) (ch *types.Channel, err error) {
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanUndeleteChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.DeletedAt == nil {
			return errors.New("channel not deleted")
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
}

func (svc *channel) SetFlag(ID uint64, flag types.ChannelMembershipFlag) (ch *types.Channel, err error) {
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

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
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanArchiveChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.ArchivedAt != nil {
			return errors.New("channel already archived")
		}

		svc.scheduleSystemMessage(ch, "<@%d> archived this channel", userID)

		if err = svc.channel.ArchiveByID(ID); err != nil {
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
	if ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	return ch, svc.db.Transaction(func() (err error) {
		var userID = repository.Identity(svc.ctx)

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanUnarchiveChannel(svc.ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.ArchivedAt == nil {
			return errors.New("channel not archived")
		}

		if err = svc.channel.UnarchiveByID(ID); err != nil {
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
	if channelID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	for _, memberID := range memberIDs {
		if memberID == 0 {
			return nil, ErrInvalidID.withStack()
		}
	}

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
		return nil, errors.New("adding members to a group is not currently supported")
	}

	if !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
		return nil, ErrNoPermissions.withStack()
	}

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
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

			if member, err = svc.cmember.Create(member); err != nil {
				return err
			}

			out = append(out, member)
		}

		return svc.flushSystemMessages()
	})
}

func (svc *channel) AddMember(channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	if channelID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	for _, memberID := range memberIDs {
		if memberID == 0 {
			return nil, ErrInvalidID.withStack()
		}
	}

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
		return nil, errors.New("adding members to a group is not currently supported")
	}

	return out, svc.db.Transaction(func() (err error) {
		if existing, err = svc.cmember.Find(&types.ChannelMemberFilter{ChannelID: channelID}); err != nil {
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
				return ErrNoPermissions.withStack()
			} else if memberID != userID && !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
				return ErrNoPermissions.withStack()
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
			}

			if exists {
				member, err = svc.cmember.Update(member)
			} else {
				member, err = svc.cmember.Create(member)
			}

			svc.event.Join(memberID, channelID)

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

			if memberID == userID && !svc.ac.CanJoinChannel(svc.ctx, ch) {
				return ErrNoPermissions.withStack()
			}
			if !svc.ac.CanManageChannelMembers(svc.ctx, ch) {
				return ErrNoPermissions.withStack()
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

	ch.CanUpdate = svc.ac.CanUpdateChannel(svc.ctx, ch)
	ch.CanArchive = svc.ac.CanArchiveChannel(svc.ctx, ch)
	ch.CanUnarchive = svc.ac.CanUnarchiveChannel(svc.ctx, ch)
	ch.CanDelete = svc.ac.CanDeleteChannel(svc.ctx, ch)
	ch.CanUndelete = svc.ac.CanUndeleteChannel(svc.ctx, ch)

	return nil
}

var _ ChannelService = &channel{}
