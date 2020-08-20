package rdbms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
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
	q = q.
		Column(fmt.Sprintf("UNIX_TIMESTAMP(DATE(%s)) timestamp", field)).
		Column("COUNT(*) AS value").
		Where(fmt.Sprintf("%s IS NOT NULL", field)).
		OrderBy("timestamp").
		GroupBy("timestamp")

	var (
		rval      = make([]uint, 0, 100)
		ts, v, i  uint
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	return rval, func() (err error) {
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&ts, &v); err != nil {
				return err
			}

			rval[2*i], rval[2*i+1] = ts, v
			i++
		}

		return
	}()
}
