package compose

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/service/values"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	systemService "github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (h helper) clearRecords() {
	h.clearNamespaces()
	h.clearModules()
	h.noError(truncateRecords(context.Background()))
}

type (
	rImportSession struct {
		Response struct {
			SessionID string `json:"sessionID"`
		} `json:"response"`
	}
)

func (h helper) makeRecordModuleWithFieldsOnNs(name string, namespace *types.Namespace, ff ...*types.ModuleField) *types.Module {
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")

	if len(ff) == 0 {
		// Default fields
		ff = types.ModuleFieldSet{
			&types.ModuleField{
				Name: "name",
			},
			&types.ModuleField{
				Name: "email",
			},
			&types.ModuleField{
				Name:  "options",
				Multi: true,
			},
			&types.ModuleField{
				Name: "description",
			},
			&types.ModuleField{
				Name: "another_record",
				Kind: "Record",
			},
		}
	}

	return h.makeModule(namespace, name, ff...)
}

func (h helper) repoMakeRecordModuleWithFields(name string, ff ...*types.ModuleField) *types.Module {
	namespace := h.makeNamespace("record testing namespace")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read", "record.value.update")

	if len(ff) == 0 {
		// Default fields
		ff = types.ModuleFieldSet{
			&types.ModuleField{
				Name: "name",
			},
			&types.ModuleField{
				Name: "email",
			},
			&types.ModuleField{
				Name:  "options",
				Multi: true,
			},
			&types.ModuleField{
				Name: "description",
			},
			&types.ModuleField{
				Name: "another_record",
				Kind: "Record",
			},
		}
	}

	return h.makeModule(namespace, name, ff...)
}

func (h helper) repoMakeRecordModuleWithFieldsRequired(name string, ff ...*types.ModuleField) *types.Module {
	namespace := h.makeNamespace("record testing namespace")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read", "record.value.update")

	if len(ff) == 0 {
		// Default fields
		ff = types.ModuleFieldSet{
			&types.ModuleField{
				Name:     "name",
				Required: true,
			},
			&types.ModuleField{
				Name: "email",
			},
			&types.ModuleField{
				Name:  "options",
				Multi: true,
			},
			&types.ModuleField{
				Name: "description",
			},
			&types.ModuleField{
				Name: "another_record",
				Kind: "Record",
			},
		}
	}

	return h.makeModule(namespace, name, ff...)
}

func (h helper) makeRecord(module *types.Module, rvs ...*types.RecordValue) *types.Record {
	recID := id.Next()
	for _, rv := range rvs {
		rv.RecordID = recID
	}

	rec := &types.Record{
		ID:          recID,
		CreatedAt:   time.Now(),
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
		// Passing the current owner in here since the tests (who care about this)
		// rely on it being set to something valid.
		OwnedBy: h.cUser.ID,

		// We are directly storing the record values here, so ensure
		// everything is formatted in the same manner as it would be
		// when stored through the service
		Values: values.Formatter().Run(module, rvs),
	}
	rec.SetModule(module)

	h.noError(dalutils.ComposeRecordCreate(context.Background(), defDal, module, rec))

	return rec
}

func (h helper) lookupRecordByID(module *types.Module, ID uint64) *types.Record {
	res, err := dalutils.ComposeRecordsFind(context.Background(), defDal, module, ID)
	h.noError(err)
	return res
}

func TestRecordRead(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recordID`, fmt.Sprintf("%d", record.ID))).
		End()
}

func TestRecordList(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")

	h.makeRecord(module)
	h.makeRecord(module)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestRecordListWithTotal(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")

	for i := 0; i < 7; i++ {
		h.makeRecord(module, &types.RecordValue{Name: "name", Value: fmt.Sprintf("%d", i+1)})
	}

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Query("incTotal", "true").
		Query("sort", "name DESC").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filter.total`, float64(7))).
		End()
}

func TestRecordListWithPaginationAndSorting(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, module.RbacResource(), "records.search")

	var aux = struct {
		Response struct {
			Filter struct {
				NextPage       *string
				PrevPage       *string
				PageNavigation []struct {
					Page   int
					Items  int
					Cursor *string
				}
			}
		}
	}{}

	for i := 0; i < 7; i++ {
		h.makeRecord(module, &types.RecordValue{Name: "name", Value: fmt.Sprintf("%d", i+1)})
	}

	// 1st page
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Query("incTotal", "true").
		Query("incPageNavigation", "true").
		Query("limit", "2").
		Query("sort", "name DESC").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.set[0].values[0].value`, "7")).
		Assert(jsonpath.Equal(`$.response.set[1].values[0].value`, "6")).
		Assert(jsonpath.Equal(`$.response.filter.total`, float64(7))).
		Assert(jsonpath.Present(`$.response.filter.pageNavigation`)).
		Assert(jsonpath.Len(`$.response.filter.pageNavigation`, 4)).
		End().
		JSON(&aux)

	h.a.Len(aux.Response.Filter.PageNavigation, 4)
	h.a.NotNil(aux.Response.Filter.PageNavigation[1].Cursor)

	// 2nd page
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Query("incTotal", "false").
		Query("incPageNavigation", "false").
		Query("limit", "2").
		Query("pageCursor", *aux.Response.Filter.PageNavigation[1].Cursor).
		Query("sort", "name DESC").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.set[0].values[0].value`, "5")).
		Assert(jsonpath.Equal(`$.response.set[1].values[0].value`, "4")).
		Assert(jsonpath.NotPresent(`$.response.filter.total`)).
		Assert(jsonpath.NotPresent(`$.response.filter.pageNavigation`)).
		End()
}

