package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	TriggerRepository interface {
		With(ctx context.Context, db *factory.DB) TriggerRepository

		FindByID(namespaceID, triggerID uint64) (*types.Trigger, error)
		Find(filter types.TriggerFilter) (set types.TriggerSet, f types.TriggerFilter, err error)
		Create(mod *types.Trigger) (*types.Trigger, error)
		Update(mod *types.Trigger) (*types.Trigger, error)
		DeleteByID(namespaceID, triggerID uint64) error
	}

	trigger struct {
		*repository
	}
)

const (
	ErrTriggerNotFound = repositoryError("TriggerNotFound")
)

func Trigger(ctx context.Context, db *factory.DB) TriggerRepository {
	return (&trigger{}).With(ctx, db)
}

func (r trigger) With(ctx context.Context, db *factory.DB) TriggerRepository {
	return &trigger{
		repository: r.repository.With(ctx, db),
	}
}

func (r trigger) table() string {
	return "compose_trigger"
}

func (r trigger) columns() []string {
	return []string{
		"id", "rel_namespace", "name",
		"actions", "enabled", "source", "rel_module",
		"created_at", "updated_at", "deleted_at",
	}
}

func (r trigger) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")
}

func (r trigger) FindByID(namespaceID, triggerID uint64) (*types.Trigger, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", triggerID)

		c = &types.Trigger{}
	)

	if namespaceID > 0 {
		query = query.Where("rel_namespace = ?", namespaceID)
	}

	return c, isFound(r.fetchOne(c, query), c.ID > 0, ErrTriggerNotFound)
}

func (r trigger) Find(filter types.TriggerFilter) (set types.TriggerSet, f types.TriggerFilter, err error) {
	f = filter

	query := r.query()

	if filter.NamespaceID > 0 {
		query = query.Where("rel_namespace = ?", filter.NamespaceID)
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

func (r trigger) Create(mod *types.Trigger) (*types.Trigger, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now().Truncate(time.Second)

	return mod, r.db().Insert(r.table(), mod)
}

func (r trigger) Update(mod *types.Trigger) (*types.Trigger, error) {
	now := time.Now().Truncate(time.Second)
	mod.UpdatedAt = &now
	return mod, r.db().Replace(r.table(), mod)
}

func (r trigger) DeleteByID(namespaceID, triggerID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		triggerID,
	)

	return err
}
