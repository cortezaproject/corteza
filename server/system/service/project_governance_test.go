package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

// newTestProjectSvc builds a bare project service (no actionlog, no access
// controller) backed by an in-memory sqlite store — sufficient for Publish
// and governance transitions, neither of which touch svc.ac or svc.actionlog
// on the success/expected-failure paths exercised below.
func newTestProjectSvc(t *testing.T) (*project, store.Storer) {
	t.Helper()
	req := require.New(t)
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	return &project{store: s}, s
}

// --- publish gate (now unconditional for every project) ---

// TestProject_Publish_RequiresApproval: publishing is blocked for every
// project until the "publish" governance step is approved — a project with no
// governance state at all has an implicit draft publish step, so it can't
// publish.
func TestProject_Publish_RequiresApproval(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "unapproved",
		Status: types.ProjectStatusDraft,
		// no governance state at all
	})

	p, err := svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
	require.Nil(t, p)
	require.ErrorContains(t, err, "publish")
	require.ErrorContains(t, err, "approved")

	// project must remain untouched (still draft)
	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Equal(t, types.ProjectStatusDraft, reloaded.Status)
}

// TestProject_Publish_SubmittedNotApproved: a publish step that is only
// submitted (not yet approved) still blocks publishing.
func TestProject_Publish_SubmittedNotApproved(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "submitted",
		Status: types.ProjectStatusDraft,
		Governance: types.ProjectGovernance{
			types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusSubmitted},
		},
	})

	_, err := svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
}

// TestProject_Publish_ApprovedSucceedsAndResetsStep: with the publish step
// approved, publishing succeeds and resets the step back to draft so the next
// publish requires a fresh approval cycle.
func TestProject_Publish_ApprovedSucceedsAndResetsStep(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "approved",
		Status: types.ProjectStatusDraft,
		Governance: types.ProjectGovernance{
			types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{
				Status:     types.ProjectGovernanceStatusApproved,
				ReviewNote: "looks good",
			},
		},
	})

	p, err := svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, types.ProjectStatusActive, p.Status)

	// the publish step must reset to draft (fresh approval required next time)
	step := p.Governance[types.ProjectGovernanceStepPublish]
	require.NotNil(t, step)
	require.Equal(t, types.ProjectGovernanceStatusDraft, step.Status)
	require.Empty(t, step.ReviewNote)

	// ...and that reset must be persisted, not just held on the in-memory return value
	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Equal(t, types.ProjectStatusActive, reloaded.Status)
	reloadedStep := reloaded.Governance[types.ProjectGovernanceStepPublish]
	require.NotNil(t, reloadedStep)
	require.Equal(t, types.ProjectGovernanceStatusDraft, reloadedStep.Status)
	require.Empty(t, reloadedStep.ReviewNote)

	// republishing immediately (without a fresh approval) must be blocked again
	_, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
}

