package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlTeamScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrTeamNotFound = repositoryError("TeamNotFound")
)

type (
	team struct{}
)

func Team() team {
	return team{}
}

func (r team) FindByID(ctx context.Context, id uint64) (*types.Team, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM teams WHERE id = ? AND " + sqlTeamScope
	mod := &types.Team{}

	return mod, isFound(db.Get(mod, sql, id), mod.ID > 0, ErrTeamNotFound)
}

func (r team) Find(ctx context.Context, filter *types.TeamFilter) ([]*types.Team, error) {
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

	return rval, db.With(ctx).Select(&rval, sql, params...)
}

func (r team) Create(ctx context.Context, mod *types.Team) (*types.Team, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	return mod, db.With(ctx).Insert("teams", mod)
}

func (r team) Update(ctx context.Context, mod *types.Team) (*types.Team, error) {
	db := factory.Database.MustGet()

	now := time.Now()
	mod.SetUpdatedAt(&now)

	return mod, db.With(ctx).Replace("teams", mod)
}

func (r team) Archive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "teams", "archived_at", time.Now(), id)
}

func (r team) Unarchive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "teams", "archived_at", nil, id)
}

func (r team) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "teams", id)
}

func (r team) Merge(ctx context.Context, id, targetTeamID uint64) error {
	return ErrNotImplemented
}

func (r team) Move(ctx context.Context, id, targetOrganisationID uint64) error {
	return ErrNotImplemented
}
