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
			EventTypes: []string{"onTimestamp", onManualEventType, "onInterval"},
		}
	)

	a.Len(trg.EventTypes, 3)
	a.True(popOnManualEventType(trg))
	a.Len(trg.EventTypes, 2)
	a.False(popOnManualEventType(trg))
	a.Len(trg.EventTypes, 2)
}

func TestPluckManualTriggers(t *testing.T) {
	var (
		a = assert.New(t)

		s = &ServerScript{
			Triggers: []*Trigger{&Trigger{
				ResourceTypes: []string{"r1", "r2"},
				EventTypes:    []string{"onTimestamp", onManualEventType, "onInterval"},
			}},
		}
	)

	a.Len(s.Triggers[0].EventTypes, 3)
	a.EqualValues(
		map[string]bool{
			"r1": true,
			"r2": true,
		},
		pluckManualTriggers(s),
	)
	a.Len(s.Triggers[0].EventTypes, 2)

	// Running again must result in empty hash
	a.EqualValues(
		map[string]bool{},
		pluckManualTriggers(s),
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

	oo, err = makeTriggerOpts(trg)
	a.NoError(err)
	a.Len(oo, 4) // 1x all resources, 1x all events, 2x constraints

	oo, err = makeTriggerOpts(&Trigger{ResourceTypes: []string{"bar"}})
	a.Error(err, "expecting to fail on trigger w/o events")

	oo, err = makeTriggerOpts(&Trigger{EventTypes: []string{"foo"}})
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
					EventTypes:    []string{"ev"},
					ResourceTypes: []string{"res"},
				},
			},
		}

		s2 = &Script{
			Triggers: []*Trigger{
				&Trigger{
					EventTypes:    []string{"foo"},
					ResourceTypes: []string{"bar"},
				},
			},
		}

		s3 = &Script{
			Triggers: []*Trigger{
				&Trigger{
					EventTypes:    []string{"not-a-match"},
					ResourceTypes: []string{"not-a-match"},
				},
			},
		}

		a     = assert.New(t)
		strip = func(b bool, _ error) bool { return b }

		f func(s *Script) (b bool, err error)
	)

	f = makeScriptFilter(Filter{})
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.True(strip(f(s3)))

	f = makeScriptFilter(Filter{ResourceTypes: []string{"res"}})
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = makeScriptFilter(Filter{EventTypes: []string{"ev"}})
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = makeScriptFilter(Filter{EventTypes: []string{"ev", "foo"}})
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.False(strip(f(s3)))
}
