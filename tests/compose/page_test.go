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
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearPages() {
	h.clearNamespaces()
	h.noError(store.TruncateComposePages(context.Background(), service.DefaultStore))
}

func (h helper) repoMakePage(ns *types.Namespace, name string) *types.Page {
	res := &types.Page{
		ID:          id.Next(),
		CreatedAt:   time.Now(),
		Title:       name,
		NamespaceID: ns.ID,
	}

	h.noError(store.CreateComposePage(context.Background(), service.DefaultStore, res))
	return res
}

func (h helper) repoMakeWeightedPage(ns *types.Namespace, name string, weight int) *types.Page {
	res := &types.Page{
		ID:          id.Next(),
		CreatedAt:   time.Now(),
		Title:       name,
		NamespaceID: ns.ID,
		Weight:      weight,
	}

	h.noError(store.CreateComposePage(context.Background(), service.DefaultStore, res))
	return res
}

func (h helper) lookupPageByID(ID uint64) *types.Page {
	res, err := store.LookupComposePageByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestPageRead(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
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
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.repoMakePage(ns, "some-page")

	cbh, err := service.DefaultPage.FindByHandle(h.secCtx(), ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestPageList(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

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
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.repoMakePage(ns, "page")
	f := h.repoMakePage(ns, "page_forbiden")

	h.deny(types.PageRBACResource.AppendID(f.ID), "read")

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
	h.clearPages()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		Header("Accept", "application/json").
		FormData("title", "some-page").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create pages")).
		End()
}

func TestPageCreate(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "page.create")

	ns := h.makeNamespace("some-namespace")

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
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		FormData("title", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this page")).
		End()
}

func TestPageUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	res := h.repoMakePage(ns, "some-page")
	h.allow(types.PageRBACResource.AppendWildcard(), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, res.ID)).
		FormData("title", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupPageByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal("changed-name", res.Title)
}

func TestPageDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this page")).
		End()
}

func TestPageDelete(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "delete")

	ns := h.makeNamespace("some-namespace")
	res := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupPageByID(res.ID)
	h.a.NotNil(res.DeletedAt)
}

func TestPageTreeRead(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
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

func TestPageLabels(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "page.create")
	h.allow(types.PageRBACResource.AppendWildcard(), "read")
	h.allow(types.PageRBACResource.AppendWildcard(), "update")
	h.allow(types.PageRBACResource.AppendWildcard(), "delete")

	var (
		ns = h.makeNamespace("some-namespace")
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Page{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/page/", ns.ID),
			types.Page{Labels: map[string]string{"foo": "bar", "bar": "42"}},
			payload,
		)
		req.NotZero(payload.ID)

		h.a.Equal(payload.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")

		ID = payload.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req     = require.New(t)
			payload = &types.Page{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(),
			t,
			fmt.Sprintf("/namespace/%d/page/%d", ns.ID, ID),
			types.Page{Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.PageSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/page/", ns.ID),
			&set,
			url.Values{"labels": []string{"baz=123"}},
		)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
