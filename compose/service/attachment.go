package service

import (
	"bytes"
	"context"
	"image"
	"image/gif"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/store"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
	"github.com/gabriel-vasile/mimetype"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180

	// using base 10, it will be less confusing for the non-techie users
	megabyte = 1_000_000
)

var (
	reMimeType = regexp.MustCompile(`\w+/[-.\w]+(?:\+[-.\w]+)?`)
)

type (
	attachment struct {
		actionlog actionlog.Recorder
		objects   objstore.Store
		ac        attachmentAccessController
		store     store.Storer
	}

	attachmentAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateNamespace(context.Context) bool
		CanReadModule(context.Context, *types.Module) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanReadRecord(context.Context, *types.Record) bool
		CanUpdateRecord(context.Context, *types.Record) bool
		CanCreateRecordOnModule(context.Context, *types.Module) bool
	}

	AttachmentService interface {
		FindByID(ctx context.Context, namespaceID, attachmentID uint64) (*types.Attachment, error)
		Find(ctx context.Context, filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		CreatePageAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, pageID uint64) (*types.Attachment, error)
		CreateRecordAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64, fieldName string) (*types.Attachment, error)
		CreateNamespaceAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker) (*types.Attachment, error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
		DeleteByID(ctx context.Context, namespaceID, attachmentID uint64) error
	}
)

func Attachment(store objstore.Store) *attachment {
	return &attachment{
		objects: store,
		ac:      DefaultAccessControl,
		store:   DefaultStore,
	}
}

func (svc attachment) Find(ctx context.Context, filter types.AttachmentFilter) (set types.AttachmentSet, f types.AttachmentFilter, err error) {
	var (
		aProps = &attachmentActionProps{filter: &filter}
	)

	err = func() error {
		if filter.NamespaceID == 0 {
			return AttachmentErrInvalidNamespaceID()
		}

		if filter.PageID > 0 {
			aProps.namespace, aProps.page, err = loadPage(ctx, svc.store, filter.NamespaceID, filter.PageID)
			if err != nil {
				return err
			} else if svc.ac.CanReadPage(ctx, aProps.page) {
				return AttachmentErrNotAllowedToReadPage()
			}
		}

		if filter.RecordID > 0 {
			aProps.namespace, aProps.module, aProps.record, err = loadRecordCombo_old(ctx, svc.store, filter.NamespaceID, filter.ModuleID, filter.RecordID)
			if err != nil {
				return err
			} else if svc.ac.CanReadRecord(ctx, aProps.record) {
				return AttachmentErrNotAllowedToReadRecord()
			}
		} else if filter.ModuleID > 0 {
			aProps.namespace, aProps.module, err = loadModuleWithNamespace(ctx, svc.store, filter.NamespaceID, filter.ModuleID)
			if err != nil {
				return err
			} else if svc.ac.CanReadRecord(ctx, aProps.record) {
				return AttachmentErrNotAllowedToReadRecord()
			}
		}

		set, f, err = store.SearchComposeAttachments(ctx, svc.store, f)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, AttachmentActionSearch, err)
}

func (svc attachment) FindByID(ctx context.Context, namespaceID, attachmentID uint64) (att *types.Attachment, err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	err = func() error {
		if attachmentID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupComposeAttachmentByID(ctx, svc.store, attachmentID); err != nil {
			return err
		}

		aProps.setAttachment(att)
		return nil
	}()

	return att, svc.recordAction(ctx, aProps, AttachmentActionLookup, err)
}

