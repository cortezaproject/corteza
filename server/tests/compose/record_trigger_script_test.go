package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
)

// Record is sent to the script and returned to the caller,
// scripts must not be triggered on records that can not be read
func TestRecordTriggerScriptForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.DenyMe(h, types.RecordRbacResource(0, 0, 0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d/trigger", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		JSON(`{"script": "some-script", "values": []}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("record.errors.notAllowedToRead")).
		End()
}

// Users that can read the record can trigger scripts on it; corredor is not
// running in tests, so the trigger returns the record untouched
func TestRecordTriggerScript(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d/trigger", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		JSON(`{"script": "some-script", "values": []}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, fmt.Sprintf("%d", record.ID))).
		End()
}
