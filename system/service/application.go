package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	application struct {
		db  *factory.DB
		ctx context.Context

		ac       applicationAccessController
		eventbus eventDispatcher

		actionlog actionlog.Recorder

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
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(ctx)

}

func (svc *application) With(ctx context.Context) ApplicationService {
	db := repository.DB(ctx)
	return &application{
		db:  db,
		ctx: ctx,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		actionlog: DefaultActionlog,

		application: repository.Application(ctx, db),
	}
}

func (svc *application) FindByID(ID uint64) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{application: &types.Application{ID: ID}}
	)

	err = svc.db.Transaction(func() error {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.application.FindByID(ID); err != nil {
			return ApplicationErrInvalidID().Wrap(err)
		}

		aaProps.setApplication(app)

		if !svc.ac.CanReadApplication(svc.ctx, app) {
			return ApplicationErrNotAllowedToRead()
		}

		return nil
	})

	return app, svc.recordAction(svc.ctx, aaProps, ApplicationActionLookup, err)
}

func (svc *application) Find(filter types.ApplicationFilter) (aa types.ApplicationSet, f types.ApplicationFilter, err error) {
	var (
		aaProps = &applicationActionProps{filter: &filter}
	)

	err = svc.db.Transaction(func() error {
		filter.IsReadable = svc.ac.FilterReadableApplications(svc.ctx)

		if filter.Deleted > rh.FilterStateExcluded {
			// If list with deleted applications is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted applications
			if !svc.ac.CanAccess(svc.ctx) {
				return ApplicationErrNotAllowedToListApplications()
			}
		}

		aa, f, err = svc.application.Find(filter)
		return err
	})

	return aa, f, svc.recordAction(svc.ctx, aaProps, ApplicationActionSearch, err)
}

func (svc *application) Create(new *types.Application) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{new: new}
	)

	err = svc.db.Transaction(func() (err error) {
		if !svc.ac.CanCreateApplication(svc.ctx) {
			return ApplicationErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.ApplicationBeforeCreate(new, nil)); err != nil {
			return
		}

		if app, err = svc.application.Create(new); err != nil {
			return
		}

		aaProps.setApplication(app)

		_ = svc.eventbus.WaitFor(svc.ctx, event.ApplicationAfterCreate(new, nil))
		return nil
	})

	return app, svc.recordAction(svc.ctx, aaProps, ApplicationActionCreate, err)
}

func (svc *application) Update(upd *types.Application) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{update: upd}
	)

	err = svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.application.FindByID(upd.ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanUpdateApplication(svc.ctx, app) {
			return ApplicationErrNotAllowedToUpdate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.ApplicationBeforeUpdate(upd, app)); err != nil {
			return
		}

		// Assign changed values
		app.Name = upd.Name
		app.Enabled = upd.Enabled
		app.Unify = upd.Unify

		if app, err = svc.application.Update(app); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.ApplicationAfterUpdate(upd, app))
		return nil
	})

	return app, svc.recordAction(svc.ctx, aaProps, ApplicationActionUpdate, err)
}

func (svc *application) Delete(ID uint64) (err error) {
	var (
		aaProps = &applicationActionProps{}
		app     *types.Application
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.application.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteApplication(svc.ctx, app) {
			return ApplicationErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.ApplicationBeforeDelete(nil, app)); err != nil {
			return
		}

		if err = svc.application.DeleteByID(ID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.ApplicationAfterDelete(nil, app))
		return nil
	})

	return svc.recordAction(svc.ctx, aaProps, ApplicationActionDelete, err)
}

func (svc *application) Undelete(ID uint64) (err error) {
	var (
		aaProps = &applicationActionProps{}
		app     *types.Application
	)

	err = svc.db.Transaction(func() (err error) {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.application.FindByID(ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteApplication(svc.ctx, app) {
			return ApplicationErrNotAllowedToUndelete()
		}

		// @todo add event
		//       if err = svc.eventbus.WaitFor(svc.ctx, event.ApplicationBeforeUndelete(nil, app)); err != nil {
		//       	return
		//       }

		if err = svc.application.UndeleteByID(ID); err != nil {
			return
		}

		// @todo add event
		//       _ = svc.eventbus.WaitFor(svc.ctx, event.ApplicationAfterUndelete(nil, app))
		return nil
	})

	return svc.recordAction(svc.ctx, aaProps, ApplicationActionDelete, err)
}
