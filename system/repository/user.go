package repository

import (
	"context"
	"time"

	"github.com/crusttech/crust/system/types"
	"github.com/titpetric/factory"
)

type (
	UserRepository interface {
		With(ctx context.Context, db *factory.DB) UserRepository

		FindUserByEmail(email string) (*types.User, error)
		FindUserByUsername(username string) (*types.User, error)
		FindUserByID(id uint64) (*types.User, error)
		FindUserBySatosaID(id string) (*types.User, error)
		FindUsers(filter *types.UserFilter) ([]*types.User, error)
		CreateUser(mod *types.User) (*types.User, error)
		UpdateUser(mod *types.User) (*types.User, error)
		SuspendUserByID(id uint64) error
		UnsuspendUserByID(id uint64) error
		DeleteUserByID(id uint64) error
	}

	user struct {
		*repository
	}
)

const (
	sqlUserColumns = "id, email, username, password, name, handle, " +
		"meta, satosa_id, rel_organisation, " +
		"created_at, updated_at, suspended_at, deleted_at"
	sqlUserScope  = "suspended_at IS NULL AND deleted_at IS NULL"
	sqlUserSelect = "SELECT " + sqlUserColumns + " FROM users WHERE " + sqlUserScope

	ErrUserNotFound = repositoryError("UserNotFound")
)

func User(ctx context.Context, db *factory.DB) UserRepository {
	return (&user{}).With(ctx, db)
}

func (r *user) With(ctx context.Context, db *factory.DB) UserRepository {
	return &user{repository: r.repository.With(ctx, db)}
}

func (r *user) FindUserByUsername(username string) (*types.User, error) {
	sql := sqlUserSelect + " AND username = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, username), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindUserBySatosaID(satosaID string) (*types.User, error) {
	sql := sqlUserSelect + " AND satosa_id = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, satosaID), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindUserByEmail(email string) (*types.User, error) {
	sql := sqlUserSelect + " AND email = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, email), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindUserByID(id uint64) (*types.User, error) {
	sql := sqlUserSelect + " AND id = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindUsers(filter *types.UserFilter) ([]*types.User, error) {
	rval := make([]*types.User, 0)
	params := make([]interface{}, 0)
	sql := sqlUserSelect

	if filter != nil {
		if filter.Query != "" {
			sql += " AND username LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY username ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *user) CreateUser(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))
	return mod, r.db().Insert("users", mod)
}

func (r *user) UpdateUser(mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, r.db().Replace("users", mod)
}

func (r *user) SuspendUserByID(id uint64) error {
	return r.updateColumnByID("users", "suspend_at", time.Now(), id)
}

func (r *user) UnsuspendUserByID(id uint64) error {
	return r.updateColumnByID("users", "suspend_at", nil, id)
}

func (r *user) DeleteUserByID(id uint64) error {
	return r.updateColumnByID("users", "deleted_at", time.Now(), id)
}
