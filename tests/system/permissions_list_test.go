package system

import (
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsList(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Get("/permissions/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response[? @.resource=="system"]`)).
		End()
}
