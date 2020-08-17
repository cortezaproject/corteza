package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertApplicationFilter(f types.ApplicationFilter) (query squirrel.SelectBuilder, err error) {
	if f.Sort == "" {
		f.Sort = "id"
	}

	query = s.QueryApplications()

	query = rh.FilterNullByState(query, "app.deleted_at", f.Deleted)

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"app.name": qs},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"app.name": f.Name})
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, s.ApplicationColumns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	return
}

func (s Store) ApplicationMetrics(ctx context.Context) (rval *types.ApplicationMetrics, err error) {
	var (
		counters = squirrel.
			Select(
				"COUNT(*) as total",
				"SUM(IF(deleted_at IS NULL, 0, 1)) as deleted",
				"SUM(IF(deleted_at IS NULL, 1, 0)) as valid",
			).
			From(s.UserTable("u"))
	)

	rval = &types.ApplicationMetrics{}

	var (
		sql, args = counters.MustSql()
		row       = s.db.QueryRowContext(ctx, sql, args...)
	)

	err = row.Scan(&rval.Total, &rval.Deleted, &rval.Valid)
	if err != nil {
		return nil, err
	}

	return
}
