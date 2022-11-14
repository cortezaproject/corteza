package system

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/id"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	helpers.DenyMe(h, types.ComponentRbacResource(), "user.create")

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

	json := h.apiInit().
		Get("/permissions/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(fmt.Sprintf(`$.response[? @.type=="%s"]`, types.ComponentResourceType))).
		End()

	fmt.Println("json: ", json.Response.Body)
}

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")
	helpers.DenyMe(h, types.ComponentRbacResource(), "user.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPermissionsReadWithFilter(t *testing.T) {
	h := newHelper(t)

	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")
	helpers.DenyMe(h, types.ComponentRbacResource(), "user.create")

	// Specific resource related rules
	testID := id.Next()
	helpers.AllowMe(h, types.UserRbacResource(testID), "read")
	helpers.AllowMe(h, types.UserRbacResource(id.Next()), "update")

	t.Log("all component-level and wildcard rules")
	h.apiInit().
		Getf("/permissions/%d/rules", h.roleID).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 2)).
		End()

	t.Log("no rules for all-users")
	h.apiInit().
		Getf("/permissions/%d/rules", h.roleID).
		Query("resource", "corteza::system:user/*").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 0)).
		End()

	t.Log("1 rule for specific user")
	h.apiInit().
		Getf("/permissions/%d/rules", h.roleID).
		Query("resource", fmt.Sprintf("corteza::system:user/%d", testID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 1)).
		End()
}

func TestPermissionsUpdate(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	h.apiInit().
		Patch(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"rules":[{"resource":"%s","operation":"user.create","access":"allow"}]}`, types.ComponentRbacResource())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(rbac.Global().FindRulesByRoleID(permDelRole), 0)

	// Setup a few fake rules for new role
	helpers.Grant(rbac.AllowRule(permDelRole, types.ComponentRbacResource(), "user.create"))

	h.a.Len(p.FindRulesByRoleID(permDelRole), 1)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", permDelRole)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure all rules for this role are deleted
	for _, r := range p.FindRulesByRoleID(permDelRole) {
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

func TestPermissionsCloneToSingleRole(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	// New role.
	roleS := h.roleID + 1
	roleT := h.roleID + 2

	h.a.Len(rbac.Global().FindRulesByRoleID(roleS), 0)
	h.a.Len(rbac.Global().FindRulesByRoleID(roleT), 0)

	// Set up a few fake rules for new role
	helpers.Grant(rbac.AllowRule(roleS, types.ComponentRbacResource(), "user.create"))

	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.update"))
	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.delete"))

	h.a.Len(p.FindRulesByRoleID(roleS), 1)
	h.a.Len(p.FindRulesByRoleID(roleT), 2)

	h.apiInit().
		Post(fmt.Sprintf("/roles/%d/rules/clone", roleS)).
		Query("cloneToRoleID", strconv.FormatUint(roleT, 10)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure all rules for role S are intact
	h.a.Len(p.FindRulesByRoleID(roleS), 1)
	// Make sure all rules for role T are cloned from role S
	h.a.Len(p.FindRulesByRoleID(roleT), 1)
}

func TestPermissionsCloneToMultipleRole(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	helpers.AllowMe(h, types.ComponentRbacResource(), "grant")

	// New role.
	roleS := h.roleID + 1
	roleT := h.roleID + 2
	roleY := h.roleID + 3

	h.a.Len(rbac.Global().FindRulesByRoleID(roleS), 0)
	h.a.Len(rbac.Global().FindRulesByRoleID(roleT), 0)
	h.a.Len(rbac.Global().FindRulesByRoleID(roleY), 0)

	// Set up a few fake rules for new role
	helpers.Grant(rbac.AllowRule(roleS, types.ComponentRbacResource(), "user.create"))

	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.update"))
	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.delete"))

	helpers.Grant(rbac.AllowRule(roleY, types.ComponentRbacResource(), "user.create"))
	helpers.Grant(rbac.AllowRule(roleY, types.ComponentRbacResource(), "user.update"))
	helpers.Grant(rbac.AllowRule(roleY, types.ComponentRbacResource(), "user.delete"))

	h.a.Len(p.FindRulesByRoleID(roleS), 1)
	h.a.Len(p.FindRulesByRoleID(roleT), 2)
	h.a.Len(p.FindRulesByRoleID(roleY), 3)

	h.apiInit().
		Post(fmt.Sprintf("/roles/%d/rules/clone", roleS)).
		Query("cloneToRoleID", strconv.FormatUint(roleT, 10)).
		Query("cloneToRoleID", strconv.FormatUint(roleY, 10)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure all rules for role S are intact
	h.a.Len(p.FindRulesByRoleID(roleS), 1)
	// Make sure all rules for role T are cloned from role S
	h.a.Len(p.FindRulesByRoleID(roleT), 1)
	// Make sure all rules for role Y are cloned from role S
	h.a.Len(p.FindRulesByRoleID(roleY), 1)
}

func TestPermissionsCloneNotAllowed(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// New role.
	roleS := h.roleID + 1
	roleT := h.roleID + 2

	h.a.Len(rbac.Global().FindRulesByRoleID(roleS), 0)
	h.a.Len(rbac.Global().FindRulesByRoleID(roleT), 0)

	// Set up a few fake rules for new role
	helpers.Grant(rbac.AllowRule(roleS, types.ComponentRbacResource(), "user.create"))

	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.update"))
	helpers.Grant(rbac.AllowRule(roleT, types.ComponentRbacResource(), "user.delete"))

	h.a.Len(p.FindRulesByRoleID(roleS), 1)
	h.a.Len(p.FindRulesByRoleID(roleT), 2)

	h.apiInit().
		Post(fmt.Sprintf("/roles/%d/rules/clone", roleS)).
		Header("Accept", "application/json").
		FormData("cloneToRoleID", strconv.FormatUint(roleT, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("role.errors.notAllowedToCloneRules")).
		End()
}
