package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/types"
)

type (
	ModuleRepository interface {
		With(ctx context.Context, db *factory.DB) ModuleRepository

		Chart(r *request.ModuleChart) (interface{}, error)

		FindByID(id uint64) (*types.Module, error)
		Find() (types.ModuleSet, error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		DeleteByID(id uint64) error

		Fields(mod *types.Module) ([]*types.ModuleField, error)
		FieldNames(mod *types.Module) ([]string, error)
	}

	module struct {
		*repository
	}
)

func Module(ctx context.Context, db *factory.DB) ModuleRepository {
	return (&module{}).With(ctx, db)
}

func (r *module) With(ctx context.Context, db *factory.DB) ModuleRepository {
	return &module{
		repository: r.repository.With(ctx, db),
	}
}

// @todo: update to accepted DeletedAt column semantics from SAM

func (r *module) FindByID(id uint64) (*types.Module, error) {
	mod := &types.Module{}
	if err := r.db().Get(mod, "SELECT * FROM crm_module WHERE id=?", id); err != nil {
		return nil, err
	}
	return mod, nil
}

func (r *module) Find() (types.ModuleSet, error) {
	mod := types.ModuleSet{}
	return mod, r.db().Select(&mod, "SELECT * FROM crm_module ORDER BY id ASC")
}

func (r *module) Create(mod *types.Module) (*types.Module, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	fields := make([]types.ModuleField, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "No fields")
	}

	for idx, v := range fields {
		v.ModuleID = mod.ID
		v.Place = idx
		if err := r.db().Replace("crm_module_form", v); err != nil {
			return nil, errors.Wrap(err, "Error adding module fields")
		}
	}
	return mod, r.db().Insert("crm_module", mod)
}

func (r *module) Update(mod *types.Module) (*types.Module, error) {
	if mod.ID == 0 {
		return nil, errors.New("Error when saving module, invalid ID")
	}
	now := time.Now()
	mod.UpdatedAt = &now

	fields := make([]types.ModuleField, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "No fields")
	}

	for idx, v := range fields {
		v.ModuleID = mod.ID
		v.Place = idx
		if err := r.db().Replace("crm_module_form", v); err != nil {
			return nil, errors.Wrap(err, "Error adding module fields")
		}
	}
	return mod, r.db().Replace("crm_module", mod)
}

func (r *module) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_module WHERE id=?", id)
	return err
}

func (r *module) Fields(mod *types.Module) ([]*types.ModuleField, error) {
	fields := make([]*types.ModuleField, 0)
	return fields, r.db().Select(&fields, "select * from crm_module_form where module_id=? order by place asc", mod.ID)
}

// FieldNames returns a slice of field names, used for ordering content row columns
func (r *module) FieldNames(mod *types.Module) ([]string, error) {
	if fields, err := r.Fields(mod); err != nil {
		return []string{}, err
	} else {
		result := make([]string, len(fields))
		for k, v := range fields {
			result[k] = v.Name
		}
		return result, nil
	}
}
