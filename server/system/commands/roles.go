package commands

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cobra"
)

func Roles(ctx context.Context, app serviceInitializer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "roles",
		Aliases: []string{"role"},
		Short:   "Role management",
	}

	cmd.AddCommand(rolesAddUser(ctx, app))
	cmd.AddCommand(rolesList(ctx, app))

	return cmd
}

func rolesAddUser(ctx context.Context, app serviceInitializer) *cobra.Command {
	return &cobra.Command{
		Use:     "useradd [role-ID-or-name-or-handle] [user-ID-or-email]",
		Aliases: []string{"assign", "add-user", "adduser", "user-add", "user"},
		Short:   "Add user to role",
		Args:    cobra.ExactArgs(2),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

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

			cli.HandleError(store.CreateRoleMember(ctx, service.DefaultStore, &types.RoleMember{
				RoleID: role.ID,
				UserID: user.ID,
			}))

			cmd.Printf("Added user [%d] %q to [%d] %q role\n", user.ID, user.Email, role.ID, role.Name)
		},
	}
}

func rolesList(ctx context.Context, app serviceInitializer) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all roles",
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			f := types.RoleFilter{Query: ""}
			if len(args) > 0 {
				f.Query = args[0]
			}

			rr, _, err := store.SearchRoles(ctx, service.DefaultStore, f)
			cli.HandleError(err)

			for _, r := range rr {
				cmd.Printf("%-20d %-30s %s\n", r.ID, r.Handle, r.Name)
			}
		},
	}
}
