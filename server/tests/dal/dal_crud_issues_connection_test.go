package dal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dal_crud_issues_connection_connectivity(t *testing.T) {
	t.Skip("needs refactoring")

	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionCRUD(h)

	aux := &dalConnectionRestRsp{}

	rsp := h.apiInit().
		Post("/system/dal/connections/").
		Body(loadRequestFromGenerics(t, "nok_connection_connectivity.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "ok_connection.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "nok_connection_connectivity.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.success.message")).
		End()
}

func Test_dal_crud_issues_connection_sensitivity(t *testing.T) {
	t.Skip("needs refactoring")

	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionCRUD(h)

	aux := &dalConnectionRestRsp{}

	rsp := h.apiInit().
		Post("/system/dal/connections/").
		Body(loadRequestFromGenerics(t, "nok_connection_sensitivity.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "ok_connection.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "nok_connection_sensitivity.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/connections/%d", aux.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.success.message")).
		End()
}
