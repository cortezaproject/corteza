package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Project struct {
		svc      projectService
		graphSvc projectGraphService
		ac       projectAccessController
	}

	projectPayload struct {
		*types.Project

		CanGrant         bool `json:"canGrant"`
		CanUpdateProject bool `json:"canUpdateProject"`
		CanDeleteProject bool `json:"canDeleteProject"`
		CanManageMembers bool `json:"canManageMembers"`
	}

	projectSetPayload struct {
		Filter types.ProjectFilter `json:"filter"`
		Set    []*projectPayload   `json:"set"`
	}

	projectMemberPayload struct {
		*types.ProjectMember

		Capabilities types.ProjectCapabilities `json:"capabilities"`
	}

	projectMemberSetPayload struct {
		Filter types.ProjectMemberFilter `json:"filter"`
		Set    []*projectMemberPayload   `json:"set"`
	}

	projectService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Project, error)
		Create(ctx context.Context, new *types.Project) (*types.Project, error)
		Update(ctx context.Context, upd *types.Project) (*types.Project, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ProjectFilter) (types.ProjectSet, types.ProjectFilter, error)

		SearchMembers(ctx context.Context, filter types.ProjectMemberFilter) (types.ProjectMemberSet, types.ProjectMemberFilter, error)
		AddMember(ctx context.Context, m *types.ProjectMember) (*types.ProjectMember, error)
		UpdateMember(ctx context.Context, m *types.ProjectMember) (*types.ProjectMember, error)
		RemoveMember(ctx context.Context, projectID, userID uint64) error

		SaveGovernanceStep(ctx context.Context, projectID uint64, stepKey string, values map[string]any) (*types.Project, error)
		TransitionGovernanceStep(ctx context.Context, projectID uint64, stepKey string, action types.ProjectGovernanceAction, note string) (*types.Project, error)
	}

	projectGraphService interface {
		Graph(ctx context.Context, projectID uint64) (*types.ProjectGraph, error)
	}

	projectAccessController interface {
		CanGrant(context.Context) bool

		CanCreateProject(context.Context) bool
		CanUpdateProject(context.Context, *types.Project) bool
		CanDeleteProject(context.Context, *types.Project) bool
		CanManageMembersOnProject(context.Context, *types.Project) bool
	}
)

func (Project) New() *Project {
	return &Project{
		svc:      service.DefaultProject,
		graphSvc: service.DefaultProjectGraph,
		ac:       service.DefaultAccessControl,
	}
}

func (ctrl *Project) List(ctx context.Context, r *request.ProjectList) (interface{}, error) {
	var (
		err error
		f   = types.ProjectFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Status:  types.ProjectStatus(r.Status),
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *Project) Create(ctx context.Context, r *request.ProjectCreate) (interface{}, error) {
	p := &types.Project{
		Handle: r.Handle,
		Status: types.ProjectStatus(r.Status),
		Config: r.Config,
		Meta:   r.Meta,
		Labels: r.Labels,
	}

	p, err := ctrl.svc.Create(ctx, p)
	return ctrl.makePayload(ctx, p, err)
}

func (ctrl *Project) Read(ctx context.Context, r *request.ProjectRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ProjectID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Project) Update(ctx context.Context, r *request.ProjectUpdate) (interface{}, error) {
	p := &types.Project{
		ID:        r.ProjectID,
		Handle:    r.Handle,
		Status:    types.ProjectStatus(r.Status),
		Config:    r.Config,
		Meta:      r.Meta,
		UpdatedAt: r.UpdatedAt,
		Labels:    r.Labels,
	}

	p, err := ctrl.svc.Update(ctx, p)
	return ctrl.makePayload(ctx, p, err)
}

func (ctrl *Project) Delete(ctx context.Context, r *request.ProjectDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ProjectID)
}

func (ctrl *Project) Undelete(ctx context.Context, r *request.ProjectUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ProjectID)
}

func (ctrl *Project) ListMembers(ctx context.Context, r *request.ProjectListMembers) (interface{}, error) {
	set, f, err := ctrl.svc.SearchMembers(ctx, types.ProjectMemberFilter{
		ProjectID: r.ProjectID,
	})
	return ctrl.makeMemberFilterPayload(set, f, err)
}

func (ctrl *Project) AddMember(ctx context.Context, r *request.ProjectAddMember) (interface{}, error) {
	m, err := ctrl.svc.AddMember(ctx, &types.ProjectMember{
		ProjectID:  r.ProjectID,
		UserID:     r.UserID,
		RolePreset: types.ProjectMemberRole(r.RolePreset),
	})
	return ctrl.makeMemberPayload(m, err)
}

func (ctrl *Project) UpdateMember(ctx context.Context, r *request.ProjectUpdateMember) (interface{}, error) {
	m, err := ctrl.svc.UpdateMember(ctx, &types.ProjectMember{
		ProjectID:  r.ProjectID,
		UserID:     r.UserID,
		RolePreset: types.ProjectMemberRole(r.RolePreset),
	})
	return ctrl.makeMemberPayload(m, err)
}

func (ctrl *Project) RemoveMember(ctx context.Context, r *request.ProjectRemoveMember) (interface{}, error) {
	return api.OK(), ctrl.svc.RemoveMember(ctx, r.ProjectID, r.UserID)
}

func (ctrl *Project) Graph(ctx context.Context, r *request.ProjectGraph) (interface{}, error) {
	return ctrl.graphSvc.Graph(ctx, r.ProjectID)
}

func (ctrl *Project) GovernanceSave(ctx context.Context, r *request.ProjectGovernanceSave) (interface{}, error) {
	p, err := ctrl.svc.SaveGovernanceStep(ctx, r.ProjectID, r.StepKey, r.Values)
	return ctrl.makePayload(ctx, p, err)
}

func (ctrl *Project) GovernanceTransition(ctx context.Context, r *request.ProjectGovernanceTransition) (interface{}, error) {
	p, err := ctrl.svc.TransitionGovernanceStep(ctx, r.ProjectID, r.StepKey, types.ProjectGovernanceAction(r.Action), r.Note)
	return ctrl.makePayload(ctx, p, err)
}

func (ctrl *Project) makePayload(ctx context.Context, p *types.Project, err error) (*projectPayload, error) {
	if err != nil || p == nil {
		return nil, err
	}

	return &projectPayload{
		Project: p,

		CanGrant:         ctrl.ac.CanGrant(ctx),
		CanUpdateProject: ctrl.ac.CanUpdateProject(ctx, p),
		CanDeleteProject: ctrl.ac.CanDeleteProject(ctx, p),
		CanManageMembers: ctrl.ac.CanManageMembersOnProject(ctx, p),
	}, nil
}

func (ctrl *Project) makeFilterPayload(ctx context.Context, nn types.ProjectSet, f types.ProjectFilter, err error) (*projectSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectSetPayload{Filter: f, Set: make([]*projectPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

func (ctrl *Project) makeMemberPayload(m *types.ProjectMember, err error) (*projectMemberPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &projectMemberPayload{
		ProjectMember: m,
		Capabilities:  m.RolePreset.Capabilities(),
	}, nil
}

func (ctrl *Project) makeMemberFilterPayload(nn types.ProjectMemberSet, f types.ProjectMemberFilter, err error) (*projectMemberSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectMemberSetPayload{Filter: f, Set: make([]*projectMemberPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makeMemberPayload(nn[i], nil)
	}

	return msp, nil
}
