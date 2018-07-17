package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlOrganisationScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrOrganisationNotFound = repositoryError("OrganisationNotFound")
)

type (
	organisation struct{}
)

func Organisation() organisation {
	return organisation{}
}

func (r organisation) FindById(ctx context.Context, id uint64) (*types.Organisation, error) {
	db := factory.Database.MustGet()

	mod := &types.Organisation{}
	if err := db.Get(mod, "SELECT * FROM organisations WHERE id = ? AND "+sqlOrganisationScope, id); err != nil {
		return nil, ErrDatabaseError
	} else if mod.ID == 0 {
		return nil, ErrOrganisationNotFound
	} else {
		return mod, nil
	}
}

func (r organisation) Find(ctx context.Context, filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	db := factory.Database.MustGet()

	var params = make([]interface{}, 0)
	sql := "SELECT * FROM organisations WHERE " + sqlOrganisationScope

	if filter != nil {
		if filter.Query != "" {
			sql += "AND ame LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	rval := make([]*types.Organisation, 0)
	if err := db.Select(&rval, sql, params...); err != nil {
		return nil, ErrDatabaseError
	} else {
		return rval, nil
	}
}

func (r organisation) Create(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("organisations", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r organisation) Update(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	db := factory.Database.MustGet()

	if err := db.Replace("organisations", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r organisation) Archive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "organisations", "archived_at", time.Now(), id)
}

func (r organisation) Unarchive(ctx context.Context, id uint64) error {
	return simpleUpdate(ctx, "organisations", "archived_at", nil, id)
}

func (r organisation) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "organisations", id)
}
