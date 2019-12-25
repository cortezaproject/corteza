package system

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)
	p := service.DefaultPermissions

	// Make sure our user can grant
	h.allow(types.SystemPermissionResource, "grant")

	// New role.
	permDelRole := h.roleID + 1

	h.a.Len(p.FindRulesByRoleID(permDelRole), 0)

	// Setup a few fake rules for new roke
	h.mockPermissions(
		permissions.AllowRule(permDelRole, types.SystemPermissionResource, "access"),
		permissions.DenyRule(permDelRole, types.SystemPermissionResource, "application.create"),
		permissions.DenyRule(permDelRole, types.SystemPermissionResource, "user.create"),
	)

	h.a.Len(p.FindRulesByRoleID(permDelRole), 3)

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
