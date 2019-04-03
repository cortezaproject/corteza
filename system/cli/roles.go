package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

func RolesReset() {
	var (
		err      error
		ctx      = context.Background()
		db       = repository.DB(ctx)
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

func RoleAssignUser(roleStr string, userStr string) {
	ctx := context.Background()
	db := repository.DB(ctx)

	// Create role and user repository.
	roleRepo := repository.Role(ctx, db)
	userRepo := repository.User(ctx, db)

	var err error

	// Try to parse roleID.
	roleID, err := strconv.ParseUint(roleStr, 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if role ID exists.
	role, err := roleRepo.FindByID(roleID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var userID uint64
	var user *types.User

	// Try to parse userID.
	userID, err = strconv.ParseUint(userStr, 10, 64)
	if err != nil {
		user, err = userRepo.FindByEmail(userStr)
	} else {
		user, err = userRepo.FindByID(userID)
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Add user to role.
	err = roleRepo.MemberAddByID(role.ID, user.ID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added user: %d %s to role: %s\n", user.ID, user.Email, role.Name)
}
