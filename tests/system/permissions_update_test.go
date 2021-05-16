package system

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ComponentRbacResource(), "grant")

	h.apiInit().
		Patch(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		JSON(`{"rules":[{"resource":"system","operation":"application.create","access":"allow"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
