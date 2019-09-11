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

func (h helper) repoPage() repository.PageRepository {
	return repository.Page(context.Background(), db())
}

func (h helper) repoMakePage(ns *types.Namespace, name string) *types.Page {
	m, err := h.
		repoPage().
		Create(&types.Page{Title: name, NamespaceID: ns.ID})
	h.a.NoError(err)

	return m
}

func TestPageRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.title`, m.Title)).
		Assert(jsonpath.Equal(`$.response.pageID`, fmt.Sprintf("%d", m.ID))).
		End()
}

func TestPageList(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoMakePage(ns, "app")
	h.repoMakePage(ns, "app")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPageCreateForbidden(t *testing.T) {
	h := newHelper(t)

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		FormData("title", "some-page").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoCreatePermissions")).
		End()
}

func TestPageCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "page.create")

	ns := h.repoMakeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		FormData("title", "some-page").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPageUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		FormData("title", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoUpdatePermissions")).
		End()
}

func TestPageUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")
	h.allow(types.PagePermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		FormData("title", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.repoPage().FindByID(ns.ID, m.ID)
	h.a.NoError(err)
	h.a.NotNil(m)
	h.a.Equal("changed-name", m.Title)
}

func TestPageDeleteForbidden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("compose.service.NoDeletePermissions")).
		End()
}

func TestPageDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "delete")

	ns := h.repoMakeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.repoPage().FindByID(ns.ID, m.ID)
	h.a.Error(err, "compose.repository.PageNotFound")
}
