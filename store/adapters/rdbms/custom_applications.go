package rdbms

import (
	"context"

	systemType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) ApplicationMetrics(ctx context.Context) (_ *systemType.ApplicationMetrics, err error) {
	var (
		aux = struct {
			Total   uint `db:"total"`
			Deleted uint `db:"deleted"`
			Valid   uint `db:"valid"`
		}{}

		query = applicationSelectQuery(s.Dialect).
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
			return s.Dialect.
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

	for id, update := range appMap {
		if !update {
			continue
		}

		if err = s.Exec(ctx, query(id, weight)); err != nil {
			return
		}

		weight++

	}

	return
}
