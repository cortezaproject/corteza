package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	report struct {
		cfg *EncoderConfig

		res *resource.Report
		rp  *types.Report
		ss  types.ReportDataSourceSet
		pp  types.ReportProjectionSet

		ux *userIndex
	}
	reportSet []*report

	reportSource struct {
		cfg *EncoderConfig

		res *resource.ReportSource
		tr  *types.ReportDataSource
	}
	reportSourceSet []*reportSource
)

// mergeReports merges b into a, prioritising a
func mergeReports(a, b *types.Report) *types.Report {
	c := a

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Meta == nil {
		c.Meta = b.Meta
	}
	if c.Sources == nil {
		c.Sources = b.Sources
	}
	if c.Projections == nil {
		c.Projections = b.Projections
	}

	if c.OwnedBy == 0 {
		c.OwnedBy = b.OwnedBy
	}
	if c.CreatedBy == 0 {
		c.CreatedBy = b.CreatedBy
	}
	if c.UpdatedBy == 0 {
		c.UpdatedBy = b.UpdatedBy
	}
	if c.DeletedBy == 0 {
		c.DeletedBy = b.DeletedBy
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}

	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}

	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return c
}

// findReport looks for the report in the resources & the store
//
// Provided resources are prioritized.
func findReport(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (wf *types.Report, err error) {
	wf = resource.FindReport(rr, ii)
	if wf != nil {
		return wf, nil
	}

	return findReportStore(ctx, s, makeGenericFilter(ii))
}

// findReportStore looks for the report in the store
func findReportStore(ctx context.Context, s store.Storer, gf genericFilter) (wf *types.Report, err error) {
	if gf.id > 0 {
		wf, err = store.LookupReportByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		wf, err = store.LookupReportByHandle(ctx, s, i)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	return nil, nil
}
