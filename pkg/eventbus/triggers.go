package eventbus

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	constraint struct {
		name  string
		op    string
		value []string
	}

	constraintSet []constraint

	Handler func(ctx context.Context, ev Event) error

	trigger struct {
		handler       Handler
		resourceTypes map[string]bool
		eventTypes    map[string]bool
		constraints   constraintSet
		weight        int
	}

	// @todo add sorting interface
	TriggerSet []*trigger

	TriggerRegOp func(t *trigger)

	eventInvokerSettable interface {
		SetInvoker(auth.Identifiable)
	}
)

// Match matches trigger with resource event
func (t trigger) Match(re Event) bool {
	if re == nil {
		return false
	}

	if len(re.ResourceType()) == 0 || !t.resourceTypes[re.ResourceType()] {
		// Expecting to have valid resource type and match at least one
		// defined resource on the trigger
		return false
	}

	if len(re.EventType()) == 0 || !t.eventTypes[re.EventType()] {
		// Expecting to have valid event type and match at least one
		// defined event on the trigger
		return false
	}

	for _, c := range t.constraints {
		// Should match all constraints
		if !re.Match(c.name, c.op, c.value...) {
			return false
		}
	}

	return true
}

func (t trigger) Handle(ctx context.Context, ev Event) error {
	defer sentry.Recover()

	if eis, ok := ev.(eventInvokerSettable); ok {
		eis.SetInvoker(auth.GetIdentityFromContext(ctx))
	}

	return t.handler(ctx, ev)
}

func NewTrigger(h Handler, ops ...TriggerRegOp) *trigger {
	var t = &trigger{
		resourceTypes: make(map[string]bool),
		eventTypes:    make(map[string]bool),
		handler:       h,
	}

	for _, op := range ops {
		op(t)
	}

	return t
}

func For(rr ...string) TriggerRegOp {
	return func(t *trigger) {
		for _, r := range rr {
			t.resourceTypes[r] = true
		}
	}
}

func On(ee ...string) TriggerRegOp {
	return func(t *trigger) {
		for _, e := range ee {
			t.eventTypes[e] = true
		}
	}
}

func Constraint(name, op string, vv ...string) TriggerRegOp {
	return func(t *trigger) {
		t.constraints = append(t.constraints, constraint{
			name:  name,
			op:    op,
			value: vv,
		})
	}
}

func Weight(weight int) TriggerRegOp {
	return func(t *trigger) {
		t.weight = weight
	}
}

// Trigger-set sorting:

func (set TriggerSet) Len() int           { return len(set) }
func (set TriggerSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set TriggerSet) Less(i, j int) bool { return set[i].weight < set[j].weight }
