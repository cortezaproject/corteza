package rdbms

import (
	"context"
	"database/sql"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	sqlx "github.com/jmoiron/sqlx/types"
)

// A temporary alias for old type that was used handle JSON encoded data
type rawJson = sqlx.JSONText

// multiDailyMetrics simplifies fetching of multiple daily metrics
//
// This function is copied from old repository and adapted to work under store,
// might be good to refactor it to fit into this store pattern a bit more
func (s Store) multiDailyMetrics(ctx context.Context, tbl exp.IdentifierExpression, expr []goqu.Expression, fields []string, mm ...*[]uint) (err error) {
	for m := 0; m < len(mm); m++ {
		*mm[m], err = s.dailyMetrics(ctx, tbl, expr, fields[m])
		if err != nil {
			return
		}
	}

	return
}

// dailyMetrics aids repositories on simple stat building queries
//
// Returns a slice of numbers (timestamp + value pairs)
//
// This function is copied from old repository and adapted to work under store,
// might be good to refactor it to fit into this store pattern a bit more
func (s Store) dailyMetrics(ctx context.Context, tbl exp.IdentifierExpression, expr []goqu.Expression, field string) ([]uint, error) {
	const format = "2006-01-02"

	var (
		rval = make([]uint, 0)

		d  time.Time
		ts string
		c  uint

		rows *sql.Rows
	)

	return rval, func() (err error) {
		daily, err := s.Dialect.AttributeCast(&dal.Attribute{Type: dal.TypeDate{}}, goqu.C(field))
		if err != nil {
			return
		}

		query := s.Dialect.GOQU().
			Select(daily, goqu.COUNT(goqu.Star()).As("value")).
			From(tbl).
			Where(exp.NewLiteralExpression("? IS NOT NULL", goqu.C(field)), goqu.And(expr...)).
			GroupBy(daily).
			Order(exp.NewOrderedExpression(daily, exp.AscDir, exp.NoNullsSortType))

		if rows, err = s.Query(ctx, query); err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&ts, &c); err != nil {
				return err
			}

			d = time.Time{}
			if len(ts) >= len(format) {
				if d, err = time.Parse(format, ts[0:len(format)]); err != nil {
					return err
				}
			}

			rval = append(rval, uint(d.Unix()), c)
		}

		return
	}()
}

// Assembles select expression that calculates total + each given ts field + valid
//
func timestampStatExpr(fields ...string) []interface{} {
	var (
		valid = []goqu.Expression{}
		ee    = []interface{}{goqu.COUNT(goqu.Star()).As("total")}

		// literal 0 and 1 values we can safely use in the query
		lit0, lit1 = goqu.L("0"), goqu.L("1")

		sum = func(field string) goqu.Expression {
			return goqu.COALESCE(
				goqu.SUM(goqu.Case().When(exp.NewLiteralExpression("? IS NOT NULL", goqu.C(field+"_at")), lit1).Else(lit0)),
				lit0,
			).As(field)
		}
	)

	for _, field := range fields {
		ee = append(ee, sum(field))
		valid = append(valid, exp.NewLiteralExpression("? IS NULL", goqu.C(field+"_at")))
	}

	return append(ee, goqu.SUM(goqu.Case().When(goqu.And(valid...), lit1).Else(lit0)).As("valid"))
}
