package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	RoleRepository interface {
		With(ctx context.Context, db *factory.DB) RoleRepository

		FindByID(id uint64) (*types.Role, error)
		FindByName(name string) (*types.Role, error)
		FindByHandle(handle string) (*types.Role, error)
		Find(filter types.RoleFilter) (types.RoleSet, types.RoleFilter, error)

		Create(mod *types.Role) (*types.Role, error)
		Update(mod *types.Role) (*types.Role, error)

		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error
		UndeleteByID(id uint64) error

		MergeByID(id, targetRoleID uint64) error
		MoveByID(id, targetOrganisationID uint64) error

		MembershipsFindByUserID(userID uint64) ([]*types.RoleMember, error)
		MemberFindByRoleID(roleID uint64) ([]*types.RoleMember, error)
		MemberAddByID(roleID, userID uint64) error
		MemberRemoveByID(roleID, userID uint64) error

		Metrics() (*types.RoleMetrics, error)
	}

	role struct {
		*repository

		// sql table reference
		roles   string
		members string
	}
)

const (
	ErrRoleNotFound = repositoryError("RoleNotFound")
)

// @todo migrate to same pattern as we have for uselang/en.jsonrs
func Role(ctx context.Context, db *factory.DB) RoleRepository {
	return (&role{}).With(ctx, db)
}

func (r *role) With(ctx context.Context, db *factory.DB) RoleRepository {
	return &role{
		repository: r.repository.With(ctx, db),
	}
}

func (r role) table() string {
	return "sys_role"
}

func (r role) tableMember() string {
	return "sys_role_member"
}

func (r role) columns() []string {
	return []string{
		"id",
		"name",
		"handle",
		"created_at",
		"updated_at",
		"archived_at",
		"deleted_at",
	}
}

func (r role) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS r")
}

func (r role) FindByID(id uint64) (*types.Role, error) {
	return r.findOneBy("id", id)
}

func (r role) FindByHandle(handle string) (*types.Role, error) {
	return r.findOneBy("handle", handle)
}

func (r role) FindByName(name string) (*types.Role, error) {
	return r.findOneBy("name", name)
}

func (r role) findOneBy(field string, value interface{}) (*types.Role, error) {
	var (
		ro = &types.Role{}
		q  = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, ro)
	)

	if err != nil {
		return nil, err
	} else if ro.ID == 0 {
		return nil, ErrRoleNotFound
	}

	return ro, nil
}

func (r *role) Find(filter types.RoleFilter) (set types.RoleSet, f types.RoleFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id"
	}

	query := r.query()

	query = rh.FilterNullByState(query, "r.deleted_at", f.Deleted)
	query = rh.FilterNullByState(query, "r.archived_at", f.Archived)

	if len(f.RoleID) > 0 {
		query = query.Where(squirrel.Eq{"r.ID": f.RoleID})
	}

	if f.MemberID > 0 {
		query = query.Where(squirrel.Expr("r.ID IN (SELECT rel_role FROM sys_role_member AS m WHERE m.rel_user = ?)", f.MemberID))
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"r.name": qs},
			squirrel.Like{"r.handle": qs},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"r.name": f.Name})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"r.handle": f.Handle})
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.PageFilter, &set)
}

func (r *role) Create(mod *types.Role) (*types.Role, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r *role) Update(mod *types.Role) (*types.Role, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	return mod, r.db().Replace(r.table(), mod)
}

func (r *role) ArchiveByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"archived_at": time.Now()}, squirrel.Eq{"id": id})
}

func (r *role) UnarchiveByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"archived_at": nil}, squirrel.Eq{"id": id})
}

func (r *role) DeleteByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": time.Now()}, squirrel.Eq{"id": id})
}

func (r *role) UndeleteByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": nil}, squirrel.Eq{"id": id})
}

func (r *role) MergeByID(id, targetRoleID uint64) error {
	return ErrNotImplemented
}

func (r *role) MoveByID(id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}

func (r *role) MembershipsFindByUserID(roleID uint64) (mm []*types.RoleMember, err error) {
	rval := make([]*types.RoleMember, 0)
	sql := "SELECT * FROM " + r.tableMember() + " WHERE rel_user = ?"
	return rval, r.db().Select(&rval, sql, roleID)
}

func (r *role) MemberFindByRoleID(roleID uint64) (mm []*types.RoleMember, err error) {
	rval := make([]*types.RoleMember, 0)
	sql := "SELECT * FROM " + r.tableMember() + " WHERE rel_role = ?"
	return rval, r.db().Select(&rval, sql, roleID)
}

func (r *role) MemberAddByID(roleID, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: roleID,
		UserID: userID,
	}
	return r.db().Replace(r.tableMember(), mod)
}

func (r *role) MemberRemoveByID(roleID, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: roleID,
		UserID: userID,
	}
	return r.db().Delete(r.tableMember(), mod, "rel_role", "rel_user")
}

// Metrics collects and returns user metrics
func (r role) Metrics() (rval *types.RoleMetrics, err error) {
	var (
		counters = squirrel.
			Select(
				"COUNT(*) as total",
				"SUM(IF(deleted_at IS NULL, 0, 1)) as deleted",
				"SUM(IF(archived_at IS NULL, 0, 1)) as archived",
				"SUM(IF(deleted_at IS NULL AND archived_at IS NULL, 1, 0)) as valid",
			).
			From(r.table() + " AS u")
	)

	rval = &types.RoleMetrics{}

	if err = rh.FetchOne(r.db(), counters, rval); err != nil {
		return
	}

	// Fetch daily metrics for created, updated, deleted and archived roles
	err = rh.MultiDailyMetrics(
		r.db(),
		squirrel.Select().From(r.table()+" AS u"),
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
		return
	}

	return
}
