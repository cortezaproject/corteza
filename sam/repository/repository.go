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

	Transactionable interface {
		BeginWith(ctx context.Context, callback BeginCallback) error
		Begin() error
		Rollback() error
		Commit() error
	}

	Contextable interface {
		WithCtx(ctx context.Context) Interfaces
	}

	Interfaces interface {
		Transactionable
		Contextable

		Attachment
		Channel
		Message
		Organisation
		Reaction
		Team
		User
	}

	BeginCallback func(r Interfaces) error
)

func New() *repository {
	return &repository{ctx: context.Background()}
}

func (r *repository) WithCtx(ctx context.Context) Interfaces {
	return &repository{ctx: ctx, tx: r.tx}
}

func (r *repository) BeginWith(ctx context.Context, callback BeginCallback) error {

	txr := &repository{ctx: ctx}

	if err := txr.Begin(); err != nil {
		return err
	}

	if err := callback(txr); err != nil {
		if err := txr.Rollback(); err != nil {
			return err
		}

		return err
	}

	return txr.Commit()
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
