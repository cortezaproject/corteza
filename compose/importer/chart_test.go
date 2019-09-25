package importer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestChartImport_CastSet(t *testing.T) {

	impFixTester(t, "chart_full_slice", func(t *testing.T, imp *Importer) {
		req := require.New(t)
		req.NotNil(imp.GetChartImporter(ns.Slug))
		req.Len(imp.GetChartImporter(ns.Slug).set, 2)
	})

	impFixTester(t,
		"chart_with_unknown_module",
		errors.New(`unknown module "un_kno_wn" referenced from chart "chart1" report config`))

	// Pre fill with module that imported chart is referring to
	imp.namespaces.Setup(ns)
	imp.GetModuleImporter(ns.Slug).set = types.ModuleSet{{NamespaceID: ns.ID, Handle: "foo"}}

	impFixTester(t, "chart_full", func(t *testing.T, chart *Chart) {
		req := require.New(t)

		req.Len(chart.set, 2)

		req.NotNil(chart.set.FindByHandle("chart1"))
		req.NotNil(chart.set.FindByHandle("chart2"))

		tc := chart.set.FindByHandle("chart1")
		req.Equal(tc.Name, "chart 1")
		req.Equal(tc.Config, types.ChartConfig{Reports: []*types.ChartConfigReport{
			{
				Filter:   "a=b",
				ModuleID: 0,
				Metrics: []map[string]interface{}{
					{
						"backgroundColor": "#e5a83b",
						"beginAtZero":     true,
						"field":           "count",
						"fixTooltips":     true,
						"label":           "Number of leads",
						"type":            "bar",
					},
				},
				Dimensions: []map[string]interface{}{
					{
						"field": "created_at",
					},
				},
			},
		}})
	})
}
