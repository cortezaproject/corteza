package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
)

var _ = errors.Wrap

type (
	recordPayload struct {
		*types.Record

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
	}

	recordSetPayload struct {
		Filter types.RecordFilter `json:"filter"`
		Set    []*recordPayload   `json:"set"`
	}

	Record struct {
		record     service.RecordService
		module     service.ModuleService
		attachment service.AttachmentService
		ac         recordAccessController
	}

	recordAccessController interface {
		CanUpdateRecord(context.Context, *types.Module) bool
		CanDeleteRecord(context.Context, *types.Module) bool
	}
)

func (Record) New() *Record {
	return &Record{
		record:     service.DefaultRecord,
		module:     service.DefaultModule,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl *Record) Report(ctx context.Context, r *request.RecordReport) (interface{}, error) {
	return ctrl.record.With(ctx).Report(r.NamespaceID, r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (ctrl *Record) List(ctx context.Context, r *request.RecordList) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	rr, filter, err := ctrl.record.With(ctx).Find(types.RecordFilter{
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Filter:      r.Filter,
		Sort:        r.Sort,
		PerPage:     r.PerPage,
		Page:        r.Page,
	})

	return ctrl.makeFilterPayload(ctx, m, rr, filter, err)
}

func (ctrl *Record) Read(ctx context.Context, r *request.RecordRead) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.With(ctx).FindByID(r.NamespaceID, r.RecordID)

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Create(ctx context.Context, r *request.RecordCreate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.With(ctx).Create(&types.Record{
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Values:      r.Values,
	})

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Update(ctx context.Context, r *request.RecordUpdate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.With(ctx).Update(&types.Record{
		ID:          r.RecordID,
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Values:      r.Values,
	})

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return resputil.OK(), ctrl.record.With(ctx).DeleteByID(r.NamespaceID, r.RecordID)
}

func (ctrl *Record) Upload(ctx context.Context, r *request.RecordUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.With(ctx).CreateRecordAttachment(
		r.NamespaceID,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.ModuleID,
		r.RecordID,
		r.FieldName,
	)

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Record) makePayload(ctx context.Context, m *types.Module, r *types.Record, err error) (*recordPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &recordPayload{
		Record: r,

		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, m),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, m),
	}, nil
}

func (ctrl Record) makeFilterPayload(ctx context.Context, m *types.Module, rr types.RecordSet, f types.RecordFilter, err error) (*recordSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &recordSetPayload{Filter: f, Set: make([]*recordPayload, len(rr))}

	for i := range rr {
		modp.Set[i], _ = ctrl.makePayload(ctx, m, rr[i], nil)
	}

	return modp, nil
}
