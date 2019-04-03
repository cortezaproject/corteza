// +build integration

package rules_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
)

func TestRules(t *testing.T) {
	Expect := func(expected rules.Access, actual rules.Access, format string, params ...interface{}) {
		Assert(t, expected == actual, format, params...)
	}

	// Create test user and role.
	user := &types.User{ID: 1337}
	role := &types.Role{ID: 123456, Name: "Test role"}

	// Write user to context.
	ctx := auth.SetIdentityToContext(context.Background(), user)

	// Connect do DB.
	db := factory.Database.MustGet()

	// Create resources interface.
	resources := rules.NewResources(ctx, db)

	// Run tests in transaction to maintain DB state.
	Error(t, db.Transaction(func() error {
		db.Exec("DELETE FROM sys_rules WHERE 1=1")
		db.Insert("sys_user", user)
		db.Insert("sys_role", role)
		db.Insert("sys_role_member", types.RoleMember{RoleID: role.ID, UserID: user.ID})

		// delete all for test roleID = 123456
		{
			err := resources.Delete(role.ID)
			NoError(t, err, "expected no error, got %+v", err)
		}

		// check that testing context allows anything
		{
			ctxAdmin := context.WithValue(ctx, "testing", true)
			Expect(rules.Allow, resources.With(ctxAdmin, db).Check("crm", "anything"), "testing context should allow anything")
		}

		// default (unset=deny), forbidden check ...:*
		{
			Expect(rules.Deny, resources.Check("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
			Expect(rules.Deny, resources.Check("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
		}

		// allow messaging:channel:2 update,delete
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:2", Operation: "update", Value: rules.Allow},
				rules.Rule{Resource: "messaging:channel:2", Operation: "delete", Value: rules.Allow},
			}
			err := resources.Grant(role.ID, list)
			NoError(t, err, "expect no error, got %+v", err)

			Expect(rules.Deny, resources.Check("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
			Expect(rules.Allow, resources.Check("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
			Expect(rules.Deny, resources.Check("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
		}

		// list grants for test role
		{
			grants, err := resources.Read(role.ID)
			NoError(t, err, "expect no error, got %+v", err)
			Assert(t, len(grants) == 2, "expected 2 grants, got %v", len(grants))

			for _, grant := range grants {
				Assert(t, grant.RoleID == role.ID, "expected RoleID == 123456, got %v", grant.RoleID)
				Assert(t, grant.Resource == "messaging:channel:2", "expected Resource == messaging:channel:2, got %s", grant.Resource)
				Assert(t, grant.Value == rules.Allow, "expected Value == Allow, got %s", grant.Value)
			}
		}

		// deny messaging:channel:1 update
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Deny},
			}
			err := resources.Grant(role.ID, list)
			NoError(t, err, "expect no error, got %+v", err)

			Expect(rules.Deny, resources.Check("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
			Expect(rules.Allow, resources.Check("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
			Expect(rules.Deny, resources.Check("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
		}

		// reset messaging:channel:1, messaging:channel:2
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Inherit},
				rules.Rule{Resource: "messaging:channel:1", Operation: "delete", Value: rules.Inherit},
				rules.Rule{Resource: "messaging:channel:2", Operation: "update", Value: rules.Inherit},
				rules.Rule{Resource: "messaging:channel:2", Operation: "delete", Value: rules.Inherit},
			}
			err := resources.Grant(role.ID, list)
			NoError(t, err, "expect no error, got %+v", err)

			Expect(rules.Deny, resources.Check("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
			Expect(rules.Deny, resources.Check("messaging:channel:2", "update"), "messaging:channel:2 update - Deny")
		}

		// [messaging:channel:*,update] - allow, [messaging:channel:1, deny]
		{
			list := []rules.Rule{
				rules.Rule{Resource: "messaging:channel:*", Operation: "update", Value: rules.Allow},
				rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Deny},
				rules.Rule{Resource: "messaging:channel:2", Operation: "update"},
				rules.Rule{Resource: "system", Operation: "organisation.create", Value: rules.Allow},
			}
			err := resources.Grant(role.ID, list)
			NoError(t, err, "expect no error, got %+v", err)

			Expect(rules.Deny, resources.Check("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
			Expect(rules.Allow, resources.Check("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
		}

		// list all by roleID
		{
			grants, err := resources.Read(role.ID)
			NoError(t, err, "expect no error, got %+v", err)
			Assert(t, len(grants) == 3, "expected grants == 3, got %v", len(grants))
		}

		// delete all by roleID
		{
			err := resources.Delete(role.ID)
			NoError(t, err, "expect no error, got %+v", err)
		}

		// list all by roleID
		{
			grants, err := resources.Read(role.ID)
			NoError(t, err, "expect no error, got %+v", err)
			Assert(t, len(grants) == 0, "expected grants == 0, got %v", len(grants))
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
