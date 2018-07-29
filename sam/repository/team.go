package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Team interface {
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
)

const (
	sqlTeamScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrTeamNotFound = repositoryError("TeamNotFound")
)

func (r *repository) FindTeamByID(id uint64) (*types.Team, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM teams WHERE id = ? AND " + sqlTeamScope
	mod := &types.Team{}

	return mod, isFound(db.Get(mod, sql, id), mod.ID > 0, ErrTeamNotFound)
}

func (r *repository) FindTeams(filter *types.TeamFilter) ([]*types.Team, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Team, 0)
	params := make([]interface{}, 0)

	sql := "SELECT * FROM teams WHERE " + sqlTeamScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, db.With(r.ctx).Select(&rval, sql, params...)
}

func (r *repository) CreateTeam(mod *types.Team) (*types.Team, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, factory.Database.MustGet().With(r.ctx).Insert("teams", mod)
}

func (r *repository) UpdateTeam(mod *types.Team) (*types.Team, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, factory.Database.MustGet().With(r.ctx).Replace("teams", mod)
}

func (r *repository) ArchiveTeamByID(id uint64) error {
	return simpleUpdate(r.ctx, "teams", "archived_at", time.Now(), id)
}

func (r *repository) UnarchiveTeamByID(id uint64) error {
	return simpleUpdate(r.ctx, "teams", "archived_at", nil, id)
}

func (r *repository) DeleteTeamByID(id uint64) error {
	return simpleDelete(r.ctx, "teams", id)
}

func (r *repository) MergeTeamByID(id, targetTeamID uint64) error {
	return ErrNotImplemented
}

func (r *repository) MoveTeamByID(id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}
