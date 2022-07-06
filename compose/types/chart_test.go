package types

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChart_decodeTranslations(t *testing.T) {
	cc := []struct {
		name string
		base *ChartConfigReport
		ccr  *ChartConfigReport
		tt   locale.ResourceTranslationIndex
	}{
		{"empty", &ChartConfigReport{}, &ChartConfigReport{}, nil},
		{
			"XAxis label",
			&ChartConfigReport{
				YAxis: map[string]interface{}{"label": ""},
			},
			&ChartConfigReport{
				YAxis: map[string]interface{}{"label": "new label"},
			},
			locale.ResourceTranslationIndex{
				"yAxis.label": &locale.ResourceTranslation{Msg: "new label"},
			},
		},
		{
			"Metric labels",
			&ChartConfigReport{
				Metrics: []map[string]interface{}{
					{"metricID": "112233"},
				},

				Dimensions: []map[string]interface{}{
					{
						"dimensionID": "223344",
						"meta": map[string]interface{}{
							"steps": []map[string]interface{}{
								{"stepID": "2233441"},
								{"stepID": "2233442"},
							},
						},
					},
					{
						"dimensionID": "443322",
						"meta": map[string]interface{}{
							"steps": []map[string]interface{}{
								{"stepID": "4433221"},
								{"stepID": "4433222"},
							},
						},
					},
				},
			},
			&ChartConfigReport{
				Metrics: []map[string]interface{}{
					{"metricID": "112233", "label": "metric label"},
				},
				Dimensions: []map[string]interface{}{
					{
						"dimensionID": "223344",
						"meta": map[string]interface{}{
							"steps": []map[string]interface{}{
								{"stepID": "2233441", "label": "Step label 1.1"},
								{"stepID": "2233442", "label": "Step label 1.2"},
							},
						},
					},
					{
						"dimensionID": "443322",
						"meta": map[string]interface{}{
							"steps": []map[string]interface{}{
								{"stepID": "4433221", "label": "Step label 2.1"},
								{"stepID": "4433222", "label": "Step label 2.2"},
							},
						},
					},
				},
			},
			locale.ResourceTranslationIndex{
				"metrics.112233.label":                       &locale.ResourceTranslation{Msg: "metric label"},
				"dimensions.223344.meta.steps.2233441.label": &locale.ResourceTranslation{Msg: "Step label 1.1"},
				"dimensions.223344.meta.steps.2233442.label": &locale.ResourceTranslation{Msg: "Step label 1.2"},
				"dimensions.443322.meta.steps.4433221.label": &locale.ResourceTranslation{Msg: "Step label 2.1"},
				"dimensions.443322.meta.steps.4433222.label": &locale.ResourceTranslation{Msg: "Step label 2.2"},
			},
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req   = require.New(t)
				chart = &Chart{Config: ChartConfig{Reports: []*ChartConfigReport{c.base}}}
			)

			chart.decodeTranslations(c.tt)
			req.Equal(c.ccr, chart.Config.Reports[0])
		})
	}
}

func TestChart_encodeTranslations(t *testing.T) {
	cc := []struct {
		name    string
		payload string
		tt      locale.ResourceTranslationSet
	}{
		{"empty", "{}", locale.ResourceTranslationSet{}},
		{
			"filled",
			`{"reports": [{
						  "YAxis": { "label": "YAxis label" },
						  "reportID": "291579520866123964",
						  "filter": "YEAR(created_at) = YEAR(NOW()) AND QUARTER(created_at) = QUARTER(NOW())",
						  "moduleID": "285374676287488188",
						  "metrics": [
							{
							  "label": "metric label",
							  "metricID": "112233"
							},
							{
							  "metricID": "223344"
							}
						  ],
						  "dimensions": [{
							  "conditions": {},
							  "field": "Status",
							  "dimensionID": "11223344",
							  "meta": {
								"steps": [
								  { "stepID": "1111", "label": "aa", "value": "23" },
								  { "stepID": "2222", "label": "bb", "value": "25" }
								]
							  },
							  "modifier": "(no grouping / buckets)"
						  }]}]}`,
			locale.ResourceTranslationSet{
				{Resource: "compose:chart/0/0", Key: "yAxis.label", Msg: "YAxis label"},
				{Resource: "compose:chart/0/0", Key: "metrics.112233.label", Msg: "metric label"},
				{Resource: "compose:chart/0/0", Key: "metrics.223344.label", Msg: ""},
				{Resource: "compose:chart/0/0", Key: "dimensions.11223344.meta.steps.1111.label", Msg: "aa"},
				{Resource: "compose:chart/0/0", Key: "dimensions.11223344.meta.steps.2222.label", Msg: "bb"},
			},
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req   = require.New(t)
				chart = &Chart{Config: ChartConfig{}}
			)

			req.NoError(json.Unmarshal([]byte(c.payload), &chart.Config))
			result := chart.encodeTranslations()
			req.Equal(c.tt, result)
		})
	}
}
