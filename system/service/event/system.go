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
	// By default we match no mather what kind of constraints we receive
	//
	// Function will be called multiple times - once for every trigger constraint
	// All should match (return true):
	//   constraint#1 AND constraint#2 AND constraint#3 ...
	//
	// When there are multiple values, Match() can decide how to treat them (OR, AND...)
	return true
}
