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

func (r channel) FindByID(ctx context.Context, id uint64) (*types.Channel, error) {
	db := factory.Database.MustGet()
	mod := &types.Channel{}
	sql := "SELECT * FROM channels WHERE id = ? AND " + sqlChannelScope

	return mod, isFound(db.With(ctx).Get(mod, sql, id), mod.ID > 0, ErrChannelNotFound)
}

func (r channel) Find(ctx context.Context, filter *types.ChannelFilter) ([]*types.Channel, error) {
	db := factory.Database.MustGet()
	params := make([]interface{}, 0)
	rval := make([]*types.Channel, 0)

	sql := "SELECT * FROM channels WHERE " + sqlChannelScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, db.With(ctx).Select(&rval, sql, params...)
}

func (r channel) Create(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	if mod.Meta == nil {
		mod.SetMeta([]byte("{}"))
	}

	return mod, db.With(ctx).Insert("channels", mod)
}

func (r channel) Update(ctx context.Context, mod *types.Channel) (*types.Channel, error) {
	db := factory.Database.MustGet()

	now := time.Now()
	mod.SetUpdatedAt(&now)

	return mod, db.With(ctx).Replace("channels", mod)
}

func (r channel) AddMember(ctx context.Context, channelID, userID uint64) error {
	sql := `INSERT INTO channel_members (rel_channel, rel_user) VALUES (?, ?)`
	return exec(factory.Database.MustGet().With(ctx).Exec(sql, channelID, userID))
}

func (r channel) RemoveMember(ctx context.Context, channelID, userID uint64) error {
	sql := `DELETE FROM channel_members WHERE rel_channel = ? AND rel_user = ?`
	return exec(factory.Database.MustGet().With(ctx).Exec(sql, channelID, userID))
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
