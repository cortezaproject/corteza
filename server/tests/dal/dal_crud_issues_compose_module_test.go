package dal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dal_crud_issues_compose_module_missing_sensitivity(t *testing.T) {
	t.Skip("needs refactoring")

	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeModuleCRUD(h)
	ns := h.createNamespace("test")
	aux := &composeModuleRestRsp{}

	rsp := h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/", ns.ID)).
		Body(loadRequestFromGenerics(t, "nok_module_missing_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	rsp = h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "ok_module.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	rsp = h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "nok_module_missing_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	rsp = h.apiInit().
		Delete(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.success.message")).
		End()
}

func Test_dal_crud_issues_compose_module_field_missing_sensitivity(t *testing.T) {
	t.Skip("needs refactoring")

	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeModuleCRUD(h)
	ns := h.createNamespace("test")
	aux := &composeModuleRestRsp{}

	rsp := h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/", ns.ID)).
		Body(loadRequestFromGenerics(t, "nok_module_missing_field_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))
	rsp = h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "ok_module.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	rsp = h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Body(loadRequestFromGenerics(t, "nok_module_missing_field_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	rsp = h.apiInit().
		Delete(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.success.message")).
		End()
}

func Test_dal_crud_issues_compose_module_nok_connection(t *testing.T) {
	t.Skip("needs refactoring")

	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeDalConnectionCRUD(h)

	ns := h.createNamespace("test")
	conn := createConnectionFromGenerics(h.secCtx(), t, "nok_connection_connectivity.json")
	aux := &composeModuleRestRsp{}

	rsp := h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/", ns.ID)).
		Body(loadRequestFromScenarioWithConnection(t, "module.json", conn.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.Len("$.response.issues", 1)).
		End()

	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", conn.ID)).
		Body(loadRequestFromGenerics(t, "ok_connection.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	rsp = h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Body(loadRequestFromScenarioWithConnection(t, "module.json", conn.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.moduleID")).
		Assert(jsonpath.NotPresent("$.response.issues")).
		End()

	rsp = h.apiInit().
		Delete(fmt.Sprintf("/compose/namespace/%d/module/%d", ns.ID, aux.Response.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
