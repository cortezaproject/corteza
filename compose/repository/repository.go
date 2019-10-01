package repository

import (
	"context"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	repository struct {
		ctx context.Context
		dbh *factory.DB
	}
)

// DB produces a contextual DB handle
func DB(ctx context.Context) *factory.DB {
	return factory.Database.MustGet("compose").With(ctx)
}

// Identity returns the User ID from context
func Identity(ctx context.Context) uint64 {
	return auth.GetIdentityFromContext(ctx).Identity()
}

// With updates repository and database contexts
func (r *repository) With(ctx context.Context, db *factory.DB) *repository {
	return &repository{
		ctx: ctx,
		dbh: db,
	}
}

// Context returns current active repository context
func (r *repository) Context() context.Context {
	return r.ctx
}

// db returns context-aware db handle
func (r *repository) db() *factory.DB {
	if r.dbh != nil {
		return r.dbh
	}
	return DB(r.ctx)
}

func (r repository) findOneInNamespaceBy(namespaceID uint64, q squirrel.SelectBuilder, eq squirrel.Eq, row interface{}) error {
	q = q.Where(eq)

	if namespaceID > 0 {
		q = q.Where("rel_namespace = ?", namespaceID)
	}

	if err := r.fetchOne(row, q); err != nil {
		row = nil
		return err
	}

	return nil
}

// Fetches single row from table
func (r repository) fetchOne(one interface{}, q squirrel.SelectBuilder) (err error) {
	return rh.FetchOne(r.db(), q, one)
}

// Counts all rows that match conditions from given query builder
func (r repository) count(q squirrel.SelectBuilder) (uint, error) {
	return rh.Count(r.db(), q)
}

// Fetches paged rows
func (r repository) fetchPaged(set interface{}, q squirrel.SelectBuilder, page, perPage uint) error {
	return rh.FetchPaged(r.db(), q, page, perPage, set)
}

func normalizePerPage(val, min, max, def uint) uint {
	return rh.NormalizePerPage(val, min, max, def)
}

func isFound(err error, valid bool, nerr error) error {
	return rh.IsFound(err, valid, nerr)
}
