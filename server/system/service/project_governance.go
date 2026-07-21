package service

import (
	"context"
	"fmt"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// Governance operations on a project's steps. Capability checks run against
// the caller's project membership (role preset), not system RBAC:
//
//	save step values          → CanWrite, step must be editable (draft/changes-requested)
//	submit (publish step)     → CanRequestApproval
//	approve / request-changes → CanGrantApproval
//
// TransitionGovernanceStep handles two shapes of step:
//
//   - Ordinary Build/Govern steps: reviewed directly, with no submit stage.
//     A granter approves the step (types.Step.Approve) or requests changes on
//     it (types.Step.Flag) in any status, at any time, including a step with
//     no entry yet. Approving clears a prior changes-requested flag.
//   - The "publish" step (types.ProjectGovernanceStepPublish): the
//     project-level approval state machine (submit → approve /
//     request-changes, types.Step.Transition). project.Publish() requires it
//     to be approved before it will run and resets it back to draft on
//     success — see service/project_revision.go.
//
// On top of that, TransitionGovernanceStep implements cross-step
// orchestration:
//
//   - "immediate send-back": flagging a non-publish step while "publish" is
//     currently submitted or approved sends "publish" back to
//     changes-requested too, since whatever was about to ship now has open
//     feedback against it. A draft/absent "publish" step is left alone —
//     there's nothing in flight to revoke.
//   - "auto-clear on resubmit": submitting "publish" resets every other step
//     currently flagged (changes-requested) back to draft with its note
//     cleared, since resubmission asserts the feedback was addressed.
//   - "approval gate": approving the "publish" step (approving the project)
//     is rejected while any other step still has changes requested — each
//     flag must be resolved by approving that step first.

func (svc *project) SaveGovernanceStep(ctx context.Context, projectID uint64, stepKey string, values map[string]any) (p *types.Project, err error) {
	var (
		paProps = &projectActionProps{project: &types.Project{ID: projectID}}
		old     *types.Project
	)

	err = func() (err error) {
		if p, err = loadProject(ctx, svc.store, projectID); err != nil {
			return
		}
		old = p.Clone()
		paProps.setProject(p)

		caps, err := svc.memberCapabilities(ctx, p)
		if err != nil {
			return err
		}
		if !caps.CanWrite {
			return ProjectErrNotAllowedToEditGovernance()
		}

		step := p.Governance.Step(stepKey)
		if !step.Editable() {
			return ProjectErrGovernanceStepLocked()
		}

		step.Values = values
		if err = svc.storeGovernance(ctx, p); err != nil {
			return
		}
		return
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionGovernanceSave, err, old, p)
}

func (svc *project) TransitionGovernanceStep(ctx context.Context, projectID uint64, stepKey string, action types.ProjectGovernanceAction, note string) (p *types.Project, err error) {
	var (
		paProps = &projectActionProps{project: &types.Project{ID: projectID}}
		old     *types.Project
	)

	err = func() (err error) {
		if p, err = loadProject(ctx, svc.store, projectID); err != nil {
			return
		}
		old = p.Clone()
		paProps.setProject(p)

		caps, err := svc.memberCapabilities(ctx, p)
		if err != nil {
			return err
		}
		if action.RequiresGrant() {
			if !caps.CanGrantApproval {
				return ProjectErrNotAllowedToEditGovernance()
			}
		} else if !caps.CanRequestApproval {
			return ProjectErrNotAllowedToEditGovernance()
		}

		if stepKey == types.ProjectGovernanceStepPublish {
			// Publish step: the project-level submit → approve /
			// request-changes machine. Approving the project (approve on the
			// publish step) is blocked while any other step still has changes
			// requested — each flag must be resolved (approved) first.
			if action == types.ProjectGovernanceActionApprove && hasFlaggedSteps(p) {
				return ProjectErrApprovalBlockedByChangesRequested()
			}
			if !p.Governance.Step(stepKey).Transition(action, note) {
				return ProjectErrInvalidGovernanceTransition()
			}
			if action == types.ProjectGovernanceActionSubmit {
				clearFlaggedSteps(p)
			}
		} else {
			// Ordinary Build/Govern step: direct review, no submit stage. A
			// granter either approves the step or requests changes on it, in
			// any status, at any time.
			switch action {
			case types.ProjectGovernanceActionApprove:
				p.Governance.Step(stepKey).Approve()
			case types.ProjectGovernanceActionRequestChanges:
				p.Governance.Step(stepKey).Flag(note)
				sendPublishBack(p, stepKey, note)
			default:
				return ProjectErrInvalidGovernanceTransition()
			}
		}

		if err = svc.storeGovernance(ctx, p); err != nil {
			return
		}
		return
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionGovernanceTransition, err, old, p)
}

// sendPublishBack is the "immediate send-back" side effect of flagging a
// non-publish step (see TransitionGovernanceStep): if "publish" is currently
// submitted or approved, its pending/granted state no longer reflects a
// project without open feedback, so it goes back to changes-requested too,
// with a note pointing at the step that was flagged. A draft or altogether
// absent "publish" step is left untouched — there's nothing in flight to
// revoke, so the map is not even auto-vivified for it.
func sendPublishBack(p *types.Project, flaggedStepKey, note string) {
	publish := p.Governance[types.ProjectGovernanceStepPublish]
	if publish == nil {
		return
	}
	if publish.Status != types.ProjectGovernanceStatusSubmitted && publish.Status != types.ProjectGovernanceStatusApproved {
		return
	}
	publish.Status = types.ProjectGovernanceStatusChangesRequested
	publish.ReviewNote = fmt.Sprintf("changes requested on step %q: %s", flaggedStepKey, note)
}

// clearFlaggedSteps is the "auto-clear on resubmit" side effect (see
// TransitionGovernanceStep): fired when "publish" is (re)submitted, it
// resets every other step currently flagged (changes-requested) back to
// draft with its note cleared, since resubmission asserts the feedback was
// addressed.
func clearFlaggedSteps(p *types.Project) {
	for key, step := range p.Governance {
		if key == types.ProjectGovernanceStepPublish {
			continue
		}
		if step.Status == types.ProjectGovernanceStatusChangesRequested {
			step.Reset()
		}
	}
}

// hasFlaggedSteps reports whether any step other than "publish" currently has
// changes requested. It backs the "approval gate": the project can't be
// approved (publish step approved) while an open change request remains on any
// step — that request must first be resolved by approving the flagged step.
func hasFlaggedSteps(p *types.Project) bool {
	for key, step := range p.Governance {
		if key == types.ProjectGovernanceStepPublish {
			continue
		}
		if step.Status == types.ProjectGovernanceStatusChangesRequested {
			return true
		}
	}
	return false
}

func (svc *project) storeGovernance(ctx context.Context, p *types.Project) error {
	p.UpdatedAt = now()
	p.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return store.UpdateProject(ctx, svc.store, p)
}

// memberCapabilities resolves the caller's capabilities on the project: an
// explicit membership wins; otherwise open-visibility projects fall back to
// the default member role. Same access model as the project scope resolver.
func (svc *project) memberCapabilities(ctx context.Context, p *types.Project) (types.ProjectCapabilities, error) {
	var (
		zero   types.ProjectCapabilities
		userID = a.GetIdentityFromContext(ctx).Identity()
	)

	m, err := store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, p.ID, userID)
	if err == nil && m != nil && m.DeletedAt == nil {
		return m.RolePreset.Capabilities(), nil
	}
	if err != nil && !errors.IsNotFound(err) {
		return zero, err
	}

	if p.Config.Visibility == types.ProjectVisibilityOpen {
		role := p.Config.DefaultMemberRole
		if !role.Valid() {
			role = types.ProjectRoleMember
		}
		return role.Capabilities(), nil
	}

	return zero, ProjectErrNotAllowedToEditGovernance()
}
