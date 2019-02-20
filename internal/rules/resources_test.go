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
		Expect(rules.Inherit, resources.IsAllowed("channel:1", "edit"), "expected inherit")
		Expect(rules.Inherit, resources.IsAllowed("channel:*", "edit"), "expected inherit")
	}

	// allow channel:2 group:2 (default deny, multi=allow)
	{
		resources.GrantByResource(2, "channel:2", []string{"edit", "delete"}, rules.Allow)
		Expect(rules.Inherit, resources.IsAllowed("channel:1", "edit"), "expected error, got nil")
		Expect(rules.Allow, resources.IsAllowed("channel:2", "edit"), "channel:2 edit, expected no error")
		Expect(rules.Allow, resources.IsAllowed("channel:*", "edit"), "channel:* edit, expected no error")
	}

	// list grants for role
	{
		grants, err := resources.ListByResource(2, "channel:2")
		NoError(t, err, "expect no error")
		Assert(t, len(grants) == 2, "expected 2 grants")
		Assert(t, grants[0].RoleID == 2, "expected RoleID == 2, got %v", grants[0].RoleID)
		Assert(t, grants[0].Resource == "channel:2", "expected Resource == channel:2, got %s", grants[0].Resource)
		Assert(t, grants[0].Operation == "delete", "expected Operation == delete, got %s", grants[0].Operation)
		Assert(t, grants[0].Value == rules.Allow, "expected Value == Allow, got %s", grants[0].Value)
	}

	// list all by role
	{
		grants, err := resources.List(2)
		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 2, "expected grants == 2, got %v", len(grants))
	}

	// deny channel:1 group:1 (explicit deny, multi=deny)
	{
		resources.GrantByResource(1, "channel:1", []string{"edit"}, rules.Deny)
		Expect(rules.Deny, resources.IsAllowed("channel:1", "edit"), "expected error, got nil")
		Expect(rules.Allow, resources.IsAllowed("channel:2", "edit"), "channel:2 edit, expected no error")
		Expect(rules.Deny, resources.IsAllowed("channel:*", "edit"), "expected error, got nil")
	}

	// reset (unset=deny)
	{
		resources.GrantByResource(2, "channel:2", []string{"edit", "delete"}, rules.Inherit)
		resources.GrantByResource(1, "channel:1", []string{"edit", "delete"}, rules.Inherit)
		Expect(rules.Inherit, resources.IsAllowed("channel:1", "edit"), "expected error, got nil")
		Expect(rules.Inherit, resources.IsAllowed("channel:*", "edit"), "expected error, got nil")
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

	// delete all by role
	{
		err := resources.Delete(2)
		NoError(t, err, "expected no error")
	}

	// list all by role
	{
		grants, err := resources.List(2)
		NoError(t, err, "expected no error")
		Assert(t, len(grants) == 0, "expected grants == 0, got %v", len(grants))
	}
}
