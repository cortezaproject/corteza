package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
)

func StartCLI(ctx context.Context) {
	var (
		db              = repository.DB(ctx)
		settingsService = settings.NewService(settings.NewRepository(db, "sys_settings"))

		cmd = &cobra.Command{Use: "system-cli"}
	)

	cmd.AddCommand(
		settingsCmd(ctx, settingsService),
		authCmd(ctx, db, settingsService),
		usersCmd(ctx, db),
		rolesCmd(ctx, db),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func exit(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
