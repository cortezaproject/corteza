package sqlite

import (
	"context"
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

	cfg, err = NewConfig("sqlite3+debug://file::memory:?cache=shared&mode=memory")
	//cfg, err = NewConfig("sqlite3://file:/tmp/test.db")
	req.NoError(err)
	db, err = rdbms.Connect(ctx, logger.MakeDebugLogger(), cfg)
	req.NoError(err)

	return
}

func TestSQLite(t *testing.T) {
	var (
		req = require.New(t)
		db  = connect(req)
	)

	defer db.Close()

	setupCodecTest(db, req)
	setupRecordSearchTest(db, req)
	test.All(t, crs.Connection(db, Dialect()))
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
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS crs_test_codec (
		id UNSIGNED BIG INT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		meta JSON,
		pID UNSIGNED BIG INT,
		pRef UNSIGNED BIG INT,
		pTimestamp_TZT TIMESTAMP,
		pTimestamp_TZF TIMESTAMP,
		pTime TIME,
		pDate DATE,
		pNumber NUMERIC,
		pText TEXT,
		pBoolean_T BOOLEAN,
		pBoolean_F BOOLEAN,
		pEnum TEXT,
		pGeometry TEXT,
		pJSON TEXT,
		pBlob BLOB,
		pUUID UUID,
		
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

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS crs_test_search`)
	req.NoError(err)
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS crs_test_search (
		id INT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		meta JSON,
		p_string TEXT,
		p_number NUMERIC,
		p_is_odd BOOLEAN,
		
		PRIMARY KEY(id )
	)
	`)
	req.NoError(err)
}
