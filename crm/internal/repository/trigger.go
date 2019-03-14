package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/types"
)

type (
	TriggerRepository interface {
		With(ctx context.Context, db *factory.DB) TriggerRepository

		FindByID(id uint64) (*types.Trigger, error)
		Find(filter types.TriggerFilter) (types.TriggerSet, error)
		Create(mod *types.Trigger) (*types.Trigger, error)
		Update(mod *types.Trigger) (*types.Trigger, error)
		DeleteByID(id uint64) error
	}

	trigger struct {
		*repository
	}
)

func Trigger(ctx context.Context, db *factory.DB) TriggerRepository {
	return (&trigger{}).With(ctx, db)
}

func (r *trigger) With(ctx context.Context, db *factory.DB) TriggerRepository {
	return &trigger{
		repository: r.repository.With(ctx, db),
	}
}

func (r *trigger) FindByID(id uint64) (*types.Trigger, error) {
	mod := &types.Trigger{}
	query := r.query().Where("id = ?", id)

	if sql, args, err := query.ToSql(); err != nil {
		return nil, err
	} else {
		return mod, r.db().Get(mod, sql, args...)
	}
}

func (r *trigger) Find(filter types.TriggerFilter) (mod types.TriggerSet, err error) {
	query := r.query()

	if filter.ModuleID > 0 {
		query = query.Where("rel_module = ?", filter.ModuleID)
	}

	if sql, args, err := query.ToSql(); err != nil {
		return nil, err
	} else {
		return mod, r.db().Select(&mod, sql, args...)
	}
}

func (r trigger) query() (query squirrel.SelectBuilder) {
	query = squirrel.Select().
		Columns(
			"id", "name", "actions", "enabled", "source", "rel_module",
			"created_at", "updated_at", "deleted_at").
		From("crm_trigger").
		OrderBy("id DESC")

	return
}

func (r *trigger) Create(mod *types.Trigger) (*types.Trigger, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("crm_trigger", mod)
}

func (r *trigger) Update(mod *types.Trigger) (*types.Trigger, error) {
	now := time.Now()
	mod.UpdatedAt = &now
	return mod, r.db().Replace("crm_trigger", mod)
}

func (r *trigger) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_trigger WHERE id = ?", id)
	return err
}
