package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res channelBase) Match(c eventbus.ConstraintMatcher) bool {
	return channelMatch(res.channel, c)
}

// Handles channel matchers
//
// This *match() fn uses a 3rd param to allow matcher chaining (see commands, channels)
func channelMatch(r *types.Channel, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "channel", "channel.name":
		return r != nil && c.Match(r.Name)
	case "channel.topic":
		return r != nil && c.Match(r.Topic)
	case "channel.type":
		return r != nil && c.Match(r.Type.String())
	}

	return false
}
