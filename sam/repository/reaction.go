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
	sql := "SELECT * FROM reactions WHERE id = ?"
	mod := &types.Reaction{}

	return mod, isFound(db.With(ctx).Get(mod, sql, id), mod.ID > 0, ErrReactionNotFound)

}

func (r reaction) FindByRange(ctx context.Context, channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Reaction, 0)
	sql := `
		SELECT * 
          FROM reactions
         WHERE rel_reaction BETWEEN ? AND ?
           AND rel_channel = ?`

	return rval, db.With(ctx).Select(&rval, sql, fromReactionID, toReactionID, channelID)
}

func (r reaction) Create(ctx context.Context, mod *types.Reaction) (*types.Reaction, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	return mod, db.With(ctx).Insert("reactions", mod)
}

func (r reaction) Delete(ctx context.Context, id uint64) error {
	db := factory.Database.MustGet()

	return exec(db.With(ctx).Exec("DELETE FROM reactions WHERE id = ?", id))
}
