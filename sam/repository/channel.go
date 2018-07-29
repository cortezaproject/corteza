package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Channel interface {
		FindChannelByID(id uint64) (*types.Channel, error)
		FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error)
		CreateChannel(mod *types.Channel) (*types.Channel, error)
		UpdateChannel(mod *types.Channel) (*types.Channel, error)
		AddChannelMember(channelID, userID uint64) error
		RemoveChannelMember(channelID, userID uint64) error
		ArchiveChannelByID(id uint64) error
		UnarchiveChannelByID(id uint64) error
		DeleteChannelByID(id uint64) error
	}
)

const (
	sqlChannelScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrChannelNotFound = repositoryError("ChannelNotFound")
)

func (r *repository) FindChannelByID(id uint64) (*types.Channel, error) {
	db := factory.Database.MustGet()
	mod := &types.Channel{}
	sql := "SELECT * FROM channels WHERE id = ? AND " + sqlChannelScope

	return mod, isFound(db.With(r.ctx).Get(mod, sql, id), mod.ID > 0, ErrChannelNotFound)
}

func (r *repository) FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error) {
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

	return rval, db.With(r.ctx).Select(&rval, sql, params...)
}

func (r *repository) CreateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, factory.Database.MustGet().With(r.ctx).Insert("channels", mod)
}

func (r *repository) UpdateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	return mod, factory.Database.MustGet().With(r.ctx).Replace("channels", mod)
}

func (r *repository) AddChannelMember(channelID, userID uint64) error {
	sql := `INSERT INTO channel_members (rel_channel, rel_user) VALUES (?, ?)`
	return exec(factory.Database.MustGet().With(r.ctx).Exec(sql, channelID, userID))
}

func (r *repository) RemoveChannelMember(channelID, userID uint64) error {
	sql := `DELETE FROM channel_members WHERE rel_channel = ? AND rel_user = ?`
	return exec(factory.Database.MustGet().With(r.ctx).Exec(sql, channelID, userID))
}

func (r *repository) ArchiveChannelByID(id uint64) error {
	return simpleUpdate(r.ctx, "channels", "archived_at", time.Now(), id)
}

func (r *repository) UnarchiveChannelByID(id uint64) error {
	return simpleUpdate(r.ctx, "channels", "archived_at", nil, id)
}

func (r *repository) DeleteChannelByID(id uint64) error {
	return simpleDelete(r.ctx, "channels", id)
}
