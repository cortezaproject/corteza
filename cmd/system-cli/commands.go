package main

import (
	"fmt"

	"github.com/spf13/cobra"

	systemCli "github.com/crusttech/crust/system/cli"
)

func setupCobra() {
	// Main command.
	rootCmd := &cobra.Command{Use: "system-cli"}

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
			systemCli.UsersList()
		},
	}
	cmdUsers.AddCommand(cmdUsersList)

	// Role management commands.
	var cmdRole = &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}
	rootCmd.AddCommand(cmdRole)

	// Reset permissions.
	var cmdRolesReset = &cobra.Command{
		Use:   "reset",
		Short: "Reset roles",
		Run: func(cmd *cobra.Command, args []string) {
			systemCli.RolesReset()
		},
	}
	cmdRole.AddCommand(cmdRolesReset)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
