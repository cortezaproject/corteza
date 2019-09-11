package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) repoRecord() repository.RecordRepository {
	return repository.Record(context.Background(), db())
}

func (h helper) repoMakeRecordModuleWithFields(name string) *types.Module {
	namespace := h.repoMakeNamespace("record testing namespace")

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.read")

	m, err := h.
		repoModule().
		Create(&types.Module{
			Name:        name,
			NamespaceID: namespace.ID,
			Fields: types.ModuleFieldSet{
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
			},
		})

	h.a.NoError(err)

	return m
}

func (h helper) repoMakeRecord(module *types.Module, name string) *types.Record {
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

	return record
}

func TestRecordRead(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module, "some-record")

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

	h.repoMakeRecord(module, "app")
	h.repoMakeRecord(module, "app")

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
		FormData("name", "some-record").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestRecordCreate(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.create")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/", module.NamespaceID, module.ID)).
		FormData("name", "some-record").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRecordUpdateForbidden(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module, "some-record")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestRecordUpdate(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module, "some-record")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		// FormData("ownerID", "5").
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
	record := h.repoMakeRecord(module, "some-record")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestRecordDelete(t *testing.T) {
	h := newHelper(t)

	module := h.repoMakeRecordModuleWithFields("record testing module")
	record := h.repoMakeRecord(module, "some-record")

	h.allow(types.ModulePermissionResource.AppendWildcard(), "record.delete")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d/record/%d", module.NamespaceID, module.ID, record.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoRecord().FindByID(module.NamespaceID, record.ID)
	h.a.Error(err, "compose.repository.RecordNotFound")
}
