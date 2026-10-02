package helpers

import (
	"context"
	"fmt"
	"os"

	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/system/types"

	// Explicitly register SQLite (not done in the app as for testing only)
	_ "github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
)

// The env var that asks for a real database, and the one that must not be
// mistaken for it.
const (
	integrationDSNEnv = "INTEGRATION_DSN"
	appDSNEnv         = "DB_DSN"
)

// What an integration suite gets when it does not ask for anything: the
// in-memory SQLite that is also the app's own default DSN. cache=shared is
// what makes the pool's connections share one database rather than each
// opening an empty one.
const inMemoryDSN = "sqlite3://file::memory:?cache=shared&mode=memory"

// The database the suites run against, which is never the one the app is
// configured with.
//
// Every suite here truncates whole tables between tests, so a suite that
// inherits DB_DSN (for example from a developer's .env) destroys the instance
// that variable points at, completely and without a word. The DSN therefore
// does not come from the app's own environment: it is in-memory unless
// INTEGRATION_DSN names something else, and naming the app's database there
// is refused rather than obeyed.
func integrationDSN() string {
	dsn := os.Getenv(integrationDSNEnv)
	if dsn == "" {
		return inMemoryDSN
	}

	if dsn == os.Getenv(appDSNEnv) {
		fmt.Fprintf(os.Stderr,
			"%s names the same database as %s.\n"+
				"These tests truncate whole tables; point %s at a database of its own, or unset it to run in memory.\n",
			integrationDSNEnv, appDSNEnv, integrationDSNEnv)
		os.Exit(1)
	}

	return dsn
}

func NewIntegrationTestApp(ctx context.Context, initTestServices func(*app.CortezaApp) error) *app.CortezaApp {
	// Enforce debug logger for tests
	logger.SetDefault(logger.MakeDebugLogger())

	var (
		a = app.New()
	)

	a.Opt = options.Init()
	a.Opt.DB.DSN = integrationDSN()

	// When running integration tests, we want to upgrade the db. Always.
	a.Opt.Upgrade.Always = true

	// Create a new JWT secret (to prevent any security weirdness)
	a.Opt.Auth.Secret = string(rand.Bytes(32))
	a.Opt.Auth.DefaultClient = ""

	a.Log = logger.Default()

	a.DefaultAuthClient = &types.AuthClient{ID: 1, Handle: "test-auth-client", Secret: "integration-tests"}

	cli.HandleError(a.InitStore(ctx))
	cli.HandleError(initTestServices(a))
	cli.HandleError(a.InitServices(ctx))
	return a
}
