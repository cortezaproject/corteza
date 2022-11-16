package event

import (
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/system/types"
)

// Match returns false if given conditions do not match event & resource internals
func (res userBase) Match(c eventbus.ConstraintMatcher) bool {
	return userMatch(res.user, c)
}

// Handles user matchers
func userMatch(r *types.User, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "user", "user.handle":
		return r != nil && c.Match(r.Handle)
	case "user.email":
		return r != nil && c.Match(r.Email)
	}

	return false
}