func TestRecordListForbiddenRecords(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	helpers.DenyMe(h, types.RecordRbacResource(0, 0, 0), "read")

	h.makeRecord(module)
	h.makeRecord(module)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Query("incTotal", "true").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		// not present because omitted when empty
		Assert(jsonpath.NotPresent(`$.response.filter.total`)).
		End()
}

func TestRecordListForbiddenFields(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, module.RbacResource(), "records.create", "records.search")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(module.NamespaceID, module.ID, module.Fields[0].ID), "record.value.read")

	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "v_name_0"}, &types.RecordValue{Name: "email", Value: "v_email_0"})
	h.makeRecord(module, &types.RecordValue{Name: "name", Value: "v_name_1"}, &types.RecordValue{Name: "email", Value: "v_email_1"})

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Query("incTotal", "true").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.filter.total`, float64(2))).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		// rec 1
		Assert(jsonpath.Len(`$.response.set[0].values`, 1)).
		Assert(jsonpath.Equal(`$.response.set[0].values[0].name`, "email")).
		// rec 2
		Assert(jsonpath.Len(`$.response.set[1].values`, 1)).
		Assert(jsonpath.Equal(`$.response.set[1].values[0].name`, "email")).
		End()
}

func TestRecordCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("record.errors.notAllowedToCreate")).
		End()
}

func TestRecordCreate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRecordCreateForbidden_forbiddenFields(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module",
		&types.ModuleField{Name: "f1", Kind: "String"},
		&types.ModuleField{Name: "f2", Kind: "String"},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.DenyMe(h, module.Fields[1].RbacResource(), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "f1", "value": "f1.v1"}, {"name": "f2", "value": "f2.v1"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertErrorP("1 issue(s) found")).
		End()
}

func TestRecordCreate_forbiddenFields(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module",
		&types.ModuleField{Name: "f1", Kind: "String"},
		&types.ModuleField{Name: "f2", Kind: "String"},
	)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.DenyMe(h, module.Fields[1].RbacResource(), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "f1", "value": "f1.v1"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRecordCreate_xss(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	var (
		ns  = h.makeNamespace("some-namespace")
		mod = h.makeModule(ns, "some-module",
			&types.ModuleField{
				Kind: "String",
				Name: "dummy",
			},
			&types.ModuleField{
				Kind: "String",
				Name: "dummyRichTextBox",
				Options: map[string]interface{}{
					"useRichTextEditor": true,
				},
			},
		)
	)

	helpers.AllowMe(h, ns.RbacResource(), "read")
	helpers.AllowMe(h, mod.RbacResource(), "read")
	helpers.AllowMe(h, mod.RbacResource(), "record.create")
	helpers.AllowMe(h, mod.Fields[0].RbacResource(), "record.value.update")
	helpers.AllowMe(h, mod.Fields[1].RbacResource(), "record.value.update")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")

	t.Run("create with rich text fields", func(t *testing.T) {
		var (
			req = require.New(t)

			payload = struct {
				Response *types.Record
			}{}

			rec = &types.Record{
				Values: types.RecordValueSet{
					&types.RecordValue{Name: "dummyRichTextBox", Value: "<img src=x onerror=alert(11111)>test"},
					&types.RecordValue{Name: "dummy", Value: "simple-text"},
				},
			}
		)

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/", ns.ID, mod.ID)).
			JSON(helpers.JSON(rec)).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Present(`$.response.values[? @.name=="dummyRichTextBox"]`)).
			Assert(jsonpath.Present(`$.response.values[? @.name=="dummy"]`)).
			Assert(jsonpath.Present(`$.response.values[? @.value=="simple-text"]`)).
			Assert(jsonpath.Present(`$.response.values[? @.value=="<img src=\"x\">test"]`)).
			End().
			JSON(&payload)

		req.NotNil(payload.Response)
		req.NotZero(payload.Response.ID)
	})
}

func TestRecordCreateWithErrors(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	fields := types.ModuleFieldSet{
		&types.ModuleField{
			ID:   id.Next(),
			Name: "name",
		},
		&types.ModuleField{
			ID:       id.Next(),
			Name:     "required",
			Required: true,
		},
	}
	module := h.repoMakeRecordModuleWithFields("record testing module", fields...)
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Assert(helpers.AssertRecordValueError(
			&types.RecordValueError{
				Kind:    "empty",
				Message: "record-field.errors.empty",
				Meta:    map[string]interface{}{"field": "required"},
			},
		)).
		End()
}

func TestRecordUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("record.errors.notAllowedToUpdate")).
		End()
}

func TestRecordUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
}

func TestRecordUpdate_missingField(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module",
		&types.ModuleField{Name: "f1", Kind: "String"},
		&types.ModuleField{Name: "f2", Kind: "String"},
	)
	record := h.makeRecord(module,
		&types.RecordValue{Name: "f1", Value: "f1.v1"},
		&types.RecordValue{Name: "f2", Value: "f2.v1"},
	)
	_ = record
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	api := h.apiInit()

	// delete f2
	api.
		Post(fmt.Sprintf("/namespace/%d/module/%d", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"name": "%s", "handle": "%s", "fields": [{ "name": "f1", "kind": "String" }]}`, module.Name, module.Handle)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// update f1
	api.
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "f1", "value": "f1.v1 (edited)"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
	h.a.Len(r.Values, 1)
	h.a.Equal("f1", r.Values[0].Name)
	h.a.Equal("f1.v1 (edited)", r.Values[0].Value)
}

