package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ModuleRepository interface {
		With(ctx context.Context, db *factory.DB) ModuleRepository

		FindByID(namespaceID, moduleID uint64) (*types.Module, error)
		FindByName(namespaceID uint64, name string) (*types.Module, error)
		FindByHandle(namespaceID uint64, handle string) (*types.Module, error)
		Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
		FindFields(moduleIDs ...uint64) (ff types.ModuleFieldSet, err error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		UpdateFields(moduleID uint64, ff types.ModuleFieldSet, hasRecords bool) (err error)
		DeleteByID(namespaceID, moduleID uint64) error
	}

	module struct {
		*repository
	}
)

const (
	ErrModuleNotFound        = repositoryError("ModuleNotFound")
	ErrModuleNameNotUnique   = repositoryError("ModuleNameNotUnique")
	ErrModuleHandleNotUnique = repositoryError("ModuleHandleNotUnique")
)

func Module(ctx context.Context, db *factory.DB) ModuleRepository {
	return (&module{}).With(ctx, db)
}

func (r module) With(ctx context.Context, db *factory.DB) ModuleRepository {
	return &module{
		repository: r.repository.With(ctx, db),
	}
}

func (r module) table() string {
	return "compose_module"
}

func (r module) tableFields() string {
	return "compose_module_field"
}

func (r module) columns() []string {
	return []string{
		"id", "rel_namespace", "handle", "name", "json",
		"created_at", "updated_at", "deleted_at",
	}
}

func (r module) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")

}

func (r module) FindByID(namespaceID, moduleID uint64) (*types.Module, error) {
	return r.findOneBy(namespaceID, "id", moduleID)
}

func (r module) FindByHandle(namespaceID uint64, handle string) (*types.Module, error) {
	return r.findOneBy(namespaceID, "LOWER(handle)", strings.ToLower(strings.TrimSpace(handle)))
}

func (r module) FindByName(namespaceID uint64, name string) (*types.Module, error) {
	return r.findOneBy(namespaceID, "LOWER(name)", strings.ToLower(strings.TrimSpace(name)))
}

func (r module) findOneBy(namespaceID uint64, field string, value interface{}) (*types.Module, error) {
	var m = &types.Module{}

	err := r.findOneInNamespaceBy(
		namespaceID,
		r.query().Columns(r.columns()...),
		squirrel.Eq{field: value},
		m,
	)

	if err == nil && m.ID == 0 {
		return nil, ErrModuleNotFound
	}

	return m, nil
}

func (r module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	f = filter

	query := r.query()

	if filter.NamespaceID > 0 {
		query = query.Where("rel_namespace = ?", filter.NamespaceID)
	}

	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("name like ?", q)
	}

	if f.Name != "" {
		query = query.Where("LOWER(name) = ?", strings.ToLower(f.Name))
	}

	if f.Handle != "" {
		query = query.Where("LOWER(handle) = ?", strings.ToLower(f.Handle))
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	if f.Count, err = r.count(query); err != nil || f.Count == 0 {
		return
	}

	query = query.
		Columns(r.columns()...).
		OrderBy("id ASC")

	return set, f, r.fetchPaged(&set, query, f.Page, f.PerPage)
}

func (r module) Create(mod *types.Module) (*types.Module, error) {
	var err error

	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now().Truncate(time.Second)

	if err = r.db().Insert(r.table(), mod); err != nil {
		return nil, err
	}

	return mod, nil
}

func (r module) Update(mod *types.Module) (*types.Module, error) {
	now := time.Now().Truncate(time.Second)
	mod.UpdatedAt = &now

	return mod, r.db().Update(r.table(), mod, "id")
}

func (r module) UpdateFields(moduleID uint64, ff types.ModuleFieldSet, hasRecords bool) error {
	if existing, err := r.FindFields(moduleID); err != nil {
		return err
	} else {
		// Remove fields that do not exist anymore
		err = existing.Walk(func(e *types.ModuleField) error {
			if ff.FindByID(e.ID) == nil {
				return r.deleteFieldByID(moduleID, e.ID)
			}

			return nil
		})

		if err != nil {
			return err
		}

		now := time.Now().Truncate(time.Second)
		for idx, f := range ff {
			if e := existing.FindByID(f.ID); e != nil {
				f.CreatedAt = e.CreatedAt

				// We do not have any other code in place that would handle changes of field name and kind, so we need
				// to reset any changes made to the field.
				// @todo remove when we are able to handle field rename & type change
				if hasRecords {
					f.Name = e.Name
					f.Kind = e.Kind
				} else {
					f.UpdatedAt = &now
				}
			} else {
				f.ID = 0
			}

			if f.ID == 0 {
				f.ID = factory.Sonyflake.NextID()
				f.CreatedAt = now
				f.UpdatedAt = nil
			}

			f.ModuleID = moduleID
			f.Place = idx
			f.DeletedAt = nil

			if err := r.db().Replace(r.tableFields(), f); err != nil {
				return errors.Wrap(err, "Error updating module fields")
			}

		}
	}

	return nil
}

func (r module) deleteFieldByID(moduleID, fieldID uint64) error {
	_, err := r.db().Exec(
		fmt.Sprintf("DELETE FROM %s WHERE rel_module = ? AND id = ?", r.tableFields()),
		moduleID,
		fieldID,
	)

	return err
}

func (r module) DeleteByID(namespaceID, moduleID uint64) error {
	_, err := r.db().Exec(
		fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?", r.table()),
		namespaceID,
		moduleID,
	)

	return err
}

func (r module) FindFields(moduleIDs ...uint64) (ff types.ModuleFieldSet, err error) {
	if len(moduleIDs) == 0 {
		return
	}

	query := `SELECT id, rel_module, place, 
                     kind, name, label, options, 
                     is_private, is_required, is_visible, is_multi, default_value, 
                     created_at, updated_at, deleted_at
                FROM %s 
               WHERE rel_module IN (?) 
                 AND deleted_at IS NULL
               ORDER BY rel_module, place`

	query = fmt.Sprintf(query, r.tableFields())

	if sql, args, err := sqlx.In(query, moduleIDs); err != nil {
		return nil, err
	} else {
		return ff, r.db().Select(&ff, sql, args...)
	}
}
