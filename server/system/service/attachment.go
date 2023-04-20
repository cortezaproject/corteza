package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/assets"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	files "github.com/cortezaproject/corteza/server/pkg/objstore"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"image"
	"image/gif"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180
	avatarWidth                = 300
	avatarHeight               = 300
)

type (
	attachment struct {
		actionlog actionlog.Recorder
		files     files.Store
		ac        attachmentAccessController
		store     store.Storer
		opt       options.AttachmentOpt
		logger    *zap.Logger
	}

	attachmentAccessController interface {
		CanManageSettings(context.Context) bool
		CanCreateApplication(context.Context) bool
	}

	AttachmentService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Attachment, error)
		Find(ctx context.Context, filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		CreateSettingsAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (*types.Attachment, error)
		CreateApplicationAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (*types.Attachment, error)
		CreateAuthAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (*types.Attachment, error)
		CreateAvatarInitialsAttachment(ctx context.Context, initials string, bgColor string, textColor string) (att *types.Attachment, err error)
		OpenOriginal(att *types.Attachment) (io.ReadSeekCloser, error)
		OpenPreview(att *types.Attachment) (io.ReadSeekCloser, error)
		DeleteByID(ctx context.Context, ID uint64) error
	}
)

func Attachment(store files.Store, opt options.AttachmentOpt, log *zap.Logger) *attachment {
	return &attachment{
		files:     store,
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		opt:       opt,
		logger:    log.Named("attachment"),
	}
}

func (svc attachment) FindByID(ctx context.Context, ID uint64) (att *types.Attachment, err error) {
	var (
		aaProps = &attachmentActionProps{}
	)

	err = func() (err error) {
		if ID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupAttachmentByID(ctx, svc.store, ID); err != nil {
			return err
		}

		aaProps.setAttachment(att)
		return nil
	}()

	return att, svc.recordAction(ctx, aaProps, AttachmentActionLookup, err)
}

func (svc attachment) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		att     *types.Attachment
		aaProps = &attachmentActionProps{attachment: &types.Attachment{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupAttachmentByID(ctx, svc.store, ID); err != nil {
			return err
		}

		att.DeletedAt = now()
		aaProps.setAttachment(att)

		return store.UpdateAttachment(ctx, svc.store, att)
	}()

	return svc.recordAction(ctx, aaProps, AttachmentActionDelete, err)
}

func (svc attachment) Find(ctx context.Context, filter types.AttachmentFilter) (aa types.AttachmentSet, f types.AttachmentFilter, err error) {
	var (
		aaProps = &attachmentActionProps{filter: &filter}
	)

	err = func() (err error) {
		aa, f, err = store.SearchAttachments(ctx, svc.store, filter)
		return err
	}()

	return aa, f, svc.recordAction(ctx, aaProps, AttachmentActionSearch, err)
}

func (svc attachment) OpenOriginal(att *types.Attachment) (io.ReadSeekCloser, error) {
	if len(att.Url) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.Url)
}

func (svc attachment) OpenPreview(att *types.Attachment) (io.ReadSeekCloser, error) {
	if len(att.PreviewUrl) == 0 {
		return nil, nil
	}

	return svc.files.Open(att.PreviewUrl)
}

