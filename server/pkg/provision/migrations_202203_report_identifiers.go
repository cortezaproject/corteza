package provision

import (
	"context"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Migrates existing reports to include the newly added resource identifiers
func migrateReportIdentifiers(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("migrating report identifiers")

	set, _, err := store.SearchReports(ctx, s, sysTypes.ReportFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return err
	}

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(r *sysTypes.Report) (err error) {
			r = setIDs(r)
			return store.UpdateReport(ctx, s, r)
		})
	})
}

// duplicate from the report service
func setIDs(r *sysTypes.Report) *sysTypes.Report {
	// scenarios
	for _, s := range r.Scenarios {
		if s.ScenarioID == 0 {
			s.ScenarioID = id.Next()
		}
	}

	// blocks
	for _, b := range r.Blocks {
		if b.BlockID == 0 {
			b.BlockID = id.Next()
		}

		// elements
		for _, elRaw := range b.Elements {
			el, ok := elRaw.(map[string]interface{})
			if !ok {
				continue
			}

			elID, ok := el["elementID"]
			sElID := cast.ToString(elID)
			if sElID != "" && sElID != "0" {
				continue
			}
			if cast.ToUint64(elID) != 0 {
				continue
			}

			el["elementID"] = strconv.FormatUint(id.Next(), 10)
		}
	}

	return r
}
