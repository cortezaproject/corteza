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
		FindByID(ctx context.Context, channelID uint64) (*types.Channel, error)
		Find(ctx context.Context, f types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error)

		Create(ctx context.Context, channel *types.Channel) (*types.Channel, error)
		Update(ctx context.Context, channel *types.Channel) (*types.Channel, error)

		FindMembers(ctx context.Context, channelID uint64) (types.ChannelMemberSet, error)

		InviteUser(ctx context.Context, channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		AddMember(ctx context.Context, channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error)
		DeleteMember(ctx context.Context, channelID uint64, memberIDs ...uint64) (err error)

		SetFlag(ctx context.Context, ID uint64, flag types.ChannelMembershipFlag) (*types.Channel, error)

		Archive(ctx context.Context, ID uint64) (*types.Channel, error)
		Unarchive(ctx context.Context, ID uint64) (*types.Channel, error)
		Delete(ctx context.Context, ID uint64) (*types.Channel, error)
		Undelete(ctx context.Context, ID uint64) (*types.Channel, error)
	}
)

const (
	settingsChannelNameLength  = 40
	settingsChannelTopicLength = 200
)

func Channel() ChannelService {
	return &channel{
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		event:     DefaultEvent,

		// System messages should be flushed at the end of each session
		sysmsgs: types.MessageSet{},
	}
}

func (svc *channel) FindByID(ctx context.Context, ID uint64) (ch *types.Channel, err error) {
	if ch, err = svc.findByID(ctx, ID); err != nil {
		return
	}

	if !svc.ac.CanReadChannel(ctx, ch) {
		return nil, ErrNoPermissions.withStack()
	}

	return
}

func (svc *channel) findByID(ctx context.Context, ID uint64) (ch *types.Channel, err error) {

	if ch, err = store.LookupMessagingChannelByID(ctx, svc.store, ID); err != nil {
		if errors.IsNotFound(err) {
			return nil, ChannelErrNotFound()
		}

		return nil, err
	}

	if err = svc.preloadExtras(ctx, svc.store, ch); err != nil {
		return nil, err
	}

	return
}

func (svc *channel) Find(ctx context.Context, filter types.ChannelFilter) (set types.ChannelSet, f types.ChannelFilter, err error) {
	filter.CurrentUserID = auth.GetIdentityFromContext(ctx).Identity()
	filter.Check = func(c *types.Channel) (b bool, e error) {
		return svc.ac.CanReadChannel(ctx, c), nil
	}

	set, f, err = store.SearchMessagingChannels(ctx, svc.store, filter)
	if err != nil {
		return
	}

	err = svc.preloadExtras(ctx, svc.store, set...)
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

	if err = svc.preloadMembers(ctx, svc.store, cc...); err != nil {
		return
	}

	err = types.ChannelSet(cc).Walk(func(c *types.Channel) error {
		return svc.setPermissionFlags(ctx, c)
	})
	if err != nil {
		return err
	}

	if err = svc.preloadUnreads(ctx, svc.store, cc...); err != nil {
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
func (svc *channel) FindMembers(ctx context.Context, channelID uint64) (out types.ChannelMemberSet, err error) {
	if _, err = svc.FindByID(ctx, channelID); err != nil {
		return
	}

	out, _, err = store.SearchMessagingChannelMembers(ctx, svc.store, types.ChannelMemberFilterChannels(channelID))
	if err != nil {
		return
	}

	return
}

func (svc *channel) Create(ctx context.Context, new *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: new}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		var chCreatorID = auth.GetIdentityFromContext(ctx).Identity()

		mm := svc.buildMemberSet(ctx, chCreatorID, new.Members...)

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

		if new.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if new.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if new.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(ctx) {
			return ChannelErrNotAllowedToCreate()
		}

		if !new.MembershipPolicy.IsValid() {
			// Reset invalid membership flag to default
			new.MembershipPolicy = types.ChannelMembershipPolicyDefault
		}

		if new.MembershipPolicy != types.ChannelMembershipPolicyDefault && !svc.ac.CanChangeChannelMembershipPolicy(ctx, new) {
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
			return svc.event.Join(ctx, m.UserID, ch.ID)
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
			svc.scheduleSystemMessage(ctx, ch, `<@%d> created %s channel`, chCreatorID, ch.Type)
		} else if len(ch.Topic) == 0 {
			svc.scheduleSystemMessage(ctx, ch, `<@%d> created %s channel **%s**`, chCreatorID, ch.Type, ch.Name)
		} else {
			svc.scheduleSystemMessage(ctx, ch, `<@%d> created %s channel **%s**, topic: %s`, chCreatorID, ch.Type, ch.Name, ch.Topic)
		}

		_ = svc.flushSystemMessages(ctx)

		// sending copy of channel to event so that members are not accidentally overwritten
		var evCh = *ch
		return svc.sendChannelEvent(ctx, &evCh)
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionCreate, err)
}

