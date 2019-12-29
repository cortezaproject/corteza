package eventbus

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	mockEvent struct {
		rType string
		eType string
		match func(name string, op string, values ...string) bool

		identity auth.Identifiable
	}
)

func (e mockEvent) ResourceType() string {
	return e.rType
}

func (e mockEvent) EventType() string {
	return e.eType
}

func (e mockEvent) Match(name string, op string, values ...string) bool {
	if e.match == nil {
		return true
	}

	return e.match(name, op, values...)
}

func (e *mockEvent) SetInvoker(identity auth.Identifiable) {
	e.identity = identity
}

func TestTrigger_Match(t *testing.T) {
	cases := []struct {
		name  string
		ops   []TriggerRegOp
		ev    Event
		match bool
	}{
		{"nil event",
			nil,
			nil,
			false,
		},
		{"empty resource",
			[]TriggerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "", eType: "bar"},
			false,
		},
		{"empty event",
			[]TriggerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "foo", eType: ""},
			false,
		},
		{"simple foo-bar test",
			[]TriggerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "foo", eType: "bar"},
			true,
		},
		{"constraint match",
			[]TriggerRegOp{For("foo"), On("bar"), Constraint("baz", "", "baz")},
			&mockEvent{
				rType: "foo",
				eType: "bar",
				match: func(name string, op string, values ...string) bool {
					return len(values) > 0 && name == values[0]
				}},
			true,
		},
		{"constraint mismatch",
			[]TriggerRegOp{For("foo"), On("bar"), Constraint("baz", "", "baz")},
			&mockEvent{
				rType: "foo",
				eType: "bar",
				match: func(name string, op string, values ...string) bool {
					return false
				}},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var trigger = NewTrigger(nil, c.ops...)
			if c.match {
				assert.True(t, trigger.Match(c.ev), "Expecting to match")

			} else {
				assert.False(t, trigger.Match(c.ev), "Expecting to not match")

			}
		})
	}
}

func TestTrigger_RegOps(t *testing.T) {
	makeTestTrigger := func(t *trigger) *trigger {
		if t.resourceTypes == nil {
			t.resourceTypes = make(map[string]bool)
		}

		if t.eventTypes == nil {
			t.eventTypes = make(map[string]bool)
		}

		return t
	}

	cases := []struct {
		name string
		exp  *trigger
		ops  []TriggerRegOp
	}{
		{
			"empty",
			makeTestTrigger(&trigger{}),
			nil,
		},
		{
			"resource types",
			makeTestTrigger(&trigger{resourceTypes: map[string]bool{"foo": true, "bar": true}}),
			[]TriggerRegOp{For("foo", "bar")},
		},
		{
			"event types",
			makeTestTrigger(&trigger{eventTypes: map[string]bool{"foo": true, "bar": true}}),
			[]TriggerRegOp{On("foo", "bar")},
		},
		{
			"weight",
			makeTestTrigger(&trigger{weight: 42}),
			[]TriggerRegOp{Weight(42)},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.exp, NewTrigger(nil, c.ops...))
		})
	}
}

func TestTriggerHandler(t *testing.T) {
	var (
		a             = assert.New(t)
		ctx           = context.Background()
		ev            = &mockEvent{}
		passedthrough bool

		trSimple = &trigger{
			handler: func(ctx context.Context, ev Event) error {
				passedthrough = true
				a.True(auth.IsSuperUser(ev.(*mockEvent).identity))
				return nil
			},
		}
	)

	a.False(passedthrough)
	a.False(auth.IsSuperUser(ev.identity))

	trSimple.Handle(auth.SetSuperUserContext(ctx), ev)

	a.True(auth.IsSuperUser(ev.identity))
	a.True(passedthrough, "expecting to pass through simple handler")
}

func TestTriggerSorting(t *testing.T) {
	var (
		a  = assert.New(t)
		tt = TriggerSet{
			NewTrigger(nil, Weight(3)),
			NewTrigger(nil, Weight(1)),
			NewTrigger(nil, Weight(2)),
		}

		w2s = func(tt TriggerSet) (out string) {
			for _, t := range tt {
				out += fmt.Sprintf("%d,", t.weight)
			}
			return
		}
	)

	a.Equal(w2s(tt), "3,1,2,")
	sort.Sort(tt)
	a.Equal(w2s(tt), "1,2,3,")
}
