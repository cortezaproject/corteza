package ql

import (
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	const WHERE = " WHERE "
	var (
		base = goqu.Dialect("sqlite").Select().Prepared(true)
		conv = Converter()

		cases = []struct {
			qry  string
			sql  string
			args []any
		}{
			{
				qry:  `foo = 1`,
				sql:  `("foo" = ?)`,
				args: []any{int64(1)},
			},
			{
				qry:  `foo = 1 AND bar = baz`,
				sql:  `(("foo" = ?) AND ("bar" = "baz"))`,
				args: []any{int64(1)},
			},
			{
				qry:  `foo = 1 OR (bar = 'baz')`,
				sql:  `(("foo" = ?) OR ("bar" = ?))`,
				args: []any{int64(1), "baz"},
			},
			{
				qry:  `foo is not null`,
				sql:  `("foo" IS NOT NULL)`,
				args: []any{},
			},
			{
				qry:  `!foo`,
				sql:  `(NOT "foo")`,
				args: []any{},
			},
			{
				qry:  `one + 2 / 3 * 4 < 10`,
				sql:  `("one" + ? / ? * ? < ?)`,
				args: []any{int64(2), int64(3), int64(4), int64(10)},
			},
			{
				qry:  `concat('foo', 'bar')`,
				sql:  `CONCAT(?, ?)`,
				args: []any{"foo", "bar"},
			},
			{
				qry:  `quarter('2022-07-21')`,
				sql:  `QUARTER(?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `year('2022-07-21')`,
				sql:  `YEAR(?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `month('2022-07-21')`,
				sql:  `MONTH(?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `date('2022-07-21')`,
				sql:  `DAY(?)`,
				args: []any{"2022-07-21"},
			},
		}
	)

	for _, c := range cases {
		t.Run(c.qry, func(t *testing.T) {
			req := require.New(t)

			t.Log(c.qry)
			ee, err := conv.Parse(c.qry)
			req.NoError(err)

			sql, args, err := base.Where(ee).ToSQL()
			req.NoError(err)

			p := strings.Index(sql, WHERE)
			req.Positive(p)

			sql = sql[p+len(WHERE):]

			req.Equal(c.sql, sql)
			req.Equal(c.args, args)
		})
	}
}
