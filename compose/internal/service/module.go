package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
)

type (
	module struct {
		db  *factory.DB
		ctx context.Context

		prmSvc PermissionsService

		moduleRepo repository.ModuleRepository
		pageRepo   repository.PageRepository
	}

	ModuleService interface {
		With(ctx context.Context) ModuleService

		FindByID(moduleID uint64) (*types.Module, error)
		Find() (types.ModuleSet, error)

		Create(module *types.Module) (*types.Module, error)
		Update(module *types.Module) (*types.Module, error)
		DeleteByID(moduleID uint64) error
	}
)

func Module() ModuleService {
	return (&module{
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc *module) With(ctx context.Context) ModuleService {
	db := repository.DB(ctx)
	return &module{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		moduleRepo: repository.Module(ctx, db),
		pageRepo:   repository.Page(ctx, db),
	}
}

func (svc *module) FindByID(id uint64) (m *types.Module, err error) {
	if m, err = svc.moduleRepo.FindByID(id); err != nil {
		return
	} else if !svc.prmSvc.CanReadModule(m) {
		return nil, errors.New("not allowed to access this module")
	}

	return
}

func (svc *module) Find() (mm types.ModuleSet, err error) {
	if mm, err = svc.moduleRepo.Find(); err != nil {
		return nil, err
	} else {
		return mm.Filter(func(m *types.Module) (bool, error) {
			return svc.prmSvc.CanReadModule(m), nil
		})
	}
}

func (svc *module) Create(mod *types.Module) (*types.Module, error) {
	if !svc.prmSvc.CanCreateModule(crmNamespace()) {
		return nil, errors.New("not allowed to create this module")
	}

	if len(mod.Fields) == 0 {
		return nil, errors.New("Error creating module: no fields")
	}

	return svc.moduleRepo.Create(mod)
}

func (svc *module) Update(module *types.Module) (m *types.Module, err error) {
	validate := func() error {
		if module.ID == 0 {
			return errors.New("Error updating module: invalid ID")
		} else if m, err = svc.moduleRepo.FindByID(module.ID); err != nil {
			return errors.Wrap(err, "Error while loading module for update")
		} else {
			if !svc.prmSvc.CanUpdateModule(m) {
				return errors.New("not allowed to update this module")
			}

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

	return m, svc.db.Transaction(func() (err error) {
		m, err = svc.moduleRepo.Update(module)
		return
	})
}

func (svc *module) DeleteByID(ID uint64) error {
	if m, err := svc.moduleRepo.FindByID(ID); err != nil {
		return errors.Wrap(err, "could not delete module")
	} else if !svc.prmSvc.CanDeleteModule(m) {
		return errors.New("not allowed to delete this module")
	}

	return svc.moduleRepo.DeleteByID(ID)
}
