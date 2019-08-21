package service

import (
	"bytes"
	"context"
	"image"
	"image/gif"
	"io"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	files "github.com/cortezaproject/corteza-server/pkg"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180
)

var (
	f2m = map[imaging.Format]string{
		imaging.JPEG: "image/jpeg",
		imaging.GIF:  "image/gif",
	}

	f2e = map[imaging.Format]string{
		imaging.JPEG: "jpg",
		imaging.GIF:  "gif",
	}
)

type (
	attachment struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac attachmentAccessController

		store   store.Store
		event   EventService
		channel ChannelService

		attachment repository.AttachmentRepository
		message    repository.MessageRepository
	}

	attachmentAccessController interface {
		CanAttachMessage(context.Context, *types.Channel) bool
	}

	AttachmentService interface {
		With(ctx context.Context) AttachmentService

		FindByID(id uint64) (*types.Attachment, error)
		Create(channelId, replyTo uint64, name string, fh *bytes.Reader, size int64, pname string, pfh *bytes.Reader, psize int64) (*types.Attachment, error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
	}
)

func Attachment(ctx context.Context, store store.Store) AttachmentService {
	return (&attachment{
		logger:  DefaultLogger.Named("attachment"),
		ac:      DefaultAccessControl,
		channel: DefaultChannel,
		store:   store,
	}).With(ctx)
}

func (svc attachment) With(ctx context.Context) AttachmentService {
	db := repository.DB(ctx)
	return &attachment{
		ctx:    ctx,
		db:     db,
		ac:     svc.ac,
		logger: svc.logger,

		store:   svc.store,
		event:   Event(ctx),
		channel: svc.channel.With(ctx),

		attachment: repository.Attachment(ctx, db),
		message:    repository.Message(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc attachment) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

func (svc attachment) FindByID(id uint64) (*types.Attachment, error) {
	return svc.attachment.FindAttachmentByID(id)
}

func (svc attachment) OpenOriginal(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.Url) == 0 {
		return nil, nil
	}

	return svc.store.Open(att.Url)
}

func (svc attachment) OpenPreview(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.PreviewUrl) == 0 {
		return nil, nil
	}

	return svc.store.Open(att.PreviewUrl)
}

func isAnimated(g *gif.GIF) bool {
	return g != nil && (g.LoopCount > 0 || len(g.Delay) > 1)
}

