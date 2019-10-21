package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	MessageFlagRepository interface {
		With(ctx context.Context, db *factory.DB) MessageFlagRepository

		FindByID(ID uint64) (*types.MessageFlag, error)
		FindByMessageIDs(IDs ...uint64) (types.MessageFlagSet, error)
		FindByFlag(messageID, userID uint64, flag string) (*types.MessageFlag, error)
		Create(mod *types.MessageFlag) (*types.MessageFlag, error)
		DeleteByID(ID uint64) error
	}

	messageFlag struct {
		*repository
	}
)

const (
	ErrMessageFlagNotFound = repositoryError("MessageFlagNotFound")
)

func MessageFlag(ctx context.Context, db *factory.DB) MessageFlagRepository {
	return (&messageFlag{}).With(ctx, db)
}

func (r messageFlag) columns() []string {
	return []string{
		"mf.id",
		"mf.rel_user",
		"mf.rel_message",
		"mf.rel_channel",
		"mf.flag",
		"mf.created_at",
	}
}

func (r messageFlag) table() string {
	return "messaging_message_flag"
}

func (r messageFlag) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS mf")
}

func (r messageFlag) queryMessagesWithFlags(flags ...string) squirrel.SelectBuilder {
	return squirrel.
		Select("mf.rel_message").
		From(r.table() + " AS mf").
		Where(squirrel.Eq{"flag": flags})
}

func (r messageFlag) With(ctx context.Context, db *factory.DB) MessageFlagRepository {
	return &messageFlag{
		repository: r.repository.With(ctx, db),
	}
}

func (r messageFlag) FindByID(ID uint64) (*types.MessageFlag, error) {
	return r.findOneBy(squirrel.Eq{"id": ID})
}

func (r messageFlag) FindByFlag(messageID, userID uint64, flag string) (*types.MessageFlag, error) {
	cnd := squirrel.Eq{
		"rel_message": messageID,
		"flag":        flag,
	}

	if userID > 0 {
		cnd["rel_user"] = userID
	}

	return r.findOneBy(cnd)
}

func (r messageFlag) findOneBy(cnd squirrel.Sqlizer) (*types.MessageFlag, error) {
	var (
		mf = &types.MessageFlag{}

		q = r.query().
			Where(cnd)

		err = rh.FetchOne(r.db(), q, mf)
	)

	if err != nil {
		return nil, err
	} else if mf.ID == 0 {
		return nil, ErrMessageFlagNotFound
	}

	return mf, nil
}

// FindByMessageIDs returns all flags by message id range
func (r messageFlag) FindByMessageIDs(IDs ...uint64) (set types.MessageFlagSet, err error) {
	if len(IDs) == 0 {
		return
	}

	return set, rh.FetchAll(r.db(), r.query().Where(squirrel.Eq{"rel_message": IDs}), &set)
}

func (r messageFlag) Create(mod *types.MessageFlag) (*types.MessageFlag, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert(r.table(), mod)
}

func (r messageFlag) DeleteByID(ID uint64) error {
	return rh.Delete(r.db(), r.table(), squirrel.Eq{"id": ID})
}
