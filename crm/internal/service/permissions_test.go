// +build integration

package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/crm/internal/repository"
	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermissions(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing", true)

	// Create user with role and add it to context.
	userSvc := systemService.TestUser(t, ctx)
	user := &systemTypes.User{
		Name:     "John Crm Doe",
		Username: "johndoe",
		SatosaID: "12345",
	}
	err := user.GeneratePassword("johndoe")
	NoError(t, err, "expected no error generating password, got %v", err)

	_, err = userSvc.Create(user, nil, "")
	NoError(t, err, "expected no error creating user, got %v", err)

	// Set Identity.
	ctx = auth.SetIdentityToContext(ctx, user)

	roleSvc := systemService.TestRole(t, ctx)
	role := &systemTypes.Role{
		Name: "Test role v1",
	}
	role, err = roleSvc.Create(role)
	NoError(t, err, "expected no error creating role, got %v", err)

	err = roleSvc.MemberAdd(role.ID, user.ID)
	NoError(t, err, "expected no error adding user to role, got %v", err)

	// Insert `grant` permission for `compose`.
	{
		db := repository.DB(ctx)
		resources := rules.NewResources(ctx, db)

		list := []rules.Rule{
			rules.Rule{Resource: types.PermissionResource, Operation: "grant", Value: rules.Allow},
		}

		err := resources.Grant(role.ID, list)
		NoError(t, err, "expected no error, got %v", err)
	}

	// Generate services.
	cleanContext := auth.SetIdentityToContext(context.Background(), user)
	permissionsSvc := Permissions().With(cleanContext)
	systemRulesSvc := systemService.TestRules(t, ctx)

	// Test `access` to compose service.
	ret := permissionsSvc.CanAccess()
	Assert(t, ret == false, "expected CanAccess == false, got %v", ret)

	// Add `access` to compose service.
	list := []rules.Rule{
		rules.Rule{Resource: types.PermissionResource, Operation: "access", Value: rules.Allow},
	}
	_, err = systemRulesSvc.Update(role.ID, list)
	NoError(t, err, "expected no error, got %v", err)

	// Test `access` to compose service.
	ret = permissionsSvc.CanAccess()
	Assert(t, ret == true, "expected CanAccess == true, got %v", ret)
}
