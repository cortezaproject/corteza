package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	Attachment interface {
		With(ctx context.Context) Attachment

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
	sqlAttachmentScope = "deleted_at IS NULL"

	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

func NewAttachment(ctx context.Context) Attachment {
	return (&attachment{}).With(ctx)
}

func (r *attachment) With(ctx context.Context) Attachment {
	return &attachment{
		repository: r.repository.With(ctx),
	}
}

func (r *attachment) FindAttachmentByID(id uint64) (*types.Attachment, error) {
	sql := "SELECT * FROM attachments WHERE id = ? AND " + sqlAttachmentScope
	mod := &types.Attachment{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrAttachmentNotFound)
}

func (r *attachment) FindAttachmentByMessageID(IDs ...uint64) (rval types.MessageAttachmentSet, err error) {
	rval = make([]*types.MessageAttachment, 0)

	if len(IDs) == 0 {
		return
	}

	sql := `SELECT a.*, rel_message
		      FROM attachments AS a
		           INNER JOIN message_attachment AS ma ON a.id = ma.rel_attachment 
		     WHERE ma.rel_message IN (?) AND ` + sqlAttachmentScope

	if sql, args, err := sqlx.In(sql, IDs); err != nil {
		return nil, err
	} else {
		return rval, r.db().Select(&rval, sql, args...)
	}
}

func (r *attachment) CreateAttachment(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("attachments", mod)
}

func (r *attachment) DeleteAttachmentByID(id uint64) error {
	return r.updateColumnByID("attachments", "deleted_at", nil, id)
}

func (r *attachment) BindAttachment(attachmentId, messageId uint64) error {
	bond := struct {
		RelAttachment uint64 `db:"rel_attachment"`
		RelMessage    uint64 `db:"rel_message"`
	}{attachmentId, messageId}

	return r.db().Insert("message_attachment", bond)
}
