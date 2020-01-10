package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
)

// Match returns false if given conditions do not match event & resource internals
func (res messagingOnInterval) Match(c eventbus.ConstraintMatcher) bool {
	// @todo this could be flippled
	//       instead of passing around raw values as strings,
	//       constraint values could/should be preparsed when creating constraint
	return scheduler.OnInterval(c.Values()...)
}

// Match returns false if given conditions do not match event & resource internals
func (res messagingOnTimestamp) Match(c eventbus.ConstraintMatcher) bool {
	// @todo this could be flippled
	//       instead of passing around raw values as strings,
	//       constraint values could/should be preparsed when creating constraint
	return scheduler.OnTimestamp(c.Values()...)
}

// Match returns false if given conditions do not match event & resource internals
func (res messagingBase) Match(c eventbus.ConstraintMatcher) bool {
	// No constraints are supported for messaging.
	return false
}
