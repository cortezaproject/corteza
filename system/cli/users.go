package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/types"
)

func users(ctx context.Context, rootCmd *cobra.Command, userService service.UserService) {
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
			uf := &types.UserFilter{
				OrderBy: "updated_at",
			}

			users, err := userService.With(ctx).Find(uf)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("                     Created    Updated    Email")
			for _, u := range users {
				upd := "---- -- --"

				if u.UpdatedAt != nil {
					upd = u.UpdatedAt.Format("2006-01-02")
				}

				fmt.Printf(
					"%20d %s %s %-100s %s\n",
					u.ID,
					u.CreatedAt.Format("2006-01-02"),
					upd,
					u.Email,
					u.Name)
			}
		},
	}
	cmdUsers.AddCommand(cmdUsersList)

}
