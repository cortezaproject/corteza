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

		FindByID(id uint64) (*types.Team, error)
		FindByMemberID(userID uint64) ([]*types.Team, error)
		Find(filter *types.TeamFilter) ([]*types.Team, error)

		Create(mod *types.Team) (*types.Team, error)
		Update(mod *types.Team) (*types.Team, error)

		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error

		MergeByID(id, targetTeamID uint64) error
		MoveByID(id, targetOrganisationID uint64) error

		MemberAddByID(id, userID uint64) error
		MemberRemoveByID(id, userID uint64) error
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

func (r *team) FindByID(id uint64) (*types.Team, error) {
	sql := "SELECT * FROM " + r.teams + " WHERE id = ? AND " + sqlTeamScope
	mod := &types.Team{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrTeamNotFound)
}

func (r *team) FindByMemberID(userID uint64) ([]*types.Team, error) {
	ids := make([]uint64, 0)
	params := make([]interface{}, 0)

	sql := "SELECT DISTINCT rel_team FROM " + r.members + "  "
	sql += "WHERE rel_user = ?"
	params = append(params, userID)

	if err := r.db().Select(&ids, sql, params...); err != nil {
		return nil, err
	}

	rval := make([]*types.Team, 0)
	for _, id := range ids {
		mod, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}
		rval = append(rval, mod)
	}

	return rval, nil
}

func (r *team) Find(filter *types.TeamFilter) ([]*types.Team, error) {
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

func (r *team) Create(mod *types.Team) (*types.Team, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.teams, mod)
}

func (r *team) Update(mod *types.Team) (*types.Team, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace(r.teams, mod)
}

func (r *team) ArchiveByID(id uint64) error {
	return r.updateColumnByID(r.teams, "archived_at", time.Now(), id)
}

func (r *team) UnarchiveByID(id uint64) error {
	return r.updateColumnByID(r.teams, "archived_at", nil, id)
}

func (r *team) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.teams, "deleted_at", time.Now(), id)
}

func (r *team) MergeByID(id, targetTeamID uint64) error {
	return ErrNotImplemented
}

func (r *team) MoveByID(id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}

func (r *team) MemberAddByID(id, userID uint64) error {
	mod := &types.TeamMember{
		TeamID: id,
		UserID: userID,
	}
	return r.db().Replace(r.members, mod)
}

func (r *team) MemberRemoveByID(id, userID uint64) error {
	mod := &types.TeamMember{
		TeamID: id,
		UserID: userID,
	}
	return r.db().Delete(r.members, mod, "rel_team", "rel_user")
}
