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

	// Assign role to user.
	var cmdUserAssignRole = &cobra.Command{
		Use:   "roleadd [userID] [roleID]",
		Short: "Assign role to user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			systemCli.RoleAssignUser(args[1], args[0])
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
			systemCli.RolesReset()
		},
	}
	cmdRole.AddCommand(cmdRolesReset)

	// Add user to role.
	var cmdRoleAddUser = &cobra.Command{
		Use:   "useradd [roleID] [userID]",
		Short: "Add user to role",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			systemCli.RoleAssignUser(args[0], args[1])
		},
	}
	cmdRole.AddCommand(cmdRoleAddUser)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
