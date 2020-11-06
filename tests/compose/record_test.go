package compose

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearRecords() {
	h.clearNamespaces()
	h.clearModules()
	h.noError(store.TruncateComposeRecords(context.Background(), service.DefaultStore, nil))
}

type (
	rImportSession struct {
		Response struct {
			SessionID string `json:"sessionID"`
		} `json:"response"`
	}
)

func (h helper) makeRecordModuleWithFieldsOnNs(name string, namespace *types.Namespace, ff ...*types.ModuleField) *types.Module {
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.read")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.read")

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

func (h helper) makeRecord(module *types.Module, rvs ...*types.RecordValue) *types.Record {
	rec := &types.Record{
		ID:          id.Next(),
		CreatedAt:   time.Now(),
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
		Values:      rvs,
	}

	h.noError(store.CreateComposeRecord(context.Background(), service.DefaultStore, module, rec))

	return rec
}

func (h helper) lookupRecordByID(module *types.Module, ID uint64) *types.Record {
	res, err := store.LookupComposeRecordByID(context.Background(), service.DefaultStore, module, ID)
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

	h.makeRecord(module)
	h.makeRecord(module)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
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
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create records")).
		End()
}

func TestRecordCreate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.create")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
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
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.create")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Expect(t).
		Assert(helpers.AssertRecordValueError(
			&types.RecordValueError{
				Kind:    "empty",
				Message: "",
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
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this record")).
		End()
}

func TestRecordUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r)
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
		Assert(helpers.AssertError("not allowed to delete this record")).
		End()
}

func TestRecordDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.makeRecord(module)

	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.delete")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r := h.lookupRecordByID(module, record.ID)
	h.a.NotNil(r.DeletedAt)
}

