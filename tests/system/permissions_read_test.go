package system

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "access")
	h.allow(types.SystemPermissionResource, "grant")
	h.deny(types.SystemPermissionResource, "application.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
