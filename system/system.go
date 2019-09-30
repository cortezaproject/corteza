package system

import (
	"context"
	"net"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/auth/external"
	"github.com/cortezaproject/corteza-server/system/commands"
	migrate "github.com/cortezaproject/corteza-server/system/db"
	"github.com/cortezaproject/corteza-server/system/grpc"
	"github.com/cortezaproject/corteza-server/system/rest"
	"github.com/cortezaproject/corteza-server/system/service"
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
			cli.HandleError(service.Init(ctx, c.Log, service.Config{
				Corredor: *c.ScriptRunner,
			}))

		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				c.InitServices(ctx, c)

				if c.ProvisionOpt.Configuration {
					cli.HandleError(provisionConfig(ctx, cmd, c))
					cli.HandleError(makeDefaultApplications(ctx, cmd, c))

					cli.HandleError(settingsAutoDiscovery(ctx, cmd, c))

					// Reload auto-configured settings
					// adding externals and oidc auto discovery depends on redirect-url setting
					cli.HandleError(service.DefaultSettings.UpdateAuthSettings(service.DefaultAuthSettings))

					cli.HandleError(authAddExternals(ctx, cmd, c))
					cli.HandleError(oidcAutoDiscovery(ctx, cmd, c))

					// Reload auto-configured settings
					cli.HandleError(service.DefaultSettings.UpdateAuthSettings(service.DefaultAuthSettings))
				}

				{
					var (
						grpcLog     = c.Log.Named("grpc-server")
						grpcLogConn = grpcLog.With(zap.String("addr", c.GRPCServerSystem.Addr))
					)

					// Temporary gRPC server initialization location
					// @todo move out of system Configure
					grpcServer := grpc.NewServer()

					ln, err := net.Listen(c.GRPCServerSystem.Network, c.GRPCServerSystem.Addr)
					if err != nil {
						grpcLogConn.Error("could not start gRPC server", zap.Error(err))
					}

					go func() {
						select {
						case <-ctx.Done():
							grpcLogConn.Debug("shutting down")
							grpcServer.GracefulStop()
							_ = ln.Close()
						}
					}()

					go func() {
						grpcLogConn.Info("Starting gRPC server")
						err := grpcServer.Serve(ln)
						grpcLogConn.Info("stopped", zap.Error(err))
					}()
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
				return commands.Importer(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Users(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Roles(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Sink(ctx, c)
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

				db = db.With(ctx).Quiet()

				return migrate.Migrate(db, c.Log)
			},
		},

		ProvisionConfig: cli.Runners{
			provisionConfig,
		},
	}
}
