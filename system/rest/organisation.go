package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/types"
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
	return ctrl.svc.org.With(ctx).FindByID(r.ID)
}

func (ctrl *Organisation) List(ctx context.Context, r *request.OrganisationList) (interface{}, error) {
	return ctrl.svc.org.With(ctx).Find(&types.OrganisationFilter{Query: r.Query})
}

func (ctrl *Organisation) Create(ctx context.Context, r *request.OrganisationCreate) (interface{}, error) {
	org := &types.Organisation{
		Name: r.Name,
	}

	return ctrl.svc.org.With(ctx).Create(org)
}

func (ctrl *Organisation) Update(ctx context.Context, r *request.OrganisationUpdate) (interface{}, error) {
	org := &types.Organisation{
		ID:   r.ID,
		Name: r.Name,
	}

	return ctrl.svc.org.With(ctx).Update(org)
}

func (ctrl *Organisation) Delete(ctx context.Context, r *request.OrganisationDelete) (interface{}, error) {
	return resputil.OK(), ctrl.svc.org.With(ctx).Delete(r.ID)
}

func (ctrl *Organisation) Archive(ctx context.Context, r *request.OrganisationArchive) (interface{}, error) {
	return resputil.OK(), ctrl.svc.org.With(ctx).Archive(r.ID)
}
