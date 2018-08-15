package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

type (
	repository struct {
		ctx context.Context

		// Current transaction
		tx *factory.DB
	}
)

// With updates repository and database contexts
func (r *repository) With(ctx context.Context) *repository {
	return &repository{
		ctx: ctx,
		tx: r.db().With(r.ctx),
	}
}

func (r *repository) Begin() error {
	return r.db().Begin()
}

func (r *repository) Commit() error {
	return errors.Wrap(r.db().Commit(), "Can not commit changes")
}

func (r *repository) Rollback() error {
	return errors.Wrap(r.db().Rollback(), "Can not rollback changes")
}

func (r *repository) db() *factory.DB {
	if r.tx == nil {
		r.tx = factory.Database.MustGet().With(r.ctx)
	}
	return r.tx
}
