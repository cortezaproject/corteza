package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/store/mysql"
	"github.com/cortezaproject/corteza-server/store/pgsql"
	"github.com/cortezaproject/corteza-server/store/sqlite"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"os"
	"testing"
)

type (
	storeInterface interface {
		storeGeneratedInterfaces
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
			init      func(ctx context.Context, dsn string) (storeInterface, error)
		}

		upgrader interface {
			Upgrade(context.Context, *zap.Logger) error
		}
	)

	var (
		ctx = context.Background()

		ss = []suite{
			{
				name:      "MySQL",
				dsnEnvKey: "RDBMS_MYSQL_DSN",
				init:      func(ctx context.Context, dsn string) (storeInterface, error) { return mysql.New(ctx, dsn) },
			},
			{
				name:      "PostgreSQL",
				dsnEnvKey: "RDBMS_PGSQL_DSN",
				init:      func(ctx context.Context, dsn string) (storeInterface, error) { return pgsql.New(ctx, dsn) },
			},
			{
				name:      "CockroachDB",
				dsnEnvKey: "RDBMS_COCKROACHDB_DSN",
				init:      nil,
			},
			{
				name:      "SQLite",
				dsnEnvKey: "RDBMS_SQLITE_DSN",
				init:      func(ctx context.Context, dsn string) (storeInterface, error) { return sqlite.New(ctx, dsn) },
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

			t.Logf("connecting to %s with %s", s.name, dsn)

			genericStore, err := s.init(ctx, dsn)
			if err != nil {
				t.Errorf("failed to initialize %s store: %s", s.name, err.Error())
				return
			}

			if up, ok := genericStore.(upgrader); ok {
				err = up.Upgrade(ctx, zap.NewNop())
				if err != nil {
					t.Errorf("failed to upgrade %s store: %s", s.name, err.Error())
					return
				}
			}

			testAllGenerated(t, genericStore)

			t.Run("ComposeRecords", func(t *testing.T) {
				var s = genericStore.(composeRecordsStore)
				testComposeRecords(t, s)
			})
		})
	}

}
