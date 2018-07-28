package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlMessageScope = "deleted_at IS NULL"

	ErrMessageNotFound = repositoryError("MessageNotFound")
)

type (
	message struct{}
)

func Message() message {
	return message{}
}

func (r message) FindByID(ctx context.Context, id uint64) (*types.Message, error) {
	db := factory.Database.MustGet()
	mod := &types.Message{}
	sql := "SELECT id, COALESCE(type,'') AS type, message, rel_user, rel_channel, COALESCE(reply_to, 0) AS reply_to FROM messages WHERE id = ? AND " + sqlMessageScope

	return mod, isFound(db.With(ctx).Get(mod, sql, id), mod.ID > 0, ErrMessageNotFound)
}

func (r message) Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error) {
	db := factory.Database.MustGet()
	params := make([]interface{}, 0)
	rval := make([]*types.Message, 0)

	sql := "SELECT id, COALESCE(type,'') AS type, message, rel_user, rel_channel, COALESCE(reply_to, 0) AS reply_to FROM messages WHERE " + sqlMessageScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND message LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	if filter.ChannelID > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelID)
	}

	if filter.FromMessageID > 0 {
		sql += " AND id > ? "
		params = append(params, filter.FromMessageID)
	}

	if filter.UntilMessageID > 0 {
		sql += " AND id < ? "
		params = append(params, filter.UntilMessageID)
	}

	sql += " ORDER BY id ASC"

	if filter.Limit > 0 {
		// @todo implement some kind of protection
		sql += " LIMIT ? "
		params = append(params, filter.Limit)
	}
	return rval, db.With(ctx).Select(&rval, sql, params...)
}

func (r message) Create(ctx context.Context, mod *types.Message) (*types.Message, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	return mod, db.With(ctx).Insert("messages", mod)
}

func (r message) Update(ctx context.Context, mod *types.Message) (*types.Message, error) {
	db := factory.Database.MustGet()

	now := time.Now()
	mod.SetUpdatedAt(&now)

	return mod, db.With(ctx).Replace("messages", mod)
}

func (r message) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "messages", id)
}
