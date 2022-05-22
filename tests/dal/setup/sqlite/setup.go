package sqlite

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	rdbmsDAL "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	"github.com/jmoiron/sqlx"
)

func Setup(dsn string) (_ dal.Connection, err error) {
	var (
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		cfg *rdbms.ConnConfig
		db  *sqlx.DB
	)

	cfg, err = sqlite.NewConfig(dsn)
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

	return rdbmsDAL.Connection(db, sqlite.Dialect()), nil
}

// remove when store support for table creation is added to CRS
//
// When support for creating DDL commands (creating tables) from DAL models and attributes
// is added, this can be removed!
func tables(ctx context.Context, db sqlx.ExecerContext) (err error) {
	return ddl.Exec(ctx, db,
		`DROP TABLE IF EXISTS crs_test_codec`,
		`CREATE TABLE IF NOT EXISTS crs_test_codec (
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
		)`,

		`DROP TABLE IF EXISTS crs_test_search`,
		`CREATE TABLE IF NOT EXISTS crs_test_search (
			id INT NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			meta JSON,
			p_string TEXT,
			p_number NUMERIC,
			p_is_odd BOOLEAN,
			
			PRIMARY KEY(id )
		)`,
	)
}
