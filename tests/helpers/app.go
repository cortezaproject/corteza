package helpers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/rand"
)

type (
	TestApp struct {
		Opt *app.Options
		Log *zap.Logger
	}
)

var _ app.Runnable = &TestApp{}

func NewIntegrationTestApp(service string, parts ...app.Runnable) app.Runnable {
	opt := app.NewOptions(service)

	// When running integration tests, we want to upgrade the db. Always.
	opt.Upgrade.Always = true

	// Create a new JWT secret (to prevent any security weirdness)
	opt.Auth.Secret = string(rand.Bytes(32))
	opt.Auth.Expiry = time.Minute

	logger.DefaultLevel.SetLevel(zap.DebugLevel)
	log := logger.MakeDebugLogger()

	testApp := app.New(parts...)

	cli.HandleError(testApp.Setup(log, opt))
	cli.HandleError(testApp.Activate(cli.Context()))
	return testApp
}

func (app *TestApp) Setup(log *zap.Logger, opt *app.Options) (err error) {
	app.Log = log
	app.Opt = opt
	return
}

func (app *TestApp) Initialize(ctx context.Context) (err error) {
	return
}

func (app *TestApp) Upgrade(ctx context.Context) (err error) {
	return
}

func (app *TestApp) Activate(ctx context.Context) (err error) {
	return
}

func (app *TestApp) Provision(ctx context.Context) (err error) {
	return
}
