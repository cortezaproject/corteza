package mssql

import (
	"context"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/dal"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"
	mssql "github.com/denisenkom/go-mssqldb"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// base for our schemas
	SCHEMA = "sqlserver"

	// debug schema with verbose logging
	// @todo this won't connect yet
	// debugSchema = SCHEMA + "+debug"
)

func init() {
	// store.Register(Connect, SCHEMA, debugSchema)
	// sql.Register(debugSchema, sqlmw.Driver(new(mssql.Driver), instrumentation.Debug()))
	store.Register(Connect, SCHEMA)
}

func Connect(ctx context.Context, dsn string) (_ store.Storer, err error) {
	cfg, err := NewConfig(dsn)
	if err != nil {
		return
	}
	db, err := connectBase(ctx, cfg)
	if err != nil {
		return
	}

	s := &rdbms.Store{
		DB: db,

		DAL: dal.Connection(db, Dialect(), DataDefiner(cfg.DBName, db)),

		Dialect: Dialect(),
		// TxRetryErrHandler: txRetryErrHandler,
		ErrorHandler: errorHandler,

		DataDefiner: DataDefiner(cfg.DBName, db),
		Ping:        db.PingContext,
	}

	s.SetDefaults()

	return s, nil
}

func connectBase(ctx context.Context, cfg *rdbms.ConnConfig) (db *sqlx.DB, err error) {
	if db, err = rdbms.Connect(ctx, logger.Default(), cfg); err != nil {
		return
	}

	// See https://stackoverflow.com/questions/2901453/sql-standard-to-escape-column-names
	if _, err = db.ExecContext(ctx, `SET QUOTED_IDENTIFIER ON`); err != nil {
		return
	}

	return
}

// NewConfig validates given DSN and ensures
// params are present and correct
//
// If param is missing it sets it to default but returns
// error in case of incorrect param value
//
// See https://github.com/go-sql-driver/mysql for available dsn params
//
func NewConfig(dsn string) (c *rdbms.ConnConfig, err error) {
	const (
		validScheme = "sqlserver"
	)
	var (
		scheme string
		u      *url.URL
	)
	c = &rdbms.ConnConfig{DriverName: scheme}

	if u, err = url.Parse(dsn); err != nil {
		return nil, err
	}

	if strings.HasPrefix(dsn, "sqlserver") {
		scheme = u.Scheme
		u.Scheme = validScheme
	}

	c = &rdbms.ConnConfig{
		DriverName:     scheme,
		DataSourceName: u.String(),
		DBName:         strings.Trim(u.Path, "/"),
		MaskedDSN:      u.Redacted(),
	}

	c.SetDefaults()

	return c, nil
}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(mssql.Error); ok {
			// https://learn.microsoft.com/en-us/sql/relational-databases/errors-events/database-engine-events-and-errors?view=sql-server-ver16
			switch implErr.Number {
			case 987, 2627: // Can't write, because of unique constraint, to table '%s'
				return store.ErrNotUnique.Wrap(implErr)
			}
		}
	}

	return err
}
