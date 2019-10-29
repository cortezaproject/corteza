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
	NamespaceRepository interface {
		With(ctx context.Context, db *factory.DB) NamespaceRepository

		FindByID(id uint64) (*types.Namespace, error)
		FindBySlug(slug string) (*types.Namespace, error)
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
	ErrNamespaceNotFound          = repositoryError("NamespaceNotFound")
	ErrNamespaceSlugNotUnique     = repositoryError("NamespaceSlugNotUnique")
	ErrNamespaceInvalidSlugFormat = repositoryError("NamespaceInvalidSlugFormat")
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
		Select(r.columns()...).
		From(r.table()).
		Where("deleted_at IS NULL")

}

func (r *namespace) With(ctx context.Context, db *factory.DB) NamespaceRepository {
	return &namespace{
		repository: r.repository.With(ctx, db),
	}
}

func (r *namespace) FindByID(namespaceID uint64) (*types.Namespace, error) {
	return r.findOneBy("id", namespaceID)
}

func (r *namespace) FindBySlug(slug string) (*types.Namespace, error) {
	return r.findOneBy("slug", slug)
}

func (r *namespace) findOneBy(field string, value interface{}) (*types.Namespace, error) {
	var (
		ns = &types.Namespace{}

		q = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, ns)
	)

	if err != nil {
		return nil, err
	} else if ns.ID == 0 {
		return nil, ErrNamespaceNotFound
	}

	return ns, nil
}

func (r *namespace) Find(filter types.NamespaceFilter) (set types.NamespaceSet, f types.NamespaceFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id ASC"
	}

	query := r.query()
	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(name)": q},
			squirrel.Like{"LOWER(slug)": q},
		})
	}

	if f.Slug != "" {
		query = query.Where(squirrel.Eq{"LOWER(slug)": strings.ToLower(f.Slug)})
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

func (r *namespace) Create(mod *types.Namespace) (*types.Namespace, error) {
	mod.ID = factory.Sonyflake.NextID()
	rh.SetCurrentTimeRounded(&mod.CreatedAt)
	mod.UpdatedAt = nil

	return mod, r.db().Insert(r.table(), mod)
}

func (r *namespace) Update(mod *types.Namespace) (*types.Namespace, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	return mod, r.db().Update(r.table(), mod, "id")
}

func (r *namespace) DeleteByID(namespaceID uint64) error {
	_, err := r.db().Exec("UPDATE "+r.table()+" SET deleted_at = NOW() WHERE id = ?", namespaceID)
	return err
}
