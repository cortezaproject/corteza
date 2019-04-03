// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/types"
)

func TestUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Run tests in transaction to maintain DB state.
	test.Error(t, db.Transaction(func() error {
		userRepo := User(context.Background(), db)
		user := &types.User{
			Name:     "John User Doe",
			Username: "johndoe",
			SatosaID: "1234",
			Meta: &types.UserMeta{
				Avatar: "123",
			},
		}
		user.GeneratePassword("johndoe")

		{
			uu, err := userRepo.Create(user)
			test.Assert(t, err == nil, "Owner.Create error: %+v", err)
			test.Assert(t, user.ID == uu.ID, "Changes were not stored")
		}

		roleRepo := Role(context.Background(), db)
		role := &types.Role{
			Name: "Test role v1",
		}

		{
			t1, err := roleRepo.Create(role)
			test.Assert(t, err == nil, "Role.Create error: %+v", err)
			test.Assert(t, role.Name == t1.Name, "Changes were not stored")

			err = roleRepo.MemberAddByID(t1.ID, user.ID)
			test.Assert(t, err == nil, "Role.MemberAddByID error: %+v", err)
		}

		{
			uu, err := userRepo.FindByID(user.ID)
			test.Assert(t, err == nil, "Owner.FindByID error: %+v", err)
			test.Assert(t, uu.Meta.Avatar == "123", "Expected avatar to be '123', got '%s'", uu.Meta.Avatar)
			test.Assert(t, len(uu.Roles) == 1, "Expected 1 role, got %d", len(uu.Roles))
		}

		{
			users, err := userRepo.Find(&types.UserFilter{Query: "John User Doe"})
			test.Assert(t, err == nil, "Owner.Find error: %+v", err)
			test.Assert(t, len(users) == 1, "Owner.Find: expected 1 user, got %d", len(users))
			test.Assert(t, len(users[0].Roles) == 1, "Owner.Find: expected 1 role, got %d", len(users[0].Roles))
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
