package rdbms

import (
	"context"
	"sort"

	systemType "github.com/crusttech/human/server/system/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) ApplicationMetrics(ctx context.Context) (_ *systemType.ApplicationMetrics, err error) {
	var (
		aux = struct {
			Total   uint `db:"total"`
			Deleted uint `db:"deleted"`
			Valid   uint `db:"valid"`
		}{}

		query = applicationSelectQuery(s.Dialect.GOQU()).
			Select(timestampStatExpr("deleted")...)
	)

	if err = s.QueryOne(ctx, query, &aux); err != nil {
		return nil, err
	}

	return &systemType.ApplicationMetrics{
		Total:   aux.Total,
		Deleted: aux.Deleted,
		Valid:   aux.Valid,
	}, nil
}

func (s Store) ReorderApplications(ctx context.Context, order []uint64) (err error) {
	var (
		apps   systemType.ApplicationSet
		appMap = map[uint64]bool{}
		weight = 1

		f = systemType.ApplicationFilter{}

		query = func(id uint64, weight int) *goqu.UpdateDataset {
			return s.Dialect.GOQU().
				Update(applicationTable).
				Set(goqu.Record{"weight": weight}).
				Where(goqu.C("id").Eq(id))
		}
	)

	if apps, _, err = s.SearchApplications(ctx, f); err != nil {
		return
	}

	for _, app := range apps {
		appMap[app.ID] = true
	}

	// honor parameter first
	for _, id := range order {
		if appMap[id] {
			appMap[id] = false

			if err = s.Exec(ctx, query(id, weight)); err != nil {
				return
			}

			weight++
		}
	}

	// Applications the caller did not name keep their relative order. Ranging
	// over appMap directly would not: Go randomises map iteration, so
	// reordering a subset silently shuffled every other application,
	// differently on each call.
	leftovers := make(systemType.ApplicationSet, 0, len(appMap))
	for _, app := range apps {
		if appMap[app.ID] {
			leftovers = append(leftovers, app)
		}
	}
	sort.SliceStable(leftovers, func(i, j int) bool { return leftovers[i].Weight < leftovers[j].Weight })

	for _, app := range leftovers {
		if err = s.Exec(ctx, query(app.ID, weight)); err != nil {
			return
		}

		weight++
	}

	return
}
