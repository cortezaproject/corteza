package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/pkg/errors"
	"io"
	"time"
)

var _ = errors.Wrap

type (
	Attachment struct {
		svc service.AttachmentService
	}

	tempAttachmentPayload struct {
		ThisIsATempSolutionForAttachmentPayload bool

		Name     string
		ModTime  time.Time
		Download bool
	}
)

func (Attachment) New(svc service.AttachmentService) *Attachment {
	return &Attachment{svc: svc}
}

func (ctrl *Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	rval := tempAttachmentPayload{Download: r.Download}

	if att, err := ctrl.svc.FindByID(r.AttachmentID); err != nil {
		return nil, err
	} else {
		rval.Name = att.Name
		rval.ModTime = att.CreatedAt
	}

	return rval, nil
}

func (ctrl *Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	return nil, errors.New("Not implemented: Attachment.preview")
}

func (ctrl Attachment) get(ID uint64, preview, download bool) (interface{}, error) {
	rval := tempAttachmentPayload{Download: download}

	if att, err := ctrl.svc.FindByID(ID); err != nil {
		return nil, err
	} else {
		// @todo update this to io.ReadSeeker when store's Open() func rval is updated
		var rs io.Reader

		if preview {
			rs, err = ctrl.svc.OpenPreview(att)
		} else {
			rs, err = ctrl.svc.OpenOriginal(att)
		}

		rval.Name = att.Name
		rval.ModTime = att.CreatedAt

		// @todo do something with rs :)
		_ = rs
	}

	return rval, nil
}
