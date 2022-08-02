package dal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func Test_dal_crud_issues_compose_record_nok_connection(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionCRUD(h)
	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "nok_connection_connectivity.json")
	module := createModuleFromGenerics(h.secCtx(), t, "ok_module.json", ns.ID, &types.ModuleConfig{
		DAL: types.ModuleConfigDAL{
			ConnectionID: connection.ID,
			Operations:   dal.FullOperations(),
		},
	})

	h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Body(loadRequestFromGenerics(t, "ok_record.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertErrorP(fmt.Sprintf("connection %d has issues", connection.ID))).
		End()
}

func Test_dal_crud_issues_compose_record_nok_model(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionCRUD(h)
	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "ok_connection.json")
	module := createModuleFromGenerics(h.secCtx(), t, "nok_module_sensitivity_level.json", ns.ID, &types.ModuleConfig{
		DAL: types.ModuleConfigDAL{
			ConnectionID: connection.ID,
			Operations:   dal.FullOperations(),
		},
	})

	h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Body(loadRequestFromGenerics(t, "ok_record.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertErrorP(fmt.Sprintf("model %d has issues", module.ID))).
		End()
}

func Test_dal_crud_issues_compose_record_ok(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionCRUD(h)
	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "ok_connection.json")
	module := createModuleFromGenerics(h.secCtx(), t, "ok_module.json", ns.ID, &types.ModuleConfig{
		DAL: types.ModuleConfigDAL{
			ConnectionID: connection.ID,
			Operations:   dal.FullOperations(),
		},
	})

	h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Body(loadRequestFromGenerics(t, "ok_record.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
