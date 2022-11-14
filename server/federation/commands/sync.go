package commands

import (
	"context"

	"github.com/spf13/cobra"
)

type (
	serviceInitializer interface {
		InitServices(ctx context.Context) error
	}
)

func Sync(ctx context.Context, app serviceInitializer) *cobra.Command {
	// Sync commands.
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync commands",
	}

	// Sync structure.
	syncStructureCmd := &cobra.Command{
		Use:   "structure",
		Short: "Sync structure",

		PreRunE: commandPreRunInitService(ctx, app),
		Run:     commandSyncStructure(ctx),
	}

	syncDataCmd := &cobra.Command{
		Use:   "data",
		Short: "Sync data",

		PreRunE: commandPreRunInitService(ctx, app),
		Run:     commandSyncData(ctx),
	}

	cmd.AddCommand(
		syncStructureCmd,
		syncDataCmd,
	)

	return cmd
}

func commandPreRunInitService(ctx context.Context, app serviceInitializer) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		return app.InitServices(ctx)
	}
}
