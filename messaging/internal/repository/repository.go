package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/auth"
)

type (
	repository struct {
		ctx context.Context
		dbh *factory.DB
	}
)

// DB produces a contextual DB handle
func DB(ctx context.Context) *factory.DB {
	return factory.Database.MustGet("messaging").With(ctx)
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
