package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Attachment interface {
		FindAttachmentByID(id uint64) (*types.Attachment, error)
		FindAttachmentByRange(channelID, fromAttachmentID, toAttachmentID uint64) ([]*types.Attachment, error)
		CreateAttachment(mod *types.Attachment) (*types.Attachment, error)
		UpdateAttachment(mod *types.Attachment) (*types.Attachment, error)
		DeleteAttachmentByID(id uint64) error
	}
)

const (
	sqlAttachmentScope = "deleted_at IS NULL"

	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

var _ Attachment = &repository{}

func (r *repository) FindAttachmentByID(id uint64) (*types.Attachment, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM attachments WHERE id = ? AND " + sqlAttachmentScope
	mod := &types.Attachment{}

	return mod, isFound(db.With(r.ctx).Get(mod, sql, id), mod.ID > 0, ErrAttachmentNotFound)
}

func (r *repository) FindAttachmentByRange(channelID, fromAttachmentID, toAttachmentID uint64) ([]*types.Attachment, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Attachment, 0)

	sql := `
		SELECT * 
          FROM attachments
         WHERE id BETWEEN ? AND ?
           AND rel_channel = ?
           AND deleted_at IS NULL`

	return rval, db.With(r.ctx).Select(&rval, sql, fromAttachmentID, toAttachmentID, channelID)
}

func (r *repository) CreateAttachment(mod *types.Attachment) (*types.Attachment, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	mod.Attachment = coalesceJson(mod.Attachment, []byte("{}"))

	return mod, factory.Database.MustGet().With(r.ctx).Insert("attachments", mod)
}

func (r *repository) UpdateAttachment(mod *types.Attachment) (*types.Attachment, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Attachment = coalesceJson(mod.Attachment, []byte("{}"))

	return mod, factory.Database.MustGet().With(r.ctx).Replace("attachments", mod)
}

func (r *repository) DeleteAttachmentByID(id uint64) error {
	return simpleDelete(r.ctx, "attachments", id)
}