func TestRecordUpdateForbidden_forbiddenFields(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module",
		&types.ModuleField{Name: "f1", Kind: "String"},
		&types.ModuleField{Name: "f2", Kind: "String"},
	)
	record := h.makeRecord(module,
		&types.RecordValue{Name: "f1", Value: "f1.v1"},
		&types.RecordValue{Name: "f2", Value: "f2.v1"},
	)

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")
	helpers.DenyMe(h, module.Fields[1].RbacResource(), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "f1", "value": "f1.v1"}, {"name": "f2", "value": "f2.v1 (edited)"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("1 issue(s) found")).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
	h.a.Equal("f2.v1", r.Values.FilterByName("f2")[0].Value)
}

func TestRecordUpdate_forbiddenFields(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module",
		&types.ModuleField{Name: "f1", Kind: "String"},
		&types.ModuleField{Name: "f2", Kind: "String"},

		// we'll test all kinds of boolean inputs
		&types.ModuleField{Name: "f-b-f-n", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-f-m", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-f-e", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-f-z", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-t-n", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-t-m", Kind: "Bool"},
		&types.ModuleField{Name: "f-b-t-v", Kind: "Bool"},
	)
	record := h.makeRecord(module,
		&types.RecordValue{Name: "f1", Value: "f1.v1"},
		&types.RecordValue{Name: "f2", Value: "f2.v1"},
		&types.RecordValue{Name: "f-b-f-n", Value: "0"}, // no-value
		&types.RecordValue{Name: "f-b-f-e", Value: "0"}, // empty
		&types.RecordValue{Name: "f-b-f-z", Value: "0"}, // zero
		&types.RecordValue{Name: "f-b-t-n", Value: "1"}, // no-value
		&types.RecordValue{Name: "f-b-t-v", Value: "1"}, // value
	)
	helpers.AllowMe(h, types.RecordRbacResource(record.NamespaceID, record.ModuleID, record.ID), "update")
	helpers.AllowMe(h, module.Fields[0].RbacResource(), "record.value.update")
	helpers.DenyMe(h, types.ModuleFieldRbacResource(record.NamespaceID, record.ModuleID, 0), "record.value.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [
			{"name": "f1", "value": "f1.v1"},
			{"name": "f2", "value": "f2.v1"},
			{"name": "f-b-f-n"},
			{"name": "f-b-f-e", "value": ""},
			{"name": "f-b-f-z", "value": "0"},
			{"name": "f-b-t-v", "value": "1"}
		]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
	h.a.Equal("f2.v1", r.Values.FilterByName("f2")[0].Value)
	h.a.Equal("", r.Values.FilterByName("f-b-f-n")[0].Value)
	h.a.Equal("", r.Values.FilterByName("f-b-f-e")[0].Value)
	h.a.Equal("", r.Values.FilterByName("f-b-f-z")[0].Value)
	h.a.Equal("1", r.Values.FilterByName("f-b-t-n")[0].Value)
	h.a.Equal("1", r.Values.FilterByName("f-b-t-v")[0].Value)
}

