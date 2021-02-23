package compose

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	h.deny(types.ComposeRBACResource, "namespace.create")

	h.apiInit().
		Get("/permissions/effective").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
