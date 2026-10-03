package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// page that allows personal layouts and a layout owned by the given user
func (h helper) repoMakePersonalPageLayout(ns *types.Namespace, ownedBy uint64) (*types.Page, *types.PageLayout) {
	pg := h.repoMakePage(ns, "some-page")
	pg.Meta.AllowPersonalLayouts = true
	h.noError(store.UpdateComposePage(context.Background(), service.DefaultStore, pg))

	ly := h.repoMakePageLayout(ns, pg, "personal")
	ly.OwnedBy = ownedBy
	h.noError(store.UpdateComposePageLayout(context.Background(), service.DefaultStore, ly))

	return pg, ly
}

// Personal layouts can be managed without permissions, but only by their owner

func TestPageLayoutPersonalCreateForOtherUser(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg, _ := h.repoMakePersonalPageLayout(ns, h.cUser.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"meta": {"title": "planted"}, "ownedBy": "%d"}`, id.Next())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToCreate")).
		End()
}

func TestPageLayoutPersonalCreateForMyself(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg, _ := h.repoMakePersonalPageLayout(ns, h.cUser.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"meta": {"title": "mine"}, "ownedBy": "%d"}`, h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPageLayoutPersonalUpdateForeign(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg, ly := h.repoMakePersonalPageLayout(ns, id.Next())
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		JSON(`{"meta": {"title": "defaced"}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToUpdate")).
		End()

	h.a.Equal("personal", h.lookupPageLayoutByID(ly.ID).Meta.Title)
}

func TestPageLayoutPersonalDeleteForeign(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg, ly := h.repoMakePersonalPageLayout(ns, id.Next())
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notAllowedToDelete")).
		End()

	h.a.Nil(h.lookupPageLayoutByID(ly.ID).DeletedAt)
}

func TestPageLayoutPersonalDeleteOwn(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	pg, ly := h.repoMakePersonalPageLayout(ns, h.cUser.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupPageLayoutByID(ly.ID).DeletedAt)
}

// A personal layout cannot be moved to another page or namespace through an
// update; it stays where it was created
func TestPageLayoutPersonalUpdateKeepsPage(t *testing.T) {
	h := newHelper(t)
	ns := h.makeNamespace("personal layout ns")
	pg, ly := h.repoMakePersonalPageLayout(ns, h.cUser.ID)
	other := h.repoMakePage(ns, "other-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"meta": {"title": "moved?"}, "pageID": "%d", "namespaceID": "%d"}`, other.ID, ns.ID+1)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	stored, err := store.LookupComposePageLayoutByID(context.Background(), service.DefaultStore, ly.ID)
	h.noError(err)
	h.a.Equal(pg.ID, stored.PageID)
	h.a.Equal(ns.ID, stored.NamespaceID)
	h.a.Equal("moved?", stored.Meta.Title)
}

func TestPageLayoutPersonalUpdateOwn(t *testing.T) {
	h := newHelper(t)
	ns := h.makeNamespace("personal layout ns")
	pg, ly := h.repoMakePersonalPageLayout(ns, h.cUser.ID)

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		JSON(`{"meta": {"title": "mine"}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	stored, err := store.LookupComposePageLayoutByID(context.Background(), service.DefaultStore, ly.ID)
	h.noError(err)
	h.a.Equal("mine", stored.Meta.Title)
}
