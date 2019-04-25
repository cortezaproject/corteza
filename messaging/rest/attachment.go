package rest

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/handlers"
	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/types"
)

var _ = errors.Wrap

type (
	Attachment struct {
		att service.AttachmentService
	}

	file struct {
		*types.Attachment
		content  io.ReadSeeker
		download bool
	}
)

func (Attachment) New() *Attachment {
	ctrl := &Attachment{}
	ctrl.att = service.DefaultAttachment
	return ctrl
}

func (ctrl *Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	if err := ctrl.isAccessible(r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.get(ctx, r.AttachmentID, false, r.Download)
}

func (ctrl *Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	if err := ctrl.isAccessible(r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.get(ctx, r.AttachmentID, true, false)
}

func (ctrl Attachment) isAccessible(attachmentID, userID uint64, signature string) error {
	if userID == 0 {
		return errors.New("missing or invalid user ID")
	}

	if attachmentID == 0 {
		return errors.New("missing or invalid attachment ID")
	}

	if auth.DefaultSigner.Verify(signature, userID, attachmentID) {
		return errors.New("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) get(ctx context.Context, ID uint64, preview, download bool) (handlers.Downloadable, error) {
	rval := &file{download: download}

	if att, err := ctrl.att.With(ctx).FindByID(ID); err != nil {
		return nil, err
	} else {
		rval.Attachment = att
		if preview {
			rval.content, err = ctrl.att.OpenPreview(att)
		} else {
			rval.content, err = ctrl.att.OpenOriginal(att)
		}

		if err != nil {
			return nil, err
		}
	}

	return rval, nil
}

func (f *file) Download() bool {
	return f.download
}

func (f *file) Name() string {
	return f.Attachment.Name
}

func (f *file) ModTime() time.Time {
	return f.Attachment.CreatedAt
}

func (f *file) Content() io.ReadSeeker {
	return f.content
}

func (f *file) Valid() bool {
	return f.content != nil
}
