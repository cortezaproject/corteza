package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/ddl"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"time"
)

// persistance layer
//
// all functions go under one struct
//   why? because it will be easier to initialize and pass around
//
// each domain will be in it's own file
//
// connection logic will be built in the persistence layer (making pkg/db obsolete)
//

type (
	schemaUpgradeGenerator interface {
		TableExists(string) bool
		CreateTable(t *ddl.Table) string
		CreateIndexes(ii ...*ddl.Index) []string
	}

	Store struct {
		config *Config

		// Schema upgrade generator converts internal upgrade config
		// to implementation specific SQL
		sug schemaUpgradeGenerator

		db dbLayer
	}

	dbLayer interface {
		sqlx.ExecerContext
		SelectContext(context.Context, interface{}, string, ...interface{}) error
		GetContext(context.Context, interface{}, string, ...interface{}) error
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	}

	ModuleFieldTypeDetector interface {
		IsBoolean() bool
		IsNumeric() bool
		IsDateTime() bool
		IsRef() bool
	}

	dbTransactionMaker interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	}
)

var log = logger.MakeDebugLogger().
	//return logger.Default().
	Named("store.rdbms").
	WithOptions(zap.AddCallerSkip(2)).
	WithOptions(zap.AddStacktrace(zap.FatalLevel))

const (
	// TxRetryHardLimit is the absolute maximum retries we'll allow
	TxRetryHardLimit = 100

	DefaultSliceCapacity = 1000

	MinEnsureFetchLimit = 10
	MaxRefetches        = 100
)

