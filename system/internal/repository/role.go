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
		MemberAddByID(roleID, userID uint64) error
		MemberRemoveByID(roleID, userID uint64) error

		Reset() error
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
	sql := "SELECT * FROM " + r.roles + " where id in (select rel_role from " + r.members + " where rel_user=?) and " + sqlRoleScope
	rval := make([]*types.Role, 0)
	if err := r.db().Select(&rval, sql, userID); err != nil {
		return nil, err
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

func (r *role) MemberAddByID(roleID, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: roleID,
		UserID: userID,
	}
	return r.db().Replace(r.members, mod)
}

func (r *role) MemberRemoveByID(roleID, userID uint64) error {
	mod := &types.RoleMember{
		RoleID: roleID,
		UserID: userID,
	}
	return r.db().Delete(r.members, mod, "rel_role", "rel_user")
}

func (r *role) Reset() error {
	sql := `REPLACE INTO sys_role (id, name, handle) VALUES
		(1, 'Everyone', 'everyone'),
		(2, 'Administrators', 'admins')
	`
	_, err := r.db().Exec(sql)
	if err != nil {
		return err
	}

	// Value: Allow (2), Deny (1), Inherit(0)
	sql = `REPLACE INTO sys_rules (rel_role, resource, operation, value) VALUES
		-- Everyone
		(1, 'system', 'user.create', 2),
		(1, 'compose:*', 'access', 2),
		(1, 'messaging:*', 'access', 2),
		-- Admins
		(2, 'compose', 'namespace.create', 2),
		(2, 'compose', 'access', 2),
		(2, 'compose', 'grant', 2),
		(2, 'compose:namespace:*', 'page.create', 2),
		(2, 'compose:namespace:*', 'read', 2),
		(2, 'compose:namespace:*', 'update', 2),
		(2, 'compose:namespace:*', 'delete', 2),
		(2, 'compose:namespace:*', 'module.create', 2),
		(2, 'compose:namespace:*', 'chart.create', 2),
		(2, 'compose:namespace:*', 'trigger.create', 2),
		(2, 'compose:chart:*', 'read', 2),
		(2, 'compose:chart:*', 'update', 2),
		(2, 'compose:chart:*', 'delete', 2),
		(2, 'compose:trigger:*', 'read', 2),
		(2, 'compose:trigger:*', 'update', 2),
		(2, 'compose:trigger:*', 'delete', 2),
		(2, 'compose:page:*', 'read', 2),
		(2, 'compose:page:*', 'update', 2),
		(2, 'compose:page:*', 'delete', 2),
		(2, 'system', 'access', 2),
		(2, 'system', 'grant', 2),
		(2, 'system', 'organisation.create', 2),
		(2, 'system', 'user.create', 2),
		(2, 'system', 'role.create', 2),
		(2, 'system:organisation:*', 'access', 2),
		(2, 'system:user:*', 'read', 2),
		(2, 'system:user:*', 'update', 2),
		(2, 'system:user:*', 'suspend', 2),
		(2, 'system:user:*', 'unsuspend', 2),
		(2, 'system:user:*', 'delete', 2),
		(2, 'system:role:*', 'read', 2),
		(2, 'system:role:*', 'update', 2),
		(2, 'system:role:*', 'delete', 2),
		(2, 'system:role:*', 'members.manage', 2),
		(2, 'messaging', 'access', 2),
		(2, 'messaging', 'grant', 2),
		(2, 'messaging', 'channel.public.create', 2),
		(2, 'messaging', 'channel.private.create', 2),
		(2, 'messaging', 'channel.group.create', 2),
		(2, 'messaging:channel:*', 'update', 2),
		(2, 'messaging:channel:*', 'leave', 2),
		(2, 'messaging:channel:*', 'read', 2),
		(2, 'messaging:channel:*', 'join', 2),
		(2, 'messaging:channel:*', 'delete', 2),
		(2, 'messaging:channel:*', 'undelete', 2),
		(2, 'messaging:channel:*', 'archive', 2),
		(2, 'messaging:channel:*', 'unarchive', 2),
		(2, 'messaging:channel:*', 'members.manage', 2),
		(2, 'messaging:channel:*', 'webhooks.manage', 2),
		(2, 'messaging:channel:*', 'attachments.manage', 2),
		(2, 'messaging:channel:*', 'message.attach', 2),
		(2, 'messaging:channel:*', 'message.update.all', 2),
		(2, 'messaging:channel:*', 'message.update.own', 2),
		(2, 'messaging:channel:*', 'message.delete.all', 2),
		(2, 'messaging:channel:*', 'message.delete.own', 2),
		(2, 'messaging:channel:*', 'message.embed', 2),
		(2, 'messaging:channel:*', 'message.send', 2),
		(2, 'messaging:channel:*', 'message.reply', 2),
		(2, 'messaging:channel:*', 'message.react', 2),
		(2, 'compose:module:*', 'read', 2),
		(2, 'compose:module:*', 'update', 2),
		(2, 'compose:module:*', 'delete', 2),
		(2, 'compose:module:*', 'record.create', 2),
		(2, 'compose:module:*', 'record.read', 2),
		(2, 'compose:module:*', 'record.update', 2),
		(2, 'compose:module:*', 'record.delete', 2)
	`
	_, err = r.db().Exec(sql)
	return err
}
