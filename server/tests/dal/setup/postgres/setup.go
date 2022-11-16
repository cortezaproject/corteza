package postgres

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/postgres"
	"github.com/jmoiron/sqlx"
)

func Setup(ctx context.Context, dsn string) (_ *sqlx.DB, err error) {
	var (
		cfg *rdbms.ConnConfig
	)

	cfg, err = postgres.NewConfig(dsn)
	if err != nil {
		return
	}

	return rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
}
