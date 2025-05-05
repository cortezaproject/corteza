package automation

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	helpers.DenyMe(h, types.ComponentRbacResource(), "workflow.create")

	h.apiInit().
		Get("/permissions/effective").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPermissionsList(t *testing.T) {
	h := newHelper(t)

	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Get("/permissions/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(fmt.Sprintf(`$.response[? @.type=="%s"]`, types.ComponentResourceType))).
		End()
}

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")
	helpers.DenyMe(h, types.ComponentRbacResource(), "workflow.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPermissionsUpdate(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Patch(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"rules":[{"resource":"%s","operation":"workflow.create","access":"allow"}]}`, types.ComponentRbacResource())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPermissionsDelete(t *testing.T) {
	ctx := context.Background()
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(mustFindRulesByRoleID(rbac.Global().FindRulesByRoleID(ctx, permDelRole)), 0)

	// Setup a few fake rules for new role
	helpers.Grant(rbac.AllowRule(permDelRole, types.ComponentRbacResource(), "workflow.create"))

	h.a.Len(mustFindRulesByRoleID(p.FindRulesByRoleID(ctx, permDelRole)), 1)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", permDelRole)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure all rules for this role are deleted
	for _, r := range mustFindRulesByRoleID(p.FindRulesByRoleID(ctx, permDelRole)) {
		h.a.True(r.Access == rbac.Inherit)
	}
}

func TestPermissionsTrace(t *testing.T) {
	h := newHelper(t)

	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Get("/permissions/trace").
		Query("roleID[]", "1").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response`)).
		End()
}

func mustFindRulesByRoleID(rr rbac.RuleSet, err error) rbac.RuleSet {
	if err != nil {
		panic(err)
	}

	return rr
}
