package rest

import (
	"github.com/pkg/errors"

	"github.com/titpetric/factory/resputil"

	"context"
	"github.com/crusttech/crust/crm/rest/server"
	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/crm/service"
)

var _ = errors.Wrap

type (
	Module struct {
		module service.ModuleService
	}
)

func (Module) New() server.ModuleAPI {
	return &Module{
		module: service.Module(),
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
	m := &types.Module{Name: r.Name}
	return s.module.With(ctx).Create(m)
}

func (s *Module) Edit(ctx context.Context, r *server.ModuleEditRequest) (interface{}, error) {
	m := &types.Module{ID: r.ID, Name: r.Name}
	return s.module.With(ctx).Update(m)
}

func (*Module) ContentList(ctx context.Context, r *server.ModuleContentListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentRead(ctx context.Context, r *server.ModuleContentReadRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentCreate(ctx context.Context, r *server.ModuleContentCreateRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentEdit(ctx context.Context, r *server.ModuleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(ctx context.Context, r *server.ModuleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
