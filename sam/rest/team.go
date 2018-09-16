package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

type (
	Team struct {
		svc struct {
			team service.TeamService
		}
	}
)

func (Team) New() *Team {
	ctrl := &Team{}
	ctrl.svc.team = service.DefaultTeam
	return ctrl
}

func (ctrl *Team) Read(ctx context.Context, r *request.TeamRead) (interface{}, error) {
	return ctrl.svc.team.FindByID(ctx, r.TeamID)
}

func (ctrl *Team) List(ctx context.Context, r *request.TeamList) (interface{}, error) {
	return ctrl.svc.team.Find(ctx, &types.TeamFilter{Query: r.Query})
}

func (ctrl *Team) Create(ctx context.Context, r *request.TeamCreate) (interface{}, error) {
	org := &types.Team{
		Name: r.Name,
	}

	return ctrl.svc.team.Create(ctx, org)
}

func (ctrl *Team) Edit(ctx context.Context, r *request.TeamEdit) (interface{}, error) {
	org := &types.Team{
		ID:   r.TeamID,
		Name: r.Name,
	}

	return ctrl.svc.team.Update(ctx, org)
}

func (ctrl *Team) Remove(ctx context.Context, r *request.TeamRemove) (interface{}, error) {
	return nil, ctrl.svc.team.Delete(ctx, r.TeamID)
}

func (ctrl *Team) Archive(ctx context.Context, r *request.TeamArchive) (interface{}, error) {
	return nil, ctrl.svc.team.Archive(ctx, r.TeamID)
}

func (ctrl *Team) Merge(ctx context.Context, r *request.TeamMerge) (interface{}, error) {
	return nil, ctrl.svc.team.Merge(ctx, r.TeamID, r.Destination)
}

func (ctrl *Team) Move(ctx context.Context, r *request.TeamMove) (interface{}, error) {
	return nil, ctrl.svc.team.Move(ctx, r.TeamID, r.Organisation_id)
}
