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

		MergeByID(id, targetRoleID uint64) error
		MoveByID(id, targetOrganisationID uint64) error

		MembershipsFindByUserID(userID uint64) ([]*types.RoleMember, error)
		MemberFindByRoleID(roleID uint64) ([]*types.RoleMember, error)
		MemberAddByID(roleID, userID uint64) error
		MemberRemoveByID(roleID, userID uint64) error
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
		From(r.table() + " AS r").
		Where(squirrel.And{
			squirrel.Eq{"deleted_at": nil},
			squirrel.Eq{"archived_at": nil},
		})

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

	if !f.IncDeleted {
		query = query.Where(squirrel.Eq{"r.deleted_at": nil})
	}

	if !f.IncArchived {
		query = query.Where(squirrel.Eq{"r.archived_at": nil})
	}

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

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
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
	return r.updateColumnByID(r.table(), "archived_at", time.Now(), id)
}

func (r *role) UnarchiveByID(id uint64) error {
	return r.updateColumnByID(r.table(), "archived_at", nil, id)
}

func (r *role) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.table(), "deleted_at", time.Now(), id)
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
