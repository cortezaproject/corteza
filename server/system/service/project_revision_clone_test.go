package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// seedCloneModule adds a module to a namespace. Store-level on purpose: the
// compose service would drag in guards and DAL wiring the copy has no opinion
// about, and the copy only ever sees what is in the store.
func seedCloneModule(t *testing.T, s store.Storer, nsID uint64, handle, name string) *composeTypes.Module {
	t.Helper()

	m := &composeTypes.Module{
		ID:          nextID(),
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

	customers := seedCloneModule(t, s, parentNs.ID, "customers", "Customers")
	invoices := seedCloneModule(t, s, parentNs.ID, "invoices", "Invoices")

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

	// Parent-revision refs with no copy: dropped, shared ones kept.
	require.Equal(t, []types.AgentAccessTAQ{{ID: sharedTAQ.ID}}, cp.Access.TAQs)
	require.Equal(t, []types.AgentAccessWorkflow{{ID: sharedWf.ID}}, cp.Access.Workflows)

	require.NoError(t, label.Load(ctx, s, cp))
	require.Equal(t, map[string]labelTypes.LabelValue{"tier": {Val: "gold"}}, cp.Labels)
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

	customers := seedCloneModule(t, s, parentNs.ID, "source-customers", "Customers")
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

	unmappable := seedCloneModule(t, s, parentNs.ID, "", "No handle")

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
