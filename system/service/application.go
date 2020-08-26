package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	application struct {
		ac        applicationAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Applications
	}

	applicationAccessController interface {
		CanAccess(context.Context) bool
		CanCreateApplication(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool

		FilterReadableApplications(ctx context.Context) *permissions.ResourceFilter
	}
)

// Application is a default application service initializer
func Application(s store.Applications, ac applicationAccessController, al actionlog.Recorder, eb eventDispatcher) *application {
	return &application{store: s, ac: ac, actionlog: al, eventbus: eb}
}

func (svc *application) LookupByID(ctx context.Context, ID uint64) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{application: &types.Application{ID: ID}}
	)

	err = func() error {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.store.LookupApplicationByID(ctx, ID); err != nil {
			return ApplicationErrInvalidID().Wrap(err)
		}

		aaProps.setApplication(app)

		if !svc.ac.CanReadApplication(ctx, app) {
			return ApplicationErrNotAllowedToRead()
		}

		return nil
	}()

	return app, svc.recordAction(ctx, aaProps, ApplicationActionLookup, err)
}

func (svc *application) Search(ctx context.Context, filter types.ApplicationFilter) (aa types.ApplicationSet, f types.ApplicationFilter, err error) {
	var (
		aaProps = &applicationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Application) (bool, error) {
		if !svc.ac.CanReadApplication(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if filter.Deleted > rh.FilterStateExcluded {
			// If list with deleted applications is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted applications
			if !svc.ac.CanAccess(ctx) {
				return ApplicationErrNotAllowedToListApplications()
			}
		}

		aa, f, err = svc.store.SearchApplications(ctx, filter)
		return err
	}()

	return aa, f, svc.recordAction(ctx, aaProps, ApplicationActionSearch, err)
}

func (svc *application) Create(ctx context.Context, new *types.Application) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateApplication(ctx) {
			return ApplicationErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.ApplicationBeforeCreate(new, nil)); err != nil {
			return
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = now()
		if err = svc.store.CreateApplication(ctx, new); err != nil {
			return
		}

		aaProps.setApplication(app)

		_ = svc.eventbus.WaitFor(ctx, event.ApplicationAfterCreate(new, nil))
		return nil
	}()

	return app, svc.recordAction(ctx, aaProps, ApplicationActionCreate, err)
}

func (svc *application) Update(ctx context.Context, upd *types.Application) (app *types.Application, err error) {
	var (
		aaProps = &applicationActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.store.LookupApplicationByID(ctx, upd.ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanUpdateApplication(ctx, app) {
			return ApplicationErrNotAllowedToUpdate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.ApplicationBeforeUpdate(upd, app)); err != nil {
			return
		}

		// Assign changed values after afterUpdate events are emitted
		app.Name = upd.Name
		app.Enabled = upd.Enabled
		app.Unify = upd.Unify
		app.UpdatedAt = nowPtr()

		if err = svc.store.UpdateApplication(ctx, app); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(ctx, event.ApplicationAfterUpdate(upd, app))
		return nil
	}()

	return app, svc.recordAction(ctx, aaProps, ApplicationActionUpdate, err)
}

func (svc *application) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &applicationActionProps{}
		app     *types.Application
	)

	err = func() (err error) {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.store.LookupApplicationByID(ctx, ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanDeleteApplication(ctx, app) {
			return ApplicationErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(ctx, event.ApplicationBeforeDelete(nil, app)); err != nil {
			return
		}

		app.DeletedAt = nowPtr()
		if err = svc.store.PartialApplicationUpdate(ctx, []string{"UpdatedAt"}, app); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.ApplicationAfterDelete(nil, app))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, ApplicationActionDelete, err)
}

func (svc *application) Undelete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &applicationActionProps{}
		app     *types.Application
	)

	err = func() (err error) {
		if ID == 0 {
			return ApplicationErrInvalidID()
		}

		if app, err = svc.store.LookupApplicationByID(ctx, ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanDeleteApplication(ctx, app) {
			return ApplicationErrNotAllowedToUndelete()
		}

		// @todo add event
		//       if err = svc.eventbus.WaitFor(ctx, event.ApplicationBeforeUndelete(nil, app)); err != nil {
		//       	return
		//       }

		app.DeletedAt = nil
		if err = svc.store.PartialApplicationUpdate(ctx, []string{"UpdatedAt"}, app); err != nil {
			return
		}

		// @todo add event
		//       _ = svc.eventbus.WaitFor(ctx, event.ApplicationAfterUndelete(nil, app))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, ApplicationActionUndelete, err)
}
