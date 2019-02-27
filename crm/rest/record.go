package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Record struct {
	record     service.RecordService
	attachment service.AttachmentService
}

func (Record) New() *Record {
	return &Record{
		record:     service.DefaultRecord,
		attachment: service.DefaultAttachment,
	}
}

func (ctrl *Record) Report(ctx context.Context, r *request.RecordReport) (interface{}, error) {
	return ctrl.record.With(ctx).Report(r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (ctrl *Record) List(ctx context.Context, r *request.RecordList) (interface{}, error) {
	return ctrl.record.With(ctx).Find(r.ModuleID, r.Filter, r.Sort, r.Page, r.PerPage)
}

func (ctrl *Record) Read(ctx context.Context, r *request.RecordRead) (interface{}, error) {
	return ctrl.record.With(ctx).FindByID(r.RecordID)
}

func (ctrl *Record) Create(ctx context.Context, r *request.RecordCreate) (interface{}, error) {
	return ctrl.record.With(ctx).Create(&types.Record{ModuleID: r.ModuleID, Values: r.Values})
}

func (ctrl *Record) Update(ctx context.Context, r *request.RecordUpdate) (interface{}, error) {
	return ctrl.record.With(ctx).Update(&types.Record{
		ID:       r.RecordID,
		ModuleID: r.ModuleID,
		Values:   r.Values})
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return resputil.OK(), ctrl.record.With(ctx).DeleteByID(r.RecordID)
}

func (ctrl *Record) Upload(ctx context.Context, r *request.RecordUpload) (interface{}, error) {
	// @todo [SECURITY] check if attachments can be added to this page
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.With(ctx).CreateRecordAttachment(
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.ModuleID,
		r.RecordID,
		r.FieldName,
	)

	if err != nil {
		return nil, err
	}

	return makeAttachmentPayload(a), nil
}
