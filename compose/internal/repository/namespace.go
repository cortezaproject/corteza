package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/compose/types"
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

func (r *namespace) With(ctx context.Context, db *factory.DB) NamespaceRepository {
	return &namespace{
		repository: r.repository.With(ctx, db),
	}
}

func (r *namespace) FindByID(id uint64) (n *types.Namespace, err error) {
	var (
		sql  string
		args []interface{}
	)

	n = &types.Namespace{}

	qb := r.query().
		Columns(r.columns()...).
		Where("id = ?", id)

	if sql, args, err = qb.ToSql(); err != nil {
		return
	}

	if err = r.db().Get(n, sql, args...); err != nil {
		return
	}

	if n == nil || n.ID == 0 {
		return nil, ErrNamespaceNotFound
	}

	return
}

func (r *namespace) Find(filter types.NamespaceFilter) (nn types.NamespaceSet, f types.NamespaceFilter, err error) {
	var (
		sql  string
		args []interface{}
	)

	f = filter

	if f.PerPage > 100 {
		f.PerPage = 100
	} else if f.PerPage == 0 {
		f.PerPage = 50
	}

	qb := r.query()
	if f.Query != "" {
		q := "%" + f.Query + "%"
		qb = qb.Where("name like ? OR slug like ?", q, q)
	}

	{
		cq := qb.Column(squirrel.Alias(squirrel.Expr("COUNT(*)"), "count"))
		if sql, args, err = cq.ToSql(); err != nil {
			return
		}

		if err = r.db().Get(&f.Count, sql, args...); err != nil {
			return
		}

		if f.Count == 0 {
			// No rows with this filter no need to continue
			return
		}
	}

	if f.Page > 0 {
		qb = qb.Offset(uint64(f.PerPage * f.Page))
	}

	qb = qb.
		Limit(uint64(f.PerPage)).
		Columns(r.columns()...).
		OrderBy("id ASC")

	if sql, args, err = qb.ToSql(); err != nil {
		return
	}

	if err = r.db().Select(&nn, sql, args...); err != nil {
		return
	}

	return
}

func (r namespace) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")

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

func (r *namespace) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM compose_namespace WHERE id=?", id)
	return err
}
