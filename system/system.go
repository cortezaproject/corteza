package system

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/commands"
	migrate "github.com/cortezaproject/corteza-server/system/db"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest"
)

const (
	system = "system"
)

func Configure() *cli.Config {
	var (
		accessControlSetup = func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
			// Calling grant directly on internal permissions service to avoid AC check for "grant"
			var p = service.DefaultPermissions
			var ac = service.DefaultAccessControl
			return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
		}
	)

	return &cli.Config{
		ServiceName: system,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				cli.HandleError(service.Init(ctx, c.Log))

				if c.ProvisionOpt.AutoSetup {
					cli.HandleError(accessControlSetup(ctx, cmd, c))

					// Run auto configuration
					commands.SettingsAutoConfigure(cmd, "", "", "", "")

					// Reload auto-configured settings
					service.DefaultAuthSettings, _ = service.DefaultSettings.LoadAuthSettings()
				}

				return
			},
		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				go service.Watchers(ctx)
				return nil
			},
		},

		ApiServerRoutes: cli.Mounters{
			rest.MountRoutes,
		},

		AdtSubCommands: cli.CommandMakers{
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Settings(ctx)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Auth(ctx)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Users(ctx)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Roles(ctx)
			},
		},

		ProvisionMigrateDatabase: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				if !c.ProvisionOpt.MigrateDatabase {
					return nil
				}

				var db, err = factory.Database.Get(system)
				if err != nil {
					return err
				}

				db = db.With(ctx)
				// Disable profiler for migrations
				db.Profiler = nil

				return migrate.Migrate(db, c.Log)
			},
		},

		ProvisionAccessControl: cli.Runners{
			accessControlSetup,
		},
	}
}
