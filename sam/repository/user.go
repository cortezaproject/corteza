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
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod := &types.User{}
	if err := db.Get(mod, "SELECT * FROM users WHERE username = ? AND "+sqlUserScope, username); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrUserNotFound
	} else {
		return mod, nil
	}
}

func (r user) FindById(ctx context.Context, id uint64) (*types.User, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod := &types.User{}
	if err := db.Get(mod, "SELECT * FROM users WHERE id = ? AND "+sqlUserScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrUserNotFound
	} else {
		return mod, nil
	}
}

func (r user) Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	var params = make([]interface{}, 0)
	sql := "SELECT * FROM users WHERE " + sqlUserScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND username LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY username ASC"

	rval := make([]*types.User, 0)
	if err := db.Select(&rval, sql, params...); err != nil {
		panic(err)
		return nil, ErrDatabaseError
	} else {
		return rval, nil
	}
}

func (r user) Create(ctx context.Context, mod *types.User) (*types.User, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("users", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r user) Update(ctx context.Context, mod *types.User) (*types.User, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	if err := db.Replace("users", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
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
