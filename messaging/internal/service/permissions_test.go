package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"

	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermissions(t *testing.T) {
	// Create test user and role.
	user := &systemTypes.User{ID: 1337}
	role := &systemTypes.Role{ID: 1234567, Name: "Admins"}

	// Write user to context.
	ctx := auth.SetIdentityToContext(context.Background(), user)

	// Connect do DB.
	db := factory.Database.MustGet()

	// Run test with savepoint.
	err := func() error {
		db.Exec("SAVEPOINT permissions_test")

		db.Insert("sys_user", user)
		db.Insert("sys_role_member", systemTypes.RoleMember{RoleID: role.ID, UserID: user.ID})

		// Insert `grant` permission for `messaging`.
		{
			db := repository.DB(ctx)
			resources := rules.NewResources(ctx, db)

			list := []rules.Rule{
				rules.Rule{Resource: "messaging", Operation: "grant", Value: rules.Allow},
			}

			err := resources.Grant(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)
		}

		// Generate services.
		channelSvc := (&channel{
			usr: systemService.User(),
			evl: Event(),
			prm: Permissions(),
		}).With(ctx)

		permissionsSvc := Permissions().With(ctx)
		systemRulesSvc := systemService.Rules().With(ctx)

		// Remove `access` to messaging service.
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging", Operation: "access", Value: rules.Deny},
			}
			_, err := systemRulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)

			// Test `access` to messaging service.
			ret := permissionsSvc.CanAccess()
			Assert(t, ret == false, "expected CanAccess == false, got %v", ret)
		}

		// Add `access` to messaging service.
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging", Operation: "access", Value: rules.Allow},
			}
			_, err := systemRulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)

			// Test `access` to messaging service.
			ret := permissionsSvc.CanAccess()
			Assert(t, ret == true, "expected CanAccess == true, got %v", ret)
		}

		// Create test channel.
		ch := &types.Channel{
			Name:  "TestChan",
			Topic: "No topic",
		}
		ch, err := channelSvc.Create(ch)
		NoError(t, err, "expected no error, got %v", err)

		// @Todo: add permission for create channel and test it.

		// Test CanReadChannel permissions. [1 - allow, 2 no permission]
		{
			ret := permissionsSvc.CanReadChannel(ch)
			Assert(t, ret == true, "expected CanReadChannel == true, got %v", ret)

			// Add [messaging:channel:*, read, deny]
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:*", Operation: "read", Value: rules.Deny},
			}
			_, err = systemRulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)

			ret = permissionsSvc.CanReadChannel(ch)
			Assert(t, ret == false, "expected CanReadChannel == false, got %v", ret)
		}

		// Test CanJoinChannel permissions [1 - deny, 2 - allow for resourceID]
		{
			// Add [messaging:channel:*, join, deny]
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:*", Operation: "join", Value: rules.Deny},
			}
			_, err = systemRulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)

			ret := permissionsSvc.CanJoinChannel(ch)
			Assert(t, ret == false, "expected CanJoinChannel == false, got %v")

			// Add [messaging:channel:ID, join, allow]
			list = []rules.Rule{
				rules.Rule{Resource: ch.Resource().String(), Operation: "join", Value: rules.Allow},
			}
			_, err = systemRulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %v", err)

			ret = permissionsSvc.CanJoinChannel(ch)
			Assert(t, ret == true, "expected CanJoinChannel == true, got %v")
		}

		// Remove test channel.
		channelSvc.Delete(ch.ID)
		return errors.New("Rollback")
	}()
	if err != nil {
		db.Exec("ROLLBACK TO SAVEPOINT permissions_test")
	}
}
