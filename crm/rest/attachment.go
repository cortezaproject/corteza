package rest

import (
	"context"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/internal/service"
	"github.com/crusttech/crust/crm/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Attachment struct {
	attachment service.AttachmentService
}

func (Attachment) New() *Attachment {
	return &Attachment{attachment: service.DefaultAttachment}
}

// Attachments returns list of all files attached to records
func (ctrl *Attachment) List(ctx context.Context, r *request.AttachmentList) (interface{}, error) {
	f := types.AttachmentFilter{
		Kind:      r.Kind,
		ModuleID:  r.ModuleID,
		RecordID:  r.RecordID,
		FieldName: r.FieldName,
		// Filter:    r.Filter,
		PerPage: r.PerPage,
		Page:    r.Page,
		// Sort:      r.Sort,
	}

	return makeRecordAttachmentSetPayload(ctrl.attachment.Find(f))
}

func (ctrl *Attachment) Details(ctx context.Context, r *request.AttachmentDetails) (interface{}, error) {
	if a, err := ctrl.attachment.FindByID(r.AttachmentID); err != nil {
		return nil, err
	} else {
		return makeAttachmentPayload(a), nil
	}
}

func (ctrl *Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	return loadAttachedFile(ctrl.attachment, r.AttachmentID, false, r.Download)
}

func (ctrl *Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	return loadAttachedFile(ctrl.attachment, r.AttachmentID, true, false)
}
