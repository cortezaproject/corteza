package service

import (
	"context"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermissions(t *testing.T) {
	ctx := context.TODO()

	// Create user with role and add it to context.
	userSvc := systemService.User().With(ctx)
	user := &systemTypes.User{
		Name:     "John Doe",
		Username: "johndoe",
		SatosaID: "1234",
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

	// Generate services.
	channelSvc := (&channel{
		usr: systemService.User(),
		evl: Event(),
		prm: Permissions(),
	}).With(ctx)

	permissionsSvc := Permissions().With(ctx)
	systemPermissionSvc := systemService.Permissions().With(ctx)

	// Test `access` to messaging service.
	ret := permissionsSvc.CanAccessMessaging()
	Assert(t, ret == false, "expected CanAccessMessaging == false, got %v", ret)

	// Add `access` to messaging service.
	list := []rules.Rule{
		rules.Rule{Resource: "messaging", Operation: "access", Value: rules.Allow},
	}
	_, err = systemPermissionSvc.Update(role.ID, list)
	NoError(t, err, "expected no error, got %v", err)

	// Test `access` to messaging service.
	ret = permissionsSvc.CanAccessMessaging()
	Assert(t, ret == true, "expected CanAccessMessaging == true, got %v", ret)

	// Create test channel.
	ch := &types.Channel{
		Name:  "TestChan",
		Topic: "No topic",
	}
	ch, err = channelSvc.Create(ch)
	NoError(t, err, "expected no error, got %v", err)

	// @Todo: add permission for create channel and test it.

	// Test CanRead permissions. [1 - no permission, 2 - allow]
	{
		ret = permissionsSvc.CanRead(ch)
		Assert(t, ret == false, "expected CanRead == false, got %v")

		// Add [messaging:channel:*, read, allow]
		list = []rules.Rule{
			rules.Rule{Resource: "messaging:channel:*", Operation: "read", Value: rules.Allow},
		}
		_, err = systemPermissionSvc.Update(role.ID, list)
		NoError(t, err, "expected no error, got %v", err)

		ret = permissionsSvc.CanRead(ch)
		Assert(t, ret == true, "expected CanRead == true, got %v")
	}

	// Test CanJoin permissions [1 - deny, 2 - allow for resourceID]
	{
		// Add [messaging:channel:*, join, deny]
		list = []rules.Rule{
			rules.Rule{Resource: "messaging:channel:*", Operation: "join", Value: rules.Deny},
		}
		_, err = systemPermissionSvc.Update(role.ID, list)
		NoError(t, err, "expected no error, got %v", err)

		ret = permissionsSvc.CanJoin(ch)
		Assert(t, ret == false, "expected CanJoin == false, got %v")

		// Add [messaging:channel:ID, join, allow]
		list = []rules.Rule{
			rules.Rule{Resource: ch.Resource().String(), Operation: "join", Value: rules.Allow},
		}
		_, err = systemPermissionSvc.Update(role.ID, list)
		NoError(t, err, "expected no error, got %v", err)

		ret = permissionsSvc.CanJoin(ch)
		Assert(t, ret == true, "expected CanJoin == true, got %v")
	}

	// Remove test channel.
	channelSvc.Delete(ch.ID)
}
