package provision

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

func migrateReports(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	reports, _, err := store.SearchReports(ctx, s, types.ReportFilter{
		Deleted: filter.StateExcluded,
	})
	if err != nil {
		return
	}

	mustMigrate := false
	for _, r := range reports {
		for _, s := range r.Sources {
			mustMigrate = mustMigrate || (s.Step != nil && s.Step.Group_legacy != nil && s.Step.Link == nil)
		}
		for _, b := range r.Blocks {
			for _, s := range b.Sources {
				mustMigrate = mustMigrate || (s != nil && s.Group_legacy != nil && s.Link == nil)
			}
		}
	}

	if !mustMigrate {
		return
	}

	ds := func(step *types.ReportStep) {
		if step.Join != nil {
			step.Kind = "link"
			step.Link = &types.ReportStepLink{
				Name:          step.Join.Name,
				LocalSource:   step.Join.LocalSource,
				LocalColumn:   step.Join.LocalColumn,
				ForeignSource: step.Join.ForeignSource,
				ForeignColumn: step.Join.ForeignColumn,
				Filter:        step.Join.Filter,
			}
			step.Join = nil
		}
		if step.Group_legacy != nil {
			step.Kind = "aggregate"
			step.Aggregate = &types.ReportStepAggregate{
				Name:    step.Group_legacy.Name,
				Source:  step.Group_legacy.Source,
				Keys:    step.Group_legacy.Keys,
				Columns: step.Group_legacy.Columns,
				Filter:  step.Group_legacy.Filter,
			}
			step.Group_legacy = nil
		}
	}

	for _, r := range reports {
		for _, s := range r.Sources {
			ds(s.Step)
		}
		for _, b := range r.Blocks {
			for _, s := range b.Sources {
				ds(s)
			}
		}
	}

	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		return store.UpdateReport(ctx, s, reports...)
	})
}
