package app

import (
	"context"
	"sync"

	authCommands "github.com/cortezaproject/corteza-server/auth/commands"
	federationCommands "github.com/cortezaproject/corteza-server/federation/commands"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/plugin"
	fakerCommands "github.com/cortezaproject/corteza-server/pkg/seeder/commands"
	"github.com/cortezaproject/corteza-server/store"
	systemCommands "github.com/cortezaproject/corteza-server/system/commands"
	"go.uber.org/zap"
)

// InitCLI function initializes basic Corteza subsystems
// and sets-up the command line interface
func (app *CortezaApp) InitCLI() {
	var (
		ctx = cli.Context()

		// path to all environmental files (or locations with .env file)
		// filled from flag values
		envs []string
	)

	app.Command = cli.RootCommand(func() (err error) {
		log := app.Log.Named("plugins")
		if app.Opt.Plugins.Enabled && len(app.Opt.Plugins.Paths) > 0 {
			log.Warn("loading", zap.String("paths", app.Opt.Plugins.Paths))

			var paths []string
			paths, err = plugin.Resolve(app.Opt.Plugins.Paths)
			log.Warn("loading", zap.Strings("resolved-paths", paths))

			app.plugins, err = plugin.Load(paths...)
			if err != nil {
				return err
			}
		} else {
			// Empty set of plugins
			app.plugins = plugin.Set{}
		}

		return err
	})

	// Environmental variables (from the env, files, see cli.LoadEnv) MUST be
	// loaded at this point!
	if len(envs) == 0 {
		envs = []string{"."}
	}

	// Loading env after the rootcommand is added, so as not to break
	// the Execute() if there is an error
	if err := cli.LoadEnv(envs...); err != nil {
		app.Log.Error("failed to load environmental variables", zap.Error(err))
		return
	}

	app.Opt = options.Init()

	app.Command.Flags().StringSliceVar(&envs, "env-file", nil,
		"Load environmental variables from files and directories containing .env file.\n"+
			"Values from loaded files DO NOT override existing variables from the environment.\n"+
			"This flag can be used multiple times, values are loaded from all provided locations.\n"+
			"If no paths are provided, corteza loads .env file from the current directory (equivalent to --env-file .)")

	serveCmd := cli.ServeCommand(func() (err error) {
		wg := &sync.WaitGroup{}

		{ // @todo refactor wait-for out of HTTP API server.
			app.HttpServer = server.New(app.Log, app.Opt)

			wg.Add(1)
			go func() {
				app.HttpServer.Serve(actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_API_REST))
				wg.Done()
			}()
		}

		{
			// @todo add other server/listeners here...
			//wg.Add(1)
			//go func(ctx context.Context) {
			//	grpcApi.Serve(actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_API_GRPC))
			//	wg.Done()
			//}(ctx)
		}

		if err = app.Activate(ctx); err != nil {
			return
		}

		app.HttpServer.Activate(app.mountHttpRoutes)

		// Wait for all servers to be done
		wg.Wait()

		app.HttpServer.Shutdown()

		return nil
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

	storeInit := func(ctx context.Context) (store.Storer, error) {
		err := app.InitStore(ctx)
		return app.Store, err
	}

	app.Command.AddCommand(
		systemCommands.Users(ctx, app),
		systemCommands.Roles(ctx, app),
		systemCommands.RBAC(ctx, storeInit),
		systemCommands.Sink(ctx, app),
		systemCommands.Settings(ctx, app),
		systemCommands.Import(ctx, storeInit),
		systemCommands.Export(ctx, storeInit),
		serveCmd,
		upgradeCmd,
		provisionCmd,
		authCommands.Command(ctx, app, storeInit),
		federationCommands.Sync(ctx, app),
		cli.EnvCommand(),
		cli.VersionCommand(),
		fakerCommands.Seeder(ctx, app),
	)

}

func (app *CortezaApp) Execute() error {
	return app.Command.Execute()
}
