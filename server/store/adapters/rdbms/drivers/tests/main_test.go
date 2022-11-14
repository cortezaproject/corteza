package tests

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/mysql"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/postgres"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"testing"
	"time"
)

type (
	// test utility struct for testing various RDBMS implementations
	connectionInfo struct {
		dsn string

		config  *rdbms.ConnConfig
		store   *rdbms.Store
		db      *sqlx.DB
		dialect drivers.Dialect

		// used for creating tables
		// @todo remove when we have proper DDL abstraction
		isSQLite   bool
		isMySQL    bool
		isPostgres bool
	}
)

const (
	DB_DSN_DEFAULT = "sqlite3://file::memory:?cache=shared&mode=memory"
)

var (
	// all connections, preloaded
	conn *connectionInfo

	log = logger.MakeDebugLogger()

	// all tests should use this context to pass
	// into DB calls to ensure commands are logged
	ctx = logger.ContextWithValue(context.Background(), log)
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func TestMain(m *testing.M) {
	cli.HandleError(connect())
	os.Exit(m.Run())
}

func connect() (err error) {
	var (
		is bool

		s store.Storer
	)

	conn = &connectionInfo{}

	if conn.dsn, is = os.LookupEnv("DB_DSN"); !is {
		conn.dsn = DB_DSN_DEFAULT
	}

	switch {
	case strings.HasPrefix(conn.dsn, sqlite.SCHEMA):
		conn.config, err = sqlite.NewConfig(conn.dsn)
		conn.dialect = sqlite.Dialect()
		conn.isSQLite = true
	case strings.HasPrefix(conn.dsn, postgres.SCHEMA):
		conn.config, err = postgres.NewConfig(conn.dsn)
		conn.dialect = postgres.Dialect()
		conn.isPostgres = true
	case strings.HasPrefix(conn.dsn, mysql.SCHEMA):
		conn.config, err = mysql.NewConfig(conn.dsn)
		conn.dialect = mysql.Dialect()
		conn.isMySQL = true
	default:
		return fmt.Errorf("unsupported DB (dns: %q)", conn.dsn)
	}

	// check each dsn by prefix and use newConnConfig from
	// appropriate driver

	if err != nil {
		return fmt.Errorf("can not create config from %s: %w", conn.dsn, err)
	}

	// need to connect in under 1 second
	toCtx, cfn := context.WithTimeout(ctx, time.Second*5)
	defer cfn()

	s, err = store.Connect(toCtx, log, conn.dsn, true)
	if err != nil {
		return fmt.Errorf("can not connect to %q: %w", conn.config.DriverName, err)
	}

	if conn.store, is = s.(*rdbms.Store); !is {
		return fmt.Errorf("can not extract *store.RDBMS from store")
	}

	if conn.db, is = conn.store.DB.(*sqlx.DB); !is {
		return fmt.Errorf("can not extract *sqlx.DB from store")
	}

	log.Debug("connected to " + conn.dsn)

	return nil
}

func exec(q string) error {
	_, err := conn.db.ExecContext(ctx, q)
	return err

}
