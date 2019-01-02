package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	ModuleRepository interface {
		With(ctx context.Context, db *factory.DB) ModuleRepository

		FindByID(id uint64) (*types.Module, error)
		Find() (types.ModuleSet, error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		DeleteByID(id uint64) error

		Fields(mod *types.Module) (types.ModuleFieldSet, error)
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

	if err := r.updateFields(mod.ID, mod.Fields); err != nil {
		return nil, err
	}

	return mod, r.db().Insert("crm_module", mod)
}

func (r *module) Update(mod *types.Module) (*types.Module, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	if err := r.updateFields(mod.ID, mod.Fields); err != nil {
		return nil, err
	}

	return mod, r.db().Replace("crm_module", mod)
}

func (r *module) updateFields(moduleID uint64, ff types.ModuleFieldSet) error {
	// @todo be more selective when deleting
	if _, err := r.db().Exec("DELETE FROM crm_module_form WHERE module_id = ?", moduleID); err != nil {
		return errors.Wrap(err, "Error updating module fields")
	}

	for idx, v := range ff {
		v.ModuleID = moduleID
		v.Place = idx
		if err := r.db().Replace("crm_module_form", v); err != nil {
			return errors.Wrap(err, "Error updating module fields")
		}
	}

	return nil
}

func (r *module) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_module WHERE id=?", id)
	return err
}

func (r *module) Fields(mod *types.Module) (ff types.ModuleFieldSet, err error) {
	return ff, r.db().Select(&ff, "select * from crm_module_form where module_id=? order by place asc", mod.ID)
}

// FieldNames returns a slice of field names, used for ordering record row columns
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
