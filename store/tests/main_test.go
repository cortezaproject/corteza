package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/mysql"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/postgres"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
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

func Test_RDBMS_SQLITE(t *testing.T) {
	testAllGenerated(t, setup(t, sqlite.Connect))
}

func Test_RDBMS_MYSQL(t *testing.T) {
	testAllGenerated(t, setup(t, mysql.Connect))
}

func Test_RDBMS_PGSQL(t *testing.T) {
	testAllGenerated(t, setup(t, postgres.Connect))
}

func Test_RDBMS_COCKROACHDB(t *testing.T) {
	testAllGenerated(t, setup(t, nil))
}

func Test_MEMORY(t *testing.T) {
	testAllGenerated(t, setup(t, nil))
}

func Test_MONGODB(t *testing.T) {
	testAllGenerated(t, setup(t, nil))
}

func Test_ELASTICSEARCH(t *testing.T) {
	testAllGenerated(t, setup(t, nil))
}

func setup(t *testing.T, connect store.ConnectorFn) (s store.Storer) {
	t.Parallel()

	if connect == nil {
		t.Skipf("connection function not set, skipping tests")
		return
	}

	var (
		err      error
		env      = t.Name()[5:] + "_DSN"
		dsn, has = os.LookupEnv(env)
		ctx      = context.Background()
	)

	if !has || len(dsn) == 0 {
		t.Skipf("no %s found, skipping tests", env)
		return
	}

	s, err = connect(ctx, dsn)
	if err != nil {
		t.Errorf("failed to initialize: %v", err)
		return
	}

	err = store.Upgrade(ctx, zap.NewNop(), s)
	if err != nil {
		t.Fatalf("failed to upgrade : %v", err)
		return
	}

	return s
}
