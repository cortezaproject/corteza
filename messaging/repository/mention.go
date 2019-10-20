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
	MentionRepository interface {
		With(ctx context.Context, db *factory.DB) MentionRepository

		FindByUserIDs(IDs ...uint64) (mm types.MentionSet, err error)
		FindByMessageIDs(IDs ...uint64) (mm types.MentionSet, err error)
		Create(m *types.Mention) (*types.Mention, error)
		DeleteByMessageID(ID uint64) error
		DeleteByID(ID uint64) error
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

func (r mention) With(ctx context.Context, db *factory.DB) MentionRepository {
	return &mention{
		repository: r.repository.With(ctx, db),
	}
}

func (r mention) table() string {
	return "messaging_mention"
}

func (r mention) columns() []string {
	return []string{
		"mm.id",
		"mm.rel_message",
		"mm.rel_channel",
		"mm.rel_user",
		"mm.rel_mentioned_by",
		"mm.created_at",
	}
}

func (r mention) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS mm")
}

func (r mention) FindByUserIDs(IDs ...uint64) (types.MentionSet, error) {
	return r.findAllBy(squirrel.Eq{"rel_user": IDs})
}

func (r mention) FindByMessageIDs(IDs ...uint64) (types.MentionSet, error) {
	return r.findAllBy(squirrel.Eq{"rel_message": IDs})
}

func (r mention) findAllBy(cnd squirrel.Sqlizer) (mm types.MentionSet, err error) {
	return mm, rh.FetchAll(r.db(), r.query().Where(cnd), &mm)
}

func (r mention) Create(m *types.Mention) (*types.Mention, error) {
	m.ID = factory.Sonyflake.NextID()
	m.CreatedAt = time.Now()
	return m, r.db().Insert(r.table(), m)
}

func (r mention) DeleteByMessageID(ID uint64) error {
	return rh.Delete(r.db(), r.table(), squirrel.Eq{"rel_message": ID})
}

func (r mention) DeleteByID(ID uint64) error {
	return rh.Delete(r.db(), r.table(), squirrel.Eq{"id": ID})
}
