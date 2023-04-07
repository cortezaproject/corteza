package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearPageLayouts() {
	h.clearNamespaces()
	h.noError(store.TruncateComposePageLayouts(context.Background(), service.DefaultStore))
}

func (h helper) repoMakePageLayout(ns *types.Namespace, pg *types.Page, title string) *types.PageLayout {
	res := &types.PageLayout{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Meta: types.PageLayoutMeta{
			Title: title,
		},
		NamespaceID: ns.ID,
		PageID:      pg.ID,
	}

	h.noError(store.CreateComposePageLayout(context.Background(), service.DefaultStore, res))
	return res
}

func (h helper) lookupPageLayoutByID(ID uint64) *types.PageLayout {
	res, err := store.LookupComposePageLayoutByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestPageLayoutRead(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")
	ly := h.repoMakePageLayout(ns, pg, "some-page-layout")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.meta.title`, ly.Meta.Title)).
		Assert(jsonpath.Equal(`$.response.pageLayoutID`, fmt.Sprintf("%d", ly.ID))).
		End()
}

func TestPageLayoutList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read", "page-layouts.search")
	ns := h.makeNamespace("some-namespace")

	pg := h.repoMakePage(ns, "some-page")

	h.repoMakePageLayout(ns, pg, "page-layout")
	ly := h.repoMakePageLayout(ns, pg, "page-layout_forbidden")

	helpers.DenyMe(h, types.PageLayoutRbacResource(ly.NamespaceID, ly.PageID, ly.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.meta.title=="page-layout_forbidden"]`)).
		End()
}

func TestPageLayoutCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{
			"meta": {"title": "some-page-layout"}
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToCreate")).
		End()
}

func TestPageLayoutCreate(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "page-layout.create")

	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")

	rsp := struct {
		Response *types.PageLayout `json:"response"`
	}{}

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{
			"meta": {"title": "some-page-layout"}
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&rsp)

	res := h.lookupPageLayoutByID(rsp.Response.ID)
	h.a.NotNil(res)
	h.a.Equal("some-page-layout", res.Meta.Title)
}

func TestPageLayoutUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")
	ly := h.repoMakePageLayout(ns, pg, "some-page-layout")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{
			"meta": {"title": "changed-name"}
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToUpdate")).
		End()
}

func TestPageLayoutUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")
	ly := h.repoMakePageLayout(ns, pg, "some-page-layout")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		JSON(fmt.Sprintf(`{
			"meta": {"title": "changed-name"}
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ly = h.lookupPageLayoutByID(ly.ID)
	h.a.NotNil(ly)
	h.a.Equal("changed-name", ly.Meta.Title)
}

func TestPageLayoutDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")
	ly := h.repoMakePageLayout(ns, pg, "some-page-layout")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToDelete")).
		End()
}

func TestPageLayoutDelete(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "read")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "delete")

	ns := h.makeNamespace("some-namespace")
	pg := h.repoMakePage(ns, "some-page")
	ly := h.repoMakePageLayout(ns, pg, "some-page-layout")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ly = h.lookupPageLayoutByID(ly.ID)
	h.a.NotNil(ly.DeletedAt)
}