func (svc attachment) DeleteByID(ctx context.Context, namespaceID, attachmentID uint64) (err error) {
	var (
		att    *types.Attachment
		aProps = &attachmentActionProps{attachment: &types.Attachment{ID: attachmentID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if attachmentID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = store.LookupComposeAttachmentByID(ctx, s, attachmentID); err != nil {
			return err
		}

		aProps.setAttachment(att)

		att.DeletedAt = now()
		return store.UpdateComposeAttachment(ctx, s, att)
	})

	return svc.recordAction(ctx, aProps, AttachmentActionDelete, err)
}

//func (svc attachment) findNamespaceByID(namespaceID uint64) (ns *types.Namespace, err error) {
//	if namespaceID == 0 {
//		return nil, AttachmentErrInvalidNamespaceID()
//	}
//
//	ns, err = svc.namespaceRepo.FindByID(namespaceID)
//	if repository.ErrNamespaceNotFound.Eq(err) {
//		return nil, AttachmentErrNamespaceNotFound()
//	}
//
//	if !svc.ac.CanReadNamespace(ctx, ns) {
//		return nil, AttachmentErrNotAllowedToReadNamespace()
//	}
//
//	return ns, nil
//}
//
//func (svc attachment) findPageByID(namespaceID, pageID uint64) (p *types.Page, err error) {
//	if pageID == 0 {
//		return nil, AttachmentErrInvalidPageID()
//	}
//
//	p, err = svc.pageRepo.FindByID(namespaceID, pageID)
//	if repository.ErrPageNotFound.Eq(err) {
//		return nil, AttachmentErrPageNotFound()
//	}
//
//	if !svc.ac.CanReadPage(ctx, p) {
//		return nil, AttachmentErrNotAllowedToReadPage()
//	}
//
//	return p, nil
//}
//
//func (svc attachment) findModuleByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
//	if moduleID == 0 {
//		return nil, AttachmentErrInvalidModuleID()
//	}
//
//	m, err = svc.moduleRepo.FindByID(namespaceID, moduleID)
//	if repository.ErrModuleNotFound.Eq(err) {
//		return nil, AttachmentErrModuleNotFound()
//	}
//
//	if !svc.ac.CanReadModule(ctx, m) {
//		return nil, AttachmentErrNotAllowedToReadModule()
//	}
//
//	return m, nil
//}
//
//func (svc attachment) findRecordByID(m *types.Module, recordID uint64) (r *types.Record, err error) {
//	if recordID == 0 {
//		return nil, AttachmentErrInvalidRecordID()
//	}
//
//	r, err = svc.recordRepo.FindByID(m.NamespaceID, recordID)
//	if repository.ErrRecordNotFound.Eq(err) {
//		return nil, AttachmentErrRecordNotFound()
//	}
//
//	if !svc.ac.CanReadRecord(ctx, m) {
//		return nil, AttachmentErrNotAllowedToReadRecord()
//	}
//
//	return r, nil
//}

func (svc attachment) OpenOriginal(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.Url) == 0 {
		return nil, nil
	}

	return svc.objects.Open(att.Url)
}

func (svc attachment) OpenPreview(att *types.Attachment) (io.ReadSeeker, error) {
	if len(att.PreviewUrl) == 0 {
		return nil, nil
	}

	return svc.objects.Open(att.PreviewUrl)
}

func (svc attachment) CreatePageAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, pageID uint64) (att *types.Attachment, err error) {
	var (
		ns *types.Namespace
		p  *types.Page

		aProps = &attachmentActionProps{
			namespace: &types.Namespace{ID: namespaceID},
			page:      &types.Page{ID: pageID},
		}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if size == 0 {
			return AttachmentErrNotAllowedToCreateEmptyAttachment()
		}

		ns, p, err = loadPage(ctx, s, namespaceID, pageID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setPage(p)

		if !svc.ac.CanUpdatePage(ctx, p) {
			return AttachmentErrNotAllowedToUpdatePage()
		}

		{
			// Verify size and type of the uploaded page attachment
			// Max size & allowed mime-types are pulled from the current settings
			var (
				maxSize      = int64(systemService.CurrentSettings.Compose.Page.Attachments.MaxSize) * megabyte
				allowedTypes = systemService.CurrentSettings.Compose.Page.Attachments.Mimetypes
				mimeType     string
			)

			if maxSize > 0 && maxSize < size {
				return AttachmentErrTooLarge().Apply(
					errors.Meta("size", size),
					errors.Meta("maxSize", maxSize),
				)
			}

			if mimeType, err = svc.extractMimetype(fh); err != nil {
				return err
			} else if !svc.checkMimeType(mimeType, allowedTypes...) {
				return AttachmentErrNotAllowedToUploadThisType()
			}
		}

		att = &types.Attachment{
			NamespaceID: namespaceID,
			Name:        strings.TrimSpace(name),
			Kind:        types.PageAttachment,
		}

		return svc.create(ctx, s, name, size, fh, att)
	})

	return att, svc.recordAction(ctx, aProps, AttachmentActionCreate, err)

}

