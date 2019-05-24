package commands

import (
	"context"
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func Users(ctx context.Context) *cobra.Command {
	// User management commands.
	cmd := &cobra.Command{
		Use:   "users",
		Short: "User management",
	}

	// List users.
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				db = factory.Database.MustGet("system")
			)

			userRepo := repository.User(ctx, db)
			uf := &types.UserFilter{
				OrderBy: "updated_at",
			}

			users, err := userRepo.Find(uf)
			if err != nil {
				exit(err)
			}

			fmt.Println("                     Created    Updated    EmailAddress")
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

	addCmd := &cobra.Command{
		Use:   "add [email]",
		Short: "Add new user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				db = factory.Database.MustGet("system")

				userRepo = repository.User(ctx, db)
				authSvc  = service.Auth(ctx)

				// @todo email validation
				user = &types.User{Email: args[0]}

				err      error
				password []byte
			)

			if user, err = userRepo.Create(user); err != nil {
				exit(err)
			}

			cmd.Printf("User created [%d].\n", user.ID)

			cmd.Print("Set password: ")
			if password, err = terminal.ReadPassword(syscall.Stdin); err != nil {
				exit(err)
			}

			if len(password) == 0 {
				// Password not set, that's ok too.
				exit(nil)
			}

			if err = authSvc.SetPassword(user.ID, string(password)); err != nil {
				exit(err)
			}
		},
	}

	pwdCmd := &cobra.Command{
		Use:   "password [email]",
		Short: "Add new user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				db = factory.Database.MustGet("system")

				userRepo = repository.User(ctx, db)
				authSvc  = service.Auth(ctx)

				user     *types.User
				err      error
				password []byte
			)

			if user, err = userRepo.FindByEmail(args[0]); err != nil {
				exit(err)
			}

			cmd.Print("Set password: ")
			if password, err = terminal.ReadPassword(syscall.Stdin); err != nil {
				exit(err)
			}

			if len(password) == 0 {
				// Password not set, that's ok too.
				exit(nil)
			}

			if err = authSvc.SetPassword(user.ID, string(password)); err != nil {
				exit(err)
			}
		},
	}

	cmd.AddCommand(
		listCmd,
		addCmd,
		pwdCmd,
	)

	return cmd
}
