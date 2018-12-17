package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	module struct {
		db  *factory.DB
		ctx context.Context

		moduleRepo repository.ModuleRepository
		pageRepo   repository.PageRepository
	}

	ModuleService interface {
		With(ctx context.Context) ModuleService

		FindByID(moduleID uint64) (*types.Module, error)
		Find() ([]*types.Module, error)

		Create(module *types.Module) (*types.Module, error)
		Update(module *types.Module) (*types.Module, error)
		DeleteByID(moduleID uint64) error

		FieldNames(mod *types.Module) ([]string, error)
	}
)

func Module() ModuleService {
	return (&module{}).With(context.Background())
}

func (s *module) With(ctx context.Context) ModuleService {
	db := repository.DB(ctx)
	return &module{
		db:         db,
		ctx:        ctx,
		moduleRepo: repository.Module(ctx, db),
		pageRepo:   repository.Page(ctx, db),
	}
}

func (s *module) FindByID(id uint64) (*types.Module, error) {
	mod, err := s.moduleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.preload(mod); err != nil {
		return nil, err
	}
	return mod, err
}

func (s *module) Find() ([]*types.Module, error) {
	return s.moduleRepo.Find()
}

func (s *module) Create(mod *types.Module) (*types.Module, error) {
	if len(mod.Fields) == 0 {
		return nil, errors.New("Error creating module: no fields")
	}
	return s.moduleRepo.Create(mod)
}

func (s *module) Update(module *types.Module) (m *types.Module, err error) {
	validate := func() error {
		if module.ID == 0 {
			return errors.New("Error updating module: invalid ID")
		} else if m, err = s.moduleRepo.FindByID(module.ID); err != nil {
			return errors.Wrap(err, "Error while loading module for update")
		} else {
			module.CreatedAt = m.CreatedAt
		}

		if len(module.Fields) == 0 {
			return errors.New("Error updating module: no fields")
		}

		return nil
	}

	if err = validate(); err != nil {
		return nil, err
	}

	return m, s.db.Transaction(func() (err error) {
		m, err = s.moduleRepo.Update(module)
		return
	})
}

func (s *module) DeleteByID(id uint64) error {
	return s.moduleRepo.DeleteByID(id)
}

func (s *module) FieldNames(mod *types.Module) ([]string, error) {
	return s.moduleRepo.FieldNames(mod)
}
