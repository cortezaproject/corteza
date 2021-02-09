package app

import (
	"context"

	federationCommands "github.com/cortezaproject/corteza-server/federation/commands"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemCommands "github.com/cortezaproject/corteza-server/system/commands"
)

// CLI function initializes basic Corteza subsystems
// and sets-up the command line interface
func (app *CortezaApp) InitCLI() {
	ctx := cli.Context()

	app.Command = cli.RootCommand(nil)

	serveCmd := cli.ServeCommand(func() (err error) {
		if err = app.Activate(ctx); err != nil {
			return
		}

		return app.Serve(ctx)
	})

	upgradeCmd := cli.UpgradeCommand(func() (err error) {
		if err = app.InitStore(ctx); err != nil {
			return
		}

		return
	})

	provisionCmd := cli.ProvisionCommand(func() (err error) {
		if err = app.Provision(ctx); err != nil {
			return
		}

		// Provisioning doesn't automatically reload rbac rules, so this is required
		rbac.Global().Reload(ctx)
		return
	})

	storeInit := func(ctx context.Context) (store.Storer, error) {
		err := app.InitStore(ctx)
		return app.Store, err
	}

	app.Command.AddCommand(
		systemCommands.Users(app),
		systemCommands.Roles(app),
		systemCommands.Auth(app, app.Opt.Auth),
		systemCommands.RBAC(app),
		systemCommands.Sink(app),
		systemCommands.Settings(),
		systemCommands.Import(storeInit),
		serveCmd,
		upgradeCmd,
		provisionCmd,
		federationCommands.Sync(app),
		cli.VersionCommand(),
	)

}

func (app *CortezaApp) Execute() error {
	return app.Command.Execute()
}
