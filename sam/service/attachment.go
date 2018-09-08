package service

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"io"
	"path"
	"strings"
)

type (
	attachment struct {
		rpo attachmentRepository
		sto attachmentStore
	}

	AttachmentService interface {
		Create(ctx context.Context, channelId uint64, name string, size int64, fh io.ReadSeeker) (*types.Attachment, error)
	}

	attachmentRepository interface {
		repository.Transactionable
		repository.Attachment
	}

	attachmentStore interface {
		Original(id uint64, ext string) string
		Preview(id uint64, ext string) string

		Save(filename string, contents io.Reader) error
		Remove(filename string) error
		Open(filename string) (io.Reader, error)
	}
)

func Attachment(store attachmentStore) *attachment {
	svc := &attachment{}
	svc.rpo = repository.New()
	svc.sto = store
	// @todo bind file store

	return svc
}

func (svc attachment) Create(ctx context.Context, channelId uint64, name string, size int64, fh io.ReadSeeker) (att *types.Attachment, err error) {
	var currentUserID uint64 = auth.GetIdentityFromContext(ctx).Identity()

	// @todo verify if current user can access this channel
	// @todo verify if current user can upload to this channel

	att = &types.Attachment{
		ID:       factory.Sonyflake.NextID(),
		UserID:   currentUserID,
		Name:     strings.TrimSpace(name),
		Mimetype: "application/octet-stream",
		Size:     size,
	}

	// Extract extension but make sure path.Ext is not confused by any leading/trailing dots
	var ext = strings.Trim(path.Ext(strings.Trim(name, ".")), ".")

	// @todo extract mimetype and update att.Mimetype

	if svc.sto != nil {
		att.Url = svc.sto.Original(att.ID, ext)
		if err = svc.sto.Save(att.Url, fh); err != nil {
			return
		}

		// Try to make preview
		svc.makePreview(att, fh)
	}

	return att, svc.rpo.BeginWith(ctx, func(r repository.Interfaces) (err error) {

		if att, err = r.CreateAttachment(att); err != nil {
			return
		}

		msg := &types.Message{
			Message:   name,
			Type:      "attachment",
			ChannelID: channelId,
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if msg, err = r.CreateMessage(msg); err != nil {
			return
		}

		if err = r.BindAttachment(att.ID, msg.ID); err != nil {
			return
		}

		return
	})
}

func (svc attachment) makePreview(att *types.Attachment, originalFh io.ReadSeeker) (err error) {
	if true {
		return
	}

	// Can and how we make a preview of this attachment?
	var ext = "jpg"
	att.PreviewUrl = svc.sto.Preview(att.ID, ext)

	return svc.sto.Save(att.PreviewUrl, originalFh)
}

var _ AttachmentService = &attachment{}
