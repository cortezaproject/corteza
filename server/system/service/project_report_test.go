package service

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type fakeIncidentSearcher struct {
	set types.ProjectIncidentSet
}

func (f fakeIncidentSearcher) Search(_ context.Context, flt types.ProjectIncidentFilter) (types.ProjectIncidentSet, types.ProjectIncidentFilter, error) {
	// single page: NextPage stays nil
	return f.set, flt, nil
}

func TestNormalizeProjectReport(t *testing.T) {
	src := projectReportSources["incident"]

	t.Run("defaults metric to count", func(t *testing.T) {
		rr := &types.ProjectReportRequest{Resource: "incident", Dimensions: []string{"status"}}
		require.NoError(t, normalizeProjectReport(rr, src))
		require.Equal(t, []string{"count"}, rr.Metrics)
	})

	t.Run("dedupes and drops empty keys", func(t *testing.T) {
		rr := &types.ProjectReportRequest{
			Resource:   "incident",
			Dimensions: []string{"status", "", "status", "severity"},
			Metrics:    []string{"count", "count"},
		}
		require.NoError(t, normalizeProjectReport(rr, src))
		require.Equal(t, []string{"status", "severity"}, rr.Dimensions)
		require.Equal(t, []string{"count"}, rr.Metrics)
	})

	t.Run("rejects unknown dimension", func(t *testing.T) {
		rr := &types.ProjectReportRequest{Resource: "incident", Dimensions: []string{"nope"}}
		require.ErrorContains(t, normalizeProjectReport(rr, src), `unknown report dimension "nope"`)
	})

	t.Run("rejects unknown metric", func(t *testing.T) {
		rr := &types.ProjectReportRequest{Resource: "incident", Metrics: []string{"nope"}}
		require.ErrorContains(t, normalizeProjectReport(rr, src), `unknown report metric "nope"`)
	})

	t.Run("review rejects severity", func(t *testing.T) {
		rr := &types.ProjectReportRequest{Resource: "review", Dimensions: []string{"severity"}}
		require.ErrorContains(t, normalizeProjectReport(rr, projectReportSources["review"]), `unknown report dimension "severity"`)
	})
}

func TestAggregateProjectReport(t *testing.T) {
	day1 := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 7, 2, 9, 0, 0, 0, time.UTC)

	samples := []reportSample{
		{createdAt: day1, dims: map[string]any{"status": "Open", "severity": "High"}},
		{createdAt: day1, dims: map[string]any{"status": "Open", "severity": "Low"}},
		{createdAt: day2, dims: map[string]any{"status": "Closed", "severity": "High"}},
	}

	t.Run("group by status", func(t *testing.T) {
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Dimensions: []string{"status"}, Metrics: []string{"count"}})
		require.Len(t, rows, 2)
		// sorted by key: Closed, Open
		require.Equal(t, "Closed", rows[0].Dimensions["status"])
		require.Equal(t, float64(1), rows[0].Metrics["count"])
		require.Equal(t, "Open", rows[1].Dimensions["status"])
		require.Equal(t, float64(2), rows[1].Metrics["count"])
	})

	t.Run("group by day", func(t *testing.T) {
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Dimensions: []string{"day"}, Metrics: []string{"count"}})
		require.Len(t, rows, 2)
		require.Equal(t, "2026-07-01", rows[0].Dimensions["day"])
		require.Equal(t, float64(2), rows[0].Metrics["count"])
		require.Equal(t, "2026-07-02", rows[1].Dimensions["day"])
		require.Equal(t, float64(1), rows[1].Metrics["count"])
	})

	t.Run("empty dimensions yields grand total", func(t *testing.T) {
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Metrics: []string{"count"}})
		require.Len(t, rows, 1)
		require.Empty(t, rows[0].Dimensions)
		require.Equal(t, float64(3), rows[0].Metrics["count"])
	})

	t.Run("windowing on created-at", func(t *testing.T) {
		from := day2
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{
			Dimensions:    []string{"status"},
			Metrics:       []string{"count"},
			FromTimestamp: &from,
		})
		require.Len(t, rows, 1)
		require.Equal(t, "Closed", rows[0].Dimensions["status"])
		require.Equal(t, float64(1), rows[0].Metrics["count"])
	})
}

func TestProjectReportReport(t *testing.T) {
	ctx := context.Background()
	svc := &projectReport{
		incident: fakeIncidentSearcher{set: types.ProjectIncidentSet{
			{Status: "Open", Severity: "High"},
			{Status: "Open", Severity: "Low"},
			{Status: "Closed", Severity: "High"},
		}},
	}

	t.Run("groups incidents by status", func(t *testing.T) {
		res, err := svc.Report(ctx, &types.ProjectReportRequest{
			Resource:   "incident",
			ProjectID:  1,
			Dimensions: []string{"status"},
		})
		require.NoError(t, err)
		require.Equal(t, "incident", res.Resource)
		require.Equal(t, []string{"count"}, res.Metrics)
		require.Len(t, res.Set, 2)
		require.Equal(t, float64(1), res.Set[0].Metrics["count"]) // Closed
		require.Equal(t, float64(2), res.Set[1].Metrics["count"]) // Open
	})

	t.Run("requires project scope", func(t *testing.T) {
		_, err := svc.Report(ctx, &types.ProjectReportRequest{Resource: "incident"})
		require.ErrorContains(t, err, "project scope")
	})

	t.Run("rejects unknown resource", func(t *testing.T) {
		_, err := svc.Report(ctx, &types.ProjectReportRequest{Resource: "nope", ProjectID: 1})
		require.ErrorContains(t, err, "unknown report resource")
	})

	t.Run("rejects inverted time range", func(t *testing.T) {
		from := time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)
		to := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
		_, err := svc.Report(ctx, &types.ProjectReportRequest{
			Resource: "incident", ProjectID: 1, FromTimestamp: &from, ToTimestamp: &to,
		})
		require.ErrorContains(t, err, "inverted")
	})
}
