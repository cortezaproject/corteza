package postgres

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/dal/test"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func connect(req *require.Assertions) (db *sqlx.DB) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		cfg *rdbms.ConnConfig
	)

	cfg, err = NewConfig("postgres+debug://darh@localhost:5432/corteza_2022_3?sslmode=disable&")
	req.NoError(err)

	db, err = rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
	req.NoError(err)

	return
}

func TestPostgres(t *testing.T) {
	var (
		req = require.New(t)
		db  = connect(req)
	)

	defer db.Close()

	setupCodecTest(db, req)
	setupRecordSearchTest(db, req)
	test.All(t, dal.Connection(db, Dialect()))
}

// remove when store support for table creation is added to CRS
//
// when support is added, test.TestRecordCodec
func setupCodecTest(db sqlx.ExecerContext, req *require.Assertions) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
	)

	_, err = db.ExecContext(ctx, `
	CREATE TEMPORARY TABLE IF NOT EXISTS crs_test_codec (
		"id" BIGINT NOT NULL,
		"created_at" TIMESTAMP NOT NULL,
		"updated_at" TIMESTAMP,
		"meta" JSON,
		"pID" BIGINT,
		"pRef" BIGINT,
		"pTimestamp_TZT" TIMESTAMPTZ,
		"pTimestamp_TZF" TIMESTAMP,
		"pTime" TIME,
		"pDate" DATE,
		"pNumber" NUMERIC,
		"pText" TEXT,
		"pBoolean_T" BOOLEAN,
		"pBoolean_F" BOOLEAN,
		"pEnum" TEXT,
		"pGeometry" TEXT,
		"pJSON" TEXT,
		"pBlob" BYTEA,
		"pUUID" UUID,
		
		PRIMARY KEY(id)
	)`)
	req.NoError(err)
}

// remove when store support for table creation is added to CRS
func setupRecordSearchTest(db sqlx.ExecerContext, req *require.Assertions) {
	var (
		err error
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
	)

	_, err = db.ExecContext(ctx, `
	CREATE TEMPORARY  TABLE IF NOT EXISTS crs_test_search (
		"id" BIGINT NOT NULL,
		"created_at" TIMESTAMPTZ NOT NULL,
		"updated_at" TIMESTAMPTZ,
		"meta" JSON,
		"p_string" TEXT,
		"p_number" NUMERIC,
		"p_is_odd" BOOLEAN,
		
		PRIMARY KEY(id )
	)
	`)
	req.NoError(err)
}
