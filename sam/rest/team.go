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
		service teamService
	}

	teamService interface {
		FindByID(context.Context, uint64) (*types.Team, error)
		Find(context.Context, *types.TeamFilter) ([]*types.Team, error)

		Create(context.Context, *types.Team) (*types.Team, error)
		Update(context.Context, *types.Team) (*types.Team, error)
		Merge(context.Context, *types.Team) error
		Move(context.Context, *types.Team) error

		deleter
		archiver
	}
)

func (Team) New() *Team {
	return &Team{}
}

func (ctrl *Team) Read(ctx context.Context, r *server.TeamReadRequest) (interface{}, error) {
	return ctrl.service.FindByID(ctx, r.TeamID)
}

func (ctrl *Team) List(ctx context.Context, r *server.TeamListRequest) (interface{}, error) {
	return ctrl.service.Find(ctx, &types.TeamFilter{Query: r.Query})
}

func (ctrl *Team) Create(ctx context.Context, r *server.TeamCreateRequest) (interface{}, error) {
	org := types.Team{}.
		New().
		SetName(r.Name)

	return ctrl.service.Create(ctx, org)
}

func (ctrl *Team) Edit(ctx context.Context, r *server.TeamEditRequest) (interface{}, error) {
	org := types.Team{}.
		New().
		SetID(r.TeamID).
		SetName(r.Name)

	return ctrl.service.Update(ctx, org)
}

func (ctrl *Team) Remove(ctx context.Context, r *server.TeamRemoveRequest) (interface{}, error) {
	return nil, ctrl.service.Delete(ctx, r.TeamID)
}

func (ctrl *Team) Archive(ctx context.Context, r *server.TeamArchiveRequest) (interface{}, error) {
	return nil, ctrl.service.Archive(ctx, r.TeamID)
}

func (ctrl *Team) Merge(ctx context.Context, r *server.TeamMergeRequest) (interface{}, error) {
	return nil, ctrl.service.Merge(ctx, &types.Team{ID: r.TeamID})
}

func (ctrl *Team) Move(ctx context.Context, r *server.TeamMoveRequest) (interface{}, error) {
	return nil, ctrl.service.Move(ctx, &types.Team{ID: r.TeamID})
}
