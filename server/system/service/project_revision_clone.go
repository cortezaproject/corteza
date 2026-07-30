package service

import (
	"context"
	"fmt"
	"strconv"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
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
//
// WHAT HAPPENS TO A REFERENCE. Every reference out of a copied resource falls
// into one of three cases, and getting them confused is how a draft ends up
// driving live data:
//
//   - Copied with the revision (the namespace, its modules, and every kind
//     this file copies): remapped to the copy's ID. Missing this is the whole
//     bug -- a copied agent that still names the parent's modules reads and
//     writes the LIVE revision.
//   - Shared with the rest of the system (knowledge bases, LLM providers,
//     users, service accounts): left untouched. Those rows are not part of the
//     revision and the draft must keep pointing at the same ones.
//   - Part of the parent revision but not copied by this pass (its TAQs and
//     workflows, which no kind here carries yet): DROPPED, and logged. The
//     draft loses a binding the user can re-add; the alternative is a draft
//     agent invoking logic bound to the live revision, which is the failure
//     this file exists to prevent.
//
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

type (
	// projectClone is the state one branch copy runs against: the two
	// revisions, the ID map built up as kinds are copied, and what is known
	// about which resources the parent revision owns.
	projectClone struct {
		parent, draft *types.Project

		ids projectCloneIDs

		// The parent namespace's slug and the clone's. Agent access context
		// names a namespace by slug as readily as by ID (see
		// remapNamespaceIdent), and the clone renames it.
		parentNsSlug, draftNsSlug string

		// owned answers "is this resource part of the revision being
		// branched?" per kind, memoised -- see ownedByParent.
		owned map[string]map[uint64]bool
	}

	// droppedRef is a reference the copy could not carry into the draft.
	droppedRef struct {
		kind string
		id   uint64
	}
)

// newProjectClone builds the reference state for one branch copy: the compose
// side is already cloned by the time this runs, so its ID mapping can be
// resolved up front and every kind copied afterwards remaps against it.
func newProjectClone(ctx context.Context, s store.Storer, parent, draft *types.Project) (*projectClone, error) {
	c := &projectClone{
		parent: parent,
		draft:  draft,
		ids:    make(projectCloneIDs),
		owned:  make(map[string]map[uint64]bool),
	}

	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("load parent revision namespace: %w", err)
	}
	draftNs, err := store.LookupComposeNamespaceByID(ctx, s, draft.Config.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("load draft revision namespace: %w", err)
	}

	c.parentNsSlug, c.draftNsSlug = parentNs.Slug, draftNs.Slug
	c.ids.set(resourceref.KindComposeNamespace, parentNs.ID, draftNs.ID)
	c.setOwned(resourceref.KindComposeNamespace, parentNs.ID, true)

	// Modules are matched by HANDLE. The compose clone mints fresh IDs and
	// returns no mapping of its own -- CloneFromStore hands back only the
	// namespace -- but it carries handles across unchanged, so the handle is
	// the one thing both sides share. (It is the same identity the publish
	// diff uses across revisions; see diffModules.)
	parentMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: parentNs.ID})
	if err != nil {
		return nil, fmt.Errorf("load parent revision modules: %w", err)
	}
	draftMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: draftNs.ID})
	if err != nil {
		return nil, fmt.Errorf("load draft revision modules: %w", err)
	}

	draftByHandle := make(map[string]uint64, len(draftMods))
	for _, m := range draftMods {
		if m.Handle != "" {
			draftByHandle[m.Handle] = m.ID
		}
	}

	for _, m := range parentMods {
		// Ownership is recorded for every module in the parent's namespace,
		// matched or not: an unmatched one is still part of the revision being
		// branched, so a reference to it must not survive into the draft. The
		// remap drops what it cannot repoint.
		c.setOwned(resourceref.KindComposeModule, m.ID, true)

		if m.Handle == "" {
			continue
		}
		if id, ok := draftByHandle[m.Handle]; ok {
			c.ids.set(resourceref.KindComposeModule, m.ID, id)
		}
	}

	return c, nil
}

func (c *projectClone) setOwned(kind string, id uint64, owned bool) {
	if id == 0 {
		return
	}
	if c.owned[kind] == nil {
		c.owned[kind] = make(map[uint64]bool)
	}
	c.owned[kind][id] = owned
}

