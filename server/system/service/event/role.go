package event

import (
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/system/types"
)

// Match returns false if given conditions do not match event & resource internals
func (res roleBase) Match(c eventbus.ConstraintMatcher) bool {
	return roleMatch(res.role, c)
}

// Handles role matchers
func roleMatch(r *types.Role, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "role", "role.handle":
		return c.Match(r.Handle)
	case "role.name":
		return c.Match(r.Name)
	}

	return false
}
