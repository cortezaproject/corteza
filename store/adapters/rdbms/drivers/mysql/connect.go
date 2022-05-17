package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/instrumentation"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ngrok/sqlmw"
)

const (
	// base for our schemas
	baseSchema = "mysql"

	// debug schema with verbose logging
	debugSchema = baseSchema + "+debug"
)

func init() {
	store.Register(Connect, baseSchema, debugSchema)
	sql.Register(debugSchema, sqlmw.Driver(new(mysql.MySQLDriver), instrumentation.Debug()))
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

	if err = connSetup(ctx, db); err != nil {
		return
	}

	s := &rdbms.Store{
		DB: db,

		Dialect:           goquDialectWrapper,
		TxRetryErrHandler: txRetryErrHandler,
		ErrorHandler:      errorHandler,

		SchemaAPI: &schema{dbName: cfg.DBName, dialect: Dialect()},
	}

	s.SetDefaults()

	return s, nil
}

// NewConfig validates given DSN and ensures
// params are present and correct
//
// If param is missing it sets it to default but returns
// error in case of incorrect param value
//
// See https://github.com/go-sql-driver/mysql for available dsn params
//
func NewConfig(in string) (*rdbms.ConnConfig, error) {
	const (
		schemeDel   = "://"
		validScheme = "mysql"
	)

	var (
		endOfSchema = strings.Index(in, schemeDel)

		c = &rdbms.ConnConfig{}
	)

	if endOfSchema > 0 && (in[:endOfSchema] == validScheme || strings.HasPrefix(in[:endOfSchema], validScheme+"+")) {
		c.DriverName = in[:endOfSchema]
		c.DataSourceName = in[endOfSchema+len(schemeDel):]
	} else {
		return nil, fmt.Errorf("expecting valid schema (mysql://) at the beginning of the DSN (%s)", in)
	}

	if pdsn, err := mysql.ParseDSN(c.DataSourceName); err != nil {
		return nil, err
	} else {
		// Ensure driver parses time
		pdsn.ParseTime = true

		if pdsn.Collation == "" {
			pdsn.Collation = "utf8mb4_general_ci"
		}

		if pdsn.Params == nil {
			pdsn.Params = make(map[string]string)
		}

		if pdsn.Params["charset"] == "" {
			pdsn.Params["charset"] = "utf8mb4"
		}

		c.DataSourceName = pdsn.FormatDSN()
		c.DBName = pdsn.DBName
	}

	c.SetDefaults()

	return c, nil
}

// Connection setup
func connSetup(ctx context.Context, db sqlx.ExecerContext) (err error) {
	// See https://dev.mysql.com/doc/refman/8.0/en/sql-mode.html#sqlmode_ansi for details
	if _, err = db.ExecContext(ctx, `SET SESSION sql_mode = 'ANSI'`); err != nil {
		return
	}

	return
}

func txRetryErrHandler(try int, err error) bool {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	var mysqlErr, ok = err.(*mysql.MySQLError)
	if !ok || mysqlErr == nil {
		return false
	}

	switch mysqlErr.Number {
	case
		1205, // lock within transaction
		1206, // total number of locks exceeded the lock table size
		1208, // not allowed while thread is holding global read lock
		1213: // deadlock found
		return true

	}

	return false
}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(*mysql.MySQLError); ok {
			// https://www.fromdual.com/de/mysql-error-codes-and-messages
			switch implErr.Number {
			case 1062: // Can't write, because of unique constraint, to table '%s'
				return store.ErrNotUnique.Wrap(implErr)
			}
		}
	}

	return err
}
