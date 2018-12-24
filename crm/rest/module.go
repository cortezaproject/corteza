package rest

import (
	"context"
	"strings"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

type (
	Module struct {
		module  service.ModuleService
		content service.RecordService
	}
)

func (Module) New() *Module {
	return &Module{
		module:  service.DefaultModule,
		content: service.DefaultRecord,
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
	}
	return s.module.With(ctx).Create(item)
}

func (s *Module) Edit(ctx context.Context, r *request.ModuleEdit) (interface{}, error) {
	item := &types.Module{
		ID:     r.ModuleID,
		Name:   r.Name,
		Fields: r.Fields,
	}
	return s.module.With(ctx).Update(item)
}

func (s *Module) RecordReport(ctx context.Context, r *request.ModuleRecordReport) (interface{}, error) {
	reportParams := &types.RecordReport{}

	if strings.TrimSpace(r.Metrics) != "" {
		reportParams.ScanMetrics(strings.Split(r.Metrics, ",")...)
	}

	if strings.TrimSpace(r.Dimensions) != "" {
		reportParams.ScanDimensions(strings.Split(r.Dimensions, ",")...)
	}

	return s.content.With(ctx).Report(r.ModuleID, reportParams)
}

func (s *Module) RecordList(ctx context.Context, r *request.ModuleRecordList) (interface{}, error) {
	return s.content.With(ctx).Find(r.ModuleID, r.Query, r.Page, r.PerPage, r.Sort)
}

func (s *Module) RecordRead(ctx context.Context, r *request.ModuleRecordRead) (interface{}, error) {
	return s.content.With(ctx).FindByID(r.RecordID)
}

func (s *Module) RecordCreate(ctx context.Context, r *request.ModuleRecordCreate) (interface{}, error) {
	item := &types.Record{
		ModuleID: r.ModuleID,
		Fields:   r.Fields,
	}
	return s.content.With(ctx).Create(item)
}

func (s *Module) RecordEdit(ctx context.Context, r *request.ModuleRecordEdit) (interface{}, error) {
	item := &types.Record{
		ID:       r.RecordID,
		ModuleID: r.ModuleID,
		Fields:   r.Fields,
	}
	return s.content.With(ctx).Update(item)
}

func (s *Module) RecordDelete(ctx context.Context, r *request.ModuleRecordDelete) (interface{}, error) {
	return resputil.OK(), s.content.With(ctx).DeleteByID(r.RecordID)
}
