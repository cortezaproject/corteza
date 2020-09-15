package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearNamespaces() {
	h.noError(store.TruncateComposeNamespaces(context.Background(), service.DefaultStore))
}

func (h helper) makeNamespace(name string) *types.Namespace {
	ns := &types.Namespace{Name: name, Slug: name}
	ns.ID = id.Next()
	ns.CreatedAt = time.Now()
	h.noError(store.CreateComposeNamespace(context.Background(), service.DefaultStore, ns))
	return ns
}

func (h helper) lookupNamespaceByID(ID uint64) *types.Namespace {
	ns, err := store.LookupComposeNamespaceByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return ns
}

func TestNamespaceRead(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, ns.Name)).
		Assert(jsonpath.Equal(`$.response.namespaceID`, fmt.Sprintf("%d", ns.ID))).
		End()
}

func TestNamespaceReadByHandle(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace-" + string(rand.Bytes(20)))

	nsbh, err := service.DefaultNamespace.With(h.secCtx()).FindByHandle(ns.Slug)

	h.noError(err)
	h.a.NotNil(nsbh)
	h.a.Equal(nsbh.ID, ns.ID)
	h.a.Equal(nsbh.Slug, ns.Slug)
}

func TestNamespaceList(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.makeNamespace("ns1")
	h.makeNamespace("ns2")

	h.apiInit().
		Get("/namespace/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestNamespaceList_filterForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.makeNamespace("namespace")
	f := h.makeNamespace("namespace_forbiden")

	h.deny(types.NamespacePermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/namespace/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="namespace_forbiden"]`)).
		End()
}

func TestNamespaceCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.apiInit().
		Post("/namespace/").
		FormData("name", "some-namespace").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create namespaces")).
		End()
}

func TestNamespaceCreate(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

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
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this namespace")).
		End()
}

func TestNamespaceUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns = h.lookupNamespaceByID(ns.ID)
	h.a.NotNil(ns)
	h.a.Equal("changed-name", ns.Name)
}

func TestNamespaceDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this namespace")).
		End()
}

func TestNamespaceDelete(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "delete")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns = h.lookupNamespaceByID(ns.ID)
	h.a.NotNil(ns.DeletedAt)
}