// TestProject_PublishGovernance_FullApprovalCycle exercises the whole
// project-level approval state machine end to end through the "publish" step
// key: submit → request-changes → resubmit → approve → publish (which resets
// the step) → publish blocked again, proving every publish needs its own fresh
// submit → approve cycle.
func TestProject_PublishGovernance_FullApprovalCycle(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "full-cycle",
		Status: types.ProjectStatusDraft,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleGovernanceOwner, // CanRequestApproval + CanGrantApproval
	})
	ctx = a.SetIdentityToContext(ctx, &types.User{ID: uid})

	stepKey := types.ProjectGovernanceStepPublish

	// draft → submit → submitted
	p, err := svc.TransitionGovernanceStep(ctx, pid, stepKey, types.ProjectGovernanceActionSubmit, "")
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusSubmitted, p.Governance[stepKey].Status)

	// publish still blocked while only submitted
	_, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)

	// submitted → request-changes → changes-requested
	p, err = svc.TransitionGovernanceStep(ctx, pid, stepKey, types.ProjectGovernanceActionRequestChanges, "needs more detail")
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, p.Governance[stepKey].Status)
	require.Equal(t, "needs more detail", p.Governance[stepKey].ReviewNote)

	// changes-requested → submit → submitted (note cleared)
	p, err = svc.TransitionGovernanceStep(ctx, pid, stepKey, types.ProjectGovernanceActionSubmit, "")
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusSubmitted, p.Governance[stepKey].Status)
	require.Empty(t, p.Governance[stepKey].ReviewNote)

	// submitted → approve → approved
	p, err = svc.TransitionGovernanceStep(ctx, pid, stepKey, types.ProjectGovernanceActionApprove, "")
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusApproved, p.Governance[stepKey].Status)

	// first publish now succeeds, and resets the step
	p, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.NoError(t, err)
	require.Equal(t, types.ProjectStatusActive, p.Status)
	require.Equal(t, types.ProjectGovernanceStatusDraft, p.Governance[stepKey].Status)

	// republishing (there's no new draft revision, so Publish still targets
	// the same, now-active project) must go through submit → approve again;
	// straight publish is blocked with a fresh draft step
	_, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
}

// --- direct per-step review (Build / Govern steps) ---

// seedProjectWithGranter creates a project plus a member with CanGrantApproval
// (and CanRequestApproval), returning the project ID and a context carrying
// that member's identity — the shared fixture for the direct-review tests.
func seedProjectWithGranter(t *testing.T, s store.Storer, handle string, governance types.ProjectGovernance) (context.Context, uint64) {
	t.Helper()
	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:         pid,
		Handle:     handle,
		Status:     types.ProjectStatusDraft,
		Governance: governance,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleGovernanceOwner, // CanRequestApproval + CanGrantApproval
	})
	return a.SetIdentityToContext(context.Background(), &types.User{ID: uid}), pid
}

// TestProject_TransitionGovernanceStep_ApproveStepFromDraft: a granter may
// approve a Build/Govern step directly, even one that has no governance entry
// yet — the entry is auto-created straight into approved, no submit stage.
func TestProject_TransitionGovernanceStep_ApproveStepFromDraft(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "approve-draft", nil)

	p, err := svc.TransitionGovernanceStep(ctx, pid, "data-model", types.ProjectGovernanceActionApprove, "")
	require.NoError(t, err)

	step := p.Governance["data-model"]
	require.NotNil(t, step)
	require.Equal(t, types.ProjectGovernanceStatusApproved, step.Status)
	require.Empty(t, step.ReviewNote)
}

// TestProject_TransitionGovernanceStep_ApproveClearsChangesRequested:
// approving a step that currently has changes requested clears the flag (note
// and status).
func TestProject_TransitionGovernanceStep_ApproveClearsChangesRequested(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "approve-clears-flag", types.ProjectGovernance{
		"data-sensitivity": &types.ProjectGovernanceStep{
			Status:     types.ProjectGovernanceStatusChangesRequested,
			ReviewNote: "redact the SSN column",
		},
	})

	p, err := svc.TransitionGovernanceStep(ctx, pid, "data-sensitivity", types.ProjectGovernanceActionApprove, "")
	require.NoError(t, err)

	step := p.Governance["data-sensitivity"]
	require.NotNil(t, step)
	require.Equal(t, types.ProjectGovernanceStatusApproved, step.Status)
	require.Empty(t, step.ReviewNote)
}

// TestProject_TransitionGovernanceStep_FlagStepAnytime: a granter may request
// changes on a step that has no governance entry yet, and at any time — the
// entry is auto-created straight into changes-requested.
func TestProject_TransitionGovernanceStep_FlagStepAnytime(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "flag-draft", nil)

	p, err := svc.TransitionGovernanceStep(ctx, pid, "data-sensitivity", types.ProjectGovernanceActionRequestChanges, "please redact the SSN column")
	require.NoError(t, err)

	step := p.Governance["data-sensitivity"]
	require.NotNil(t, step)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, step.Status)
	require.Equal(t, "please redact the SSN column", step.ReviewNote)
}

