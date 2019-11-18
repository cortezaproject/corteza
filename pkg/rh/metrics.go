package rh

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"
)

// DailyMetrics aids repositories on simple stat building queries
//
// Returns a slice of numbers (timestamp + value pairs
func DailyMetrics(db *factory.DB, q squirrel.SelectBuilder, field string) (rval []uint, err error) {
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
		GroupBy("timestamp")

	if err = FetchAll(db, q, &aux); err != nil {
		return
	}

	rval = make([]uint, 2*len(aux))
	for i := 0; i < len(aux); i++ {
		rval[2*i], rval[2*i+1] = aux[i].Timestamp, aux[i].Value
	}

	return
}

// MultiDailyMetrics simplifies fetching of multiple daily metrics
func MultiDailyMetrics(db *factory.DB, q squirrel.SelectBuilder, fields []string, mm ...*[]uint) (err error) {
	for m := 0; m < len(mm); m++ {
		*mm[m], err = DailyMetrics(db, q, fields[m])
		if err != nil {
			return
		}
	}

	return
}
