package app

import (
	"github.com/cortezaproject/corteza-server/pkg/cli"
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

		return
	})

	app.Command.AddCommand(
		systemCommands.Users(app),
		systemCommands.Roles(app),
		systemCommands.Auth(app),
		systemCommands.RBAC(app),
		systemCommands.Sink(app),
		serveCmd,
		upgradeCmd,
		provisionCmd,
		cli.VersionCommand(),
	)

}

func (app *CortezaApp) Execute() error {
	return app.Command.Execute()
}
