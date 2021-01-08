package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	files "github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180
)

type (
	attachment struct {
		actionlog actionlog.Recorder

		ac attachmentAccessController

		files files.Store
		store store.Storer
		event EventService
	}

	attachmentAccessController interface {
		CanAttachMessage(context.Context, *types.Channel) bool
	}

	AttachmentService interface {
		FindByID(ctx context.Context, id uint64) (*types.Attachment, error)
		CreateMessageAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, channelId, replyTo uint64) (*types.Attachment, error)
		OpenOriginal(ctx context.Context, att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(ctx context.Context, att *types.Attachment) (io.ReadSeeker, error)
	}
)

func Attachment(store files.Store) AttachmentService {
	return &attachment{
		ac:    DefaultAccessControl,
		files: store,
		store: DefaultStore,

		actionlog: DefaultActionlog,
		event:     DefaultEvent,
	}
}

func (svc attachment) FindByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return store.LookupMessagingAttachmentByID(ctx, svc.store, id)
}

func (svc attachment) OpenOriginal(_ context.Context, att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.Url) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.Url)
}

func (svc attachment) OpenPreview(_ context.Context, att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.PreviewUrl) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.PreviewUrl)
}

func (svc attachment) CreateMessageAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, channelID, replyTo uint64) (att *types.Attachment, err error) {
	var (
		aProps = &attachmentActionProps{channel: &types.Channel{ID: channelID}, replyTo: replyTo}

		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
		ch            *types.Channel
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {

		if ch, err = store.LookupMessagingChannelByID(ctx, s, channelID); err != nil {
			if errors.IsNotFound(err) {
				return AttachmentErrChannelNotFound()
			}
		}

		aProps.setChannel(ch)

		if !svc.ac.CanAttachMessage(ctx, ch) {
			return AttachmentErrNotAllowedToAttachToChannel()
		}

		att = &types.Attachment{
			OwnerID: currentUserID,
			Name:    strings.TrimSpace(name),
		}

		err = svc.create(ctx, name, size, fh, att)
		if err != nil {
			return err
		}

		if err = store.CreateMessagingAttachment(ctx, s, att); err != nil {
			return err
		}

		msg := &types.Message{
			ID:         nextID(),
			Attachment: att,
			Message:    name,
			Type:       types.MessageTypeAttachment,
			ChannelID:  channelID,
			ReplyTo:    replyTo,
			UserID:     currentUserID,
			CreatedAt:  *now(),
		}

		if strings.HasPrefix(att.Meta.Original.Mimetype, "image/") {
			msg.Type = types.MessageTypeInlineImage
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if err = store.CreateMessagingMessage(ctx, s, msg); err != nil {
			return
		}

		aProps.setMessageID(msg.ID)

		ma := &types.MessageAttachment{AttachmentID: att.ID, MessageID: msg.ID}
		if err = store.CreateMessagingMessageAttachment(ctx, s, ma); err != nil {
			return
		}

		return svc.sendEvent(ctx, msg)
	})

	return att, svc.recordAction(ctx, aProps, AttachmentActionCreate, err)
}

func (svc attachment) create(ctx context.Context, name string, size int64, fh io.ReadSeeker, att *types.Attachment) (err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	att.ID = nextID()
	att.CreatedAt = *now()

	if svc.files == nil {
		return fmt.Errorf("can not create attachment: store handler not set")
	}

	aProps.setName(name)
	aProps.setSize(size)

	// Extract extension but make sure path.Ext is not confused by any leading/trailing dots
	att.Meta.Original.Extension = strings.Trim(path.Ext(strings.Trim(name, ".")), ".")

	att.Meta.Original.Size = size
	if att.Meta.Original.Mimetype, err = svc.extractMimetype(fh); err != nil {
		return AttachmentErrFailedToExtractMimeType(aProps).Wrap(err)
	}

	att.Url = svc.files.Original(att.ID, att.Meta.Original.Extension)
	aProps.setUrl(att.Url)

	if err = svc.files.Save(att.Url, fh); err != nil {
		return AttachmentErrFailedToStoreFile(aProps).Wrap(err)
	}

	// Process image: extract width, height, make preview
	err = svc.processImage(fh, att)
	if err != nil {
		return AttachmentErrFailedToProcessImage(aProps).Wrap(err)
	}

	return nil
}

func (svc attachment) extractMimetype(file io.ReadSeeker) (mimetype string, err error) {
	if _, err = file.Seek(0, 0); err != nil {
		return
	}

	// Make sure we rewind when we're done
	defer file.Seek(0, 0)

	// See http.DetectContentType about 512 bytes
	var buf = make([]byte, 512)
	if _, err = file.Read(buf); err != nil {
		return
	}

	return http.DetectContentType(buf), nil
}

func (svc attachment) processImage(original io.ReadSeeker, att *types.Attachment) (err error) {
	if !strings.HasPrefix(att.Meta.Original.Mimetype, "image/") {
		// Only supporting previews from images (for now)
		return
	}

	var (
		preview       image.Image
		opts          []imaging.EncodeOption
		format        imaging.Format
		previewFormat imaging.Format
		animated      bool
		f2m           = map[imaging.Format]string{
			imaging.JPEG: "image/jpeg",
			imaging.GIF:  "image/gif",
		}

		f2e = map[imaging.Format]string{
			imaging.JPEG: "jpg",
			imaging.GIF:  "gif",
		}
	)

	if _, err = original.Seek(0, 0); err != nil {
		return
	}

	if format, err = imaging.FormatFromExtension(att.Meta.Original.Extension); err != nil {
		return fmt.Errorf("Could not get format from extension '%s': %w", att.Meta.Original.Extension, err)
	}

	previewFormat = format

	if imaging.JPEG == format {
		// Rotate image if needed
		// if preview, _, err = exiffix.Decode(original); err != nil {
		// 	//return errors.Wrapf(err, "Could not decode EXIF from JPEG")
		// }
		preview, _, _ = exiffix.Decode(original)
	}

	if imaging.GIF == format {
		// Decode all and check loops & delay to determine if GIF is animated or not
		if cfg, err := gif.DecodeAll(original); err == nil {
			animated = cfg.LoopCount > 0 || len(cfg.Delay) > 1

			// Use first image for the preview
			preview = cfg.Image[0]
		} else {
			return fmt.Errorf("could not decode gif config: %w", err)
		}

	} else {
		// Use GIF preview for GIFs and JPEG for everything else!
		previewFormat = imaging.JPEG

		// Store with a bit lower quality
		opts = append(opts, imaging.JPEGQuality(85))
	}

	// In case of JPEG we decode the image and rotate it beforehand
	// other cases are handled here
	if preview == nil {
		if preview, err = imaging.Decode(original); err != nil {
			return fmt.Errorf("could not decode original image: %w", err)
		}
	}

	var width, height = preview.Bounds().Max.X, preview.Bounds().Max.Y
	att.SetOriginalImageMeta(width, height, animated)

	if width > attachmentPreviewMaxWidth && width > height {
		// Landscape does not fit
		preview = imaging.Resize(preview, attachmentPreviewMaxWidth, 0, imaging.Lanczos)
	} else if height > attachmentPreviewMaxHeight {
		// Height does not fit
		preview = imaging.Resize(preview, 0, attachmentPreviewMaxHeight, imaging.Lanczos)
	}

	// Get dimensions from the preview
	width, height = preview.Bounds().Max.X, preview.Bounds().Max.Y

	var buf = &bytes.Buffer{}
	if err = imaging.Encode(buf, preview, previewFormat, opts...); err != nil {
		return
	}

	meta := att.SetPreviewImageMeta(width, height, false)
	meta.Size = int64(buf.Len())
	meta.Mimetype = f2m[previewFormat]
	meta.Extension = f2e[previewFormat]

	// Can and how we make a preview of this attachment?
	att.PreviewUrl = svc.files.Preview(att.ID, meta.Extension)

	return svc.files.Save(att.PreviewUrl, buf)
}

// Sends message to event loop
//
// It also preloads user
func (svc attachment) sendEvent(ctx context.Context, msg *types.Message) (err error) {
	return svc.event.Message(ctx, msg)
}
