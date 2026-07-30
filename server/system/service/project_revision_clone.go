package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/spf13/cast"
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
//   - Part of the parent revision but not copied by this pass (classic
//     automation workflows, which no kind here carries yet): DROPPED, and
//     logged. The draft loses a binding the user can re-add; the alternative is
//     a draft agent invoking logic bound to the live revision, which is the
//     failure this file exists to prevent.
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

// lookupAny resolves a reference whose KIND the caller cannot name.
//
// RBAC rule resources are the case this exists for: a rule's path segments are
// the ids of the resource and of each of its parents, one level per segment,
// and which level is which depends on the resource type. IDs are snowflakes
// though -- unique across every kind -- so an id can be resolved without
// knowing what it names, and a segment that matches nothing is left alone.
func (m projectCloneIDs) lookupAny(old uint64) (new uint64, remapped bool) {
	if old == 0 {
		return old, false
	}
	for _, kk := range m {
		if n, ok := kk[old]; ok {
			return n, true
		}
	}
	return old, false
}

// kindProjectAiSystem is a clone-local id-map key. AI systems are not graph
// resources and have no resourceref kind of their own, but FRIA scenarios name
// them by ID and so must remap through the same machinery as everything else.
const kindProjectAiSystem = "human::system:project-ai-system"

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

		// The TAQ copies, kept so their agent references can be repointed once
		// agents have been copied. TAQs and agents reference EACH OTHER, so one
		// of the two directions cannot be resolved in a single ordered pass --
		// see cloneProjectResources.
		copiedTAQs []*automationTypes.NgAutomation
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
	// Connections are named by TAQ steps; TAQs are named by agents and
	// chatbots; agents are named by chatbots; roles' RBAC rules name everything,
	// including the compose resources the namespace clone already minted.
	if err = svc.cloneProjectConnections(ctx, s, c); err != nil {
		return fmt.Errorf("copy configured connections into revision: %w", err)
	}

	if err = svc.cloneProjectTAQs(ctx, s, c); err != nil {
		return fmt.Errorf("copy automations into revision: %w", err)
	}

	if err = svc.cloneProjectAgents(ctx, s, c); err != nil {
		return fmt.Errorf("copy agents into revision: %w", err)
	}

	// The one reference the ordering cannot serve. A TAQ step may invoke an
	// agent (agentRun/agentPrompt) while an agent's tool list names TAQs, so
	// whichever kind is copied first has a reference the map cannot answer yet.
	// TAQs go first -- agents carry far more references and copying them
	// half-remapped would be the worse trade -- and their agent refs are
	// repointed here, once the agent copies exist.
	if err = svc.repointTAQAgentRefs(ctx, s, c); err != nil {
		return fmt.Errorf("repoint automation agent references: %w", err)
	}

	if err = svc.cloneProjectChatbots(ctx, s, c); err != nil {
		return fmt.Errorf("copy chatbots into revision: %w", err)
	}

	if err = svc.cloneProjectAiSystems(ctx, s, c); err != nil {
		return fmt.Errorf("copy AI systems into revision: %w", err)
	}

	if err = svc.cloneProjectReviews(ctx, s, c); err != nil {
		return fmt.Errorf("copy reviews into revision: %w", err)
	}

	if err = svc.cloneProjectRoles(ctx, s, c); err != nil {
		return fmt.Errorf("copy roles into revision: %w", err)
	}

	return
}

