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

	db.Insert("sys_user", user)
	var i uint64 = 0
	for i < 5 {
		db.Insert("sys_role", types.Role{ID: i, Name: fmt.Sprintf("Role %d", i)})
		i++
	}
	db.Insert("sys_role_member", types.RoleMember{RoleID: 1, UserID: user.ID})
	db.Insert("sys_role_member", types.RoleMember{RoleID: 2, UserID: user.ID})

	Expect := func(expected rules.Access, actual rules.Access, format string, params ...interface{}) {
		Assert(t, expected == actual, format, params...)
	}

	resources := rules.NewResources(ctx, db)

	// default (unset=deny)
	{
		Expect(rules.Inherit, resources.IsAllowed("channel:1", "update"), "expected inherit")
		Expect(rules.Inherit, resources.IsAllowed("channel:*", "update"), "expected inherit")
	}

	// allow channel:2 group:2 (default deny, multi=allow)
	{
		list := []rules.Rule{
			rules.Rule{Resource: "channel:2", Operation: "update", Value: rules.Allow},
			rules.Rule{Resource: "channel:2", Operation: "delete", Value: rules.Allow},
		}

		resources.Grant(2, list)
		Expect(rules.Inherit, resources.IsAllowed("channel:1", "update"), "expected error, got nil")
		Expect(rules.Allow, resources.IsAllowed("channel:2", "update"), "channel:2 update, expected no error")
		Expect(rules.Allow, resources.IsAllowed("channel:*", "update"), "channel:* update, expected no error")
	}

	// list grants for role 2
	{
		grants, err := resources.List(2)
		NoError(t, err, "expect no error")
		Assert(t, len(grants) == 2, "expected 2 grants")

		for _, grant := range grants {
			Assert(t, grant.RoleID == 2, "expected RoleID == 2, got %v", grant.RoleID)
			Assert(t, grant.Resource == "channel:2", "expected Resource == channel:2, got %s", grant.Resource)
			// Assert(t, grant.Operation == "delete", "expected Operation == delete, got %s", grant.Operation)
			Assert(t, grant.Value == rules.Allow, "expected Value == Allow, got %s", grant.Value)
		}
	}

	// deny channel:1 group:1 (explicit deny, multi=deny)
	{
		list := []rules.Rule{
			rules.Rule{Resource: "channel:1", Operation: "update", Value: rules.Deny},
		}
		resources.Grant(1, list)
		Expect(rules.Deny, resources.IsAllowed("channel:1", "update"), "expected error, got nil")
		Expect(rules.Allow, resources.IsAllowed("channel:2", "update"), "channel:2 update, expected no error")
		Expect(rules.Deny, resources.IsAllowed("channel:*", "update"), "expected error, got nil")
	}

	// reset (unset=deny)
	{
		list1 := []rules.Rule{
			rules.Rule{Resource: "channel:1", Operation: "update", Value: rules.Inherit},
			rules.Rule{Resource: "channel:1", Operation: "delete", Value: rules.Inherit},
		}
		resources.Grant(1, list1)

		list2 := []rules.Rule{
			rules.Rule{Resource: "channel:2", Operation: "update", Value: rules.Inherit},
			rules.Rule{Resource: "channel:2", Operation: "delete", Value: rules.Inherit},
		}
		resources.Grant(2, list2)

		Expect(rules.Inherit, resources.IsAllowed("channel:1", "update"), "expected error, got nil")
		Expect(rules.Inherit, resources.IsAllowed("channel:*", "update"), "expected error, got nil")
	}

	// Grant by roleID
	{
		list := []rules.Rule{
			rules.Rule{Resource: "channel:*", Operation: "update", Value: rules.Allow},
			rules.Rule{Resource: "channel:1", Operation: "update", Value: rules.Deny},
			rules.Rule{Resource: "channel:2", Operation: "update"},
			rules.Rule{Resource: "system", Operation: "organisation.create", Value: rules.Allow},
		}
		err := resources.Grant(2, list)
		NoError(t, err, "expected no error")
	}

	// list all by roleID
	{
		grants, err := resources.List(2)

		fmt.Println(grants)

		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 3, "expected grants == 3, got %v", len(grants))
	}

	// delete all by roleID
	{
		err := resources.Delete(2)
		NoError(t, err, "expected no error")
	}

	// list all by roleID
	{
		grants, err := resources.List(2)
		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 0, "expected grants == 0, got %v", len(grants))
	}
}
