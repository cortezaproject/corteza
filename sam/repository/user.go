package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	User interface {
		FindUserByUsername(username string) (*types.User, error)
		FindUserByID(id uint64) (*types.User, error)
		FindUsers(filter *types.UserFilter) ([]*types.User, error)
		CreateUser(mod *types.User) (*types.User, error)
		UpdateUser(mod *types.User) (*types.User, error)
		SuspendUserByID(id uint64) error
		UnsuspendUserByID(id uint64) error
		DeleteUserByID(id uint64) error
	}
)

const (
	sqlUserScope  = "suspended_at IS NULL AND deleted_at IS NULL"
	sqlUserSelect = "SELECT * FROM users WHERE " + sqlUserScope

	ErrUserNotFound = repositoryError("UserNotFound")
)

func (r *repository) FindUserByUsername(username string) (*types.User, error) {
	sql := "SELECT * FROM users WHERE username = ? AND " + sqlUserScope
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, username), mod.ID > 0, ErrUserNotFound)
}

func (r *repository) FindUserByID(id uint64) (*types.User, error) {
	sql := "SELECT * FROM users WHERE id = ? AND " + sqlUserScope
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrUserNotFound)
}

func (r *repository) FindUsers(filter *types.UserFilter) ([]*types.User, error) {
	rval := make([]*types.User, 0)
	params := make([]interface{}, 0)
	sql := "SELECT * FROM users WHERE " + sqlUserScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND username LIKE ?"
			params = append(params, filter.Query+"%")
		}

		if filter.MembersOfChannel > 0 {
			sql += " AND id IN (SELECT rel_user FROM channel_members WHERE rel_channel = ?)"
			params = append(params, filter.MembersOfChannel)
		}
	}

	sql += " ORDER BY username ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *repository) CreateUser(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, r.db().Insert("users", mod)
}

func (r *repository) UpdateUser(mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, r.db().Replace("users", mod)
}

func (r *repository) SuspendUserByID(id uint64) error {
	return simpleUpdate(r.ctx, "users", "suspend_at", time.Now(), id)
}

func (r *repository) UnsuspendUserByID(id uint64) error {
	return simpleUpdate(r.ctx, "users", "suspend_at", nil, id)
}

func (r *repository) DeleteUserByID(id uint64) error {
	return simpleDelete(r.ctx, "users", id)
}
