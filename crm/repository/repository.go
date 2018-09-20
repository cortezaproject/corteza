package repository

import (
	"context"
	"github.com/titpetric/factory"
)

type (
	repository struct {
		ctx context.Context

		// Get database handle
		dbh func(ctxs ...context.Context) *factory.DB
	}
)

var _db *factory.DB

// DB returns a repository-wide singleton DB handle
func DB(ctxs ...context.Context) *factory.DB {
	if _db == nil {
		_db = factory.Database.MustGet()
	}
	for _, ctx := range ctxs {
		_db = _db.With(ctx)
		break
	}
	return _db
}

// With updates repository and database contexts
func (r *repository) With(ctx context.Context) *repository {
	res := &repository{
		ctx: ctx,
		dbh: DB,
	}
	if r != nil {
		res.dbh = r.dbh
	}
	return res
}

// db returns context-aware db handle
func (r *repository) db() *factory.DB {
	return r.dbh(r.ctx)
}
