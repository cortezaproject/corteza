package mysql

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	rdbmsdal "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/jmoiron/sqlx"
)

func init() {
	dal.Register(dalConnector, baseSchema, debugSchema)
}

func dalConnector(ctx context.Context, dsn string, cc ...capabilities.Capability) (_ dal.StoreConnection, err error) {
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

	return rdbmsdal.Connection(db, Dialect()), nil
}
