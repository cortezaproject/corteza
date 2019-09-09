package messaging

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func (h helper) repoChannel() repository.ChannelRepository {
	var (
		ctx = context.Background()
		db  = factory.Database.MustGet("messaging").With(ctx)
	)

	return repository.Channel(ctx, db)
}

func (h helper) repoChMember() repository.ChannelMemberRepository {
	var (
		ctx = context.Background()
		db  = factory.Database.MustGet("messaging").With(ctx)
	)

	return repository.ChannelMember(ctx, db)
}

func (h helper) repoMakePublicCh() *types.Channel {
	ch, err := h.repoChannel().Create(&types.Channel{
		Name: "Test channel " + time.Now().String(),
		Type: types.ChannelTypePublic,
	})

	h.a.NoError(err)
	return ch
}

func (h helper) repoMakePrivateCh() *types.Channel {
	ch, err := h.repoChannel().Create(&types.Channel{
		Name: "Test channel " + time.Now().String(),
		Type: types.ChannelTypePrivate,
	})

	h.a.NoError(err)
	return ch
}

func (h helper) repoMakeMember(ch *types.Channel, u *sysTypes.User) *types.ChannelMember {
	m, err := h.
		repoChMember().
		Create(&types.ChannelMember{ChannelID: ch.ID, UserID: h.cUser.ID, Type: types.ChannelMembershipTypeMember})
	h.a.NoError(err)

	return m
}

func (h helper) repoChAssertNotMember(ch *types.Channel, u *sysTypes.User) {
	mm, err := h.repoChMember().Find(&types.ChannelMemberFilter{ChannelID: ch.ID, MemberID: h.cUser.ID})
	h.a.NoError(err)
	h.a.NotNil(mm)
	h.a.NotContains(mm.AllMemberIDs(), u.ID, "not expecting to find a member")
}

func (h helper) repoChAssertMember(ch *types.Channel, u *sysTypes.User, typ types.ChannelMembershipType) {
	mm, err := h.repoChMember().Find(&types.ChannelMemberFilter{ChannelID: ch.ID, MemberID: h.cUser.ID})

	h.a.NoError(err)
	h.a.NotNil(mm)
	h.a.NotNil(mm.FindByUserID(u.ID), "expecting to find a member")
	h.a.Equal(typ, mm.FindByUserID(u.ID).Type, "expecting to find a member")
}
