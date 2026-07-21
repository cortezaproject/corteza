package service

import (
	"context"

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

		if !p.Governance.Step(stepKey).Transition(action, note) {
			return ProjectErrInvalidGovernanceTransition()
		}

		if err = svc.storeGovernance(ctx, p); err != nil {
			return
		}
		return
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionGovernanceTransition, err, old, p)
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
