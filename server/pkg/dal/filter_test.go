package dal

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestInternalFilterMergeFilters_expression(t *testing.T) {
	t.Run("b no expr", func(t *testing.T) {
		a, err := toInternalFilter(internalFilter{expression: "a == b"})
		require.NoError(t, err)
		b, err := toInternalFilter(internalFilter{expression: ""})
		require.NoError(t, err)

		c := a.mergeFilters(b)
		require.Equal(t, "a == b", c.expression)
		require.Equal(t, "eq(a, b)", c.expParsed.String())
	})

	t.Run("a no expr", func(t *testing.T) {
		a, err := toInternalFilter(internalFilter{expression: ""})
		require.NoError(t, err)
		b, err := toInternalFilter(internalFilter{expression: "a == b"})
		require.NoError(t, err)

		c := a.mergeFilters(b)
		require.Equal(t, "a == b", c.expression)
		require.Equal(t, "eq(a, b)", c.expParsed.String())
	})

	t.Run("both missing", func(t *testing.T) {
		a, err := toInternalFilter(internalFilter{expression: ""})
		require.NoError(t, err)
		b, err := toInternalFilter(internalFilter{expression: ""})
		require.NoError(t, err)

		c := a.mergeFilters(b)
		require.Equal(t, "", c.expression)
		require.Equal(t, "<nil>", c.expParsed.String())
	})

	t.Run("combine", func(t *testing.T) {
		a, err := toInternalFilter(internalFilter{expression: "a == b"})
		require.NoError(t, err)
		b, err := toInternalFilter(internalFilter{expression: "c == d"})
		require.NoError(t, err)

		c := a.mergeFilters(b)
		require.Equal(t, "(a == b) && (c == d)", c.expression)
		require.Equal(t, "and(group(eq(a, b)), group(eq(c, d)))", c.expParsed.String())
	})
}

func TestInternalFilterMergeFilters_constraints(t *testing.T) {
	t.Run("b no constr", func(t *testing.T) {
		a := internalFilter{constraints: map[string][]any{"a": {"b"}}}
		b := internalFilter{constraints: nil}

		c := a.mergeFilters(b)
		require.Equal(t, map[string][]any{"a": {"b"}}, c.constraints)
	})

	t.Run("a no constr", func(t *testing.T) {
		a := internalFilter{constraints: nil}
		b := internalFilter{constraints: map[string][]any{"a": {"b"}}}

		c := a.mergeFilters(b)
		require.Equal(t, map[string][]any{"a": {"b"}}, c.constraints)
	})

	t.Run("both missing", func(t *testing.T) {
		a := internalFilter{constraints: nil}
		b := internalFilter{constraints: nil}

		c := a.mergeFilters(b)
		// Use this assertion because equal asserts type also
		require.Nil(t, c.constraints)
	})

	t.Run("combine", func(t *testing.T) {
		a := internalFilter{constraints: map[string][]any{"a": {"b"}}}
		b := internalFilter{constraints: map[string][]any{"a": {"c"}, "d": {"e"}}}

		c := a.mergeFilters(b)
		require.Equal(t, map[string][]any{"a": {"b", "c"}, "d": {"e"}}, c.constraints)
	})
}

func TestInternalFilterMergeFilters_stateConstraints(t *testing.T) {
	t.Run("b no constr", func(t *testing.T) {
		a := internalFilter{stateConstraints: map[string]filter.State{"a": filter.StateInclusive}}
		b := internalFilter{stateConstraints: nil}

		c := a.mergeFilters(b)
		require.Equal(t, map[string]filter.State{"a": filter.StateInclusive}, c.stateConstraints)
	})

	t.Run("a no constr", func(t *testing.T) {
		a := internalFilter{stateConstraints: nil}
		b := internalFilter{stateConstraints: map[string]filter.State{"a": filter.StateInclusive}}

		c := a.mergeFilters(b)
		require.Equal(t, map[string]filter.State{"a": filter.StateInclusive}, c.stateConstraints)
	})

	t.Run("both missing", func(t *testing.T) {
		a := internalFilter{stateConstraints: nil}
		b := internalFilter{stateConstraints: nil}

		c := a.mergeFilters(b)
		// Use this assertion because equal asserts type also
		require.Nil(t, c.stateConstraints)
	})

	t.Run("combine", func(t *testing.T) {
		a := internalFilter{stateConstraints: map[string]filter.State{"a": filter.StateInclusive}}
		b := internalFilter{stateConstraints: map[string]filter.State{"a": filter.StateExcluded, "d": filter.StateExclusive}}

		c := a.mergeFilters(b)
		require.Equal(t, map[string]filter.State{"a": filter.StateExcluded, "d": filter.StateExclusive}, c.stateConstraints)
	})
}

