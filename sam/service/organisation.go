package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	organisation struct {
		rpo organisationRepository
	}

	organisationRepository interface {
		repository.Transactionable
		repository.Organisation
	}
)

func Organisation() *organisation {
	return &organisation{rpo: repository.New()}
}

func (svc organisation) FindByID(ctx context.Context, id uint64) (*types.Organisation, error) {
	// @todo: permission check if current user can read organisation
	return svc.rpo.FindOrganisationByID(id)
}

func (svc organisation) Find(ctx context.Context, filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	// @todo: permission check to return only organisations that organisation has access to
	// @todo: actual searching not just a full select
	return svc.rpo.FindOrganisations(filter)
}

func (svc organisation) Create(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.CreateOrganisation(mod)
}

func (svc organisation) Update(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.UpdateOrganisation(mod)
}

func (svc organisation) Delete(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can remove organisation
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	return svc.rpo.DeleteOrganisationByID(id)
}

func (svc organisation) Archive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can archive organisation
	return svc.rpo.ArchiveOrganisationByID(id)
}

func (svc organisation) Unarchive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can unarchive organisation
	return svc.rpo.UnarchiveOrganisationByID(id)
}
