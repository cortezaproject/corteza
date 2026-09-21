package compose

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// Resources must only be reachable through the namespace they belong to

func TestModuleReadForeignNamespace(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("some-namespace")
	foreign := h.makeNamespace("foreign-namespace")
	m := h.makeModule(foreign, "some-module")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.notFound")).
		End()
}

func TestPageReadForeignNamespace(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	ns := h.makeNamespace("some-namespace")
	foreign := h.makeNamespace("foreign-namespace")
	p := h.repoMakePage(foreign, "some-page")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d", ns.ID, p.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page.errors.notFound")).
		End()
}

func TestPageLayoutReadForeignNamespace(t *testing.T) {
	h := newHelper(t)
	h.clearPageLayouts()

	ns := h.makeNamespace("some-namespace")
	foreign := h.makeNamespace("foreign-namespace")
	pg := h.repoMakePage(foreign, "some-page")
	ly := h.repoMakePageLayout(foreign, pg, "some-layout")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageLayoutRbacResource(0, 0, 0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/page/%d/layout/%d", ns.ID, pg.ID, ly.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page-layout.errors.notFound")).
		End()
}
