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
		service organisationService
	}

	organisationService interface {
		FindByID(context.Context, uint64) (*types.Organisation, error)
		Find(context.Context, *types.OrganisationFilter) ([]*types.Organisation, error)

		Create(context.Context, *types.Organisation) (*types.Organisation, error)
		Update(context.Context, *types.Organisation) (*types.Organisation, error)

		deleter
		archiver
	}
)

func (Organisation) New() *Organisation {
	return &Organisation{}
}

func (ctrl *Organisation) Read(ctx context.Context, r *server.OrganisationReadRequest) (interface{}, error) {
	return ctrl.service.FindByID(ctx, r.ID)
}

func (ctrl *Organisation) List(ctx context.Context, r *server.OrganisationListRequest) (interface{}, error) {
	return ctrl.service.Find(ctx, &types.OrganisationFilter{Query: r.Query})
}

func (ctrl *Organisation) Create(ctx context.Context, r *server.OrganisationCreateRequest) (interface{}, error) {
	org := types.Organisation{}.
		New().
		SetName(r.Name)

	return ctrl.service.Create(ctx, org)
}

func (ctrl *Organisation) Edit(ctx context.Context, r *server.OrganisationEditRequest) (interface{}, error) {
	org := types.Organisation{}.
		New().
		SetID(r.ID).
		SetName(r.Name)

	return ctrl.service.Update(ctx, org)
}

func (ctrl *Organisation) Remove(ctx context.Context, r *server.OrganisationRemoveRequest) (interface{}, error) {
	return nil, ctrl.service.Delete(ctx, r.ID)
}

func (ctrl *Organisation) Archive(ctx context.Context, r *server.OrganisationArchiveRequest) (interface{}, error) {
	return nil, ctrl.service.Archive(ctx, r.ID)
}
