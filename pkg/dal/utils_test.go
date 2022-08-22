package dal

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestCompareValues(t *testing.T) {
	n := time.Now()
	m := n.Add(time.Second)

	tcc := []struct {
		name string
		a    any
		b    any
		out  int
	}{
		{
			name: "two ints; eq",
			a:    10,
			b:    10,
			out:  0,
		}, {
			name: "two ints; lt",
			a:    9,
			b:    10,
			out:  -1,
		}, {
			name: "two ints; gt",
			a:    10,
			b:    9,
			out:  1,
		},

		{
			name: "two strings; eq",
			a:    "aa",
			b:    "aa",
			out:  0,
		}, {
			name: "two strings; lt",
			a:    "a",
			b:    "aa",
			out:  -1,
		}, {
			name: "two strings; gt",
			a:    "aa",
			b:    "a",
			out:  1,
		},

		{
			name: "two uints; eq",
			a:    uint(10),
			b:    uint(10),
			out:  0,
		}, {
			name: "two uints; lt",
			a:    uint(9),
			b:    uint(10),
			out:  -1,
		}, {
			name: "two uints; gt",
			a:    uint(10),
			b:    uint(9),
			out:  1,
		},

		{
			name: "two floats; eq",
			a:    float64(10),
			b:    float64(10),
			out:  0,
		}, {
			name: "two floats; lt",
			a:    float64(9),
			b:    float64(10),
			out:  -1,
		}, {
			name: "two floats; gt",
			a:    float64(10),
			b:    float64(9),
			out:  1,
		},

		{
			name: "two times; eq",
			a:    n,
			b:    n,
			out:  0,
		}, {
			name: "two times; lt",
			a:    n,
			b:    m,
			out:  -1,
		}, {
			name: "two times; gt",
			a:    m,
			b:    n,
			out:  1,
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, compareValues(c.a, c.b))
		})
	}
}

func TestCompareGetters(t *testing.T) {
	tcc := []struct {
		name string
		a    ValueGetter
		b    ValueGetter
		attr string
		out  int
	}{{
		name: "eq",
		a:    simpleRow{"a": 10},
		b:    simpleRow{"a": 10},
		attr: "a",
		out:  0,
	}, {
		name: "lt",
		a:    simpleRow{"a": 9},
		b:    simpleRow{"a": 10},
		attr: "a",
		out:  -1,
	}, {
		name: "gt",
		a:    simpleRow{"a": 10},
		b:    simpleRow{"a": 9},
		attr: "a",
		out:  1,
	}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, compareGetters(c.a, c.b, c.a.CountValues(), c.b.CountValues(), c.attr))
		})
	}
}

func TestConstraintsToExpression(t *testing.T) {
	tcc := []struct {
		name string
		cc   map[string][]any
		out  []string
	}{{
		name: "one constraint, one value",
		cc: map[string][]any{
			"k1": {"v1"},
		},
		out: []string{`k1 == 'v1'`},
	}, {
		name: "one constraint, multiple values",
		cc: map[string][]any{
			"k1": {"v1", 10, true},
		},
		out: []string{`k1 == 'v1' || k1 == 10 || k1 == true`},
	}, {
		name: "multiple constraints, multiple values",
		cc: map[string][]any{
			"k1": {"v1", 10, true},
			"k2": {"v2", 42, false},
		},
		out: []string{
			`(k1 == 'v1' || k1 == 10 || k1 == true)`,
			`(k2 == 'v2' || k2 == 42 || k2 == false)`,
		},
	}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			// Do it like so since map order is not defined and the test would be flaky
			got := constraintsToExpression(c.cc)
			for _, o := range c.out {
				require.Contains(t, got, o)
			}
		})
	}
}

func TestStateConstraintsToExpression(t *testing.T) {
	tcc := []struct {
		name string
		cc   map[string]filter.State
		out  []string
	}{{
		name: "exclude these ones",
		cc: map[string]filter.State{
			"k1": filter.StateExcluded,
			"k2": filter.StateExcluded,
		},
		out: []string{`k1 == null`, `k2 == null`},
	}, {
		name: "only these ones",
		cc: map[string]filter.State{
			"k1": filter.StateExclusive,
			"k2": filter.StateExclusive,
		},
		out: []string{`k1 != null`, `k2 != null`},
	}, {
		name: "both",
		cc: map[string]filter.State{
			"k1": filter.StateExclusive,
			"k2": filter.StateExclusive,
		},
		out: []string{},
	}, {
		name: "mix and match",
		cc: map[string]filter.State{
			"k1": filter.StateExcluded,
			"k2": filter.StateInclusive,
			"k3": filter.StateExclusive,
		},
		out: []string{
			`k1 == null`,
			`k3 != null`,
		},
	}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			// Do it like so since map order is not defined and the test would be flaky
			got := stateConstraintsToExpression(c.cc)
			for _, o := range c.out {
				require.Contains(t, got, o)
			}
		})
	}
}
