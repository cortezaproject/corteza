package compose

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

type (
	rImportSession struct {
		Response struct {
			SessionID string `json:"sessionID"`
		} `json:"response"`
	}
)

func (h helper) repoRecord() repository.RecordRepository {
	return repository.Record(context.Background(), db())
}

func (h helper) repoMakeRecordModuleWithFieldsOnNs(name string, namespace *types.Namespace, ff ...*types.ModuleField) *types.Module {
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.read")

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

	return h.repoMakeModule(namespace, name, ff...)
}

func (h helper) repoMakeRecordModuleWithFields(name string, ff ...*types.ModuleField) *types.Module {
	namespace := h.repoMakeNamespace("record testing namespace")

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.read")

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

	return h.repoMakeModule(namespace, name, ff...)
}

func (h helper) repoMakeRecord(module *types.Module, rvs ...*types.RecordValue) *types.Record {
	record, err := h.
		repoRecord().
		Create(&types.Record{
			ModuleID:    module.ID,
			NamespaceID: module.NamespaceID,
			// @todo createdAt should be auto-populated by record repo!
			//       (same as we do everywhere else)
			CreatedAt: time.Now(),
		})
	h.a.NoError(err)

	err = h.repoRecord().UpdateValues(record.ID, rvs)
	h.a.NoError(err)

	return record
}

func TestRecordRead(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module)

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

	module := h.repoMakeRecordModuleWithFields("record testing module")

	h.repoMakeRecord(module)
	h.repoMakeRecord(module)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRecordCreateForbidden(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "val"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create records")).
		End()
}

func TestRecordCreate(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.create")

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

	fields := types.ModuleFieldSet{
		&types.ModuleField{
			Name: "name",
		},
		&types.ModuleField{
			Name:     "required",
			Required: true,
		},
	}
	module := h.repoMakeRecordModuleWithFields("record testing module", fields...)
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.create")

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

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module)

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this record")).
		End()
}

func TestRecordUpdate(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module)
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		JSON(fmt.Sprintf(`{"values": [{"name": "name", "value": "changed-val"}]}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	r, err := h.repoRecord().FindByID(module.NamespaceID, record.ID)
	h.a.NoError(err)
	h.a.NotNil(r)
	// h.a.Equal(5, r.OwnedBy)
}

func TestRecordDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module)

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this record")).
		End()
}

func TestRecordDelete(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module)

	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.delete")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoRecord().FindByID(module.NamespaceID, record.ID)
	h.a.Error(err, "record does not exist")
}

func TestRecordExport(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record export module")
	for i := 0; i < 10; i++ {
		h.repoMakeRecord(module, &types.RecordValue{Name: "name", Value: fmt.Sprintf("d%d", i)})
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
	h.a.NoError(err)
	h.a.Equal("name\nd0\nd1\nd2\nd3\nd4\nd5\nd6\nd7\nd8\nd9\n", string(b))
}

func (h helper) apiInitRecordImport(api *apitest.APITest, url, f string, file []byte) *apitest.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", f)
	h.a.NoError(err)

	_, err = part.Write(file)
	h.a.NoError(err)
	h.a.NoError(writer.Close())

	return api.
		Post(url).
		Body(body.String()).
		ContentType(writer.FormDataContentType()).
		Expect(h.t).
		Status(http.StatusOK)
}

func (h helper) apiRunRecordImport(api *apitest.APITest, url, b string) *apitest.Response {
	return api.
		Patch(url).
		JSON(b).
		Expect(h.t).
		Status(http.StatusOK)
}

func TestRecordImportInit(t *testing.T) {
	h := newHelper(t)

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

	module := h.repoMakeRecordModuleWithFields("record import init module")
	url := fmt.Sprintf("/namespace/%d/module/%d/record/import", module.NamespaceID, module.ID)
	h.apiInitRecordImport(h.apiInit(), url, "invalid", []byte("nope")).
		Assert(helpers.AssertError("compose.service.RecordImportFormatNotSupported")).
		End()
}

func TestRecordImportRun(t *testing.T) {
	h := newHelper(t)

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

	module := h.repoMakeRecordModuleWithFields("record import run module")
	h.apiRunRecordImport(h.apiInit(), fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID), `{"fields":{"fname":"name","femail":"email"},"onError":"fail"}`).
		Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
		End()
}

func TestRecordImportImportProgress(t *testing.T) {
	h := newHelper(t)

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

	module := h.repoMakeRecordModuleWithFields("record import module")
	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d/record/import/123", module.NamespaceID, module.ID)).
		Expect(h.t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.RecordImportSessionNotFound")).
		End()
}
