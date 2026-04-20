package helpers

import (
	"context"

	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/system/types"

	// Explicitly register SQLite (not done in the app as for testing only)
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
)

func NewIntegrationTestApp(ctx context.Context, initTestServices func(*app.HumanApp) error) *app.HumanApp {
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
	a.Opt.Auth.DefaultClient = ""

	a.Log = logger.Default()

	a.DefaultAuthClient = &types.AuthClient{ID: 1, Handle: "test-auth-client", Secret: "integration-tests"}

	cli.HandleError(a.InitStore(ctx))
	cli.HandleError(initTestServices(a))
	cli.HandleError(a.InitServices(ctx))
	return a
}
