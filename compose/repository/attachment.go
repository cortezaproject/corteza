package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	AttachmentRepository interface {
		With(ctx context.Context, db *factory.DB) AttachmentRepository

		Find(filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		FindByID(namespaceID, attachmentID uint64) (*types.Attachment, error)
		Create(mod *types.Attachment) (*types.Attachment, error)
		DeleteByID(namespaceID, attachmentID uint64) error
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
	return "compose_attachment"
}

func (r attachment) columns() []string {
	return []string{
		"a.id", "a.rel_namespace", "a.rel_owner", "a.kind",
		"a.url", "a.preview_url",
		"a.name",
		"a.meta",
		"a.created_at", "a.updated_at", "a.deleted_at",
	}
}

func (r attachment) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table() + " AS a").
		Where("a.deleted_at IS NULL")

}

func (r attachment) FindByID(namespaceID, attachmentID uint64) (*types.Attachment, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("a.id = ?", attachmentID)

		a = &types.Attachment{}
	)

	if namespaceID > 0 {
		query = query.Where("a.rel_namespace = ?", namespaceID)
	}

	return a, isFound(r.fetchOne(a, query), a.ID > 0, ErrAttachmentNotFound)
}

func (r attachment) Find(filter types.AttachmentFilter) (set types.AttachmentSet, f types.AttachmentFilter, err error) {
	f = filter
	// f.PerPage = normalizePerPage(f.PerPage, 5, 100, 50)

	query := r.query().
		Where(squirrel.Eq{"a.kind": f.Kind})

	if filter.NamespaceID > 0 {
		query = query.Where("a.rel_namespace = ?", filter.NamespaceID)
	}

	switch f.Kind {
	case types.PageAttachment:
		// @todo implement filtering by page
		if f.PageID > 0 {
			err = errors.New("filtering by pageID not implemented")
			return
		}
	case types.RecordAttachment:
		query = query.
			Join("compose_record_value AS v ON (v.ref = a.id)")

		if f.ModuleID > 0 {
			query = query.
				Join("compose_record AS r ON (r.id = v.record_id)").
				Where(squirrel.Eq{"r.module_id": f.ModuleID})
		}

		if f.RecordID > 0 {
			query = query.Where(squirrel.Eq{"v.record_id": f.RecordID})
		}

		if f.FieldName != "" {
			query = query.Where(squirrel.Eq{"v.name": f.FieldName})
		}
	default:
		err = errors.New("unsupported kind value")
		return
	}

	if f.Filter != "" {
		err = errors.New("filtering by filter not implemented")
		return
	}

	if f.Count, err = r.count(query); err != nil || f.Count == 0 {
		return
	}

	query = query.
		Columns(r.columns()...).
		OrderBy("id ASC")

	return set, f, r.fetchPaged(&set, query, f.Page, f.PerPage)
}

func (r attachment) Create(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r attachment) DeleteByID(namespaceID, attachmentID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		attachmentID,
	)

	return err
}
