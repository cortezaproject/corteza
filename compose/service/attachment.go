package service

import (
	"bytes"
	"context"
	"image"
	"image/gif"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/edwvee/exiffix"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/store"
)

const (
	attachmentPreviewMaxWidth  = 320
	attachmentPreviewMaxHeight = 180
)

type (
	attachment struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		store store.Store

		ac attachmentAccessController

		pageRepo      repository.PageRepository
		moduleRepo    repository.ModuleRepository
		recordRepo    repository.RecordRepository
		namespaceRepo repository.NamespaceRepository

		attachmentRepo repository.AttachmentRepository
	}

	attachmentAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanReadRecord(context.Context, *types.Module) bool
		CanUpdateRecord(context.Context, *types.Module) bool
		CanCreateRecord(context.Context, *types.Module) bool
	}

	AttachmentService interface {
		With(ctx context.Context) AttachmentService

		FindByID(namespaceID, attachmentID uint64) (*types.Attachment, error)
		Find(filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		CreatePageAttachment(namespaceID uint64, name string, size int64, fh io.ReadSeeker, pageID uint64) (*types.Attachment, error)
		CreateRecordAttachment(namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64, fieldName string) (*types.Attachment, error)
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
		DeleteByID(namespaceID, attachmentID uint64) error
	}
)

func Attachment(store store.Store) AttachmentService {
	return (&attachment{
		store: store,
		ac:    DefaultAccessControl,
	}).With(context.Background())
}

func (svc attachment) With(ctx context.Context) AttachmentService {
	db := repository.DB(ctx)
	return &attachment{
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac:    svc.ac,
		store: svc.store,

		attachmentRepo: repository.Attachment(ctx, db),
		recordRepo:     repository.Record(ctx, db),
		pageRepo:       repository.Page(ctx, db),
		moduleRepo:     repository.Module(ctx, db),
		namespaceRepo:  repository.Namespace(ctx, db),
	}
}

func (svc attachment) FindByID(namespaceID, attachmentID uint64) (att *types.Attachment, err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	err = svc.db.Transaction(func() error {
		if attachmentID == 0 {
			return AttachmentErrInvalidID()
		}

		if att, err = svc.attachmentRepo.FindByID(namespaceID, attachmentID); err != nil {
			return err
		}

		aProps.setAttachment(att)
		return nil
	})

	return att, svc.recordAction(svc.ctx, aProps, AttachmentActionLookup, err)
}

func (svc attachment) DeleteByID(namespaceID, attachmentID uint64) (err error) {
	var (
		aProps = &attachmentActionProps{attachment: &types.Attachment{ID: attachmentID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if attachmentID == 0 {
			return AttachmentErrInvalidID()
		}

		if aProps.attachment, err = svc.attachmentRepo.FindByID(namespaceID, attachmentID); err != nil {
			return err
		}

		return svc.attachmentRepo.DeleteByID(namespaceID, attachmentID)
	})

	return svc.recordAction(svc.ctx, aProps, AttachmentActionDelete, err)
}

func (svc attachment) Find(filter types.AttachmentFilter) (set types.AttachmentSet, f types.AttachmentFilter, err error) {
	var (
		aProps = &attachmentActionProps{filter: &filter}
	)

	err = svc.db.Transaction(func() error {
		if filter.NamespaceID == 0 {
			return AttachmentErrInvalidNamespaceID()
		}

		if filter.PageID > 0 {
			if aProps.page, err = svc.findPageByID(filter.NamespaceID, filter.PageID); err != nil {
				return err
			}
		}

		if filter.ModuleID > 0 {
			var m *types.Module
			if m, err = svc.findModuleByID(filter.NamespaceID, filter.ModuleID); err != nil {
				return err
			}

			aProps.setModule(m)
			if filter.RecordID > 0 {
				if aProps.record, err = svc.findRecordByID(m, filter.RecordID); err != nil {
					return err
				}
			}
		} else if filter.RecordID > 0 {
			return AttachmentErrInvalidModuleID()
		}

		set, f, err = svc.attachmentRepo.Find(filter)
		return err
	})

	return set, f, svc.recordAction(svc.ctx, aProps, AttachmentActionSearch, err)
}

func (svc attachment) findNamespaceByID(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, AttachmentErrInvalidNamespaceID()
	}

	ns, err = svc.namespaceRepo.FindByID(namespaceID)
	if repository.ErrNamespaceNotFound.Eq(err) {
		return nil, AttachmentErrNamespaceNotFound()
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, AttachmentErrNotAllowedToReadNamespace()
	}

	return ns, nil
}

func (svc attachment) findPageByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	if pageID == 0 {
		return nil, AttachmentErrInvalidPageID()
	}

	p, err = svc.pageRepo.FindByID(namespaceID, pageID)
	if repository.ErrPageNotFound.Eq(err) {
		return nil, AttachmentErrPageNotFound()
	}

	if !svc.ac.CanReadPage(svc.ctx, p) {
		return nil, AttachmentErrNotAllowedToReadPage()
	}

	return p, nil
}

