package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	ErrReactionNotFound = repositoryError("ReactionNotFound")
)

type (
	reaction struct{}
)

func Reaction() reaction {
	return reaction{}
}

func (r reaction) FindByID(ctx context.Context, id uint64) (*types.Reaction, error) {
	db := factory.Database.MustGet()

	mod := &types.Reaction{}
	if err := db.GetContext(ctx, mod, "SELECT * FROM reactions WHERE id = ?", id); err != nil {
		return nil, err
	} else if mod.ID == 0 {
		return nil, ErrReactionNotFound
	} else {
		return mod, nil
	}
}

func (r reaction) FindByRange(ctx context.Context, channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error) {
	db := factory.Database.MustGet()

	sql := `
		SELECT * 
          FROM reactions
         WHERE rel_reaction BETWEEN ? AND ?
           AND rel_channel = ?`

	rval := make([]*types.Reaction, 0)
	if err := db.Select(&rval, sql, fromReactionID, toReactionID, channelID); err != nil {
		return nil, err
	}

	return rval, nil
}

func (r reaction) Create(ctx context.Context, mod *types.Reaction) (*types.Reaction, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	if err := db.Insert("reactions", mod); err != nil {
		return nil, err
	} else {
		return mod, nil
	}
}

func (r reaction) Delete(ctx context.Context, id uint64) error {
	db := factory.Database.MustGet()

	if _, err := db.ExecContext(ctx, "DELETE FROM reactions WHERE id = ?", id); err != nil {
		return err
	} else {
		return nil
	}
}
