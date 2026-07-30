package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/crusttech/human/server/store/adapters/rdbms/dal"

	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms"
	"github.com/crusttech/human/server/store/adapters/rdbms/instrumentation"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ngrok/sqlmw"
)

const (
	// base for our schemas
	SCHEMA = "postgres"

	// debug schema with verbose logging
	debugSchema = SCHEMA + "+debug"
)

func init() {
	store.Register(Connect, SCHEMA, debugSchema)
	sql.Register(debugSchema, sqlmw.Driver(new(pq.Driver), instrumentation.Debug()))
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

		DAL: dal.Connection(db, Dialect(), DataDefiner(cfg.DBName, db)),

		Dialect:      Dialect(),
		ErrorHandler: errorHandler,

		DataDefiner: DataDefiner(cfg.DBName, db),
		Ping:        db.PingContext,
	}

	s.SetDefaults()

	return s, nil
}

// NewConfig validates given DSN and ensures
// params are present and correct
func NewConfig(dsn string) (c *rdbms.ConnConfig, err error) {
	const (
		validScheme = "postgres"
	)
	var (
		scheme string
		u      *url.URL
	)

	if u, err = url.Parse(dsn); err != nil {
		return nil, err
	}

	if strings.HasPrefix(dsn, "postgres") {
		scheme = u.Scheme
		u.Scheme = validScheme
	} else {
		return nil, fmt.Errorf("expecting valid schema (postgres://) at the beginning of the DSN")
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

// firstNonEmptyStr returns the first non-empty value, so a unique-violation
// message degrades to a readable placeholder when the driver leaves the table
// or constraint name blank rather than printing an empty pair of quotes.
func firstNonEmptyStr(vv ...string) string {
	for _, v := range vv {
		if v != "" {
			return v
		}
	}

	return ""
}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(*pq.Error); ok {
			switch implErr.Code.Name() {
			case "unique_violation":
				// Name the table and constraint postgres actually rejected.
				// Without them this is indistinguishable from the service-level
				// check in checkXConstraints, and the two have entirely
				// different causes: a stale index the model no longer declares
				// vs a genuine duplicate.
				return store.ErrNotUniqueOn(
					firstNonEmptyStr(implErr.Table, "record"),
					firstNonEmptyStr(implErr.Constraint, "a unique constraint"),
				).Wrap(implErr)
			}
		}
	}

	return err
}
