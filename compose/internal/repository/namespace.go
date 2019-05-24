package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	NamespaceRepository interface {
		With(ctx context.Context, db *factory.DB) NamespaceRepository

		FindByID(id uint64) (*types.Namespace, error)
		Find(filter types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
		Create(mod *types.Namespace) (*types.Namespace, error)
		Update(mod *types.Namespace) (*types.Namespace, error)
		DeleteByID(id uint64) error
	}

	namespace struct {
		*repository
	}
)

const (
	ErrNamespaceNotFound = repositoryError("NamespaceNotFound")
)

func Namespace(ctx context.Context, db *factory.DB) NamespaceRepository {
	return (&namespace{}).With(ctx, db)
}

func (r namespace) table() string {
	return "compose_namespace"
}

func (r namespace) columns() []string {
	return []string{
		"id",
		"name",
		"slug",
		"enabled",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

func (r namespace) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")

}

func (r *namespace) With(ctx context.Context, db *factory.DB) NamespaceRepository {
	return &namespace{
		repository: r.repository.With(ctx, db),
	}
}

func (r *namespace) FindByID(namespaceID uint64) (*types.Namespace, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", namespaceID)

		n = &types.Namespace{}
	)

	return n, isFound(r.fetchOne(n, query), n.ID > 0, ErrNamespaceNotFound)
}

func (r *namespace) Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	f = filter

	query := r.query()
	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("name like ? OR slug like ?", q, q)
	}

	if f.Count, err = r.count(query); err != nil || f.Count == 0 {
		return
	}

	if f.Page > 0 {
		query = query.Offset(uint64(f.PerPage * f.Page))
	}

	query = query.
		Columns(r.columns()...).
		OrderBy("id ASC")

	return set, f, r.fetchPaged(&set, query, f.Page, f.PerPage)
}

func (r *namespace) Create(mod *types.Namespace) (*types.Namespace, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.table(), mod)
}

func (r *namespace) Update(mod *types.Namespace) (*types.Namespace, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	return mod, r.db().Replace(r.table(), mod)
}

func (r *namespace) DeleteByID(namespaceID uint64) error {
	_, err := r.db().Exec("UPDATE "+r.table()+" SET deleted_at = NOW() WHERE id = ?", namespaceID)
	return err
}
