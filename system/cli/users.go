package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

func usersCmd(ctx context.Context, db *factory.DB) *cobra.Command {
	// User management commands.
	cmd := &cobra.Command{
		Use:   "users",
		Short: "User management",
	}

	// List users.
	cmdUsersList := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Run: func(cmd *cobra.Command, args []string) {
			userRepo := repository.User(ctx, db)
			uf := &types.UserFilter{
				OrderBy: "updated_at",
			}

			users, err := userRepo.Find(uf)
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

	cmd.AddCommand(cmdUsersList)

	return cmd
}
