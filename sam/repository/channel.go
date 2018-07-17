package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlChannelScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrChannelNotFound = repositoryError("ChannelNotFound")
)

type (
	channel struct{}
)

func Channel() channel {
	return channel{}
}

func (r channel) FindById(ctx context.Context, id uint64) (*types.Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod := &types.Channel{}
	if err := db.Get(mod, "SELECT * FROM channels WHERE id = ? AND "+sqlChannelScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrChannelNotFound
	} else {
		return mod, nil
	}
}

func (r channel) Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	var params = make([]interface{}, 0)
	sql := "SELECT * FROM channels WHERE " + sqlChannelScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	rval := make([]*types.Channel, 0)
	if err := db.Select(&rval, sql, params...); err != nil {
		return nil, ErrDatabaseError
	} else {
		return rval, nil
	}
}

func (r channel) Create(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("channels", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r channel) Update(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	if err := db.Replace("channels", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r channel) Archive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "channels", "archived_at", time.Now(), id)
}

func (r channel) Unarchive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "channels", "archived_at", nil, id)
}

func (r channel) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "channels", id)
}
