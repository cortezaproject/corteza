package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow("compose", "grant")

	h.apiInit().
		Patch(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		JSON(`{"rules":[{"resource":"compose","operation":"namespace.create","access":"allow"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
