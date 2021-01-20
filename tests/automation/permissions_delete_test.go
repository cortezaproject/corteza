package automation

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	h.allow(types.AutomationRBACResource, "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(rbac.Global().FindRulesByRoleID(permDelRole), 0)

	// Setup a few fake rules for new roke
	h.mockPermissions(
		rbac.AllowRule(permDelRole, types.AutomationRBACResource, "access"),
		rbac.DenyRule(permDelRole, types.AutomationRBACResource, "workflow.create"),
	)

	h.a.Len(p.FindRulesByRoleID(permDelRole), 2)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", permDelRole)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure everything is deleted
	rr, _ := p.FindRulesByRoleID(permDelRole).Filter(func(r *rbac.Rule) (b bool, e error) {
		return r.Access != rbac.Inherit, nil
	})

	h.a.Empty(rr)
}
