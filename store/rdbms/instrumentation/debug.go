package instrumentation

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/ngrok/sqlmw"
	"go.uber.org/zap"
)

type (
	// Debug instrumentation wrapper for sql driver
	//
	// All calls (with context) are wrapped (using github.com/ngrok/sqlmw package)
	// and that allows us to log all traffic going in and out of the database
	debug struct {
		sqlmw.NullInterceptor
		fallback *zap.Logger
	}
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

func Debug() *debug {
	return &debug{fallback: zap.NewNop()}
}

// Returns logger from context with additional options
//
// fn tries checks context or uses fallback logger
func (ld debug) log(ctx context.Context) *zap.Logger {
	return logger.
		ContextValue(ctx, logger.Default(), ld.fallback).
		WithOptions(
			// get all the way back to RDBMS layer
			// @todo not sure that all functions below
			//       have the same stack trace, might
			//       need adjustment on individual fn bellow
			zap.AddCallerSkip(9),
		)
}

func (debug) argToZapFields(args []driver.NamedValue) []zap.Field {
	var (
		name string
		out  = make([]zap.Field, len(args))
	)
	for i := range args {
		name = "arg."
		if args[i].Name != "" {
			name += args[i].Name
		} else {
			name += fmt.Sprintf("%d", args[i].Ordinal)
		}

		out[i] = zap.Any(name, args[i].Value)
	}

	return out
}

func (ld debug) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	var (
		startedAt = now()
		tx, err   = conn.BeginTx(ctx, txOpts)
		log       = ld.log(ctx).With(zap.Duration("duration", time.Since(startedAt)))
		message   = "conn/begin"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return tx, nil
}

func (ld debug) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	var (
		startedAt      = now()
		statement, err = conn.PrepareContext(ctx, query)
		log            = ld.log(ctx).With(
			zap.Duration("duration", time.Since(startedAt)),
			zap.Any("query", query),
		)
		message = "conn/prepare"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return statement, nil
}

func (ld debug) ConnPing(ctx context.Context, conn driver.Pinger) error {
	var (
		startedAt = now()
		err       = conn.Ping(ctx)
		log       = ld.log(ctx).With(zap.Duration("duration", time.Since(startedAt)))
		message   = "conn/ping"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return err
	}

	log.Debug(message)
	return nil
}

func (ld debug) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	var (
		startedAt   = now()
		result, err = conn.ExecContext(ctx, query, args)
		log         = ld.log(ctx).With(
			zap.Duration("duration", time.Since(startedAt)),
			zap.Any("query", query),
		).With(ld.argToZapFields(args)...)
		message = "conn/exec"
	)

	if err != nil {
		if errors.Is(err, driver.ErrSkip) {
			// Not actually an error; just informing caller
			// that it should take an alternative path
			return nil, err
		}

		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return result, nil
}

func (ld debug) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	var (
		startedAt = now()
		rows, err = conn.QueryContext(ctx, query, args)
		log       = ld.log(ctx).With(
			zap.String("query", query),
			zap.Duration("duration", time.Since(startedAt)),
		).With(ld.argToZapFields(args)...)
		message = "statement/query"
	)

	if err != nil {
		if errors.Is(err, driver.ErrSkip) {
			// Not actually an error; just informing caller
			// that it should take an alternative path
			return nil, err
		}

		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return rows, nil
}

func (ld debug) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, _ string, args []driver.NamedValue) (driver.Result, error) {
	var (
		startedAt   = now()
		result, err = stmt.ExecContext(ctx, args)
		log         = ld.log(ctx).With(
			zap.Duration("duration", time.Since(startedAt)),
		).With(ld.argToZapFields(args)...)
		message = "statement/exec"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return result, nil
}

func (ld debug) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, _ string, args []driver.NamedValue) (driver.Rows, error) {
	var (
		startedAt = now()
		rows, err = stmt.QueryContext(ctx, args)
		log       = ld.log(ctx).With(
			zap.Duration("duration", time.Since(startedAt)),
		).With(ld.argToZapFields(args)...)
		message = "statement/query"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return nil, err
	}

	log.Debug(message)
	return rows, nil
}

func (ld debug) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	var (
		startedAt = now()
		err       = stmt.Close()
		log       = ld.log(ctx).With(zap.Duration("duration", time.Since(startedAt)))
		message   = "statement/close"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return err
	}

	log.Debug(message)
	return nil
}

func (ld debug) TxCommit(ctx context.Context, tx driver.Tx) error {
	var (
		startedAt = now()
		err       = tx.Commit()
		log       = ld.log(ctx).With(zap.Duration("duration", time.Since(startedAt)))
		message   = "tx/commit"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return err
	}

	log.Debug(message)
	return nil
}

func (ld debug) TxRollback(ctx context.Context, tx driver.Tx) error {
	var (
		startedAt = now()
		err       = tx.Rollback()
		log       = ld.log(ctx).With(zap.Duration("duration", time.Since(startedAt)))
		message   = "tx/rollback"
	)

	if err != nil {
		log.Error(message, zap.Error(err))
		return err
	}

	log.Warn(message)
	return nil
}
