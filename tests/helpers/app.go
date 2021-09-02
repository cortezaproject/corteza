package helpers

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	// Explicitly register SQLite (not done in the app as for testing only)
	_ "github.com/cortezaproject/corteza-server/store/sqlite3"
)

func NewIntegrationTestApp(ctx context.Context, initTestServices func(*app.CortezaApp) error) *app.CortezaApp {
	// Enforce debug logger for tests
	logger.SetDefault(logger.MakeDebugLogger())

	var (
		a = app.New()
	)

	a.Opt = options.Init()

	// When running integration tests, we want to upgrade the db. Always.
	a.Opt.Upgrade.Always = true

	// Create a new JWT secret (to prevent any security weirdness)
	a.Opt.Auth.Secret = string(rand.Bytes(32))
	a.Opt.Auth.Expiry = time.Minute
	a.Opt.Auth.DefaultClient = ""

	a.Log = logger.Default()

	cli.HandleError(a.InitStore(ctx))
	cli.HandleError(initTestServices(a))
	cli.HandleError(a.InitServices(ctx))
	return a
}
