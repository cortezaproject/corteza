package compose

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestRecordCreate_batch(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("batch testing namespace")
	module := h.repoMakeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.repoMakeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.create")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [], "records": [{"refField": "another_record", "set": [{"moduleID": "%d", "values": []}]}]}`, childModule.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.recordID`)).
		Assert(jsonpath.Len(`$.response.records`, 1)).
		Assert(jsonpath.Present(`$.response.records[0].recordID`)).
		End()
}

func TestRecordUpdate_batch(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("batch testing namespace")
	module := h.repoMakeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.repoMakeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.update")

	record := h.repoMakeRecord(module)
	childRecord := h.repoMakeRecord(childModule, &types.RecordValue{Name: "another_record", Value: strconv.FormatUint(record.ID, 10), Ref: record.ID})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "Some Name"}], "records": [{"refField": "another_record", "set": [{"moduleID": "%d", "recordID": "%d", "values": [{"name": "name", "value": "Another Name"}]}]}]}`, childModule.ID, childRecord.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, strconv.FormatUint(record.ID, 10))).
		Assert(jsonpath.Len(`$.response.records`, 1)).
		Assert(jsonpath.Equal(`$.response.records[0].recordID`, strconv.FormatUint(childRecord.ID, 10))).
		End()
}

func TestRecordMixed_batch(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("batch testing namespace")
	module := h.repoMakeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.repoMakeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.update")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.create")

	record := h.repoMakeRecord(module)
	childRecord := h.repoMakeRecord(childModule, &types.RecordValue{Name: "another_record", Value: strconv.FormatUint(record.ID, 10), Ref: record.ID})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "Some Name"}], "records": [{"refField": "another_record", "set": [{"moduleID": "%d", "values": [{"name": "name", "value": "Added Name"}]},{"moduleID": "%d", "recordID": "%d", "values": [{"name": "name", "value": "Another Name"}]}]}]}`, childModule.ID, childModule.ID, childRecord.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, strconv.FormatUint(record.ID, 10))).
		Assert(jsonpath.Len(`$.response.records`, 2)).
		Assert(jsonpath.Present(`$.response.records[0].recordID`)).
		Assert(jsonpath.Equal(`$.response.records[1].recordID`, strconv.FormatUint(childRecord.ID, 10))).
		End()
}
