package service

import (
	"context"
	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	module struct {
		repository moduleRepository
	}

	moduleRepository interface {
		FindByID(ctx context.Context, moduleID uint64) (*types.Module, error)
		Find(ctx context.Context) ([]*types.Module, error)

		Create(ctx context.Context, module *types.Module) (*types.Module, error)
		Update(ctx context.Context, module *types.Module) (*types.Module, error)
		DeleteByID(ctx context.Context, moduleID uint64) error
	}
)

func Module() module {
	return module{
		repository: repository.Module(),
	}
}

func (svc module) FindByID(ctx context.Context, id uint64) (*types.Module, error) {
	return svc.repository.FindByID(ctx, id)
}

func (svc module) Find(ctx context.Context) ([]*types.Module, error) {
	return svc.repository.Find(ctx)
}

func (svc module) Create(ctx context.Context, mod *types.Module) (*types.Module, error) {
	return svc.repository.Create(ctx, mod)
}

func (svc module) Update(ctx context.Context, mod *types.Module) (*types.Module, error) {
	return svc.repository.Update(ctx, mod)
}

func (svc module) DeleteByID(ctx context.Context, id uint64) error {
	return svc.repository.DeleteByID(ctx, id)
}
