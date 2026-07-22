package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
)

type chart struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        chartAccessController
	services  *chartServices
}

func (svc *chart) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.Chart, err error) {
	var (
		aProps = &chartActionProps{chart: &types.Chart{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, namespaceID, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ChartActionLookup, err)
}

func (svc *chart) Search(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error) {
	var (
		aProps = &chartActionProps{filter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ChartActionSearch, err)
}

func (svc *chart) Create(ctx context.Context, new *types.Chart) (res *types.Chart, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &chartActionProps{chart: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ChartActionCreate, err)
}

func (svc *chart) Update(ctx context.Context, upd *types.Chart) (res *types.Chart, err error) {
	var (
		aProps = &chartActionProps{changed: upd}
		old    *types.Chart
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChart(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChart(res)
		aProps.setChanged(res)
		old = res.Clone()

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ChartErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ChartErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.NamespaceID = upd.NamespaceID
		res.Name = upd.Name
		res.UpdatedAt = now()

		if err = store.UpdateComposeChart(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ChartActionUpdate, err, old, res)
}

func (svc *chart) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &chartActionProps{}
		res    *types.Chart
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChart(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChart(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onDelete(ctx, s, namespaceID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ChartActionDelete, err)
}

func (svc *chart) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &chartActionProps{}
		res    *types.Chart
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChart(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChart(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onUndelete(ctx, s, namespaceID, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ChartActionUndelete, err)
}

func loadChart(ctx context.Context, s store.ComposeCharts, ID uint64) (res *types.Chart, err error) {
	if ID == 0 {
		return nil, ChartErrInvalidID()
	}

	if res, err = store.LookupComposeChartByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ChartErrNotFound()
	}

	return
}

// toLabeledCharts converts to []label.LabeledResource
func toLabeledCharts(set []*types.Chart) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *chart) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
