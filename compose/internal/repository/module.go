package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/compose/types"
)

type (
	ModuleRepository interface {
		With(ctx context.Context, db *factory.DB) ModuleRepository

		FindByID(namespaceID, moduleID uint64) (*types.Module, error)
		Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
		FindFields(moduleIDs ...uint64) (ff types.ModuleFieldSet, err error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		DeleteByID(namespaceID, moduleID uint64) error
	}

	module struct {
		*repository
	}
)

const (
	ErrModuleNotFound = repositoryError("ModuleNotFound")
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

func (r module) columns() []string {
	return []string{
		"id", "rel_namespace", "name", "json",
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
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", moduleID)

		c = &types.Module{}
	)

	if namespaceID > 0 {
		query = query.Where("rel_namespace = ?", namespaceID)
	}

	return c, isFound(r.fetchOne(c, query), c.ID > 0, ErrModuleNotFound)
}

func (r module) Find(filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error) {
	f = filter
	f.PerPage = normalizePerPage(f.PerPage, 5, 100, 50)

	query := r.query()

	if filter.NamespaceID > 0 {
		query = query.Where("a.rel_namespace = ?", filter.NamespaceID)
	}

	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("name like ?", q)
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
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	if err := r.updateFields(mod.ID, mod.Fields); err != nil {
		return nil, err
	}

	return mod, r.db().Insert(r.table(), mod)
}

func (r module) Update(mod *types.Module) (*types.Module, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	if err := r.updateFields(mod.ID, mod.Fields); err != nil {
		return nil, err
	}

	return mod, r.db().Replace(r.table(), mod)
}

func (r module) updateFields(moduleID uint64, ff types.ModuleFieldSet) error {
	// @todo be more selective when deleting
	if _, err := r.db().Exec("DELETE FROM compose_module_form WHERE module_id = ?", moduleID); err != nil {
		return errors.Wrap(err, "Error updating module fields")
	}

	for idx, v := range ff {
		v.ModuleID = moduleID
		v.Place = idx
		if err := r.db().Replace("compose_module_form", v); err != nil {
			return errors.Wrap(err, "Error updating module fields")
		}
	}

	return nil
}

func (r module) DeleteByID(namespaceID, attachmentID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		attachmentID,
	)

	return err
}

func (r module) FindFields(moduleIDs ...uint64) (ff types.ModuleFieldSet, err error) {
	if len(moduleIDs) == 0 {
		return
	}

	if sql, args, err := sqlx.In("SELECT * FROM compose_module_form WHERE module_id IN (?) ORDER BY module_id AND place", moduleIDs); err != nil {
		return nil, err
	} else {
		return ff, r.db().Select(&ff, sql, args...)
	}
}
