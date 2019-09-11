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

func (h helper) repoNamespace() repository.NamespaceRepository {
	return repository.Namespace(context.Background(), db())
}

func (h helper) repoMakeNamespace(name string) *types.Namespace {
	ns, err := h.
		repoNamespace().
		Create(&types.Namespace{Name: name})
	h.a.NoError(err)

	return ns
}

func TestNamespaceRead(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, ns.Name)).
		Assert(jsonpath.Equal(`$.response.namespaceID`, fmt.Sprintf("%d", ns.ID))).
		End()
}

func TestNamespaceList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeNamespace("app")
	h.repoMakeNamespace("app")

	h.apiInit().
		Get("/namespace/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestNamespaceCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/namespace/").
		FormData("name", "some-namespace").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestNamespaceCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ComposePermissionResource, "namespace.create")

	h.apiInit().
		Post("/namespace/").
		FormData("name", "some-namespace").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestNamespaceUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestNamespaceUpdate(t *testing.T) {
	h := newHelper(t)
	ns := h.repoMakeNamespace("some-namespace")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns, err := h.repoNamespace().FindByID(ns.ID)
	h.a.NoError(err)
	h.a.NotNil(ns)
	h.a.Equal("changed-name", ns.Name)
}

func TestNamespaceDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestNamespaceDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "delete")

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns, err := h.repoNamespace().FindByID(ns.ID)
	h.a.Error(err, "compose.repository.NamespaceNotFound")
}
