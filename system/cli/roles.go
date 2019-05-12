package cli

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

func rolesCmd(ctx context.Context, db *factory.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}

	addUserCmd := &cobra.Command{
		Use:   "useradd [role-ID-or-name-or-handle] [user-ID-or-email]",
		Short: "Add user to role",
		Args:  cobra.ExactArgs(2),
		Run:   rolesUserAddCmd(ctx, db),
	}

	cmd.AddCommand(addUserCmd)

	return cmd
}

func rolesUserAddCmd(ctx context.Context, db *factory.DB) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Create role and user repository.
		var (
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
			exit(cmd, err)
		} else if len(rr) == 1 {
			role = rr[0]
		} else if len(rr) > 1 {
			exit(cmd, errors.Errorf("too many roles found with name %q", roleStr))
		} else if role == nil {
			if ID, err = strconv.ParseUint(roleStr, 10, 64); err != nil {
				// Could not parse ID out of role string
				return
			} else if role, err = roleRepo.FindByID(ID); err != nil {
				return
			}
		}

		if user, err = userRepo.FindByEmail(userStr); repository.ErrUserNotFound.Eq(err) {
			exit(cmd, err)
		} else if user == nil || user.ID == 0 {
			if ID, err = strconv.ParseUint(userStr, 10, 64); err != nil {
				exit(cmd, err)
			} else if user, err = userRepo.FindByID(ID); err != nil {
				exit(cmd, err)
			}
		}

		// Add user to role.
		if err = roleRepo.MemberAddByID(role.ID, user.ID); err != nil {
			exit(cmd, err)
		}

		cmd.Printf("Added user [%d] %q to [%d] %q role\n", user.ID, user.Email, role.ID, role.Name)
	}
}
