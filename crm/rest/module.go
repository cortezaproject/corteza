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

func (c *Module) List(r *server.ModuleListRequest) (interface{}, error) {
	return c.service.Find(context.TODO())
}

func (c *Module) Edit(r *server.ModuleEditRequest) (interface{}, error) {
	m := types.Module{}.New()
	m.SetID(r.ID).SetName(r.Name)

	if m.GetID() > 0 {
		return c.service.Update(context.TODO(), m)
	}

	return c.service.Create(context.TODO(), m)
}

func (*Module) ContentList(r *server.ModuleContentListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentEdit(r *server.ModuleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(r *server.ModuleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
