package system

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/commands"
	migrate "github.com/cortezaproject/corteza-server/system/db"
	"github.com/cortezaproject/corteza-server/system/internal/auth/external"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest"
)

const (
	system = "system"
)

func Configure() *cli.Config {
	var (
		servicesInitialized bool
	)

	return &cli.Config{
		ServiceName: system,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				return
			},
		},

		InitServices: func(ctx context.Context, c *cli.Config) {
			if servicesInitialized {
				return
			}
			servicesInitialized = true

			// storagePath := options.EnvString("", "SYSTEM_STORAGE_PATH", "var/store")
			cli.HandleError(service.Init(ctx, c.Log))

		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				c.InitServices(ctx, c)

				if c.ProvisionOpt.AutoSetup {
					cli.HandleError(accessControlSetup(ctx, cmd, c))
					cli.HandleError(makeDefaultApplications(ctx, cmd, c))
					cli.HandleError(discoverSettings(ctx, cmd, c))

					// Run auto configuration
					commands.SettingsAutoConfigure(cmd)

					// Reload auto-configured settings
					service.DefaultAuthSettings, _ = service.DefaultSettings.LoadAuthSettings()
				}

				// Initialize external authentication (from default settings)
				external.Init()
				go service.Watchers(ctx)
				return nil
			},
		},

		ApiServerRoutes: cli.Mounters{
			rest.MountRoutes,
		},

		AdtSubCommands: cli.CommandMakers{
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Settings(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Auth(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Users(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Roles(ctx, c)
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
