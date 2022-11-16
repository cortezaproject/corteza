package cockroach

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/rdbms"
	"github.com/cortezaproject/corteza/server/store/rdbms/instrumentation"
	"github.com/lib/pq"
	"github.com/ngrok/sqlmw"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

type (
	Store struct {
		*rdbms.Store
	}
)

const (
	storeTypeName      = "cockroachdb"
	storeTypeWithDebug = "cockroachdb+debug"
)

// For now, CockroachDB implementation is disabled due to some limitation from cockroachDB!
//
//	Here is why could not work on this further:
//	* Prefix/Expression Indexes.
//		I.e. we are not able to use "lower" expression while creating/updating indexes
//	* Not able to use Prefix/Expression into upsert.
//		I.e. also not able to use "lower"(for most case) expression in upsert
//
// @note we could enable this once cockroachDB overcome above limitations :)
func init() {
	store.Register(Connect, storeTypeName, storeTypeWithDebug)
	sql.Register(storeTypeWithDebug, sqlmw.Driver(new(pq.Driver), instrumentation.Debug()))
}

func Connect(ctx context.Context, dsn string) (store.Storer, error) {
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

	ql.QueryEncoder = &QueryEncoder{}

	return s, nil
}

func (s *Store) Upgrade(ctx context.Context, log *zap.Logger) (err error) {
	if err = (&rdbms.Schema{}).Upgrade(ctx, NewUpgrader(log, s)); err != nil {
		return fmt.Errorf("can not upgrade cockroach schema: %w", err)
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

	if strings.HasPrefix(dsn, storeTypeName) {
		scheme = validScheme
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