func (svc *channel) buildMemberSet(ctx context.Context, owner uint64, members ...uint64) (mm types.ChannelMemberSet) {
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

func (svc *channel) Update(ctx context.Context, upd *types.Channel) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{changed: upd}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
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

		if ch, err = svc.FindByID(ctx, upd.ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUpdateChannel(ctx, ch) {
			return ChannelErrNotAllowedToUpdate()
		}

		if upd.Type.IsValid() && ch.Type != upd.Type {
			if upd.Type == types.ChannelTypePublic && !svc.ac.CanCreatePublicChannel(ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			if upd.Type == types.ChannelTypePrivate && !svc.ac.CanCreatePrivateChannel(ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			if upd.Type == types.ChannelTypeGroup && !svc.ac.CanCreateGroupChannel(ctx) {
				return ChannelErrNotAllowedToUpdate()
			}

			changed = true
		}

		var chUpdatorId = auth.GetIdentityFromContext(ctx).Identity()

		if len(upd.Name) > 0 && ch.Name != upd.Name {
			if settingsChannelNameLength > 0 && len(upd.Name) > settingsChannelNameLength {
				return fmt.Errorf("channel name (%d characters) too long (max: %d)", len(upd.Name), settingsChannelNameLength)
			} else if ch.Name != "" {
				svc.scheduleSystemMessage(ctx, upd, "<@%d> renamed channel **%s** (was: %s)", chUpdatorId, upd.Name, ch.Name)
			} else {
				svc.scheduleSystemMessage(ctx, upd, "<@%d> set channel name to **%s**", chUpdatorId, upd.Name)
			}

			ch.Name = upd.Name
			changed = true
		}

		if len(upd.Topic) > 0 && ch.Topic != upd.Topic {
			if settingsChannelTopicLength > 0 && len(upd.Topic) > settingsChannelTopicLength {
				return fmt.Errorf("channel topic (%d characters) too long (max: %d)", len(upd.Topic), settingsChannelTopicLength)
			} else if ch.Topic != "" {
				svc.scheduleSystemMessage(ctx, upd, "<@%d> changed channel topic: %s (was: %s)", chUpdatorId, upd.Topic, ch.Topic)
			} else {
				svc.scheduleSystemMessage(ctx, upd, "<@%d> set channel topic to %s", chUpdatorId, upd.Topic)
			}

			ch.Topic = upd.Topic
			changed = true
		}

		if ch.MembershipPolicy != upd.MembershipPolicy && !svc.ac.CanChangeChannelMembershipPolicy(ctx, ch) {
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

		_ = svc.flushSystemMessages(ctx)

		return svc.sendChannelEvent(ctx, ch)
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionUpdate, err)

}

func (svc *channel) Delete(ctx context.Context, ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}

		var userID = auth.GetIdentityFromContext(ctx).Identity()

		if ch, err = svc.findByID(ctx, ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanDeleteChannel(ctx, ch) {
			return ChannelErrNotAllowedToDelete()
		}

		if ch.DeletedAt != nil {
			return ChannelErrAlreadyDeleted()
		}

		// Set deletedAt timestamp so that our clients can react properly...
		ch.DeletedAt = now()

		svc.scheduleSystemMessage(ctx, ch, "<@%d> deleted this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		_ = svc.sendChannelEvent(ctx, ch)
		_ = svc.flushSystemMessages(ctx)
		return nil
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionDelete, err)
}

func (svc *channel) Undelete(ctx context.Context, ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}

		var userID = auth.GetIdentityFromContext(ctx).Identity()

		if ch, err = svc.findByID(ctx, ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUndeleteChannel(ctx, ch) {
			return ChannelErrNotAllowedToUndelete()
		}

		if ch.DeletedAt == nil {
			return ChannelErrNotDeleted()
		}

		ch.DeletedAt = nil

		svc.scheduleSystemMessage(ctx, ch, "<@%d> undeleted this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		_ = svc.flushSystemMessages(ctx)
		return svc.sendChannelEvent(ctx, ch)
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionUndelete, err)
}

func (svc *channel) SetFlag(ctx context.Context, ID uint64, flag types.ChannelMembershipFlag) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{flag: string(flag)}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var membership *types.ChannelMember
		var userID = auth.GetIdentityFromContext(ctx).Identity()

		if ch, err = svc.FindByID(ctx, ID); err != nil {
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

		return svc.flushSystemMessages(ctx)

		// Setting a flag on a channel is a private thing,
		// no need to send channel event back to everyone
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionSetFlag, err)

}

func (svc *channel) Archive(ctx context.Context, ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var userID = auth.GetIdentityFromContext(ctx).Identity()

		if ch, err = svc.findByID(ctx, ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanArchiveChannel(ctx, ch) {
			return ChannelErrNotAllowedToUndelete()
		}

		if ch.ArchivedAt != nil {
			return ChannelErrAlreadyArchived()
		}

		ch.ArchivedAt = now()

		svc.scheduleSystemMessage(ctx, ch, "<@%d> archived this channel", userID)

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		_ = svc.flushSystemMessages(ctx)
		return svc.sendChannelEvent(ctx, ch)
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionArchive, err)
}

func (svc *channel) Unarchive(ctx context.Context, ID uint64) (ch *types.Channel, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return ChannelErrInvalidID()
		}
		var userID = auth.GetIdentityFromContext(ctx).Identity()

		if ch, err = svc.findByID(ctx, ID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if !svc.ac.CanUnarchiveChannel(ctx, ch) {
			return ErrNoPermissions.withStack()
		}

		if ch.ArchivedAt == nil {
			return ChannelErrNotArchived()
		}

		ch.ArchivedAt = nil

		if err = store.UpdateMessagingChannel(ctx, s, ch); err != nil {
			return
		}

		svc.scheduleSystemMessage(ctx, ch, "<@%d> unarchived this channel", userID)

		_ = svc.flushSystemMessages(ctx)
		return svc.sendChannelEvent(ctx, ch)
	})

	return ch, svc.recordAction(ctx, aProps, ChannelActionUnarchive, err)
}

