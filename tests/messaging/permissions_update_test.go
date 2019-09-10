package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow("messaging", "access")
	h.allow("messaging", "grant")
	h.deny("messaging", "channel.group.create")

	h.apiInit().
		Patch(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		JSON(`{"rules":[{"resource":"messaging","operation":"channel.group.create","access":"allow"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
