package service

import (
	"context"
	"fmt"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
	"github.com/crusttech/crust/store"
	"github.com/titpetric/factory"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	attachment struct {
		rpo attachmentRepository
		sto store.Store

		config struct {
			url        string
			previewUrl string
		}
	}

	AttachmentService interface {
		FindByID(id uint64) (*types.Attachment, error)
		Create(ctx context.Context, channelId uint64, name string, size int64, fh io.ReadSeeker) (*types.Attachment, error)
		LoadFromMessages(ctx context.Context, mm types.MessageSet) (err error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
	}

	attachmentRepository interface {
		repository.Transactionable
		repository.Attachment
	}
)

func Attachment(store store.Store) *attachment {
	svc := &attachment{}

	svc.config.url = "/attachment/%d/%s"
	svc.config.previewUrl = "/attachment/%d/%s/preview"
	svc.rpo = repository.New()
	svc.sto = store

	return svc
}

func (svc attachment) FindByID(id uint64) (*types.Attachment, error) {
	return svc.rpo.FindAttachmentByID(id)
}

func (svc attachment) OpenOriginal(att *types.Attachment) (io.ReadSeeker, error) {
	return svc.sto.Open(att.Url)
}

func (svc attachment) OpenPreview(att *types.Attachment) (io.ReadSeeker, error) {
	return svc.sto.Open(att.PreviewUrl)

}

func (svc attachment) LoadFromMessages(ctx context.Context, mm types.MessageSet) (err error) {
	var ids []uint64
	mm.Walk(func(m *types.Message) error {
		if m.Type == "attachment" {
			ids = append(ids, m.ID)
		}
		return nil
	})

	if set, err := svc.rpo.FindAttachmentByMessageID(ids...); err != nil {
		return err
	} else {
		return set.Walk(func(a *types.MessageAttachment) error {
			if a.MessageID > 0 {
				if m := mm.FindById(a.MessageID); m != nil {
					m.Attachment = &a.Attachment

					m.Attachment.Url = svc.url(&a.Attachment)
					m.Attachment.PreviewUrl = svc.previewUrl(&a.Attachment)
				}
			}

			return nil
		})
	}
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

	log.Printf("Processing uploaded file (name: %s, size: %d, mime: %s)", att.Name, att.Size, att.Mimetype)

	if svc.sto != nil {
		att.Url = svc.sto.Original(att.ID, ext)
		if err = svc.sto.Save(att.Url, fh); err != nil {
			log.Print(err.Error())
			return
		}

		// Try to make preview
		svc.makePreview(att, fh)
	}

	log.Printf("File %s stored as %s", att.Name, att.Url)

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

		log.Printf("File %s (id: %s) attached to message (id: %d)", att.Name, att.ID, msg.ID)

		return
	})
}

// Generates URL to a location
func (svc attachment) url(att *types.Attachment) string {
	return fmt.Sprintf(svc.config.url, att.ID, url.PathEscape(att.Name))
}

// Generates URL to a location
func (svc attachment) previewUrl(att *types.Attachment) string {
	return fmt.Sprintf(svc.config.previewUrl, att.ID, url.PathEscape(att.Name))
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
