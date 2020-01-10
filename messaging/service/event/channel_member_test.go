package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelMemberMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &channelMemberBase{
			channel: &types.Channel{Name: "ChanChan"},
			member:  &types.ChannelMember{Type: types.ChannelMembershipTypeOwner},
		}

		cOwn = eventbus.MustMakeConstraint("channel-member.type", "eq", "owner")
		cChn = eventbus.MustMakeConstraint("channel", "eq", "ChanChan")
	)

	a.True(channelMemberMatch(res.member, cOwn))

	a.True(res.Match(cOwn))
	a.True(res.Match(cChn))
}