func TestInternalFilterMergeFilters_order(t *testing.T) {
	t.Run("b no order", func(t *testing.T) {
		a := internalFilter{orderBy: filter.SortExprSet{{Column: "a"}}}
		b := internalFilter{orderBy: nil}

		c := a.mergeFilters(b)
		require.Equal(t, filter.SortExprSet{{Column: "a"}}, c.orderBy)
	})

	t.Run("a no order", func(t *testing.T) {
		a := internalFilter{orderBy: nil}
		b := internalFilter{orderBy: filter.SortExprSet{{Column: "a"}}}

		c := a.mergeFilters(b)
		require.Equal(t, filter.SortExprSet{{Column: "a"}}, c.orderBy)
	})

	t.Run("both missing", func(t *testing.T) {
		a := internalFilter{orderBy: nil}
		b := internalFilter{orderBy: nil}

		c := a.mergeFilters(b)
		// Use this assertion because equal asserts type also
		require.Nil(t, c.orderBy)
	})

	t.Run("combine", func(t *testing.T) {
		a := internalFilter{orderBy: filter.SortExprSet{{Column: "a"}}}
		b := internalFilter{orderBy: filter.SortExprSet{{Column: "b"}}}

		c := a.mergeFilters(b)
		require.Equal(t, filter.SortExprSet{{Column: "a"}, {Column: "b"}}, c.orderBy)
	})
}

func TestInternalFilterMergeFilters_cursor(t *testing.T) {
	t.Run("b no cursor", func(t *testing.T) {
		a := internalFilter{cursor: &filter.PagingCursor{ROrder: true}}
		b := internalFilter{orderBy: nil}

		c := a.mergeFilters(b)
		require.Nil(t, c.cursor)
	})

	t.Run("a no constr", func(t *testing.T) {
		a := internalFilter{orderBy: nil}
		b := internalFilter{cursor: &filter.PagingCursor{ROrder: true}}

		c := a.mergeFilters(b)
		require.Equal(t, &filter.PagingCursor{ROrder: true}, c.cursor)
	})

	t.Run("both missing", func(t *testing.T) {
		a := internalFilter{orderBy: nil}
		b := internalFilter{orderBy: nil}

		c := a.mergeFilters(b)
		// Use this assertion because equal asserts type also
		require.Nil(t, c.cursor)
	})

	t.Run("combine", func(t *testing.T) {
		a := internalFilter{cursor: &filter.PagingCursor{ROrder: true}}
		b := internalFilter{cursor: &filter.PagingCursor{ROrder: false}}

		c := a.mergeFilters(b)
		require.Equal(t, &filter.PagingCursor{ROrder: false}, c.cursor)
	})
}

func TestInternalFilterEmpty(t *testing.T) {
	tcc := []struct {
		name string
		f    internalFilter
		want bool
	}{
		{
			name: "completely empty",
			f:    internalFilter{},
			want: true,
		},
		{
			name: "only constraints",
			f: internalFilter{
				constraints: map[string][]any{
					"a": {"b"},
				},
			},
			want: false,
		},
		{
			name: "only state constraints",
			f: internalFilter{
				stateConstraints: map[string]filter.State{
					"a": filter.StateInclusive,
				},
			},
			want: false,
		},
		{
			name: "only expression",
			f: internalFilter{
				expression: "true",
			},
			want: false,
		},
		{
			name: "only order",
			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "a"}},
			},
			want: false,
		},
		{
			name: "only cursor",
			f: internalFilter{
				cursor: &filter.PagingCursor{ROrder: true},
			},
			want: false,
		},
		{
			name: "all",
			f: internalFilter{
				constraints: map[string][]any{
					"a": {"b"},
				},
				stateConstraints: map[string]filter.State{
					"a": filter.StateInclusive,
				},
				orderBy: filter.SortExprSet{{Column: "a"}},
				cursor:  &filter.PagingCursor{ROrder: true},
			},
			want: false,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.f.empty()
			require.Equal(t, tc.want, got)
		})
	}
}
