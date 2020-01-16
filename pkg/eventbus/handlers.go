package eventbus

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	HandlerFn func(ctx context.Context, ev Event) error

	handler struct {
		handler       HandlerFn
		resourceTypes map[string]bool
		eventTypes    map[string]bool
		constraints   constraintSet
		weight        int
	}

	// @todo add sorting interface
	HandlerSet []*handler

	HandlerRegOp func(t *handler)

	eventInvokerSettable interface {
		SetInvoker(auth.Identifiable)
	}
)

// Match matches handler with resource event
func (t handler) Match(re Event) bool {
	if re == nil {
		return false
	}

	if len(re.ResourceType()) == 0 || !t.resourceTypes[re.ResourceType()] {
		// Expecting to have valid resource type and match at least one
		// defined resource on the handler
		return false
	}

	if len(re.EventType()) == 0 || !t.eventTypes[re.EventType()] {
		// Expecting to have valid event type and match at least one
		// defined event on the handler
		return false
	}

	for _, c := range t.constraints {
		// Should match all constraints
		if !re.Match(c) {
			return false
		}
	}

	return true
}

func (t handler) Handle(ctx context.Context, ev Event) error {
	defer sentry.Recover()

	if eis, ok := ev.(eventInvokerSettable); ok {
		eis.SetInvoker(auth.GetIdentityFromContext(ctx))
	}

	return t.handler(ctx, ev)
}

func NewHandler(h HandlerFn, ops ...HandlerRegOp) *handler {
	var t = &handler{
		resourceTypes: make(map[string]bool),
		eventTypes:    make(map[string]bool),
		handler:       h,
	}

	for _, op := range ops {
		op(t)
	}

	return t
}

func For(rr ...string) HandlerRegOp {
	return func(t *handler) {
		for _, r := range rr {
			t.resourceTypes[r] = true
		}
	}
}

func On(ee ...string) HandlerRegOp {
	return func(t *handler) {
		for _, e := range ee {
			t.eventTypes[e] = true
		}
	}
}

func Constraint(c ConstraintMatcher) HandlerRegOp {
	return func(t *handler) {
		t.constraints = append(t.constraints, c)
	}
}

func Weight(weight int) HandlerRegOp {
	return func(t *handler) {
		t.weight = weight
	}
}

// handler sorting:

func (set HandlerSet) Len() int           { return len(set) }
func (set HandlerSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set HandlerSet) Less(i, j int) bool { return set[i].weight < set[j].weight }