func (svc attachment) CreateRecordAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64, fieldName string) (att *types.Attachment, err error) {
	var (
		ns *types.Namespace
		m  *types.Module
		r  *types.Record

		aProps = &attachmentActionProps{
			namespace: &types.Namespace{ID: namespaceID},
			module:    &types.Module{ID: moduleID},
			record:    &types.Record{ID: recordID},
		}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if size == 0 {
			return AttachmentErrNotAllowedToCreateEmptyAttachment()
		}

		ns, m, err = loadModuleWithNamespace(ctx, s, namespaceID, moduleID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

		if recordID > 0 {
			// Uploading to existing record
			//
			// To allow upload (attachment creation) user must have permissions to
			// alter that record

			r, err = store.LookupComposeRecordByID(ctx, s, m, recordID)
			if err != nil {
				return err
			}

			aProps.setRecord(r)

			if !svc.ac.CanUpdateRecord(ctx, r) {
				return AttachmentErrNotAllowedToUpdateRecord()
			}
		} else {
			// Uploading to non-existing record
			//
			// To allow upload (attachment creation) user must have permissions to
			// create records
			if !svc.ac.CanCreateRecordOnModule(ctx, m) {
				return AttachmentErrNotAllowedToCreateRecords()
			}
		}

		{
			// Verify size and type of the uploaded record attachment
			// Max size & allowed mime-types are pulled from the current settings
			var (
				maxSize      = int64(systemService.CurrentSettings.Compose.Record.Attachments.MaxSize) * megabyte
				allowedTypes = systemService.CurrentSettings.Compose.Record.Attachments.Mimetypes
				mimeType     string
			)

			f := m.Fields.FindByName(fieldName)
			if f == nil || f.Kind != "File" {
				return AttachmentErrInvalidModuleField().Apply(
					errors.Meta("fieldName", fieldName),
				)
			}

			if aux := f.Options.Int64("maxSize"); aux > 0 {
				maxSize = aux * megabyte
			}
			if aux := f.Options.String("mimetypes"); len(aux) > 0 {
				allowedTypes = strings.Split(aux, ",")
			}

			if maxSize > 0 && maxSize < size {
				return AttachmentErrTooLarge().Apply(
					errors.Meta("size", size),
					errors.Meta("maxSize", maxSize),
				)
			}

			if mimeType, err = svc.extractMimetype(fh); err != nil {
				return err
			} else if !svc.checkMimeType(mimeType, allowedTypes...) {
				return AttachmentErrNotAllowedToUploadThisType().Apply(errors.Meta("mimetype", mimeType))
			}
		}

		att = &types.Attachment{
			NamespaceID: namespaceID,
			Name:        strings.TrimSpace(name),
			Kind:        types.RecordAttachment,
		}

		return svc.create(ctx, s, name, size, fh, att)
	})

	return att, svc.recordAction(ctx, aProps, AttachmentActionCreate, err)
}

func (svc attachment) CreateNamespaceAttachment(ctx context.Context, name string, size int64, fh io.ReadSeeker) (att *types.Attachment, err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if size == 0 {
			return AttachmentErrNotAllowedToCreateEmptyAttachment()
		}

		if !svc.ac.CanCreateNamespace(ctx) {
			return AttachmentErrNotAllowedToUpdateNamespace()
		}

		{
			// Verify file size and
			var (
				// use max-file-size from page attachments for now
				maxSize  = int64(systemService.CurrentSettings.Compose.Page.Attachments.MaxSize) * megabyte
				mimeType string
			)

			if maxSize > 0 && maxSize < size {
				return AttachmentErrTooLarge().Apply(
					errors.Meta("size", size),
					errors.Meta("maxSize", maxSize),
				)
			}

			if mimeType, err = svc.extractMimetype(fh); err != nil {
				return err
			} else if !svc.checkMimeType(mimeType, "image/png", "image/gif", "image/jpeg") {
				return AttachmentErrNotAllowedToUploadThisType()
			}
		}

		att = &types.Attachment{
			Name: strings.TrimSpace(name),
			Kind: types.NamespaceAttachment,
		}

		// @todo limit upload on image/* only!

		return svc.create(ctx, s, name, size, fh, att)
	})

	return att, svc.recordAction(ctx, aProps, AttachmentActionCreate, err)
}

