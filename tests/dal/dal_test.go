package dal

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/tests/dal/setup/mysql"
	"github.com/cortezaproject/corteza-server/tests/dal/setup/postgres"
	"github.com/cortezaproject/corteza-server/tests/dal/setup/sqlite"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"testing"
)

// Tests DAL functionalities for all supported databases
//
// Tests for drivers are executed if DSN env variable (DAL_TEST_DSN_<driver>) is present.
// Example: DAL_TEST_DSN_MYSQL, DAL_TEST_DSN_POSTGRES
//
// Multiple space delimited DSNs are supported.
//
// Tests scans current and two parent folders for presence of .env file
// and loads the first one found.
//
//
func TestDAL(t *testing.T) {
	helpers.RecursiveDotEnvLoad()

	var (
		// enrich context with debug logger
		//
		// this will enable us to log driver commands
		// when +debug is used on DSN schema
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		conn dal.Connection
		err  error

		drivers = []struct {
			name    string
			dsn     string
			connect func(string) (dal.Connection, error)
		}{
			{
				name:    "sqlite",
				dsn:     "sqlite3+debug://file::memory:?cache=shared&mode=memory",
				connect: sqlite.Setup,
			},
			{
				name:    "mysql",
				dsn:     os.Getenv("DAL_TEST_DSN_MYSQL"),
				connect: mysql.Setup,
			},
			{
				name:    "postgres",
				dsn:     os.Getenv("DAL_TEST_DSN_POSTGRES"),
				connect: postgres.Setup,
			},
		}
	)

	for _, driver := range drivers {
		t.Run(driver.name, func(t *testing.T) {
			if driver.dsn == "" {
				t.Skip("DSN for DAL test not set")
			}

			for _, dsn := range strings.Split(driver.dsn, " ") {
				t.Run("", func(t *testing.T) {
					t.Log("Connecting to ", dsn)
					if conn, err = driver.connect(dsn); err != nil {
						t.Fatal(err)
					}

					t.Run("RecordCodec", func(t *testing.T) { RecordCodec(t, ctx, conn) })
					t.Run("RecordSearch", func(t *testing.T) { RecordSearch(t, ctx, conn) })
				})
			}

		})
	}
}
