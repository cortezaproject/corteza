package compose

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)
	p := rbac.Global()

	// Make sure our user can grant
	h.allow(types.ComponentRbacResource(), "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(rbac.Global().FindRulesByRoleID(permDelRole), 0)

	// Setup a few fake rules for new roke
	h.mockPermissions(
		rbac.DenyRule(permDelRole, types.ComponentRbacResource(), "namespace.create"),
	)

	h.a.Len(p.FindRulesByRoleID(permDelRole), 1)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", permDelRole)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure everything is deleted
	for _, r := range p.FindRulesByRoleID(permDelRole) {
		h.a.True(r.Access == rbac.Inherit)
	}
}
