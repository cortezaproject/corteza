package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlUserScope  = "suspended_at IS NULL AND deleted_at IS NULL"
	sqlUserSelect = "SELECT * FROM users WHERE " + sqlUserScope

	ErrUserNotFound = repositoryError("UserNotFound")
)

type (
	user struct{}
)

func User() user {
	return user{}
}

func (r user) FindByUsername(ctx context.Context, username string) (*types.User, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM users WHERE username = ? AND " + sqlUserScope
	mod := &types.User{}

	return mod, isFound(db.Get(mod, sql, username), mod.ID > 0, ErrUserNotFound)
}

func (r user) FindByID(ctx context.Context, id uint64) (*types.User, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM users WHERE id = ? AND " + sqlUserScope
	mod := &types.User{}

	return mod, isFound(db.Get(mod, sql, id), mod.ID > 0, ErrUserNotFound)
}

func (r user) Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error) {
	db := factory.Database.MustGet()
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

	return rval, db.With(ctx).Select(&rval, sql, params...)
}

func (r user) Create(ctx context.Context, mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, factory.Database.MustGet().With(ctx).Insert("users", mod)
}

func (r user) Update(ctx context.Context, mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, factory.Database.MustGet().With(ctx).Replace("users", mod)
}

func (r user) Suspend(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "users", "suspend_at", time.Now(), id)
}

func (r user) Unsuspend(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "users", "suspend_at", nil, id)
}

func (r user) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "users", id)
}
