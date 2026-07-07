package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type applicationServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *application) FindByID(ctx context.Context, ID uint64) (res *types.Application, err error) {
	var (
		aProps = &applicationActionProps{application: &types.Application{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadApplication(ctx, svc.store, ID); err != nil {
			return ApplicationErrInvalidID().Wrap(err)
		}

		aProps.setApplication(res)

		if !svc.ac.CanReadApplication(ctx, res) {
			return ApplicationErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ApplicationActionLookup, err)
}

func (svc *application) Search(ctx context.Context, filter types.ApplicationFilter) (set types.ApplicationSet, f types.ApplicationFilter, err error) {
	var (
		aProps = &applicationActionProps{filter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchApplications(ctx) {
			return ApplicationErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ApplicationActionSearch, err)
}

func (svc *application) Create(ctx context.Context, new *types.Application) (res *types.Application, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &applicationActionProps{application: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateApplication(ctx) {
			return ApplicationErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateApplication(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ApplicationActionCreate, err)
}

func (svc *application) Update(ctx context.Context, upd *types.Application) (res *types.Application, err error) {
	var (
		aProps = &applicationActionProps{update: upd}
		old    *types.Application
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadApplication(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setApplication(res)

		if !svc.ac.CanUpdateApplication(ctx, res) {
			return ApplicationErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ApplicationErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Name = upd.Name
		res.Enabled = upd.Enabled
		res.Weight = upd.Weight
		res.UpdatedAt = now()

		if err = store.UpdateApplication(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ApplicationActionUpdate, err, old, res)
}

func (svc *application) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &applicationActionProps{}
		res    *types.Application
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadApplication(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setApplication(res)

		if !svc.ac.CanDeleteApplication(ctx, res) {
			return ApplicationErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateApplication(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ApplicationActionDelete, err)
}

func (svc *application) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &applicationActionProps{}
		res    *types.Application
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadApplication(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setApplication(res)

		if !svc.ac.CanDeleteApplication(ctx, res) {
			return ApplicationErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		if err = store.UpdateApplication(ctx, svc.store, res); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, ApplicationActionUndelete, err)
}

func loadApplication(ctx context.Context, s store.Applications, ID uint64) (res *types.Application, err error) {
	if ID == 0 {
		return nil, ApplicationErrInvalidID()
	}

	if res, err = store.LookupApplicationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ApplicationErrNotFound()
	}

	return
}

// toLabeledApplications converts to []label.LabeledResource
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

func (svc *application) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *application) scopeServices(ctx context.Context) *applicationServices {
	return &applicationServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *application) Flag(ctx context.Context, app *types.Application, ownedBy uint64, f string) (err error) {
	var (
		aProps = &applicationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onFlag(ctx, aProps, app, ownedBy, f)
		return err
	}()

	return svc.recordAction(ctx, aProps, ApplicationActionFlag, err)
}

func (svc *application) Unflag(ctx context.Context, app *types.Application, ownedBy uint64, f string) (err error) {
	var (
		aProps = &applicationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onUnflag(ctx, aProps, app, ownedBy, f)
		return err
	}()

	return svc.recordAction(ctx, aProps, ApplicationActionUnflag, err)
}

func (svc *application) Reorder(ctx context.Context, order []uint64) (err error) {
	var (
		aProps = &applicationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onReorder(ctx, aProps, order)
		return err
	}()

	return svc.recordAction(ctx, aProps, ApplicationActionReorder, err)
}