// cloneProjectConnections copies the parent's configured connections.
//
// Copied first because everything downstream can name one: a TAQ step invokes a
// connection's generated function as conn_{configurationID}_{operation}, so a
// draft whose TAQs still carried the parent's configuration IDs would call the
// LIVE connection -- with the live credentials, against the live system.
func (svc *project) cloneProjectConnections(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	cc, _, err := store.SearchConfiguredConnections(ctx, s, types.ConfiguredConnectionFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	// human/connection-id and human/connection-revision are stamped by
	// beforeCreate (configured_connection.go) and live in the label table, which
	// the store search does not join in. Without loading them the copies would
	// carry an empty label set and lose the link back to the connection
	// definition they were configured from.
	if err = label.Load(ctx, s, toLabeledConfiguredConnections(cc)...); err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range cc {
		cp := *src.Clone()

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		// The one reference a configured connection holds into the revision:
		// which namespace its discovered data model was written into.
		cp.Config.NamespaceID, _ = c.ids.get(resourceref.KindComposeNamespace, cp.Config.NamespaceID)
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		if err = store.CreateConfiguredConnection(ctx, s, &cp); err != nil {
			return err
		}
		if err = label.Create(ctx, s, &cp); err != nil {
			return err
		}

		c.ids.set(resourceref.KindConfiguredConnection, src.ID, cp.ID)
	}

	return
}

// cloneProjectTAQs copies the parent's project automations.
//
// Disabled: StateInclusive because a project TAQ is created disabled and only
// enabled once its logic is built (the same reason the project graph asks for
// it) -- the default excludes them, which would silently drop every TAQ a
// revision has not switched on yet.
func (svc *project) cloneProjectTAQs(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	aa, _, err := store.SearchAutomationNgAutomations(ctx, s, automationTypes.NgAutomationFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
		Disabled:  filter.StateInclusive,
	})
	if err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range aa {
		// Clone() covers Meta and Scope -- both POINTERS, so a plain struct copy
		// would have the draft and the parent share them. It does NOT cover
		// Steps, which is a slice of POINTERS, and steps are exactly what the
		// remap below rewrites: without the deep copy, repointing the draft's
		// steps would rewrite the parent revision's own automation in place.
		cp := *src.Clone()
		cp.Steps = cloneTAQSteps(src.Steps)

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		c.remapTAQSteps(&cp)

		if err = store.CreateAutomationNgAutomation(ctx, s, &cp); err != nil {
			return err
		}

		c.ids.set(resourceref.KindNgAutomation, src.ID, cp.ID)
		c.copiedTAQs = append(c.copiedTAQs, &cp)
	}

	return
}

// cloneTAQSteps deep-copies a step set, including each step's argument list.
// Everything the remap writes to has to be the copy's own.
func cloneTAQSteps(ss automationTypes.NgAutomationStepSet) automationTypes.NgAutomationStepSet {
	if ss == nil {
		return nil
	}

	out := make(automationTypes.NgAutomationStepSet, 0, len(ss))
	for _, s := range ss {
		if s == nil {
			out = append(out, nil)
			continue
		}

		cp := *s
		if s.Arguments != nil {
			cp.Arguments = make([]*automationTypes.Expr, 0, len(s.Arguments))
			for _, arg := range s.Arguments {
				if arg == nil {
					cp.Arguments = append(cp.Arguments, nil)
					continue
				}
				argCp := *arg
				cp.Arguments = append(cp.Arguments, &argCp)
			}
		}
		out = append(out, &cp)
	}

	return out
}

