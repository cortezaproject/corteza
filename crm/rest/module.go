package rest

import (
	"github.com/pkg/errors"

	"github.com/titpetric/factory/resputil"

	"context"
	"github.com/crusttech/crust/crm/rest/server"
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

func (Module) New() server.ModuleAPI {
	return &Module{
		module:  service.Module(),
		content: service.Content(),
	}
}

func (s *Module) List(ctx context.Context, r *server.ModuleListRequest) (interface{}, error) {
	return s.module.With(ctx).Find()
}

func (s *Module) Read(ctx context.Context, r *server.ModuleReadRequest) (interface{}, error) {
	return s.module.With(ctx).FindByID(r.ID)
}

func (s *Module) Delete(ctx context.Context, r *server.ModuleDeleteRequest) (interface{}, error) {
	return resputil.OK(), s.module.With(ctx).DeleteByID(r.ID)
}

func (s *Module) Create(ctx context.Context, r *server.ModuleCreateRequest) (interface{}, error) {
	return s.module.With(ctx).Create(
		&types.Module{Name: r.Name},
	)
}

func (s *Module) Edit(ctx context.Context, r *server.ModuleEditRequest) (interface{}, error) {
	return s.module.With(ctx).Update(
		&types.Module{ID: r.ID, Name: r.Name},
	)
}

func (s *Module) ContentList(ctx context.Context, r *server.ModuleContentListRequest) (interface{}, error) {
	return s.content.With(ctx).Find()
}

func (s *Module) ContentRead(ctx context.Context, r *server.ModuleContentReadRequest) (interface{}, error) {
	return s.content.With(ctx).FindByID(r.ID)
}

func (s *Module) ContentCreate(ctx context.Context, r *server.ModuleContentCreateRequest) (interface{}, error) {
	item := &types.Content{
		ModuleID: r.Module,
	}
	fields := &item.Fields
	if err := fields.Scan(r.Payload); err != nil {
		return nil, err
	}
	return s.content.With(ctx).Create(item)
}

func (s *Module) ContentEdit(ctx context.Context, r *server.ModuleContentEditRequest) (interface{}, error) {
	item := &types.Content{
		ID: r.ID,
		ModuleID: r.Module,
	}
	fields := &item.Fields
	if err := fields.Scan(r.Payload); err != nil {
		return nil, err
	}
	return s.content.With(ctx).Update(item)
}

func (s *Module) ContentDelete(ctx context.Context, r *server.ModuleContentDeleteRequest) (interface{}, error) {
	return resputil.OK(), s.content.With(ctx).DeleteByID(r.ID)
}
