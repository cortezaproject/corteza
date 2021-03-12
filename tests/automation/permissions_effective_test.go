package automation

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	h.allow(types.AutomationRBACResource, "access")
	h.deny(types.AutomationRBACResource, "workflow.create")

	h.apiInit().
		Get("/permissions/effective").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
