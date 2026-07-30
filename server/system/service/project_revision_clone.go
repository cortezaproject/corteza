package service

import (
	"context"
	"fmt"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// Copying a revision's system-scoped resources into its draft.
//
// CreateRevision clones the compose namespace via the compose envoy path, which
// carries modules, pages, layouts and charts. Nothing carried the resources
// scoped by PROJECT rather than by namespace -- agents, chatbots, automations,
// connections, roles -- so until 2026-07-30 a branch produced a draft holding
// the parent's data model and none of its logic or access model. The draft's
// own deployment plan then reported every one of them as a deletion, with the
// user having changed nothing.
//
// WHY THIS IS HAND-WRITTEN. Envoy is the obvious tool: rewiring references
// across an ID remap is exactly what it does, and it is how the namespace
// clone works. It cannot do it for system resources. Two independent blockers,
// both verified 2026-07-30:
//
//   - Scope resolution is generated for compose only: getScopeNodes outside
//     compose is a stub, so a project-scoped decode silently reads every row
//     of that type in the store rather than the project's.
//   - There is no matchup at all outside compose. The matchup<X> functions
//     that decide whether an encoded node updates an existing row or creates
//     one are generated only for compose (see matchupModules in
//     compose/envoy/store_encode.gen.go, which searches EVERY row unscoped and
//     matches on any identifier). The system component generates none, so an
//     encode here has no defined create-vs-update behaviour to rely on.
//
// CORRECTION (2026-07-30, second audit): an earlier version of this comment
// cited "matchupAgents in system/envoy/store_encode.gen.go" as proof that a
// copy would MOVE its source. That function does not exist -- the mechanism
// described is compose's, and it was never verified against agents. The
// conclusion below stands on the two blockers above; the original evidence for
// it did not.
//
// Both are fixable only by generating real scope support for system resources,
// which rewrites the encode path of every system type. That is a deliberate
// infrastructure project, not part of a branch-copy fix.
//
// ORDER MATTERS. Resources are copied in dependency order so that by the time
// a resource is copied, everything it can reference has already been copied and
// has an entry in the id map. See cloneProjectResources.

// projectCloneIDs maps a parent revision's resource IDs to their copies in the
// draft, per resource kind. References are remapped through it as each kind is
// copied.
//
// Keyed by the resourceref kind constants rather than by Go type so the map can
// carry compose IDs too: refs from system resources into the cloned namespace
// (a trigger constraint naming a module, say) resolve against the same map.
type projectCloneIDs map[string]map[uint64]uint64

func (m projectCloneIDs) set(kind string, old, new uint64) {
	if old == 0 || new == 0 {
		return
	}
	if m[kind] == nil {
		m[kind] = make(map[uint64]uint64)
	}
	m[kind][old] = new
}

// get resolves a reference to its copy. A reference that has no copy is
// returned unchanged and reported as not remapped, which is the correct
// outcome for the many refs that point at resources living OUTSIDE the project
// (knowledge bases, LLM providers, users) -- those are shared, not copied, and
// the draft must keep pointing at the same row.
func (m projectCloneIDs) get(kind string, old uint64) (new uint64, remapped bool) {
	if kk, ok := m[kind]; ok {
		if n, ok := kk[old]; ok {
			return n, true
		}
	}
	return old, false
}

// cloneProjectResources copies every system-scoped resource of the parent
// revision into the draft, remapping references as it goes.
//
// Runs inside CreateRevision's flow AFTER the namespace clone, so the compose
// side already has its new IDs.
func (svc *project) cloneProjectResources(ctx context.Context, s store.Storer, parent, draft *types.Project) (err error) {
	ids := make(projectCloneIDs)

	// Dependency order: a resource may only reference kinds copied before it.
	// Agents reference automations and workflows; chatbots reference agents;
	// roles' RBAC rules reference everything.
	if err = svc.cloneProjectAgents(ctx, s, parent, draft, ids); err != nil {
		return fmt.Errorf("copy agents into revision: %w", err)
	}

	return
}

// cloneProjectAgents copies the parent's agents into the draft.
//
// Handles are carried over unchanged -- deliberately. The publish diff
// identifies resources across revisions by kind + handle (diffSources in
// project_revision.go), so a renamed copy would read as a removal plus an
// addition and the diff would be useless. That is only possible because the
// agent handle index is scoped per project (see agent.cue's
// unique_handle_per_project); it was global until 2026-07-30, which made a
// same-handle copy impossible to store at all.
func (svc *project) cloneProjectAgents(ctx context.Context, s store.Storer, parent, draft *types.Project, ids projectCloneIDs) (err error) {
	aa, _, err := store.SearchAgents(ctx, s, types.AgentFilter{
		ProjectID: parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range aa {
		cp := *src

		cp.ID = nextID()
		cp.ProjectID = draft.ID
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		// Slices are shared with the source by the struct copy above; anything
		// that gets rewritten below must be detached first or the parent
		// revision's agent is mutated in place.
		cp.Meta.SidebarRoles = append([]string(nil), src.Meta.SidebarRoles...)

		if err = store.CreateAgent(ctx, s, &cp); err != nil {
			return err
		}

		ids.set(resourceref.KindAgent, src.ID, cp.ID)
	}

	return
}