// remapTAQSteps repoints one copied automation's steps at the draft.
//
// Two kinds of reference live in a step:
//
//   - The step's own ref, for connection-generated functions: the configured
//     connection is encoded IN the function name as conn_{id}_{operation} (see
//     stepResourceRefs in automation/types/resource_refs.go), so the id has to
//     be rewritten inside the string.
//   - Constant arguments bound to a resource-typed parameter -- the module a
//     composeRecords step writes to, the namespace it works in, the agent it
//     invokes. Left alone, a draft TAQ writes to the LIVE module.
//
// Trigger constraints are deliberately untouched: they name modules by HANDLE,
// and handles survive the clone unchanged.
func (c *projectClone) remapTAQSteps(cp *automationTypes.NgAutomation) {
	for _, step := range cp.Steps {
		if step == nil {
			continue
		}

		step.Ref = c.remapConnFunctionRef(step.Ref)

		kinds := automationTypes.NgAutomationStepParamKinds(step)
		if len(kinds) == 0 {
			continue
		}

		for _, arg := range step.Arguments {
			if arg == nil || arg.Value == nil {
				continue
			}

			kind, is := kinds[arg.Target]
			if !is {
				continue
			}

			id, err := strconv.ParseUint(cast.ToString(arg.Value), 10, 64)
			if err != nil {
				// A handle, an expression, or anything else that is not an ID:
				// handles survive the clone, so there is nothing to rewrite.
				continue
			}

			if mapped, ok := c.ids.get(kind, id); ok {
				// Written back as a string on purpose. Every ID in this system
				// is a snowflake, which does not survive JSON's float64, so the
				// stored form is always a string -- and these values round-trip
				// through JSON on every read.
				arg.Value = strconv.FormatUint(mapped, 10)
			}
		}
	}
}

// remapConnFunctionRef rewrites the configured-connection ID embedded in a
// connection-generated function reference (conn_{id}_{operation}).
func (c *projectClone) remapConnFunctionRef(ref string) string {
	rest, found := strings.CutPrefix(ref, "conn_")
	if !found {
		return ref
	}

	idStr, op, found := strings.Cut(rest, "_")
	if !found {
		return ref
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ref
	}

	mapped, ok := c.ids.get(resourceref.KindConfiguredConnection, id)
	if !ok {
		return ref
	}

	return fmt.Sprintf("conn_%d_%s", mapped, op)
}

// repointTAQAgentRefs closes the TAQ/agent reference cycle: the TAQ copies were
// stored before any agent had been copied, so their agentRun/agentPrompt steps
// still named the parent's agents. Now that the agent id map is complete, those
// arguments are rewritten and the changed rows updated.
//
// Runs over the in-memory copies rather than re-reading them, so it sees the
// same structs the create wrote.
func (svc *project) repointTAQAgentRefs(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	if len(c.ids[resourceref.KindAgent]) == 0 {
		return nil
	}

	for _, cp := range c.copiedTAQs {
		changed := false

		for _, step := range cp.Steps {
			if step == nil {
				continue
			}

			kinds := automationTypes.NgAutomationStepParamKinds(step)
			for _, arg := range step.Arguments {
				if arg == nil || arg.Value == nil || kinds[arg.Target] != resourceref.KindAgent {
					continue
				}

				id, parseErr := strconv.ParseUint(cast.ToString(arg.Value), 10, 64)
				if parseErr != nil {
					continue
				}
				if mapped, ok := c.ids.get(resourceref.KindAgent, id); ok {
					arg.Value = strconv.FormatUint(mapped, 10)
					changed = true
				}
			}
		}

		if !changed {
			continue
		}

		if err = store.UpdateAutomationNgAutomation(ctx, s, cp); err != nil {
			return err
		}
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

		c.reportDroppedRefs("agent", cp.ID, cp.Handle, dropped)

		c.ids.set(resourceref.KindAgent, src.ID, cp.ID)
	}

	return
}

