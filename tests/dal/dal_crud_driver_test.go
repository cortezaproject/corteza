package dal

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dal_crud_driver_list(t *testing.T) {
	h := newHelperT(t)

	h.apiInit().
		Get("/system/dal/drivers/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.GreaterThan("$.response.set", 0)).
		Assert(jsonpath.Present("$.response.set[0].operations")).
		End()
}
