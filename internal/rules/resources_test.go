package rules_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
)

func TestRules(t *testing.T) {
	user := &types.User{ID: 1337}
	ctx := auth.SetIdentityToContext(context.Background(), user)

	db := factory.Database.MustGet()

	roleID := uint64(123456)

	db.Insert("sys_user", user)
	db.Insert("sys_role", types.Role{ID: roleID, Name: fmt.Sprintf("Role %d", roleID)})
	db.Insert("sys_role_member", types.RoleMember{RoleID: roleID, UserID: user.ID})

	Expect := func(expected rules.Access, actual rules.Access, format string, params ...interface{}) {
		Assert(t, expected == actual, format, params...)
	}

	resources := rules.NewResources(ctx, db)

	// delete all for test roleID = 123456
	{
		err := resources.Delete(roleID)
		NoError(t, err, "expected no error")
	}

	// default (unset=deny), forbidden check ...:*
	{
		Expect(rules.Inherit, resources.IsAllowed("messaging:channel:1", "update"), "messaging:channel:1 update - Inherit")
		Expect(rules.Deny, resources.IsAllowed("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
	}

	// allow messaging:channel:2 update,delete
	{
		list := []rules.Rule{
			rules.Rule{Resource: "messaging:channel:2", Operation: "update", Value: rules.Allow},
			rules.Rule{Resource: "messaging:channel:2", Operation: "delete", Value: rules.Allow},
		}
		err := resources.Grant(roleID, list)
		NoError(t, err, "expect no error")

		Expect(rules.Inherit, resources.IsAllowed("messaging:channel:1", "update"), "messaging:channel:1 update - Inherit")
		Expect(rules.Allow, resources.IsAllowed("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
		Expect(rules.Deny, resources.IsAllowed("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
	}

	// list grants for test role
	{
		grants, err := resources.Read(roleID)
		NoError(t, err, "expect no error")
		Assert(t, len(grants) == 2, "expected 2 grants")

		for _, grant := range grants {
			Assert(t, grant.RoleID == roleID, "expected RoleID == 123456, got %v", grant.RoleID)
			Assert(t, grant.Resource == "messaging:channel:2", "expected Resource == messaging:channel:2, got %s", grant.Resource)
			Assert(t, grant.Value == rules.Allow, "expected Value == Allow, got %s", grant.Value)
		}
	}

	// deny messaging:channel:1 update
	{
		list := []rules.Rule{
			rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Deny},
		}
		err := resources.Grant(roleID, list)
		NoError(t, err, "expect no error")

		Expect(rules.Deny, resources.IsAllowed("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
		Expect(rules.Allow, resources.IsAllowed("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
		Expect(rules.Deny, resources.IsAllowed("messaging:channel:*", "update"), "messaging:channel:* update - Deny")
	}

	// reset messaging:channel:1, messaging:channel:2
	{
		list := []rules.Rule{
			rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Inherit},
			rules.Rule{Resource: "messaging:channel:1", Operation: "delete", Value: rules.Inherit},
			rules.Rule{Resource: "messaging:channel:2", Operation: "update", Value: rules.Inherit},
			rules.Rule{Resource: "messaging:channel:2", Operation: "delete", Value: rules.Inherit},
		}
		err := resources.Grant(roleID, list)
		NoError(t, err, "expect no error")

		Expect(rules.Inherit, resources.IsAllowed("messaging:channel:1", "update"), "messaging:channel:1 update - Inherit")
		Expect(rules.Inherit, resources.IsAllowed("messaging:channel:2", "update"), "messaging:channel:2 update - Inherit")
	}

	// [messaging:channel:*,update] - allow, [messaging:channel:1, deny]
	{
		list := []rules.Rule{
			rules.Rule{Resource: "messaging:channel:*", Operation: "update", Value: rules.Allow},
			rules.Rule{Resource: "messaging:channel:1", Operation: "update", Value: rules.Deny},
			rules.Rule{Resource: "messaging:channel:2", Operation: "update"},
			rules.Rule{Resource: "system", Operation: "organisation.create", Value: rules.Allow},
		}
		err := resources.Grant(roleID, list)
		NoError(t, err, "expected no error")

		Expect(rules.Deny, resources.IsAllowed("messaging:channel:1", "update"), "messaging:channel:1 update - Deny")
		Expect(rules.Allow, resources.IsAllowed("messaging:channel:2", "update"), "messaging:channel:2 update - Allow")
	}

	// list all by roleID
	{
		grants, err := resources.Read(roleID)
		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 3, "expected grants == 3, got %v", len(grants))
	}

	// delete all by roleID
	{
		err := resources.Delete(roleID)
		NoError(t, err, "expected no error")
	}

	// list all by roleID
	{
		grants, err := resources.Read(roleID)
		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 0, "expected grants == 0, got %v", len(grants))
	}
}
