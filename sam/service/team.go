package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	team struct {
		ctx  context.Context
		team repository.Team
	}

	TeamService interface {
		With(ctx context.Context) TeamService

		FindByID(teamID uint64) (*types.Team, error)
		Find(filter *types.TeamFilter) ([]*types.Team, error)

		Create(team *types.Team) (*types.Team, error)
		Update(team *types.Team) (*types.Team, error)
		Merge(teamID, targetTeamID uint64) error
		Move(teamID, organisationID uint64) error

		deleter
		archiver
	}
)

func Team() *team {
	return &team{
		ctx:  context.Background(),
		team: repository.NewTeam(context.Background()),
	}
}

func (svc *team) With(ctx context.Context) TeamService {
	return &team{
		ctx:  ctx,
		team: svc.team.With(ctx),
	}
}

func (svc *team) FindByID(id uint64) (*types.Team, error) {
	// @todo: permission check if current user has access to this team
	return svc.team.FindTeamByID(id)
}

func (svc *team) Find(filter *types.TeamFilter) ([]*types.Team, error) {
	// @todo: permission check to return only teams that current user has access to
	return svc.team.FindTeams(filter)
}

func (svc *team) Create(mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team

	return svc.team.CreateTeam(mod)
}

func (svc *team) Update(mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team
	// @todo: make sure archived & deleted entries can not be edited

	return svc.team.UpdateTeam(mod)
}

func (svc *team) Delete(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.team.DeleteTeamByID(id)
}

func (svc *team) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.team.ArchiveTeamByID(id)
}

func (svc *team) Unarchive(id uint64) error {
	// @todo: permissions check if current user can unarchive team
	// @todo: make history accessible
	// @todo: notify users that team has been unarchived
	return svc.team.UnarchiveTeamByID(id)
}

func (svc *team) Merge(id, targetTeamID uint64) error {
	// @todo: permission check if current user can merge team
	return svc.team.MergeTeamByID(id, targetTeamID)
}

func (svc *team) Move(id, targetOrganisationID uint64) error {
	// @todo: permission check if current user can move team to another organisation
	return svc.team.MoveTeamByID(id, targetOrganisationID)
}

var _ TeamService = &team{}
