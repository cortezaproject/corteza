package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlAttachmentScope = "deleted_at IS NULL"

	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

type (
	attachment struct{}
)

func Attachment() attachment {
	return attachment{}
}

func (r attachment) FindByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM attachments WHERE id = ? AND " + sqlAttachmentScope
	mod := &types.Attachment{}

	return mod, isFound(db.With(ctx).Get(mod, sql, id), mod.ID > 0, ErrAttachmentNotFound)
}

func (r attachment) FindByRange(ctx context.Context, channelID, fromAttachmentID, toAttachmentID uint64) ([]*types.Attachment, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Attachment, 0)

	sql := `
		SELECT * 
          FROM attachments
         WHERE id BETWEEN ? AND ?
           AND rel_channel = ?
           AND deleted_at IS NULL`

	return rval, db.With(ctx).Select(&rval, sql, fromAttachmentID, toAttachmentID, channelID)
}

func (r attachment) Create(ctx context.Context, mod *types.Attachment) (*types.Attachment, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	mod.Attachment = coalesceJson(mod.Attachment, []byte("{}"))

	return mod, factory.Database.MustGet().With(ctx).Insert("attachments", mod)
}

func (r attachment) Update(ctx context.Context, mod *types.Attachment) (*types.Attachment, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Attachment = coalesceJson(mod.Attachment, []byte("{}"))

	return mod, factory.Database.MustGet().With(ctx).Replace("attachments", mod)
}

func (r attachment) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "attachments", id)
}
