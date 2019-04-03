package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
)

func Init(ctx context.Context) {
	// Main command.
	rootCmd := &cobra.Command{Use: "system-cli"}
	db := repository.DB(ctx)

	settingsService := settings.NewService(settings.NewRepository(db, "sys_settings"))

	Settings(rootCmd, settingsService)

	ExternalAuth(ctx, rootCmd, settingsService)

	// @todo move cmd setup lines below to similar structure as Settings()

	// User management commands.
	var cmdUsers = &cobra.Command{
		Use:   "users",
		Short: "User management",
	}
	rootCmd.AddCommand(cmdUsers)

	// List users.
	var cmdUsersList = &cobra.Command{
		Use:   "list",
		Short: "List users",
		Run: func(cmd *cobra.Command, args []string) {
			UsersList()
		},
	}
	cmdUsers.AddCommand(cmdUsersList)

	roles(ctx, rootCmd, db)

	err := rootCmd.Execute()
	if err != nil {
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