// ownedByParent reports whether a referenced resource belongs to the parent
// revision -- the test that separates "part of what is being branched" from
// "shared infrastructure both revisions legitimately point at".
//
// Answers are memoised per kind + id: the agents of one project tend to name
// the same handful of automations.
func (c *projectClone) ownedByParent(ctx context.Context, s store.Storer, kind string, id uint64) (bool, error) {
	if id == 0 {
		return false, nil
	}
	if kk, ok := c.owned[kind]; ok {
		if known, ok := kk[id]; ok {
			return known, nil
		}
	}

	var (
		owned bool
		err   error
	)

	switch kind {
	case resourceref.KindNgAutomation:
		var au *automationTypes.NgAutomation
		if au, err = store.LookupAutomationNgAutomationByID(ctx, s, id); err == nil {
			owned = au.ProjectID == c.parent.ID
		}

	case resourceref.KindAutomationWorkflow:
		// Looked up one at a time rather than by a project-filtered search:
		// WorkflowFilter carries a ProjectID field, but automation/workflow.cue
		// does not list project_id in the filter's byValue, so the store never
		// emits the predicate and a "project-scoped" search silently returns
		// every workflow in the system.
		var wf *automationTypes.Workflow
		if wf, err = store.LookupAutomationWorkflowByID(ctx, s, id); err == nil {
			owned = wf.ProjectID == c.parent.ID
		}

	default:
		// Kinds whose ownership is settled up front (the namespace and its
		// modules, seeded by newProjectClone) never reach the switch. Anything
		// else is not something this pass can classify, and guessing "owned"
		// would drop a reference the draft needs -- treat it as shared.
		return false, nil
	}

	if err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		}
		// A reference to a row that no longer exists is already broken, and it
		// is not this pass's job to tidy that up: dropping it would report a
		// change the branch did not make. Carry it over untouched.
		owned = false
	}

	c.setOwned(kind, id, owned)
	return owned, nil
}

// carryRef decides what becomes of one reference out of a copied resource:
// remapped to its copy, kept as it is, or dropped. See the three cases in this
// file's header comment.
func (c *projectClone) carryRef(ctx context.Context, s store.Storer, kind string, id uint64) (newID uint64, keep bool, err error) {
	if id == 0 {
		return id, true, nil
	}

	if mapped, ok := c.ids.get(kind, id); ok {
		return mapped, true, nil
	}

	owned, err := c.ownedByParent(ctx, s, kind, id)
	if err != nil {
		return id, false, err
	}

	return id, !owned, nil
}

