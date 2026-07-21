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

func TestProject_Publish_GatedRequiresApproval(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "gated-unapproved",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeGated,
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

func TestProject_Publish_GatedSubmittedNotApproved(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "gated-submitted",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeGated,
		Governance: types.ProjectGovernance{
			types.ProjectGovernanceStepPublish: &types.ProjectGovernanceStep{Status: types.ProjectGovernanceStatusSubmitted},
		},
	})

	_, err := svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.Error(t, err)
}

func TestProject_Publish_GatedApprovedSucceedsAndResetsStep(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "gated-approved",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeGated,
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

func TestProject_Publish_FreeModeUnaffectedByGovernance(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid := nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "free-mode",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeFree,
		// no governance state at all
	})

	p, err := svc.Publish(ctx, pid, types.PublishRequest{Confirm: true})
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, types.ProjectStatusActive, p.Status)
	require.Empty(t, p.Governance)
}

// TestProject_PublishGovernance_FullApprovalCycle exercises the whole
// governance state machine end to end through the "publish" step key:
// submit → request-changes → resubmit → approve → publish (which resets the
// step) → submit → approve → publish again, proving every publish needs its
// own fresh submit → approve cycle.
func TestProject_PublishGovernance_FullApprovalCycle(t *testing.T) {
	ctx := context.Background()
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "full-cycle",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeGated,
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

	// NOTE: the project is active now, so a real second publish also needs
	// draft status again — that's the CreateRevision flow's job, out of scope
	// here. This test only asserts the governance half of the cycle: fresh
	// submit → approve is required again after a reset.
}

// seedGatedProjectWithGranter creates a gated-mode project plus a member
// with CanGrantApproval (and CanRequestApproval), returning the project ID
// and a context carrying that member's identity — the shared fixture for
// the flag-anytime tests below.
func seedGatedProjectWithGranter(t *testing.T, s store.Storer, handle string, governance types.ProjectGovernance) (context.Context, uint64) {
	t.Helper()
	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:         pid,
		Handle:     handle,
		Status:     types.ProjectStatusDraft,
		Mode:       types.ProjectModeGated,
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

// TestProject_TransitionGovernanceStep_FlagDraftStepAnytime covers product
// decision 1: a granter may request changes on a step that has no
// governance entry yet, and at any time (not just "submitted") — the entry
// is auto-created straight into changes-requested.
func TestProject_TransitionGovernanceStep_FlagDraftStepAnytime(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedGatedProjectWithGranter(t, s, "flag-draft", nil)

	p, err := svc.TransitionGovernanceStep(ctx, pid, "data-sensitivity", types.ProjectGovernanceActionRequestChanges, "please redact the SSN column")
	require.NoError(t, err)

	step := p.Governance["data-sensitivity"]
	require.NotNil(t, step)
	require.Equal(t, types.ProjectGovernanceStatusChangesRequested, step.Status)
	require.Equal(t, "please redact the SSN column", step.ReviewNote)
}

// TestProject_TransitionGovernanceStep_FlagSendsBackSubmittedPublish covers
// product decisions 1+2: flagging a non-publish step while "publish" is
// submitted sends "publish" back to changes-requested too, atomically with
// the flag.
func TestProject_TransitionGovernanceStep_FlagSendsBackSubmittedPublish(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedGatedProjectWithGranter(t, s, "flag-sends-back-submitted", types.ProjectGovernance{
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

// TestProject_TransitionGovernanceStep_FlagRevokesApprovedPublish covers
// product decision 2's other trigger status: flagging a non-publish step
// while "publish" is already approved revokes that approval.
func TestProject_TransitionGovernanceStep_FlagRevokesApprovedPublish(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedGatedProjectWithGranter(t, s, "flag-revokes-approved", types.ProjectGovernance{
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

// TestProject_TransitionGovernanceStep_FlagNoSideEffectWhenPublishDraft
// covers the negative case of product decision 2: a draft (or absent)
// "publish" step has nothing in flight, so flagging another step must not
// touch it at all.
func TestProject_TransitionGovernanceStep_FlagNoSideEffectWhenPublishDraft(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedGatedProjectWithGranter(t, s, "flag-no-side-effect", nil)

	p, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionRequestChanges, "needs work")
	require.NoError(t, err)

	// no publish entry was even created
	_, exists := p.Governance[types.ProjectGovernanceStepPublish]
	require.False(t, exists)
}

// TestProject_TransitionGovernanceStep_ResubmitClearsFlaggedSteps covers
// product decision 3: submitting "publish" resets every other
// changes-requested step back to draft with its note cleared, but leaves
// draft/other-status steps alone.
func TestProject_TransitionGovernanceStep_ResubmitClearsFlaggedSteps(t *testing.T) {
	svc, s := newTestProjectSvc(t)
	ctx, pid := seedGatedProjectWithGranter(t, s, "resubmit-clears-flags", types.ProjectGovernance{
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

// TestProject_TransitionGovernanceStep_FlagRejectedInFreeMode covers product
// decision 4's free-mode carve-out: per-step flagging has no meaning without
// approval concepts, so it's rejected outright for free-mode projects.
func TestProject_TransitionGovernanceStep_FlagRejectedInFreeMode(t *testing.T) {
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "free-mode-flag",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeFree,
	})
	seedMember(t, s, &types.ProjectMember{
		ID:         nextID(),
		ProjectID:  pid,
		UserID:     uid,
		RolePreset: types.ProjectRoleGovernanceOwner,
	})
	ctx := a.SetIdentityToContext(context.Background(), &types.User{ID: uid})

	_, err := svc.TransitionGovernanceStep(ctx, pid, "summary", types.ProjectGovernanceActionRequestChanges, "no-op")
	require.Error(t, err)

	// nothing was persisted
	reloaded, err := store.LookupProjectByID(ctx, s, pid)
	require.NoError(t, err)
	require.Empty(t, reloaded.Governance)
}

// TestProject_TransitionGovernanceStep_FlagRejectedWithoutGrantCapability
// covers the capability requirement for flag-anytime: it needs
// CanGrantApproval just like the existing request-changes action, so a
// member with only CanRequestApproval (Developer) is rejected.
func TestProject_TransitionGovernanceStep_FlagRejectedWithoutGrantCapability(t *testing.T) {
	svc, s := newTestProjectSvc(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{
		ID:     pid,
		Handle: "flag-needs-grant",
		Status: types.ProjectStatusDraft,
		Mode:   types.ProjectModeGated,
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
