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

		FindByEmail(email string) (*types.User, error)
		FindByUsername(username string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindBySatosaID(id string) (*types.User, error)
		Find(filter *types.UserFilter) ([]*types.User, error)

		Create(mod *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)
		SuspendByID(id uint64) error
		UnsuspendByID(id uint64) error
		DeleteByID(id uint64) error
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

func (r *user) FindByUsername(username string) (*types.User, error) {
	sql := sqlUserSelect + " AND username = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, username), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindBySatosaID(satosaID string) (*types.User, error) {
	sql := sqlUserSelect + " AND satosa_id = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, satosaID), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindByEmail(email string) (*types.User, error) {
	sql := sqlUserSelect + " AND email = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, email), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindByID(id uint64) (*types.User, error) {
	sql := sqlUserSelect + " AND id = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrUserNotFound)
}

func (r *user) Find(filter *types.UserFilter) ([]*types.User, error) {
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

func (r *user) Create(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert("users", mod)
}

func (r *user) Update(mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	return mod, r.db().Replace("users", mod)
}

func (r *user) SuspendByID(id uint64) error {
	return r.updateColumnByID("users", "suspend_at", time.Now(), id)
}

func (r *user) UnsuspendByID(id uint64) error {
	return r.updateColumnByID("users", "suspend_at", nil, id)
}

func (r *user) DeleteByID(id uint64) error {
	return r.updateColumnByID("users", "deleted_at", time.Now(), id)
}
