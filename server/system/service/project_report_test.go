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

	// lastFilter records the filter passed on the most recent Search call,
	// so tests can assert on what the service threaded through (e.g. RevisionID).
	lastFilter *types.ProjectIncidentFilter
}

func (f *fakeIncidentSearcher) Search(_ context.Context, flt types.ProjectIncidentFilter) (types.ProjectIncidentSet, types.ProjectIncidentFilter, error) {
	f.lastFilter = &flt
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

	t.Run("backlog-item supports priority/category, rejects severity", func(t *testing.T) {
		rr := &types.ProjectReportRequest{Resource: "backlog-item", Dimensions: []string{"priority", "category"}}
		require.NoError(t, normalizeProjectReport(rr, projectReportSources["backlog-item"]))

		rr = &types.ProjectReportRequest{Resource: "backlog-item", Dimensions: []string{"severity"}}
		require.ErrorContains(t, normalizeProjectReport(rr, projectReportSources["backlog-item"]), `unknown report dimension "severity"`)
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

func TestAggregateProjectReportOpenOverdue(t *testing.T) {
	defer func(orig func() time.Time) { projectReportNowFn = orig }(projectReportNowFn)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	projectReportNowFn = func() time.Time { return now }

	day1 := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)

	t.Run("open counts everything not Completed", func(t *testing.T) {
		samples := []reportSample{
			{createdAt: day1, dims: map[string]any{"status": "Open"}},
			{createdAt: day1, dims: map[string]any{"status": "InProgress"}},
			{createdAt: day1, dims: map[string]any{"status": "Completed"}},
		}
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Metrics: []string{"open"}})
		require.Len(t, rows, 1)
		require.Equal(t, float64(2), rows[0].Metrics["open"])
	})

	t.Run("overdue counts open rows with a past ISO DateDue", func(t *testing.T) {
		samples := []reportSample{
			{createdAt: day1, dateDue: "2020-01-01", dims: map[string]any{"status": "Open"}},
			{createdAt: day1, dateDue: "2099-01-01", dims: map[string]any{"status": "Open"}},
		}
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Metrics: []string{"count", "open", "overdue"}})
		require.Len(t, rows, 1)
		require.Equal(t, float64(2), rows[0].Metrics["count"])
		require.Equal(t, float64(2), rows[0].Metrics["open"])
		require.Equal(t, float64(1), rows[0].Metrics["overdue"])
	})

	t.Run("overdue skips unparseable or empty DateDue", func(t *testing.T) {
		samples := []reportSample{
			{createdAt: day1, dateDue: "", dims: map[string]any{"status": "Open"}},
			{createdAt: day1, dateDue: "not-a-date", dims: map[string]any{"status": "Open"}},
		}
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Metrics: []string{"overdue"}})
		require.Len(t, rows, 1)
		require.Equal(t, float64(0), rows[0].Metrics["overdue"])
	})

	t.Run("overdue not counted when Completed", func(t *testing.T) {
		samples := []reportSample{
			{createdAt: day1, dateDue: "2020-01-01", dims: map[string]any{"status": "Completed"}},
		}
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{Metrics: []string{"open", "overdue"}})
		require.Len(t, rows, 1)
		require.Equal(t, float64(0), rows[0].Metrics["open"])
		require.Equal(t, float64(0), rows[0].Metrics["overdue"])
	})

	t.Run("combined metrics with a dimension", func(t *testing.T) {
		samples := []reportSample{
			{createdAt: day1, dateDue: "2020-01-01", dims: map[string]any{"status": "Open", "severity": "High"}},
			{createdAt: day1, dateDue: "2099-01-01", dims: map[string]any{"status": "Open", "severity": "High"}},
			{createdAt: day1, dateDue: "2020-01-01", dims: map[string]any{"status": "Completed", "severity": "High"}},
			{createdAt: day1, dateDue: "2020-01-01", dims: map[string]any{"status": "Open", "severity": "Low"}},
		}
		rows := aggregateProjectReport(samples, &types.ProjectReportRequest{
			Dimensions: []string{"severity"},
			Metrics:    []string{"count", "open", "overdue"},
		})
		require.Len(t, rows, 2)
		// sorted by key: High, Low
		require.Equal(t, "High", rows[0].Dimensions["severity"])
		require.Equal(t, float64(3), rows[0].Metrics["count"])
		require.Equal(t, float64(2), rows[0].Metrics["open"])
		require.Equal(t, float64(1), rows[0].Metrics["overdue"])
		require.Equal(t, "Low", rows[1].Dimensions["severity"])
		require.Equal(t, float64(1), rows[1].Metrics["count"])
		require.Equal(t, float64(1), rows[1].Metrics["open"])
		require.Equal(t, float64(1), rows[1].Metrics["overdue"])
	})
}

