package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Match returns false if given conditions do not match event & resource internals
func (res sinkBase) Match(c eventbus.ConstraintMatcher) bool {
	if !sinkMatch(res.request, c) {
		return false
	}

	return false
}

// Handles sink matchers
func sinkMatch(r *types.SinkRequest, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "url", "request.url":
		return c.Match(r.RequestURL)
	}

	return true
}
