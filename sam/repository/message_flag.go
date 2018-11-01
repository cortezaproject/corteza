package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	MessageFlagRepository interface {
		With(ctx context.Context, db *factory.DB) MessageFlagRepository

		FindByID(ID uint64) (*types.MessageFlag, error)
		FindByMessageIDs(IDs ...uint64) ([]*types.MessageFlag, error)
		FindByFlag(messageID, userID uint64, flag string) (*types.MessageFlag, error)
		Create(mod *types.MessageFlag) (*types.MessageFlag, error)
		DeleteByID(ID uint64) error
	}

	reaction struct {
		*repository
	}
)

const (
	ErrMessageFlagNotFound = repositoryError("MessageFlagNotFound")
)

func MessageFlag(ctx context.Context, db *factory.DB) MessageFlagRepository {
	return (&reaction{}).With(ctx, db)
}

func (r *reaction) With(ctx context.Context, db *factory.DB) MessageFlagRepository {
	return &reaction{
		repository: r.repository.With(ctx, db),
	}
}

func (r *reaction) FindByID(ID uint64) (*types.MessageFlag, error) {
	sql := "SELECT * FROM message_flags WHERE id = ?"
	mod := &types.MessageFlag{}
	return mod, isFound(r.db().Get(mod, sql, ID), mod.ID > 0, ErrMessageFlagNotFound)
}

func (r *reaction) FindByFlag(messageID, userID uint64, flag string) (*types.MessageFlag, error) {
	args := []interface{}{messageID, flag}
	sql := "SELECT * FROM message_flags WHERE rel_message = ? AND flag = ? "

	if userID > 0 {
		sql += "AND rel_user = ? "
		args = append(args, userID)
	}

	mod := &types.MessageFlag{}
	return mod, isFound(r.db().Get(mod, sql, args...), mod.ID > 0, ErrMessageFlagNotFound)
}

// FindByMessageRange returns all flags by message id range
func (r *reaction) FindByMessageIDs(IDs ...uint64) ([]*types.MessageFlag, error) {
	rval := make([]*types.MessageFlag, 0)

	sql := `SELECT * FROM message_flags WHERE rel_message IN (?)`

	if sql, args, err := sqlx.In(sql, IDs); err != nil {
		return nil, err
	} else {
		return rval, r.db().Select(&rval, sql, args...)
	}
}

func (r *reaction) Create(mod *types.MessageFlag) (*types.MessageFlag, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert("message_flags", mod)
}

func (r *reaction) DeleteByID(ID uint64) error {
	return exec(r.db().Exec("DELETE FROM message_flags WHERE id = ?", ID))
}
