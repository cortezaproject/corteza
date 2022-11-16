package postgres

import (
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
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
				qry:  `quarter('2022-07-21')`,
				sql:  `EXTRACT(QUARTER FROM TIMESTAMP $1)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `year('2022-07-21')`,
				sql:  `EXTRACT(YEAR FROM TIMESTAMP $1)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `month('2022-07-21')`,
				sql:  `EXTRACT(MONTH FROM TIMESTAMP $1)`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `timestamp('2022-07-21')`,
				sql:  `$1::TIMESTAMPTZ`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `date('2022-07-21')`,
				sql:  `$1::DATE`,
				args: []any{"2022-07-21"},
			},
			{
				qry:  `time('2022-07-21 12:41')`,
				sql:  `DATE_TRUNC('second', $1::TIME)::TIME`,
				args: []any{"2022-07-21 12:41"},
			},
			{
				qry:  `date_format('2022-07-21','%a')`,
				sql:  `TO_CHAR($1::TIMESTAMPTZ, $2::TEXT)`,
				args: []any{"2022-07-21", "Dy"},
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
