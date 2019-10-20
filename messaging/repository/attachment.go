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
	AttachmentRepository interface {
		With(ctx context.Context, db *factory.DB) AttachmentRepository

		FindAttachmentByID(id uint64) (*types.Attachment, error)
		FindAttachmentByMessageID(IDs ...uint64) (types.MessageAttachmentSet, error)

		CreateAttachment(mod *types.Attachment) (*types.Attachment, error)
		DeleteAttachmentByID(id uint64) error

		BindAttachment(attachmentId, messageId uint64) error
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
	return "messaging_attachment"
}

func (r attachment) tableMessage() string {
	return "messaging_message_attachment"
}

func (r attachment) columns() []string {
	return []string{
		"a.id",
		"a.rel_user",
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

func (r attachment) FindAttachmentByID(ID uint64) (*types.Attachment, error) {
	return r.findOneBy("id", ID)
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

func (r attachment) FindAttachmentByMessageID(IDs ...uint64) (rval types.MessageAttachmentSet, err error) {
	rval = types.MessageAttachmentSet{}

	if len(IDs) == 0 {
		return
	}

	query := r.query().
		Columns("ma.rel_message").
		Join(r.tableMessage() + " AS ma ON (a.id = ma.rel_attachment)").
		Where(squirrel.Eq{"rel_message": IDs})

	return rval, rh.FetchAll(r.db(), query, &rval)
}

func (r attachment) CreateAttachment(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r attachment) DeleteAttachmentByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": time.Now()}, squirrel.Eq{"id": ID})
}

func (r attachment) BindAttachment(attachmentId, messageId uint64) error {
	bond := struct {
		RelAttachment uint64 `db:"rel_attachment"`
		RelMessage    uint64 `db:"rel_message"`
	}{attachmentId, messageId}

	return r.db().Insert(r.tableMessage(), bond)
}