func Connect(ctx context.Context, cfg *Config) (s *Store, err error) {
	if err = cfg.ParseExtra(); err != nil {
		return nil, err
	}

	cfg.SetDefaults()
	s = &Store{
		config: cfg,
	}

	if err = s.Connect(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// WithTx spins up new store instance with transaction
func (s *Store) withTx(tx dbLayer) *Store {
	return &Store{
		config: s.config,
		sug:    s.sug,
		db:     tx,
	}
}

// Temporary solution for logging
func (s *Store) log(ctx context.Context) *zap.Logger {
	// @todo extract logger from context
	return log
}

func (s *Store) Connect(ctx context.Context) error {
	s.log(ctx).Debug("opening connection", zap.String("driver", s.config.DriverName), zap.String("dsn", s.config.MaskedDSN()))
	db, err := sqlx.Open(s.config.DriverName, s.config.DataSourceName)
	healthcheck.Defaults().Add(dbHealthcheck(db), "Store/RDBMS/"+s.config.DriverName)
	if err != nil {
		return err
	}

	s.log(ctx).Debug(
		"setting connection parameters",
		zap.Int("MaxOpenConns", s.config.MaxOpenConns),
		zap.Duration("MaxLifetime", s.config.ConnMaxLifetime),
		zap.Int("MaxIdleConns", s.config.MaxIdleConns),
	)

	db.SetMaxOpenConns(s.config.MaxOpenConns)
	db.SetConnMaxLifetime(s.config.ConnMaxLifetime)
	db.SetMaxIdleConns(s.config.MaxIdleConns)

	if err = s.tryToConnect(ctx, db); err != nil {
		return err
	}

	s.db = db
	return err
}

func (s Store) tryToConnect(ctx context.Context, db *sqlx.DB) error {
	var (
		connErrCh = make(chan error, 1)
		patience  = time.Now().Add(s.config.ConnTryPatience)
	)

	//defer close(connErrCh)

	go func() {
		defer sentry.Recover()

		var (
			err error
			try = 0
		)

		for {
			try++

			if s.config.ConnTryMax <= try {
				connErrCh <- fmt.Errorf("could not connect in %d tries", try)
				return
			}

			if err = db.PingContext(ctx); err != nil {

				if time.Now().After(patience) {
					// don't make too much fuss
					// if we're in patience mode
					s.log(ctx).Warn(
						"could not connect to the database",
						zap.Error(err),
						zap.Int("try", try),
						zap.Float64("delay", s.config.ConnTryBackoffDelay.Seconds()),
					)
				}

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(s.config.ConnTryBackoffDelay):
					//	Wait before next try
					continue
				}
			}

			s.log(ctx).Debug("connected to the database")
			break
		}

		connErrCh <- err
	}()

	to := s.config.ConnTryTimeout * time.Duration(s.config.ConnTryMax*2)
	select {
	case err := <-connErrCh:
		return err
	case <-time.After(to):
		// Wait before next try
		return fmt.Errorf("timedout after %ds", to.Seconds())
	case <-ctx.Done():
		return fmt.Errorf("connection cancelled")
	}
}

func dbHealthcheck(db *sqlx.DB) func(ctx context.Context) error {
	return db.PingContext
}

func (s Store) Query(ctx context.Context, q squirrel.Sqlizer) (*sql.Rows, error) {
	var (
		start            = time.Now()
		query, args, err = q.ToSql()
	)

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	s.log(ctx).Debug(query, zap.Any("args", args), zap.Duration("duration", time.Now().Sub(start)))

	return s.db.QueryContext(ctx, query, args...)
}

// QueryRow returns row instead of filling in the passed struct
func (s Store) QueryRow(ctx context.Context, q squirrel.SelectBuilder) (*sql.Row, error) {
	var (
		start            = time.Now()
		query, args, err = q.ToSql()
	)

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	s.log(ctx).Debug(query, zap.Any("args", args), zap.Duration("duration", time.Now().Sub(start)))

	return s.db.QueryRowContext(ctx, query, args...), nil
}

func (s Store) Exec(ctx context.Context, sqlizer squirrel.Sqlizer) error {
	var (
		start            = time.Now()
		query, args, err = sqlizer.ToSql()
	)

	if err != nil {
		return err
	}

	s.log(ctx).Debug(query, zap.Any("args", args), zap.Duration("duration", time.Now().Sub(start)))

	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

func (s Store) Tx(ctx context.Context, fn func(context.Context, store.Storer) error) error {
	return tx(ctx, s.db, s.config, nil, func(ctx context.Context, tx dbLayer) error {
		return fn(ctx, s.withTx(tx))
	})
}

func (s Store) Truncate(ctx context.Context, table string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM "+table)
	return err
}

// SelectBuilder is a shorthand for squirrel.selectBuilder
//
// Sets passed table & columns and configured placeholder format
func (s Store) SelectBuilder(table string, cc ...string) squirrel.SelectBuilder {
	return squirrel.Select(cc...).From(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

// InsertBuilder is a shorthand for squirrel.insertBuilder
//
// Sets passed table and configured placeholder format
func (s Store) InsertBuilder(table string) squirrel.InsertBuilder {
	return squirrel.Insert(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

// UpdateBuilder is a shorthand for squirrel.updateBuilder
//
// Sets passed table and configured placeholder format
func (s Store) UpdateBuilder(table string) squirrel.UpdateBuilder {
	return squirrel.Update(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

// DeleteBuilder is a shorthand for squirrel.deleteBuilder
//
// Sets passed table and configured placeholder format
func (s Store) DeleteBuilder(table string) squirrel.DeleteBuilder {
	return squirrel.Delete(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

func (s Store) DB() dbLayer {
	return s.db
}

func (s Store) Config() *Config {
	return s.config
}

// column preprocessor logic to modify db value before using it in condition filter
//
// It checks registered ColumnPreprocessors from config
// and then the standard set
//
// No preprocessor ("") is intentionally checked after checking registered list of preprocessors
func (s Store) preprocessColumn(col string, p string) string {
	if fn, has := s.config.ColumnPreprocessors[p]; has {
		return fn(col, p)
	}

	switch p {
	case "":
		return col
	case "lower":
		return fmt.Sprintf("LOWER(%s)", col)
	default:
		panic(fmt.Sprintf("unknown preprocessor %q used for column %q", p, col))
	}
}

// value preprocessor logic to modify input value before using it in condition filters
//
// It checks registered ValuePreprocessors from config
// and then the standard set
//
// No preprocessor ("") is intentionally checked after checking registered list of preprocessors
func (s Store) preprocessValue(val interface{}, p string) interface{} {
	if fn, has := s.config.ValuePreprocessors[p]; has {
		return fn(val, p)
	}

	switch p {
	case "":
		return val
	case "lower":
		if str, ok := val.(string); ok {
			return strings.ToLower(str)
		}
		panic(fmt.Sprintf("preprocessor %q not compatible with type %T (value: %v)", p, val, val))

	default:
		panic(fmt.Sprintf("unknown preprocessor %q used for value %v", p, val))
	}
}

// SqlFunctionHandler calls configured sql function handler if set
// otherwise returns passed arguments directly
func (s Store) SqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	if s.config.SqlFunctionHandler == nil {
		return f, nil
	}

	return s.config.SqlFunctionHandler(f)
}

// FieldToColumnTypeCaster calls configured field type caster if set
// otherwise returns passed arguments directly
func (s Store) FieldToColumnTypeCaster(f ModuleFieldTypeDetector, i ql.Ident) (ql.Ident, error) {
	var err error
	i.Value, err = s.config.CastModuleFieldToColumnType(f, i.Value)
	return i, err
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
			return nil
		}

		if lastTaskErr = task(ctx, tx); lastTaskErr == nil {
			// Task completed successfully
			return tx.Commit()
		}

		// No matter the cause of the error, rollback the transaction
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction (tries: %d) on error %v: %w", try, lastTaskErr, rollbackErr)
		}

		// Call the the configured transaction retry error handler
		// if this particular error should be retried or not
		//
		// We can not generalize here because different store implementation have
		// different errors and we need to act accordingly
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

func setCursorCond(q squirrel.SelectBuilder, cursor *filter.PagingCursor) squirrel.SelectBuilder {
	if cursor != nil && len(cursor.Keys()) > 0 {
		const cursorTpl = `(%s) %s (?%s)`
		op := ">"
		if cursor.Reverse {
			op = "<"
		}

		pred := fmt.Sprintf(cursorTpl, strings.Join(cursor.Keys(), ", "), op, strings.Repeat(", ?", len(cursor.Keys())-1))
		q = q.Where(pred, cursor.Values()...)
	}

	return q
}

func setOrderBy(q squirrel.SelectBuilder, sort filter.SortExprSet, ss ...string) (squirrel.SelectBuilder, error) {
	var (
		sortable = slice.ToStringBoolMap(ss)
		sqlSort  = make([]string, len(sort))
	)
	for i, c := range sort {
		if sortable[c.Column] {
			sqlSort[i] = sort[i].Column
		} else {
			return q, fmt.Errorf("column %q is not sortable", c.Column)
		}

		if sort[i].Descending {
			sqlSort[i] += " DESC"
		}
	}

	return q.OrderBy(sqlSort...), nil
}

// TxNoRetry - Transaction retry handler
//
// Only returns false so transactions will never retry
func TxNoRetry(int, error) bool             { return false }
func ErrHandlerFallthrough(err error) error { return err }
