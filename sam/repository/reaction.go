package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
)

const (
	sqlReactionScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrReactionNotFound = repositoryError("ReactionNotFound")
)

type (
	reaction struct{}
)

func Reaction() reaction {
	return reaction{}
}

func (r reaction) FindById(ctx context.Context, id uint64) (*types.Reaction, error) {
	db := factory.Database.MustGet()

	mod := &types.Reaction{}
	if err := db.GetContext(ctx, mod, "SELECT * FROM reactions WHERE id = ? AND "+sqlReactionScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrReactionNotFound
	} else {
		return mod, nil
	}
}

func (r reaction) FindByRange(ctx context.Context, channelId, fromReactionId, toReactionId uint64) ([]*types.Reaction, error) {
	db := factory.Database.MustGet()

	sql := `
		SELECT * 
          FROM reactions
         WHERE rel_reaction BETWEEN ? AND ?
           AND rel_channel = ?`

	rval := make([]*types.Reaction, 0)
	if err := db.Select(&rval, sql, fromReactionId, toReactionId, channelId); err != nil {
		return nil, ErrDatabaseError
	}

	return rval, nil
}

func (r reaction) Create(ctx context.Context, mod *types.Reaction) (*types.Reaction, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("reactions", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r reaction) Delete(ctx context.Context, id uint64) error {
	db := factory.Database.MustGet()

	if _, err := db.ExecContext(ctx, "DELETE FROM reactions WHERE id = ?", id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}
