package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	"net/http"
	"net/url"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/store"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	attachment struct {
		ctx context.Context

		attachment repository.Attachment
		message    repository.Message
		store      store.Store

		config struct {
			url        string
			previewUrl string
		}
	}

	AttachmentService interface {
		With(ctx context.Context) AttachmentService

		FindByID(id uint64) (*types.Attachment, error)
		Create(channelId uint64, name string, size int64, fh io.ReadSeeker) (*types.Attachment, error)
		LoadFromMessages(mm types.MessageSet) (err error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
	}
)

func Attachment(store store.Store) *attachment {
	svc := &attachment{
		ctx:        context.Background(),
		attachment: repository.NewAttachment(context.Background()),
		message:    repository.NewMessage(context.Background()),
		store:      store,
	}
	svc.config.url = "/attachment/%d/%s"
	svc.config.previewUrl = "/attachment/%d/%s/preview"
	return svc
}

func (svc *attachment) With(ctx context.Context) AttachmentService {
	return &attachment{
		ctx:        ctx,
		attachment: svc.attachment.With(ctx),
		message:    svc.message.With(ctx),
		store:      svc.store,
		config:     svc.config,
	}
}

func (svc *attachment) FindByID(id uint64) (*types.Attachment, error) {
	return svc.attachment.FindAttachmentByID(id)
}

func (svc *attachment) OpenOriginal(att *types.Attachment) (io.ReadSeeker, error) {
	return svc.store.Open(att.Url)
}

func (svc *attachment) OpenPreview(att *types.Attachment) (io.ReadSeeker, error) {
	return svc.store.Open(att.PreviewUrl)

}

func (svc *attachment) LoadFromMessages(mm types.MessageSet) (err error) {
	var ids []uint64
	mm.Walk(func(m *types.Message) error {
		if m.Type == types.MessageTypeAttachment || m.Type == types.MessageTypeInlineImage {
			ids = append(ids, m.ID)
		}
		return nil
	})

	if set, err := svc.attachment.FindAttachmentByMessageID(ids...); err != nil {
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

func (svc *attachment) Create(channelId uint64, name string, size int64, fh io.ReadSeeker) (att *types.Attachment, err error) {
	var currentUserID uint64 = repository.Identity(svc.ctx)

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

	if err := svc.extractMeta(att, fh); err != nil {
		// @todo logmeta extraction failure
	}

	log.Printf("Processing uploaded file (name: %s, size: %d, mime: %s)", att.Name, att.Size, att.Mimetype)

	if svc.store != nil {
		att.Url = svc.store.Original(att.ID, ext)
		if err = svc.store.Save(att.Url, fh); err != nil {
			log.Print(err.Error())
			return
		}

		// Try to make preview
		svc.makePreview(att, fh)
	}

	log.Printf("File %s stored as %s", att.Name, att.Url)

	return att, repository.DB().Transaction(func() (err error) {

		if att, err = svc.attachment.CreateAttachment(att); err != nil {
			return
		}

		msg := &types.Message{
			Message:   name,
			Type:      types.MessageTypeAttachment,
			ChannelID: channelId,
			UserID:    currentUserID,
		}

		if strings.HasPrefix(att.Mimetype, "image/") {
			msg.Type = types.MessageTypeInlineImage
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if msg, err = svc.message.CreateMessage(msg); err != nil {
			return
		}

		if err = svc.attachment.BindAttachment(att.ID, msg.ID); err != nil {
			return
		}

		log.Printf("File %s (id: %d) attached to message (id: %d)", att.Name, att.ID, msg.ID)

		return
	})
}

// Generates URL to a location
func (svc *attachment) url(att *types.Attachment) string {
	return fmt.Sprintf(svc.config.url, att.ID, url.PathEscape(att.Name))
}

// Generates URL to a location
func (svc *attachment) previewUrl(att *types.Attachment) string {
	return fmt.Sprintf(svc.config.previewUrl, att.ID, url.PathEscape(att.Name))
}

func (svc *attachment) extractMeta(att *types.Attachment, file io.ReadSeeker) (err error) {
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

func (svc *attachment) makePreview(att *types.Attachment, original io.ReadSeeker) (err error) {
	if true {
		return
	}

	// Can and how we make a preview of this attachment?
	var ext = "jpg"
	att.PreviewUrl = svc.store.Preview(att.ID, ext)

	return svc.store.Save(att.PreviewUrl, original)
}

var _ AttachmentService = &attachment{}
