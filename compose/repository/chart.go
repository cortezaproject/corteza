package repository

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	ChartRepository interface {
		With(ctx context.Context, db *factory.DB) ChartRepository

		FindByID(namespaceID, chartID uint64) (*types.Chart, error)
		FindByHandle(namespaceID uint64, handle string) (c *types.Chart, err error)
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
	ErrChartNotFound        = repositoryError("ChartNotFound")
	ErrChartHandleNotUnique = repositoryError("ChartHandleNotUnique")
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
		"id",
		"rel_namespace",
		"handle",
		"name",
		"config",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r chart) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table()).
		Where("deleted_at IS NULL")
}

func (r chart) FindByID(namespaceID, chartID uint64) (*types.Chart, error) {
	return r.findOneBy(namespaceID, "id", chartID)
}

func (r chart) FindByHandle(namespaceID uint64, handle string) (*types.Chart, error) {
	return r.findOneBy(namespaceID, "LOWER(handle)", strings.ToLower(strings.TrimSpace(handle)))
}

func (r chart) findOneBy(namespaceID uint64, field string, value interface{}) (*types.Chart, error) {
	var (
		c = &types.Chart{}

		q = r.query().
			Where(squirrel.Eq{field: value, "rel_namespace": namespaceID})

		err = rh.FetchOne(r.db(), q, c)
	)

	if err != nil {
		return nil, err
	} else if c.ID == 0 {
		return nil, ErrChartNotFound
	}

	return c, nil
}

func (r chart) Find(filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id ASC"
	}

	query := r.query()

	if filter.NamespaceID > 0 {
		query = query.Where(squirrel.Eq{"rel_namespace": filter.NamespaceID})
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(name)": q},
		})
	}

	if f.Handle != "" {
		query = query.Where("LOWER(handle) = LOWER(?)", f.Handle)
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
}

func (r chart) Create(mod *types.Chart) (*types.Chart, error) {
	mod.ID = factory.Sonyflake.NextID()
	rh.SetCurrentTimeRounded(&mod.CreatedAt)
	mod.UpdatedAt = nil

	return mod, r.db().Insert(r.table(), mod)
}

func (r chart) Update(mod *types.Chart) (*types.Chart, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	return mod, r.db().Update(r.table(), mod, "id")
}

func (r chart) DeleteByID(namespaceID, chartID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		chartID,
	)

	return err
}
