package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	Reaction interface {
		With(ctx context.Context) Reaction

		FindReactionByID(id uint64) (*types.Reaction, error)
		FindReactionsByRange(channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error)
		CreateReaction(mod *types.Reaction) (*types.Reaction, error)
		DeleteReactionByID(id uint64) error
	}

	reaction struct {
		*repository
	}
)

const (
	ErrReactionNotFound = repositoryError("ReactionNotFound")
)

func NewReaction(ctx context.Context) Reaction {
	return (&reaction{}).With(ctx)
}

func (r *reaction) With(ctx context.Context) Reaction {
	return &reaction{
		repository: r.repository.With(ctx),
	}
}

func (r *reaction) FindReactionByID(id uint64) (*types.Reaction, error) {
	sql := "SELECT * FROM reactions WHERE id=?"
	mod := &types.Reaction{}
	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrReactionNotFound)
}

func (r *reaction) FindReactionsByRange(channelID, fromReactionID, toReactionID uint64) ([]*types.Reaction, error) {
	rval := make([]*types.Reaction, 0)
	sql := `SELECT * FROM reactions WHERE rel_reaction BETWEEN ? AND ? AND rel_channel=?`
	return rval, r.db().Select(&rval, sql, fromReactionID, toReactionID, channelID)
}

func (r *reaction) CreateReaction(mod *types.Reaction) (*types.Reaction, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert("reactions", mod)
}

func (r *reaction) DeleteReactionByID(id uint64) error {
	return exec(r.db().Exec("DELETE FROM reactions WHERE id=?", id))
}
