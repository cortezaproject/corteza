package rdbms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertUserFilter(f types.UserFilter) (query squirrel.SelectBuilder, err error) {
	if f.Sort == "" {
		f.Sort = "id"
	}

	query = s.QueryUsers()

	// Returns user filter (flt) wrapped in IF() function with cnd as condition (when cnd != nil)
	whereMasked := func(cnd *permissions.ResourceFilter, flt squirrel.Sqlizer) squirrel.Sqlizer {
		if cnd != nil {
			return rh.SquirrelFunction("IF", cnd, flt, squirrel.Expr("false"))
		} else {
			return flt
		}
	}

	query = rh.FilterNullByState(query, "usr.deleted_at", f.Deleted)
	query = rh.FilterNullByState(query, "usr.suspended_at", f.Suspended)

	if len(f.UserID) > 0 {
		query = query.Where(squirrel.Eq{"usr.ID": f.UserID})
	}

	if len(f.RoleID) > 0 {
		or := squirrel.Or{}
		// Due to lack of support for more exotic expressions (slice of values inside subquery)
		// we'll use set of OR expressions as a workaround
		for _, roleID := range f.RoleID {
			or = append(or, squirrel.Expr("usr.ID IN (SELECT rel_user FROM sys_role_member WHERE rel_role = ?)", roleID))
		}

		query = query.Where(or)
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"usr.username": qs},
			squirrel.Like{"usr.handle": qs},
			whereMasked(f.IsEmailUnmaskable, squirrel.Like{"usr.email": qs}),
			whereMasked(f.IsNameUnmaskable, squirrel.Like{"usr.name": qs}),
		})
	}

	if f.Email != "" {
		query = query.Where(whereMasked(f.IsEmailUnmaskable, squirrel.Eq{"usr.email": f.Email}))
	}

	if f.Username != "" {
		query = query.Where(squirrel.Eq{"usr.username": f.Username})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"usr.handle": f.Handle})
	}

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"usr.kind": f.Kind})
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, s.UserColumns()...); err != nil {
		return
	}

	query = query.OrderBy(orderBy...)
	return
}

func (s Store) CountUsers(ctx context.Context, f types.UserFilter) (uint, error) {
	if q, err := s.convertUserFilter(f); err != nil {
		return 0, fmt.Errorf("could not count users: %w", err)
	} else {
		return Count(ctx, s.db, q)
	}
}