func TestRecordUpdate_refUnchanged(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	namespace := h.makeNamespace("record testing namespace")

	// mods
	mRef := h.makeRecordModuleWithFieldsOnNs("record testing module", namespace)
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", namespace,
		&types.ModuleField{
			Name: "name",
			Kind: "String",
		},
		&types.ModuleField{
			Name: "ref",
			Kind: "Record",
			Options: types.ModuleFieldOptions{
				"moduleID": strconv.FormatUint(mRef.ID, 10),
			},
		},
	)

	// Records
	rRef := h.makeRecord(mRef)
	record := h.makeRecord(module,
		&types.RecordValue{
			Name:  "name",
			Value: "value; name",
		},
		&types.RecordValue{
			Name:  "ref",
			Value: strconv.FormatUint(rRef.ID, 10),
			Ref:   rRef.ID,
		},
	)

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}, {"name": "ref", "value": "%d"}]}`, rRef.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.Equal(rRef.ID, r.Values.Get("ref", 0).Ref)
	h.a.Equal(strconv.FormatUint(rRef.ID, 10), r.Values.Get("ref", 0).Value)
	h.a.NotNil(r)
}

func TestRecordUpdate_refChanged(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	namespace := h.makeNamespace("record testing namespace")

	// mods
	mRef := h.makeRecordModuleWithFieldsOnNs("record testing module (ref)", namespace)
	module := h.makeRecordModuleWithFieldsOnNs("record testing module", namespace,
		&types.ModuleField{
			Name: "name",
			Kind: "String",
		},
		&types.ModuleField{
			Name: "ref",
			Kind: "Record",
			Options: types.ModuleFieldOptions{
				"moduleID": strconv.FormatUint(mRef.ID, 10),
			},
		},
	)

	// Records
	rRef := h.makeRecord(mRef)
	rRef2 := h.makeRecord(mRef)
	record := h.makeRecord(module,
		&types.RecordValue{
			Name:  "name",
			Value: "value; name",
		},
		&types.RecordValue{
			Name:  "ref",
			Value: strconv.FormatUint(rRef.ID, 10),
			Ref:   rRef.ID,
		},
	)

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	h.apiInit().
		Debug().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}, {"name": "ref", "value": "%d"}]}`, rRef2.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.Equal(rRef2.ID, r.Values.Get("ref", 0).Ref)
	h.a.Equal(strconv.FormatUint(rRef2.ID, 10), r.Values.Get("ref", 0).Value)
	h.a.NotNil(r)
}

func TestRecordUpdate_deleteOld(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module, &types.RecordValue{Name: "name", Value: "test name"}, &types.RecordValue{Name: "email", Value: "test@email.tld"})
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "email", "value": "test@email.tld"}]}`)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.values`, 1)).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
	h.a.Empty(r.Values.FilterByName("name"))
}

func TestRecordDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("record.errors.notAllowedToDelete")).
		End()
}

func TestRecordDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r.DeletedAt)
}

func TestRecordAttachment(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	namespace := h.makeNamespace("record attachment testing namespace")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read", "record.value.update")

	const maxSizeLimit = 1

	module := h.makeModule(
		namespace,
		"module",
		&types.ModuleField{Name: "no_constraints", Kind: "File"},
		&types.ModuleField{
			Name:    "max_size",
			Kind:    "File",
			Options: types.ModuleFieldOptions{"maxSize": maxSizeLimit},
		},
		&types.ModuleField{
			Name: "img_only",
			Kind: "File",
			Options: types.ModuleFieldOptions{
				"maxSize":   maxSizeLimit,
				"mimetypes": "image/gif, image/png, image/jpeg",
			},
		},
		&types.ModuleField{Name: "str", Kind: "String"},
		&types.ModuleField{
			Name: "csv_only",
			Kind: "File",
			Options: types.ModuleFieldOptions{
				"mimetypes": "text/csv",
			},
		},
	)

	xxlBlob := bytes.Repeat([]byte("0"), maxSizeLimit*1_000_000+1)

	testImgFh, err := os.ReadFile("./testdata/test.png")
	h.noError(err)

	testCsvFh, err := os.ReadFile("./testdata/test.csv")
	h.noError(err)

	defer func() {
		// reset settings after we're done
		systemService.CurrentSettings.Compose.Record.Attachments.MaxSize = 0
		systemService.CurrentSettings.Compose.Record.Attachments.Mimetypes = nil
	}()

	systemService.CurrentSettings.Compose.Record.Attachments.MaxSize = maxSizeLimit
	systemService.CurrentSettings.Compose.Record.Attachments.Mimetypes = []string{}

	cc := []struct {
		name  string
		file  []byte
		fname string
		mtype string
		form  map[string]string
		test  func(*http.Response, *http.Request) error
	}{
		{
			"empty file",
			[]byte(""),
			"empty",
			"plain/text",
			map[string]string{"fieldName": "no_constraints"},
			helpers.AssertError("attachment.errors.notAllowedToCreateEmptyAttachment"),
		},
		{
			"no file",
			nil,
			"empty",
			"plain/text",
			map[string]string{"fieldName": "no_constraints"},
			helpers.AssertError("attachment.errors.notAllowedToCreateEmptyAttachment"),
		},
		{
			"no field",
			[]byte("."),
			"dot",
			"plain/text",
			nil,
			helpers.AssertError("attachment.errors.invalidModuleField"),
		},
		{
			"invalid field",
			[]byte("."),
			"dot",
			"plain/text",
			map[string]string{"fieldName": "str"},
			helpers.AssertError("attachment.errors.invalidModuleField"),
		},
		{
			"valid upload, no constraints",
			[]byte("."),
			"dot",
			"plain/text",
			map[string]string{"fieldName": "no_constraints"},
			helpers.AssertNoErrors,
		},
		{
			"global max size - over sized",
			xxlBlob,
			"numbers",
			"plain/text",
			map[string]string{"fieldName": "no_constraints"},
			helpers.AssertError("attachment.errors.tooLarge"),
		},
		{
			"field max size - ok",
			[]byte("12345"),
			"numbers",
			"plain/text",
			map[string]string{"fieldName": "max_size"},
			helpers.AssertNoErrors,
		},
		{
			"field max size - over sized",
			xxlBlob,
			"numbers",
			"plain/text",
			map[string]string{"fieldName": "max_size"},
			helpers.AssertError("attachment.errors.tooLarge"),
		},
		{
			"global mimetype - invalid",
			testImgFh,
			"numbers.gif",
			"image/gif",
			map[string]string{"fieldName": "no_constraints"},
			helpers.AssertError("attachment.errors.failedToProcessImage"),
		},
		{
			"field mimetype - ok",
			testImgFh,
			"image.png",
			"image/gif",
			map[string]string{"fieldName": "img_only"},
			helpers.AssertNoErrors,
		},
		{
			"field mimetype - invalid",
			testImgFh,
			"image.png",
			"image/gif",
			map[string]string{"fieldName": "img_only"},
			helpers.AssertNoErrors,
		},
		{
			"csv file - ok",
			testCsvFh,
			"testCSV",
			"text/csv",
			map[string]string{"fieldName": "csv_only"},
			helpers.AssertNoErrors,
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			h.t = t

			helpers.InitFileUpload(t, h.apiInit(),
				fmt.Sprintf("/namespace/%d/module/%d/record/attachment", module.NamespaceID, module.ID),
				c.form,
				c.file,
				c.fname,
				c.mtype,
			).
				Status(http.StatusOK).
				Assert(c.test).
				End()

		})
	}
}

