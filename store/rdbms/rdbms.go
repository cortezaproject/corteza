package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/ddl"
	"github.com/jmoiron/sqlx"
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
	txRetryOnErrHandler func(int, error) bool
	columnPreprocFn     func(string, string) string
	valuePreprocFn      func(interface{}, string) interface{}
	errorHandler        func(error) error
	triggerKey          string

	schemaUpgradeGenerator interface {
		TableExists(string) bool
		CreateTable(t *ddl.Table) string
		CreateIndexes(ii ...*ddl.Index) []string
	}

	rowScanner interface {
		Scan(...interface{}) error
	}

	TriggerHandlers map[triggerKey]interface{}

	Config struct {
		DriverName     string
		DataSourceName string
		DBName         string

		PlaceholderFormat squirrel.PlaceholderFormat

		// These 3 are passed directly to connection
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
		MaxIdleConns    int

		// Disable transactions
		TxDisabled bool

		// How many times should we retry failed transaction?
		TxMaxRetries int

		// TxRetryErrHandler should return true if transaction should be retried
		//
		// Because retry algorithm varies between concrete rdbms implementations
		//
		// Handler must return true if failed transaction should be replied
		// and false if we're safe to terminate it
		TxRetryErrHandler txRetryOnErrHandler

		ColumnPreprocessors map[string]columnPreprocFn
		ValuePreprocessors  map[string]valuePreprocFn

		ErrorHandler errorHandler

		// Implementations can override internal RDBMS row scanners
		RowScanners map[string]interface{}

		// Different store backend implementation might handle upsert differently...
		UpsertBuilder func(*Config, string, store.Payload, ...string) (squirrel.InsertBuilder, error)

		// TriggerHandlers handle various exceptions that can not be handled generally within RDBMS package.
		// see triggerKey type and defined constants to see where the hooks are and how can they be called
		TriggerHandlers TriggerHandlers

		// UniqueConstraintCheck flag controls if unique constraints should be explicitly checked within
		// store or is this handled inside the storage
		//
		//
		UniqueConstraintCheck bool

		// FunctionHandler takes care of translation & transformation of (sql) functions
		// and their parameters
		//
		// Functions are used in filters and aggregations
		SqlFunctionHandler func(f ql.Function) (ql.ASTNode, error)

		CastModuleFieldToColumnType func(field ModuleFieldTypeDetector, ident ql.Ident) (ql.Ident, error)
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

const (
	// This is the absolute maximum retries we'll allow
	TxRetryHardLimit = 100

	DefaultSliceCapacity = 1000

	MinRefetchLimit = 10
	MaxRefetches    = 100
)

func New(ctx context.Context, cfg *Config) (*Store, error) {
	var s = &Store{
		config: cfg,
	}

	if s.config.PlaceholderFormat == nil {
		s.config.PlaceholderFormat = squirrel.Question
	}

	if s.config.TxMaxRetries == 0 {
		s.config.TxMaxRetries = TxRetryHardLimit
	}

	if s.config.TxRetryErrHandler == nil {
		// Default transaction retry handler
		s.config.TxRetryErrHandler = TxNoRetry
	}

	if s.config.ErrorHandler == nil {
		s.config.ErrorHandler = ErrHandlerFallthrough
	}

	if s.config.UpsertBuilder == nil {
		s.config.UpsertBuilder = UpsertBuilder
	}

	if s.config.MaxIdleConns == 0 {
		// Same as default in the db/sql
		s.config.MaxIdleConns = 2
	}

	if s.config.TriggerHandlers == nil {
		s.config.TriggerHandlers = TriggerHandlers{}
	}

	if err := s.Connect(ctx); err != nil {
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

func (s *Store) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.config.DriverName, s.config.DataSourceName)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(s.config.MaxOpenConns)
	db.SetConnMaxLifetime(s.config.ConnMaxLifetime)
	db.SetMaxIdleConns(s.config.MaxIdleConns)

	s.db = db
	return err
}

func (s Store) Query(ctx context.Context, q squirrel.SelectBuilder) (*sql.Rows, error) {
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	return s.db.QueryContext(ctx, query, args...)
}

// QueryRow returns row instead of filling in the passed struct
func (s Store) QueryRow(ctx context.Context, q squirrel.SelectBuilder) (*sql.Row, error) {
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	return s.db.QueryRowContext(ctx, query, args...), nil
}

func (s Store) Exec(ctx context.Context, sqlizer squirrel.Sqlizer) error {
	query, args, err := sqlizer.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

func (s Store) Tx(ctx context.Context, fn func(context.Context, store.Storable) error) error {
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
	return s.config.CastModuleFieldToColumnType(f, i)
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

// TxNoRetry - Transaction retry handler
//
// Only returns false so transactions will never retry
func TxNoRetry(int, error) bool             { return false }
func ErrHandlerFallthrough(err error) error { return err }
