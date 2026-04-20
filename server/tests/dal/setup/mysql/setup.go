package mysql

import (
	"context"

	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/store/adapters/rdbms"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/mysql"
	"github.com/jmoiron/sqlx"
)

func Setup(ctx context.Context, dsn string) (_ *sqlx.DB, err error) {
	var (
		cfg *rdbms.ConnConfig
	)

	cfg, err = mysql.NewConfig(dsn)
	if err != nil {
		return
	}

	return rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
}
