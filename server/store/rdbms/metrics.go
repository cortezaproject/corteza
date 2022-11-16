package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"time"
)

// multiDailyMetrics simplifies fetching of multiple daily metrics
//
// This function is copied from old repository and adapted to work under store,
// might be good to refactor it to fit into this store pattern a bit more
func (s Store) multiDailyMetrics(ctx context.Context, q squirrel.SelectBuilder, fields []string, mm ...*[]uint) (err error) {
	for m := 0; m < len(mm); m++ {
		*mm[m], err = s.dailyMetrics(ctx, q, fields[m])
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
func (s Store) dailyMetrics(ctx context.Context, q squirrel.SelectBuilder, field string) ([]uint, error) {
	sqlfn, err := s.SqlFunctionHandler(ql.Function{Name: "DATE", Arguments: ql.ASTSet{ql.Ident{Value: field}}})
	if err != nil {
		return nil, err
	}

	q = q.
		Column(sqlfn).
		Column("COUNT(*) AS value").
		Where(fmt.Sprintf("%s IS NOT NULL", field))

	// since orderBy & groupBy do not accept Sqlizer, we'll convert it up front
	// args are omitted because there should be none
	// error is checked when adding sqlfn to Column fn
	sqlfnStr, _, _ := sqlfn.ToSql()

	q = q.
		OrderBy(sqlfnStr).
		GroupBy(sqlfnStr)

	const (
		cap = 100
	)

	var (
		rval = make([]uint, 0, cap)

		d  time.Time
		ts string
		c  uint

		rows *sql.Rows
	)

	return rval, func() (err error) {
		rows, err = s.Query(ctx, q)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&ts, &c); err != nil {
				return err
			}

			d = time.Time{}
			if ts != "" {
				if d, err = time.Parse(time.RFC3339, ts); err != nil {
					return err
				}
			}

			rval = append(rval, uint(d.Unix()), c)
		}

		return
	}()
}
