package pgsql

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

type (
	Store struct {
		*rdbms.Store
	}
)

func init() {
	store.Register(Connect, "postgresql", "postgres", "pgsql")
}

func Connect(ctx context.Context, dsn string) (store.Storable, error) {
	var (
		err error
		cfg *rdbms.Config

		s = new(Store)
	)

	if cfg, err = ProcDataSourceName(dsn); err != nil {
		return nil, err
	}

	cfg.PlaceholderFormat = squirrel.Dollar
	cfg.ErrorHandler = errorHandler
	cfg.SqlFunctionHandler = sqlFunctionHandler
	cfg.CastModuleFieldToColumnType = fieldToColumnTypeCaster

	if s.Store, err = rdbms.Connect(ctx, cfg); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Upgrade(ctx context.Context, log *zap.Logger) (err error) {
	if err = (&rdbms.Schema{}).Upgrade(ctx, NewUpgrader(log, s)); err != nil {
		return fmt.Errorf("can not upgrade postgresql schema: %w", err)
	}

	return nil
}

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
