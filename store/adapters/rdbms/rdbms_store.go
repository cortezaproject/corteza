package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exec"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type (
	sqlizer interface {
		ToSQL() (string, []interface{}, error)
	}

	scanner interface {
		Scan(...interface{}) error
	}

	txRetryOnErrHandler func(int, error) bool

	Functions struct {
		// returns lower case text
		LOWER func(interface{}) exp.SQLFunctionExpression

		// returns date part of the input (YYYY-MM-DD)
		DATE func(interface{}) exp.SQLFunctionExpression
	}

	Store struct {
		DB sqlx.ExtContext

		// DAL connection
		DAL dal.Connection

		// Logger for connection
		Logger *zap.Logger

		// data definer interface use for schema information lookups and modification
		DataDefiner ddl.DataDefiner

		Dialect goqu.DialectWrapper

		// set to -1 to disable transactions
		TxRetryLimit int

		// TxRetryErrHandler should return true if transaction should be retried
		//
		// Because retry algorithm varies between concrete rdbms implementations
		//
		// Handler must return true if failed transaction should be replied
		// and false if we're safe to terminate it
		TxRetryErrHandler txRetryOnErrHandler

		ErrorHandler store.ErrorHandler

		Functions *Functions

		// additional (per-resource-type) filters used when searching
		// these filters can modify expression used for querying the database
		Filters *extendedFilters
	}
)

// SetDefaults fn sets all defaults that need to be set
func (s *Store) SetDefaults() {
	if s.TxRetryErrHandler == nil {
		s.TxRetryErrHandler = func(i int, err error) bool { return false }
	}

	if s.ErrorHandler == nil {
		s.ErrorHandler = func(err error) error { return err }
	}

	if s.Functions == nil {
		s.Functions = DefaultFunctions()
	}

	if s.Filters == nil {
		s.Filters = DefaultFilters()
	}
}

// WithTx spins up new store instance with transaction
func (s *Store) withTx(tx sqlx.ExtContext) *Store {
	return &Store{
		DB: tx,

		Logger:            s.Logger,
		DataDefiner:       s.DataDefiner,
		Dialect:           s.Dialect,
		TxRetryErrHandler: s.TxRetryErrHandler,
		ErrorHandler:      s.ErrorHandler,
		Functions:         s.Functions,
		Filters:           s.Filters,
	}
}

func (s Store) Exec(ctx context.Context, q sqlizer) error {
	var (
		query, args, err = q.ToSQL()
	)

	if err != nil {
		return fmt.Errorf("could not build query: %w", err)
	}

	_, err = s.DB.ExecContext(ctx, query, args...)
	return store.HandleError(err, s.ErrorHandler)
}

func (s Store) Query(ctx context.Context, q sqlizer) (*sql.Rows, error) {
	var (
		rr *sql.Rows

		query, args, err = q.ToSQL()
	)

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	rr, err = s.DB.QueryContext(ctx, query, args...)
	if err = store.HandleError(err, s.ErrorHandler); err != nil {
		return nil, err
	}

	return rr, nil
}

func (s Store) QueryOne(ctx context.Context, q sqlizer, dst interface{}) (err error) {
	var (
		rows *sql.Rows
	)

	rows, err = s.Query(ctx, q)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return store.ErrNotFound.Stack(1)
	}

	return exec.NewScanner(rows).ScanStruct(dst)
}

// ToDalConn uses store as DAL connection
func (s Store) ToDalConn() dal.Connection {
	return s.DAL
}

func DefaultFunctions() *Functions {
	f := &Functions{}

	if f.LOWER == nil {
		f.LOWER = func(value interface{}) exp.SQLFunctionExpression {
			return goqu.Func("LOWER", value)
		}
	}

	if f.DATE == nil {
		f.DATE = func(value interface{}) exp.SQLFunctionExpression {
			return goqu.Func("DATE", value)
		}
	}

	return f
}
