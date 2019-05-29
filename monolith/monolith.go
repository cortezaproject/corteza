package monolith

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system"
)

func Configure() *cli.Config {
	cmp := compose.Configure()
	msg := messaging.Configure()
	sys := system.Configure()

	cmp.Init()
	msg.Init()
	sys.Init()

	// Combines all three services/apps and makes them run as one monolith app
	return &cli.Config{
		ServiceName: "",

		RootCommandDBSetup: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				cli.HandleError(cmp.RootCommandDBSetup.Run(ctx, cmd, cmp))
				cli.HandleError(msg.RootCommandDBSetup.Run(ctx, cmd, msg))
				cli.HandleError(sys.RootCommandDBSetup.Run(ctx, cmd, sys))
				return
			},
		},

		RootCommandName: "corteza-server",
		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				cli.HandleError(cmp.RootCommandPreRun.Run(ctx, cmd, cmp))
				cli.HandleError(msg.RootCommandPreRun.Run(ctx, cmd, msg))
				cli.HandleError(sys.RootCommandPreRun.Run(ctx, cmd, sys))
				return
			},
		},

		ApiServerRoutes: cli.Mounters{
			func(r chi.Router) {
				r.Route("/compose", cmp.ApiServerRoutes.MountRoutes)
				r.Route("/messaging", msg.ApiServerRoutes.MountRoutes)
				r.Route("/system", sys.ApiServerRoutes.MountRoutes)
			},
		},

		AdtSubCommands: cli.CommandMakers{
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				if cc := cmp.AdtSubCommands; len(cc) > 0 {
					sub := &cobra.Command{Use: "compose", Short: "Commands from compose service"}
					sub.AddCommand(cc.Make(ctx, c)...)
					return sub
				}

				return nil
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				if cc := msg.AdtSubCommands; len(cc) > 0 {
					sub := &cobra.Command{Use: "messaging", Short: "Commands from messaging service"}
					sub.AddCommand(cc.Make(ctx, c)...)
					return sub
				}

				return nil
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				if cc := sys.AdtSubCommands; len(cc) > 0 {
					sub := &cobra.Command{Use: "system", Short: "Commands from system service"}
					sub.AddCommand(cc.Make(ctx, c)...)
					return sub
				}

				return nil
			},
		},

		ProvisionMigrateDatabase: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				cli.HandleError(cmp.ProvisionMigrateDatabase.Run(ctx, cmd, cmp))
				cli.HandleError(msg.ProvisionMigrateDatabase.Run(ctx, cmd, msg))
				cli.HandleError(sys.ProvisionMigrateDatabase.Run(ctx, cmd, sys))
				return
			},
		},

		ProvisionAccessControl: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				cli.HandleError(cmp.ProvisionAccessControl.Run(ctx, cmd, cmp))
				cli.HandleError(msg.ProvisionAccessControl.Run(ctx, cmd, msg))
				cli.HandleError(sys.ProvisionAccessControl.Run(ctx, cmd, sys))
				return
			},
		},
	}
}
