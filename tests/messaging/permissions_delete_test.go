package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsDelete(t *testing.T) {
	h := newHelper(t)

	h.a.Empty(p.FindRulesByRoleID(h.roleID))

	h.allow("messaging", "access")
	h.allow("messaging", "grant")
	h.deny("messaging", "channel.group.create")

	h.a.Len(p.FindRulesByRoleID(h.roleID), 3)

	h.apiInit().
		Delete(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Make sure everything is deleted
	rr, _ := p.FindRulesByRoleID(h.roleID).Filter(func(r *permissions.Rule) (b bool, e error) {
		return r.Access != permissions.Inherit, nil
	})

	h.a.Empty(rr)
}
