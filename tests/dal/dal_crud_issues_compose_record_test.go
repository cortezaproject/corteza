package dal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func Test_dal_crud_issues_compose_record_nok_connection(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connection.create")
	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connections.search")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "read")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "update")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "delete")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "delete")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read")

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "nok_connection_connectivity.json")
	module := createModuleFromGenerics(h.secCtx(), t, "ok_module.json", ns.ID, &types.ModelConfig{
		ConnectionID: connection.ID,
		Capabilities: capabilities.FullCapabilities(),
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

	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connection.create")
	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connections.search")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "read")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "update")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "delete")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "delete")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read")

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "ok_connection.json")
	module := createModuleFromGenerics(h.secCtx(), t, "nok_module_sensitivity_level.json", ns.ID, &types.ModelConfig{
		ConnectionID: connection.ID,
		Capabilities: capabilities.FullCapabilities(),
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

	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connection.create")
	helpers.AllowMe(h, systemTypes.ComponentRbacResource(), "dal-connections.search")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "read")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "update")
	helpers.AllowMe(h, systemTypes.DalConnectionRbacResource(0), "delete")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "delete")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read")

	ns := h.createNamespace("test")

	connection := createConnectionFromGenerics(h.secCtx(), t, "ok_connection.json")
	module := createModuleFromGenerics(h.secCtx(), t, "ok_module.json", ns.ID, &types.ModelConfig{
		ConnectionID: connection.ID,
		Capabilities: capabilities.FullCapabilities(),
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
