package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"

	"github.com/crusttech/crust/crm/internal/repository"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermissions(t *testing.T) {
	ctx := context.TODO()

	// Create user with role and add it to context.
	userSvc := systemService.User().With(ctx)
	user := &systemTypes.User{
		Name:     "John Crm Doe",
		Username: "johndoe",
		SatosaID: "12345",
	}
	err := user.GeneratePassword("johndoe")
	NoError(t, err, "expected no error generating password, got %v", err)

	_, err = userSvc.Create(user)
	NoError(t, err, "expected no error creating user, got %v", err)

	roleSvc := systemService.Role().With(ctx)
	role := &systemTypes.Role{
		Name: "Test role v1",
	}
	role, err = roleSvc.Create(role)
	NoError(t, err, "expected no error creating role, got %v", err)

	err = roleSvc.MemberAdd(role.ID, user.ID)
	NoError(t, err, "expected no error adding user to role, got %v", err)

	// Set Identity.
	ctx = auth.SetIdentityToContext(ctx, user)

	// Insert `grant` permission for `compose`.
	{
		db := repository.DB(ctx)
		resources := rules.NewResources(ctx, db)

		list := []rules.Rule{
			rules.Rule{Resource: "compose", Operation: "grant", Value: rules.Allow},
		}

		err := resources.Grant(role.ID, list)
		NoError(t, err, "expected no error, got %v", err)
	}

	// Generate services.
	permissionsSvc := Permissions().With(ctx)
	systemRulesSvc := systemService.Rules().With(ctx)

	// Test `access` to compose service.
	ret := permissionsSvc.CanAccess()
	Assert(t, ret == true, "expected CanAccess == true, got %v", ret)

	// Add `access` to compose service.
	list := []rules.Rule{
		rules.Rule{Resource: "compose", Operation: "access", Value: rules.Allow},
	}
	_, err = systemRulesSvc.Update(role.ID, list)
	NoError(t, err, "expected no error, got %v", err)

	// Test `access` to compose service.
	ret = permissionsSvc.CanAccess()
	Assert(t, ret == true, "expected CanAccess == true, got %v", ret)
}
