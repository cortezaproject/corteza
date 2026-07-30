package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// seedCloneModule adds a module to a project's namespace. Store-level on
// purpose: the compose service would drag in guards and DAL wiring the copy has
// no opinion about, and the copy only ever sees what is in the store.
//
// ProjectID is set the way compose sets it (from the namespace's project), so
// the clone has something to carry across wrongly if the repoint regresses.
func seedCloneModule(t *testing.T, s store.Storer, projectID, nsID uint64, handle, name string) *composeTypes.Module {
	t.Helper()

	m := &composeTypes.Module{
		ID:          nextID(),
		ProjectID:   projectID,
		NamespaceID: nsID,
		Handle:      handle,
		Name:        name,
		CreatedAt:   time.Now(),
	}
	require.NoError(t, store.CreateComposeModule(context.Background(), s, m))
	return m
}

// seedCloneTAQ adds an ng-automation. projectID 0 stands for one that lives
// outside any project — shared infrastructure the draft must keep using.
func seedCloneTAQ(t *testing.T, s store.Storer, projectID uint64, handle string) *automationTypes.NgAutomation {
	t.Helper()

	au := &automationTypes.NgAutomation{
		ID:        nextID(),
		ProjectID: projectID,
		Handle:    handle,
		Meta:      &automationTypes.NgAutomationMeta{Short: handle},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAutomationNgAutomation(context.Background(), s, au))
	return au
}

func seedCloneWorkflow(t *testing.T, s store.Storer, projectID uint64, handle string) *automationTypes.Workflow {
	t.Helper()

	wf := &automationTypes.Workflow{
		ID:        nextID(),
		ProjectID: projectID,
		Handle:    handle,
		Meta:      &automationTypes.WorkflowMeta{Name: handle},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAutomationWorkflow(context.Background(), s, wf))
	return wf
}

// loadCopiedAgent returns the single agent copied into a revision.
func loadCopiedAgent(t *testing.T, s store.Storer, projectID uint64) *types.Agent {
	t.Helper()

	aa, _, err := store.SearchAgents(context.Background(), s, types.AgentFilter{ProjectID: projectID})
	require.NoError(t, err)
	require.Len(t, aa, 1, "the revision must hold exactly one copied agent")
	return aa[0]
}

