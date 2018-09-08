package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Attachment interface {
		FindAttachmentByID(id uint64) (*types.Attachment, error)
		CreateAttachment(mod *types.Attachment) (*types.Attachment, error)
		DeleteAttachmentByID(id uint64) error
		BindAttachment(attachmentId, messageId uint64) error
	}
)

const (
	sqlAttachmentScope = "deleted_at IS NULL"

	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

var _ Attachment = &repository{}

func (r *repository) FindAttachmentByID(id uint64) (*types.Attachment, error) {
	sql := "SELECT * FROM attachments WHERE id = ? AND " + sqlAttachmentScope
	mod := &types.Attachment{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrAttachmentNotFound)
}

func (r *repository) FindAttachmentByMessageID(id uint64) (*types.Attachment, error) {
	sql := "SELECT a.* " +
		"     FROM attachments AS a" +
		"          INNER JOIN message_attachment AS ma ON a.id = ma.rel_attachment " +
		"    WHERE ma.rel_message = ? AND " + sqlAttachmentScope
	mod := &types.Attachment{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrAttachmentNotFound)
}

func (r *repository) FindAttachmentByRange(channelID, fromAttachmentID, toAttachmentID uint64) ([]*types.Attachment, error) {
	rval := make([]*types.Attachment, 0)

	sql := `
		SELECT * 
          FROM attachments
         WHERE id BETWEEN ? AND ?
           AND rel_channel = ?
           AND deleted_at IS NULL`

	return rval, r.db().Select(&rval, sql, fromAttachmentID, toAttachmentID, channelID)
}

func (r *repository) CreateAttachment(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("attachments", mod)
}

func (r *repository) DeleteAttachmentByID(id uint64) error {
	return r.updateColumnByID("attachments", "deleted_at", nil, id)
}

func (r *repository) BindAttachment(attachmentId, messageId uint64) error {
	bond := struct {
		RelAttachment uint64 `db:"rel_attachment"`
		RelMessage    uint64 `db:"rel_message"`
	}{attachmentId, messageId}

	return r.db().Insert("message_attachment", bond)
}
