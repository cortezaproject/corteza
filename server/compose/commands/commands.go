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

func Base(ctx context.Context, app serviceInitializer) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "compose",
		Aliases: []string{"cmp"},
	}

	cmd.AddCommand(
		Records(ctx, app),
	)

	return
}
