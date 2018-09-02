package rest

import (
	"github.com/pkg/errors"

	"github.com/titpetric/factory/resputil"

	"context"
	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type (
	Module struct {
		module  service.ModuleService
		content service.ContentService
	}
)

func (Module) New(module service.ModuleService, content service.ContentService) *Module {
	return &Module{module, content}
}

func (s *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	return s.module.With(ctx).Find()
}

func (s *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	return s.module.With(ctx).FindByID(r.ID)
}

func (s *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	return resputil.OK(), s.module.With(ctx).DeleteByID(r.ID)
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
		ID:     r.ID,
		Name:   r.Name,
		Fields: r.Fields,
	}
	return s.module.With(ctx).Update(item)
}

func (s *Module) ContentList(ctx context.Context, r *request.ModuleContentList) (interface{}, error) {
	return s.content.With(ctx).Find()
}

func (s *Module) ContentRead(ctx context.Context, r *request.ModuleContentRead) (interface{}, error) {
	return s.content.With(ctx).FindByID(r.ID)
}

func (s *Module) ContentCreate(ctx context.Context, r *request.ModuleContentCreate) (interface{}, error) {
	item := &types.Content{
		ModuleID: r.Module,
		Fields:   r.Fields,
	}
	return s.content.With(ctx).Create(item)
}

func (s *Module) ContentEdit(ctx context.Context, r *request.ModuleContentEdit) (interface{}, error) {
	item := &types.Content{
		ID:       r.ID,
		ModuleID: r.Module,
		Fields:   r.Fields,
	}
	return s.content.With(ctx).Update(item)
}

func (s *Module) ContentDelete(ctx context.Context, r *request.ModuleContentDelete) (interface{}, error) {
	return resputil.OK(), s.content.With(ctx).DeleteByID(r.ID)
}
