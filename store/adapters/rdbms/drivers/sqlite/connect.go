package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/instrumentation"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/ngrok/sqlmw"
)

func init() {
	store.Register(Connect, "sqlite3", "sqlite3+debug")
	sql.Register("sqlite3+debug", sqlmw.Driver(new(sqlite3.SQLiteDriver), instrumentation.Debug()))
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

		Dialect:      goqu.Dialect("sqlite3"),
		ErrorHandler: errorHandler,

		TxRetryLimit: -1,
		//TxRetryErrHandler: txRetryErrHandler,

		SchemaAPI: &schema{},
	}

	s.SetDefaults()

	return s, nil
}

func ConnectInMemory(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, "sqlite3://file::memory:?cache=shared&mode=memory")
}

func ConnectInMemoryWithDebug(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, "sqlite3+debug://file::memory:?cache=shared&mode=memory")
}

// NewConfig validates given DSN and ensures
// params are present and correct
func NewConfig(in string) (*rdbms.ConnConfig, error) {
	const (
		schemeDel   = "://"
		validScheme = "sqlite3"
	)

	var (
		endOfSchema = strings.Index(in, schemeDel)

		cfg = &rdbms.ConnConfig{}
	)

	if endOfSchema > 0 && (in[:endOfSchema] == validScheme || strings.HasPrefix(in[:endOfSchema], validScheme+"+")) {
		cfg.DriverName = in[:endOfSchema]
		cfg.DataSourceName = in[endOfSchema+len(schemeDel):]
	} else {
		return nil, fmt.Errorf("expecting valid schema (sqlite3://) at the beginning of the DSN (%s)", in)
	}

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
