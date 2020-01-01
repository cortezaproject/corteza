package corredor

import (
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopOnManualEventType(t *testing.T) {
	var (
		a = assert.New(t)

		trg = &Trigger{
			Events: []string{"onTimestamp", onManualEventType, "onInterval"},
		}
	)

	a.Len(trg.Events, 3)
	a.True(popOnManualEventType(trg))
	a.Len(trg.Events, 2)
	a.False(popOnManualEventType(trg))
	a.Len(trg.Events, 2)
}

func TestPluckManualTriggers(t *testing.T) {
	var (
		a = assert.New(t)

		s = &ServerScript{
			Triggers: []*Trigger{&Trigger{
				Resources: []string{"r1", "r2"},
				Events:    []string{"onTimestamp", onManualEventType, "onInterval"},
				RunAs:     "moi",
			}},
		}
	)

	a.Len(s.Triggers[0].Events, 3)
	a.EqualValues(
		map[string]string{
			"r1": "moi",
			"r2": "moi",
		},
		pluckManualTriggers(s),
	)
	a.Len(s.Triggers[0].Events, 2)

	// Running again must result in empty hash
	a.EqualValues(
		map[string]string{},
		pluckManualTriggers(s),
	)
}

func TestTriggerOptsMaking(t *testing.T) {
	var (
		a = assert.New(t)

		trg = &Trigger{
			Resources: []string{"r1", "r2"},
			Events:    []string{"onTimestamp", onManualEventType, "onInterval"},
			Constraints: []*TConstraint{
				&TConstraint{Name: "some1", Op: "eq", Value: []string{"other"}},
				&TConstraint{Name: "some2", Op: "eq", Value: []string{"other"}},
			},
		}

		oo  []eventbus.TriggerRegOp
		err error
	)

	oo, err = makeTriggerOpts(trg)
	a.NoError(err)
	a.Len(oo, 4) // 1x all resources, 1x all events, 2x constraints

	oo, err = makeTriggerOpts(&Trigger{Resources: []string{"bar"}})
	a.Error(err, "expecting to fail on trigger w/o events")

	oo, err = makeTriggerOpts(&Trigger{Events: []string{"foo"}})
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

func TestScriptFilterMaker(t *testing.T) {
	var (
		s1 = &Script{
			Triggers: []*Trigger{
				&Trigger{
					Events:    []string{"ev"},
					Resources: []string{"res"},
				},
			},
		}

		s2 = &Script{
			Triggers: []*Trigger{
				&Trigger{
					Events:    []string{"foo"},
					Resources: []string{"bar"},
				},
			},
		}

		s3 = &Script{
			Triggers: []*Trigger{
				&Trigger{
					Events:    []string{"not-a-match"},
					Resources: []string{"not-a-match"},
				},
			},
		}

		a     = assert.New(t)
		strip = func(b bool, _ error) bool { return b }

		f func(s *Script) (b bool, err error)
	)

	f = makeScriptFilter(ManualScriptFilter{})
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.True(strip(f(s3)))

	f = makeScriptFilter(ManualScriptFilter{ResourceTypes: []string{"res"}})
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = makeScriptFilter(ManualScriptFilter{EventTypes: []string{"ev"}})
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = makeScriptFilter(ManualScriptFilter{EventTypes: []string{"ev", "foo"}})
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.False(strip(f(s3)))
}
