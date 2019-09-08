// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/test"
	"github.com/cortezaproject/corteza-server/system/types"
)

func TestUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet("system")

	// Run tests in transaction to maintain DB state.
	test.Error(t, db.Transaction(func() error {
		userRepo := User(context.Background(), db)
		user := &types.User{
			Name:     "John User Doe",
			Username: "johndoe",
			Meta: &types.UserMeta{
				Avatar: "123",
			},
		}
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
		}

		{
			users, _, err := userRepo.Find(types.UserFilter{Query: "John User Doe"})
			test.Assert(t, err == nil, "Owner.Find error: %+v", err)
			test.Assert(t, len(users) == 1, "Owner.Find: expected 1 user, got %d", len(users))
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
