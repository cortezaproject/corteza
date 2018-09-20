package service

import (
	"context"

	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	organisation struct {
		ctx context.Context
		rpo repository.Organisation
	}

	OrganisationService interface {
		With(ctx context.Context) OrganisationService

		FindByID(organisationID uint64) (*types.Organisation, error)
		Find(filter *types.OrganisationFilter) ([]*types.Organisation, error)

		Create(organisation *types.Organisation) (*types.Organisation, error)
		Update(organisation *types.Organisation) (*types.Organisation, error)

		deleter
		archiver
	}
)

func Organisation() *organisation {
	return &organisation{
		ctx: context.Background(),
		rpo: repository.NewOrganisation(context.Background()),
	}
}

func (svc *organisation) With(ctx context.Context) OrganisationService {
	return &organisation{
		ctx: ctx,
		rpo: svc.rpo.With(ctx),
	}
}

func (svc *organisation) FindByID(id uint64) (*types.Organisation, error) {
	// @todo: permission check if current user can read organisation
	return svc.rpo.FindOrganisationByID(id)
}

func (svc *organisation) Find(filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	// @todo: permission check to return only organisations that organisation has access to
	// @todo: actual searching not just a full select
	return svc.rpo.FindOrganisations(filter)
}

func (svc *organisation) Create(mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.CreateOrganisation(mod)
}

func (svc *organisation) Update(mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.rpo.UpdateOrganisation(mod)
}

func (svc *organisation) Delete(id uint64) error {
	// @todo: permissions check if current user can remove organisation
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	return svc.rpo.DeleteOrganisationByID(id)
}

func (svc *organisation) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can archive organisation
	return svc.rpo.ArchiveOrganisationByID(id)
}

func (svc *organisation) Unarchive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can unarchive organisation
	return svc.rpo.UnarchiveOrganisationByID(id)
}

var _ OrganisationService = &organisation{}
