package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Organisation struct {
		svc struct {
			org service.OrganisationService
		}
	}
)

func (Organisation) New() *Organisation {
	ctrl := &Organisation{}
	ctrl.svc.org = service.DefaultOrganisation
	return ctrl
}

func (ctrl *Organisation) Read(ctx context.Context, r *request.OrganisationRead) (interface{}, error) {
	return ctrl.svc.org.FindByID(ctx, r.ID)
}

func (ctrl *Organisation) List(ctx context.Context, r *request.OrganisationList) (interface{}, error) {
	return ctrl.svc.org.Find(ctx, &types.OrganisationFilter{Query: r.Query})
}

func (ctrl *Organisation) Create(ctx context.Context, r *request.OrganisationCreate) (interface{}, error) {
	org := &types.Organisation{
		Name: r.Name,
	}

	return ctrl.svc.org.Create(ctx, org)
}

func (ctrl *Organisation) Edit(ctx context.Context, r *request.OrganisationEdit) (interface{}, error) {
	org := &types.Organisation{
		ID:   r.ID,
		Name: r.Name,
	}

	return ctrl.svc.org.Update(ctx, org)
}

func (ctrl *Organisation) Remove(ctx context.Context, r *request.OrganisationRemove) (interface{}, error) {
	return nil, ctrl.svc.org.Delete(ctx, r.ID)
}

func (ctrl *Organisation) Archive(ctx context.Context, r *request.OrganisationArchive) (interface{}, error) {
	return nil, ctrl.svc.org.Archive(ctx, r.ID)
}
