package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store/bulk"
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

	schemaUpgradeGenerator interface {
		TableExists(string) bool
		CreateTable(t *ddl.Table) string
		CreateIndexes(ii ...*ddl.Index) []string
	}

	rowScanner interface {
		Scan(...interface{}) error
	}

	Config struct {
		DriverName     string
		DataSourceName string
		DBName         string

		PlaceholderFormat squirrel.PlaceholderFormat

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

		// Implementations can override internal RDBMS row scanners
		RowScanners map[string]interface{}
	}

	Store struct {
		config *Config

		// Schema upgrade generator converts internal upgrade config
		// to implementation specific SQL
		sug schemaUpgradeGenerator

		db *sqlx.DB
	}
)

const (
	// This is the absolute maximum retries we'll allow
	TxRetryHardLimit = 100

	DefaultSliceCapacity = 1000
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

//func Instrumentation(log *zap.Logger) {
//	logger := instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
//		//spew.Dump(msg, keyvals)
//		log.With(zap.Any("kv", keyvals)).Info(msg)
//	})
//
//	sql.Register(
//		"mysql+instrumented",
//		instrumentedsql.WrapDriver(&mysql.MySQLDriver{}, instrumentedsql.WithLogger(logger)))
//
//	sql.Register(
//		"postgres+instrumented",
//		instrumentedsql.WrapDriver(&pq.Driver{}, instrumentedsql.WithLogger(logger)))
//}

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

	if err := s.Connect(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, s.config.DriverName, s.config.DataSourceName)
	return err
}

// Select is a shorthand for squirrel.SelectBuilder
//
// Sets passed table & columns and configured placeholder format
func (s Store) Select(table string, cc ...string) squirrel.SelectBuilder {
	return squirrel.Select(cc...).From(table).PlaceholderFormat(s.config.PlaceholderFormat)
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

// Insert is a shorthand for squirrel.InsertBuilder
//
// Sets passed table and configured placeholder format
func (s Store) Insert(table string) squirrel.InsertBuilder {
	return squirrel.Insert(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

// Update is a shorthand for squirrel.UpdateBuilder
//
// Sets passed table and configured placeholder format
func (s Store) Update(table string) squirrel.UpdateBuilder {
	return squirrel.Update(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

// Delete is a shorthand for squirrel.DeleteBuilder
//
// Sets passed table and configured placeholder format
func (s Store) Delete(table string) squirrel.DeleteBuilder {
	return squirrel.Delete(table).PlaceholderFormat(s.config.PlaceholderFormat)
}

func (s Store) DB() *sqlx.DB {
	return s.db
}

func (s Store) Config() *Config {
	return s.config
}

// Bulk returns channel that accepts jobs and executes them inside a transaction
//
// Note: This is experimental function!
// Final version might not return channel directly
//
func (s Store) Bulk(ctx context.Context) chan bulk.Job {
	jc := make(chan bulk.Job)

	go Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		var job bulk.Job

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case job = <-jc:
				if err = job.Do(ctx, s); err != nil {
					return
				}
			}
		}
	})

	return jc
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

func ExecuteSqlizer(ctx context.Context, db sqlx.ExecerContext, sqlizer squirrel.Sqlizer) error {
	query, args, err := sqlizer.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, query, args...)

	return err
}

func Truncate(ctx context.Context, db sqlx.ExecerContext, table string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM "+table)
	return err
}

// Tx begins a new db transaction and handles it's retries when possible
//
// It utilizes configured transaction error handlers and max-retry limits
// to determine if and how many times transaction should be retried
func Tx(ctx context.Context, db *sqlx.DB, cfg *Config, txOpt *sql.TxOptions, task func(*sqlx.Tx) error) error {
	var (
		lastTaskErr error
		err         error
		tx          *sqlx.Tx
		try         = 1
	)

	for {
		try++

		// Start transaction
		tx, err = db.BeginTxx(ctx, txOpt)
		if err != nil {
			return nil
		}

		if lastTaskErr = task(tx); lastTaskErr == nil {
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

		// Tx error handlers can take current number of tries into account and
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
func TxNoRetry(int, error) bool { return false }
