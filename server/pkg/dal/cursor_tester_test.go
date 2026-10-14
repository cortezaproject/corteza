package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/stretchr/testify/require"
)

// The in-Go cursor test keeps the same rows the SQL cursor condition keeps:
// NULLs sort first ascending and last descending, a NULL on either side never
// errors, and a string value with quotes in it is compared as is
func TestGenericRowTester_cursor(t *testing.T) {
	ctx := context.Background()

	row := func(id int, a any) *Row {
		r := (&Row{}).WithValue("id", 0, id)
		if a != nil {
			r = r.WithValue("a", 0, a)
		} else {
			r.SetValue("a", 0, nil)
		}
		return r
	}

	rows := []*Row{row(1, nil), row(2, nil), row(3, 1), row(4, 1), row(5, 2)}

	cursor := func(a any, id int, desc bool) *filter.PagingCursor {
		pc := &filter.PagingCursor{}
		pc.Set("a", a, desc)
		pc.Set("id", id, desc)
		return pc
	}

	tcc := []struct {
		name string
		cur  *filter.PagingCursor
		keep []int
	}{
		{"asc after a=1,id=3", cursor(1, 3, false), []int{4, 5}},
		{"asc after a=nil,id=1", cursor(nil, 1, false), []int{2, 3, 4, 5}},
		{"desc after a=1,id=4", cursor(1, 4, true), []int{3, 1, 2}},
		{"desc after a=nil,id=2", cursor(nil, 2, true), []int{1}},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			f := internalFilter{
				orderBy: filter.SortExprSet{{Column: "a", Descending: tc.cur.Desc()[0]}, {Column: "id", Descending: tc.cur.Desc()[0]}},
				cursor:  tc.cur,
			}

			tester, err := prepareGenericRowTester(f)
			req.NoError(err)

			var kept []int
			for _, r := range rows {
				ok, err := tester.Test(ctx, r)
				req.NoError(err)
				if ok {
					id, _ := r.GetValue("id", 0)
					kept = append(kept, id.(int))
				}
			}

			req.ElementsMatch(tc.keep, kept)
		})
	}

	t.Run("string value with quotes", func(t *testing.T) {
		req := require.New(t)

		pc := &filter.PagingCursor{}
		pc.Set("s", `a "quoted" \ value`, false)
		f := internalFilter{orderBy: filter.SortExprSet{{Column: "s"}}, cursor: pc}

		tester, err := prepareGenericRowTester(f)
		req.NoError(err)

		ok, err := tester.Test(ctx, (&Row{}).WithValue("s", 0, "b"))
		req.NoError(err)
		req.True(ok)
	})
}