func (svc attachment) CreateSettingsAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (att *types.Attachment, err error) {
	var (
		aaProps       = &attachmentActionProps{}
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = func() (err error) {
		if !svc.ac.CanManageSettings(ctx) {
			return AttachmentErrNotAllowedToCreate()
		}

		att = &types.Attachment{
			OwnerID: currentUserID,
			Name:    strings.TrimSpace(name),
			Kind:    types.AttachmentKindSettings,
		}

		aaProps.setAttachment(att)

		if labels != nil {
			att.Meta.Labels = labels
		}

		if err = svc.create(ctx, name, size, fh, att); err != nil {
			return err
		}

		return err
	}()

	return att, svc.recordAction(ctx, aaProps, AttachmentActionCreate, err)
}

func (svc attachment) CreateApplicationAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (att *types.Attachment, err error) {
	var (
		aaProps       = &attachmentActionProps{}
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = func() (err error) {
		att = &types.Attachment{
			OwnerID: currentUserID,
			Name:    strings.TrimSpace(name),
			Kind:    types.AttachmentKindSettings,
		}

		aaProps.setAttachment(att)

		if !svc.ac.CanCreateApplication(ctx) {
			return AttachmentErrNotAllowedToCreate()
		}

		if labels != nil {
			att.Meta.Labels = labels
		}

		if err = svc.create(ctx, name, size, fh, att); err != nil {
			return err
		}

		return err
	}()

	return att, svc.recordAction(ctx, aaProps, AttachmentActionCreate, err)
}

func (svc attachment) CreateAuthAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker, labels map[string]string) (att *types.Attachment, err error) {
	var (
		aaProps       = &attachmentActionProps{}
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	// check the file type
	ext := filepath.Ext(name)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, AttachmentErrInvalidAvatarFileType()
	}

	//check the file size
	if size > svc.opt.AvatarMaxFileSize {
		return nil, AttachmentErrInvalidAvatarFileSize()
	}

	err = func() (err error) {
		att = &types.Attachment{
			OwnerID: currentUserID,
			Name:    strings.TrimSpace(name),
			Kind:    types.AttachmentKindAvatar,
		}

		aaProps.setAttachment(att)

		if labels != nil {
			att.Meta.Labels = labels
		}

		if err = svc.create(ctx, name, size, fh, att); err != nil {
			return err
		}

		return err
	}()

	return att, svc.recordAction(ctx, aaProps, AttachmentActionCreate, err)
}

func (svc attachment) CreateAvatarInitialsAttachment(ctx context.Context, initials string, bgColor string, textColor string) (att *types.Attachment, err error) {
	var (
		aaProps       = &attachmentActionProps{}
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	if bgColor == "" {
		bgColor = svc.opt.AvatarInitialsBackgroundColor
	}

	if textColor == "" {
		textColor = svc.opt.AvatarInitialsColor
	}

	// Set initials image background color and
	// draw the initials text in the center of the image
	dc := gg.NewContext(avatarWidth, avatarHeight)
	dc.SetHexColor(bgColor)
	dc.Clear()

	fontBytes, err := svc.processFontsFile()
	if err != nil {
		return nil, err
	}

	// Get the font face properties
	f, _ := truetype.Parse(fontBytes)

	face := truetype.NewFace(f, &truetype.Options{
		Size:    120,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// set initials text color
	dc.SetFontFace(face)
	dc.SetHexColor(textColor)

	textWidth, textHeight := dc.MeasureString(initials)
	dc.DrawString(initials, (avatarWidth-textWidth)/2, (avatarHeight+textHeight)/2.2)

	// Generate and save the initials image in PNG format
	err = func() (err error) {
		att = &types.Attachment{
			OwnerID: currentUserID,
			Name:    initials,
			Kind:    types.AttachmentKindAvatarInitials,
		}

		att.ID = nextID()
		att.CreatedAt = *now()
		att.Meta.Original.Extension = "png"
		att.Meta.Labels = map[string]string{"key": types.AttachmentKindAvatarInitials}

		if att.Meta.Original.Image == nil {
			att.Meta.Original.Image = &types.AttachmentImageMeta{}
		}

		att.Meta.Original.Image.Initial = initials
		att.Meta.Original.Image.InitialColor = textColor
		att.Meta.Original.Image.BackgroundColor = bgColor

		att.Url = svc.files.Original(att.ID, att.Meta.Original.Extension)
		aaProps.setUrl(att.Url)

		aaProps.setAttachment(att)

		var buf = &bytes.Buffer{}
		if err = imaging.Encode(buf, dc.Image(), imaging.PNG); err != nil {
			return err
		}

		if err = svc.files.Save(att.Url, buf); err != nil {
			return AttachmentErrFailedToStoreFile(aaProps).Wrap(err)
		}

		if err = store.CreateAttachment(ctx, svc.store, att); err != nil {
			return
		}

		return nil
	}()

	return att, nil
}

func (svc attachment) create(ctx context.Context, name string, size int64, fh io.ReadSeeker, att *types.Attachment) (err error) {
	var (
		avatar  image.Image
		aaProps = &attachmentActionProps{}
	)

	att.ID = nextID()
	att.CreatedAt = *now()

	if svc.files == nil {
		return fmt.Errorf("cannot create attachment: store handler not set")
	}

	if size == 0 {
		return AttachmentErrNotAllowedToCreateEmptyAttachment(aaProps)
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

	// process avatar
	if att.Kind == types.AttachmentKindAvatar {
		if avatar, err = imaging.Decode(fh); err != nil {
			return fmt.Errorf("could not decode original avatar: %w", err)
		}

		avatar = imaging.Resize(avatar, avatarWidth, avatarHeight, imaging.Lanczos)

		var buf = &bytes.Buffer{}
		if err = imaging.Encode(buf, avatar, imaging.JPEG); err != nil {
			return
		}

		if err = svc.files.Save(att.Url, buf); err != nil {
			return AttachmentErrFailedToStoreFile(aaProps).Wrap(err)
		}

		if err = store.CreateAttachment(ctx, svc.store, att); err != nil {
			return
		}

		return nil
	}

	if err = svc.files.Save(att.Url, fh); err != nil {
		return AttachmentErrFailedToStoreFile(aaProps).Wrap(err)
	}

	// Process image: extract width, height, make preview
	err = svc.processImage(fh, att)
	if err != nil {
		return AttachmentErrFailedToProcessImage(aaProps).Wrap(err)
	}

	if err = store.CreateAttachment(ctx, svc.store, att); err != nil {
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
	if !strings.HasPrefix(att.Meta.Original.Mimetype, "image/") || att.Meta.Original.Mimetype == "image/x-icon" {
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
		return fmt.Errorf("could not get format from extension '%s': %w", att.Meta.Original.Extension, err)
	}

	previewFormat = format

	if imaging.JPEG == format {
		// Rotate image if needed
		// if preview, _, err = exiffix.Decode(original); err != nil {
		// 	return fmt.Errorf("Could not decode EXIF from JPEG", err)
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

// processFontsFile checks if the file exists and has the correct file extension,
// It validates the file path provided in the AVATAR_INITIALS_FONT_PATH environment variable,
// then reads and returns the file content
func (svc attachment) processFontsFile() (fontBytes []byte, err error) {
	ext := strings.ToLower(filepath.Ext(svc.opt.AvatarInitialsFontPath))
	if ext != ".ttf" {
		err = fmt.Errorf("invalid font file extension, please provide a truetype font (.ttf) file")
		svc.logger.Error(err.Error())
		return nil, AttachmentErrInvalidAvatarGenerateFontFile().Wrap(err)
	}

	assetsFile := assets.Files(svc.logger, "")
	fontFile, err := assetsFile.Open(svc.opt.AvatarInitialsFontPath)
	if err != nil {
		err = fmt.Errorf("%w : please ensure that the correct AVATAR_INITIALS_FONT_PATH is set", err)
		svc.logger.Error(err.Error())
		return nil, AttachmentErrInvalidAvatarGenerateFontFile().Wrap(err)
	}

	fontBytes, err = io.ReadAll(fontFile)
	if err != nil {
		err = fmt.Errorf("failed to read font file: %w", err)
		svc.logger.Error(err.Error())
		return nil, AttachmentErrInvalidAvatarGenerateFontFile().Wrap(err)
	}

	return fontBytes, nil
}
