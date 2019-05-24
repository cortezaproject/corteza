package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	application struct {
		db  *factory.DB
		ctx context.Context

		ac applicationAccessController

		application repository.ApplicationRepository
	}

	applicationAccessController interface {
		CanCreateApplication(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
	}

	ApplicationService interface {
		With(ctx context.Context) ApplicationService

		FindByID(applicationID uint64) (*types.Application, error)
		Find() (types.ApplicationSet, error)

		Create(application *types.Application) (*types.Application, error)
		Update(application *types.Application) (*types.Application, error)
		DeleteByID(id uint64) error
	}
)

func Application(ctx context.Context) ApplicationService {
	return (&application{
		ac: DefaultAccessControl,
	}).With(ctx)

}

func (svc *application) With(ctx context.Context) ApplicationService {
	db := repository.DB(ctx)
	return &application{
		db:          db,
		ctx:         ctx,
		ac:          svc.ac,
		application: repository.Application(ctx, db),
	}
}

func (svc *application) FindByID(id uint64) (*types.Application, error) {
	app, err := svc.application.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadApplication(svc.ctx, app) {
		return nil, errors.New("Not allowed to access application")
	}

	return app, nil
}

func (svc *application) Find() (types.ApplicationSet, error) {
	apps, err := svc.application.Find()
	if err != nil {
		return nil, err
	}

	ret := []*types.Application{}
	for _, app := range apps {
		if svc.ac.CanReadApplication(svc.ctx, app) {
			ret = append(ret, app)
		} //
	}
	return ret, nil
}

func (svc *application) Create(mod *types.Application) (*types.Application, error) {
	if !svc.ac.CanCreateApplication(svc.ctx) {
		return nil, errors.New("Not allowed to create application")
	}
	return svc.application.Create(mod)
}

func (svc *application) Update(mod *types.Application) (t *types.Application, err error) {
	if !svc.ac.CanUpdateApplication(svc.ctx, mod) {
		return nil, errors.New("Not allowed to update application")
	}

	// @todo: make sure archived & deleted entries can not be edited

	return t, svc.db.Transaction(func() (err error) {
		if t, err = svc.application.FindByID(mod.ID); err != nil {
			return
		}

		// Assign changed values
		t.Name = mod.Name
		t.Enabled = mod.Enabled
		t.Unify = mod.Unify

		if t, err = svc.application.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc *application) DeleteByID(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that application has been removed (remove from web UI)

	app := &types.Application{ID: id}
	if !svc.ac.CanDeleteApplication(svc.ctx, app) {
		return errors.New("Not allowed to delete application")
	}
	return svc.application.DeleteByID(id)
}

var _ ApplicationService = &application{}
