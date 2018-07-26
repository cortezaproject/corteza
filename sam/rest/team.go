package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Team struct {
		svc teamService
	}

	teamService interface {
		FindByID(ctx context.Context, teamID uint64) (*types.Team, error)
		Find(ctx context.Context, filter *types.TeamFilter) ([]*types.Team, error)

		Create(ctx context.Context, team *types.Team) (*types.Team, error)
		Update(ctx context.Context, team *types.Team) (*types.Team, error)
		Merge(ctx context.Context, teamID, targetTeamID uint64) error
		Move(ctx context.Context, teamID, organisationID uint64) error

		deleter
		archiver
	}
)

func (Team) New(teamSvc teamService) *Team {
	var ctrl = &Team{}
	ctrl.svc = teamSvc
	return ctrl
}

func (ctrl *Team) Read(ctx context.Context, r *server.TeamReadRequest) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.TeamID)
}

func (ctrl *Team) List(ctx context.Context, r *server.TeamListRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.TeamFilter{Query: r.Query})
}

func (ctrl *Team) Create(ctx context.Context, r *server.TeamCreateRequest) (interface{}, error) {
	org := types.Team{}.
		New().
		SetName(r.Name)

	return ctrl.svc.Create(ctx, org)
}

func (ctrl *Team) Edit(ctx context.Context, r *server.TeamEditRequest) (interface{}, error) {
	org := types.Team{}.
		New().
		SetID(r.TeamID).
		SetName(r.Name)

	return ctrl.svc.Update(ctx, org)
}

func (ctrl *Team) Remove(ctx context.Context, r *server.TeamRemoveRequest) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.TeamID)
}

func (ctrl *Team) Archive(ctx context.Context, r *server.TeamArchiveRequest) (interface{}, error) {
	return nil, ctrl.svc.Archive(ctx, r.TeamID)
}

func (ctrl *Team) Merge(ctx context.Context, r *server.TeamMergeRequest) (interface{}, error) {
	return nil, ctrl.svc.Merge(ctx, r.TeamID, r.Destination)
}

func (ctrl *Team) Move(ctx context.Context, r *server.TeamMoveRequest) (interface{}, error) {
	return nil, ctrl.svc.Move(ctx, r.TeamID, r.Organisation_id)
}
