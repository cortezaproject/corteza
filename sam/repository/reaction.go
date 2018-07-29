package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Reaction interface {
		FindReactionByID(id uint64) (*types.Reaction, error)
		FindReactionsByRange(channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error)
		CreateReaction(mod *types.Reaction) (*types.Reaction, error)
		DeleteReactionByID(id uint64) error
	}
)

const (
	ErrReactionNotFound = repositoryError("ReactionNotFound")
)

func (r *repository) FindReactionByID(id uint64) (*types.Reaction, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM reactions WHERE id = ?"
	mod := &types.Reaction{}

	return mod, isFound(db.With(r.ctx).Get(mod, sql, id), mod.ID > 0, ErrReactionNotFound)

}

func (r *repository) FindReactionsByRange(channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Reaction, 0)
	sql := `
		SELECT * 
          FROM reactions
         WHERE rel_reaction BETWEEN ? AND ?
           AND rel_channel = ?`

	return rval, db.With(r.ctx).Select(&rval, sql, fromReactionID, toReactionID, channelID)
}

func (r *repository) CreateReaction(mod *types.Reaction) (*types.Reaction, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, factory.Database.MustGet().With(r.ctx).Insert("reactions", mod)
}

func (r *repository) DeleteReactionByID(id uint64) error {
	db := factory.Database.MustGet()
	return exec(db.With(r.ctx).Exec("DELETE FROM reactions WHERE id = ?", id))
}
