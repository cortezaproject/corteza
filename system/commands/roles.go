package commands

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cobra"
)

func Roles(ctx context.Context, app serviceInitializer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}

	addUserCmd := &cobra.Command{
		Use:     "useradd [role-ID-or-name-or-handle] [user-ID-or-email]",
		Short:   "Add user to role",
		Args:    cobra.ExactArgs(2),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				roleStr, userStr = args[0], args[1]

				role *types.Role
				user *types.User

				err error
			)

			role, err = service.DefaultRole.FindByAny(ctx, roleStr)
			cli.HandleError(err)

			user, err = service.DefaultUser.FindByAny(ctx, userStr)
			cli.HandleError(err)

			err = service.DefaultRole.MemberAdd(ctx, role.ID, user.ID)
			cli.HandleError(err)

			cmd.Printf("Added user [%d] %q to [%d] %q role\n", user.ID, user.Email, role.ID, role.Name)
		},
	}

	cmd.AddCommand(addUserCmd)

	return cmd
}
