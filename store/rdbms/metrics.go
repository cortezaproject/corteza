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
func (s Store) dailyMetrics(ctx context.Context, q squirrel.SelectBuilder, field string) (rval []uint, err error) {
	var (
		aux = make([]struct {
			Timestamp uint
			Value     uint
		}, 0)
	)

	q = q.
		Column(fmt.Sprintf("UNIX_TIMESTAMP(DATE(%s)) timestamp", field)).
		Column("COUNT(*) AS value").
		Where(fmt.Sprintf("%s IS NOT NULL", field)).
		OrderBy("timestamp").
		GroupBy("timestamp")

	sql, args := q.MustSql()

	if err = s.db.SelectContext(ctx, &aux, sql, args...); err != nil {
		return
	}

	rval = make([]uint, 2*len(aux))
	for i := 0; i < len(aux); i++ {
		rval[2*i], rval[2*i+1] = aux[i].Timestamp, aux[i].Value
	}

	return
}
