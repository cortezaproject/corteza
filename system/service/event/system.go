package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
)

// Match returns false if given conditions do not match event & resource internals
func (res systemOnInterval) Match(c eventbus.ConstraintMatcher) bool {
	return scheduler.OnInterval(c.Values()...)
}

// Match returns false if given conditions do not match event & resource internals
func (res systemOnTimestamp) Match(c eventbus.ConstraintMatcher) bool {
	return scheduler.OnTimestamp(c.Values()...)
}

// Match returns false if given conditions do not match event & resource internals
func (res systemBase) Match(c eventbus.ConstraintMatcher) bool {
	// No constraints are supported for system.
	return false
}
