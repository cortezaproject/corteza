package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)
	p := service.DefaultPermissions

	// Make sure our user can grant
	h.allow(types.ComposePermissionResource, "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(service.DefaultPermissions.FindRulesByRoleID(permDelRole), 0)

	// Setup a few fake rules for new roke
	h.mockPermissions(
		permissions.AllowRule(permDelRole, types.ComposePermissionResource, "access"),
		permissions.DenyRule(permDelRole, types.ComposePermissionResource, "namespace.create"),
	)

	h.a.Len(p.FindRulesByRoleID(permDelRole), 2)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", permDelRole)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure everything is deleted
	rr, _ := p.FindRulesByRoleID(permDelRole).Filter(func(r *permissions.Rule) (b bool, e error) {
		return r.Access != permissions.Inherit, nil
	})

	h.a.Empty(rr)
}
