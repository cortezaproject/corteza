package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res roleMemberBase) Match(c eventbus.ConstraintMatcher) bool {
	return eventbus.MatchFirst(
		func() bool { return userMatch(res.user, c) },
		func() bool { return roleMatch(res.role, c) },
	)
}
