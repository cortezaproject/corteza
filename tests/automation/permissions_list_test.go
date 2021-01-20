package automation

import (
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
)

func TestPermissionsList(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Get("/permissions/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response[? @.resource=="automation"]`)).
		End()
}
