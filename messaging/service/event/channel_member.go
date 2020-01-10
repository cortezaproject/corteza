package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res channelMemberBase) Match(c eventbus.ConstraintMatcher) bool {
	return channelMatch(res.channel, c, channelMemberMatch(res.member, c))
}

// Handles channel member matchers
func channelMemberMatch(r *types.ChannelMember, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "channel-member.type":
		return c.Match(string(r.Type))
	}

	return false
}
