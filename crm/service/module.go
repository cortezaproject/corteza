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
		FindById(context.Context, uint64) (*types.Module, error)
		Find(context.Context) ([]*types.Module, error)

		Create(context.Context, *types.Module) (*types.Module, error)
		Update(context.Context, *types.Module) (*types.Module, error)
		DeleteById(context.Context, uint64) error
	}
)

func Module() module {
	return module{
		repository: repository.Module(),
	}
}

func (svc module) FindById(ctx context.Context, id uint64) (*types.Module, error) {
	return svc.repository.FindById(ctx, id)
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

func (svc module) DeleteById(ctx context.Context, id uint64) error {
	return svc.repository.DeleteById(ctx, id)
}
