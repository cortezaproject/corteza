package service

import (
	"bytes"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/id"
	files "github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	"io"
	"net/http"
	"path"
	"strings"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180
)

type (
	attachment struct {
		ctx       context.Context
		actionlog actionlog.Recorder
		files     files.Store
		ac        attachmentAccessController
		store     store.Storer
	}

	attachmentAccessController interface {
		CanManageSettings(context.Context) bool
	}

	AttachmentService interface {
		With(ctx context.Context) AttachmentService

		FindByID(ID uint64) (*types.Attachment, error)
		Find(filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		CreateSettingsAttachment(name string, size int64, fh io.ReadSeeker, labels map[string]string) (*types.Attachment, error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
		DeleteByID(ID uint64) error
	}
)

func Attachment(store files.Store) AttachmentService {
	return (&attachment{
		files:     store,
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
		store:     DefaultNgStore,
	}).With(context.Background())
}

func (svc attachment) With(ctx context.Context) AttachmentService {
	return &attachment{
		ctx: ctx,

		actionlog: svc.actionlog,
		ac:        svc.ac,

		files: svc.files,

		store: DefaultNgStore,
	}
}

func (svc attachment) FindByID(ID uint64) (att *types.Attachment, err error) {
	var (
		aaProps = &attachmentActionProps{}
	)

	err = func() (err error) {
		if ID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupAttachmentByID(svc.ctx, svc.store, ID); err != nil {
			return err
		}

		aaProps.setAttachment(att)
		return nil
	}()

	return att, svc.recordAction(svc.ctx, aaProps, AttachmentActionLookup, err)
}

func (svc attachment) DeleteByID(ID uint64) (err error) {
	var (
		att     *types.Attachment
		aaProps = &attachmentActionProps{attachment: &types.Attachment{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupAttachmentByID(svc.ctx, svc.store, ID); err != nil {
			return err
		}

		att.DeletedAt = nowPtr()
		aaProps.setAttachment(att)

		return store.UpdateAttachment(svc.ctx, svc.store, att)
	}()

	return svc.recordAction(svc.ctx, aaProps, AttachmentActionDelete, err)
}

func (svc attachment) Find(filter types.AttachmentFilter) (aa types.AttachmentSet, f types.AttachmentFilter, err error) {
	var (
		aaProps = &attachmentActionProps{filter: &filter}
	)

	err = func() (err error) {
		aa, f, err = store.SearchAttachments(svc.ctx, svc.store, filter)
		return err
	}()

	return aa, f, svc.recordAction(svc.ctx, aaProps, AttachmentActionSearch, err)
}

func (svc attachment) OpenOriginal(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.Url) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.Url)
}

func (svc attachment) OpenPreview(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.PreviewUrl) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.PreviewUrl)
}

func (svc attachment) CreateSettingsAttachment(name string, size int64, fh io.ReadSeeker, labels map[string]string) (att *types.Attachment, err error) {
	var (
		aaProps       = &attachmentActionProps{}
		currentUserID = intAuth.GetIdentityFromContext(svc.ctx).Identity()
	)

	err = func() (err error) {
		if !svc.ac.CanManageSettings(svc.ctx) {
			return AttachmentErrNotAllowedToCreate()
		}

		att = &types.Attachment{
			ID:      id.Next(),
			OwnerID: currentUserID,
			Name:    strings.TrimSpace(name),
			Kind:    types.AttachmentKindSettings,
		}

		aaProps.setAttachment(att)

		if labels != nil {
			att.Meta.Labels = labels
		}

		if err = svc.create(name, size, fh, att); err != nil {
			return err
		}

		return err
	}()

	return att, svc.recordAction(svc.ctx, aaProps, AttachmentActionCreate, err)
}

func (svc attachment) create(name string, size int64, fh io.ReadSeeker, att *types.Attachment) (err error) {
	var (
		aaProps = &attachmentActionProps{}
	)

	if svc.files == nil {
		return errors.New("can not create attachment: store handler not set")
	}

	aaProps.setName(name)
	aaProps.setSize(size)

	// Extract extension but make sure path.Ext is not confused by any leading/trailing dots
	att.Meta.Original.Extension = strings.Trim(path.Ext(strings.Trim(name, ".")), ".")

	att.Meta.Original.Size = size
	if att.Meta.Original.Mimetype, err = svc.extractMimetype(fh); err != nil {
		return AttachmentErrFailedToExtractMimeType(aaProps).Wrap(err)
	}

	att.Url = svc.files.Original(att.ID, att.Meta.Original.Extension)
	aaProps.setUrl(att.Url)

	if err = svc.files.Save(att.Url, fh); err != nil {
		return AttachmentErrFailedToStoreFile(aaProps).Wrap(err)
	}

	// Process image: extract width, height, make preview
	err = svc.processImage(fh, att)
	if err != nil {
		return AttachmentErrFailedToProcessImage(aaProps).Wrap(err)
	}

	if err = store.CreateAttachment(svc.ctx, svc.store, att); err != nil {
		return
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
		return errors.Wrapf(err, "could not get format from extension '%s'", att.Meta.Original.Extension)
	}

	previewFormat = format

	if imaging.JPEG == format {
		// Rotate image if needed
		// if preview, _, err = exiffix.Decode(original); err != nil {
		// 	return errors.Wrapf(err, "Could not decode EXIF from JPEG")
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
			return errors.Wrapf(err, "could not decode gif config")
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
			return errors.Wrapf(err, "could not decode original image")
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
