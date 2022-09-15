package mysql

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	rdbmsdal "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
)

func init() {
	dal.RegisterConnector(dalConnector, SCHEMA, debugSchema)
}

func dalConnector(ctx context.Context, dsn string, cc ...dal.Operation) (_ dal.Connection, err error) {
	cfg, err := NewConfig(dsn)
	if err != nil {
		return
	}

	// @todo rework the config building a bit; this will do for now
	if cfg.ConnTryMax >= 99 {
		cfg.ConnTryMax = 2
	}

	db, err := connectBase(ctx, cfg)

	if err != nil {
		return
	}
	return rdbmsdal.Connection(db, Dialect(), DataDefiner(cfg.DBName, db), cc...), nil
}
