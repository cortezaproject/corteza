package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	team struct {
		rpo teamRepository
	}

	teamRepository interface {
		repository.Transactionable
		repository.Team
	}
)

func Team() *team {
	return &team{rpo: repository.New()}
}

func (svc team) FindByID(ctx context.Context, id uint64) (*types.Team, error) {
	// @todo: permission check if current user has access to this team
	return svc.rpo.FindTeamByID(id)
}

func (svc team) Find(ctx context.Context, filter *types.TeamFilter) ([]*types.Team, error) {
	// @todo: permission check to return only teams that current user has access to
	return svc.rpo.FindTeams(filter)
}

func (svc team) Create(ctx context.Context, mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team

	return svc.rpo.CreateTeam(mod)
}

func (svc team) Update(ctx context.Context, mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.UpdateTeam(mod)
}

func (svc team) Delete(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.rpo.DeleteTeamByID(id)
}

func (svc team) Archive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.rpo.ArchiveTeamByID(id)
}

func (svc team) Unarchive(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can unarchive team
	// @todo: make history accessible
	// @todo: notify users that team has been unarchived
	return svc.rpo.UnarchiveTeamByID(id)
}

func (svc team) Merge(ctx context.Context, id, targetTeamID uint64) error {
	// @todo: permission check if current user can merge team
	return svc.rpo.MergeTeamByID(id, targetTeamID)
}

func (svc team) Move(ctx context.Context, id, targetOrganisationID uint64) error {
	// @todo: permission check if current user can move team to another organisation
	return svc.rpo.MoveTeamByID(id, targetOrganisationID)
}
