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

// Return context-aware db handle
func (r *repository) db() *factory.DB {
	return r.dbh(r.ctx)
}
