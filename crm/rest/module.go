package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

type (
	Module struct {
		module service.ModuleService
		record service.RecordService
	}
)

func (Module) New() *Module {
	return &Module{
		module: service.DefaultModule,
		record: service.DefaultRecord,
	}
}

func (s *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	return s.module.With(ctx).Find()
}

func (s *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	return s.module.With(ctx).FindByID(r.ModuleID)
}

func (s *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	return resputil.OK(), s.module.With(ctx).DeleteByID(r.ModuleID)
}

func (s *Module) Create(ctx context.Context, r *request.ModuleCreate) (interface{}, error) {
	item := &types.Module{
		Name:   r.Name,
		Fields: r.Fields,
		Meta:   r.Meta,
	}
	return s.module.With(ctx).Create(item)
}

func (s *Module) Edit(ctx context.Context, r *request.ModuleEdit) (interface{}, error) {
	item := &types.Module{
		ID:     r.ModuleID,
		Name:   r.Name,
		Fields: r.Fields,
		Meta:   r.Meta,
	}
	return s.module.With(ctx).Update(item)
}

func (s *Module) RecordReport(ctx context.Context, r *request.ModuleRecordReport) (interface{}, error) {
	return s.record.With(ctx).Report(r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (s *Module) RecordList(ctx context.Context, r *request.ModuleRecordList) (interface{}, error) {
	return s.record.With(ctx).Find(r.ModuleID, r.Filter, r.Sort, r.Page, r.PerPage)
}

func (s *Module) RecordRead(ctx context.Context, r *request.ModuleRecordRead) (interface{}, error) {
	return s.record.With(ctx).FindByID(r.RecordID)
}

func (s *Module) RecordCreate(ctx context.Context, r *request.ModuleRecordCreate) (interface{}, error) {
	return s.record.With(ctx).Create(&types.Record{ModuleID: r.ModuleID, Values: r.Values})
}

func (s *Module) RecordEdit(ctx context.Context, r *request.ModuleRecordEdit) (interface{}, error) {
	return s.record.With(ctx).Update(&types.Record{ModuleID: r.ModuleID, Values: r.Values})
}

func (s *Module) RecordDelete(ctx context.Context, r *request.ModuleRecordDelete) (interface{}, error) {
	return resputil.OK(), s.record.With(ctx).DeleteByID(r.RecordID)
}
