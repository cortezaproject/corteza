package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	team struct {
		db  *factory.DB
		ctx context.Context

		team repository.TeamRepository
	}

	TeamService interface {
		With(ctx context.Context) TeamService

		FindByID(teamID uint64) (*types.Team, error)
		Find(filter *types.TeamFilter) ([]*types.Team, error)

		Create(team *types.Team) (*types.Team, error)
		Update(team *types.Team) (*types.Team, error)
		Merge(teamID, targetTeamID uint64) error
		Move(teamID, organisationID uint64) error

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error

		MemberAdd(teamID, userID uint64) error
		MemberRemove(teamID, userID uint64) error
	}
)

func Team() TeamService {
	return (&team{}).With(context.Background())
}

func (svc *team) With(ctx context.Context) TeamService {
	db := repository.DB(ctx)
	return &team{
		db:   db,
		ctx:  ctx,
		team: repository.Team(ctx, db),
	}
}

func (svc *team) FindByID(id uint64) (*types.Team, error) {
	// @todo: permission check if current user has access to this team
	return svc.team.FindByID(id)
}

func (svc *team) Find(filter *types.TeamFilter) ([]*types.Team, error) {
	// @todo: permission check to return only teams that current user has access to
	return svc.team.Find(filter)
}

func (svc *team) Create(mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team

	return svc.team.Create(mod)
}

func (svc *team) Update(mod *types.Team) (*types.Team, error) {
	// @todo: permission check if current user can add/edit team
	// @todo: make sure archived & deleted entries can not be edited

	return svc.team.Update(mod)
}

func (svc *team) Delete(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.team.DeleteByID(id)
}

func (svc *team) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that team has been removed (remove from web UI)
	// @todo: permissions check if current user can remove team
	return svc.team.ArchiveByID(id)
}

func (svc *team) Unarchive(id uint64) error {
	// @todo: permissions check if current user can unarchive team
	// @todo: make history accessible
	// @todo: notify users that team has been unarchived
	return svc.team.UnarchiveByID(id)
}

func (svc *team) Merge(id, targetTeamID uint64) error {
	// @todo: permission check if current user can merge team
	return svc.team.MergeByID(id, targetTeamID)
}

func (svc *team) Move(id, targetOrganisationID uint64) error {
	// @todo: permission check if current user can move team to another organisation
	return svc.team.MoveByID(id, targetOrganisationID)
}

func (svc *team) MemberAdd(id, userID uint64) error {
	// @todo: permission check if current user can add user in to a team
	return svc.team.MemberAddByID(id, userID)
}

func (svc *team) MemberRemove(id, userID uint64) error {
	// @todo: permission check if current user can remove user from a team
	return svc.team.MemberRemoveByID(id, userID)
}

var _ TeamService = &team{}
