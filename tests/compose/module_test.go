package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) repoModule() repository.ModuleRepository {
	return repository.Module(context.Background(), db())
}

func (h helper) repoMakeModule(ns *types.Namespace, name string, ff ...*types.ModuleField) *types.Module {
	return h.repoSaveModule(&types.Module{Name: name, NamespaceID: ns.ID, Fields: ff})
}

func (h helper) repoSaveModule(mod *types.Module) *types.Module {
	m, err := h.
		repoModule().
		Create(mod)
	h.a.NoError(err)

	err = h.repoModule().UpdateFields(m.ID, m.Fields, false)
	h.a.NoError(err)

	return m
}

func TestModuleRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, m.Name)).
		Assert(jsonpath.Equal(`$.response.moduleID`, fmt.Sprintf("%d", m.ID))).
		End()
}

func TestModuleReadByHandle(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	c := h.repoMakeModule(ns, "some-module")

	cbh, err := service.DefaultModule.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.a.NoError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestModuleList(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoMakeModule(ns, "app")
	h.repoMakeModule(ns, "app")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleListQuery(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoSaveModule(&types.Module{
		Name:        "name",
		Handle:      "handle",
		NamespaceID: ns.ID,
	})

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Query("query", "handle").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.set != null")).
		End()
}

func TestModuleList_filterForbiden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoMakeModule(ns, "module")
	f := h.repoMakeModule(ns, "module_forbiden")

	h.deny(types.ModulePermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="module_forbiden"]`)).
		End()
}

func TestModuleCreateForbidden(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestModuleCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "module.create")

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestModuleUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.updatedAt")).
		End()

	m, err := h.repoModule().FindByID(ns.ID, m.ID)
	h.a.NoError(err)
	h.a.NotNil(m)
	h.a.Equal("changed-name", m.Name)
}

func TestModuleFieldsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module", &types.ModuleField{Kind: "String", Name: "existing"})
	h.allow(types.ModulePermissionResource.AppendWildcard(), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "existing_edited", "kind": "Number" }, { "name": "new", "kind": "DateTime" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ff, err := h.repoModule().FindFields(m.ID)
	h.a.NoError(err)
	h.a.NotNil(ff)
	h.a.Len(ff, 2)

	h.a.NotNil(ff[0].UpdatedAt)
	h.a.Equal(ff[0].Name, "existing_edited")
	h.a.Equal(ff[0].Kind, "Number")
	h.a.Nil(ff[1].UpdatedAt)
	h.a.Equal(ff[1].Name, "new")
	h.a.Equal(ff[1].Kind, "DateTime")
}

func TestModuleFieldsPreventUpdate_ifRecordExists(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module", &types.ModuleField{Kind: "String", Name: "existing"})
	h.repoMakeRecord(m, &types.RecordValue{Name: "existing", Value: "value"})
	h.allow(types.ModulePermissionResource.AppendWildcard(), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "existing_edited", "kind": "Number" }, { "name": "new", "kind": "DateTime" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ff, err := h.repoModule().FindFields(m.ID)
	h.a.NoError(err)
	h.a.NotNil(ff)
	h.a.Len(ff, 2)

	h.a.Nil(ff[0].UpdatedAt)
	h.a.Equal(ff[0].Name, "existing")
	h.a.Equal(ff[0].Kind, "String")
	h.a.Nil(ff[1].UpdatedAt)
	h.a.Equal(ff[1].Name, "new")
	h.a.Equal(ff[1].Kind, "DateTime")
}

func TestModuleDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestModuleDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "delete")

	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.repoModule().FindByID(ns.ID, m.ID)
	h.a.Error(err, "compose.repository.ModuleNotFound")
}
