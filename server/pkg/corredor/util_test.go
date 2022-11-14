package corredor

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapManualTriggers(t *testing.T) {
	var (
		a = assert.New(t)
	)

	a.EqualValues(
		map[string]bool{
			"r1": true,
			"r2": true,
		},
		mapExplicitTriggers(&ServerScript{
			Triggers: []*Trigger{&Trigger{
				ResourceTypes: []string{"r1", "r2"},
				EventTypes:    []string{"onTimestamp", onManualEventType, "onInterval"},
			}},
		}),
	)

	a.EqualValues(
		map[string]bool{},
		mapExplicitTriggers(&ServerScript{
			Triggers: []*Trigger{&Trigger{
				ResourceTypes: []string{"r1", "r2"},
				EventTypes:    []string{"onTimestamp", "onInterval"},
			}},
		}),
	)
}

func TestTriggerOptsMaking(t *testing.T) {
	var (
		a = assert.New(t)

		trg = &Trigger{
			ResourceTypes: []string{"r1", "r2"},
			EventTypes:    []string{"onTimestamp", onManualEventType, "onInterval"},
			Constraints: []*TConstraint{
				&TConstraint{Name: "some1", Op: "eq", Value: []string{"other"}},
				&TConstraint{Name: "some2", Op: "eq", Value: []string{"other"}},
			},
		}

		oo  []eventbus.HandlerRegOp
		err error
	)

	oo, err = triggerToHandlerOps(trg)
	a.NoError(err)
	a.Len(oo, 4) // 1x all resources, 1x all events, 2x constraints

	oo, err = triggerToHandlerOps(&Trigger{ResourceTypes: []string{"bar"}})
	a.Error(err, "expecting to fail on trigger w/o events")

	oo, err = triggerToHandlerOps(&Trigger{EventTypes: []string{"foo"}})
	a.Error(err, "expecting to fail on trigger w/o resources")
}

func TestArgEncoding(t *testing.T) {
	var (
		a = assert.New(t)

		args = map[string]string{}
	)

	a.NoError(encodeArguments(args, "foo", "string"))
	a.NoError(encodeArguments(args, "bar", 42))
	a.NoError(encodeArguments(args, "baz", true))
	a.NoError(encodeArguments(args, "obj", struct{ A string }{A: "A"}))
	a.Error(encodeArguments(args, "func", func() {}))

	a.EqualValues(map[string]string{
		"foo": `"string"`,
		"bar": `42`,
		"baz": `true`,
		"obj": `{"A":"A"}`,
	}, args)
}
