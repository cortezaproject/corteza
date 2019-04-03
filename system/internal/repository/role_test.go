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

func TestRole(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	test.Error(t, db.Transaction(func() error {
		userRepo := User(context.Background(), db)
		user := &types.User{
			Name:     "John Role Doe",
			Username: "johndoe",
		}
		user.GeneratePassword("johndoe")

		{
			u1, err := userRepo.Create(user)
			test.Assert(t, err == nil, "Owner.Create error: %+v", err)
			test.Assert(t, user.ID == u1.ID, "Changes were not stored")
		}

		roleRepo := Role(context.Background(), db)
		role := &types.Role{
			Name: "Test role v1",
		}

		{
			t1, err := roleRepo.Create(role)
			test.Assert(t, err == nil, "Role.Create error: %+v", err)
			test.Assert(t, role.Name == t1.Name, "Changes were not stored")
		}

		{
			role.Name = "Test role v2"
			t1, err := roleRepo.Update(role)
			test.Assert(t, err == nil, "Role.Update error: %+v", err)
			test.Assert(t, role.Name == t1.Name, "Changes were not stored")
		}

		{
			t1, err := roleRepo.FindByID(role.ID)
			test.Assert(t, err == nil, "Role.FindByID error: %+v", err)
			test.Assert(t, role.Name == t1.Name, "Changes were not stored")
		}

		{
			aa, err := roleRepo.Find(&types.RoleFilter{Query: role.Name})
			test.Assert(t, err == nil, "Role.Find error: %+v", err)
			test.Assert(t, len(aa) > 0, "No results found")
		}

		{
			err := roleRepo.ArchiveByID(role.ID)
			test.Assert(t, err == nil, "Role.ArchiveByID error: %+v", err)
		}

		{
			err := roleRepo.UnarchiveByID(role.ID)
			test.Assert(t, err == nil, "Role.UnarchiveByID error: %+v", err)
		}

		{
			err := roleRepo.MemberAddByID(role.ID, user.ID)
			test.Assert(t, err == nil, "Role.MemberAddByID error: %+v", err)
		}

		{
			roles, err := roleRepo.FindByMemberID(user.ID)
			test.Assert(t, err == nil, "Role.FindByMemberID error: %+v", err)
			test.Assert(t, len(roles) > 0, "No results found")
		}

		{
			roles, err := roleRepo.FindByMemberID(0)
			test.Assert(t, err == nil, "Role.FindByMemberID error: %+v", err)
			test.Assert(t, len(roles) == 0, "Results found")
		}

		{
			err := roleRepo.MemberRemoveByID(role.ID, user.ID)
			test.Assert(t, err == nil, "Role.MemberRemoveByID error: %+v", err)
		}

		{
			err := roleRepo.DeleteByID(role.ID)
			test.Assert(t, err == nil, "Role.DeleteByID error: %+v", err)
		}

		{
			err := userRepo.DeleteByID(user.ID)
			test.Assert(t, err == nil, "Owner.DeleteByID error: %+v", err)
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
