package provision

import (
	"context"
	"fmt"
	cmpTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"strings"
)

// Migrates resource translations from the resource
// struct to the dedicated store (table)
//
// While doing this, we also modify some resource substructure:
//  - chart reports (assign report IDs)
//
// Note: we will migrate all translations to current default language
// If you do not like that, shut down Corteza after migrations and fix this directly in the store
func migratePost202203ResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("migrating post 202203 resource locales")

	var (
		migrated = make(map[string]bool)
	)

	set, _, err := store.SearchResourceTranslations(ctx, s, sysTypes.ResourceTranslationFilter{})
	set.Walk(func(r *sysTypes.ResourceTranslation) error {
		var pos = strings.Index(r.Resource, "/")
		if pos < 0 {
			return nil
		}

		migrated[r.Resource[0:pos]] = true
		return nil
	})

	if !migrated[cmpTypes.ChartResourceTranslationType] {
		if err = migrateComposeChartResourceTranslations(ctx, log, s); err != nil {
			return
		}
	}

	return
}

// migrate resource translations for compose chart
func migrateComposeChartResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) error {
	var (
		tt sysTypes.ResourceTranslationSet
	)
	set, _, err := store.SearchComposeCharts(ctx, s, cmpTypes.ChartFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return err
	}

	log.Info("migrating compose charts", zap.Int("count", len(set)))

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(res *cmpTypes.Chart) (err error) {
			if tt, err = convertComposeChartReportsTranslations(res); err != nil {
				return err
			}

			if err = store.CreateResourceTranslation(ctx, s, tt...); err != nil {
				return
			}

			if err = store.UpdateComposeChart(ctx, s, res); err != nil {
				return
			}

			return
		})
	})
}

// collects translations for compose chart blocks and alters them (adding block ID)
func convertComposeChartReportsTranslations(res *cmpTypes.Chart) (sysTypes.ResourceTranslationSet, error) {
	var tt sysTypes.ResourceTranslationSet

	for i, r := range res.Config.Reports {
		reportID := id.Next()
		res.Config.Reports[i].ReportID = reportID

		if _, ok := r.YAxis["label"]; ok {
			tt = append(tt,
				makeResourceTranslation(res, "yAxis.label", r.YAxis["label"].(string)),
			)
		}

		// Ensure chart report metric IDs
		for j, m := range r.Metrics {
			metricID := id.Next()
			res.Config.Reports[i].Metrics[j]["metricID"] = fmt.Sprintf("%d", metricID)
			cmx := fmt.Sprintf("metrics.%d.", metricID)

			if _, ok := m["label"]; ok {
				tt = append(tt,
					makeResourceTranslation(res, cmx+"label", m["label"].(string)),
				)
			}
		}
	}

	return tt, nil
}
