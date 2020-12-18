package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/mysql"
	"github.com/cortezaproject/corteza-server/store/postgres"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

var (
	now = func() *time.Time {
		n := time.Now().Round(time.Second)
		return &n
	}
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func Test_Store(t *testing.T) {
	type (
		suite struct {
			name      string
			dsnEnvKey string
			init      store.ConnectorFn
		}
	)

	logger.SetDefault(logger.MakeDebugLogger())

	var (
		ctx = context.Background()

		ss = []suite{
			{
				name:      "MySQL",
				dsnEnvKey: "RDBMS_MYSQL_DSN",
				init:      mysql.Connect,
			},
			{
				name:      "PostgreSQL",
				dsnEnvKey: "RDBMS_PGSQL_DSN",
				init:      postgres.Connect,
			},
			{
				name:      "CockroachDB",
				dsnEnvKey: "RDBMS_COCKROACHDB_DSN",
				init:      nil,
			},
			{
				name:      "SQLite",
				dsnEnvKey: "RDBMS_SQLITE_DSN",
				init:      sqlite3.Connect,
			},
			{
				name:      "InMemory",
				dsnEnvKey: "MEMORY_DSN",
				init:      nil,
			},
			{
				name:      "MongoDB",
				dsnEnvKey: "MONGODB_DSN",
				init:      nil,
			},
			{
				name:      "ElasticSearch",
				dsnEnvKey: "ELASTICSEARCH_DSN",
				init:      nil,
			},
		}
	)

	for _, s := range ss {
		t.Run(s.name, func(t *testing.T) {
			dsn, has := os.LookupEnv(s.dsnEnvKey)
			if !has {
				t.Skipf("no %s found, skipping %s store tests", s.dsnEnvKey, s.name)
				return
			}

			genericStore, err := s.init(ctx, dsn)
			if err != nil {
				t.Errorf("failed to initialize %s store: %s", s.name, err.Error())
				return
			}

			err = store.Upgrade(ctx, zap.NewNop(), genericStore)
			if err != nil {
				t.Errorf("failed to upgrade %s store: %s", s.name, err.Error())
				return
			}

			testAllGenerated(t, genericStore)
		})
	}

}
