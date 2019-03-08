package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	application struct {
		db  *factory.DB
		ctx context.Context

		application repository.ApplicationRepository
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

func Application() ApplicationService {
	return (&application{}).With(context.Background())
}

func (svc *application) With(ctx context.Context) ApplicationService {
	db := repository.DB(ctx)
	return &application{
		db:          db,
		ctx:         ctx,
		application: repository.Application(ctx, db),
	}
}

func (svc *application) FindByID(id uint64) (*types.Application, error) {
	// @todo: permission check if current user has access to this application
	return svc.application.FindByID(id)
}

func (svc *application) Find() (types.ApplicationSet, error) {
	// @todo: permission check to return only applications that current user has access to
	return svc.application.Find()
}

func (svc *application) Create(mod *types.Application) (*types.Application, error) {
	// @todo: permission check if current user can add/edit application

	return svc.application.Create(mod)
}

func (svc *application) Update(mod *types.Application) (t *types.Application, err error) {
	// @todo: permission check if current user can add/edit application
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
	// @todo: permissions check if current user can remove application
	return svc.application.DeleteByID(id)
}

var _ ApplicationService = &application{}
