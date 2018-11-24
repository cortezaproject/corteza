package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"
)

type (
	TeamRepository interface {
		With(ctx context.Context, db *factory.DB) TeamRepository

		FindTeamByID(id uint64) (*types.Team, error)
		FindTeams(filter *types.TeamFilter) ([]*types.Team, error)
		CreateTeam(mod *types.Team) (*types.Team, error)
		UpdateTeam(mod *types.Team) (*types.Team, error)
		ArchiveTeamByID(id uint64) error
		UnarchiveTeamByID(id uint64) error
		DeleteTeamByID(id uint64) error
		MergeTeamByID(id, targetTeamID uint64) error
		MoveTeamByID(id, targetOrganisationID uint64) error
	}

	team struct {
		*repository

		// sql table reference
		teams   string
		members string
	}
)

const (
	sqlTeamScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrTeamNotFound = repositoryError("TeamNotFound")
)

func Team(ctx context.Context, db *factory.DB) TeamRepository {
	return (&team{}).With(ctx, db)
}

func (r *team) With(ctx context.Context, db *factory.DB) TeamRepository {
	return &team{
		repository: r.repository.With(ctx, db),
		teams:      "sys_team",
		members:    "sys_team_member",
	}
}

func (r *team) FindTeamByID(id uint64) (*types.Team, error) {
	sql := "SELECT * FROM " + r.teams + " WHERE id = ? AND " + sqlTeamScope
	mod := &types.Team{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrTeamNotFound)
}

func (r *team) FindTeams(filter *types.TeamFilter) ([]*types.Team, error) {
	rval := make([]*types.Team, 0)
	params := make([]interface{}, 0)

	sql := "SELECT * FROM " + r.teams + " WHERE " + sqlTeamScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *team) CreateTeam(mod *types.Team) (*types.Team, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.teams, mod)
}

func (r *team) UpdateTeam(mod *types.Team) (*types.Team, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace(r.teams, mod)
}

func (r *team) ArchiveTeamByID(id uint64) error {
	return r.updateColumnByID(r.teams, "archived_at", time.Now(), id)
}

func (r *team) UnarchiveTeamByID(id uint64) error {
	return r.updateColumnByID(r.teams, "archived_at", nil, id)
}

func (r *team) DeleteTeamByID(id uint64) error {
	return r.updateColumnByID(r.teams, "deleted_at", time.Now(), id)
}

func (r *team) MergeTeamByID(id, targetTeamID uint64) error {
	return ErrNotImplemented
}

func (r *team) MoveTeamByID(id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}
