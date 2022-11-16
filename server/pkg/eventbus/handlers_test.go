package eventbus

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cortezaproject/corteza/server/pkg/auth"
)

type (
	mockEvent struct {
		rType string
		eType string
		match func(c ConstraintMatcher) bool

		identity auth.Identifiable
	}
)

func (e mockEvent) ResourceType() string {
	return e.rType
}

func (e mockEvent) EventType() string {
	return e.eType
}

func (e mockEvent) Match(c ConstraintMatcher) bool {
	if e.match == nil {
		return true
	}

	return e.match(c)
}

func (e *mockEvent) SetInvoker(identity auth.Identifiable) {
	e.identity = identity
}

func TestHandler_Match(t *testing.T) {
	cases := []struct {
		name  string
		ops   []HandlerRegOp
		ev    Event
		match bool
	}{
		{"nil event",
			nil,
			nil,
			false,
		},
		{"empty resource",
			[]HandlerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "", eType: "bar"},
			false,
		},
		{"empty event",
			[]HandlerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "foo", eType: ""},
			false,
		},
		{"simple foo-bar test",
			[]HandlerRegOp{For("foo"), On("bar")},
			&mockEvent{rType: "foo", eType: "bar"},
			true,
		},
		{"ConstraintMatcher match",
			[]HandlerRegOp{For("foo"), On("bar"), Constraint(MustMakeConstraint("baz", "=", "baz"))},
			&mockEvent{
				rType: "foo",
				eType: "bar",
				match: func(c ConstraintMatcher) bool {
					return len(c.Values()) > 0 && c.Name() == c.Values()[0]
				}},
			true,
		},
		{"ConstraintMatcher mismatch",
			[]HandlerRegOp{For("foo"), On("bar"), Constraint(MustMakeConstraint("baz", "=", "baz"))},
			&mockEvent{
				rType: "foo",
				eType: "bar",
				match: func(c ConstraintMatcher) bool {
					return false
				}},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var handler = NewHandler(nil, c.ops...)
			if c.match {
				assert.True(t, handler.Match(c.ev), "Expecting to match")

			} else {
				assert.False(t, handler.Match(c.ev), "Expecting to not match")

			}
		})
	}
}

func TestHandler_RegOps(t *testing.T) {
	makeTestHandler := func(t *handler) *handler {
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
		exp  *handler
		ops  []HandlerRegOp
	}{
		{
			"empty",
			makeTestHandler(&handler{}),
			nil,
		},
		{
			"resource types",
			makeTestHandler(&handler{resourceTypes: map[string]bool{"foo": true, "bar": true}}),
			[]HandlerRegOp{For("foo", "bar")},
		},
		{
			"event types",
			makeTestHandler(&handler{eventTypes: map[string]bool{"foo": true, "bar": true}}),
			[]HandlerRegOp{On("foo", "bar")},
		},
		{
			"weight",
			makeTestHandler(&handler{weight: 42}),
			[]HandlerRegOp{Weight(42)},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.exp, NewHandler(nil, c.ops...))
		})
	}
}

func TestHandlerHandler(t *testing.T) {
	var (
		identity uint64 = 42

		a             = assert.New(t)
		ctx           = context.Background()
		ev            = &mockEvent{}
		passedthrough bool

		trSimple = &handler{
			handler: func(ctx context.Context, ev Event) error {
				passedthrough = true
				a.Equal(identity, ev.(*mockEvent).identity.Identity())
				return nil
			},
		}
	)

	a.False(passedthrough)
	a.Nil(ev.identity)

	a.NoError(trSimple.Handle(auth.SetIdentityToContext(ctx, auth.Authenticated(identity)), ev))

	a.NotNil(ev.identity)
	a.Equal(identity, ev.identity.Identity())
	a.True(passedthrough, "expecting to pass through simple handler")
}

func TestHandlerSorting(t *testing.T) {
	var (
		a  = assert.New(t)
		tt = HandlerSet{
			NewHandler(nil, Weight(3)),
			NewHandler(nil, Weight(1)),
			NewHandler(nil, Weight(2)),
		}

		w2s = func(tt HandlerSet) (out string) {
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