func (svc attachment) create(ctx context.Context, s store.ComposeAttachments, name string, size int64, fh io.ReadSeeker, att *types.Attachment) (err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	// preset attachment ID because we need ref for storage
	att.ID = nextID()
	att.CreatedAt = *now()

	if att.OwnerID == 0 {
		att.OwnerID = auth.GetIdentityFromContext(ctx).Identity()
	}

	if svc.objects == nil {
		return errors.Internal("cannot create attachment: store handler not set")
	}

	if size == 0 {
		return AttachmentErrNotAllowedToCreateEmptyAttachment()
	}

	aProps.setName(name)
	aProps.setSize(size)

	// Extract extension but make sure path.Ext is not confused by any leading/trailing dots
	att.Meta.Original.Extension = strings.Trim(path.Ext(strings.Trim(name, ".")), ".")

	att.Meta.Original.Size = size
	if att.Meta.Original.Mimetype, err = svc.extractMimetype(fh); err != nil {
		return AttachmentErrFailedToExtractMimeType(aProps).Wrap(err)
	}

	att.Url = svc.objects.Original(att.ID, att.Meta.Original.Extension)
	aProps.setUrl(att.Url)

	if err = svc.objects.Save(att.Url, fh); err != nil {
		return AttachmentErrFailedToStoreFile(aProps).Wrap(err)
	}

	// Process image: extract width, height, make preview
	err = svc.processImage(fh, att)
	if err != nil {
		return AttachmentErrFailedToProcessImage(aProps).Wrap(err)
	}

	if err = store.CreateComposeAttachment(ctx, s, att); err != nil {
		return
	}

	return nil
}

func (svc attachment) extractMimetype(file io.ReadSeeker) (mType string, err error) {
	if _, err = file.Seek(0, 0); err != nil {
		return
	}

	// Make sure we rewind when we're done
	defer func(file io.ReadSeeker, offset int64, whence int) {
		_, _ = file.Seek(offset, whence)
	}(file, 0, 0)

	var mime *mimetype.MIME
	if mime, err = mimetype.DetectReader(file); err != nil {
		return
	}

	return mime.String(), nil
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
		return errors.Internal("could not get format from extension '%s'", att.Meta.Original.Extension).Wrap(err)
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
			return errors.Internal("Could not decode gif config").Wrap(err)
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
			return errors.Internal("Could not decode original image").Wrap(err)
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
	att.PreviewUrl = svc.objects.Preview(att.ID, meta.Extension)

	return svc.objects.Save(att.PreviewUrl, buf)
}

func (attachment) checkMimeType(test string, vv ...string) bool {
	if len(vv) == 0 {
		// return true if there are no type constraints to check against
		return true
	}

	for _, v := range vv {
		v = strings.TrimSpace(v)

		if !reMimeType.MatchString(v) {
			continue
		}

		if v == test {
			return true
		}
	}

	return false
}

var _ AttachmentService = &attachment{}

// loadRecordCombo_old Loads namespace, module and record
// @todo temporary
func loadRecordCombo_old(ctx context.Context, s store.Storer, namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	if ns, m, err = loadModuleWithNamespace(ctx, s, namespaceID, moduleID); err != nil {
		return
	}

	if r, err = store.LookupComposeRecordByID(ctx, s, m, recordID); err != nil {
		return
	}

	if r.ModuleID != moduleID {
		return nil, nil, nil, RecordErrInvalidModuleID()
	}

	return
}
