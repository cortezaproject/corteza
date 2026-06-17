package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/label"
)

func (svc *chart) FindByID(ctx context.Context, namespaceID uint64, ID uint64) (res *types.Chart, err error) {
	var (
		aProps = &chartActionProps{chart: &types.Chart{ID: ID}}
	)

	err = func() error {
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ChartActionCreate, err)
}

func (svc *chart) Update(ctx context.Context, upd *types.Chart) (res *types.Chart, err error) {
	var (
		aProps = &chartActionProps{changed: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ChartActionUpdate, err)
}

func (svc *chart) DeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &chartActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, namespaceID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ChartActionDelete, err)
}

func (svc *chart) UndeleteByID(ctx context.Context, namespaceID uint64, ID uint64) (err error) {
	var (
		aProps = &chartActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, namespaceID, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ChartActionUndelete, err)
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
