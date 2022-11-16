package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (s Store) convertRoleFilter(f types.RoleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.rolesSelectBuilder()

	query = filter.StateCondition(query, "rl.deleted_at", f.Deleted)
	query = filter.StateCondition(query, "rl.archived_at", f.Archived)

	if len(f.RoleID) > 0 {
		query = query.Where(squirrel.Eq{"rl.ID": f.RoleID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"rl.id": f.LabeledIDs})
	}

	if f.MemberID > 0 {
		query = query.Where(squirrel.Expr("rl.ID IN (SELECT rel_role FROM role_members AS m WHERE m.rel_user = ?)", f.MemberID))
	}

	if f.Query != "" {
		qs := "%" + f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"rl.name": qs},
			squirrel.Like{"rl.handle": qs},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"rl.name": f.Name})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"rl.handle": f.Handle})
	}

	return
}

func (s Store) RoleMetrics(ctx context.Context) (*types.RoleMetrics, error) {
	var (
		counters = squirrel.
				Select(
				"COUNT(*) as total",
				"SUM(CASE WHEN deleted_at   IS NULL                          THEN 0 ELSE 1 END) as deleted",
				"SUM(CASE WHEN archived_at  IS NULL                          THEN 0 ELSE 1 END) as archived",
				"SUM(CASE WHEN deleted_at   IS NULL AND archived_at  IS NULL THEN 1 ELSE 0 END) as valid",
			).
			PlaceholderFormat(s.config.PlaceholderFormat).
			From(s.roleTable("u"))

		rval     = &types.RoleMetrics{}
		row, err = s.QueryRow(ctx, counters)
	)

	if err != nil {
		return nil, err
	}

	err = row.Scan(&rval.Total, &rval.Deleted, &rval.Archived, &rval.Valid)
	if err != nil {
		return nil, err
	}

	// Fetch daily metrics for created, updated, deleted and suspended users
	err = s.multiDailyMetrics(
		ctx,
		squirrel.
			Select().
			PlaceholderFormat(s.config.PlaceholderFormat).
			From(s.roleTable("u")),
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"archived_at",
		},
		&rval.DailyCreated,
		&rval.DailyUpdated,
		&rval.DailyDeleted,
		&rval.DailyArchived,
	)

	if err != nil {
		return nil, err
	}

	return rval, nil
}
