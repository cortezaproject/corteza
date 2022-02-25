package compose

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
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

	res.Blocks = types.PageBlocks{
		{BlockID: 1},
		{BlockID: 2},
	}

	res.Config.NavItem.Icon = &types.PageConfigIcon{
		Type:  "type",
		Src:   "src",
		Style: map[string]string{"sty": "le"},
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

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.title`, m.Title)).
		Assert(jsonpath.Equal(`$.response.pageID`, fmt.Sprintf("%d", m.ID))).
		Assert(jsonpath.Equal(`$.response.config.navItem.icon.src`, `src`)).
		Assert(jsonpath.Len(`$.response.blocks`, 2)).
		End()
}

func TestPageReadByHandle(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
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

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
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

func TestPageList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	ns := h.makeNamespace("some-namespace")

	h.repoMakePage(ns, "page")
	f := h.repoMakePage(ns, "page_forbidden")

	helpers.DenyMe(h, types.PageRbacResource(f.NamespaceID, f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.title=="page_forbidden"]`)).
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
		Assert(helpers.AssertError("page.errors.notAllowedToCreate")).
		End()
}

func TestPageCreate(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "page.create")

	ns := h.makeNamespace("some-namespace")

	rsp := struct {
		Response *types.Page `json:"response"`
	}{}

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{
			"title": "some-page",
			"config":{"navItem":{"icon":{"src":"my-icon"}}},
			"blocks":[{"blockID": "1"}]
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&rsp)

	res := h.lookupPageByID(rsp.Response.ID)
	h.a.NotNil(res)
	h.a.Equal("some-page", res.Title)
	h.a.Len(res.Blocks, 1)
	h.a.NotNil(res.Config.NavItem.Icon)
	h.a.Equal("my-icon", res.Config.NavItem.Icon.Src)
}

func TestPageUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	ns := h.makeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		FormData("title", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page.errors.notAllowedToUpdate")).
		End()
}

func TestPageUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	ns := h.makeNamespace("some-namespace")
	res := h.repoMakePage(ns, "some-page")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "update")

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

func TestPageUpdateWithBlocksAndConfig(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	ns := h.makeNamespace("some-namespace")
	res := h.repoMakePage(ns, "some-page")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, res.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{
			"title": "changed-name",
			"config":{"navItem":{"icon":{"src":"my-icon"}}},
			"blocks":[{"blockID": "1"}]
		}`)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupPageByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal("changed-name", res.Title)
	h.a.NotNil(res.Config.NavItem.Icon)
	h.a.Equal("my-icon", res.Config.NavItem.Icon.Src)
	h.a.Len(res.Blocks, 1)
}

func TestPageReorder(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.PageRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	ns := h.makeNamespace("some-namespace")
	res := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/reorder", ns.ID, res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestPageDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.repoMakePage(ns, "some-page")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page.errors.notAllowedToDelete")).
		End()
}

func TestPageDelete(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "delete")

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

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
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

func TestPageAttachment(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	ns := h.makeNamespace("page attachment testing namespace")
	page := h.repoMakePage(ns, "some-page")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read", "update")

	xxlBlob := bytes.Repeat([]byte("0"), 1_000_001)

	testImgFh, err := os.ReadFile("./testdata/test.png")
	h.noError(err)

	defer func() {
		// reset settings after we're done
		systemService.CurrentSettings.Compose.Page.Attachments.MaxSize = 0
		systemService.CurrentSettings.Compose.Page.Attachments.Mimetypes = nil
	}()

	// one megabyte limit
	systemService.CurrentSettings.Compose.Page.Attachments.MaxSize = 1
	systemService.CurrentSettings.Compose.Page.Attachments.Mimetypes = []string{
		"application/octet-stream",
	}

	cc := []struct {
		name  string
		file  []byte
		fname string
		mtype string
		form  map[string]string
		test  func(*http.Response, *http.Request) error
	}{
		{
			"empty file",
			[]byte(""),
			"empty",
			"plain/text",
			map[string]string{},
			helpers.AssertError("attachment.errors.notAllowedToCreateEmptyAttachment"),
		},
		{
			"no file",
			nil,
			"empty",
			"plain/text",
			map[string]string{},
			helpers.AssertError("attachment.errors.notAllowedToCreateEmptyAttachment"),
		},
		{
			"valid upload, no constraints",
			[]byte("."),
			"dot",
			"plain/text",
			map[string]string{},
			helpers.AssertNoErrors,
		},
		{
			"global max size - over sized",
			xxlBlob,
			"numbers",
			"plain/text",
			map[string]string{},
			helpers.AssertError("attachment.errors.tooLarge"),
		},
		{
			"global mimetype - invalid",
			testImgFh,
			"numbers.gif",
			"image/gif",
			map[string]string{},
			helpers.AssertError("attachment.errors.notAllowedToUploadThisType"),
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			h.t = t

			helpers.InitFileUpload(t, h.apiInit(),
				fmt.Sprintf("/namespace/%d/page/%d/attachment", page.NamespaceID, page.ID),
				c.form,
				c.file,
				c.fname,
				c.mtype,
			).
				Status(http.StatusOK).
				Assert(c.test).
				End()

		})
	}
}

func TestPageLabels(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "pages.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "page.create")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "delete")

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
