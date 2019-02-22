package service

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	internalAuth "github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"
	systemRepos "github.com/crusttech/crust/system/repository"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermission(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	ctx := context.TODO()

	// Create user for test.
	userRepo := systemRepos.User(ctx, factory.Database.MustGet())
	user := &systemTypes.User{
		Name:     "John Doe",
		Username: "johndoe",
		SatosaID: "1234",
	}
	err := user.GeneratePassword("johndoe")
	NoError(t, err, "expected no error generating password, got %v", err)

	_, err = userRepo.Create(user)
	NoError(t, err, "expected no error creating user, got %v", err)

	// Create role for test and add user
	roleRepo := systemRepos.Role(ctx, factory.Database.MustGet())
	role := &systemTypes.Role{
		Name: "Test role v1",
	}
	_, err = roleRepo.Create(role)
	NoError(t, err, "expected no error creating role, got %v", err)

	err = roleRepo.MemberAddByID(role.ID, user.ID)
	NoError(t, err, "expected no error adding user to role, got %v", err)

	// Set Identity.
	ctx = internalAuth.SetIdentityToContext(ctx, user)

	// Create permission service.
	permissionSvc := Permission().With(ctx)

	// Update rules for test role.
	{
		list := []rules.Rule{
			rules.Rule{Resource: "messaging:channel:*", Operation: "message.update.all", Value: rules.Allow},
			rules.Rule{Resource: "messaging:channel:1", Operation: "message.update.all", Value: rules.Deny},
			rules.Rule{Resource: "messaging:channel:2", Operation: "message.update.all"},
			rules.Rule{Resource: "system", Operation: "organisation.create", Value: rules.Allow},
			rules.Rule{Resource: "system:organisation:*", Operation: "access", Value: rules.Allow},
			rules.Rule{Resource: "messaging:channel", Operation: "message.update.all", Value: rules.Allow},
		}
		_, err := permissionSvc.Update(role.ID, list)
		NoError(t, err, "expected no error, setting rules")
	}

	// Update with invalid roles
	{
		list := []rules.Rule{
			rules.Rule{Resource: "nosystem:channel:*", Operation: "message.update.all", Value: rules.Allow},
		}
		_, err := permissionSvc.Update(role.ID, list)
		Error(t, err, "expected error")

		list = []rules.Rule{
			rules.Rule{Resource: "messaging:noresource:1", Operation: "message.update.all", Value: rules.Deny},
		}
		_, err = permissionSvc.Update(role.ID, list)
		Error(t, err, "expected error")

		list = []rules.Rule{
			rules.Rule{Resource: "messaging:channel:", Operation: "message.update.all"},
		}
		_, err = permissionSvc.Update(role.ID, list)
		Error(t, err, "expected error")

		list = []rules.Rule{
			rules.Rule{Resource: "system:organisation:*", Operation: "invalid", Value: rules.Allow},
		}
		_, err = permissionSvc.Update(role.ID, list)
		Error(t, err, "expected error")
	}

	// Read rules for test role.
	{
		ret, err := permissionSvc.Read(role.ID)
		NoError(t, err, "expected no error, setting rules")

		rules := ret.([]rules.Rule)

		Assert(t, len(rules) == 5, "expected len(rules) == 5, got %v", len(rules))
	}

	// Delete rules for test role.
	{
		_, err := permissionSvc.Delete(role.ID)
		NoError(t, err, "expected no error, setting rules")
	}

	// Read rules for test role.
	{
		ret, err := permissionSvc.Read(role.ID)
		NoError(t, err, "expected no error, setting rules")

		rules := ret.([]rules.Rule)

		Assert(t, len(rules) == 0, "expected len(rules) == 0, got %v", len(rules))
	}

	// List defined permissions.
	{
		ret, err := permissionSvc.List()
		NoError(t, err, "expected no error, setting rules")

		perms := ret.([]systemTypes.Permission)

		Assert(t, len(perms) > 0, "expected len(rules) > 0, got %v", len(perms))
	}
}
