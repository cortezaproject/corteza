package actionlog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReportRequestNormalize(t *testing.T) {
	t.Run("defaults metrics to count", func(t *testing.T) {
		rr := ReportRequest{Dimensions: []string{"resource"}}
		require.NoError(t, rr.Normalize())
		require.Equal(t, []string{"count"}, rr.Metrics)
	})

	t.Run("dedupes and drops empty keys", func(t *testing.T) {
		rr := ReportRequest{
			Dimensions: []string{"resource", "", "resource", "action"},
			Metrics:    []string{"count", "count", "errors"},
		}
		require.NoError(t, rr.Normalize())
		require.Equal(t, []string{"resource", "action"}, rr.Dimensions)
		require.Equal(t, []string{"count", "errors"}, rr.Metrics)
	})

	t.Run("rejects unknown dimension", func(t *testing.T) {
		rr := ReportRequest{Dimensions: []string{"nope"}}
		require.ErrorContains(t, rr.Normalize(), `unknown report dimension "nope"`)
	})

	t.Run("rejects unknown metric", func(t *testing.T) {
		rr := ReportRequest{Metrics: []string{"nope"}}
		require.ErrorContains(t, rr.Normalize(), `unknown report metric "nope"`)
	})

	t.Run("allows empty dimensions for global totals", func(t *testing.T) {
		rr := ReportRequest{}
		require.NoError(t, rr.Normalize())
		require.Empty(t, rr.Dimensions)
	})
}

func TestReportRegistryColumnsResolvable(t *testing.T) {
	for key, d := range reportDimensions {
		require.Equal(t, key, d.Key, "dimension registry key mismatch")
		require.NotEmpty(t, d.Column, "dimension %q needs a column", key)
	}

	for key, m := range reportMetrics {
		require.Equal(t, key, m.Key, "metric registry key mismatch")
		if m.Kind != ReportMetricCount {
			require.NotEmpty(t, m.Column, "metric %q needs a column", key)
		}
	}
}
