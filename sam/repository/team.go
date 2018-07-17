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

func (r team) FindById(ctx context.Context, id uint64) (*types.Team, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod := &types.Team{}
	if err := db.Get(mod, "SELECT * FROM teams WHERE id = ? AND "+sqlTeamScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrTeamNotFound
	} else {
		return mod, nil
	}
}

func (r team) Find(ctx context.Context, filter *types.TeamFilter) ([]*types.Team, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	var params = make([]interface{}, 0)
	sql := "SELECT * FROM teams WHERE " + sqlTeamScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	rval := make([]*types.Team, 0)
	if err := db.Select(&rval, sql, params...); err != nil {
		return nil, ErrDatabaseError
	} else {
		return rval, nil
	}
}

func (r team) Create(ctx context.Context, mod *types.Team) (*types.Team, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("teams", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r team) Update(ctx context.Context, mod *types.Team) (*types.Team, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, ErrDatabaseError
	}

	if err := db.Replace("teams", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
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

func (r team) Merge(ctx context.Context, id, targetTeamId uint64) error {
	return ErrNotImplemented
}

func (r team) Move(ctx context.Context, id, targetOrganisationId uint64) error {
	return ErrNotImplemented
}
