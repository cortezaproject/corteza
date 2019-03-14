package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
)

func TestUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Run tests in transaction to maintain DB state.
	Error(t, db.Transaction(func() error {
		userRepo := User(context.Background(), db)
		user := &types.User{
			Name:     "John User Doe",
			Username: "johndoe",
			SatosaID: "1234",
		}
		user.GeneratePassword("johndoe")

		{
			uu, err := userRepo.Create(user)
			assert(t, err == nil, "Owner.Create error: %+v", err)
			assert(t, user.ID == uu.ID, "Changes were not stored")
		}

		roleRepo := Role(context.Background(), db)
		role := &types.Role{
			Name: "Test role v1",
		}

		{
			t1, err := roleRepo.Create(role)
			assert(t, err == nil, "Role.Create error: %+v", err)
			assert(t, role.Name == t1.Name, "Changes were not stored")

			err = roleRepo.MemberAddByID(t1.ID, user.ID)
			assert(t, err == nil, "Role.MemberAddByID error: %+v", err)
		}

		{
			uu, err := userRepo.FindByID(user.ID)
			assert(t, err == nil, "Owner.FindByID error: %+v", err)
			assert(t, len(uu.Roles) == 1, "Expected 1 role, got %d", len(uu.Roles))
		}

		{
			users, err := userRepo.Find(&types.UserFilter{Query: "John User Doe"})
			assert(t, err == nil, "Owner.Find error: %+v", err)
			assert(t, len(users) == 1, "Owner.Find: expected 1 user, got %d", len(users))
			assert(t, len(users[0].Roles) == 1, "Owner.Find: expected 1 role, got %d", len(users[0].Roles))
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
