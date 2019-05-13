package cli

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/types"
)

func usersCmd(ctx context.Context, db *factory.DB) *cobra.Command {
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
			userRepo := repository.User(ctx, db)
			uf := &types.UserFilter{
				OrderBy: "updated_at",
			}

			users, err := userRepo.Find(uf)
			if err != nil {
				exit(cmd, err)
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
				userRepo = repository.User(ctx, db)
				authSvc  = service.Auth(ctx)

				// @todo email validation
				user = &types.User{Email: args[0]}

				err      error
				password []byte
			)

			if user, err = userRepo.Create(user); err != nil {
				exit(cmd, err)
			}

			cmd.Printf("User created [%d].\n", user.ID)

			cmd.Print("Set password: ")
			if password, err = terminal.ReadPassword(syscall.Stdin); err != nil {
				exit(cmd, err)
			}

			if len(password) == 0 {
				// Password not set, that's ok too.
				exit(cmd, nil)
			}

			if err = authSvc.SetPassword(user.ID, string(password)); err != nil {
				exit(cmd, err)
			}
		},
	}

	pwdCmd := &cobra.Command{
		Use:   "password [email]",
		Short: "Add new user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				userRepo = repository.User(ctx, db)
				authSvc  = service.Auth(ctx)

				user     *types.User
				err      error
				password []byte
			)

			if user, err = userRepo.FindByEmail(args[0]); err != nil {
				exit(cmd, err)
			}

			cmd.Print("Set password: ")
			if password, err = terminal.ReadPassword(syscall.Stdin); err != nil {
				exit(cmd, err)
			}

			if len(password) == 0 {
				// Password not set, that's ok too.
				exit(cmd, nil)
			}

			if err = authSvc.SetPassword(user.ID, string(password)); err != nil {
				exit(cmd, err)
			}
		},
	}

	jwtCmd := &cobra.Command{
		Use:   "jwt [email-or-id]",
		Short: "Generates new JWT for a user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				userRepo = repository.User(ctx, db)
				roleRepo = repository.Role(ctx, db)
				// authSvc  = service.Auth(ctx)

				user *types.User
				err  error
				ID   uint64
				rr   types.RoleSet

				userStr = args[0]
			)

			if user, err = userRepo.FindByEmail(userStr); repository.ErrUserNotFound.Eq(err) {
				if regexp.MustCompile(`/^\d+$/`).MatchString(userStr) {
					if ID, err = strconv.ParseUint(userStr, 10, 64); err == nil {
						user, err = userRepo.FindByID(ID)
					}
				}
			}

			if err == nil {
				rr, err = roleRepo.FindByMemberID(user.ID)
			}

			if err != nil {
				exit(cmd, err)
			}

			user.SetRoles(rr.IDs())

			cmd.Println(auth.DefaultJwtHandler.Encode(user))
		},
	}

	cmd.AddCommand(
		listCmd,
		addCmd,
		pwdCmd,
		jwtCmd,
	)

	return cmd
}
