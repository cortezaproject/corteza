package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// Governance operations on the project build pipeline. Capability checks run
// against the caller's project membership (role preset), not system RBAC:
//
//	save step values  → CanWrite, step must be editable (draft/changes-requested)
//	submit / recall   → CanRequestApproval
//	approve / request-changes / reopen → CanGrantApproval

func (svc *project) SaveGovernanceStep(ctx context.Context, projectID uint64, stepKey string, values map[string]any) (p *types.Project, err error) {
	var paProps = &projectActionProps{project: &types.Project{ID: projectID}}

	err = func() (err error) {
		if p, err = loadProject(ctx, svc.store, projectID); err != nil {
			return
		}
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
		return svc.storeGovernance(ctx, p)
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionGovernanceSave, err)
}

func (svc *project) TransitionGovernanceStep(ctx context.Context, projectID uint64, stepKey string, action types.ProjectGovernanceAction, note string) (p *types.Project, err error) {
	var paProps = &projectActionProps{project: &types.Project{ID: projectID}}

	err = func() (err error) {
		if p, err = loadProject(ctx, svc.store, projectID); err != nil {
			return
		}
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

		return svc.storeGovernance(ctx, p)
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionGovernanceTransition, err)
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