func TestRecordExport(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record export module")
	for i := 0; i < 10; i++ {
		h.makeRecord(module, &types.RecordValue{Name: "name", Value: fmt.Sprintf("d%d", i), Place: uint(i)})
	}

	// we'll not use standard asserts (AssertNoErrors) here,
	// because we're not returning JSON errors.
	r := h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/export.csv", module.NamespaceID, module.ID)).
		Query("fields", "name").
		Expect(t).
		Status(http.StatusOK).
		End()

	b, err := ioutil.ReadAll(r.Response.Body)
	h.noError(err)
	h.a.Equal("name\nd0\nd1\nd2\nd3\nd4\nd5\nd6\nd7\nd8\nd9\n", string(b))
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
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import init module")
	tests := []struct {
		Name    string
		Content string
	}{
		{
			Name:    "f1.csv",
			Content: "name,email\nv1,v2\n",
		},
		{
			Name:    "f1.json",
			Content: `{"name":"v1","email":"v2"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
			h.apiInitRecordImport(h.apiInit(), url, test.Name, []byte(test.Content)).
				Assert(jsonpath.Present("$.response.sessionID")).
				Assert(jsonpath.Present(`$.response.fields.name==""`)).
				Assert(jsonpath.Present(`$.response.fields.email==""`)).
				Assert(jsonpath.Present("$.response.progress")).
				Assert(jsonpath.Present("$.response.progress.entryCount==1")).
				End()
		})
	}
}

func TestRecordImportInit_invalidFileFormat(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import init module")
	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	h.apiInitRecordImport(h.apiInit(), url, "invalid", []byte("nope")).
		Assert(helpers.AssertError("compose.service.RecordImportFormatNotSupported")).
		End()
}

func TestRecordImportRun(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import run module")
	tests := []struct {
		Name    string
		Content string
	}{
		{
			Name:    "f1.csv",
			Content: "fname,femail\nv1,v2\n",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
			rsp := &rImportSession{}
			api := h.apiInit()

			r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
			r.JSON(rsp)

			h.apiRunRecordImport(api, fmt.Sprintf("%s/%s", url, rsp.Response.SessionID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
				Assert(helpers.AssertNoErrors).
				Assert(jsonpath.Present("$.response.progress")).
				Assert(jsonpath.Present(`$.response.fields.fname=="name"`)).
				Assert(jsonpath.Present(`$.response.fields.femail=="email"`)).
				End()
		})
	}
}

func TestRecordImportRun_sessionNotFound(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import run module")
	h.apiRunRecordImport(h.apiInit(), fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
		Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
		End()
}

func TestRecordImportImportProgress(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import session module")
	tests := []struct {
		Name    string
		Content string
	}{
		{
			Name:    "f1.csv",
			Content: "fname,femail\nv1,v2\n",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
			rsp := &rImportSession{}
			api := h.apiInit()

			r := h.apiInitRecordImport(api, url, test.Name, []byte(test.Content)).End()
			r.JSON(rsp)

			api.Get(fmt.Sprintf("%s/%s", url, rsp.Response.SessionID)).
				Expect(h.t).
				Status(http.StatusOK).
				Assert(helpers.AssertNoErrors).
				Assert(jsonpath.Present("$.response.progress")).
				End()
		})
	}
}

func TestRecordImportImportProgress_sessionNotFound(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("record import module")
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID)).
		Header("Accept", "application/json").
		Expect(h.t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
		End()
}

func TestRecordFieldModulePermissionCheck(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	// make a standard module, and prevent (DENY) current user to
	// read from "name" and update "email" fields
	module := h.repoMakeRecordModuleWithFields("record testing module")
	h.deny(module.Fields.FindByName("name").RBACResource(), "record.value.read")
	h.deny(module.Fields.FindByName("email").RBACResource(), "record.value.update")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.create")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.update")

	record := h.makeRecord(
		module,
		&types.RecordValue{Name: "name", Value: "should not be readable"},
		&types.RecordValue{Name: "email", Value: "should not be writable"},
	)

	// Fetching record should work as before but without read-protected fields
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
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
					Assert(helpers.AssertError("1 issue(s) found")).
					End()
			})

			t.Run("field:name", func(t *testing.T) {
				// Try to change name, (not readable), expect it to work
				b().JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-name"}]}`)).
					Expect(t).
					Status(http.StatusOK).
					Assert(helpers.AssertNoErrors).
					End()
			})

			t.Run("field:description", func(t *testing.T) {
				// Try to change description, (no perm. rules), expect it to work
				b().JSON(fmt.Sprintf(`{"values": [{"name": "description", "value": "changed-description"}]}`)).
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.create")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.update")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "record.read")

	var (
		ns  = h.makeNamespace("some-namespace")
		mod = h.makeModule(ns, "some-module", &types.ModuleField{Kind: "String", Name: "dummy"})
		ID  uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req = require.New(t)

			payload = struct {
				Response *types.Record
			}{}

			rec = &types.Record{
				Values: types.RecordValueSet{&types.RecordValue{Name: "dummy", Value: "dummy"}},
				Labels: map[string]string{
					"foo": "bar",
					"bar": "42",
				},
			}
		)

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/", ns.ID, mod.ID)).
			JSON(helpers.JSON(rec)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&payload)

		req.NotNil(payload.Response)
		req.NotZero(payload.Response.ID)

		h.a.Equal(payload.Response.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Response.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Response.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.Response.LabelResourceKind(), payload.Response.ID),
			"response must match stored labels")

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
				Labels: map[string]string{
					"foo": "baz",
					"baz": "123",
				},
			}
		)

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", ns.ID, mod.ID, ID)).
			JSON(helpers.JSON(rec)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End().
			JSON(&payload)
		req.NotZero(payload.Response.ID)

		// disabled for now
		//req.Nil(payload.Response.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Response.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Response.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Response.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Response.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.Response.LabelResourceKind(), payload.Response.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.RecordSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/module/%d/record/", ns.ID, mod.ID),
			&set, url.Values{"labels": []string{"baz=123"}},
		)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