func (svc attachment) Create(channelId, replyTo uint64, name string, fh *bytes.Reader, size int64, pname string, pfh *bytes.Reader, psize int64) (att *types.Attachment, err error) {
	if svc.store == nil {
		return nil, errors.New("Can not create attachment: store handler not set")
	}

	currentUserID := repository.Identity(svc.ctx)

	if ch, err := svc.channel.FindByID(channelId); err != nil {
		return nil, err
	} else if !svc.ac.CanAttachMessage(svc.ctx, ch) {
		return nil, ErrNoPermissions.withStack()
	}

	att = &types.Attachment{
		ID:     factory.Sonyflake.NextID(),
		UserID: currentUserID,
		Name:   strings.TrimSpace(name),
	}

	log := svc.log(
		zap.String("name", att.Name),
		zap.Int64("size", att.Meta.Original.Size),
	)

	att.Meta.Original.Extension, _ = files.ExtractExtFromURL(name)
	att.Meta.Original.Size = size
	if att.Meta.Original.Mimetype, err = svc.extractMimetype(fh); err != nil {
		log.Error("could not extract mime-type", zap.Error(err))
		return
	}

	if att.Meta.Preview == nil {
		// initial set in case preview meta is not defined yet
		att.SetPreviewImageMeta(0, 0, false)
	}

	if pfh != nil {
		att.Meta.Preview.Extension, _ = files.ExtractExtFromURL(pname)
		att.Meta.Preview.Size = psize
		if att.Meta.Preview.Mimetype, err = svc.extractMimetype(pfh); err != nil {
			log.Error("could not extract mime-type", zap.Error(err))
			return
		}
	} else {
		att.Meta.Preview.Extension = att.Meta.Original.Extension
		att.Meta.Preview.Size = att.Meta.Original.Size
		att.Meta.Preview.Mimetype = att.Meta.Original.Mimetype
	}

	att.Url = svc.store.Original(att.ID, att.Meta.Original.Extension)
	if err = svc.store.Save(att.Url, fh); err != nil {
		log.Error("could not store file", zap.Error(err))
		return
	}

	// Only support processing for images
	if strings.HasPrefix(att.Meta.Original.Mimetype, "image/") {
		// Original
		format, err := imaging.FormatFromExtension(att.Meta.Original.Extension)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get format from extension '%s'", att.Meta.Original.Extension)
		}
		img, g, err := svc.loadMedia(fh, format)
		if err != nil {
			return nil, err
		}
		animated := isAnimated(g)
		att.Meta.Original.SetImageMeta(img.Bounds().Max.X, img.Bounds().Max.Y, animated)

		// Preview
		format, err = imaging.FormatFromExtension(att.Meta.Preview.Extension)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get format from extension '%s'", att.Meta.Preview.Extension)
		}
		if pfh != nil {
			img, g, err := svc.loadMedia(pfh, format)
			if err != nil {
				return nil, err
			}
			animated := isAnimated(g)
			att.Meta.Preview.SetImageMeta(img.Bounds().Max.X, img.Bounds().Max.Y, animated)
			att.PreviewUrl = svc.store.Preview(att.ID, att.Meta.Preview.Extension)
			svc.store.Save(att.PreviewUrl, pfh)
		} else {
			bb, width, height, size, t, ext, err := svc.makePreview(img, g, format)
			if err != nil {
				return nil, err
			}
			att.Meta.Preview.Size = size
			att.Meta.Preview.Mimetype = t
			att.Meta.Preview.Extension = ext
			att.Meta.Preview.SetImageMeta(width, height, false)
			att.PreviewUrl = svc.store.Preview(att.ID, ext)
			svc.store.Save(att.PreviewUrl, bb)
		}
	}

	return att, svc.db.Transaction(func() (err error) {
		if att, err = svc.attachment.CreateAttachment(att); err != nil {
			return
		}

		msg := &types.Message{
			Attachment: att,
			Message:    name,
			Type:       types.MessageTypeAttachment,
			ChannelID:  channelId,
			ReplyTo:    replyTo,
			UserID:     currentUserID,
		}

		if strings.HasPrefix(att.Meta.Original.Mimetype, "image/") {
			msg.Type = types.MessageTypeInlineImage
		}

		// Create the first message, doing this directly with repository to circumvent
		// message service constraints
		if msg, err = svc.message.Create(msg); err != nil {
			return
		}

		if err = svc.attachment.BindAttachment(att.ID, msg.ID); err != nil {
			return
		}

		return svc.sendEvent(msg)
	})
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

func (svc attachment) loadMedia(fh io.ReadSeeker, format imaging.Format) (img image.Image, g *gif.GIF, err error) {
	defer fh.Seek(0, 0)
	if _, err = fh.Seek(0, 0); err != nil {
		return
	}

	if imaging.JPEG == format {
		if img, _, err = exiffix.Decode(fh); err != nil {
			return
		}
	}

	if imaging.GIF == format {
		// Decode all and check loops & delay to determine if GIF is animated or not
		if g, err = gif.DecodeAll(fh); err == nil {
			img = g.Image[0]
		} else {
			return
		}
	}

	if img == nil {
		if img, _, err = exiffix.Decode(fh); err != nil {
			return
		}
	}

	return
}

func (svc attachment) makePreview(img image.Image, g *gif.GIF, format imaging.Format) (buf *bytes.Buffer, width, height int, size int64, mimetype, extension string, err error) {
	var (
		processedFormat imaging.Format
		opts            []imaging.EncodeOption
	)

	processedFormat = format
	if format != imaging.GIF {
		processedFormat = imaging.JPEG
		opts = append(opts, imaging.JPEGQuality(85))
	}

	width, height = img.Bounds().Max.X, img.Bounds().Max.Y
	if width > attachmentPreviewMaxWidth && width > height {
		// Landscape does not fit
		img = imaging.Resize(img, attachmentPreviewMaxWidth, 0, imaging.Lanczos)
	} else if height > attachmentPreviewMaxHeight {
		// Height does not fit
		img = imaging.Resize(img, 0, attachmentPreviewMaxHeight, imaging.Lanczos)
	}
	width, height = img.Bounds().Max.X, img.Bounds().Max.Y
	buf = &bytes.Buffer{}
	if err = imaging.Encode(buf, img, processedFormat, opts...); err != nil {
		return
	}
	size = int64(buf.Len())
	mimetype = f2m[processedFormat]
	extension = f2e[processedFormat]
	return
}

// Sends message to event loop
//
// It also preloads user
func (svc attachment) sendEvent(msg *types.Message) (err error) {
	return svc.event.Message(msg)
}