func TestRecordExport(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// module := h.repoMakeRecordModuleWithFields("record export module")
	// expected := "id,name\n"
	// for i := 0; i < 10; i++ {
	// 	r := h.makeRecord(module, &types.RecordValue{Name: "name", Value: fmt.Sprintf("d%d", i), Place: uint(i)})
	// 	expected += fmt.Sprintf("%d,d%d\n", r.ID, i)
	// }

	// // we'll not use standard asserts (AssertNoErrors) here,
	// // because we're not returning JSON errors.
	// r := h.apiInit().
	// 	Get(fmt.Sprintf("/namespace/%d/module/%d/record/export.csv", module.NamespaceID, module.ID)).
	// 	Query("fields", "name").
	// 	Header("Accept", "application/json").
	// 	Expect(t).
	// 	Status(http.StatusOK).
	// 	End()

	// b, err := ioutil.ReadAll(r.Response.Body)
	// h.noError(err)
	// h.a.Equal(expected, string(b))
}

func (h helper) apiInitRecordImport(api *apitest.APITest, url, f string, file []byte) *apitest.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", f)
	h.noError(err)

	_, err = part.Write(file)
	h.noError(err)
	h.noError(writer.Close())

	return api.
		Post(url).
		Header("Accept", "application/json").
		Body(body.String()).
		ContentType(writer.FormDataContentType()).
		Expect(h.t).
		Status(http.StatusOK)
}

func (h helper) apiRunRecordImport(api *apitest.APITest, url, b string) *apitest.Response {
	return api.
		Patch(url).
		Header("Accept", "application/json").
		JSON(b).
		Expect(h.t).
		Status(http.StatusOK)
}

