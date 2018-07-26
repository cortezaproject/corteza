package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Organisation struct {
		svc organisationService
	}

	organisationService interface {
		FindByID(ctx context.Context, organisationID uint64) (*types.Organisation, error)
		Find(ctx context.Context, filter *types.OrganisationFilter) ([]*types.Organisation, error)

		Create(ctx context.Context, organisation *types.Organisation) (*types.Organisation, error)
		Update(ctx context.Context, organisation *types.Organisation) (*types.Organisation, error)

		deleter
		archiver
	}
)

func (Organisation) New(organisationSvc organisationService) *Organisation {
	var ctrl = &Organisation{}
	ctrl.svc = organisationSvc
	return ctrl
}

func (ctrl *Organisation) Read(ctx context.Context, r *server.OrganisationReadRequest) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.ID)
}

func (ctrl *Organisation) List(ctx context.Context, r *server.OrganisationListRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.OrganisationFilter{Query: r.Query})
}

func (ctrl *Organisation) Create(ctx context.Context, r *server.OrganisationCreateRequest) (interface{}, error) {
	org := types.Organisation{}.
		New().
		SetName(r.Name)

	return ctrl.svc.Create(ctx, org)
}

func (ctrl *Organisation) Edit(ctx context.Context, r *server.OrganisationEditRequest) (interface{}, error) {
	org := types.Organisation{}.
		New().
		SetID(r.ID).
		SetName(r.Name)

	return ctrl.svc.Update(ctx, org)
}

func (ctrl *Organisation) Remove(ctx context.Context, r *server.OrganisationRemoveRequest) (interface{}, error) {
	return nil, ctrl.svc.Delete(ctx, r.ID)
}

func (ctrl *Organisation) Archive(ctx context.Context, r *server.OrganisationArchiveRequest) (interface{}, error) {
	return nil, ctrl.svc.Archive(ctx, r.ID)
}
