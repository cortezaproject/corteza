package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"github.com/titpetric/factory"
)

type (
	Module interface {
		With(ctx context.Context) Module

		FindByID(id uint64) (*types.Module, error)
		Find() ([]*types.Module, error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		DeleteByID(id uint64) error
	}

	module struct {
		*repository
	}
)

func NewModule(ctx context.Context) Module {
	return (&module{}).With(ctx)
}

func (r *module) With(ctx context.Context) Module {
	return &module{
		repository: r.repository.With(ctx),
	}
}

// @todo: update to accepted DeletedAt column semantics from SAM

func (r *module) FindByID(id uint64) (*types.Module, error) {
	mod := &types.Module{}
	return mod, r.db().Get(mod, "SELECT * FROM crm_module WHERE id=?", id)
}

func (r *module) Find() ([]*types.Module, error) {
	mod := make([]*types.Module, 0)
	return mod, r.db().Select(&mod, "SELECT * FROM crm_module ORDER BY id ASC")
}

func (r *module) Create(mod *types.Module) (*types.Module, error) {
	mod.ID = factory.Sonyflake.NextID()
	return mod, r.db().Insert("crm_module", mod)
}

func (r *module) Update(mod *types.Module) (*types.Module, error) {
	return mod, r.db().Replace("crm_module", mod)
}

func (r *module) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_module WHERE id=?", id)
	return err
}
