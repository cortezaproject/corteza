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

	settingsService := settings.NewService(settings.NewRepository(repository.DB(ctx), "sys_settings"))

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

	// Assign role to user.
	var cmdUserAssignRole = &cobra.Command{
		Use:   "roleadd [userID] [roleID]",
		Short: "Assign role to user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			RoleAssignUser(args[1], args[0])
		},
	}
	cmdUsers.AddCommand(cmdUserAssignRole)

	// Role management commands.
	var cmdRole = &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}
	rootCmd.AddCommand(cmdRole)

	// Reset roles.
	var cmdRolesReset = &cobra.Command{
		Use:   "reset",
		Short: "Reset roles",
		Run: func(cmd *cobra.Command, args []string) {
			RolesReset()
		},
	}
	cmdRole.AddCommand(cmdRolesReset)

	// Add user to role.
	var cmdRoleAddUser = &cobra.Command{
		Use:   "useradd [roleID] [userID]",
		Short: "Add user to role",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			RoleAssignUser(args[0], args[1])
		},
	}
	cmdRole.AddCommand(cmdRoleAddUser)

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
