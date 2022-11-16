package rdbms

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (s Store) convertUserFilter(f types.UserFilter) (query squirrel.SelectBuilder, err error) {
	query = s.usersSelectBuilder()

	query = filter.StateCondition(query, "usr.deleted_at", f.Deleted)
	query = filter.StateCondition(query, "usr.suspended_at", f.Suspended)

	if len(f.UserID) > 0 {
		query = query.Where(squirrel.Eq{"usr.id": f.UserID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"usr.id": f.LabeledIDs})
	}

	if len(f.RoleID) > 0 {
		or := squirrel.Or{}
		// Due to lack of support for more exotic expressions (slice of values inside subquery)
		// we'll use set of OR expressions as a workaround
		for _, roleID := range f.RoleID {
			or = append(or, squirrel.Expr("usr.ID IN (SELECT rel_user FROM role_members WHERE rel_role = ?)", roleID))
		}

		query = query.Where(or)
	}

	if f.Query != "" {
		qs := "%" + f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"usr.username": qs},
			squirrel.Like{"usr.handle": qs},

			// do a lookup on a potentially masked fields
			// this should be filtered out in the check function
			squirrel.Like{"usr.email": qs},
			squirrel.Like{"usr.name": qs},
		})
	}

	if f.Email != "" {
		query = query.Where(squirrel.Eq{"usr.email": f.Email})
	}

	if f.Username != "" {
		query = query.Where(squirrel.Eq{"usr.username": f.Username})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"usr.handle": f.Handle})
	}

	if !f.AllKinds {
		// When not explicitly requested to search all kids of users,
		// always limit to one kind
		query = query.Where(squirrel.Eq{"usr.kind": f.Kind})
	}

	return
}

func (s Store) CountUsers(ctx context.Context, f types.UserFilter) (uint, error) {
	if q, err := s.convertUserFilter(f); err != nil {
		return 0, fmt.Errorf("could not count users: %w", err)
	} else {
		return Count(ctx, s.db, q)
	}
}

func (s Store) UserMetrics(ctx context.Context) (*types.UserMetrics, error) {
	var (
		counters = squirrel.
				Select(
				"COUNT(*) as total",
				"SUM(CASE WHEN deleted_at   IS NULL                          THEN 0 ELSE 1 END) as deleted",
				"SUM(CASE WHEN suspended_at IS NULL                          THEN 0 ELSE 1 END) as suspended",
				"SUM(CASE WHEN deleted_at   IS NULL AND suspended_at IS NULL THEN 1 ELSE 0 END) as valid",
			).
			PlaceholderFormat(s.config.PlaceholderFormat).
			From(s.userTable("u")).
			Where(squirrel.Eq{"kind": ""})

		row, err = s.QueryRow(ctx, counters)
		rval     = &types.UserMetrics{}
	)

	if err != nil {
		return nil, err
	}

	err = row.Scan(&rval.Total, &rval.Deleted, &rval.Suspended, &rval.Valid)
	if err != nil {
		return nil, err
	}

	// Fetch daily metrics for created, updated, deleted and suspended users
	err = s.multiDailyMetrics(
		ctx,
		squirrel.
			Select().
			PlaceholderFormat(s.config.PlaceholderFormat).
			From(s.userTable("u")).
			Where(squirrel.Eq{"kind": ""}),
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"suspended_at",
		},
		&rval.DailyCreated,
		&rval.DailyUpdated,
		&rval.DailyDeleted,
		&rval.DailySuspended,
	)

	if err != nil {
		return nil, err
	}

	return rval, nil
}
