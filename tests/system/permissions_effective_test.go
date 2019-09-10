package system

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "access")
	h.deny(types.SystemPermissionResource, "application.create")

	h.apiInit().
		Get("/permissions/effective").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
