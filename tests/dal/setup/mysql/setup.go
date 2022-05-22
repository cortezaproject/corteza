package mysql

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	rdbmsDAL "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

func Setup(dsn string) (_ dal.Connection, err error) {
	var (
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		cfg *rdbms.ConnConfig
		db  *sqlx.DB
	)

	cfg, err = mysql.NewConfig(dsn)
	if err != nil {
		return
	}

	db, err = rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
	if err != nil {
		return
	}

	if err = tables(ctx, db); err != nil {
		return
	}

	return rdbmsDAL.Connection(db, mysql.Dialect()), nil
}

// remove when store support for table creation is added to CRS
//
// When support for creating DDL commands (creating tables) from DAL models and attributes
// is added, this can be removed!
func tables(ctx context.Context, db sqlx.ExecerContext) (err error) {
	return ddl.Exec(ctx, db,
		`DROP TABLE IF EXISTS crs_test_codec`,
		strings.ReplaceAll(
			`CREATE TABLE IF NOT EXISTS crs_test_codec (
			"id" BIGINT UNSIGNED,
			"created_at" DATETIME,
			"updated_at" DATETIME,
			"meta" JSON,
			"pID" BIGINT UNSIGNED,
			"pRef" BIGINT UNSIGNED,
			"pTimestamp_TZT" DATETIME,
			"pTimestamp_TZF" DATETIME,
			"pTime" TIME,
			"pDate" DATE,
			"pNumber" NUMERIC,
			"pText" TEXT,
			"pBoolean_T" BOOLEAN,
			"pBoolean_F" BOOLEAN,
			"pEnum" TEXT,
			"pGeometry" TEXT,
			"pJSON" TEXT,
			"pBlob" BLOB,
			"pUUID" VARCHAR(36),
			
			PRIMARY KEY("id")
		)`, "\"", "`"),
		`DROP TABLE IF EXISTS crs_test_search`,
		strings.ReplaceAll(
			`CREATE TABLE IF NOT EXISTS crs_test_search (
			"id" INT NOT NULL,
			"created_at" DATETIME,
			"updated_at" DATETIME,
			"meta" JSON,
			"p_string" TEXT,
			"p_number" NUMERIC,
			"p_is_odd" BOOLEAN,
			
			PRIMARY KEY("id")
		)`, "\"", "`"),
	)
}
