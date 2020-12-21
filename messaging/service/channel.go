package service

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	channel struct {
		ctx   context.Context
		event EventService
		ac    applicationAccessController

		actionlog actionlog.Recorder
		store     store.Storer

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
	return (&channel{
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
	}).With(ctx)
}

func (svc *channel) With(ctx context.Context) ChannelService {
	return &channel{
		ctx: ctx,

		event: Event(ctx),
		ac:    DefaultAccessControl,

		actionlog: DefaultActionlog,
		store:     svc.store,

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

	if ch, err = store.LookupMessagingChannelByID(svc.ctx, svc.store, ID); err != nil {
		if errors.IsNotFound(err) {
			return nil, ChannelErrNotFound()
		}

		return nil, err
	}

	if err = svc.preloadExtras(svc.ctx, svc.store, ch); err != nil {
		return nil, err
	}

	return
}

func (svc *channel) Find(filter types.ChannelFilter) (set types.ChannelSet, f types.ChannelFilter, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(svc.ctx).Identity()
	filter.Check = func(c *types.Channel) (b bool, e error) {
		return svc.ac.CanReadChannel(svc.ctx, c), nil
	}

	set, f, err = store.SearchMessagingChannels(svc.ctx, svc.store, filter)
	if err != nil {
		return
	}

	err = svc.preloadExtras(svc.ctx, svc.store, set...)
	if err != nil {
		return
	}

	return
}

// preloadExtras pre-loads channel's members, views
func (svc *channel) preloadExtras(ctx context.Context, s store.Storer, cc ...*types.Channel) (err error) {
	if len(cc) == 0 {
		return nil
	}

	if err = svc.preloadMembers(svc.ctx, svc.store, cc...); err != nil {
		return
	}

	if err = types.ChannelSet(cc).Walk(svc.setPermissionFlags); err != nil {
		return err
	}

	if err = svc.preloadUnreads(svc.ctx, svc.store, cc...); err != nil {
		return
	}

	return
}

func (channel) preloadMembers(ctx context.Context, s store.Storer, cc ...*types.Channel) (err error) {
	if len(cc) == 0 {
		return nil
	}

	var (
		userID = auth.GetIdentityFromContext(ctx).Identity()
		mm     types.ChannelMemberSet
		f      = types.ChannelMemberFilterChannels(types.ChannelSet(cc).IDs()...)
	)

	// Load membership info of all channels
	if mm, _, err = store.SearchMessagingChannelMembers(ctx, s, f); err != nil {
		return
	} else {
		err = types.ChannelSet(cc).Walk(func(ch *types.Channel) error {
			ch.Members = mm.MembersOf(ch.ID)
			ch.Member = mm.FindByChannelID(ch.ID).FindByUserID(userID)
			return nil
		})
	}

	return
}

func (channel) preloadUnreads(ctx context.Context, s store.Storer, cc ...*types.Channel) (err error) {
	if len(cc) == 0 {
		return nil
	}

	var (
		unread types.UnreadSet
		userID = auth.GetIdentityFromContext(ctx).Identity()
	)

	unread, err = store.CountMessagingUnread(ctx, s, userID, 0)
	if err != nil {
		return
	}

	_ = types.ChannelSet(cc).Walk(func(ch *types.Channel) error {
		ch.Unread = unread.FindByChannelId(ch.ID)
		return nil
	})

	unread, err = store.CountMessagingUnreadThreads(ctx, s, userID, 0)
	if err != nil {
		return
	}

	_ = types.ChannelSet(cc).Walk(func(ch *types.Channel) error {
		var u = unread.FindByChannelId(ch.ID)

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

	return nil
}

// FindMembers loads all members (and full users) for a specific channel
func (svc *channel) FindMembers(channelID uint64) (out types.ChannelMemberSet, err error) {
	if _, err = svc.FindByID(channelID); err != nil {
		return
	}

	out, _, err = store.SearchMessagingChannelMembers(svc.ctx, svc.store, types.ChannelMemberFilterChannels(channelID))
	if err != nil {
		return
	}

	return
}

func (svc *channel) Create(new *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: new}
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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
			ch, err = store.LookupMessagingChannelByMemberSet(ctx, s, mm.AllMemberIDs()...)
			if err == nil {
				if err = svc.preloadExtras(ctx, s, ch); err != nil {
					return
				}

				if !ch.CanObserve {
					return ChannelErrNotAllowedToRead()
				} else {
					// Group already exists so let's just return it
					return nil
				}
			} else if !errors.IsNotFound(err) {
				return
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
			ID:               nextID(),
			Name:             new.Name,
			Topic:            new.Topic,
			Type:             new.Type,
			MembershipPolicy: new.MembershipPolicy,
			CreatorID:        chCreatorID,
			CreatedAt:        *now(),
		}

		// Save the channel
		if err = store.CreateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		err = mm.Walk(func(m *types.ChannelMember) (err error) {
			// Assign channel ID to membership
			m.ChannelID = ch.ID

			// Create member
			if err = svc.createMember(ctx, s, m); err != nil {
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
		ch.Member = mm.FindByUserID(chCreatorID)

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

		// sending copy of channel to event so that members are not accidentally overwritten
		var evCh = *ch
		return svc.sendChannelEvent(&evCh)
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

func (svc *channel) Update(upd *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: upd}
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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
				return fmt.Errorf("channel name (%d characters) too long (max: %d)", len(upd.Name), settingsChannelNameLength)
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
				return fmt.Errorf("channel topic (%d characters) too long (max: %d)", len(upd.Topic), settingsChannelTopicLength)
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

		ch.UpdatedAt = now()

		// Save the updated channel
		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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
		}

		// Set deletedAt timestamp so that our clients can react properly...
		ch.DeletedAt = now()

		svc.scheduleSystemMessage(ch, "<@%d> deleted this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		ch.DeletedAt = nil

		svc.scheduleSystemMessage(ch, "<@%d> undeleted this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		_ = svc.flushSystemMessages()
		return svc.sendChannelEvent(ch)
	})

	return ch, svc.recordAction(svc.ctx, aProps, ChannelActionUndelete, err)
}

func (svc *channel) SetFlag(ID uint64, flag types.ChannelMembershipFlag) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{flag: string(flag)}
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var membership *types.ChannelMember
		var userID = auth.GetIdentityFromContext(svc.ctx).Identity()

		if ch, err = svc.FindByID(ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		f := types.ChannelMemberFilter{ChannelID: []uint64{ch.ID}, MemberID: []uint64{userID}}
		if members, _, err := store.SearchMessagingChannelMembers(ctx, s, f); err != nil {
			return err
		} else if len(members) == 1 {
			membership = members[0]
			membership.Flag = flag
		}

		if membership == nil {
			return ChannelErrNotMember()
		}

		ch.Member = membership
		ch.Member.UpdatedAt = now()
		if err = store.UpdateMessagingChannelMember(ctx, s, membership); err != nil {
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		ch.ArchivedAt = now()

		svc.scheduleSystemMessage(ch, "<@%d> archived this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		ch.ArchivedAt = nil

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		f := types.ChannelMemberFilterChannels(channelID)
		if existing, _, err = store.SearchMessagingChannelMembers(ctx, s, f); err != nil {
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
				CreatedAt: *now(),
				ChannelID: channelID,
				UserID:    memberID,
				Type:      types.ChannelMembershipTypeInvitee,
			}

			if err = store.CreateMessagingChannelMember(ctx, s, member); err != nil {
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

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		f := types.ChannelMemberFilterChannels(channelID)
		if existing, _, err = store.SearchMessagingChannelMembers(ctx, s, f); err != nil {
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
				member.UpdatedAt = now()
				err = store.UpdateMessagingChannelMember(ctx, s, member)
			} else {
				err = svc.createMember(ctx, s, member)
			}

			if err != nil {
				return err
			}

			svc.event.Join(memberID, channelID)

			out = append(out, member)

			ch.Member = member
			ch.Members = out.AllMemberIDs()
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
func (svc channel) createMember(ctx context.Context, s store.Storer, m *types.ChannelMember) (err error) {
	m.CreatedAt = *now()

	if err = store.CreateMessagingChannelMember(ctx, s, m); err != nil {
		return
	}

	// Create zero-count unread record
	if err = store.PresetMessagingUnread(ctx, s, m.ChannelID, 0, m.UserID); err != nil {
		return
	}

	return
}

func (svc *channel) DeleteMember(channelID uint64, memberIDs ...uint64) (err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(svc.ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		f := types.ChannelMemberFilterChannels(channelID)
		if existing, _, err = store.SearchMessagingChannelMembers(ctx, s, f); err != nil {
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

			if err = store.DeleteMessagingChannelMemberByChannelIDUserID(ctx, s, channelID, memberID); err != nil {
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
		msg.ID = nextID()
		msg.CreatedAt = *now()

		if err = store.CreateMessagingMessage(svc.ctx, svc.store, msg); err != nil {
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
			f := types.ChannelMemberFilterChannels(ch.ID)
			if mm, _, err := store.SearchMessagingChannelMembers(svc.ctx, svc.store, f); err != nil {
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
