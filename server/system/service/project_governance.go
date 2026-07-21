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
//	save step values  → CanWrite, step must be editable (draft/changes-requested)
//	submit / recall   → CanRequestApproval
//	approve / request-changes / reopen → CanGrantApproval
//
// Per-section approval gates are gone: every step but one only ever uses
// SaveGovernanceStep to persist form values and stays in the draft status,
// so it never locks. The exception is the "publish" step key
// (types.ProjectGovernanceStepPublish): in gated mode, project.Publish()
// requires it to be approved before it will run, and resets it back to
// draft on success — see service/project_revision.go.
//
// On top of that, TransitionGovernanceStep implements one more piece of
// cross-step orchestration, all gated-mode only:
//
//   - "flag anytime": a request-changes on any step other than "publish"
//     bypasses that step's own submitted-only rule (types.Step.Flag instead
//     of types.Step.Transition) — a granter can raise it on any step, in any
//     status, at any time, including a step with no entry yet.
//   - "immediate send-back": flagging a non-publish step while "publish" is
//     currently submitted or approved sends "publish" back to
//     changes-requested too, since whatever was about to ship now has open
//     feedback against it. A draft/absent "publish" step is left alone —
//     there's nothing in flight to revoke.
//   - "auto-clear on resubmit": submitting "publish" resets every other step
//     currently flagged (changes-requested) back to draft with its note
//     cleared, since resubmission asserts the feedback was addressed.
//
// Free-mode projects have no approval concepts at all, so the flag-anytime
// path is rejected outright for them (same types.ProjectModeGated check
// project.Publish() uses to decide whether the publish gate applies).

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

		if action == types.ProjectGovernanceActionRequestChanges && stepKey != types.ProjectGovernanceStepPublish {
			// Flag-anytime: a granter may request changes on any non-publish
			// step regardless of its current status, but only in gated mode
			// — free-mode projects have no approval concepts to flag.
			if p.Mode != types.ProjectModeGated {
				return ProjectErrInvalidGovernanceTransition()
			}
			p.Governance.Step(stepKey).Flag(note)
			sendPublishBack(p, stepKey, note)
		} else if !p.Governance.Step(stepKey).Transition(action, note) {
			return ProjectErrInvalidGovernanceTransition()
		}

		if action == types.ProjectGovernanceActionSubmit && stepKey == types.ProjectGovernanceStepPublish {
			clearFlaggedSteps(p)
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
