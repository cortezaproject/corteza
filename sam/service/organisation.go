package service

import (
	"context"
	"github.com/crusttech/crust/sam/repository"
	"github.com/crusttech/crust/sam/types"
)

type (
	organisation struct {
		repository organisationRepository
	}

	organisationRepository interface {
		FindById(context.Context, uint64) (*types.Organisation, error)
		Find(context.Context, *types.OrganisationFilter) ([]*types.Organisation, error)

		Create(context.Context, *types.Organisation) (*types.Organisation, error)
		Update(context.Context, *types.Organisation) (*types.Organisation, error)

		deleter
		archiver
	}
)

func Organisation() *organisation {
	return &organisation{repository: repository.Organisation()}
}

func (svc organisation) FindById(ctx context.Context, id uint64) (*types.Organisation, error) {
	// @todo: permission check if current user can read organisation
	return svc.repository.FindById(ctx, id)
}

func (svc organisation) Find(ctx context.Context, filter *types.OrganisationFilter) ([]*types.Organisation, error) {
	// @todo: permission check to return only organisations that organisation has access to
	// @todo: actual searching not just a full select
	return svc.repository.Find(ctx, filter)
}

func (svc organisation) Create(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.repository.Create(ctx, mod)
}

func (svc organisation) Update(ctx context.Context, mod *types.Organisation) (*types.Organisation, error) {
	// @todo: permission check if current user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	return svc.repository.Update(ctx, mod)
}

func (svc organisation) Delete(ctx context.Context, id uint64) error {
	// @todo: permissions check if current user can remove organisation
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	return svc.repository.Delete(ctx, id)
}

func (svc organisation) Archive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can archive organisation
	return svc.repository.Archive(ctx, id)
}

func (svc organisation) Unarchive(ctx context.Context, id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that organisation has been removed (remove from web UI)
	// @todo: permissions check if current user can unarchive organisation
	return svc.repository.Unarchive(ctx, id)
}
