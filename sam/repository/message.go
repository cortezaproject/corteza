package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
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

	mod := &types.Message{}
	if err := db.GetContext(ctx, mod, "SELECT * FROM messages WHERE id = ? AND "+sqlMessageScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrMessageNotFound
	} else {
		return mod, nil
	}
}

func (r message) Find(ctx context.Context, filter *types.MessageFilter) ([]*types.Message, error) {
	db := factory.Database.MustGet()

	var params = make([]interface{}, 0)
	sql := "SELECT * FROM messages WHERE " + sqlMessageScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	rval := make([]*types.Message, 0)
	if err := db.SelectContext(ctx, &rval, sql, params...); err != nil {
		return nil, ErrDatabaseError
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
