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
	conn struct {
		label string
		dsn   string

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

	// callback funnction used in eachDB() test utility function
	connFn func(*testing.T, *conn) error
)

const (
	// connection environment variable prefix
	envVarPrefix = "TEST_STORE_ADAPTERS_RDBMS_"
)

var (
	// all connections, preloaded
	connections []*conn
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func TestMain(m *testing.M) {
	var (
		err error
	)

	connections, err = getConnections()
	cli.HandleError(err)

	os.Exit(m.Run())
}

// eachConnection calls connFun with a new connection to each configured
// test database
func eachDB(t *testing.T, fn connFn) {
	for _, c := range connections {
		t.Run(c.label, func(t *testing.T) {
			if err := fn(t, c); err != nil {
				t.Fatalf("%v", err)
			}
		})
	}
}

func getConnections() (configs []*conn, err error) {
	var (
		log = logger.Default()
		ctx = context.Background()

		s store.Storer

		is bool
	)

	for _, pair := range os.Environ() {
		if !strings.HasPrefix(pair, envVarPrefix) {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}

		configs = append(configs, &conn{
			label: kv[0][len(envVarPrefix):],
			dsn:   kv[1],
		})
	}

	if len(configs) == 0 {
		// append in-memory sqlite3 to have at least on implementation to test
		configs = append(configs, &conn{
			label: "sqlite-in-memory",
			dsn:   "sqlite3://file::memory:?cache=shared&mode=memory",
		})
	}

	for _, c := range configs {
		err = func() error {
			// check each dsn by prefix and use newConnConfig from
			// appropriate driver
			switch {
			case strings.HasPrefix(c.dsn, sqlite.SCHEMA):
				c.config, err = sqlite.NewConfig(c.dsn)
				c.dialect = sqlite.Dialect()
				c.isSQLite = true
			case strings.HasPrefix(c.dsn, postgres.SCHEMA):
				c.config, err = postgres.NewConfig(c.dsn)
				c.dialect = postgres.Dialect()
				c.isPostgres = true
			case strings.HasPrefix(c.dsn, mysql.SCHEMA):
				c.config, err = mysql.NewConfig(c.dsn)
				c.dialect = mysql.Dialect()
				c.isMySQL = true
			default:
				return fmt.Errorf("unsupported DB: %s", c.dsn)
			}

			if err != nil {
				return fmt.Errorf("can not create config from %s: %w", c.dsn, err)
			}

			// need to connect in under 1 second
			toCtx, cfn := context.WithTimeout(ctx, time.Second*5)
			defer cfn()

			s, err = store.Connect(toCtx, log, c.dsn, true)
			if err != nil {
				return fmt.Errorf("can not connect to %q: %w", c.config.DriverName, err)
			}

			if c.store, is = s.(*rdbms.Store); !is {
				return fmt.Errorf("can not extract *store.RDBMS from store")
			}

			if c.db, is = c.store.DB.(*sqlx.DB); !is {
				return fmt.Errorf("can not extract *sqlx.DB from store")
			}

			return nil
		}()

		if err != nil {
			return nil, err
		}
	}

	return
}

func exec(c *conn, q string) error {
	_, err := c.db.ExecContext(context.Background(), q)
	return err

}