func (svc attachment) findModuleByID(namespaceID, moduleID uint64) (m *types.Module, err error) {
	if moduleID == 0 {
		return nil, AttachmentErrInvalidModuleID()
	}

	m, err = svc.moduleRepo.FindByID(namespaceID, moduleID)
	if repository.ErrModuleNotFound.Eq(err) {
		return nil, AttachmentErrModuleNotFound()
	}

	if !svc.ac.CanReadModule(svc.ctx, m) {
		return nil, AttachmentErrNotAllowedToReadModule()
	}

	return m, nil
}

func (svc attachment) findRecordByID(m *types.Module, recordID uint64) (r *types.Record, err error) {
	if recordID == 0 {
		return nil, AttachmentErrInvalidRecordID()
	}

	r, err = svc.recordRepo.FindByID(m.NamespaceID, recordID)
	if repository.ErrRecordNotFound.Eq(err) {
		return nil, AttachmentErrRecordNotFound()
	}

	if !svc.ac.CanReadRecord(svc.ctx, m) {
		return nil, AttachmentErrNotAllowedToReadRecord()
	}

	return r, nil
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

func (svc attachment) CreatePageAttachment(namespaceID uint64, name string, size int64, fh io.ReadSeeker, pageID uint64) (att *types.Attachment, err error) {
	var (
		ns *types.Namespace
		p  *types.Page

		aProps = &attachmentActionProps{
			namespace: &types.Namespace{ID: namespaceID},
			page:      &types.Page{ID: pageID},
		}
	)

	err = svc.db.Transaction(func() error {
		ns, err = svc.findNamespaceByID(namespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)

		p, err = svc.findPageByID(namespaceID, pageID)
		if err != nil {
			return err
		}

		aProps.setPage(p)

		if !svc.ac.CanUpdatePage(svc.ctx, p) {
			return AttachmentErrNotAllowedToUpdatePage()
		}

		att = &types.Attachment{
			NamespaceID: namespaceID,
			Name:        strings.TrimSpace(name),
			Kind:        types.PageAttachment,
		}

		return svc.create(name, size, fh, att)
	})

	return att, svc.recordAction(svc.ctx, aProps, AttachmentActionCreate, err)

}
func (svc attachment) CreateRecordAttachment(namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64, fieldName string) (att *types.Attachment, err error) {
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

	err = svc.db.Transaction(func() (err error) {
		ns, err = svc.findNamespaceByID(namespaceID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)

		m, err = svc.findModuleByID(namespaceID, moduleID)
		if err != nil {
			return err
		}

		aProps.setModule(m)

		if recordID > 0 {
			// Uploading to existing record
			//
			// To allow upload (attachment creation) user must have permissions to
			// alter that record

			r, err = svc.findRecordByID(m, recordID)
			if err != nil {
				return err
			}

			aProps.setRecord(r)

			if !svc.ac.CanUpdateRecord(svc.ctx, m) {
				return AttachmentErrNotAllowedToUpdateRecord()
			}
		} else {
			// Uploading to non-existing record
			//
			// To allow upload (attachment creation) user must have permissions to
			// create records
			if !svc.ac.CanCreateRecord(svc.ctx, m) {
				return AttachmentErrNotAllowedToCreateRecords()
			}
		}

		att = &types.Attachment{
			NamespaceID: namespaceID,
			Name:        strings.TrimSpace(name),
			Kind:        types.RecordAttachment,
		}

		return svc.create(name, size, fh, att)
	})

	return att, svc.recordAction(svc.ctx, aProps, AttachmentActionCreate, err)

}

func (svc attachment) create(name string, size int64, fh io.ReadSeeker, att *types.Attachment) (err error) {
	var (
		aProps = &attachmentActionProps{}
	)

	// preset attachment ID because we need ref for storage
	att.ID = factory.Sonyflake.NextID()

	if att.OwnerID == 0 {
		att.OwnerID = auth.GetIdentityFromContext(svc.ctx).Identity()
	}

	if svc.store == nil {
		return errors.New("can not create attachment: store handler not set")
	}

	aProps.setName(name)
	aProps.setSize(size)

	// Extract extension but make sure path.Ext is not confused by any leading/trailing dots
	att.Meta.Original.Extension = strings.Trim(path.Ext(strings.Trim(name, ".")), ".")

	att.Meta.Original.Size = size
	if att.Meta.Original.Mimetype, err = svc.extractMimetype(fh); err != nil {
		return AttachmentErrFailedToExtractMimeType(aProps).Wrap(err)
	}

	att.Url = svc.store.Original(att.ID, att.Meta.Original.Extension)
	aProps.setUrl(att.Url)

	if err = svc.store.Save(att.Url, fh); err != nil {
		return AttachmentErrFailedToStoreFile(aProps).Wrap(err)
	}

	// Process image: extract width, height, make preview
	err = svc.processImage(fh, att)
	if err != nil {
		return AttachmentErrFailedToProcessImage(aProps).Wrap(err)
	}

	if att, err = svc.attachmentRepo.Create(att); err != nil {
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
		return errors.Wrapf(err, "Could not get format from extension '%s'", att.Meta.Original.Extension)
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
			return errors.Wrapf(err, "Could not decode gif config")
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
			return errors.Wrapf(err, "Could not decode original image")
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
	att.PreviewUrl = svc.store.Preview(att.ID, meta.Extension)

	return svc.store.Save(att.PreviewUrl, buf)
}

var _ AttachmentService = &attachment{}
