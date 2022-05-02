package mysql

import (
	"context"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/crs/test"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/crs"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func connect(req *require.Assertions) (db *sqlx.DB) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		cfg *rdbms.ConnConfig
	)

	cfg, err = NewConfig("mysql+debug://crust:crust@tcp(localhost:3306)/crust_store_test?collation=utf8mb4_general_ci&charset=utf8mb4")
	req.NoError(err)

	db, err = rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
	req.NoError(err)

	return
}

func TestMySQL(t *testing.T) {
	var (
		req = require.New(t)
		db  = connect(req)
	)

	defer db.Close()

	setupCodecTest(db, req)
	setupRecordSearchTest(db, req)
	test.All(t, crs.Connection(db, &dialect{}))
}

// remove when store support for table creation is added to CRS
//
// when support is added, test.TestRecordCodec
func setupCodecTest(db sqlx.ExecerContext, req *require.Assertions) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
	)

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS crs_test_codec`)
	req.NoError(err)

	tblCreate := `
		CREATE TABLE IF NOT EXISTS crs_test_codec (
			"id" BIGINT UNSIGNED NOT NULL,
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
		)`

	_, err = db.ExecContext(ctx, strings.ReplaceAll(tblCreate, "\"", "`"))
	req.NoError(err)
}

// remove when store support for table creation is added to CRS
func setupRecordSearchTest(db sqlx.ExecerContext, req *require.Assertions) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
	)

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS crs_test_search`)
	req.NoError(err)
	tblCreate := `
	CREATE TABLE IF NOT EXISTS crs_test_search (
		"id" INT NOT NULL,
		"created_at" DATETIME,
		"updated_at" DATETIME,
		"meta" JSON,
		"p_string" TEXT,
		"p_number" NUMERIC,
		"p_is_odd" BOOLEAN,
		
		PRIMARY KEY("id")
	)
	`

	_, err = db.ExecContext(ctx, strings.ReplaceAll(tblCreate, "\"", "`"))
	req.NoError(err)
}
