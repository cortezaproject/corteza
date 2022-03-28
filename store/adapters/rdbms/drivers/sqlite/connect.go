package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/instrumentation"
	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/ngrok/sqlmw"
)

func init() {
	store.Register(Connect, "sqlite3", "sqlite3+debug")
	sql.Register("sqlite3+debug", sqlmw.Driver(new(sqlite3.SQLiteDriver), instrumentation.Debug()))
}

func Connect(ctx context.Context, dsn string) (_ store.Storer, err error) {
	var (
		cfg *rdbms.Config
		s   *rdbms.Store
	)
	if cfg, err = ProcDataSourceName(dsn); err != nil {
		return
	}

	//cfg.PlaceholderFormat = squirrel.Dollar
	cfg.TxRetryErrHandler = txRetryErrHandler
	cfg.ErrorHandler = errorHandler

	// Using transactions in SQLite causes table locking
	// @todo there must be a better way to go around this
	cfg.TxDisabled = true
	//cfg.SqlFunctionHandler = sqlFunctionHandler
	//cfg.ASTFormatter = sqlASTFormatter
	//cfg.CastModuleFieldToColumnType = fieldToColumnTypeCaster

	// Set to zero
	// Otherwise SQLite (in-memory) disconnects
	// and all tables and data is lost
	cfg.ConnMaxLifetime = 0

	if s, err = rdbms.Connect(ctx, cfg); err != nil {
		return
	}

	cfg.Upgrader = NewUpgrader(s)

	//ql.QueryEncoder = &QueryEncoder{}

	return s, nil
}

func ConnectInMemory(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, "sqlite3://file::memory:?cache=shared&mode=memory")
}

func ConnectInMemoryWithDebug(ctx context.Context) (s store.Storer, err error) {
	return Connect(ctx, "sqlite3+debug://file::memory:?cache=shared&mode=memory")
}

//func (s *Store) Upgrade(ctx context.Context, log *zap.Logger) (err error) {
//	if err = (&rdbms.Schema{}).Upgrade(ctx, NewUpgrader(log, s)); err != nil {
//		return fmt.Errorf("cannot upgrade sqlite schema: %w", err)
//	}
//
//	return nil
//}

// ProcDataSourceName validates given DSN and ensures
// params are present and correct
func ProcDataSourceName(in string) (*rdbms.Config, error) {
	const (
		schemeDel   = "://"
		validScheme = "sqlite3"
	)

	var (
		endOfSchema = strings.Index(in, schemeDel)

		c = &rdbms.Config{
			// note that we're not using vanilla SQLite
			// see dialect.go
			Dialect: goqu.Dialect("sqlite3"),
		}
	)

	if endOfSchema > 0 && (in[:endOfSchema] == validScheme || strings.HasPrefix(in[:endOfSchema], validScheme+"+")) {
		c.DriverName = in[:endOfSchema]
		c.DataSourceName = in[endOfSchema+len(schemeDel):]
	} else {
		return nil, fmt.Errorf("expecting valid schema (sqlite3://) at the beginning of the DSN (%s)", in)
	}

	return c, nil
}

func txRetryErrHandler(try int, err error) bool {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	var sqliteErr, ok = err.(sqlite3.Error)
	if !ok {
		return false
	}

	switch sqliteErr.Code {
	case sqlite3.ErrLocked:
		return true

	}

	return false
}

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
