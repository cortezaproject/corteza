package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"
)

type (
	RoleRepository interface {
		With(ctx context.Context, db *factory.DB) RoleRepository

		FindByID(id uint64) (*types.Role, error)
		FindByMemberID(userID uint64) ([]*types.Role, error)
		Find(filter *types.RoleFilter) ([]*types.Role, error)

		Create(mod *types.Role) (*types.Role, error)
		Update(mod *types.Role) (*types.Role, error)

		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error

		MergeByID(id, targetRoleID uint64) error
		MoveByID(id, targetOrganisationID uint64) error

		MemberFindByRoleID(roleID uint64) ([]*types.RoleMember, error)
		MemberAddByID(id, userID uint64) error
		MemberRemoveByID(id, userID uint64) error
	}

	role struct {
		*repository

		// sql table reference
		roles   string
		members string
	}
)

const (
	sqlRoleScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrRoleNotFound = repositoryError("RoleNotFound")
)

func Role(ctx context.Context, db *factory.DB) RoleRepository {
	return (&role{}).With(ctx, db)
}

func (r *role) With(ctx context.Context, db *factory.DB) RoleRepository {
	return &role{
		repository: r.repository.With(ctx, db),
		roles:      "sys_role",
		members:    "sys_role_member",
	}
}

func (r *role) FindByID(id uint64) (*types.Role, error) {
	sql := "SELECT * FROM " + r.roles + " WHERE id = ? AND " + sqlRoleScope
	mod := &types.Role{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrRoleNotFound)
}

func (r *role) FindByMemberID(userID uint64) ([]*types.Role, error) {
	ids := make([]uint64, 0)
	params := make([]interface{}, 0)

	sql := "SELECT DISTINCT rel_role FROM " + r.members + "  "
	sql += "WHERE rel_user = ?"
	params = append(params, userID)

	if err := r.db().Select(&ids, sql, params...); err != nil {
		return nil, err
	}

	rval := make([]*types.Role, 0)
	for _, id := range ids {
		mod, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}
		rval = append(rval, mod)
	}

	return rval, nil
}

func (r *role) Find(filter *types.RoleFilter) ([]*types.Role, error) {
	rval := make([]*types.Role, 0)
	params := make([]interface{}, 0)

	sql := "SELECT * FROM " + r.roles + " WHERE " + sqlRoleScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *role) Create(mod *types.Role) (*types.Role, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.roles, mod)
}

func (r *role) Update(mod *types.Role) (*types.Role, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace(r.roles, mod)
}

func (r *role) ArchiveByID(id uint64) error {
	return r.updateColumnByID(r.roles, "archived_at", time.Now(), id)
}

func (r *role) UnarchiveByID(id uint64) error {
	return r.updateColumnByID(r.roles, "archived_at", nil, id)
}

func (r *role) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.roles, "deleted_at", time.Now(), id)
}

func (r *role) MergeByID(id, targetRoleID uint64) error {
	return ErrNotImplemented
}

func (r *role) MoveByID(id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}

func (r *role) MemberFindByRoleID(roleID uint64) (mm []*types.RoleMember, err error) {
	rval := make([]*types.RoleMember, 0)
	sql := "SELECT * FROM " + r.members + " WHERE rel_role = ?"
	return rval, r.db().Select(&rval, sql, roleID)
}

func (r *role) MemberAddByID(id, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: id,
		UserID: userID,
	}
	return r.db().Replace(r.members, mod)
}

func (r *role) MemberRemoveByID(id, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: id,
		UserID: userID,
	}
	return r.db().Delete(r.members, mod, "rel_role", "rel_user")
}