// TestCreateRevision_AgentReferencesFollowTheDraft is the bug this file exists
// for: agents were copied into the draft with every reference left pointing at
// the parent revision, so a draft agent read and wrote the LIVE namespace and
// modules — worse than not copying it at all.
//
// It pins all three outcomes a reference can have (see project_revision_clone.go):
// copied-with-the-revision refs are remapped, shared refs are left alone, and
// refs into the parent revision that nothing copies yet are dropped.
func TestCreateRevision_AgentReferencesFollowTheDraft(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-clone-refs")
	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	require.NoError(t, err)

	customers := seedCloneModule(t, s, parent.ID, parentNs.ID, "customers", "Customers")
	invoices := seedCloneModule(t, s, parent.ID, parentNs.ID, "invoices", "Invoices")

	ownedTAQ := seedCloneTAQ(t, s, parent.ID, "clone-refs-owned-taq")
	sharedTAQ := seedCloneTAQ(t, s, 0, "clone-refs-shared-taq")
	ownedWf := seedCloneWorkflow(t, s, parent.ID, "clone-refs-owned-wf")
	sharedWf := seedCloneWorkflow(t, s, 0, "clone-refs-shared-wf")

	// An LLM provider is shared across projects, so its ID must survive the
	// copy untouched — the counter-example to everything else asserted here.
	const sharedProviderID = uint64(918273645)

	src := &types.Agent{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "clone-refs-agent",
		Status:    "active",
		// A source that has been edited a few times; the copy is a new agent
		// and must not inherit the counter.
		Revision:  7,
		Meta:      types.AgentMeta{Short: "Reference carrier"},
		Execution: types.AgentExecution{Model: types.AgentExecutionModel{LLMProviderID: sharedProviderID}},
		Access: types.AgentAccess{
			Context: types.AgentAccessContext{Namespace: parentNs.Slug, Module: customers.Handle},
			Tools: []types.AgentAccessTool{
				{
					Name:  "compose_record_search",
					Allow: []types.AgentAccessAllow{{NamespaceID: parentNs.ID, ModuleIDs: types.AgentAccessIDList{customers.ID, invoices.ID}}},
				},
				{
					// No module ids: covers the whole namespace.
					Name:  "compose_module_lookup",
					Allow: []types.AgentAccessAllow{{NamespaceID: parentNs.ID}},
				},
			},
			TAQs:      []types.AgentAccessTAQ{{ID: ownedTAQ.ID}, {ID: sharedTAQ.ID}},
			Workflows: []types.AgentAccessWorkflow{{ID: ownedWf.ID}, {ID: sharedWf.ID}},
		},
		Labels:    map[string]labelTypes.LabelValue{"tier": {Val: "gold"}},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAgent(ctx, s, src))
	require.NoError(t, label.Create(ctx, s, src))

	rev, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	draftNs, err := store.LookupComposeNamespaceByID(ctx, s, rev.Config.NamespaceID)
	require.NoError(t, err)

	draftMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: draftNs.ID})
	require.NoError(t, err)
	draftModByHandle := make(map[string]uint64, len(draftMods))
	for _, m := range draftMods {
		draftModByHandle[m.Handle] = m.ID
	}
	require.Len(t, draftModByHandle, 2)

	cp := loadCopiedAgent(t, s, rev.ID)

	require.Equal(t, rev.ID, cp.ProjectID)
	require.NotEqual(t, src.ID, cp.ID)
	require.Equal(t, src.Handle, cp.Handle, "the publish diff matches across revisions on kind + handle")
	require.Equal(t, 1, cp.Revision, "a copy is a new agent on its first version")

	// Compose refs: remapped onto the clone.
	require.Equal(t, draftNs.Slug, cp.Access.Context.Namespace)
	require.Equal(t, customers.Handle, cp.Access.Context.Module, "module handles survive the clone unchanged")

	require.Equal(t, draftNs.ID, cp.Access.Tools[0].Allow[0].NamespaceID)
	require.Equal(t,
		types.AgentAccessIDList{draftModByHandle["customers"], draftModByHandle["invoices"]},
		cp.Access.Tools[0].Allow[0].ModuleIDs,
	)

	// A namespace-wide entry stays namespace-wide, on the draft's namespace.
	require.Equal(t, draftNs.ID, cp.Access.Tools[1].Allow[0].NamespaceID)
	require.Empty(t, cp.Access.Tools[1].Allow[0].ModuleIDs)

	// Shared refs: untouched.
	require.Equal(t, sharedProviderID, cp.Execution.Model.LLMProviderID)

	// The parent's own TAQ is copied by the branch now, so the agent's binding
	// follows the copy rather than being dropped — the shared one still points
	// at the same row it always did.
	draftTAQs, _, err := store.SearchAutomationNgAutomations(ctx, s, automationTypes.NgAutomationFilter{
		ProjectID: rev.ID,
		Disabled:  filter.StateInclusive,
	})
	require.NoError(t, err)
	require.Len(t, draftTAQs, 1)
	require.NotEqual(t, ownedTAQ.ID, draftTAQs[0].ID)
	require.Equal(t, ownedTAQ.Handle, draftTAQs[0].Handle, "the publish diff matches on kind + handle")

	require.Equal(t,
		[]types.AgentAccessTAQ{{ID: draftTAQs[0].ID}, {ID: sharedTAQ.ID}},
		cp.Access.TAQs,
	)

	// Workflows are still not copied by this pass, so a parent-owned one is
	// dropped rather than left pointing at the live revision's logic.
	require.Equal(t, []types.AgentAccessWorkflow{{ID: sharedWf.ID}}, cp.Access.Workflows)

	require.NoError(t, label.Load(ctx, s, cp))
	require.Equal(t, map[string]labelTypes.LabelValue{"tier": {Val: "gold"}}, cp.Labels)
}

// TestCreateRevision_ClonedResourcesBelongToTheDraft covers the other half of
// what a branch has to repoint. The envoy clone re-points every child at the
// new namespace but copies ProjectID across verbatim, so the draft's modules
// still claimed to belong to the parent — which made the draft look empty to
// everything that enumerates a project by ProjectID (the resource graph, and
// the deployment plan built on it) and made the parent look like it held every
// resource twice.
func TestCreateRevision_ClonedResourcesBelongToTheDraft(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-clone-scope")
	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	require.NoError(t, err)

	seedCloneModule(t, s, parent.ID, parentNs.ID, "scope-customers", "Customers")

	rev, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	draftMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftMods, 1, "the draft's own modules must be findable by its project id")
	require.Equal(t, rev.Config.NamespaceID, draftMods[0].NamespaceID)

	parentMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{ProjectID: parent.ID})
	require.NoError(t, err)
	require.Len(t, parentMods, 1, "the parent must not gain a second copy of its own module")
	require.Equal(t, parentNs.ID, parentMods[0].NamespaceID)
}