func TestRecordImportInit(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// module := h.repoMakeRecordModuleWithFields("record import init module")
	// tests := []struct {
	// 	Name    string
	// 	Content string
	// }{
	// 	{
	// 		Name:    "f1.csv",
	// 		Content: "name,email\nv1,v2\n",
	// 	},
	// 	{
	// 		Name:    "f1.json",
	// 		Content: `{"name":"v1","email":"v2"}` + "\n",
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.Name, func(t *testing.T) {
	// 		url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	// 		h.apiInitRecordImport(h.apiInit(), url, test.Name, []byte(test.Content)).
	// 			Assert(jsonpath.Present("$.response.sessionID")).
	// 			Assert(jsonpath.Present(`$.response.fields.name==""`)).
	// 			Assert(jsonpath.Present(`$.response.fields.email==""`)).
	// 			Assert(jsonpath.Present("$.response.progress")).
	// 			Assert(jsonpath.Present("$.response.progress.entryCount==1")).
	// 			End()
	// 	})
	// }
}

func TestRecordImportInit_invalidFileFormat(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// 	h := newHelper(t)
	// 	h.clearRecords()

	// 	module := h.repoMakeRecordModuleWithFields("record import init module")
	// 	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	// 	h.apiInitRecordImport(h.apiInit(), url, "invalid", []byte("nope")).
	// 		Assert(helpers.AssertError("compose.service.RecordImportFormatNotSupported")).
	// 		End()
	// }

	// func TestRecordImportRun(t *testing.T) {
	// 	h := newHelper(t)
	// 	h.clearRecords()
	// 	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")

	// 	module := h.repoMakeRecordModuleWithFields("record import run module")
	// 	tests := []struct {
	// 		Name    string
	// 		Content string
	// 	}{
	// 		{
	// 			Name:    "f1.csv",
	// 			Content: "fname,femail\nv1,v2\n",
	// 		},
	// 	}

	// 	for _, test := range tests {
	// 		t.Run(test.Name, func(t *testing.T) {
	// 			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	// 			rsp := &rImportSession{}
	// 			api := h.apiInit()

	// 			r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
	// 			r.JSON(rsp)

	// 			h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
	// 				Assert(helpers.AssertNoErrors).
	// 				Assert(jsonpath.Present("$.response.progress")).
	// 				Assert(jsonpath.Present(`$.response.fields.fname=="name"`)).
	// 				Assert(jsonpath.Present(`$.response.fields.femail=="email"`)).
	// 				End()
	// 		})
	// 	}
}

func TestRecordImportRun_sessionNotFound(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// module := h.repoMakeRecordModuleWithFields("record import run module")
	// h.apiRunRecordImport(h.apiInit(), fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
	// 	Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
	// 	End()
}

// @todo revert whe we add import RBAC operations
// func TestRecordImportRunForbidden(t *testing.T) {
// 	h := newHelper(t)
// 	h.clearRecords()
// 	helpers.DenyMe(h, types.ModuleRbacResource(0, 0), "record.create")

// 	module := h.repoMakeRecordModuleWithFields("record import run module")
// 	tests := []struct {
// 		Name    string
// 		Content string
// 	}{
// 		{
// 			Name:    "f1.csv",
// 			Content: "fname,femail\nv1,v2\n",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
// 			rsp := &rImportSession{}
// 			api := h.apiInit()

// 			r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
// 			r.JSON(rsp)

// 			h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
// 				Assert(helpers.AssertErrorP("record.errors.notAllowedToCreate for module")).
// 				End()
// 		})
// 	}
// }

// @todo revert whe we add import RBAC operations
// func TestRecordImportRunForbidden_field(t *testing.T) {
// 	h := newHelper(t)
// 	h.clearRecords()
// 	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
// 	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")

// 	module := h.repoMakeRecordModuleWithFields("record import run module")

// 	f := module.Fields.FindByName("name")
// 	helpers.DenyMe(h, f.RbacResource(), "record.value.update")

// 	tests := []struct {
// 		Name    string
// 		Content string
// 	}{
// 		{
// 			Name:    "f1.csv",
// 			Content: "fname,femail\nv1,v2\n",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			h.t = t
// 			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
// 			rsp := &rImportSession{}
// 			api := h.apiInit()

// 			r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
// 			r.JSON(rsp)

// 			h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
// 				Assert(helpers.AssertErrorP("1 issue(s) found")).
// 				End()
// 		})
// 	}
// }

func TestRecordImportRunFieldError_missing(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()
	// helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")

	// module := h.repoMakeRecordModuleWithFieldsRequired("record import run module")

	// tests := []struct {
	// 	Name    string
	// 	Content string
	// }{
	// 	{
	// 		Name:    "f1.csv",
	// 		Content: "fname,femail\n,v2\n",
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.Name, func(t *testing.T) {
	// 		url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	// 		rsp := &rImportSession{}
	// 		api := h.apiInit()

	// 		r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
	// 		r.JSON(rsp)

	// 		h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"femail":"email"},"onError":"skip"}`).
	// 			End()

	// 		api.Get(fmt.Sprintf("%s/%s", url, rsp.Response.SessionID)).
	// 			Expect(h.t).
	// 			Status(http.StatusOK).
	// 			Assert(helpers.AssertNoErrors).
	// 			Assert(jsonpath.Present("$.response.progress.failLog.errors[\"empty field name\"]")).
	// 			End()
	// 	})
	// }
}

func TestRecordImportImportProgress(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// module := h.repoMakeRecordModuleWithFields("record import session module")
	// tests := []struct {
	// 	Name    string
	// 	Content string
	// }{
	// 	{
	// 		Name:    "f1.csv",
	// 		Content: "fname,femail\nv1,v2\n",
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.Name, func(t *testing.T) {
	// 		url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	// 		rsp := &rImportSession{}
	// 		api := h.apiInit()

	// 		r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
	// 		r.JSON(rsp)

	// 		api.Get(fmt.Sprintf("%s/%s", url, rsp.Response.SessionID)).
	// 			Expect(h.t).
	// 			Status(http.StatusOK).
	// 			Assert(helpers.AssertNoErrors).
	// 			Assert(jsonpath.Present("$.response.progress")).
	// 			End()
	// 	})
	// }
}

func TestRecordImportImportProgress_sessionNotFound(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// module := h.repoMakeRecordModuleWithFields("record import module")
	// h.apiInit().
	// 	Get(fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID)).
	// 	Header("Accept", "application/json").
	// 	Expect(h.t).
	// 	Status(http.StatusOK).
	// 	Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
	// 	End()
}

func TestRecordFieldModulePermissionCheck(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	// make a standard module, and prevent (DENY) current user to
	// read from "name" and update "email" fields
	module := h.repoMakeRecordModuleWithFields("record testing module")
	helpers.AllowMe(h, module.RbacResource(), "records.create", "records.search")
	helpers.DenyMe(h, module.Fields.FindByName("name").RbacResource(), "record.value.read")
	helpers.DenyMe(h, module.Fields.FindByName("email").RbacResource(), "record.value.update")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update")

	record := h.makeRecord(
		module,
		&types.RecordValue{Name: "name", Value: "should not be readable"},
		&types.RecordValue{Name: "email", Value: "should not be writable"},
	)

	// Fetching record should work as before but without read-protected fields
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		// should not return name
		Assert(jsonpath.NotPresent(`$.response.values[? @.name=="name"]`)).
		// should return email
		Assert(jsonpath.Present(`$.response.values[? @.name=="email"]`)).
		End()

	// Searching records should work as before but without read-protected fields
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		// should not return name
		Assert(jsonpath.NotPresent(`$.response.set[0].values[? @.name=="name"]`)).
		// should return email
		Assert(jsonpath.Present(`$.response.set[0].values[? @.name=="email"]`)).
		End()

	bb := map[string]func() *apitest.Request{
		"update": func() *apitest.Request {
			return h.apiInit().
				Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID))
		},

		"create": func() *apitest.Request {
			return h.apiInit().
				Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID))
		},
	}

	for name, b := range bb {
		t.Run(name, func(t *testing.T) {
			t.Run("field:email", func(t *testing.T) {
				// Try to change email (not writable!), expect error...
				b().JSON(fmt.Sprintf(`{"values": [{"name": "email", "value": "changed-email"}]}`)).
					Header("Accept", "application/json").
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertErrorP("1 issue(s) found")).
					End()
			})

			t.Run("field:name", func(t *testing.T) {
				// Try to change name, (not readable), expect it to work
				b().JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-name"}]}`)).
					Header("Accept", "application/json").
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End()
			})

			t.Run("field:description", func(t *testing.T) {
				// Try to change description, (no perm. rules), expect it to work
				b().JSON(fmt.Sprintf(`{"values": [{"name": "description", "value": "changed-description"}]}`)).
					Header("Accept", "application/json").
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End()
			})
		})
	}
}

