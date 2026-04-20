package event

import "github.com/crusttech/human/server/pkg/eventbus"

// Match returns false if given conditions do not match event & resource internals
func (res authBase) Match(c eventbus.ConstraintMatcher) bool {
	return userMatch(res.user, c)
}
