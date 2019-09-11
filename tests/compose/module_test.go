package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) repoModule() repository.ModuleRepository {
	return repository.Module(context.Background(), db())
}

func (h helper) repoMakeModule(ns *types.Namespace, name string) *types.Module {
	m, err := h.
		repoModule().
		Create(&types.Module{Name: name, NamespaceID: ns.ID})
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
		End()

	m, err := h.repoModule().FindByID(ns.ID, m.ID)
	h.a.NoError(err)
	h.a.NotNil(m)
	h.a.Equal("changed-name", m.Name)
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
