package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	pkgdal "github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/instrumentation"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/ngrok/sqlmw"
)

const (
	// base for our schemas
	SCHEMA = "sqlite3"

	// alternative s hema with custom driver
	altSchema = SCHEMA + "+alt"

	// debug schema with verbose logging
	debugSchema = SCHEMA + "+debug"
)

var (
	// make a custom driver with REGEX support
	driver = &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) (err error) {
			// register regexp function and use Go's regexp fn
			if err = conn.RegisterFunc("regexp", regexp.MatchString, true); err != nil {
				return
			}

			return
		},
	}
)

func init() {
	// register alter driver
	sql.Register(altSchema, driver)
	// register drbug driver
	sql.Register(debugSchema, sqlmw.Driver(driver, instrumentation.Debug()))

	store.Register(Connect, SCHEMA, altSchema, debugSchema)
}

func Connect(ctx context.Context, dsn string) (_ store.Storer, err error) {
	var (
		db  *sqlx.DB
		cfg *rdbms.ConnConfig
	)

	if cfg, err = NewConfig(dsn); err != nil {
		return
	}

	if db, err = rdbms.Connect(ctx, logger.Default(), cfg); err != nil {
		return
	}

	s := &rdbms.Store{
		DB: db,

		DAL: dal.Connection(db, Dialect(), DataDefiner(cfg.DBName, db), pkgdal.FullOperations()...),

		Dialect:      goquDialectWrapper,
		ErrorHandler: errorHandler,

		TxRetryLimit: -1,
		//TxRetryErrHandler: txRetryErrHandler,

		//SchemaAPI: &schema{},
	}

	s.SetDefaults()

	return s, nil
}

func ConnectInMemory(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, SCHEMA+"://file::memory:?cache=shared&mode=memory")
}

func ConnectInMemoryWithDebug(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, debugSchema+"://file::memory:?cache=shared&mode=memory")
}

// NewConfig validates given DSN and ensures
// params are present and correct
func NewConfig(in string) (*rdbms.ConnConfig, error) {
	const (
		schemaDel = "://"
	)

	var (
		cfg = &rdbms.ConnConfig{
			DriverName: altSchema,
		}
	)

	switch {
	case strings.HasPrefix(in, SCHEMA+schemaDel), strings.HasPrefix(in, altSchema+schemaDel):
	// no special handlign
	case strings.HasPrefix(in, debugSchema+schemaDel):
		cfg.DriverName = debugSchema
	default:
		return nil, fmt.Errorf("expecting valid schema (sqlite3://) at the beginning of the DSN")
	}

	// reassemble DSN with base schema
	cfg.DataSourceName = in[strings.Index(in, schemaDel)+len(schemaDel):]

	// Set to zero
	// Otherwise SQLite (in-memory) disconnects
	// and all tables and data is lost
	cfg.ConnMaxLifetime = 0

	cfg.SetDefaults()

	return cfg, nil
}

// Transactions are disabled on SQLite
//func txRetryErrHandler(try int, err error) bool {
//	for errors.Unwrap(err) != nil {
//		err = errors.Unwrap(err)
//	}
//
//	var sqliteErr, ok = err.(sqlite3.Error)
//	if !ok {
//		return false
//	}
//
//	switch sqliteErr.Code {
//	case sqlite3.ErrLocked:
//		return true
//
//	}
//
//	return false
//}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(sqlite3.Error); ok {
			switch implErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return store.ErrNotUnique.Wrap(implErr)
			}
		}
	}

	return err
}