func TestRecordLabels(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "record.create", "records.search")
	helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "update", "read")

	var (
		ns  = h.makeNamespace("some-namespace")
		mod = h.makeModule(ns, "some-module", &types.ModuleField{Kind: "String", Name: "dummy"})
		ID  uint64
	)

	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "record.create", "read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read")
	helpers.AllowMe(h, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update")

	t.Run("create", func(t *testing.T) {
		var (
			req = require.New(t)

			payload = struct {
				Response *types.Record
			}{}

			rec = &types.Record{
				Values: types.RecordValueSet{&types.RecordValue{Name: "dummy", Value: "dummy"}},
				Meta: map[string]any{
					"foo": "bar",
					"bar": "42",
				},
			}
		)

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/", ns.ID, mod.ID)).
			JSON(helpers.JSON(rec)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&payload)

		req.NotNil(payload.Response)
		req.NotZero(payload.Response.ID)

		h.a.Equal(payload.Response.Meta["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Response.Meta["bar"], "42",
			"labels must contain bar with value 42")

		ID = payload.Response.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)

			payload = struct {
				Response *types.Record
			}{}

			rec = &types.Record{
				ID:     ID,
				Values: types.RecordValueSet{&types.RecordValue{Name: "dummy", Value: "dummy"}},
				Meta: map[string]any{
					"foo": "baz",
					"baz": "123",
				},
			}
		)

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", ns.ID, mod.ID, ID)).
			JSON(helpers.JSON(rec)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&payload)
		req.NotZero(payload.Response.ID)

		// disabled for now
		//req.Nil(payload.Response.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Response.Meta["foo"], "baz",
			"meta must contain foo with value baz")
		req.NotContains(payload.Response.Meta, "bar",
			"meta must not contain bar")
		req.Equal(payload.Response.Meta["baz"], "123",
			"meta must contain baz with value 123")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)

			payload = struct {
				Response struct {
					Set types.RecordSet
				}
			}{}

			get = func(meta string, p any) {
				h.apiInit().
					Getf("/namespace/%d/module/%d/record/", ns.ID, mod.ID).
					Header("Accept", "application/json").
					Query("meta", meta).
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End().
					JSON(p)
			}
		)

		get("baz=123", &payload)
		t.Log("is record included in search result?")
		req.NotEmpty(payload.Response.Set)
		req.Len(payload.Response.Set, 1)
		req.NotNil(payload.Response.Set.FindByID(ID))

		t.Log("is meta included")
		req.NotNil(payload.Response.Set.FindByID(ID).Meta)

		get("k2342341241241244=bar", &payload)
		t.Log("no records with this meta constraints should exist")
		req.Empty(payload.Response.Set)

		get("baz", &payload)
		t.Log("one record should be found")
		req.Len(payload.Response.Set, 1)
		req.NotNil(payload.Response.Set.FindByID(ID))
	})
}

