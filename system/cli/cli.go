package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/internal/service"
)

func Init(ctx context.Context) {
	// Main command.
	rootCmd := &cobra.Command{Use: "system-cli"}
	db := repository.DB(ctx)

	settingsService := settings.NewService(settings.NewRepository(db, "sys_settings"))

	Settings(rootCmd, settingsService)

	ExternalAuth(ctx, rootCmd, settingsService)

	users(ctx, rootCmd, service.DefaultUser)

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
