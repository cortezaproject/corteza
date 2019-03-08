package rest

import (
	"context"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Application struct {
	svc struct {
		application service.ApplicationService
	}
}

func (Application) New() *Application {
	ctrl := &Application{}
	ctrl.svc.application = service.DefaultApplication

	return ctrl
}

func (ctrl *Application) List(ctx context.Context, r *request.ApplicationList) (interface{}, error) {
	return ctrl.svc.application.With(ctx).Find()
}

func (ctrl *Application) Create(ctx context.Context, r *request.ApplicationCreate) (interface{}, error) {
	app := &types.Application{
		Name:    r.Name,
		Enabled: r.Enabled,
	}

	if r.Unify != nil {
		app.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(app.Unify); err != nil {
			return nil, err
		}
	}

	return ctrl.svc.application.With(ctx).Create(app)
}

func (ctrl *Application) Update(ctx context.Context, r *request.ApplicationUpdate) (interface{}, error) {
	app := &types.Application{
		ID:      r.ApplicationID,
		Name:    r.Name,
		Enabled: r.Enabled,
	}

	if r.Unify != nil {
		app.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(app.Unify); err != nil {
			return nil, err
		}
	}

	return ctrl.svc.application.With(ctx).Update(app)
}

func (ctrl *Application) Read(ctx context.Context, r *request.ApplicationRead) (interface{}, error) {
	return ctrl.svc.application.With(ctx).FindByID(r.ApplicationID)
}

func (ctrl *Application) Delete(ctx context.Context, r *request.ApplicationDelete) (interface{}, error) {
	return nil, ctrl.svc.application.With(ctx).DeleteByID(r.ApplicationID)
}
