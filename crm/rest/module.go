package rest

import (
	"github.com/pkg/errors"

	"context"
	"github.com/crusttech/crust/crm/rest/server"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type (
	Module struct {
		service moduleService
	}

	moduleService interface {
		FindById(context.Context, uint64) (*types.Module, error)
		Find(context.Context) ([]*types.Module, error)

		Create(context.Context, *types.Module) (*types.Module, error)
		Update(context.Context, *types.Module) (*types.Module, error)
		Delete(context.Context, *types.Module) error
	}
)

func (Module) New() *Module {
	return &Module{
		service: service.Module(),
	}
}

func (c *Module) List(ctx context.Context, r *server.ModuleListRequest) (interface{}, error) {
	return c.service.Find(ctx)
}

func (c *Module) Edit(ctx context.Context, r *server.ModuleEditRequest) (interface{}, error) {
	m := types.Module{}.New()
	m.SetID(r.ID).SetName(r.Name)

	if m.GetID() > 0 {
		return c.service.Update(ctx, m)
	}

	return c.service.Create(ctx, m)
}

func (*Module) ContentList(ctx context.Context, r *server.ModuleContentListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentEdit(ctx context.Context, r *server.ModuleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(ctx context.Context, r *server.ModuleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
