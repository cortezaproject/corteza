package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	internalAuth "github.com/crusttech/crust/internal/auth"
	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
)

func TestRules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	// Create test user and role.
	user := &types.User{ID: 1337}
	role := &types.Role{ID: 123456, Name: "Test role"}

	// Write user to context.
	ctx := internalAuth.SetIdentityToContext(context.Background(), user)

	// Connect do DB.
	db := factory.Database.MustGet()

	// Create resources interface.
	resources := internalRules.NewResources(ctx, db)

	// Run tests in transaction to maintain DB state.
	Error(t, db.Transaction(func() error {
		db.Delete("sys_rules", "1=1")
		db.Insert("sys_user", user)
		db.Insert("sys_role", role)
		db.Insert("sys_role_member", types.RoleMember{RoleID: role.ID, UserID: user.ID})

		// delete all for test roleID = 123456
		{
			err := resources.Delete(role.ID)
			NoError(t, err, "expected no error, got %+v", err)
		}

		// Create rules service.
		rulesSvc := Rules().With(ctx)

		// Update rules for test role, with error.
		{
			list := []internalRules.Rule{
				internalRules.Rule{Resource: "messaging:channel:1", Operation: "message.update.all", Value: internalRules.Allow},
			}
			_, err := rulesSvc.Update(role.ID, list)
			Error(t, err, "expected error == No Allow rule for messaging")
		}

		// Insert `grant` permission for `messaging` and `system`.
		{
			list := []internalRules.Rule{
				internalRules.Rule{Resource: "system", Operation: "grant", Value: internalRules.Allow},
				internalRules.Rule{Resource: "messaging", Operation: "grant", Value: internalRules.Allow},
			}

			err := resources.Grant(role.ID, list)
			NoError(t, err, "expected no error, got %v+", err)
		}

		// List possible permissions with `messaging` and `system` grants.
		{
			ret, err := rulesSvc.List()
			NoError(t, err, "expected no error, got %+v", err)

			perms := ret.([]types.Permission)

			Assert(t, len(perms) > 0, "expected len(rules) > 0, got %v", len(perms))
		}

		// Update rules for test role.
		{
			list := []internalRules.Rule{
				internalRules.Rule{Resource: "messaging:channel:*", Operation: "message.update.all", Value: internalRules.Allow},
				internalRules.Rule{Resource: "messaging:channel:1", Operation: "message.update.all", Value: internalRules.Deny},
				internalRules.Rule{Resource: "messaging:channel:2", Operation: "message.update.all"},
				internalRules.Rule{Resource: "system", Operation: "organisation.create", Value: internalRules.Allow},
				internalRules.Rule{Resource: "system:organisation:*", Operation: "access", Value: internalRules.Allow},
				internalRules.Rule{Resource: "messaging:channel", Operation: "message.update.all", Value: internalRules.Allow},
			}
			_, err := rulesSvc.Update(role.ID, list)
			NoError(t, err, "expected no error, got %+v", err)
		}

		// Update with invalid roles
		{
			list := []internalRules.Rule{
				internalRules.Rule{Resource: "nosystem:channel:*", Operation: "message.update.all", Value: internalRules.Allow},
			}
			_, err := rulesSvc.Update(role.ID, list)
			Error(t, err, "expected error")

			list = []internalRules.Rule{
				internalRules.Rule{Resource: "messaging:noresource:1", Operation: "message.update.all", Value: internalRules.Deny},
			}
			_, err = rulesSvc.Update(role.ID, list)
			Error(t, err, "expected error")

			list = []internalRules.Rule{
				internalRules.Rule{Resource: "messaging:channel:", Operation: "message.update.all"},
			}
			_, err = rulesSvc.Update(role.ID, list)
			Error(t, err, "expected error")

			list = []internalRules.Rule{
				internalRules.Rule{Resource: "system:organisation:*", Operation: "invalid", Value: internalRules.Allow},
			}
			_, err = rulesSvc.Update(role.ID, list)
			Error(t, err, "expected error")
		}

		// Read rules for test role.
		{
			ret, err := rulesSvc.Read(role.ID)
			NoError(t, err, "expected no error, got %+v", err)

			rules := ret.([]internalRules.Rule)

			Assert(t, len(rules) == 7, "expected len(rules) == 7, got %v", len(rules))
		}

		// Delete rules for test role.
		{
			_, err := rulesSvc.Delete(role.ID)
			NoError(t, err, "expected no error, got %+v", err)
		}

		// Read rules for test role.
		{
			ret, err := rulesSvc.Read(role.ID)
			NoError(t, err, "expected no error, got %+v", err)

			rules := ret.([]internalRules.Rule)

			Assert(t, len(rules) == 0, "expected len(rules) == 0, got %v", len(rules))
		}

		// List possible permissions with no grants.
		{
			ret, err := rulesSvc.List()
			NoError(t, err, "expected no error, got %+v", err)

			perms := ret.([]types.Permission)

			Assert(t, len(perms) == 0, "expected len(rules) == 0, got %v", len(perms))
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
