package event

import "github.com/cortezaproject/corteza-server/pkg/eventbus"

// Match returns false if given conditions do not match event & resource internals
func (res recordBase) Match(c eventbus.ConstraintMatcher) bool {
	return namespaceMatch(res.namespace, c, moduleMatch(res.module, c, false))
}
