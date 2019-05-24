package commands

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

func Roles(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}

	addUserCmd := &cobra.Command{
		Use:   "useradd [role-ID-or-name-or-handle] [user-ID-or-email]",
		Short: "Add user to role",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Create role and user repository.
			var (
				db = factory.Database.MustGet("system")

				roleStr, userStr = args[0], args[1]

				roleRepo = repository.Role(ctx, db)
				userRepo = repository.User(ctx, db)

				rr   []*types.Role
				role *types.Role
				user *types.User
				ID   uint64

				err error
			)

			// Try to find role by name and by ID
			if rr, err = roleRepo.Find(&types.RoleFilter{Query: roleStr}); err != nil {
				exit(err)
			} else if len(rr) == 1 {
				role = rr[0]
			} else if len(rr) > 1 {
				exit(errors.Errorf("too many roles found with name %q", roleStr))
			} else if role == nil {
				if ID, err = strconv.ParseUint(roleStr, 10, 64); err != nil {
					// Could not parse ID out of role string
					return
				} else if role, err = roleRepo.FindByID(ID); err != nil {
					return
				}
			}

			if user, err = userRepo.FindByEmail(userStr); repository.ErrUserNotFound.Eq(err) {
				exit(err)
			} else if user == nil || user.ID == 0 {
				if ID, err = strconv.ParseUint(userStr, 10, 64); err != nil {
					exit(err)
				} else if user, err = userRepo.FindByID(ID); err != nil {
					exit(err)
				}
			}

			// Add user to role.
			if err = roleRepo.MemberAddByID(role.ID, user.ID); err != nil {
				exit(err)
			}

			cmd.Printf("Added user [%d] %q to [%d] %q role\n", user.ID, user.Email, role.ID, role.Name)
		},
	}

	cmd.AddCommand(addUserCmd)

	return cmd
}
