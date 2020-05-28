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

func (h helper) repoMakeWeightedPage(ns *types.Namespace, name string, weight int) *types.Page {
	m, err := h.
		repoPage().
		Create(&types.Page{Title: name, NamespaceID: ns.ID, Weight: weight})
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

func TestPageReadByHandle(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	c := h.repoMakePage(ns, "some-page")

	cbh, err := service.DefaultPage.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.a.NoError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
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

func TestPageList_filterForbiden(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")

	h.repoMakePage(ns, "page")
	f := h.repoMakePage(ns, "page_forbiden")

	h.deny(types.PagePermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.title=="page_forbiden"]`)).
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
		Assert(helpers.AssertError("not allowed to create pages")).
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
		Assert(helpers.AssertError("not allowed to update this page")).
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
		Assert(helpers.AssertError("not allowed to delete this page")).
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
	h.a.Error(err, "page does not exist")
}

func TestPageTreeRead(t *testing.T) {
	h := newHelper(t)

	h.allow(types.NamespacePermissionResource.AppendWildcard(), "read")
	h.allow(types.PagePermissionResource.AppendWildcard(), "read")
	ns := h.repoMakeNamespace("some-namespace")
	h.repoMakeWeightedPage(ns, "p1", 1)
	h.repoMakeWeightedPage(ns, "p4", 4)
	h.repoMakeWeightedPage(ns, "p3", 3)
	h.repoMakeWeightedPage(ns, "p2", 2)

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/tree", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response[0].title`, "p1")).
		Assert(jsonpath.Equal(`$.response[1].title`, "p2")).
		Assert(jsonpath.Equal(`$.response[2].title`, "p3")).
		Assert(jsonpath.Equal(`$.response[3].title`, "p4")).
		End()
}
