package rdbms

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertApplicationFilter(f types.ApplicationFilter) (query squirrel.SelectBuilder, err error) {
	query = s.applicationsSelectBuilder()

	query = filter.StateCondition(query, "app.deleted_at", f.Deleted)

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"app.id": f.LabeledIDs})
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"app.name": qs},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"app.name": f.Name})
	}

	return
}

func (s Store) ApplicationMetrics(ctx context.Context) (*types.ApplicationMetrics, error) {
	var (
		counters = squirrel.
				Select(
				"COUNT(*) as total",
				"SUM(CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END) as deleted",
				"SUM(CASE WHEN deleted_at IS NULL THEN 1 ELSE 0 END) as valid",
			).
			PlaceholderFormat(s.config.PlaceholderFormat).
			From(s.applicationTable("u"))

		rval     = &types.ApplicationMetrics{}
		row, err = s.QueryRow(ctx, counters)
	)

	if err != nil {
		return nil, err
	}

	err = row.Scan(&rval.Total, &rval.Deleted, &rval.Valid)
	if err != nil {
		return nil, err
	}

	return rval, nil
}

func (s Store) ReorderApplications(ctx context.Context, order []uint64) (err error) {
	var (
		apps   types.ApplicationSet
		appMap = map[uint64]bool{}
		weight = 1

		f = types.ApplicationFilter{}
	)

	if apps, _, err = s.SearchApplications(ctx, f); err != nil {
		return
	}

	for _, app := range apps {
		appMap[app.ID] = true
	}

	// honor parameter first
	for _, pageID := range order {
		if appMap[pageID] {
			appMap[pageID] = false
			err = s.execUpdateApplications(ctx,
				squirrel.Eq{"app.id": pageID},
				store.Payload{"weight": weight})

			if err != nil {
				return
			}

			weight++
		}
	}

	for pageID, update := range appMap {
		if !update {
			continue
		}

		err = s.execUpdateApplications(ctx,
			squirrel.Eq{"app.id": pageID},
			store.Payload{"weight": weight})

		if err != nil {
			return
		}

		weight++

	}

	return
}
