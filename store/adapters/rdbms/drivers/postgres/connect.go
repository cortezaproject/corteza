package postgres

import (
	"context"
	"database/sql"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/instrumentation"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/lib/pq"
	"github.com/ngrok/sqlmw"
)

func init() {
	store.Register(Connect, "postgres", "postgres+debug")
	sql.Register("postgres+debug", sqlmw.Driver(new(pq.Driver), instrumentation.Debug()))
}

func Connect(ctx context.Context, dsn string) (_ store.Storer, err error) {
	var (
		cfg *rdbms.Config
		s   *rdbms.Store
	)

	if cfg, err = ProcDataSourceName(dsn); err != nil {
		return
	}

	cfg.ErrorHandler = errorHandler

	if s, err = rdbms.Connect(ctx, cfg); err != nil {
		return
	}

	cfg.Upgrader = NewUpgrader(s)
	return s, nil
}

//func (s *Store) Upgrade(ctx context.Context, log *zap.Logger) (err error) {
//	if err = (&rdbms.Schema{}).Upgrade(ctx, NewUpgrader(log, s)); err != nil {
//		return fmt.Errorf("cannot upgrade postgresql schema: %w", err)
//	}
//
//	return nil
//}

// ProcDataSourceName validates given DSN and ensures
// params are present and correct
func ProcDataSourceName(dsn string) (c *rdbms.Config, err error) {
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
	}

	return &rdbms.Config{
		DriverName:     scheme,
		DataSourceName: u.String(),
		DBName:         strings.Trim(u.Path, "/"),
		Dialect:        goqu.Dialect("postgres"),
	}, nil
}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(*pq.Error); ok {
			switch implErr.Code.Name() {
			case "unique_violation":
				return store.ErrNotUnique.Wrap(implErr)
			}
		}
	}

	return err
}
