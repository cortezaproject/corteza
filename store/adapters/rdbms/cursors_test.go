package rdbms

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
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
			"((f1 IS NOT NULL AND FALSE) OR (f1 > ?))",
			`((("f1" IS NOT NULL) AND FALSE) OR ("f1" > ?))`,
			[]any{1},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 2, false)
				c.Set("f2", 3, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND FALSE) OR (f1 > ?)) OR (((f1 IS NULL AND FALSE) OR f1 = ?) AND ((f2 IS NOT NULL AND FALSE) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND FALSE) OR ("f1" > ?)) OR (((("f1" IS NULL) AND FALSE) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND FALSE) OR ("f2" > ?))))`,
			[]any{2, 2, 3},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 4, false)
				c.LThen = true
				return c
			}(),
			"((f1 IS NOT NULL AND FALSE) OR (f1 > ?))",
			`((("f1" IS NOT NULL) AND FALSE) OR ("f1" > ?))`,
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
			"(((f1 IS NOT NULL AND FALSE) OR (f1 > ?)) OR (((f1 IS NULL AND FALSE) OR f1 = ?) AND ((f2 IS NOT NULL AND FALSE) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND FALSE) OR ("f1" > ?)) OR (((("f1" IS NULL) AND FALSE) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND FALSE) OR ("f2" > ?))))`,
			[]any{5, 5, 6},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", 7, false)
				c.Set("f2", nil, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND FALSE) OR (f1 > ?)) OR (((f1 IS NULL AND FALSE) OR f1 = ?) AND ((f2 IS NOT NULL AND TRUE) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND FALSE) OR ("f1" > ?)) OR (((("f1" IS NULL) AND FALSE) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND TRUE) OR ("f2" > ?))))`,
			[]any{7, 7, nil},
		},
		{
			func() *filter.PagingCursor {
				c := &filter.PagingCursor{}
				c.Set("f1", nil, false)
				c.Set("f2", 8, false)
				return c
			}(),
			"(((f1 IS NOT NULL AND TRUE) OR (f1 > ?)) OR (((f1 IS NULL AND TRUE) OR f1 = ?) AND ((f2 IS NOT NULL AND FALSE) OR (f2 > ?))))",
			`(((("f1" IS NOT NULL) AND TRUE) OR ("f1" > ?)) OR (((("f1" IS NULL) AND TRUE) OR ("f1" = ?)) AND ((("f2" IS NOT NULL) AND FALSE) OR ("f2" > ?))))`,
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
