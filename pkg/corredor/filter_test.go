package corredor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterResourcePrefixing(t *testing.T) {
	f := &Filter{
		ResourceTypePrefixes: []string{
			"system",
			"system:one",
			"two",
			"ui:admin:foo",
			"ui:admin",
		},
	}

	f.procRTPrefixes("system")

	assert.New(t).Equal(
		f.ResourceTypePrefixes,
		[]string{
			"system",
			"system:one",
			"system:two",
			"ui:admin:foo",
			"ui:admin",
		},
	)

	f = &Filter{}

	f.procRTPrefixes("system")

	assert.New(t).Equal(
		f.ResourceTypePrefixes,
		[]string{
			"system",
			"ui:",
		},
	)
}

func TestResourceTypesChecker(t *testing.T) {
	var (
		a = assert.New(t)
	)

	// nothing requested in the filter
	a.True(Filter{ResourceTypes: []string{}}.checkResourceTypes("c"))

	// no resources
	a.False(Filter{ResourceTypes: []string{"a", "b"}}.checkResourceTypes())

	// no resources, no filter
	a.True(Filter{}.checkResourceTypes())

	a.True(Filter{ResourceTypes: []string{"a"}}.checkResourceTypes("a", "b"))

	a.False(Filter{ResourceTypes: []string{"c"}}.checkResourceTypes("a", "b"))
}

func TestEventTypesChecker(t *testing.T) {
	var (
		a = assert.New(t)
	)

	// nothing requested in the filter
	a.True(Filter{EventTypes: []string{}}.checkEventTypes("c"))

	// no resources
	a.False(Filter{EventTypes: []string{"a", "b"}}.checkEventTypes())

	// no resources, no filter
	a.True(Filter{}.checkEventTypes())

	a.True(Filter{EventTypes: []string{"a"}}.checkEventTypes("a", "b"))

	a.False(Filter{EventTypes: []string{"c"}}.checkEventTypes("a", "b"))
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

	f = Filter{}.makeFilterFn()
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.True(strip(f(s3)))

	f = Filter{ResourceTypes: []string{"res"}}.makeFilterFn()
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = Filter{EventTypes: []string{"ev"}}.makeFilterFn()
	a.True(strip(f(s1)))
	a.False(strip(f(s2)))
	a.False(strip(f(s3)))

	f = Filter{EventTypes: []string{"ev", "foo"}}.makeFilterFn()
	a.True(strip(f(s1)))
	a.True(strip(f(s2)))
	a.False(strip(f(s3)))
}
