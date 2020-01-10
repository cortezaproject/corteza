package event

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

// Match returns false if given conditions do not match event & resource internals
func (res pageBase) Match(c eventbus.ConstraintMatcher) bool {
	return namespaceMatch(res.namespace, c, pageMatch(res.page, c, false))
}

// Handles namespace matchers
func pageMatch(r *types.Page, c eventbus.ConstraintMatcher, def bool) bool {
	switch c.Name() {
	case "page", "page.handle":
		return c.Match(r.Handle)
	case "page.title":
		return c.Match(r.Title)
	}

	return def
}
