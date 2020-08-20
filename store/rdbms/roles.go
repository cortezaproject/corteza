package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertRoleFilter(f types.RoleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.rolesSelectBuilder()

	query = rh.FilterNullByState(query, "rl.deleted_at", f.Deleted)
	query = rh.FilterNullByState(query, "rl.archived_at", f.Archived)

	if len(f.RoleID) > 0 {
		query = query.Where(squirrel.Eq{"rl.ID": f.RoleID})
	}

	if f.MemberID > 0 {
		query = query.Where(squirrel.Expr("rl.ID IN (SELECT rel_role FROM sys_role_member AS m WHERE m.rel_user = ?)", f.MemberID))
	}

	if f.Query != "" {
		qs := f.Query + "%"
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
				"SUM(IF(deleted_at IS NULL, 0, 1)) as deleted",
				"SUM(IF(archived_at IS NULL, 0, 1)) as archived",
				"SUM(IF(deleted_at IS NULL AND archived_at IS NULL, 1, 0)) as valid",
			).
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
		squirrel.Select().From(s.roleTable("u")),
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"suspended_at",
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
