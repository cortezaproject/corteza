package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	application struct {
		ac        applicationAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Storer
	}

	applicationAccessController interface {
		CanCreateApplication(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
	}
)

// Application is a default application service initializer
func Application(s store.Storer, ac applicationAccessController, al actionlog.Recorder, eb eventDispatcher) *application {
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

		if app, err = store.LookupApplicationByID(ctx, svc.store, ID); err != nil {
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

func (svc *application) Search(ctx context.Context, af types.ApplicationFilter) (aa types.ApplicationSet, f types.ApplicationFilter, err error) {
	var (
		aaProps = &applicationActionProps{filter: &af}
	)

	// For each fetched item, store backend will check if it is valid or not
	af.Check = func(res *types.Application) (bool, error) {
		if !svc.ac.CanReadApplication(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if af.Deleted > filter.StateExcluded {
			// If list with deleted applications is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted applications
			//if !svc.ac.CanAccess(ctx) {
			//	return ApplicationErrNotAllowedToListApplications()
			//}
		}

		if len(af.Labels) > 0 {
			af.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Application{}.LabelResourceKind(),
				af.Labels,
			)

			if err != nil {
				return err
			}
		}

		if aa, f, err = store.SearchApplications(ctx, svc.store, af); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledApplications(aa)...); err != nil {
			return err
		}

		return nil

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
		new.CreatedAt = *now()

		if new.Unify == nil {
			new.Unify = &types.ApplicationUnify{}
		}

		if err = store.CreateApplication(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		app = new

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

		if app, err = store.LookupApplicationByID(ctx, svc.store, upd.ID); err != nil {
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
		app.UpdatedAt = now()

		if upd.Unify != nil {
			app.Unify = upd.Unify
		}

		if err = store.UpdateApplication(ctx, svc.store, app); err != nil {
			return err
		}

		if label.Changed(app.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			app.Labels = upd.Labels
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

		if app, err = store.LookupApplicationByID(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanDeleteApplication(ctx, app) {
			return ApplicationErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(ctx, event.ApplicationBeforeDelete(nil, app)); err != nil {
			return
		}

		app.DeletedAt = now()
		if err = store.UpdateApplication(ctx, svc.store, app); err != nil {
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

		if app, err = store.LookupApplicationByID(ctx, svc.store, ID); err != nil {
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
		if err = store.UpdateApplication(ctx, svc.store, app); err != nil {
			return
		}

		// @todo add event
		//       _ = svc.eventbus.WaitFor(ctx, event.ApplicationAfterUndelete(nil, app))
		return nil
	}()

	return svc.recordAction(ctx, aaProps, ApplicationActionUndelete, err)
}

// toLabeledApplications converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledApplications(set []*types.Application) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
