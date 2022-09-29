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
		}, {
			name: "two ints; a nil",
			a:    nil,
			b:    10,
			out:  -1,
		}, {
			name: "two ints; b nil",
			a:    10,
			b:    nil,
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
		}, {
			name: "two strings; a nil",
			a:    nil,
			b:    "aa",
			out:  -1,
		}, {
			name: "two strings; b nil",
			a:    "aa",
			b:    nil,
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
		}, {
			name: "two uints; a nil",
			a:    nil,
			b:    uint(10),
			out:  -1,
		}, {
			name: "two uints; b nil",
			a:    uint(10),
			b:    nil,
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
		}, {
			name: "two floats; a nil",
			a:    nil,
			b:    float64(10),
			out:  -1,
		}, {
			name: "two floats; b nil",
			a:    float64(10),
			b:    nil,
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
		}, {
			name: "two times; a nil",
			a:    nil,
			b:    n,
			out:  -1,
		}, {
			name: "two times; b nil",
			a:    n,
			b:    nil,
			out:  1,
		},

		{
			name: "two nils",
			a:    nil,
			b:    nil,
			out:  0,
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
		name: "single eq",
		a:    simpleRow{"a": 10},
		b:    simpleRow{"a": 10},
		attr: "a",
		out:  0,
	}, {
		name: "single lt",
		a:    simpleRow{"a": 9},
		b:    simpleRow{"a": 10},
		attr: "a",
		out:  -1,
	}, {
		name: "single gt",
		a:    simpleRow{"a": 10},
		b:    simpleRow{"a": 9},
		attr: "a",
		out:  1,
	},

		{
			name: "multi eq both empty",
			a:    &Row{},
			b:    &Row{},
			attr: "a",
			out:  0,
		}, {
			name: "multi eq same values",
			a:    (&Row{}).WithValue("a", 0, 1),
			b:    (&Row{}).WithValue("a", 0, 1),
			attr: "a",
			out:  0,
		}, {
			name: "multi lt a less items",
			a:    (&Row{}),
			b:    (&Row{}).WithValue("a", 0, 1),
			attr: "a",
			out:  -1,
		}, {
			name: "multi lt a item less",
			a:    (&Row{}).WithValue("a", 0, 0),
			b:    (&Row{}).WithValue("a", 0, 1),
			attr: "a",
			out:  -1,
		}, {
			name: "multi gt a more items",
			a:    (&Row{}).WithValue("a", 0, 1),
			b:    (&Row{}),
			attr: "a",
			out:  1,
		}, {
			name: "multi gt a item more",
			a:    (&Row{}).WithValue("a", 0, 1),
			b:    (&Row{}).WithValue("a", 0, 0),
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

func TestRowComparator(t *testing.T) {
	// @todo add some more extreme cases
	tcc := []struct {
		name string
		a    ValueGetter
		b    ValueGetter
		ss   filter.SortExprSet
		less bool
	}{
		// Simple one col cases
		{
			name: "single column simple asc less",
			a:    simpleRow{"a": 10},
			b:    simpleRow{"a": 20},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}},
			less: true,
		},
		{
			name: "single column simple asc more",
			a:    simpleRow{"a": 20},
			b:    simpleRow{"a": 10},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}},
			less: false,
		},
		{
			name: "single column simple asc equal",
			a:    simpleRow{"a": 10},
			b:    simpleRow{"a": 10},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}},
			less: false,
		},
		{
			name: "single column simple desc less",
			a:    simpleRow{"a": 20},
			b:    simpleRow{"a": 10},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}},
			less: true,
		},
		{
			name: "single column simple desc more",
			a:    simpleRow{"a": 10},
			b:    simpleRow{"a": 20},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}},
			less: false,
		},
		{
			name: "single column simple desc equal",
			a:    simpleRow{"a": 10},
			b:    simpleRow{"a": 10},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}},
			less: false,
		},

		// basic 2 col cases
		{
			name: "two column asc less first priority",
			a:    simpleRow{"a": 10, "b": 100},
			b:    simpleRow{"a": 20, "b": 1},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}, {Column: "b", Descending: false}},
			less: true,
		},
		{
			name: "two column asc less first equal",
			a:    simpleRow{"a": 10, "b": 1},
			b:    simpleRow{"a": 10, "b": 2},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}, {Column: "b", Descending: false}},
			less: true,
		},
		{
			name: "two column asc equal",
			a:    simpleRow{"a": 10, "b": 1},
			b:    simpleRow{"a": 10, "b": 1},
			ss:   filter.SortExprSet{{Column: "a", Descending: false}, {Column: "b", Descending: false}},
			less: false,
		},
		{
			name: "two column desc less first priority",
			a:    simpleRow{"a": 20, "b": 1},
			b:    simpleRow{"a": 10, "b": 100},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}, {Column: "b", Descending: true}},
			less: true,
		},
		{
			name: "two column desc less first equal",
			a:    simpleRow{"a": 10, "b": 2},
			b:    simpleRow{"a": 10, "b": 1},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}, {Column: "b", Descending: true}},
			less: true,
		},
		{
			name: "two column desc equal",
			a:    simpleRow{"a": 10, "b": 1},
			b:    simpleRow{"a": 10, "b": 1},
			ss:   filter.SortExprSet{{Column: "a", Descending: true}, {Column: "b", Descending: true}},
			less: false,
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {

			less := makeRowComparator(c.ss...)(c.a, c.b)
			require.Equal(t, c.less, less)
		})
	}
}

func TestRowResetting(t *testing.T) {
	r := &Row{}

	gv := func(ident string, pos uint) any {
		v, err := r.GetValue(ident, pos)
		require.NoError(t, err)
		return v
	}

	r.SetValue("a", 0, 1)
	r.SetValue("a", 0, 2)
	r.SetValue("a", 1, 1)

	require.Equal(t, 2, gv("a", 0))
	require.Equal(t, 1, gv("a", 1))
	require.Equal(t, uint(2), r.counters["a"])

	r.Reset()
	require.Equal(t, uint(0), r.counters["a"])

	r.SetValue("a", 0, 3)
	r.SetValue("b", 0, 4)
	require.Equal(t, 3, gv("a", 0))
	require.Equal(t, uint(1), r.counters["a"])
	require.Equal(t, 4, gv("b", 0))
	require.Equal(t, uint(1), r.counters["b"])
}
