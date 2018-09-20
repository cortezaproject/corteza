package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
)

type (
	repository struct {
		ctx context.Context

		// Get database handle
		dbh func(ctxs ...context.Context) *factory.DB
	}
)

var (
_db *factory.DB
_ctx context.Context
)

// DB returns a repository-wide singleton DB handle
func DB(ctxs ...context.Context) *factory.DB {
	if _db == nil {
		_db = factory.Database.MustGet()
	}
	for _, ctx := range ctxs {
		_db = _db.With(ctx)
		_ctx = ctx
		break
	}
	return _db
}

func Identity(ctx context.Context) uint64 {
	return auth.GetIdentityFromContext(ctx).Identity()
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

// Context returns current active repository context
func (r *repository) Context() context.Context {
	return r.ctx
}


// db returns context-aware db handle
func (r *repository) db() *factory.DB {
	return r.dbh(r.ctx)
}
