package compose

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ComposePermissionResource, "access")
	h.allow(types.ComposePermissionResource, "grant")
	h.deny(types.ComposePermissionResource, "namespace.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
