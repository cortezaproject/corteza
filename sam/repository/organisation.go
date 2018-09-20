package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	Organisation interface {
		With(ctx context.Context) Organisation

		FindOrganisationByID(id uint64) (*types.Organisation, error)
		FindOrganisations(filter *types.OrganisationFilter) ([]*types.Organisation, error)
		CreateOrganisation(mod *types.Organisation) (*types.Organisation, error)
		UpdateOrganisation(mod *types.Organisation) (*types.Organisation, error)
		ArchiveOrganisationByID(id uint64) error
		UnarchiveOrganisationByID(id uint64) error
		DeleteOrganisationByID(id uint64) error
	}

	organisation struct {
		*repository
	}
)

const (
	sqlOrganisationScope = "deleted_at IS NULL AND archived_at IS NULL"

	ErrOrganisationNotFound = repositoryError("OrganisationNotFound")
)

func NewOrganisation(ctx context.Context) Organisation {
	return (&organisation{}).With(ctx)
}

func (r *organisation) With(ctx context.Context) Organisation {
	return &organisation{
		repository: r.repository.With(ctx),
	}
}

func (r *organisation) FindOrganisationByID(id uint64) (*types.Organisation, error) {
	sql := "SELECT * FROM organisations WHERE id = ? AND " + sqlOrganisationScope
	mod := &types.Organisation{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrOrganisationNotFound)
}

func (r *organisation) FindOrganisations(filter *types.OrganisationFilter) ([]*types.Organisation, error) {
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

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *organisation) CreateOrganisation(mod *types.Organisation) (*types.Organisation, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("organisations", mod)
}

func (r *organisation) UpdateOrganisation(mod *types.Organisation) (*types.Organisation, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace("organisations", mod)
}

func (r *organisation) ArchiveOrganisationByID(id uint64) error {
	return r.updateColumnByID("organisations", "archived_at", time.Now(), id)
}

func (r *organisation) UnarchiveOrganisationByID(id uint64) error {
	return r.updateColumnByID("organisations", "archived_at", nil, id)
}

func (r *organisation) DeleteOrganisationByID(id uint64) error {
	return r.updateColumnByID("organisations", "deleted_at", nil, id)
}
