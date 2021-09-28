package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newReportFromResource(res *resource.Report, cfg *EncoderConfig) resourceState {
	return &report{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *report) Prepare(ctx context.Context, pl *payload) (err error) {
	// Reset old identifiers
	n.res.Res.ID = 0

	// Try to get the original report
	n.rp, err = findReportStore(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.rp != nil {
		n.res.Res.ID = n.rp.ID
	}
	return nil
}

func (n *report) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.rp != nil && n.rp.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.rp.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Sys users
	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	res.CreatedBy = pl.invokerID
	if us != nil {
		if us.OwnedBy != nil {
			res.OwnedBy = us.OwnedBy.UserID
		}
		if us.CreatedBy != nil {
			res.CreatedBy = us.CreatedBy.UserID
		}
		if us.UpdatedBy != nil {
			res.UpdatedBy = us.UpdatedBy.UserID
		}
		if us.DeletedBy != nil {
			res.DeletedBy = us.DeletedBy.UserID
		}
	}

	res.Sources = make(types.ReportDataSourceSet, 0, 10)
	for _, rp := range n.res.Sources {
		res.Sources = append(res.Sources, rp.Res)
	}

	res.Blocks = make(types.ReportBlockSet, 0, 10)
	for _, rp := range n.res.Blocks {
		res.Blocks = append(res.Blocks, rp.Res)
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to automation/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh report
	if !exists {
		return store.CreateReport(ctx, pl.s, res)
	}

	// Update existing report
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeReports(n.rp, res)

	case resource.MergeRight:
		res = mergeReports(res, n.rp)
	}

	err = store.UpdateReport(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
