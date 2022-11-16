package apigw

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

func Test_route_reload(t *testing.T) {
	var (
		ctx, h, s = setupScenario(t)
	)

	r, err := store.LookupApigwRouteByEndpoint(ctx, s, "/test")
	h.a.NoError(err)

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Body("60").
		End()

	helpers.AllowMe(h, r.RbacResource(), "update")
	h.apiInit().
		Put(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		FormData("endpoint", "/test").
		FormData("method", "GET").
		FormData("enabled", "false").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusNotFound).
		End()

	h.apiInit().
		Put(fmt.Sprintf("/apigw/route/%d", r.ID)).
		Header("Accept", "application/json").
		FormData("endpoint", "/test").
		FormData("method", "GET").
		FormData("enabled", "true").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Body("60").
		End()
}
