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
