package event

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res messageBase) Match(c eventbus.ConstraintMatcher) bool {
	return eventbus.MatchFirst(
		func() bool { return messageMatch(res.message, c) },
		func() bool { return channelMatch(res.channel, c) },
	)
}

// Handles message matchers
func messageMatch(r *types.Message, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "message":
		return c != nil && c.Match(r.Message)
	case "message.type":
		return c != nil && c.Match(string(r.Type))
	}

	return false
}
