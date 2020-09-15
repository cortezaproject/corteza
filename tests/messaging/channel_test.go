package messaging

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"time"
)

func (h helper) repoMakePublicCh() *types.Channel {
	ch := &types.Channel{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Name:      "Test channel " + time.Now().String(),
		Type:      types.ChannelTypePublic,
	}

	h.a.NoError(store.CreateMessagingChannel(context.Background(), service.DefaultStore, ch))
	return ch
}

func (h helper) repoMakePrivateCh() *types.Channel {
	ch := &types.Channel{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Name:      "Test channel " + time.Now().String(),
		Type:      types.ChannelTypePrivate,
	}

	h.a.NoError(store.CreateMessagingChannel(context.Background(), service.DefaultStore, ch))
	return ch
}

func (h helper) repoMakeMember(ch *types.Channel, u *sysTypes.User) *types.ChannelMember {
	m := &types.ChannelMember{
		CreatedAt: time.Now(),
		ChannelID: ch.ID,
		UserID:    u.ID,
		Type:      types.ChannelMembershipTypeMember,
	}
	h.a.NoError(store.CreateMessagingChannelMember(context.Background(), service.DefaultStore, m))
	return m
}

func (h helper) lookupChMembership(ch *types.Channel) types.ChannelMemberSet {
	f := types.ChannelMemberFilter{ChannelID: []uint64{ch.ID}}

	mm, _, err := store.SearchMessagingChannelMembers(context.Background(), service.DefaultStore, f)
	h.a.NoError(err)
	return mm
}

func (h helper) repoChAssertNotMember(ch *types.Channel, u *sysTypes.User) {
	h.a.NotContains(h.lookupChMembership(ch).AllMemberIDs(), u.ID, "not expecting to find a member")
}

func (h helper) repoChAssertMember(ch *types.Channel, u *sysTypes.User, typ types.ChannelMembershipType) {
	mm := h.lookupChMembership(ch)
	h.a.NotNil(mm.FindByUserID(u.ID), "expecting to find a member")
	h.a.Equal(typ, mm.FindByUserID(u.ID).Type, "expecting to find a member")
}