func TestParseReportDueDate(t *testing.T) {
	t.Run("ISO date", func(t *testing.T) {
		tm, ok := parseReportDueDate("2026-07-16")
		require.True(t, ok)
		require.Equal(t, 2026, tm.Year())
		require.Equal(t, time.July, tm.Month())
		require.Equal(t, 16, tm.Day())
	})

	t.Run("RFC3339", func(t *testing.T) {
		tm, ok := parseReportDueDate("2026-07-16T10:00:00Z")
		require.True(t, ok)
		require.Equal(t, 2026, tm.Year())
	})

	t.Run("dotted day.month.year", func(t *testing.T) {
		tm, ok := parseReportDueDate("16.07.2026")
		require.True(t, ok)
		require.Equal(t, 16, tm.Day())
		require.Equal(t, time.July, tm.Month())
	})

	t.Run("slashed month/day/year", func(t *testing.T) {
		tm, ok := parseReportDueDate("07/16/2026")
		require.True(t, ok)
		require.Equal(t, 16, tm.Day())
		require.Equal(t, time.July, tm.Month())
	})

	t.Run("empty is never overdue", func(t *testing.T) {
		_, ok := parseReportDueDate("")
		require.False(t, ok)
	})

	t.Run("unparseable is never overdue", func(t *testing.T) {
		_, ok := parseReportDueDate("garbage")
		require.False(t, ok)
	})
}

func TestProjectReportReport(t *testing.T) {
	ctx := context.Background()
	incidentSearcher := &fakeIncidentSearcher{set: types.ProjectIncidentSet{
		{Status: "Open", Severity: "High"},
		{Status: "Open", Severity: "Low"},
		{Status: "Closed", Severity: "High"},
	}}
	svc := &projectReport{incident: incidentSearcher}

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

	t.Run("zero RevisionID stays chain-wide: no RevisionID constraint on the filter", func(t *testing.T) {
		_, err := svc.Report(ctx, &types.ProjectReportRequest{
			Resource: "incident", ProjectID: 1, Dimensions: []string{"status"},
		})
		require.NoError(t, err)
		require.NotNil(t, incidentSearcher.lastFilter)
		require.Zero(t, incidentSearcher.lastFilter.RevisionID)
	})

	t.Run("non-zero RevisionID is threaded through to the underlying filter", func(t *testing.T) {
		_, err := svc.Report(ctx, &types.ProjectReportRequest{
			Resource: "incident", ProjectID: 1, RevisionID: 42, Dimensions: []string{"status"},
		})
		require.NoError(t, err)
		require.NotNil(t, incidentSearcher.lastFilter)
		require.Equal(t, uint64(42), incidentSearcher.lastFilter.RevisionID)
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

	t.Run("rejects unknown metric with open/overdue listed as available", func(t *testing.T) {
		_, err := svc.Report(ctx, &types.ProjectReportRequest{Resource: "incident", ProjectID: 1, Metrics: []string{"nope"}})
		require.ErrorContains(t, err, "unknown report metric")
		require.ErrorContains(t, err, "open")
		require.ErrorContains(t, err, "overdue")
	})

	t.Run("combined count/open/overdue metrics grouped by status", func(t *testing.T) {
		defer func(orig func() time.Time) { projectReportNowFn = orig }(projectReportNowFn)
		projectReportNowFn = func() time.Time { return time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC) }

		combinedSvc := &projectReport{
			incident: &fakeIncidentSearcher{set: types.ProjectIncidentSet{
				{Status: "Open", Severity: "High", DateDue: "2020-01-01"},      // open, overdue
				{Status: "Open", Severity: "Low", DateDue: "2099-01-01"},       // open, not overdue
				{Status: "Completed", Severity: "High", DateDue: "2020-01-01"}, // closed, never overdue
			}},
		}

		res, err := combinedSvc.Report(ctx, &types.ProjectReportRequest{
			Resource:   "incident",
			ProjectID:  1,
			Dimensions: []string{"status"},
			Metrics:    []string{"count", "open", "overdue"},
		})
		require.NoError(t, err)
		require.Len(t, res.Set, 2)

		// sorted by key: Completed, Open
		require.Equal(t, "Completed", res.Set[0].Dimensions["status"])
		require.Equal(t, float64(1), res.Set[0].Metrics["count"])
		require.Equal(t, float64(0), res.Set[0].Metrics["open"])
		require.Equal(t, float64(0), res.Set[0].Metrics["overdue"])

		require.Equal(t, "Open", res.Set[1].Dimensions["status"])
		require.Equal(t, float64(2), res.Set[1].Metrics["count"])
		require.Equal(t, float64(2), res.Set[1].Metrics["open"])
		require.Equal(t, float64(1), res.Set[1].Metrics["overdue"])
	})
}