// TestCreateRevision_AgentCopyLeavesTheSourceAlone guards the deep copy. The
// remap writes into Access.Tools[].Allow and the TAQ/workflow slices, which a
// plain struct copy still shares with the source — so remapping the copy would
// silently repoint the PARENT revision's agent at the draft, i.e. break the
// live project as a side effect of branching it.
func TestCreateRevision_AgentCopyLeavesTheSourceAlone(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-clone-source")
	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	require.NoError(t, err)

	customers := seedCloneModule(t, s, parent.ID, parentNs.ID, "source-customers", "Customers")
	ownedTAQ := seedCloneTAQ(t, s, parent.ID, "clone-source-owned-taq")

	src := &types.Agent{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "clone-source-agent",
		Status:    "active",
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{
				Name:  "compose_record_search",
				Allow: []types.AgentAccessAllow{{NamespaceID: parentNs.ID, ModuleIDs: types.AgentAccessIDList{customers.ID}}},
			}},
			TAQs: []types.AgentAccessTAQ{{ID: ownedTAQ.ID}},
		},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAgent(ctx, s, src))

	_, err = svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	reloaded, err := store.LookupAgentByID(ctx, s, src.ID)
	require.NoError(t, err)

	require.Equal(t, parentNs.ID, reloaded.Access.Tools[0].Allow[0].NamespaceID)
	require.Equal(t, types.AgentAccessIDList{customers.ID}, reloaded.Access.Tools[0].Allow[0].ModuleIDs)
	require.Equal(t, []types.AgentAccessTAQ{{ID: ownedTAQ.ID}}, reloaded.Access.TAQs)
}

// TestCreateRevision_AgentAllowNarrowsRatherThanWidens covers the one case
// where dropping a module id is not enough: an allow entry with an EMPTY
// module list covers every module in its namespace (checkAllow, in
// system/agentic/policy), so emptying a scoped entry would hand the copy
// BROADER access than the source had. The entry has to go instead.
//
// The unmappable module here has no handle, which is the only thing the two
// namespaces are matched on — so it is provably part of the parent revision
// and provably has no counterpart in the draft.
func TestCreateRevision_AgentAllowNarrowsRatherThanWidens(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-clone-narrow")
	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	require.NoError(t, err)

	unmappable := seedCloneModule(t, s, parent.ID, parentNs.ID, "", "No handle")

	src := &types.Agent{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "clone-narrow-agent",
		Status:    "active",
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{
				Name:  "compose_record_search",
				Allow: []types.AgentAccessAllow{{NamespaceID: parentNs.ID, ModuleIDs: types.AgentAccessIDList{unmappable.ID}}},
			}},
		},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAgent(ctx, s, src))

	rev, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	cp := loadCopiedAgent(t, s, rev.ID)
	require.Empty(t, cp.Access.Tools[0].Allow,
		"an allow entry left with no modules must be dropped, not kept as a namespace-wide grant")
}

