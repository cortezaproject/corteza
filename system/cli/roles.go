package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

func rolesCmd(ctx context.Context, db *factory.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Role management",
	}

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset roles",
		Run:   rolesResetCmd(ctx, db),
	}

	addUserCmd := &cobra.Command{
		Use:   "useradd [role-ID-or-name-or-handle] [user-ID-or-email]",
		Short: "Add user to role",
		Args:  cobra.ExactArgs(2),
		Run:   rolesUserAddCmd(ctx, db),
	}

	cmd.AddCommand(resetCmd, addUserCmd)

	return cmd
}

func rolesResetCmd(ctx context.Context, db *factory.DB) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var (
			err      error
			roleRepo = repository.Role(ctx, db)
			ruleRepo = rules.NewResources(ctx, db)
		)

		// Recreates core roles (will not reset memberships!)
		err = types.RoleSet{
			&types.Role{ID: 1, Handle: "everyone", Name: "Everyone"},
			&types.Role{ID: 2, Handle: "admins", Name: "Administrators"},
		}.Walk(func(r *types.Role) error {
			// Update will exec REPLACE and take care of creation as well
			r.CreatedAt = time.Now()
			_, err := roleRepo.Update(r)
			return err
		})

		if err != nil {
			fmt.Printf("could not reset roles: %v", err)
			return
		}

		rules := []rules.Rule{
			{1, "system", "user.create", 2},
			{1, "compose", "access", 2},
			{1, "messaging", "access", 2},
			{2, "compose", "namespace.create", 2},
			{2, "compose", "access", 2},
			{2, "compose", "grant", 2},
			{2, "compose:namespace:*", "page.create", 2},
			{2, "compose:namespace:*", "read", 2},
			{2, "compose:namespace:*", "update", 2},
			{2, "compose:namespace:*", "delete", 2},
			{2, "compose:namespace:*", "module.create", 2},
			{2, "compose:namespace:*", "chart.create", 2},
			{2, "compose:namespace:*", "trigger.create", 2},
			{2, "compose:chart:*", "read", 2},
			{2, "compose:chart:*", "update", 2},
			{2, "compose:chart:*", "delete", 2},
			{2, "compose:trigger:*", "read", 2},
			{2, "compose:trigger:*", "update", 2},
			{2, "compose:trigger:*", "delete", 2},
			{2, "compose:page:*", "read", 2},
			{2, "compose:page:*", "update", 2},
			{2, "compose:page:*", "delete", 2},
			{2, "system", "access", 2},
			{2, "system", "grant", 2},
			{2, "system", "settings.read", 2},
			{2, "system", "settings.manage", 2},
			{2, "system", "organisation.create", 2},
			{2, "system", "user.create", 2},
			{2, "system", "role.create", 2},
			{2, "system:organisation:*", "access", 2},
			{2, "system:user:*", "read", 2},
			{2, "system:user:*", "update", 2},
			{2, "system:user:*", "suspend", 2},
			{2, "system:user:*", "unsuspend", 2},
			{2, "system:user:*", "delete", 2},
			{2, "system:role:*", "read", 2},
			{2, "system:role:*", "update", 2},
			{2, "system:role:*", "delete", 2},
			{2, "system:role:*", "members.manage", 2},
			{2, "messaging", "access", 2},
			{2, "messaging", "grant", 2},
			{2, "messaging", "channel.public.create", 2},
			{2, "messaging", "channel.private.create", 2},
			{2, "messaging", "channel.group.create", 2},
			{2, "messaging:channel:*", "update", 2},
			{2, "messaging:channel:*", "leave", 2},
			{2, "messaging:channel:*", "read", 2},
			{2, "messaging:channel:*", "join", 2},
			{2, "messaging:channel:*", "delete", 2},
			{2, "messaging:channel:*", "undelete", 2},
			{2, "messaging:channel:*", "archive", 2},
			{2, "messaging:channel:*", "unarchive", 2},
			{2, "messaging:channel:*", "members.manage", 2},
			{2, "messaging:channel:*", "webhooks.manage", 2},
			{2, "messaging:channel:*", "attachments.manage", 2},
			{2, "messaging:channel:*", "message.attach", 2},
			{2, "messaging:channel:*", "message.update.all", 2},
			{2, "messaging:channel:*", "message.update.own", 2},
			{2, "messaging:channel:*", "message.delete.all", 2},
			{2, "messaging:channel:*", "message.delete.own", 2},
			{2, "messaging:channel:*", "message.embed", 2},
			{2, "messaging:channel:*", "message.send", 2},
			{2, "messaging:channel:*", "message.reply", 2},
			{2, "messaging:channel:*", "message.react", 2},
			{2, "compose:module:*", "read", 2},
			{2, "compose:module:*", "update", 2},
			{2, "compose:module:*", "delete", 2},
			{2, "compose:module:*", "record.create", 2},
			{2, "compose:module:*", "record.read", 2},
			{2, "compose:module:*", "record.update", 2},
			{2, "compose:module:*", "record.delete", 2},
		}

		var purgedRoles = map[uint64]bool{}
		for _, r := range rules {
			if !purgedRoles[r.RoleID] {
				if err = ruleRepo.Delete(r.RoleID); err != nil {
					fmt.Printf("could not reset rules: %v", err)
					return
				}

				purgedRoles[r.RoleID] = true
			}
		}

		// Recreates rules for core roles
		err = ruleRepo.Reset(rules)

		if err != nil {
			fmt.Printf("could not reset rules: %v", err)
		}
	}
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

		if user, err = userRepo.FindByEmail(userStr); err != nil && err != repository.ErrUserNotFound {
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
