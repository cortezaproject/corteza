package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertRoleFilter(f types.RoleFilter) (query squirrel.SelectBuilder, err error) {
	if f.Sort == "" {
		f.Sort = "id"
	}

	query = s.QueryRoles()

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

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, s.RoleColumns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	return
}

func (s Store) roleMemberTable() string {
	return "sys_role_member"
}

//func (s *Store) MembershipsFindByUserID(roleID uint64) (mm []*types.RoleMember, err error) {
//	rval := make([]*types.RoleMember, 0)
//	sql := "SELECT * FROM " + rl.tableMember() + " WHERE rel_user = ?"
//	return rval, rl.db().Select(&rval, sql, roleID)
//}
//
//func (s *Store) MemberFindByRoleID(roleID uint64) (mm []*types.RoleMember, err error) {
//	rval := make([]*types.RoleMember, 0)
//	sql := "SELECT * FROM " + rl.tableMember() + " WHERE rel_role = ?"
//	return rval, rl.db().Select(&rval, sql, roleID)
//}
//
//func (s *Store) MemberAddByID(roleID, userID uint64) error {
//	mod := &types.RoleMember{
//		RoleID: roleID,
//		UserID: userID,
//	}
//	return rl.db().Replace(rl.tableMember(), mod)
//}
//
//func (s *Store) MemberRemoveByID(roleID, userID uint64) error {
//	mod := &types.RoleMember{
//		RoleID: roleID,
//		UserID: userID,
//	}
//	return rl.db().Delete(rl.tableMember(), mod, "rel_role", "rel_user")
//}

//// Metrics collects and returns user metrics
//func (s *Store) Metrics() (rval *types.RoleMetrics, err error) {
//	var (
//		counters = squirrel.
//			Select(
//				"COUNT(*) as total",
//				"SUM(IF(deleted_at IS NULL, 0, 1)) as deleted",
//				"SUM(IF(archived_at IS NULL, 0, 1)) as archived",
//				"SUM(IF(deleted_at IS NULL AND archived_at IS NULL, 1, 0)) as valid",
//			).
//			From(s.roleTable() + " AS u")
//	)
//
//	rval = &types.RoleMetrics{}
//
//	if err = rh.FetchOne(s.db(), counters, rval); err != nil {
//		return
//	}
//
//	// Fetch daily metrics for created, updated, deleted and archived roles
//	err = rh.MultiDailyMetrics(
//		s.db(),
//		squirrel.Select().From(s.roleTable()+" AS u"),
//		[]string{
//			"created_at",
//			"updated_at",
//			"deleted_at",
//			"archived_at",
//		},
//		&rval.DailyCreated,
//		&rval.DailyUpdated,
//		&rval.DailyDeleted,
//		&rval.DailyArchived,
//	)
//
//	if err != nil {
//		return
//	}
//
//	return
//}
