package rdbms

import (
	"context"

	systemType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) CountUsers(ctx context.Context, f systemType.UserFilter) (c uint, _ error) {
	var (
		aux = struct {
			Count uint `db:"count"`
		}{}
		expr, _, err = s.config.Filters.User(f)

		query = s.config.Dialect.
			From(userTable).
			Select(goqu.COUNT(goqu.Star()).As("count"))
	)

	if err != nil {
		return
	}

	if err = s.QueryOne(ctx, query.Where(expr...).Limit(1), &aux); err != nil {
		return
	}

	return aux.Count, nil
}

func (s Store) UserMetrics(ctx context.Context) (m *systemType.UserMetrics, err error) {
	var (
		aux = struct {
			Total     uint `db:"total"`
			Deleted   uint `db:"deleted"`
			Valid     uint `db:"valid"`
			Suspended uint `db:"suspended"`
		}{}

		query = userSelectQuery(s.config.Dialect).
			Select(timestampStatExpr("deleted", "suspended")...)
	)

	if err = s.QueryOne(ctx, query, &aux); err != nil {
		return nil, err
	}

	m = &systemType.UserMetrics{
		Total:     aux.Total,
		Valid:     aux.Valid,
		Deleted:   aux.Deleted,
		Suspended: aux.Suspended,
	}

	// Fetch daily metrics for created, updated, deleted and suspended users
	err = s.multiDailyMetrics(
		ctx,
		userTable,
		[]goqu.Expression{
			goqu.C("kind").Eq(""),
		},
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"suspended_at",
		},
		&m.DailyCreated,
		&m.DailyUpdated,
		&m.DailyDeleted,
		&m.DailySuspended,
	)

	if err != nil {
		return nil, err
	}

	return
}
