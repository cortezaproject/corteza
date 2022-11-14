package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
)

type (
	dbTransactionMaker interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	}
)

func (s *Store) Tx(ctx context.Context, fn func(context.Context, store.Storer) error) error {
	if s.TxRetryLimit < 0 {
		return fn(ctx, s)
	}

	return txHandler(ctx, s.DB, s.TxRetryLimit, s.TxRetryErrHandler, func(ctx context.Context, tx sqlx.ExtContext) error {
		return fn(ctx, s.withTx(tx))
	})
}

// tx begins a new db transaction and handles retries when possible
//
// It utilizes configured transaction error handlers and max-retry limits
// to determine if and how many times transaction should be retried
//
func txHandler(ctx context.Context, dbc interface{}, max int, reh txRetryOnErrHandler, task func(context.Context, sqlx.ExtContext) error) error {
	var (
		lastTaskErr error
		err         error
		db          *sqlx.DB
		tx          *sqlx.Tx
		try         = 1
	)

	switch dbc.(type) {
	case dbTransactionMaker:
		// we can make a transaction!
		db = dbc.(*sqlx.DB)
	case sqlx.ExtContext:
		// Already in a transaction, run the given task and finish
		return task(ctx, dbc.(sqlx.ExtContext))
	default:
		return fmt.Errorf("could not use the db connection for transaction")
	}

	for {
		try++

		// Start transaction
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}

		if lastTaskErr = task(ctx, tx); lastTaskErr == nil {
			// Task completed successfully
			return tx.Commit()
		}

		// No matter the cause of the error, rollback the transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction (tries: %d) on error %v: %w", try, lastTaskErr, rollbackErr)
		}

		if errors.IsAny(lastTaskErr) {
			// not a store-internal/transaction error and can be returned right away,
			// no need to re-run the transaction
			return lastTaskErr
		}

		// Call the configured transaction retry error handler
		// if this particular error should be retried or not
		//
		// We cannot generalize here because different store implementation have
		// different errors; we need to act accordingly
		if reh != nil && !reh(try, lastTaskErr) {
			return fmt.Errorf("failed to complete transaction: %w", lastTaskErr)
		}

		// tx error handlers can take current number of tries into account and
		// break the retry-loop earlier, but that might not be always the case
		//
		// We'll check the configured and hard-limit maximums
		if try >= max || try >= TxRetryHardLimit {
			return fmt.Errorf("failed to perform transaction (tries: %d), last error: %w", try, lastTaskErr)
		}

		// Sleep (with a bit of kickback) before doing next retry
		time.Sleep(50 * time.Duration(try*50))
	}
}