func TestRecordReports(t *testing.T) {
	t.Skip("@todo not yet refactored")

	// h := newHelper(t)
	// h.clearRecords()

	// helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "records.search")
	// helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	// helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "record.create")
	// helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")

	// var (
	// 	ns  = h.makeNamespace("some-namespace")
	// 	mod = h.makeModule(ns, "some-module",
	// 		&types.ModuleField{
	// 			Kind:    "Number",
	// 			Name:    "n_float",
	// 			Options: types.ModuleFieldOptions{"precision": 2},
	// 		},
	// 		&types.ModuleField{
	// 			Kind:    "Number",
	// 			Name:    "n_int",
	// 			Options: types.ModuleFieldOptions{"precision": 0},
	// 		},
	// 		&types.ModuleField{
	// 			Kind:    "Number",
	// 			Name:    "n_int_multi",
	// 			Multi:   true,
	// 			Options: types.ModuleFieldOptions{"precision": 0},
	// 		},
	// 	)
	// )

	// h.makeRecord(mod,
	// 	&types.RecordValue{Name: "n_float", Value: "1.1"},
	// 	&types.RecordValue{Name: "n_int", Value: "1"},
	// 	&types.RecordValue{Name: "n_int_multi", Value: "1"},
	// )

	// h.makeRecord(mod,
	// 	&types.RecordValue{Name: "n_float", Value: "2.3"},
	// 	&types.RecordValue{Name: "n_int", Value: "2"},
	// 	&types.RecordValue{Name: "n_int_multi", Value: "1"},
	// 	&types.RecordValue{Name: "n_int_multi", Value: "2", Place: 1},
	// 	&types.RecordValue{Name: "n_int_multi", Value: "3", Place: 2},
	// )

	// t.Run("base metrics", func(t *testing.T) {
	// 	tcc := []struct {
	// 		op         string
	// 		expCount   float64
	// 		expFloat   float64
	// 		expInteger float64
	// 		expMultInt float64
	// 	}{
	// 		{
	// 			op:         "COUNT",
	// 			expCount:   2,
	// 			expFloat:   2,
	// 			expInteger: 2,
	// 			expMultInt: 4, // counting multi values as well
	// 		},
	// 		{
	// 			op:         "SUM",
	// 			expCount:   2,
	// 			expFloat:   3.4,
	// 			expInteger: 3,
	// 			expMultInt: 7, // summing multi values as well
	// 		},
	// 		{
	// 			op:         "MAX",
	// 			expCount:   2,
	// 			expFloat:   2.3,
	// 			expInteger: 2,
	// 			expMultInt: 3, // all values, even the last one
	// 		},
	// 		{
	// 			op:         "MIN",
	// 			expCount:   2,
	// 			expFloat:   1.1,
	// 			expInteger: 1,
	// 			expMultInt: 1,
	// 		},
	// 		{
	// 			op:         "AVG",
	// 			expCount:   2,
	// 			expFloat:   1.7,
	// 			expInteger: 1.5,
	// 			expMultInt: 1.75, // all values!
	// 		},
	// 		// @todo
	// 		// {
	// 		// 	op: "STD",
	// 		// 	expFloat: 0,
	// 		// 	expInteger: 0,
	// 		// },
	// 	}

	// 	for _, tc := range tcc {
	// 		t.Run("basic operations; float; "+tc.op, func(t *testing.T) {
	// 			h.apiInit().
	// 				Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", mod.NamespaceID, mod.ID)).
	// 				Query("metrics", tc.op+"(n_float) as rp").
	// 				Query("dimensions", "DATE_FORMAT(created_at,'Y-01-01')").
	// 				Header("Accept", "application/json").
	// 				Expect(t).
	// 				Status(http.StatusOK).
	// 				Assert(jsonpath.Len(`$.response`, 1)).
	// 				Assert(jsonpath.Equal(`$.response[0].count`, tc.expCount)).
	// 				Assert(jsonpath.Equal(`$.response[0].rp`, tc.expFloat)).
	// 				End()
	// 		})
	// 		t.Run("basic operations; int; "+tc.op, func(t *testing.T) {
	// 			h.apiInit().
	// 				Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", mod.NamespaceID, mod.ID)).
	// 				Query("metrics", tc.op+"(n_int) as rp").
	// 				Query("dimensions", "DATE_FORMAT(created_at,'Y-01-01')").
	// 				Header("Accept", "application/json").
	// 				Expect(t).
	// 				Status(http.StatusOK).
	// 				Assert(jsonpath.Len(`$.response`, 1)).
	// 				Assert(jsonpath.Equal(`$.response[0].count`, tc.expCount)).
	// 				Assert(jsonpath.Equal(`$.response[0].rp`, tc.expInteger)).
	// 				End()
	// 		})
	// 		t.Run("basic operations; int multi-value-field; "+tc.op, func(t *testing.T) {
	// 			h.apiInit().
	// 				Get(fmt.Sprintf("/namespace/%d/module/%d/record/report", mod.NamespaceID, mod.ID)).
	// 				Query("metrics", tc.op+"(n_int_multi) as rp").
	// 				Query("dimensions", "DATE_FORMAT(created_at,'Y-01-01')").
	// 				Header("Accept", "application/json").
	// 				Expect(t).
	// 				Status(http.StatusOK).
	// 				Assert(jsonpath.Len(`$.response`, 1)).
	// 				Assert(jsonpath.Equal(`$.response[0].count`, tc.expCount)).
	// 				Assert(jsonpath.Equal(`$.response[0].rp`, tc.expMultInt)).
	// 				End()
	// 		})
	// 	}
	// })
}
