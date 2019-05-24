package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MentionRepository interface {
		With(ctx context.Context, db *factory.DB) MentionRepository

		FindByUserIDs(IDs ...uint64) (mm types.MentionSet, err error)
		FindByMessageIDs(IDs ...uint64) (mm types.MentionSet, err error)
		Create(m *types.Mention) (*types.Mention, error)
		DeleteByMessageID(ID uint64) error
		DeleteByID(ID uint64) error

		CountMentions(userID uint64) (c int, err error)
		ChangeMention(userID, target uint64) error
	}

	mention struct {
		*repository
	}
)

var (
	ErrMentionNotFound = repositoryError("MentionNotFound")
)

func Mention(ctx context.Context, db *factory.DB) MentionRepository {
	return (&mention{}).With(ctx, db)
}

func (r *mention) With(ctx context.Context, db *factory.DB) MentionRepository {
	return &mention{
		repository: r.repository.With(ctx, db),
	}
}

func (r *mention) FindByUserIDs(IDs ...uint64) (types.MentionSet, error) {
	return r.findByIDs("rel_user", IDs...)
}

func (r *mention) FindByMessageIDs(IDs ...uint64) (types.MentionSet, error) {
	return r.findByIDs("rel_message", IDs...)
}

func (r *mention) findByIDs(col string, IDs ...uint64) (mm types.MentionSet, err error) {
	mm = types.MentionSet{}

	if len(IDs) == 0 {
		return
	}

	sql := fmt.Sprintf(`SELECT * FROM messaging_mention WHERE %s IN (?)`, col)

	if sql, args, err := sqlx.In(sql, IDs); err != nil {
		return nil, err
	} else {
		return mm, r.db().Select(&mm, sql, args...)
	}
}

func (r *mention) Create(m *types.Mention) (*types.Mention, error) {
	m.ID = factory.Sonyflake.NextID()
	m.CreatedAt = time.Now()
	return m, r.db().Insert("messaging_mention", m)
}

func (r *mention) DeleteByMessageID(ID uint64) error {
	return exec(r.db().Exec("DELETE FROM messaging_mention WHERE rel_message = ?", ID))
}

func (r *mention) DeleteByID(ID uint64) error {
	return exec(r.db().Exec("DELETE FROM messaging_mention WHERE id = ?", ID))
}

func (r *mention) CountMentions(userID uint64) (c int, err error) {
	return c, r.db().Get(&c,
		"SELECT COUNT(*) FROM messaging_mention WHERE rel_user = ?",
		userID)
}

func (r *mention) ChangeMention(userID, target uint64) error {
	_, err := r.db().Exec("UPDATE messaging_mention SET rel_user = ? WHERE rel_user = ?", target, userID)
	return err
}
