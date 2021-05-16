package compose

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"strconv"
	"testing"
)

func TestRecordCreate_batch(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	ns := h.makeNamespace("batch testing namespace")
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.makeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModuleRbacResource(0, 0), "record.create")

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
	h.clearRecords()

	ns := h.makeNamespace("batch testing namespace")
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.makeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModuleRbacResource(0, 0), "record.update")

	record := h.makeRecord(module)
	childRecord := h.makeRecord(childModule, &types.RecordValue{Name: "another_record", Value: strconv.FormatUint(record.ID, 10), Ref: record.ID})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "Some Name"}], "records": [{"refField": "another_record", "set": [{"moduleID": "%d", "recordID": "%d", "values": [{"name": "name", "value": "Another Name"}]}]}]}`, childModule.ID, childRecord.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, strconv.FormatUint(record.ID, 10))).
		Assert(jsonpath.Len(`$.response.records`, 1)).
		Assert(jsonpath.Equal(`$.response.records[0].recordID`, strconv.FormatUint(childRecord.ID, 10))).
		End()
}

func TestRecordDelete_batch(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	ns := h.makeNamespace("batch testing namespace")
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.makeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModuleRbacResource(0, 0), "record.update")
	h.allow(types.ModuleRbacResource(0, 0), "record.delete")

	record := h.makeRecord(module)
	childRecord := h.makeRecord(childModule, &types.RecordValue{Name: "another_record", Value: strconv.FormatUint(record.ID, 10), Ref: record.ID})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "Some Name"}], "records": [{"refField": "another_record", "set": [{"deletedAt": "2020-05-15T15:01:02+02:00", "moduleID": "%d", "recordID": "%d", "values": [{"name": "name", "value": "Another Name"}]}]}]}`, childModule.ID, childRecord.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, strconv.FormatUint(record.ID, 10))).
		Assert(jsonpath.Len(`$.response.records`, 1)).
		Assert(jsonpath.Equal(`$.response.records[0].recordID`, strconv.FormatUint(childRecord.ID, 10))).
		End()

	record = h.lookupRecordByID(module, childRecord.ID)
	h.a.NotNil(record.DeletedAt)
}

func TestRecordMixed_batch(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	ns := h.makeNamespace("batch testing namespace")
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", ns)
	childModule := h.makeRecordModuleWithFieldsOnNs("record testing module child", ns)
	h.allow(types.ModuleRbacResource(0, 0), "record.update")
	h.allow(types.ModuleRbacResource(0, 0), "record.create")

	record := h.makeRecord(module)
	childRecord := h.makeRecord(childModule, &types.RecordValue{Name: "another_record", Value: strconv.FormatUint(record.ID, 10), Ref: record.ID})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
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
