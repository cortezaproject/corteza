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

func (s Store) Tx(ctx context.Context, fn func(context.Context, store.Storer) error) error {
	return tx(ctx, s.db, s.config, nil, func(ctx context.Context, tx dbLayer) error {
		return fn(ctx, s.withTx(tx))
	})
}

// tx begins a new db transaction and handles it's retries when possible
//
// It utilizes configured transaction error handlers and max-retry limits
// to determine if and how many times transaction should be retried
//
func tx(ctx context.Context, dbCandidate interface{}, cfg *Config, txOpt *sql.TxOptions, task func(context.Context, dbLayer) error) error {
	if cfg.TxDisabled {
		return task(ctx, dbCandidate.(dbLayer))
	}

	var (
		lastTaskErr error
		err         error
		db          *sqlx.DB
		tx          *sqlx.Tx
		try         = 1
	)

	switch dbCandidate.(type) {
	case dbTransactionMaker:
		// we can make a transaction, yay
		db = dbCandidate.(*sqlx.DB)
	case dbLayer:
		// Already in a transaction, run the given task and finish
		return task(ctx, dbCandidate.(dbLayer))
	default:
		return fmt.Errorf("could not use the db connection for transaction")
	}

	for {
		try++

		// Start transaction
		tx, err = db.BeginTxx(ctx, txOpt)
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
		// different errors, and we need to act accordingly
		if !cfg.TxRetryErrHandler(try, lastTaskErr) {
			return fmt.Errorf("failed to complete transaction: %w", lastTaskErr)
		}

		// tx error handlers can take current number of tries into account and
		// break the retry-loop earlier, but that might not be always the case
		//
		// We'll check the configured and hard-limit maximums
		if try >= cfg.TxMaxRetries || try >= TxRetryHardLimit {
			return fmt.Errorf("failed to perform transaction (tries: %d), last error: %w", try, lastTaskErr)
		}

		// Sleep (with a bit of kickback) before doing next retry
		time.Sleep(50 * time.Duration(try*50))
	}
}
