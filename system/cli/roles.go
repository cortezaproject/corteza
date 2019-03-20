package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

func RolesReset() {
	ctx := context.Background()
	db := repository.DB(ctx)

	err := repository.Role(ctx, db).Reset()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Everyone and Administrators role were reset.")
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
