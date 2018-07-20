package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
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

func (r message) FindById(ctx context.Context, id uint64) (*types.Message, error) {
	db := factory.Database.MustGet()

	sql := "SELECT id, COALESCE(type,'') AS type, message, rel_user, rel_channel, COALESCE(reply_to, 0) AS reply_to FROM messages WHERE id = ? AND " + sqlMessageScope

	mod := &types.Message{}
	if err := db.GetContext(ctx, mod, sql, id); err != nil {
		return nil, errors.Wrap(err, ErrDatabaseError.String())
	} else if mod.ID == 0 {
		return nil, ErrMessageNotFound
	} else {
		return mod, nil
	}
}

func (r message) Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error) {
	db := factory.Database.MustGet()

	var params = make([]interface{}, 0)
	sql := "SELECT id, COALESCE(type,'') AS type, message, rel_user, rel_channel, COALESCE(reply_to, 0) AS reply_to FROM messages WHERE " + sqlMessageScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	if filter.ChannelId > 0 {
		sql += " AND rel_channel = ? "
		params = append(params, filter.ChannelId)
	}

	if filter.LastMessageId > 0 {
		sql += " AND id > ? "
		params = append(params, filter.LastMessageId)
	}

	sql += " ORDER BY id ASC"

	rval := make([]*types.Message, 0)
	if err := db.SelectContext(ctx, &rval, sql, params...); err != nil {
		return nil, errors.Wrap(err, ErrDatabaseError.String())
	} else {
		return rval, nil
	}
}

func (r message) Create(ctx context.Context, mod *types.Message) (*types.Message, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("messages", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r message) Update(ctx context.Context, mod *types.Message) (*types.Message, error) {
	db := factory.Database.MustGet()

	if err := db.Replace("messages", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r message) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "messages", id)
}
