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

func (r organisation) FindByID(ctx context.Context, id uint64) (*types.Organisation, error) {
	db := factory.Database.MustGet()
	sql := "SELECT * FROM organisations WHERE id = ? AND " + sqlOrganisationScope
	mod := &types.Organisation{}

	return mod, isFound(db.Get(mod, sql, id), mod.ID > 0, ErrOrganisationNotFound)
}

func (r organisation) Find(ctx context.Context, filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	db := factory.Database.MustGet()
	rval := make([]*types.Organisation, 0)
	params := make([]interface{}, 0)
	sql := "SELECT * FROM organisations WHERE " + sqlOrganisationScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, db.With(ctx).Select(&rval, sql, params...)
}

func (r organisation) Create(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, factory.Database.MustGet().With(ctx).Insert("organisations", mod)
}

func (r organisation) Update(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, factory.Database.MustGet().With(ctx).Replace("organisations", mod)
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
