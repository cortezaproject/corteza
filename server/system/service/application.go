package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/flag"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (LookupByID, Create, Update, Delete, Undelete, loadApplication
// and toLabeledApplications) is generated in application.gen.go from
// system/application.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// before-create / before-update hooks the generated Create / Update call into,
// and the custom bodies the generated wrappers delegate to: onSearch (flag
// filtering + enrichment), onFlag / onUnflag (via checkFlag) and onReorder.

type (
	application struct {
		ac        applicationAccessController
		eventbus  eventDispatcher
		actionlog actionlog.Recorder
		store     store.Storer
	}

	applicationAccessController interface {
		CanCreateApplication(context.Context) bool
		CanSearchApplications(context.Context) bool
		CanSelfApplicationFlag(context.Context) bool
		CanGlobalApplicationFlag(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
	}
)

// Application is a default application service initializer
func Application(s store.Storer, ac applicationAccessController, al actionlog.Recorder, eb eventDispatcher) *application {
	return &application{store: s, ac: ac, actionlog: al, eventbus: eb}
}

// beforeCreate runs after the access check and before the generated Create
// assigns the ID / timestamps and persists. It defaults the Unify config.
func (svc *application) beforeCreate(ctx context.Context, new *types.Application) error {
	if new.Unify == nil {
		new.Unify = &types.ApplicationUnify{}
	}

	return nil
}

// beforeUpdate runs after the stale-data guard and before the generated Update
// copies the mutable fields onto the loaded record. Unify is merged here because
// the original copies it only when provided (it must not be nulled out otherwise).
func (svc *application) beforeUpdate(ctx context.Context, upd, existing *types.Application) error {
	if upd.Unify != nil {
		existing.Unify = upd.Unify
	}

	return nil
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction + the standard CanSearchApplications
// check; the flag filtering / enrichment (flag.Search / flag.Load), which the
// template does not emit, lives here.
func (svc *application) onSearch(ctx context.Context, af types.ApplicationFilter, aaProps *applicationActionProps) (aa types.ApplicationSet, f types.ApplicationFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	af.Check = func(res *types.Application) (bool, error) {
		if !svc.ac.CanReadApplication(ctx, res) {
			return false, nil
		}

		return true, nil
	}

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
			return
		}

		// labels specified but no labeled resources found
		if len(af.LabeledIDs) == 0 {
			return
		}
	}

	if len(af.Flags) > 0 {
		af.FlaggedIDs, err = flag.Search(
			ctx,
			svc.store,
			a.GetIdentityFromContext(ctx).Identity(),
			(&types.Application{}).FlagResourceKind(),
			af.Flags...,
		)

		if err != nil {
			return
		}

		// flags specified byt no flagged resources found
		if len(af.FlaggedIDs) == 0 {
			return
		}
	}

	if aa, f, err = store.SearchApplications(ctx, svc.store, af); err != nil {
		return
	}

	if err = label.Load(ctx, svc.store, toLabeledApplications(aa)...); err != nil {
		return
	}

	if err = flag.Load(ctx, svc.store, f.IncFlags, a.GetIdentityFromContext(ctx).Identity(), toFlaggedApplications(aa)...); err != nil {
		return
	}

	return
}

func (svc *application) onFlag(ctx context.Context, _ *applicationActionProps, app *types.Application, ownedBy uint64, f string) (err error) {
	if err = svc.checkFlag(ctx, ownedBy); err != nil {
		return err
	}

	return flag.Create(ctx, svc.store, app, ownedBy, f)
}

func (svc *application) onUnflag(ctx context.Context, _ *applicationActionProps, app *types.Application, ownedBy uint64, f string) (err error) {
	if err = svc.checkFlag(ctx, ownedBy); err != nil {
		return err
	}

	return flag.Delete(ctx, svc.store, app, ownedBy, f)
}

func (svc *application) checkFlag(ctx context.Context, ownedBy uint64) error {
	if ownedBy == 0 {
		if !svc.ac.CanGlobalApplicationFlag(ctx) {
			return ApplicationErrNotAllowedToManageFlagGlobal()
		}
	} else {
		if ownedBy != a.GetIdentityFromContext(ctx).Identity() || !svc.ac.CanSelfApplicationFlag(ctx) {
			return ApplicationErrNotAllowedToManageFlag()
		}
	}

	return nil
}

func (svc *application) onReorder(ctx context.Context, aProps *applicationActionProps, order []uint64) (err error) {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		for _, id := range order {
			// This access control creates an aux application so we don't have to fetch them
			// from the store; the ID is the only thing that matters...
			auxApp := &types.Application{
				ID: id,
			}

			if !svc.ac.CanUpdateApplication(ctx, auxApp) {
				aProps.application = auxApp
				return ApplicationErrNotAllowedToUpdate(aProps)
			}
		}

		return store.ReorderApplications(ctx, s, order)
	})
}

// toFlaggedApplications converts to []flag.FlaggedResource
//
// This function is auto-generated.
func toFlaggedApplications(set []*types.Application) []flag.FlaggedResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]flag.FlaggedResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
