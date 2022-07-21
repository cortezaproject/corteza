package sqlite

import (
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// @todo Ql functions should be under store/tests so it can be tested across all drivers along with generated tests.
// 		for now, Its test coverage is limited per driver.
func TestConverter(t *testing.T) {
	const SELECT = "SELECT "
	var (
		conv = ql.Converter(ql.RefHandler(dialect.ExprHandler))

		cases = []struct {
			qry  string
			sql  string
			args []any
		}{
			{
				qry:  `now()`,
				sql:  `DATE('NOW')`,
				args: []any{},
			},
			{
				qry:  `quarter('2022-07-21')`,
				sql:  `(CAST(STRFTIME('%m', ?) AS INTEGER) + 2) / 3`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `year('2022-07-21')`,
				sql:  `STRFTIME('%Y', ?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `month('2022-07-21')`,
				sql:  `STRFTIME('%m', ?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `date('2022-07-21')`,
				sql:  `STRFTIME('%Y-%m-%dT00:00:00Z', ?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `datetime('2022-07-21')`,
				sql:  `DATETIME(?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `timestamp('2022-07-21')`,
				sql:  `DATETIME(?)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `date_format('2022-07-21','%d')`,
				sql:  `STRFTIME(?, ?)`,
				args: []any{"%d", "2022-07-21"},
			},
		}
	)

	for _, c := range cases {
		t.Run(c.qry, func(t *testing.T) {
			req := require.New(t)

			ee, err := conv.Parse(c.qry)
			req.NoError(err)

			sql, args, err := dialect.GOQU().Select(ee).ToSQL()
			req.NoError(err)

			p := strings.Index(sql, SELECT)
			req.Zero(p)

			sql = sql[p+len(SELECT):]

			req.Equal(c.sql, sql)
			req.Equal(c.args, args)
		})
	}

}