// cloneProjectChatbots copies the parent's chatbots into the draft.
//
// Handles are carried over unchanged for the same reason agents' are -- the
// publish diff matches on kind + handle -- which is why the chatbot handle index
// moved from global to per-project (system/chatbot.cue).
//
// The WIDGET KEY does not follow that rule. It is the public embed id a customer
// has pasted into their site, and its uniqueness is global on purpose, so the
// copy takes a fresh throwaway key. The canonical key is handed over at PUBLISH,
// when the outgoing chatbot's key is suffixed out of the way and the newly live
// copy inherits it -- see swapChatbotWidgetKeys. That is what keeps an embed
// working across a publish without ever having two chatbots claim one key.
func (svc *project) cloneProjectChatbots(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	cc, _, err := store.SearchChatbots(ctx, s, types.ChatbotFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	if err = label.Load(ctx, s, toLabeledChatbots(cc)...); err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range cc {
		// Chatbot.Clone() copies Handoff and Styling but NOT Scenarios, which is
		// a slice -- and scenarios are what the remap below rewrites. Without
		// this the draft and the parent would share the backing array and
		// repointing the draft's agent bindings would repoint the parent's.
		cp := *src.Clone()
		cp.Scenarios = append(types.ChatbotScenarios(nil), src.Scenarios...)

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		if cp.WidgetKey, err = generateChatbotWidgetKey(); err != nil {
			return err
		}
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		var dropped []droppedRef
		if dropped, err = c.remapChatbotRefs(ctx, s, &cp); err != nil {
			return err
		}

		if err = store.CreateChatbot(ctx, s, &cp); err != nil {
			return err
		}
		if err = label.Create(ctx, s, &cp); err != nil {
			return err
		}

		c.reportDroppedRefs("chatbot", cp.ID, cp.Handle, dropped)

		c.ids.set(resourceref.KindChatbot, src.ID, cp.ID)
	}

	return
}

// remapChatbotRefs repoints a copied chatbot's agent bindings and automation
// hooks at the draft.
func (c *projectClone) remapChatbotRefs(ctx context.Context, s store.Storer, cp *types.Chatbot) (dropped []droppedRef, err error) {
	if dropped, err = c.remapChatbotHook(ctx, s, &cp.Handoff.Automation.OnRequested, dropped); err != nil {
		return nil, err
	}
	if dropped, err = c.remapChatbotHook(ctx, s, &cp.Handoff.Automation.OnAccepted, dropped); err != nil {
		return nil, err
	}

	for i := range cp.Scenarios {
		sc := &cp.Scenarios[i]

		var (
			mapped uint64
			keep   bool
		)
		if mapped, keep, err = c.carryRef(ctx, s, resourceref.KindAgent, sc.AgentID); err != nil {
			return nil, err
		}
		if !keep {
			// A scenario pointing at an agent of the parent revision that no
			// copy exists for is cleared rather than left in place: an unbound
			// scenario is visibly incomplete in the editor, whereas one still
			// naming the live agent looks correct and is not.
			dropped = append(dropped, droppedRef{resourceref.KindAgent, sc.AgentID})
			sc.AgentID = 0
		} else {
			sc.AgentID = mapped
		}

		if dropped, err = c.remapChatbotHook(ctx, s, &sc.Automation.Before, dropped); err != nil {
			return nil, err
		}
		if dropped, err = c.remapChatbotHook(ctx, s, &sc.Automation.After, dropped); err != nil {
			return nil, err
		}
	}

	return dropped, nil
}

// remapChatbotHook rewrites one automation hook. The hook names its automation
// as the string "corteza::automation:ng-automation/{id}" (see
// chatbotSession.invokeAutomation), so the ID is rewritten inside the ref.
func (c *projectClone) remapChatbotHook(
	ctx context.Context,
	s store.Storer,
	h *types.ChatbotAutomationHook,
	dropped []droppedRef,
) ([]droppedRef, error) {
	if h.Automation == "" {
		return dropped, nil
	}

	prefix := resourceref.KindNgAutomation + "/"
	idStr, found := strings.CutPrefix(h.Automation, prefix)
	if !found {
		return dropped, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return dropped, nil
	}

	mapped, keep, err := c.carryRef(ctx, s, resourceref.KindNgAutomation, id)
	if err != nil {
		return nil, err
	}
	if !keep {
		h.Automation = ""
		return append(dropped, droppedRef{resourceref.KindNgAutomation, id}), nil
	}

	h.Automation = prefix + strconv.FormatUint(mapped, 10)
	return dropped, nil
}

// cloneProjectAiSystems copies the parent's AI systems, the resource entries
// that say what each one is made of, and the FRIA scenarios assessed against
// them.
//
// AI systems and FRIA scenarios are PER REVISION -- their ProjectID is the
// revision's own id, not the chain root the work items use -- so a branch that
// did not copy them produced a draft whose whole Govern step was empty, and
// whose deployment plan the reviewer had nothing to read against.
func (svc *project) cloneProjectAiSystems(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	ss, _, err := store.SearchProjectAiSystems(ctx, s, types.ProjectAiSystemFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	for _, src := range ss {
		cp := *src.Clone()

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		cp.CreatedAt = *now()
		cp.UpdatedAt, cp.DeletedAt = nil, nil

		if err = store.CreateProjectAiSystem(ctx, s, &cp); err != nil {
			return err
		}

		c.ids.set(kindProjectAiSystem, src.ID, cp.ID)

		if err = svc.cloneProjectAiSystemEntries(ctx, s, c, src.ID, cp.ID); err != nil {
			return err
		}
	}

	return svc.cloneProjectFriaScenarios(ctx, s, c)
}

// cloneProjectAiSystemEntries copies the membership rows that record which
// project resources make up one AI system.
//
// Each entry names its resource as a "kind/id" string, so it is remapped through
// the same id map every other reference uses. An entry naming something the copy
// does not own -- an entry for a resource kind this pass does not carry, say --
// is dropped: the alternative is an AI system in the draft that claims a
// resource belonging to the live revision, which is precisely the misstatement
// the FRIA record must not make.
func (svc *project) cloneProjectAiSystemEntries(
	ctx context.Context,
	s store.Storer,
	c *projectClone,
	srcSystemID, newSystemID uint64,
) (err error) {
	ee, _, err := store.SearchProjectAiSystemEntrys(ctx, s, types.ProjectAiSystemEntryFilter{
		ProjectAiSystemID: srcSystemID,
	})
	if err != nil {
		return err
	}

	for _, src := range ee {
		ref, ok := c.remapResourceRef(src.ResourceRef)
		if !ok {
			if DefaultLogger != nil {
				DefaultLogger.Warn(
					"revision copy dropped an AI system entry that could not be repointed at the draft",
					zap.Uint64("projectID", c.draft.ID),
					zap.Uint64("projectAiSystemID", newSystemID),
					zap.String("resourceRef", src.ResourceRef),
				)
			}
			continue
		}

		cp := *src.Clone()
		cp.ID = nextID()
		cp.ProjectAiSystemID = newSystemID
		cp.ResourceRef = ref
		cp.CreatedAt = *now()

		if err = store.CreateProjectAiSystemEntry(ctx, s, &cp); err != nil {
			return err
		}
	}

	return
}

// cloneProjectFriaScenarios copies the FRIA scenarios, rebinding each to the
// copied AI system it assesses.
func (svc *project) cloneProjectFriaScenarios(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	ff, _, err := store.SearchProjectFriaScenarios(ctx, s, types.ProjectFriaScenarioFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range ff {
		cp := *src.Clone()

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		// A scenario is an assessment OF an AI system; left pointing at the
		// parent's it would be filed under the live revision's system and vanish
		// from the draft's own FRIA.
		cp.AiSystemID, _ = c.ids.get(kindProjectAiSystem, src.AiSystemID)
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		if err = store.CreateProjectFriaScenario(ctx, s, &cp); err != nil {
			return err
		}
	}

	return
}

// cloneProjectReviews copies the reviews bound to the parent revision.
//
// A review is a work item: its ProjectID is the CHAIN ROOT, shared by every
// revision, and RevisionID is what ties it to one of them. So the copy keeps
// ProjectID and rebinds RevisionID -- the reverse of every other kind here.
//
// Reviews with no revision (RevisionID 0) are deliberately not copied: those are
// chain-wide by design and the draft already sees them.
func (svc *project) cloneProjectReviews(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	rr, _, err := store.SearchProjectReviews(ctx, s, types.ProjectReviewFilter{
		RevisionID: c.parent.ID,
		Deleted:    filter.StateExcluded,
	})
	if err != nil {
		return err
	}

	invoker := a.GetIdentityFromContext(ctx).Identity()

	for _, src := range rr {
		cp := *src.Clone()

		cp.ID = nextID()
		cp.RevisionID = c.draft.ID
		cp.CreatedAt = *now()
		cp.CreatedBy = invoker
		cp.UpdatedAt, cp.UpdatedBy = nil, 0
		cp.DeletedAt, cp.DeletedBy = nil, 0

		if err = store.CreateProjectReview(ctx, s, &cp); err != nil {
			return err
		}
	}

	return
}

// cloneProjectRoles copies the parent's project roles, their members and their
// RBAC rules.
//
// Roles are PER REVISION, deliberately -- unlike project membership, which is
// chain-wide. A revision's access model is part of what a revision IS: it is
// reviewed with the revision and deployed with it, and a draft that edited the
// live revision's permissions would be changing production before anyone
// approved it.
//
// Handles are the one thing that cannot be carried over: they are minted as
// proj_<projectID>_<name> by the wizard and a role handle is not project-scoped,
// so the copy re-mints its own with the draft's id. That is why the publish diff
// compares project roles on the SUFFIX -- see diffSources.
func (svc *project) cloneProjectRoles(ctx context.Context, s store.Storer, c *projectClone) (err error) {
	rr, _, err := store.SearchRoles(ctx, s, types.RoleFilter{
		ProjectID: c.parent.ID,
		Deleted:   filter.StateExcluded,
		Archived:  filter.StateInclusive,
	})
	if err != nil {
		return err
	}
	if len(rr) == 0 {
		return nil
	}

	if err = label.Load(ctx, s, toLabeledRoles(rr)...); err != nil {
		return err
	}

	// Rules carry no project scope and the store filter cannot narrow by role,
	// so they are loaded once and grouped in memory -- the same shape the
	// project graph uses.
	rules, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	if err != nil {
		return err
	}
	byRole := make(map[uint64]rbac.RuleSet, len(rr))
	for _, rule := range rules {
		byRole[rule.RoleID] = append(byRole[rule.RoleID], rule)
	}

	for _, src := range rr {
		cp := *src.Clone()

		cp.ID = nextID()
		cp.ProjectID = c.draft.ID
		cp.Handle = projectRoleHandle(c.draft.ID, src.Handle)
		cp.CreatedAt = *now()
		cp.UpdatedAt, cp.DeletedAt = nil, nil

		if err = store.CreateRole(ctx, s, &cp); err != nil {
			return err
		}
		if err = label.Create(ctx, s, &cp); err != nil {
			return err
		}

		c.ids.set(resourceref.KindRole, src.ID, cp.ID)

		if err = svc.cloneRoleMembers(ctx, s, src.ID, cp.ID); err != nil {
			return err
		}
		if err = svc.cloneRoleRules(ctx, s, c, byRole[src.ID], cp.ID); err != nil {
			return err
		}
	}

	return
}

// cloneRoleMembers carries a copied role's membership over.
//
// A role without its members is not a copy of the access model, it is an empty
// shell that grants nothing -- and the draft would report every project user as
// gone. Membership is a user reference, which is shared infrastructure, so the
// user ids are carried across untouched.
func (svc *project) cloneRoleMembers(ctx context.Context, s store.Storer, srcRoleID, newRoleID uint64) (err error) {
	mm, _, err := store.SearchRoleMembers(ctx, s, types.RoleMemberFilter{RoleID: srcRoleID})
	if err != nil {
		return err
	}

	for _, m := range mm {
		if err = store.CreateRoleMember(ctx, s, &types.RoleMember{
			RoleID:   newRoleID,
			Resource: m.Resource,
		}); err != nil {
			return err
		}
	}

	return
}

// cloneRoleRules rewrites a role's RBAC rules onto the draft's own resources.
//
// A rule names its subject as an rbac resource string -- "kind/parentID/id",
// with a level per path segment -- so every numeric segment is looked up in the
// id map. IDs are snowflakes and unique across kinds, so a segment can be
// resolved without knowing which level it is.
//
// A segment the copy does not own is left as it is: plenty of rules grant on
// resources outside the revision (a component-level grant, a shared connection),
// and rewriting or dropping those would silently change what the role can do.
func (svc *project) cloneRoleRules(
	ctx context.Context,
	s store.Storer,
	c *projectClone,
	rules rbac.RuleSet,
	newRoleID uint64,
) (err error) {
	for _, src := range rules {
		cp := rbac.Rule{
			RoleID:    newRoleID,
			Resource:  c.remapRbacResource(src.Resource),
			Operation: src.Operation,
			Access:    src.Access,
		}

		if err = store.CreateRbacRule(ctx, s, &cp); err != nil {
			return err
		}
	}

	return
}

// projectRoleHandle re-mints a project role handle for the given revision.
//
// Wizard-created handles are proj_<projectID>_<name> (see the projects store in
// the webapp), so the revision id is stripped and replaced. A handle in any other
// shape is taken whole as the suffix rather than guessed at -- it still has to be
// made revision-unique, because role handles are globally unique.
func projectRoleHandle(projectID uint64, handle string) string {
	suffix := handle

	if rest, found := strings.CutPrefix(handle, "proj_"); found {
		if _, name, found := strings.Cut(rest, "_"); found {
			suffix = name
		}
	}

	return fmt.Sprintf("proj_%d_%s", projectID, suffix)
}

// remapResourceRef rewrites a "kind/id" resource reference through the id map.
// Reports false when the reference names something of the parent revision that
// this pass did not copy, so the caller can drop it.
func (c *projectClone) remapResourceRef(ref string) (string, bool) {
	kind, idStr, found := strings.Cut(ref, "/")
	if !found {
		return ref, true
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ref, true
	}

	if mapped, ok := c.ids.get(kind, id); ok {
		return kind + "/" + strconv.FormatUint(mapped, 10), true
	}

	// Not in the map. Either it is shared with the rest of the system (keep) or
	// it belongs to the parent revision and nothing copied it (drop).
	if c.owned[kind][id] {
		return ref, false
	}

	return ref, true
}

// remapRbacResource rewrites every ID segment of an rbac resource string.
func (c *projectClone) remapRbacResource(res string) string {
	pp := strings.Split(res, "/")
	if len(pp) < 2 {
		return res
	}

	for i := 1; i < len(pp); i++ {
		id, err := strconv.ParseUint(pp[i], 10, 64)
		if err != nil {
			// A wildcard, or an empty segment. Both stay.
			continue
		}
		if mapped, ok := c.ids.lookupAny(id); ok {
			pp[i] = strconv.FormatUint(mapped, 10)
		}
	}

	return strings.Join(pp, "/")
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

	// TAQs are copied by cloneProjectTAQs, which runs first, so these remap
	// through the id map. Classic workflows are not copied by any pass yet, so a
	// parent-owned one is dropped instead -- the third case in this file's
	// header comment. Both go through carryRef, which is what makes the two
	// behave differently without either loop knowing about it.
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
func (c *projectClone) reportDroppedRefs(kind string, id uint64, handle string, dropped []droppedRef) {
	if len(dropped) == 0 || DefaultLogger == nil {
		return
	}

	refs := make([]string, len(dropped))
	for i, d := range dropped {
		refs[i] = fmt.Sprintf("%s/%d", d.kind, d.id)
	}

	DefaultLogger.Warn(
		"revision copy dropped references that could not be repointed at the draft",
		zap.Uint64("projectID", c.draft.ID),
		zap.String("resourceKind", kind),
		zap.Uint64("resourceID", id),
		zap.String("resourceHandle", handle),
		zap.Strings("refs", refs),
	)
}
