package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
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
	}

	module struct {
		*repository
	}
)

const (
	sqlModuleColumns = `
		id, name, json, 
		created_at, updated_at, deleted_at
	`
	sqlModuleSelect = `
		SELECT ` + sqlModuleColumns + ` FROM crm_module WHERE deleted_at IS NULL
	`
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

func (r *module) FindByID(id uint64) (mod *types.Module, err error) {
	mod = &types.Module{}

	if err = r.db().Get(mod, sqlModuleSelect+" AND id = ? ", id); err != nil {
		return
	}

	if mod.Fields, err = r.fields(id); err != nil {
		return
	}

	return
}

func (r *module) Find() (mm types.ModuleSet, err error) {
	if err = r.db().Select(&mm, sqlModuleSelect+" ORDER BY id ASC"); err != nil {
		return
	}

	var ff types.ModuleFieldSet
	if ff, err = r.fields(mm.IDs()...); err != nil {
		return
	} else {
		_ = ff.Walk(func(f *types.ModuleField) error {
			mm.FindByID(f.ModuleID).Fields = append(mm.FindByID(f.ModuleID).Fields, f)
			return nil
		})
	}

	return mm, nil
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

func (r *module) fields(IDs ...uint64) (ff types.ModuleFieldSet, err error) {
	if len(IDs) == 0 {
		return
	}

	if sql, args, err := sqlx.In("SELECT * FROM crm_module_form WHERE module_id IN (?) ORDER BY module_id AND place", IDs); err != nil {
		return nil, err
	} else {
		return ff, r.db().Select(&ff, sql, args...)
	}
}

// // FieldNames returns a slice of field names, used for ordering record row columns
// func (r *module) FieldNames(mod *types.Module) ([]string, error) {
// 	if fields, err := r.Fields(mod.ID); err != nil {
// 		return []string{}, err
// 	} else {
// 		result := make([]string, len(fields))
// 		for k, v := range fields {
// 			result[k] = v.Name
// 		}
// 		return result, nil
// 	}
// }
