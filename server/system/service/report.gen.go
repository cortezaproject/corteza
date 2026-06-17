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
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *report) FindByID(ctx context.Context, ID uint64) (res *types.Report, err error) {
	var (
		aProps = &reportActionProps{report: &types.Report{ID: ID}}
	)

	err = func() error {
		if res, err = loadReport(ctx, svc.store, ID); err != nil {
			return ReportErrInvalidID().Wrap(err)
		}

		aProps.setReport(res)

		if !svc.ac.CanReadReport(ctx, res) {
			return ReportErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ReportActionLookup, err)
}

func (svc *report) Search(ctx context.Context, filter types.ReportFilter) (set types.ReportSet, f types.ReportFilter, err error) {
	var (
		aProps = &reportActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Report) (bool, error) {
		if !svc.ac.CanReadReport(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchReports(ctx) {
			return ReportErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Report{}.LabelResourceKind(),
				filter.Labels,
			)
			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchReports(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledReports(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ReportActionSearch, err)
}

func (svc *report) Create(ctx context.Context, new *types.Report) (res *types.Report, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &reportActionProps{report: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateReport(ctx) {
			return ReportErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateReport(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ReportActionCreate, err)
}

func (svc *report) Update(ctx context.Context, upd *types.Report) (res *types.Report, err error) {
	var (
		aProps = &reportActionProps{update: upd}
	)

	err = func() (err error) {
		if res, err = loadReport(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aProps.setReport(res)

		if !svc.ac.CanUpdateReport(ctx, res) {
			return ReportErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ReportErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Handle = upd.Handle
		res.Meta = upd.Meta
		res.Scenarios = upd.Scenarios
		res.Sources = upd.Sources
		res.Blocks = upd.Blocks
		res.UpdatedAt = now()

		if err = store.UpdateReport(ctx, svc.store, res); err != nil {
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

	return res, svc.recordAction(ctx, aProps, ReportActionUpdate, err)
}

func (svc *report) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reportActionProps{}
		res    *types.Report
	)

	err = func() (err error) {
		if res, err = loadReport(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setReport(res)

		if !svc.ac.CanDeleteReport(ctx, res) {
			return ReportErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateReport(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ReportActionDelete, err)
}

func (svc *report) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &reportActionProps{}
		res    *types.Report
	)

	err = func() (err error) {
		if res, err = loadReport(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setReport(res)

		if !svc.ac.CanDeleteReport(ctx, res) {
			return ReportErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		if err = store.UpdateReport(ctx, svc.store, res); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, ReportActionUndelete, err)
}

func loadReport(ctx context.Context, s store.Reports, ID uint64) (res *types.Report, err error) {
	if ID == 0 {
		return nil, ReportErrInvalidID()
	}

	if res, err = store.LookupReportByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ReportErrNotFound()
	}

	return
}

// toLabeledReports converts to []label.LabeledResource
func toLabeledReports(set []*types.Report) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
