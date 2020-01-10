package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res roleMemberBase) Match(c eventbus.ConstraintMatcher) bool {
	return userMatch(res.user, c, roleMatch(res.role, c, false))
}