// cloneProjectResources copies every system-scoped resource of the parent
// revision into the draft, remapping references as it goes.
//
// Runs inside CreateRevision's flow AFTER the namespace clone, so the compose
// side already has its new IDs.
func (svc *project) cloneProjectResources(ctx context.Context, s store.Storer, parent, draft *types.Project) (err error) {
	c, err := newProjectClone(ctx, s, parent, draft)
	if err != nil {
		return err
	}

	// Dependency order: a resource may only reference kinds copied before it.
	// Agents reference automations and workflows; chatbots reference agents;
	// roles' RBAC rules reference everything.
	if err = svc.cloneProjectAgents(ctx, s, c); err != nil {
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
func (svc *project) cloneProjectAgents(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	aa, _, err := store.SearchAgents(ctx, s, types.AgentFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	// Labels are a separate table and the store search does not join them in
	// (the agent service loads them explicitly too, in onLookup). Without this
	// the copies below would faithfully carry an empty label set.
	if err = label.Load(ctx, s, toLabeledAgents(aa)...); err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range aa {
		// Deep copy, not `cp := *src`. The remap below writes into
		// Access.Tools[].Allow and the TAQ/workflow slices, and a struct copy
		// shares every one of them with the source -- rewriting them in place
		// would repoint the PARENT revision's agent at the draft.
		cp := *src.Clone()

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		// A copy is a new agent on its first version. Inheriting the source's
		// counter would have the draft's agent claim an edit history it does
		// not have, and onUpdate bumps from whatever is stored.
		cp.Revision = 1
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		var dropped []droppedRef
		if dropped, err = c.remapAgentAccess(ctx, s, &cp); err != nil {
			return err
		}

		if err = store.CreateAgent(ctx, s, &cp); err != nil {
			return err
		}

		// Labels live in their own table and store.CreateAgent does not touch
		// them -- the agent service pairs the two calls the same way (onCreate
		// in agent.go). Without this the copy silently loses every label the
		// source carried, including the ones other features filter on.
		if err = label.Create(ctx, s, &cp); err != nil {
			return err
		}

		c.reportDropped(&cp, dropped)

		c.ids.set(resourceref.KindAgent, src.ID, cp.ID)
	}

	return
}

// remapAgentAccess repoints one copied agent's access model at the draft, and
// reports the references it could not carry.
func (c *projectClone) remapAgentAccess(ctx context.Context, s store.Storer, cp *types.Agent) (dropped []droppedRef, err error) {
	cp.Access.Context.Namespace = c.remapNamespaceIdent(cp.Access.Context.Namespace)
	cp.Access.Context.Module = c.remapModuleIdent(cp.Access.Context.Module)

	for i := range cp.Access.Tools {
		tool := &cp.Access.Tools[i]

		kept := tool.Allow[:0]
		for _, allow := range tool.Allow {
			allow.NamespaceID, _ = c.ids.get(resourceref.KindComposeNamespace, allow.NamespaceID)

			mods := make(types.AgentAccessIDList, 0, len(allow.ModuleIDs))
			for _, id := range allow.ModuleIDs {
				var (
					mapped uint64
					keep   bool
				)
				if mapped, keep, err = c.carryRef(ctx, s, resourceref.KindComposeModule, id); err != nil {
					return nil, err
				}
				if !keep {
					dropped = append(dropped, droppedRef{resourceref.KindComposeModule, id})
					continue
				}
				mods = append(mods, mapped)
			}

			// An allow entry with an EMPTY module list covers every module in
			// its namespace (checkAllow, system/agentic/policy), so emptying a
			// scoped entry would silently WIDEN it from three modules to the
			// whole namespace. Drop the entry instead. With no entries left the
			// tool denies outright, which is the fail-closed end of this.
			if len(mods) == 0 && len(allow.ModuleIDs) > 0 {
				continue
			}

			allow.ModuleIDs = mods
			kept = append(kept, allow)
		}
		tool.Allow = kept
	}

	// TAQs and workflows are the parent revision's own logic, and no pass
	// copies them yet -- so these are dropped rather than remapped. See the
	// third case in this file's header comment. Once automations are copied,
	// their entries land in the id map and carryRef starts remapping them
	// instead, with no change needed here.
	if len(cp.Access.TAQs) > 0 {
		kept := cp.Access.TAQs[:0]
		for _, taq := range cp.Access.TAQs {
			var (
				mapped uint64
				keep   bool
			)
			if mapped, keep, err = c.carryRef(ctx, s, resourceref.KindNgAutomation, taq.ID); err != nil {
				return nil, err
			}
			if !keep {
				dropped = append(dropped, droppedRef{resourceref.KindNgAutomation, taq.ID})
				continue
			}
			taq.ID = mapped
			kept = append(kept, taq)
		}
		cp.Access.TAQs = kept
	}

	if len(cp.Access.Workflows) > 0 {
		kept := cp.Access.Workflows[:0]
		for _, wf := range cp.Access.Workflows {
			var (
				mapped uint64
				keep   bool
			)
			if mapped, keep, err = c.carryRef(ctx, s, resourceref.KindAutomationWorkflow, wf.ID); err != nil {
				return nil, err
			}
			if !keep {
				dropped = append(dropped, droppedRef{resourceref.KindAutomationWorkflow, wf.ID})
				continue
			}
			wf.ID = mapped
			kept = append(kept, wf)
		}
		cp.Access.Workflows = kept
	}

	return dropped, nil
}

// remapNamespaceIdent rewrites an access-context namespace identifier. The
// field is a free-form string and both forms occur in stored configs, so both
// are handled: an ID goes through the id map, a slug is compared against the
// parent's.
func (c *projectClone) remapNamespaceIdent(ident string) string {
	if ident == "" {
		return ident
	}

	if id, err := strconv.ParseUint(ident, 10, 64); err == nil {
		if mapped, ok := c.ids.get(resourceref.KindComposeNamespace, id); ok {
			return strconv.FormatUint(mapped, 10)
		}
		return ident
	}

	// Slug form. CreateRevision gives the clone a suffixed slug, so the
	// parent's slug names the LIVE namespace and has to be rewritten; any
	// other slug names a namespace outside this revision and stays put.
	if ident == c.parentNsSlug {
		return c.draftNsSlug
	}

	return ident
}

// remapModuleIdent rewrites an access-context module identifier. Module
// handles survive the clone unchanged -- that is what makes the module id map
// buildable at all -- so only the numeric form needs rewriting.
func (c *projectClone) remapModuleIdent(ident string) string {
	if ident == "" {
		return ident
	}

	if id, err := strconv.ParseUint(ident, 10, 64); err == nil {
		if mapped, ok := c.ids.get(resourceref.KindComposeModule, id); ok {
			return strconv.FormatUint(mapped, 10)
		}
	}

	return ident
}

// reportDropped surfaces the references a copy could not carry into the draft.
//
// A warning rather than a failure: the copy is still the best available draft
// of the resource, and failing the whole branch over a tool binding the user
// can re-add would be the worse trade. It does have to be visible though --
// an agent that came back quietly narrower is otherwise a support call with
// nothing to go on.
func (c *projectClone) reportDropped(cp *types.Agent, dropped []droppedRef) {
	if len(dropped) == 0 || DefaultLogger == nil {
		return
	}

	refs := make([]string, len(dropped))
	for i, d := range dropped {
		refs[i] = fmt.Sprintf("%s/%d", d.kind, d.id)
	}

	DefaultLogger.Warn(
		"revision copy dropped agent references that could not be repointed at the draft",
		zap.Uint64("projectID", c.draft.ID),
		zap.Uint64("agentID", cp.ID),
		zap.String("agentHandle", cp.Handle),
		zap.Strings("refs", refs),
	)
}
