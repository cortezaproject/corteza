package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	ChartRepository interface {
		With(ctx context.Context, db *factory.DB) ChartRepository

		FindByID(id uint64) (*types.Chart, error)
		Find() (types.ChartSet, error)
		Create(mod *types.Chart) (*types.Chart, error)
		Update(mod *types.Chart) (*types.Chart, error)
		DeleteByID(id uint64) error
	}

	chart struct {
		*repository
	}
)

const sqlChartColumns = "id, name, config, " +
	"created_at, updated_at, deleted_at"

const sqlChartSelect = "SELECT " + sqlChartColumns + " FROM crm_chart"

func Chart(ctx context.Context, db *factory.DB) ChartRepository {
	return (&chart{}).With(ctx, db)
}

func (r *chart) With(ctx context.Context, db *factory.DB) ChartRepository {
	return &chart{
		repository: r.repository.With(ctx, db),
	}
}

func (r *chart) FindByID(id uint64) (*types.Chart, error) {
	mod := &types.Chart{}
	if err := r.db().Get(mod, sqlChartSelect+" WHERE id = ?", id); err != nil {
		return nil, err
	}
	return mod, nil
}

func (r *chart) Find() (types.ChartSet, error) {
	mod := types.ChartSet{}
	return mod, r.db().Select(&mod, sqlChartSelect+" ORDER BY id ASC")
}

func (r *chart) Create(mod *types.Chart) (*types.Chart, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("crm_chart", mod)
}

func (r *chart) Update(mod *types.Chart) (*types.Chart, error) {
	now := time.Now()
	mod.UpdatedAt = &now
	return mod, r.db().Replace("crm_chart", mod)
}

func (r *chart) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_chart WHERE id = ?", id)
	return err
}
