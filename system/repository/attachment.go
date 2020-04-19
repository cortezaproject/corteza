package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	AttachmentRepository interface {
		With(ctx context.Context, db *factory.DB) AttachmentRepository

		Find(filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		FindByID(attachmentID uint64) (*types.Attachment, error)
		Create(mod *types.Attachment) (*types.Attachment, error)
		DeleteByID(attachmentID uint64) error
	}

	attachment struct {
		*repository
	}
)

const (
	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

func Attachment(ctx context.Context, db *factory.DB) AttachmentRepository {
	return (&attachment{}).With(ctx, db)
}

func (r attachment) With(ctx context.Context, db *factory.DB) AttachmentRepository {
	return &attachment{
		repository: r.repository.With(ctx, db),
	}
}

func (r attachment) table() string {
	return "sys_attachment"
}

func (r attachment) columns() []string {
	return []string{
		"a.id",
		"a.rel_owner",
		"a.kind",
		"a.url",
		"a.preview_url",
		"a.name",
		"a.meta",
		"a.created_at",
		"a.updated_at",
		"a.deleted_at",
	}
}

func (r attachment) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS a").
		Where("a.deleted_at IS NULL")

}

func (r attachment) FindByID(attachmentID uint64) (*types.Attachment, error) {
	return r.findOneBy("id", attachmentID)
}

func (r attachment) findOneBy(field string, value interface{}) (*types.Attachment, error) {
	var (
		p = &types.Attachment{}
		q = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, p)
	)

	if err != nil {
		return nil, err
	} else if p.ID == 0 {
		return nil, ErrAttachmentNotFound
	}

	return p, nil
}

func (r attachment) Find(filter types.AttachmentFilter) (set types.AttachmentSet, f types.AttachmentFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id ASC"
	}

	query := r.query().
		Where(squirrel.Eq{"a.kind": f.Kind})

	if f.Filter != "" {
		err = errors.New("filtering by filter not implemented")
		return
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.PageFilter, &set)
}

func (r attachment) Create(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r attachment) DeleteByID(attachmentID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE id = ?",
		attachmentID,
	)

	return err
}
