package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearModules() {
	h.clearNamespaces()
	h.noError(store.TruncateComposeModules(context.Background(), service.DefaultStore))
	h.noError(store.TruncateComposeModuleFields(context.Background(), service.DefaultStore))
}

func (h helper) makeModule(ns *types.Namespace, name string, ff ...*types.ModuleField) *types.Module {
	return h.createModule(&types.Module{
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      ff,
		CreatedAt:   time.Now(),
	})
}

func (h helper) createModule(res *types.Module) *types.Module {
	res.ID = id.Next()
	res.CreatedAt = time.Now()
	h.noError(store.CreateComposeModule(context.Background(), service.DefaultStore, res))

	_ = res.Fields.Walk(func(f *types.ModuleField) error {
		f.ID = id.Next()
		f.ModuleID = res.ID
		f.CreatedAt = time.Now()
		return nil
	})

	h.noError(store.CreateComposeModuleField(context.Background(), service.DefaultStore, res.Fields...))

	return res
}

func (h helper) lookupModuleByID(ID uint64) *types.Module {
	res, err := store.LookupComposeModuleByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)

	res.Fields, _, err = store.SearchComposeModuleFields(context.Background(), service.DefaultStore, types.ModuleFieldFilter{ModuleID: []uint64{ID}})
	h.noError(err)

	return res
}

func TestModuleRead(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

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
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.makeModule(ns, "some-module")

	cbh, err := service.DefaultModule.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestModuleList(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeModule(ns, "app")
	h.makeModule(ns, "app")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleListQuery(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.createModule(&types.Module{
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
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeModule(ns, "module")
	f := h.makeModule(ns, "module_forbiden")

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
	h.clearModules()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create modules")).
		End()
}

func TestModuleCreate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "module.create")

	ns := h.makeNamespace("some-namespace")

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
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this module")).
		End()
}

func TestModuleUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.updatedAt")).
		End()

	m, err := store.LookupComposeModuleByID(context.Background(), service.DefaultStore, m.ID)
	h.noError(err)
	h.a.NotNil(m)
	h.a.Equal("changed-name", m.Name)
}

func TestModuleFieldsUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
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

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 2)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "existing_edited")
	h.a.Equal(m.Fields[0].Kind, "Number")
	h.a.Nil(m.Fields[1].UpdatedAt)
	h.a.Equal(m.Fields[1].Name, "new")
	h.a.Equal(m.Fields[1].Kind, "DateTime")
}

func TestModuleFieldsPreventUpdate_ifRecordExists(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
	h.makeRecord(m, &types.RecordValue{Name: "existing", Value: "value"})
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

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 2)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "existing")
	h.a.Equal(m.Fields[0].Kind, "String")
	h.a.Nil(m.Fields[1].UpdatedAt)
	h.a.Equal(m.Fields[1].Name, "new")
	h.a.Equal(m.Fields[1].Kind, "DateTime")
}

func TestModuleDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this module")).
		End()
}

func TestModuleDelete(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "read")
	h.allow(types.ModulePermissionResource.AppendWildcard(), "delete")

	ns := h.makeNamespace("some-namespace")
	res := h.makeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupModuleByID(res.ID)
	h.a.NotNil(res.DeletedAt)
}