// TestCreateRevision_CopiesEveryProjectScopedKind walks a branch end to end
// across every kind the copy carries, and checks the one thing that matters for
// each of them: that its references name the DRAFT's resources.
//
// The failure this guards against is not "the draft is missing something" — it
// is worse than that. A draft holding a copy that still names the parent's
// connection, module or agent looks complete in the editor and drives the LIVE
// revision the moment it runs.
func TestCreateRevision_CopiesEveryProjectScopedKind(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-clone-kinds")
	parentNs, err := store.LookupComposeNamespaceByID(ctx, s, parent.Config.NamespaceID)
	require.NoError(t, err)

	orders := seedCloneModule(t, s, parent.ID, parentNs.ID, "orders", "Orders")

	// A configured connection, with the labels beforeCreate stamps on it.
	conn := &types.ConfiguredConnection{
		ID:        nextID(),
		ProjectID: parent.ID,
		Name:      "Warehouse",
		Status:    "installed",
		Config:    types.ConfiguredConnectionConfig{NamespaceID: parentNs.ID},
		Labels:    map[string]labelTypes.LabelValue{"human/connection-id": {Val: "42"}},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateConfiguredConnection(ctx, s, conn))
	require.NoError(t, label.Create(ctx, s, conn))

	agent := &types.Agent{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "kinds-agent",
		Status:    "active",
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAgent(ctx, s, agent))

	// A TAQ carrying all three reference shapes a step can hold: a connection
	// encoded in the function name, a module named by a constant argument, and
	// an agent — which is the reference the copy order cannot serve directly,
	// since agents are copied after TAQs.
	taq := &automationTypes.NgAutomation{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "kinds-taq",
		Meta:      &automationTypes.NgAutomationMeta{Short: "kinds-taq"},
		Steps: automationTypes.NgAutomationStepSet{
			{ID: nextID(), Ref: fmt.Sprintf("conn_%d_list", conn.ID)},
			{ID: nextID(), Ref: "composeRecordsSearch", Arguments: []*automationTypes.Expr{
				{Target: "module", Value: strconv.FormatUint(orders.ID, 10)},
				{Target: "namespace", Value: strconv.FormatUint(parentNs.ID, 10)},
			}},
			{ID: nextID(), Ref: "agentRun", Arguments: []*automationTypes.Expr{
				{Target: "agentID", Value: strconv.FormatUint(agent.ID, 10)},
			}},
		},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateAutomationNgAutomation(ctx, s, taq))

	hook := resourceref.KindNgAutomation + "/" + strconv.FormatUint(taq.ID, 10)
	chatbot := &types.Chatbot{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "kinds-chatbot",
		Name:      "Kinds",
		WidgetKey: "kinds-original-widget-key",
		Handoff: types.ChatbotHandoff{
			Automation: types.ChatbotHandoffAutomation{
				OnRequested: types.ChatbotAutomationHook{Automation: hook},
			},
		},
		Scenarios: types.ChatbotScenarios{{
			ID:      "s1",
			Type:    "agent",
			AgentID: agent.ID,
			Automation: types.ChatbotScenarioAutomation{
				After: types.ChatbotAutomationHook{Automation: hook},
			},
		}},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateChatbot(ctx, s, chatbot))

	aiSystem := &types.ProjectAiSystem{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "kinds-system",
		RiskClass: "high",
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateProjectAiSystem(ctx, s, aiSystem))
	require.NoError(t, store.CreateProjectAiSystemEntry(ctx, s, &types.ProjectAiSystemEntry{
		ID:                nextID(),
		ProjectAiSystemID: aiSystem.ID,
		ResourceRef:       resourceref.KindAgent + "/" + strconv.FormatUint(agent.ID, 10),
		CreatedAt:         time.Now(),
	}))
	require.NoError(t, store.CreateProjectFriaScenario(ctx, s, &types.ProjectFriaScenario{
		ID:         nextID(),
		ProjectID:  parent.ID,
		AiSystemID: aiSystem.ID,
		Title:      "Wrongful refusal",
		Severity:   "high",
		CreatedAt:  time.Now(),
	}))

	require.NoError(t, store.CreateProjectReview(ctx, s, &types.ProjectReview{
		ID:         nextID(),
		ProjectID:  parent.RootProjectID(),
		RevisionID: parent.ID,
		Title:      "Pre-deployment review",
		Status:     "pending",
		CreatedAt:  time.Now(),
	}))

	role := &types.Role{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    fmt.Sprintf("proj_%d_reviewer", parent.ID),
		Name:      "Reviewer",
		Meta:      &types.RoleMeta{},
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateRole(ctx, s, role))
	require.NoError(t, store.CreateRoleMember(ctx, s, &types.RoleMember{RoleID: role.ID, Resource: "12345"}))
	require.NoError(t, store.CreateRbacRule(ctx, s,
		&rbac.Rule{
			RoleID:    role.ID,
			Resource:  fmt.Sprintf("corteza::compose:module/%d/%d", parentNs.ID, orders.ID),
			Operation: "read",
			Access:    rbac.Allow,
		},
		// Names nothing the revision owns: it must survive untouched.
		&rbac.Rule{RoleID: role.ID, Resource: "corteza::system:user/*", Operation: "read", Access: rbac.Allow},
	))

	rev, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	draftNs, err := store.LookupComposeNamespaceByID(ctx, s, rev.Config.NamespaceID)
	require.NoError(t, err)
	draftMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: draftNs.ID})
	require.NoError(t, err)
	require.Len(t, draftMods, 1)
	draftOrders := draftMods[0]

	// Connection: repointed at the draft's namespace, labels carried.
	draftConns, _, err := store.SearchConfiguredConnections(ctx, s, types.ConfiguredConnectionFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftConns, 1)
	draftConn := draftConns[0]
	require.NotEqual(t, conn.ID, draftConn.ID)
	require.Equal(t, draftNs.ID, draftConn.Config.NamespaceID)
	require.NoError(t, label.Load(ctx, s, draftConn))
	require.Equal(t, "42", draftConn.Labels["human/connection-id"].Val)

	// TAQ: every step reference follows the draft, including the agent one,
	// which is only resolvable after the agents pass has run.
	draftTAQs, _, err := store.SearchAutomationNgAutomations(ctx, s, automationTypes.NgAutomationFilter{
		ProjectID: rev.ID,
		Disabled:  filter.StateInclusive,
	})
	require.NoError(t, err)
	require.Len(t, draftTAQs, 1)
	draftTAQ := draftTAQs[0]
	require.Equal(t, taq.Handle, draftTAQ.Handle)

	draftAgent := loadCopiedAgent(t, s, rev.ID)

	require.Equal(t, fmt.Sprintf("conn_%d_list", draftConn.ID), draftTAQ.Steps[0].Ref)
	require.Equal(t, strconv.FormatUint(draftOrders.ID, 10), draftTAQ.Steps[1].Arguments[0].Value)
	require.Equal(t, strconv.FormatUint(draftNs.ID, 10), draftTAQ.Steps[1].Arguments[1].Value)
	require.Equal(t, strconv.FormatUint(draftAgent.ID, 10), draftTAQ.Steps[2].Arguments[0].Value)

	// …and the parent's own TAQ is left exactly as it was. The steps are a slice
	// of POINTERS, so a shallow copy would have rewritten these in place.
	srcTAQ, err := store.LookupAutomationNgAutomationByID(ctx, s, taq.ID)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("conn_%d_list", conn.ID), srcTAQ.Steps[0].Ref)
	require.Equal(t, strconv.FormatUint(orders.ID, 10), srcTAQ.Steps[1].Arguments[0].Value)
	require.Equal(t, strconv.FormatUint(agent.ID, 10), srcTAQ.Steps[2].Arguments[0].Value)

	// Chatbot: same handle (the publish diff needs it), a FRESH widget key (the
	// live one still holds the canonical one), agent + hooks repointed.
	draftBots, _, err := store.SearchChatbots(ctx, s, types.ChatbotFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftBots, 1)
	draftBot := draftBots[0]
	require.Equal(t, chatbot.Handle, draftBot.Handle)
	require.NotEmpty(t, draftBot.WidgetKey)
	require.NotEqual(t, chatbot.WidgetKey, draftBot.WidgetKey, "the embed id stays with the live chatbot until publish")
	require.Equal(t, draftAgent.ID, draftBot.Scenarios[0].AgentID)

	draftHook := resourceref.KindNgAutomation + "/" + strconv.FormatUint(draftTAQ.ID, 10)
	require.Equal(t, draftHook, draftBot.Handoff.Automation.OnRequested.Automation)
	require.Equal(t, draftHook, draftBot.Scenarios[0].Automation.After.Automation)

	srcBot, err := store.LookupChatbotByID(ctx, s, chatbot.ID)
	require.NoError(t, err)
	require.Equal(t, agent.ID, srcBot.Scenarios[0].AgentID, "the copy must not repoint the live chatbot")

	// AI system, its entries and the FRIA scenarios assessed against it: all
	// per-revision, all rebound to the draft's own rows.
	draftSystems, _, err := store.SearchProjectAiSystems(ctx, s, types.ProjectAiSystemFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftSystems, 1)
	draftSystem := draftSystems[0]
	require.Equal(t, aiSystem.Handle, draftSystem.Handle)

	draftEntries, _, err := store.SearchProjectAiSystemEntrys(ctx, s, types.ProjectAiSystemEntryFilter{
		ProjectAiSystemID: draftSystem.ID,
	})
	require.NoError(t, err)
	require.Len(t, draftEntries, 1)
	require.Equal(t,
		resourceref.KindAgent+"/"+strconv.FormatUint(draftAgent.ID, 10),
		draftEntries[0].ResourceRef,
	)

	draftScenarios, _, err := store.SearchProjectFriaScenarios(ctx, s, types.ProjectFriaScenarioFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftScenarios, 1)
	require.Equal(t, draftSystem.ID, draftScenarios[0].AiSystemID)

	// Review: bound to the draft, still filed under the chain root.
	draftReviews, _, err := store.SearchProjectReviews(ctx, s, types.ProjectReviewFilter{RevisionID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftReviews, 1)
	require.Equal(t, parent.RootProjectID(), draftReviews[0].ProjectID)

	// Role: re-minted handle, members carried, rules rewritten onto the draft's
	// module — and the rule naming nothing in the revision left alone.
	draftRoles, _, err := store.SearchRoles(ctx, s, types.RoleFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftRoles, 1)
	draftRole := draftRoles[0]
	require.Equal(t, fmt.Sprintf("proj_%d_reviewer", rev.ID), draftRole.Handle)

	members, _, err := store.SearchRoleMembers(ctx, s, types.RoleMemberFilter{RoleID: draftRole.ID})
	require.NoError(t, err)
	require.Len(t, members, 1)
	require.Equal(t, "12345", members[0].Resource)

	rules, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	require.NoError(t, err)
	draftRules := make(map[string]bool)
	for _, r := range rules {
		if r.RoleID == draftRole.ID {
			draftRules[r.Resource] = true
		}
	}
	require.True(t, draftRules[fmt.Sprintf("corteza::compose:module/%d/%d", draftNs.ID, draftOrders.ID)],
		"a rule on a copied module must follow the copy")
	require.True(t, draftRules["corteza::system:user/*"],
		"a rule naming nothing the revision owns must survive unchanged")
}

