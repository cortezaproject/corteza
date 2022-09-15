package postgres

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	rdbmsdal "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/jmoiron/sqlx"
)

func init() {
	dal.RegisterConnector(dalConnector, SCHEMA, debugSchema)
}

func dalConnector(ctx context.Context, dsn string, cc ...dal.Operation) (_ dal.Connection, err error) {
	var (
		db  *sqlx.DB
		cfg *rdbms.ConnConfig
	)

	if cfg, err = NewConfig(dsn); err != nil {
		return
	}

	// @todo rework the config building a bit; this will do for now
	if cfg.ConnTryMax >= 99 {
		cfg.ConnTryMax = 2
	}

	if db, err = rdbms.Connect(ctx, logger.Default(), cfg); err != nil {
		return
	}

	return rdbmsdal.Connection(db, Dialect(), DataDefiner(cfg.DBName, db), cc...), nil
}
