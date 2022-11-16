package event

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res moduleBase) Match(c eventbus.ConstraintMatcher) bool {
	return eventbus.MatchFirst(
		func() bool { return moduleMatch(res.module, c) },
		func() bool { return namespaceMatch(res.namespace, c) },
	)
}

// Handles module matchers
func moduleMatch(r *types.Module, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "module", "module.handle":
		return c.Match(r.Handle)
	case "module.name":
		return c.Match(r.Name)
	}

	return false
}
