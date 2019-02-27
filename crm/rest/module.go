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

func (ctrl *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	return ctrl.module.With(ctx).Find()
}

func (ctrl *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	return ctrl.module.With(ctx).FindByID(r.ModuleID)
}

func (ctrl *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	return resputil.OK(), ctrl.module.With(ctx).DeleteByID(r.ModuleID)
}

func (ctrl *Module) Create(ctx context.Context, r *request.ModuleCreate) (interface{}, error) {
	item := &types.Module{
		Name:   r.Name,
		Fields: r.Fields,
		Meta:   r.Meta,
	}
	return ctrl.module.With(ctx).Create(item)
}

func (ctrl *Module) Update(ctx context.Context, r *request.ModuleUpdate) (interface{}, error) {
	item := &types.Module{
		ID:     r.ModuleID,
		Name:   r.Name,
		Fields: r.Fields,
		Meta:   r.Meta,
	}
	return ctrl.module.With(ctx).Update(item)
}
