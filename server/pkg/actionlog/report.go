package actionlog

import (
	"fmt"
	"sort"
	"strings"
)

type (
	ReportRequest struct {
		Dimensions []string
		Metrics    []string

		Filter
	}

	ReportRow struct {
		Dimensions map[string]any     `json:"dimensions"`
		Metrics    map[string]float64 `json:"metrics"`
	}

	ReportRowSet []*ReportRow

	ReportResult struct {
		Dimensions  []string          `json:"dimensions"`
		Metrics     []string          `json:"metrics"`
		ActorLabels map[string]string `json:"actorLabels,omitempty"`
		Set         ReportRowSet      `json:"set"`
	}

	ReportDimensionKind string
	ReportMetricKind    string

	ReportDimension struct {
		Key    string              `json:"key"`
		Kind   ReportDimensionKind `json:"kind"`
		Column string              `json:"-"`
	}

	ReportMetric struct {
		Key    string           `json:"key"`
		Kind   ReportMetricKind `json:"kind"`
		Column string           `json:"-"`
	}
)

const (
	// group by raw column value
	ReportDimensionColumn ReportDimensionKind = "column"
	// group by an ID reference column; values are serialized as strings
	ReportDimensionRef ReportDimensionKind = "ref"
	// group by a timestamp column truncated to date
	ReportDimensionDate ReportDimensionKind = "date"

	// COUNT(*)
	ReportMetricCount ReportMetricKind = "count"
	// COUNT(DISTINCT column)
	ReportMetricCountDistinct ReportMetricKind = "countDistinct"
	// count of rows where column is a non-empty string
	ReportMetricCountNotEmpty ReportMetricKind = "countNotEmpty"
)

// Report output is defined by these two registries; adding a new dimension or
// metric over an existing kind is a single entry here. New kinds additionally
// need an expression builder in the rdbms store adapter.
var (
	reportDimensions = map[string]ReportDimension{
		"resource": {Key: "resource", Kind: ReportDimensionColumn, Column: "resource"},
		"action":   {Key: "action", Kind: ReportDimensionColumn, Column: "action"},
		"origin":   {Key: "origin", Kind: ReportDimensionColumn, Column: "request_origin"},
		"severity": {Key: "severity", Kind: ReportDimensionColumn, Column: "severity"},
		"actor":    {Key: "actor", Kind: ReportDimensionRef, Column: "actor_id"},
		"day":      {Key: "day", Kind: ReportDimensionDate, Column: "ts"},
	}

	reportMetrics = map[string]ReportMetric{
		"count":  {Key: "count", Kind: ReportMetricCount},
		"actors": {Key: "actors", Kind: ReportMetricCountDistinct, Column: "actor_id"},
		"errors": {Key: "errors", Kind: ReportMetricCountNotEmpty, Column: "error"},
	}
)

func ReportDimensionByKey(key string) (ReportDimension, bool) {
	d, ok := reportDimensions[key]
	return d, ok
}

func ReportMetricByKey(key string) (ReportMetric, bool) {
	m, ok := reportMetrics[key]
	return m, ok
}

// Normalize dedupes dimension/metric keys, applies the default metric and
// rejects unknown keys so lower layers can trust the request.
func (rr *ReportRequest) Normalize() error {
	rr.Dimensions = dedupeKeys(rr.Dimensions)
	for _, key := range rr.Dimensions {
		if _, ok := reportDimensions[key]; !ok {
			return fmt.Errorf("unknown report dimension %q (available: %s)", key, availableKeys(reportDimensions))
		}
	}

	if len(rr.Metrics) == 0 {
		rr.Metrics = []string{"count"}
	}

	rr.Metrics = dedupeKeys(rr.Metrics)
	for _, key := range rr.Metrics {
		if _, ok := reportMetrics[key]; !ok {
			return fmt.Errorf("unknown report metric %q (available: %s)", key, availableKeys(reportMetrics))
		}
	}

	return nil
}

func dedupeKeys(kk []string) []string {
	var (
		seen = make(map[string]bool, len(kk))
		out  = make([]string, 0, len(kk))
	)

	for _, k := range kk {
		if k == "" || seen[k] {
			continue
		}

		seen[k] = true
		out = append(out, k)
	}

	return out
}

func availableKeys[T any](m map[string]T) string {
	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}

	sort.Strings(kk)
	return strings.Join(kk, ", ")
}
