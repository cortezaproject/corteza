package event

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = eventbus.ConstraintMaker

// Match returns false if given conditions do not match event & resource internals
func (res queueBase) Match(c eventbus.ConstraintMatcher) bool {
	return queueMatch(res.payload, c)
}

func queueMatch(r *types.QueueMessage, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "message.queue":
		return c.Match(r.Queue)
	}

	return false
}
