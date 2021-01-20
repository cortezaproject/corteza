package automation

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow(types.AutomationRBACResource, "access")
	h.allow(types.AutomationRBACResource, "grant")
	h.deny(types.AutomationRBACResource, "workflow.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
