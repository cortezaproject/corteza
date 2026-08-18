package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// projectLockFixture puts a namespace under a project of the given status and
// grants everything the create endpoints below ask for, so the only thing left
// that can refuse them is the project-state lock.
func (h helper) projectLockFixture(status sysTypes.ProjectStatus, label string) (*types.Namespace, *types.Page) {
	proj := &sysTypes.Project{
		ID:     id.Next(),
		Handle: "lock-" + label + "-" + rs(),
		Status: status,
	}
	h.noError(store.CreateProject(context.Background(), service.DefaultStore, proj))

	ns := h.makeNamespace("lock-namespace-" + rs())
	ns.ProjectID = proj.ID
	h.noError(store.UpdateComposeNamespace(context.Background(), service.DefaultStore, ns))

	pg := h.repoMakePage(ns, "lock-page")

	helpers.AllowMe(h, types.NamespaceRbacResource(ns.ID), "read", "pages.search", "charts.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(ns.ID), "page.create", "chart.create")
	helpers.AllowMe(h, types.PageRbacResource(ns.ID, 0), "read", "page-layout.create")

	return ns, pg
}

// Page, chart and page-layout creation is refused once the owning project
// leaves draft, the way module creation already is. Update and Delete get the
// lock from the generated wrapper's svc.guard call; Create has no such hook and
// carries the check in its own onCreate.
func TestProjectLockRefusesCreateOnActiveProject(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns, pg := h.projectLockFixture(sysTypes.ProjectStatusActive, "active")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		JSON(`{ "title": "blocked", "handle": "blocked_page", "visible": true }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("project.errors.locked")).
		End()

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		JSON(`{ "name": "blocked", "handle": "blocked_chart" }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("project.errors.locked")).
		End()

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		JSON(`{ "handle": "blocked_layout", "meta": { "title": "blocked" } }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("project.errors.locked")).
		End()
}

// The same three creates against a DRAFT project's namespace all succeed, so
// the refusals above are the lock talking and not a missing permission.
func TestProjectLockAllowsCreateOnDraftProject(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns, pg := h.projectLockFixture(sysTypes.ProjectStatusDraft, "draft")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/", ns.ID)).
		JSON(`{ "title": "allowed", "handle": "allowed_page", "visible": true }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/chart/", ns.ID)).
		JSON(`{ "name": "allowed", "handle": "allowed_chart" }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/page/%d/layout/", ns.ID, pg.ID)).
		JSON(`{ "handle": "allowed_layout", "meta": { "title": "allowed" } }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
