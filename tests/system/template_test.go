package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) repoMakeTemplate(ss ...string) *types.Template {
	var res = &types.Template{
		ID:        id.Next(),
		CreatedAt: time.Now(),
	}

	if len(ss) > 0 {
		res.Handle = ss[0]
	} else {
		res.Handle = "n_" + rs()
	}
	if len(ss) > 1 {
		res.Template = ss[1]
	}
	if len(ss) > 2 {
		res.Type = types.DocumentType(ss[2])
	}

	h.a.NoError(store.CreateTemplate(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupTemplateByID(ID uint64) *types.Template {
	res, err := store.LookupTemplateByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestTemplateRead(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()

	u := h.repoMakeTemplate()

	h.apiInit().
		Get(fmt.Sprintf("/template/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.handle`, u.Handle)).
		Assert(jsonpath.Equal(`$.response.templateID`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestTemplateList(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()

	h.repoMakeTemplate(rs())
	h.repoMakeTemplate(rs())

	h.apiInit().
		Get("/template/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestTemplateList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()

	// @todo this can be a problematic test because it leaves
	//       behind templates that are not denied this context
	//       db purge might be needed

	h.repoMakeTemplate("template")
	f := h.repoMakeTemplate()

	h.deny(f.RbacResource(), "read")

	h.apiInit().
		Get("/template/").
		Query("handle", f.Handle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.handle=="%s"]`, f.Handle))).
		End()
}

func TestTemplateCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()

	h.apiInit().
		Post("/template/").
		Header("Accept", "application/json").
		FormData("handle", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create templates")).
		End()
}

func TestTemplateCreate(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.ComponentRbacResource(), "template.create")

	h.apiInit().
		Post("/template/").
		FormData("handle", rs()).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestTemplateUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	u := h.repoMakeTemplate()

	h.apiInit().
		Put(fmt.Sprintf("/template/%d", u.ID)).
		Header("Accept", "application/json").
		FormData("handle", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this template")).
		End()
}

func TestTemplateUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	res := h.repoMakeTemplate()
	h.allow(types.TemplateRbacResource(0), "update")

	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/template/%d", res.ID)).
		Header("Accept", "application/json").
		FormData("handle", newHandle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupTemplateByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newHandle, res.Handle)
}

func TestTemplateDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	u := h.repoMakeTemplate()

	h.apiInit().
		Delete(fmt.Sprintf("/template/%d", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this template")).
		End()
}

func TestTemplateDelete(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.TemplateRbacResource(0), "delete")

	res := h.repoMakeTemplate()

	h.apiInit().
		Delete(fmt.Sprintf("/template/%d", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupTemplateByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestTemplateUndelete(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.TemplateRBACResource.AppendWildcard(), "delete")

	res := h.repoMakeTemplate()

	h.apiInit().
		Post(fmt.Sprintf("/template/%d/undelete", res.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupTemplateByID(res.ID)
	h.a.NotNil(res)
	h.a.Nil(res.DeletedAt)
}

func TestTemplateRenderForbiden(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.deny(types.TemplateRbacResource(0), "render")

	res := h.repoMakeTemplate("rendering", "Hello, {{.interpolate}}", "text/plain")

	h.apiInit().
		Post(fmt.Sprintf("/template/%d/render/testing.txt", res.ID)).
		JSON(`{"variables": {"interpolate": "world!"}}`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to render this template")).
		End()
}

func TestTemplateRenderDriverUndefined(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.TemplateRbacResource(0), "render")

	res := h.repoMakeTemplate("rendering", "Hello, {{.interpolate}}", "text/notexisting")

	h.apiInit().
		Post(fmt.Sprintf("/template/%d/render/testing.txt", res.ID)).
		JSON(`{"variables": {"interpolate": "world!"}}`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("rendering failed: driver not found")).
		End()
}

func TestTemplateRenderPlain(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.TemplateRbacResource(0), "render")

	res := h.repoMakeTemplate("rendering", "Hello, {{.interpolate}}", "text/plain")

	h.apiInit().
		Post(fmt.Sprintf("/template/%d/render/testing.txt", res.ID)).
		JSON(`{"variables": {"interpolate": "world!"}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertBody("Hello, world!")).
		End()
}

func TestTemplateRenderHTML(t *testing.T) {
	h := newHelper(t)
	h.clearTemplates()
	h.allow(types.TemplateRbacResource(0), "render")

	res := h.repoMakeTemplate("rendering", "<h1>Hello, {{.interpolate}}</h1>", "text/html")

	h.apiInit().
		Post(fmt.Sprintf("/template/%d/render/testing.html", res.ID)).
		JSON(`{"variables": {"interpolate": "world!"}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertBody("<h1>Hello, world!</h1>")).
		End()
}
