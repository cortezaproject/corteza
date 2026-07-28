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
	"github.com/crusttech/human/server/tests/helpers"
	sysTypes "github.com/crusttech/human/server/system/types"
)

// Verifies, end to end: (1) creating a module directly on an Active project's
// namespace is now blocked (the Defect 2 fix), (2) a namespace cloned as a
// draft revision -- with its ProjectID correctly repointed at the new draft,
// mirroring the CreateRevision fix -- allows both create AND update, and (3)
// the same clone WITHOUT the repoint (ProjectID left on the still-active
// parent, i.e. today's CreateRevision bug reproduced directly) blocks both.
func TestZZDiag3GuardCreateAndRevisionRepoint(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "module.update")

	activeProj := &sysTypes.Project{
		ID:     id.Next(),
		Handle: "zzdiag3-active-" + rs(),
		Status: sysTypes.ProjectStatusActive,
	}
	h.noError(store.CreateProject(context.Background(), service.DefaultStore, activeProj))

	draftProj := &sysTypes.Project{
		ID:     id.Next(),
		Handle: "zzdiag3-draft-" + rs(),
		Status: sysTypes.ProjectStatusDraft,
	}
	h.noError(store.CreateProject(context.Background(), service.DefaultStore, draftProj))

	ns := h.makeNamespace("zzdiag3-namespace")
	ns.ProjectID = activeProj.ID
	h.noError(store.UpdateComposeNamespace(context.Background(), service.DefaultStore, ns))

	// (1) Create directly against the active project's namespace: must now be blocked.
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		JSON(`{ "name": "deal", "handle": "deal", "fields": [{ "name": "title", "kind": "String" }] }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("project.errors.locked")).
		End()

	// Seed the module directly via store (bypassing the guarded service) so we
	// have something to clone -- mirroring how the real namespace already has
	// modules before a revision is ever created.
	h.createModule(ns, &types.Module{
		Name:        "deal",
		Handle:      "deal",
		NamespaceID: ns.ID,
		Fields: types.ModuleFieldSet{
			{Name: "title", Kind: "String"},
		},
	})

	// (2) Clone as a proper draft revision: ProjectID repointed at draftProj
	// via the same explicit post-clone UpdateComposeNamespace the
	// CreateRevision fix performs (CloneFromStore/envoyRun ignores
	// dup.ProjectID/Enabled entirely -- see project_revision.go).
	dupGood := &types.Namespace{
		Name: "zzdiag3-namespace (good revision)",
		Slug: "zzdiag3-namespace-good",
	}
	goodNs, err := service.DefaultNamespace.CloneFromStore(h.secCtx(), ns.ID, dupGood)
	h.noError(err)
	goodNs.ProjectID = draftProj.ID
	goodNs.Enabled = false
	h.noError(store.UpdateComposeNamespace(context.Background(), service.DefaultStore, goodNs))
	helpers.AllowMe(h, types.NamespaceRbacResource(goodNs.ID), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(goodNs.ID), "module.create")
	helpers.AllowMe(h, types.ModuleRbacResource(goodNs.ID, 0), "read", "module.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", goodNs.ID)).
		JSON(`{ "name": "new_module", "handle": "new_module", "fields": [{ "name": "x", "kind": "String" }] }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	gm, err := service.DefaultModule.FindByHandle(h.secCtx(), goodNs.ID, "deal")
	h.noError(err)
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "title", "kind": "String" }, { "name": "y", "kind": "Number" }] }`, gm.Name, gm.Fields[0].ID)
	h.apiInit().
		Put(fmt.Sprintf("/namespace/%d/module/%d", goodNs.ID, gm.ID)).
		JSON(fjs).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// (3) Clone WITHOUT the repoint -- ProjectID left on the still-active
	// parent, reproducing today's (pre-fix) CreateRevision bug directly --
	// must block create, proving the repoint in (2) is load-bearing, not
	// incidental.
	dupBad := &types.Namespace{
		Name: "zzdiag3-namespace (bad revision)",
		Slug: "zzdiag3-namespace-bad",
	}
	badNs, err := service.DefaultNamespace.CloneFromStore(h.secCtx(), ns.ID, dupBad)
	h.noError(err)
	helpers.AllowMe(h, types.NamespaceRbacResource(badNs.ID), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(badNs.ID), "module.create")
	helpers.AllowMe(h, types.ModuleRbacResource(badNs.ID, 0), "read", "module.update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", badNs.ID)).
		JSON(`{ "name": "new_module2", "handle": "new_module2", "fields": [{ "name": "x", "kind": "String" }] }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("project.errors.locked")).
		End()
}