// TestProject_TransitionGovernanceStep_FlagSendsBackSubmittedPublish: flagging
// a non-publish step while "publish" is submitted sends "publish" back to
// changes-requested too, atomically with the flag.
func TestProject_TransitionGovernanceStep_FlagSendsBackSubmittedPublish(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "flag-sends-back-submitted", types.ProjectGovernance{
		types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusSubmitted},
	})

	p, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionRequestChanges, "typo in the title")
	require.NoError(t, err)

	flagged := p.Governance["summary"]
	require.NotNil(t, flagged)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, flagged.Status)
	require.Equal(t, "typo in the title", flagged.ReviewNote)

	publish := p.Governance[types.ProjectGovernanceStepPublish]
	require.NotNil(t, publish)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, publish.Status)
	require.Contains(t, publish.ReviewNote, "summary")
	require.Contains(t, publish.ReviewNote, "typo in the title")

	// persisted, not just held on the in-memory return value
	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, reloaded.Governance[types.ProjectGovernanceStepPublish].Status)
}

// TestProject_TransitionGovernanceStep_FlagRevokesApprovedPublish: flagging a
// non-publish step while "publish" is already approved revokes that approval,
// and publishing is blocked again as a result.
func TestProject_TransitionGovernanceStep_FlagRevokesApprovedPublish(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "flag-revokes-approved", types.ProjectGovernance{
		types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusApproved},
	})

	p, err := svc.TransitionGovernanceStep(ctx, pid, "data-sensitivity", types.ProjectGovernanceActionRequestChanges, "missing DPIA")
	require.NoError(t, err)

	publish := p.Governance[types.ProjectGovernanceStepPublish]
	require.NotNil(t, publish)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, publish.Status)
	require.Contains(t, publish.ReviewNote, "data-sensitivity")

	// publish must now be blocked again despite having been approved a moment ago
	_, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
}

// TestProject_TransitionGovernanceStep_FlagNoSideEffectWhenPublishDraft: a
// draft (or absent) "publish" step has nothing in flight, so flagging another
// step must not touch it at all.
func TestProject_TransitionGovernanceStep_FlagNoSideEffectWhenPublishDraft(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "flag-no-side-effect", nil)

	p, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionRequestChanges, "needs work")
	require.NoError(t, err)

	// no publish entry was even created
	_, exists := p.Governance[types.ProjectGovernanceStepPublish]
	require.False(t, exists)
}

// TestProject_TransitionGovernanceStep_ResubmitClearsFlaggedSteps: submitting
// "publish" resets every other changes-requested step back to draft with its
// note cleared, but leaves draft/other-status steps alone.
func TestProject_TransitionGovernanceStep_ResubmitClearsFlaggedSteps(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "resubmit-clears-flags", types.ProjectGovernance{
		types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusDraft},
		"summary":                          &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusChangesRequested, ReviewNote: "flagged 1"},
		"data-sensitivity":                 &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusChangesRequested, ReviewNote: "flagged 2"},
		"untouched-draft":                  &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusDraft},
	})

	p, err := svc.TransitionGovernanceStep(ctx, pid, types.ProjectGovernanceStepPublish, types.ProjectGovernanceActionSubmit, "")
	require.NoError(t, err)

	require.Equal(t, types.ProjectGovernanceStatusSubmitted, p.Governance[types.ProjectGovernanceStepPublish].Status)

	require.Equal(t, types.ProjectGovernanceStatusDraft, p.Governance["summary"].Status)
	require.Empty(t, p.Governance["summary"].ReviewNote)

	require.Equal(t, types.ProjectGovernanceStatusDraft, p.Governance["data-sensitivity"].Status)
	require.Empty(t, p.Governance["data-sensitivity"].ReviewNote)

	require.Equal(t, types.ProjectGovernanceStatusDraft, p.Governance["untouched-draft"].Status)
}

