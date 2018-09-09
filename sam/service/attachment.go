package service

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"io"
	"net/http"
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

func (svc attachment) extractMeta(att *types.Attachment, file io.ReadSeeker) (err error) {
	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	// Make sure we rewind...
	defer file.Seek(0, 0)

	// See http.DetectContentType about 512 bytes
	var buf = make([]byte, 512)
	if _, err = file.Read(buf); err != nil {
		return
	}

	att.Mimetype = http.DetectContentType(buf)

	// @todo compare mime with extension (or better, enforce extension from mimetype)
	//if extensions, err := mime.ExtensionsByType(att.Mimetype); err == nil {
	//	extensions[0]
	//}

	// @todo extract image info so we can provide additional features if needed
	//if strings.HasPrefix(att.Mimetype, "image/gif") {
	//	if cfg, err := gif.DecodeAll(file); err == nil {
	//		m.Width = cfg.Config.Width
	//		m.Height = cfg.Config.Height
	//		m.Animated = cfg.LoopCount > 0 || len(cfg.Delay) > 1
	//	}
	//} else if strings.HasPrefix(att.Mimetype, "image") {
	//	if cfg, _, err := image.DecodeConfig(file); err == nil {
	//		m.Width = cfg.Width
	//		m.Height = cfg.Height
	//	}
	//}

	return
}

func (svc attachment) makePreview(att *types.Attachment, original io.ReadSeeker) (err error) {
	if true {
		return
	}

	// Can and how we make a preview of this attachment?
	var ext = "jpg"
	att.PreviewUrl = svc.sto.Preview(att.ID, ext)

	return svc.sto.Save(att.PreviewUrl, original)
}

var _ AttachmentService = &attachment{}
