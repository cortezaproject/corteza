package rdbms

import (
	"fmt"
	"testing"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"
)

func Test_buildCursorCond(t *testing.T) {
	tests := []struct {
		cursor *filter.PagingCursor
		sql    string
		esql   string
		args   []any
	}{
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 1, false)
				return c
			}(),
			"((f1 IS NOT NULL AND 1=0) OR (f1 > ?))",
			`((("f1" IS NOT NULL) AND 1=0) OR ("f1" > ?))`,
			[]any{1},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 2, false)
				c.Set("f2", 3, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND 1=0) OR (f1 > ?)) OR (((f1 IS NULL AND 1=0) OR f1 = ?) AND ((f2 IS NOT NULL AND 1=0) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND 1=0) OR ("f1" > ?)) OR (((("f1" IS NULL) AND 1=0) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND 1=0) OR ("f2" > ?))))`,
			[]any{2, 2, 3},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 4, false)
				c.LThen = true
				return c
			}(),
			"((f1 IS NOT NULL AND 1=0) OR (f1 > ?))",
			`((("f1" IS NOT NULL) AND 1=0) OR ("f1" > ?))`,
			[]any{4},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 5, false)
				c.Set("f2", 6, false)
				c.LThen = true
				return c
			}(),
			"(((f1 IS NOT NULL AND 1=0) OR (f1 > ?)) OR (((f1 IS NULL AND 1=0) OR f1 = ?) AND ((f2 IS NOT NULL AND 1=0) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND 1=0) OR ("f1" > ?)) OR (((("f1" IS NULL) AND 1=0) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND 1=0) OR ("f2" > ?))))`,
			[]any{5, 5, 6},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 7, false)
				c.Set("f2", nil, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND 1=0) OR (f1 > ?)) OR (((f1 IS NULL AND 1=0) OR f1 = ?) AND ((f2 IS NOT NULL AND 1=1) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND 1=0) OR ("f1" > ?)) OR (((("f1" IS NULL) AND 1=0) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND 1=1) OR ("f2" > ?))))`,
			[]any{7, 7, nil},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", nil, false)
				c.Set("f2", 8, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND 1=1) OR (f1 > ?)) OR (((f1 IS NULL AND 1=1) OR f1 = ?) AND ((f2 IS NOT NULL AND 1=0) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND 1=1) OR ("f1" > ?)) OR (((("f1" IS NULL) AND 1=1) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND 1=0) OR ("f2" > ?))))`,
			[]any{nil, nil, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.cursor.String(), func(t *testing.T) {
			var (
				req = require.New(t)

				sql, args, err = CursorCondition(tt.cursor, nil, nil).ToSQL()
			)

			req.NoError(err)
			req.Equal(tt.sql, sql)
			req.Equal(tt.args, args)

			ee, err := CursorExpression(tt.cursor, nil, nil)
			req.NoError(err)
			sql, args, err = goqu.Dialect("sqlite3").Select().Where(ee).ToSQL()
			req.NoError(err)
			req.Equal(tt.esql, sql[15:])
		})
	}
}

func Test_buildCursorCondIsNull(t *testing.T) {
	var (
		req = require.New(t)
		cur = &filter.PagingCursor{}
	)

	cur.SetModifier("deletedAt", 1, false, filter.ISNULL, "deletedAt")
	cur.Set("name", "b", false)

	sql, args, err := CursorCondition(cur, nil, map[string]string{"deletedat": "deleted_at", "name": "name"}).ToSQL()
	req.NoError(err)
	req.Equal("((((CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END) IS NOT NULL AND 1=0) OR ((CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END) > ?)) OR ((((CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END) IS NULL AND 1=0) OR (CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END) = ?) AND ((name IS NOT NULL AND 1=0) OR (name > ?))))", sql)
	req.Equal([]any{1, 1, "b"}, args)
}

// The expression builder serves resources sorted on a JSON value; it has to
// honour the same modifiers the raw-SQL builder does.
func Test_cursorExpressionModifiers(t *testing.T) {
	var (
		req = require.New(t)
		cur = &filter.PagingCursor{}
	)

	cur.SetModifier("deletedAt", 0, false, filter.ISNULL, "deletedAt")
	cur.SetModifier("", "2026-01-01", true, filter.COALESCE, "deletedAt", "updatedAt")

	ee, err := CursorExpression(cur, nil, nil)
	req.NoError(err)
	sql, args, err := goqu.Dialect("sqlite3").Select().Where(ee).ToSQL()
	req.NoError(err)
	req.Contains(sql, `((CASE WHEN "deletedAt" IS NULL THEN 0 ELSE 1 END) > ?)`)
	req.Contains(sql, `(COALESCE("deletedAt", "updatedAt") < ?)`)
	req.Equal("[0 0 2026-01-01]", fmt.Sprint(args))
}
