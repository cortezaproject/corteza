package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	OrganisationRepository interface {
		With(ctx context.Context, db *factory.DB) OrganisationRepository

		FindByID(id uint64) (*types.Organisation, error)
		Find(filter *types.OrganisationFilter) ([]*types.Organisation, error)
		Create(mod *types.Organisation) (*types.Organisation, error)
		Update(mod *types.Organisation) (*types.Organisation, error)
		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error
	}

	organisation struct {
		*repository

		// sql table reference
		organisations string
	}
)

const (
	sqlOrganisationScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrOrganisationNotFound = repositoryError("OrganisationNotFound")
)

// @todo migrate to same pattern as we have for users
func Organisation(ctx context.Context, db *factory.DB) OrganisationRepository {
	return (&organisation{}).With(ctx, db)
}

func (r *organisation) With(ctx context.Context, db *factory.DB) OrganisationRepository {
	return &organisation{
		repository:    r.repository.With(ctx, db),
		organisations: "sys_organisation",
	}
}

func (r *organisation) FindByID(id uint64) (*types.Organisation, error) {
	sql := "SELECT * FROM " + r.organisations + " WHERE id = ? AND " + sqlOrganisationScope
	mod := &types.Organisation{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrOrganisationNotFound)
}

func (r *organisation) Find(filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	rval := make([]*types.Organisation, 0)
	params := make([]interface{}, 0)
	sql := "SELECT * FROM " + r.organisations + " WHERE " + sqlOrganisationScope

	if filter != nil {
		if filter.Query != "" {
			sql += " AND name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *organisation) Create(mod *types.Organisation) (*types.Organisation, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.organisations, mod)
}

func (r *organisation) Update(mod *types.Organisation) (*types.Organisation, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace(r.organisations, mod)
}

func (r *organisation) ArchiveByID(id uint64) error {
	return r.updateColumnByID(r.organisations, "archived_at", time.Now(), id)
}

func (r *organisation) UnarchiveByID(id uint64) error {
	return r.updateColumnByID(r.organisations, "archived_at", nil, id)
}

func (r *organisation) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.organisations, "deleted_at", time.Now(), id)
}
