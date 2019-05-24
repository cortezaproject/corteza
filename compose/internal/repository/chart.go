package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ChartRepository interface {
		With(ctx context.Context, db *factory.DB) ChartRepository

		FindByID(namespaceID, chartID uint64) (*types.Chart, error)
		Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)
		Create(mod *types.Chart) (*types.Chart, error)
		Update(mod *types.Chart) (*types.Chart, error)
		DeleteByID(namespaceID, chartID uint64) error
	}

	chart struct {
		*repository
	}
)

const (
	ErrChartNotFound = repositoryError("ChartNotFound")
)

func Chart(ctx context.Context, db *factory.DB) ChartRepository {
	return (&chart{}).With(ctx, db)
}

func (r chart) With(ctx context.Context, db *factory.DB) ChartRepository {
	return &chart{
		repository: r.repository.With(ctx, db),
	}
}

func (r chart) table() string {
	return "compose_chart"
}

func (r chart) columns() []string {
	return []string{
		"id", "rel_namespace", "name", "config",
		"created_at", "updated_at", "deleted_at",
	}
}

func (r chart) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")
}

func (r chart) FindByID(namespaceID, chartID uint64) (*types.Chart, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", chartID)

		c = &types.Chart{}
	)

	if namespaceID > 0 {
		query = query.Where("rel_namespace = ?", namespaceID)
	}

	return c, isFound(r.fetchOne(c, query), c.ID > 0, ErrChartNotFound)
}

func (r chart) Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
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

func (r chart) Create(mod *types.Chart) (*types.Chart, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now().Truncate(time.Second)

	return mod, r.db().Insert(r.table(), mod)
}

func (r chart) Update(mod *types.Chart) (*types.Chart, error) {
	now := time.Now().Truncate(time.Second)
	mod.UpdatedAt = &now
	return mod, r.db().Replace(r.table(), mod)
}

func (r chart) DeleteByID(namespaceID, chartID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		chartID,
	)

	return err
}