// TestSwapChatbotWidgetKeys is the promise a publish makes to every embed
// already pasted into a customer's site: the key keeps working.
//
// A widget key is globally unique, so the revision copy could not be created
// holding it and took a throwaway one. Publishing hands the canonical key over —
// retiring the outgoing chatbot's copy of it FIRST, since otherwise the two
// would collide on the unique index mid-swap.
func TestSwapChatbotWidgetKeys(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-widget-keys")

	live := &types.Chatbot{
		ID:        nextID(),
		ProjectID: parent.ID,
		Handle:    "support",
		WidgetKey: "canonical-embed-key",
		CreatedAt: time.Now(),
	}
	require.NoError(t, store.CreateChatbot(ctx, s, live))

	rev, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	draftBots, _, err := store.SearchChatbots(ctx, s, types.ChatbotFilter{ProjectID: rev.ID})
	require.NoError(t, err)
	require.Len(t, draftBots, 1)
	require.NotEqual(t, live.WidgetKey, draftBots[0].WidgetKey)

	require.NoError(t, svc.swapChatbotWidgetKeys(ctx, s, parent, rev))

	promoted, err := store.LookupChatbotByID(ctx, s, draftBots[0].ID)
	require.NoError(t, err)
	require.Equal(t, "canonical-embed-key", promoted.WidgetKey,
		"the embed the customer pasted must resolve to the revision going live")

	retired, err := store.LookupChatbotByID(ctx, s, live.ID)
	require.NoError(t, err)
	require.NotEqual(t, "canonical-embed-key", retired.WidgetKey)
	require.Contains(t, retired.WidgetKey, "canonical-embed-key-deprecated-")

	// The key really did move: looking it up finds exactly the new chatbot.
	byKey, err := store.LookupChatbotByWidgetKey(ctx, s, "canonical-embed-key")
	require.NoError(t, err)
	require.Equal(t, promoted.ID, byKey.ID)
}

// TestDiffSources_ProjectRolesMatchOnTheirSuffix pins the publish plan against a
// false alarm. A project role's handle embeds the revision id, so a copy can
// never keep the original — and comparing raw handles reported every role as
// removed and re-added on a branch that had not touched one of them, which is
// exactly the kind of noise that trains a reviewer to skim the plan.
func TestDiffSources_ProjectRolesMatchOnTheirSuffix(t *testing.T) {
	old := []*GraphSource{{Kind: "role", Handle: "proj_100_reviewer", Name: "Reviewer"}}
	new := []*GraphSource{{Kind: "role", Handle: "proj_200_reviewer", Name: "Reviewer"}}

	require.Empty(t, diffSources(old, new), "the same role under two revisions is not a change")

	// A role that really was added still reads as one.
	added := []*GraphSource{
		{Kind: "role", Handle: "proj_200_reviewer", Name: "Reviewer"},
		{Kind: "role", Handle: "proj_200_approver", Name: "Approver"},
	}
	changes := diffSources(old, added)
	require.Len(t, changes, 1)
	require.Equal(t, types.ProjectChangeOpAdded, changes[0].Op)
	require.Equal(t, "Approver", changes[0].Name)
}
