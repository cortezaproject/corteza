package repository

import (
	"context"
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
	tx := r.tx
	if tx == nil {
		tx = factory.Database.MustGet().With(ctx)
	}

	txr := &repository{ctx: ctx, tx: tx}

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
	// @todo implementation
	return r.tx.Begin()
}

func (r *repository) Commit() error {
	// @todo implementation
	return r.tx.Commit()
}

func (r *repository) Rollback() error {
	// @todo implementation
	return r.tx.Rollback()
}

func (r *repository) db() *factory.DB {
	return r.tx.With(r.ctx)
}