func (svc *channel) InviteUser(ctx context.Context, channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if channelID == 0 {
			return ChannelErrInvalidID()
		}

		for _, memberID := range memberIDs {
			if memberID == 0 {
				return ChannelErrInvalidID()
			}
		}

		var (
			userID   = auth.GetIdentityFromContext(ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		out = types.ChannelMemberSet{}

		if ch, err = svc.FindByID(ctx, channelID); err != nil {
			return
		}

		aProps.setChannel(ch)

		if ch.Type == types.ChannelTypeGroup {
			return ChannelErrUnableToManageGroupMembers()
		}

		if !svc.ac.CanManageChannelMembers(ctx, ch) {
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

			svc.scheduleSystemMessage(ctx, ch, "<@%d> invited <@%d> to the channel", userID, memberID)

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

		return svc.flushSystemMessages(ctx)
	})

	return out, svc.recordAction(ctx, aProps, ChannelActionInviteMember, err)
}

func (svc *channel) AddMember(ctx context.Context, channelID uint64, memberIDs ...uint64) (out types.ChannelMemberSet, err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if channelID == 0 {
			return ChannelErrInvalidID()
		}

		for _, memberID := range memberIDs {
			if memberID == 0 {
				return ChannelErrInvalidID()
			}
		}

		var (
			userID   = auth.GetIdentityFromContext(ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		out = types.ChannelMemberSet{}

		if ch, err = svc.FindByID(ctx, channelID); err != nil {
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

			if memberID == userID && !svc.ac.CanJoinChannel(ctx, ch) {
				return ChannelErrNotAllowedToJoin()
			} else if memberID != userID && !svc.ac.CanManageChannelMembers(ctx, ch) {
				return ChannelErrNotAllowedToManageMembers()
			}

			if !exists {
				if userID == memberID {
					svc.scheduleSystemMessage(ctx, ch, "<@%d> joined", memberID)
				} else {
					svc.scheduleSystemMessage(ctx, ch, "<@%d> added <@%d> to the channel", userID, memberID)
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

			svc.event.Join(ctx, memberID, channelID)

			out = append(out, member)

			ch.Member = member
			ch.Members = out.AllMemberIDs()
		}

		// Push channel to all members
		if err = svc.sendChannelEvent(ctx, ch); err != nil {
			return
		}

		return svc.flushSystemMessages(ctx)
	})

	return out, svc.recordAction(ctx, aProps, ChannelActionAddMember, err)
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

func (svc *channel) DeleteMember(ctx context.Context, channelID uint64, memberIDs ...uint64) (err error) {
	var (
		aProps = &channelActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			userID   = auth.GetIdentityFromContext(ctx).Identity()
			ch       *types.Channel
			existing types.ChannelMemberSet
		)

		if ch, err = svc.FindByID(ctx, channelID); err != nil {
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

			if memberID == userID && !svc.ac.CanLeaveChannel(ctx, ch) {
				return ChannelErrNotAllowedToPart()
			} else if memberID != userID && !svc.ac.CanManageChannelMembers(ctx, ch) {
				return ChannelErrNotAllowedToManageMembers()
			}

			if userID == memberID {
				svc.scheduleSystemMessage(ctx, ch, "<@%d> left the channel", memberID)
			} else {
				svc.scheduleSystemMessage(ctx, ch, "<@%d> removed from the channel", memberID)
			}

			if err = store.DeleteMessagingChannelMemberByChannelIDUserID(ctx, s, channelID, memberID); err != nil {
				return err
			}

			_ = svc.event.Part(ctx, memberID, channelID)
		}

		return svc.flushSystemMessages(ctx)
	})

	return svc.recordAction(ctx, aProps, ChannelActionRemoveMember, err)

}

func (svc *channel) scheduleSystemMessage(ctx context.Context, ch *types.Channel, format string, a ...interface{}) {
	svc.sysmsgs = append(svc.sysmsgs, &types.Message{
		ChannelID: ch.ID,
		Message:   fmt.Sprintf(format, a...),
		Type:      types.MessageTypeChannelEvent,
	})
}

// Flushes sys message stack, stores them into repo & pushes them into event loop
func (svc *channel) flushSystemMessages(ctx context.Context) (err error) {
	defer func() {
		svc.sysmsgs = types.MessageSet{}
	}()

	return svc.sysmsgs.Walk(func(msg *types.Message) error {
		msg.ID = nextID()
		msg.CreatedAt = *now()

		if err = store.CreateMessagingMessage(ctx, svc.store, msg); err != nil {
			return err
		} else {
			return svc.event.Message(ctx, msg)
		}
	})
}

// Sends channel event
func (svc *channel) sendChannelEvent(ctx context.Context, ch *types.Channel) (err error) {
	if ch.DeletedAt == nil && ch.ArchivedAt == nil {
		// Looks like a valid channel

		// Preload members, if needed
		if len(ch.Members) == 0 || ch.Member == nil {
			f := types.ChannelMemberFilterChannels(ch.ID)
			if mm, _, err := store.SearchMessagingChannelMembers(ctx, svc.store, f); err != nil {
				return err
			} else {
				ch.Members = mm.AllMemberIDs()
				ch.Member = mm.FindByUserID(auth.GetIdentityFromContext(ctx).Identity())
			}
		}
	}

	if err = svc.setPermissionFlags(ctx, ch); err != nil {
		return
	}

	if err = svc.event.Channel(ctx, ch); err != nil {
		return
	}

	return nil
}

func (svc *channel) setPermissionFlags(ctx context.Context, ch *types.Channel) (err error) {
	ch.CanJoin = svc.ac.CanJoinChannel(ctx, ch)
	ch.CanPart = svc.ac.CanLeaveChannel(ctx, ch)
	ch.CanObserve = svc.ac.CanReadChannel(ctx, ch)
	ch.CanSendMessages = svc.ac.CanSendMessage(ctx, ch)

	ch.CanDeleteMessages = svc.ac.CanDeleteMessages(ctx, ch)
	ch.CanDeleteOwnMessages = svc.ac.CanDeleteOwnMessages(ctx, ch)
	ch.CanUpdateMessages = svc.ac.CanUpdateMessages(ctx, ch)
	ch.CanUpdateOwnMessages = svc.ac.CanUpdateOwnMessages(ctx, ch)
	ch.CanChangeMembers = svc.ac.CanManageChannelMembers(ctx, ch)
	// @todo migrate to proper change-membership-policy action check
	ch.CanChangeMembershipPolicy = svc.ac.CanChangeChannelMembershipPolicy(ctx, ch)

	ch.CanUpdate = svc.ac.CanUpdateChannel(ctx, ch)
	ch.CanArchive = svc.ac.CanArchiveChannel(ctx, ch)
	ch.CanUnarchive = svc.ac.CanUnarchiveChannel(ctx, ch)
	ch.CanDelete = svc.ac.CanDeleteChannel(ctx, ch)
	ch.CanUndelete = svc.ac.CanUndeleteChannel(ctx, ch)

	return nil
}

var _ ChannelService = &channel{}
