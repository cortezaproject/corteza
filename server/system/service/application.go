package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/flag"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

type (
	applicationServices struct {
		eventbus eventDispatcher
	}

	applicationAccessController interface {
		CanCreateApplication(context.Context) bool
		CanSearchApplications(context.Context) bool
		CanSelfApplicationFlag(context.Context) bool
		CanGlobalApplicationFlag(context.Context) bool
		CanReadApplication(context.Context, *types.Application) bool
		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
		CanManageSourceOnApplication(context.Context, *types.Application) bool
	}
)

// ApplicationKindCustom marks an application whose UI is its own HTML source,
// shown by the app view in a sandbox, rather than a shell section or a link.
const ApplicationKindCustom = "custom"

// ApplicationSourceMaxSize caps the HTML a custom application may hold.
const ApplicationSourceMaxSize = 256 * 1024

// Application is a default application service initializer
func Application(s store.Storer, ac applicationAccessController, al actionlog.Recorder, eb eventDispatcher) *application {
	return &application{
		store:     s,
		ac:        ac,
		actionlog: al,
		services:  &applicationServices{eventbus: eb},
	}
}

func (svc *application) onSearch(ctx context.Context, af types.ApplicationFilter, aProps *applicationActionProps) (aa types.ApplicationSet, f types.ApplicationFilter, err error) {
	af.Check = func(res *types.Application) (bool, error) {
		if !svc.ac.CanReadApplication(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	if af.Deleted > filter.StateExcluded {
		// If list with deleted applications is requested
		// user must have access permissions to system (ie: is admin)
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

	err = flag.Load(ctx, svc.store, f.IncFlags, a.GetIdentityFromContext(ctx).Identity(), toFlaggedApplications(aa)...)
	return
}

func (svc *application) beforeCreate(ctx context.Context, new *types.Application) error {
	if err := svc.services.eventbus.WaitFor(ctx, event.ApplicationBeforeCreate(new, nil)); err != nil {
		return err
	}

	if new.Meta == nil {
		new.Meta = &types.ApplicationMeta{}
	}

	if new.Unify == nil {
		new.Unify = &types.ApplicationUnify{}
	}

	if new.OwnerID == 0 {
		new.OwnerID = a.GetIdentityFromContext(ctx).Identity()
	}

	return nil
}

func (svc *application) beforeUpdate(ctx context.Context, upd, res *types.Application) error {
	if err := svc.services.eventbus.WaitFor(ctx, event.ApplicationBeforeUpdate(upd, res)); err != nil {
		return err
	}

	if upd.Meta != nil {
		res.Meta = upd.Meta
	}

	if upd.Unify != nil {
		res.Unify = upd.Unify
	}

	return nil
}

// onSetSource replaces a custom application's HTML and the meta the app view
// reads instead of the source. The whole document is replaced every time;
// the store keeps no history of it.
func (svc *application) onSetSource(ctx context.Context, aProps *applicationActionProps, app *types.Application, source string, meta *types.ApplicationSourceMeta) (err error) {
	if app == nil || app.ID == 0 {
		return ApplicationErrInvalidID()
	}

	if app, err = loadApplication(ctx, svc.store, app.ID); err != nil {
		return err
	}

	aProps.setApplication(app)

	if !svc.ac.CanManageSourceOnApplication(ctx, app) {
		return ApplicationErrNotAllowedToManageSource()
	}

	if app.Unify == nil || app.Unify.Kind != ApplicationKindCustom {
		return ApplicationErrNotCustom()
	}

	if len(source) > ApplicationSourceMaxSize {
		return ApplicationErrSourceTooLarge()
	}

	if meta == nil {
		meta = &types.ApplicationSourceMeta{}
	}

	sum := sha256.Sum256([]byte(source))
	meta.Hash = hex.EncodeToString(sum[:])
	meta.Size = len(source)
	meta.UpdatedAt = now()
	meta.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	app.Source = source
	app.SourceMeta = meta
	app.UpdatedAt = now()

	return store.UpdateApplication(ctx, svc.store, app)
}

func (svc *application) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		aaProps = &applicationActionProps{}
		app     *types.Application
	)

	err = func() (err error) {
		if app, err = loadApplication(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanDeleteApplication(ctx, app) {
			return ApplicationErrNotAllowedToDelete()
		}

		if err = svc.services.eventbus.WaitFor(ctx, event.ApplicationBeforeDelete(nil, app)); err != nil {
			return
		}

		app.DeletedAt = now()
		if err = store.UpdateApplication(ctx, svc.store, app); err != nil {
			return
		}

		_ = svc.services.eventbus.WaitFor(ctx, event.ApplicationAfterDelete(nil, app))
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
		if app, err = loadApplication(ctx, svc.store, ID); err != nil {
			return
		}

		aaProps.setApplication(app)

		if !svc.ac.CanDeleteApplication(ctx, app) {
			return ApplicationErrNotAllowedToUndelete()
		}

		app.DeletedAt = nil
		if err = store.UpdateApplication(ctx, svc.store, app); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aaProps, ApplicationActionUndelete, err)
}

func (svc *application) onFlag(ctx context.Context, _ *applicationActionProps, app *types.Application, ownedBy uint64, f string) error {
	if err := svc.checkFlag(ctx, ownedBy); err != nil {
		return err
	}

	return flag.Create(ctx, svc.store, app, ownedBy, f)
}

func (svc *application) onUnflag(ctx context.Context, _ *applicationActionProps, app *types.Application, ownedBy uint64, f string) error {
	if err := svc.checkFlag(ctx, ownedBy); err != nil {
		return err
	}

	return flag.Delete(ctx, svc.store, app, ownedBy, f)
}

func (svc *application) onReorder(ctx context.Context, aProps *applicationActionProps, order []uint64) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		for _, id := range order {
			auxApp := &types.Application{ID: id}
			if !svc.ac.CanUpdateApplication(ctx, auxApp) {
				aProps.application = auxApp
				return ApplicationErrNotAllowedToUpdate(aProps)
			}
		}
		return store.ReorderApplications(ctx, s, order)
	})
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
