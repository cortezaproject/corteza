package service

import (
	"context"
	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	module struct {
		repository repository.Module
	}

	ModuleService interface {
		With(ctx context.Context) ModuleService

		FindByID(moduleID uint64) (*types.Module, error)
		Find() ([]*types.Module, error)

		Create(module *types.Module) (*types.Module, error)
		Update(module *types.Module) (*types.Module, error)
		DeleteByID(moduleID uint64) error
	}
)

func Module() ModuleService {
	return &module{
		repository: repository.NewModule(context.Background()),
	}
}

func (s *module) With(ctx context.Context) ModuleService {
	return &module{
		repository: s.repository.With(ctx),
	}
}

func (s *module) FindByID(id uint64) (*types.Module, error) {
	return s.repository.FindByID(id)
}

func (s *module) Find() ([]*types.Module, error) {
	return s.repository.Find()
}

func (s *module) Create(mod *types.Module) (*types.Module, error) {
	return s.repository.Create(mod)
}

func (s *module) Update(mod *types.Module) (*types.Module, error) {
	return s.repository.Update(mod)
}

func (s *module) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
