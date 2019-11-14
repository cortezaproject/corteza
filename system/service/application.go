package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/repository"
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
		CanAccess(context.Context) bool
		CanCreateApplication(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool

		FilterReadableApplications(ctx context.Context) *permissions.ResourceFilter
	}

	ApplicationService interface {
		With(ctx context.Context) ApplicationService

		FindByID(applicationID uint64) (*types.Application, error)
		Find(types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error)

		Create(application *types.Application) (*types.Application, error)
		Update(application *types.Application) (*types.Application, error)
		Delete(uint64) error
		Undelete(uint64) error
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

func (svc *application) FindByID(ID uint64) (app *types.Application, err error) {
	if ID == 0 {
		return nil, ErrInvalidID
	}

	if app, err = svc.application.FindByID(ID); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadApplication(svc.ctx, app) {
		return nil, ErrNoPermissions.withStack()
	}

	return app, nil
}

func (svc *application) Find(f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error) {
	f.IsReadable = svc.ac.FilterReadableApplications(svc.ctx)

	if f.Deleted > rh.FilterStateExcluded {
		// If list with deleted applications is requested
		// user must have access permissions to system (ie: is admin)
		//
		// not the best solution but ATM it allows us to have at least
		// some kind of control over who can see deleted applications
		if !svc.ac.CanAccess(svc.ctx) {
			return nil, f, ErrNoPermissions.withStack()
		}
	}

	return svc.application.Find(f)
}

func (svc *application) Create(mod *types.Application) (*types.Application, error) {
	if !svc.ac.CanCreateApplication(svc.ctx) {
		return nil, ErrNoPermissions.withStack()
	}
	return svc.application.Create(mod)
}

func (svc *application) Update(mod *types.Application) (t *types.Application, err error) {
	if !svc.ac.CanUpdateApplication(svc.ctx, mod) {
		return nil, ErrNoPermissions.withStack()
	}

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

func (svc *application) Delete(ID uint64) error {
	app := &types.Application{ID: ID}
	if !svc.ac.CanDeleteApplication(svc.ctx, app) {
		return ErrNoPermissions.withStack()
	}
	return svc.application.DeleteByID(ID)
}

func (svc *application) Undelete(ID uint64) error {
	app := &types.Application{ID: ID}
	if !svc.ac.CanDeleteApplication(svc.ctx, app) {
		return ErrNoPermissions.withStack()
	}
	return svc.application.UndeleteByID(ID)
}

var _ ApplicationService = &application{}