// TestProject_TransitionGovernanceStep_ApproveProjectBlockedWhileStepFlagged:
// approving the project (approve on the publish step) is rejected while any
// other step still has changes requested; once that step is approved, the
// project approval goes through.
func TestProject_TransitionGovernanceStep_ApproveProjectBlockedWhileStepFlagged(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedProjectWithGranter(t, s, "approve-blocked-by-flag", types.ProjectGovernance{
		types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusSubmitted},
		"data-sensitivity":                 &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusChangesRequested, ReviewNote: "fix this"},
	})

	// project approval blocked while data-sensitivity is flagged
	_, err := svc.TransitionGovernanceStep(ctx, pid, types.ProjectGovernanceStepPublish, types.ProjectGovernanceActionApprove, "")
	require.Error(t, err)

	// publish step must still be submitted (unchanged), and nothing published
	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusSubmitted, reloaded.Governance[types.ProjectGovernanceStepPublish].Status)

	// resolve the flag by approving that step
	_, err = svc.TransitionGovernanceStep(ctx, pid, "data-sensitivity", types.ProjectGovernanceActionApprove, "")
	require.NoError(t, err)

	// now the project can be approved
	p, err := svc.TransitionGovernanceStep(ctx, pid, types.ProjectGovernanceStepPublish, types.ProjectGovernanceActionApprove, "")
	require.NoError(t, err)
	require.Equal(t, types.ProjectGovernanceStatusApproved, p.Governance[types.ProjectGovernanceStepPublish].Status)

	// and publishing now succeeds
	p, err = svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.NoError(t, err)
	require.Equal(t, types.ProjectStatusActive, p.Status)
}

// --- capability rejections ---

// TestProject_TransitionGovernanceStep_FlagRejectedWithoutGrantCapability:
// requesting changes needs CanGrantApproval, so a member with only
// CanRequestApproval (Developer) is rejected and nothing is persisted.
func TestProject_TransitionGovernanceStep_FlagRejectedWithoutGrantCapability(t *testing.T) {
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "flag-needs-grant",
		Status: types.ProjectStatusDraft,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleDeveloper, // CanRequestApproval only, no CanGrantApproval
	})
	ctx := a.SetIdentityToContext(context.Background(), &types.User{ID: uid})

	_, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionRequestChanges, "no-op")
	require.Error(t, err)

	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Empty(t, reloaded.Governance)
}

// TestProject_TransitionGovernanceStep_ApproveRejectedWithoutGrantCapability:
// approving a step needs CanGrantApproval, so a Developer is rejected and
// nothing is persisted.
func TestProject_TransitionGovernanceStep_ApproveRejectedWithoutGrantCapability(t *testing.T) {
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "approve-needs-grant",
		Status: types.ProjectStatusDraft,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleDeveloper, // CanRequestApproval only, no CanGrantApproval
	})
	ctx := a.SetIdentityToContext(context.Background(), &types.User{ID: uid})

	_, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionApprove, "")
	require.Error(t, err)

	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Empty(t, reloaded.Governance)
}

// TestProject_TransitionGovernanceStep_SubmitRejectedWithoutRequestCapability:
// submitting the publish step needs CanRequestApproval, so a member with
// neither approval capability (JuniorDeveloper) is rejected.
func TestProject_TransitionGovernanceStep_SubmitRejectedWithoutRequestCapability(t *testing.T) {
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "submit-needs-request",
		Status: types.ProjectStatusDraft,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleJuniorDeveloper, // CanWrite only, no request/grant
	})
	ctx := a.SetIdentityToContext(context.Background(), &types.User{ID: uid})

	_, err := svc.TransitionGovernanceStep(ctx, pid, types.ProjectGovernanceStepPublish, types.ProjectGovernanceActionSubmit, "")
	require.Error(t, err)

	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Empty(t, reloaded.Governance)
}
