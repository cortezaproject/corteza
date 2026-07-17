package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/types"
)

// projectReport aggregates project-scoped category rows (incident, task,
// feature, privacy, review) into grouped report rows for dashboard and
// per-category charts.
//
// Categories live on the generic codegen store with Search-only access, so v1
// aggregates in-memory: it pages the full project-scoped rowset through each
// category service (which enforces RBAC + project scope for us), projects each
// row to its dimension values and groups+counts. Per-project category volumes
// are small, so this is adequate until a metric needs SQL pushdown.

const projectReportPageSize = 1000

type (
	projectReport struct {
		incident    projectReportIncidentSearcher
		task        projectReportTaskSearcher
		feature     projectReportFeatureSearcher
		privacy     projectReportPrivacySearcher
		review      projectReportReviewSearcher
		backlogItem projectReportBacklogItemSearcher
	}

	projectReportIncidentSearcher interface {
		Search(context.Context, types.ProjectIncidentFilter) (types.ProjectIncidentSet, types.ProjectIncidentFilter, error)
	}
	projectReportTaskSearcher interface {
		Search(context.Context, types.ProjectTaskFilter) (types.ProjectTaskSet, types.ProjectTaskFilter, error)
	}
	projectReportFeatureSearcher interface {
		Search(context.Context, types.ProjectFeatureFilter) (types.ProjectFeatureSet, types.ProjectFeatureFilter, error)
	}
	projectReportPrivacySearcher interface {
		Search(context.Context, types.ProjectPrivacyFilter) (types.ProjectPrivacySet, types.ProjectPrivacyFilter, error)
	}
	projectReportReviewSearcher interface {
		Search(context.Context, types.ProjectReviewFilter) (types.ProjectReviewSet, types.ProjectReviewFilter, error)
	}
	projectReportBacklogItemSearcher interface {
		Search(context.Context, types.ProjectBacklogItemFilter) (types.ProjectBacklogItemSet, types.ProjectBacklogItemFilter, error)
	}

	// reportSample is the neutral projected form the aggregator groups over:
	// the created-at timestamp (for the "day" dimension and windowing), the
	// pre-extracted column dimension values (dims always includes "status",
	// used for the "open"/"overdue" metrics regardless of the requested
	// grouping), and the raw DateDue string used for "overdue".
	reportSample struct {
		createdAt time.Time
		dims      map[string]any
		dateDue   string
	}

	// projectReportSource declares the dimensions/metrics a category supports
	// and how to load its project-scoped rows as samples.
	projectReportSource struct {
		dimensions map[string]bool
		metrics    map[string]bool
		load       func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error)
	}
)

// v1 metrics are the same for every category: all 5 category types
// (incident, task, feature, privacy, review) carry both Status and DateDue,
// so "open" and "overdue" are declared uniformly rather than per-source.
var projectReportMetrics = map[string]bool{"count": true, "open": true, "overdue": true}

// projectReportCompletedStatus is the single terminal status value used to
// derive the "open" metric across every category. Kept here as the one
// server-side source of truth, mirroring the frontend's isOpenStatus in
// client/web/unify/src/sections/project/stores/events.js (status !==
// 'Completed'); if that rule ever needs per-category nuance, this is the
// only place that has to change.
const projectReportCompletedStatus = "Completed"

// isOpenReportStatus reports whether a category row's Status counts as
// "open" for the "open"/"overdue" metrics.
func isOpenReportStatus(status string) bool {
	return status != projectReportCompletedStatus
}

// projectReportNowFn is the report's clock, indirected so tests can stub
// "now" for deterministic "overdue" assertions.
var projectReportNowFn = time.Now

// reportDueDateFormats is the best-effort set of layouts tried against the
// free-form DateDue string field (there is no schema/validation on it — it's
// plain text on all 5 category types). ISO 8601 date-only is tried first
// since that's what the "date" input widgets on the frontend emit
// (client/web/unify/src/sections/project/config/eventForm.js), then RFC3339
// in case a full timestamp lands there, then two locale-formatted layouts
// (day.month.year and month/day/year) mirroring the format list already used
// for free-form date strings elsewhere in the codebase
// (pkg/envoy/resource/util.go toTime). A value that matches none of these —
// including an empty string — is treated as "never overdue", not as an
// error: DateDue is unvalidated user input and this metric is best-effort.
var reportDueDateFormats = []string{
	"2006-01-02",
	time.RFC3339,
	"02.01.2006",
	"01/02/2006",
}

// parseReportDueDate parses DateDue against reportDueDateFormats in order,
// returning ok=false for empty or unparseable values.
func parseReportDueDate(v string) (time.Time, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, false
	}

	for _, f := range reportDueDateFormats {
		if t, err := time.Parse(f, v); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

// isReportDueOverdue reports whether a raw DateDue value parses to a point
// in time before now. Unparseable/empty values are never overdue.
func isReportDueOverdue(dateDue string, now time.Time) bool {
	t, ok := parseReportDueDate(dateDue)
	if !ok {
		return false
	}
	return t.Before(now)
}

// projectReportSources is the registry: one entry per category. Adding a
// dimension is a single map key here plus its extraction in the projector.
var projectReportSources = map[string]projectReportSource{
	"incident": {
		dimensions: map[string]bool{"status": true, "severity": true, "risk": true, "type": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectIncidentSet, *filter.PagingCursor, error) {
					set, f, err := svc.incident.Search(ctx, types.ProjectIncidentFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectIncident) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "severity": r.Severity, "risk": r.Risk, "type": r.IncidentType,
					}}
				},
			)
		},
	},
	"task": {
		dimensions: map[string]bool{"status": true, "severity": true, "risk": true, "type": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectTaskSet, *filter.PagingCursor, error) {
					set, f, err := svc.task.Search(ctx, types.ProjectTaskFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectTask) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "severity": r.Severity, "risk": r.Risk, "type": r.TaskType,
					}}
				},
			)
		},
	},
	"feature": {
		dimensions: map[string]bool{"status": true, "severity": true, "risk": true, "type": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectFeatureSet, *filter.PagingCursor, error) {
					set, f, err := svc.feature.Search(ctx, types.ProjectFeatureFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectFeature) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "severity": r.Severity, "risk": r.Risk, "type": r.FeatureType,
					}}
				},
			)
		},
	},
	"privacy": {
		dimensions: map[string]bool{"status": true, "severity": true, "risk": true, "type": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectPrivacySet, *filter.PagingCursor, error) {
					set, f, err := svc.privacy.Search(ctx, types.ProjectPrivacyFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectPrivacy) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "severity": r.Severity, "risk": r.Risk, "type": r.RequestType,
					}}
				},
			)
		},
	},
	"review": {
		// review has no severity/risk; it carries scope + review type instead.
		dimensions: map[string]bool{"status": true, "type": true, "scope": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectReviewSet, *filter.PagingCursor, error) {
					set, f, err := svc.review.Search(ctx, types.ProjectReviewFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectReview) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "type": r.ReviewType, "scope": r.Scope,
					}}
				},
			)
		},
	},
	"backlog-item": {
		// backlog items are sub-issues linked to a category item; they carry
		// priority + category instead of severity/risk/type.
		dimensions: map[string]bool{"status": true, "priority": true, "category": true, "day": true},
		metrics:    projectReportMetrics,
		load: func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error) {
			return collectReportSamples(
				func(cursor *filter.PagingCursor) (types.ProjectBacklogItemSet, *filter.PagingCursor, error) {
					set, f, err := svc.backlogItem.Search(ctx, types.ProjectBacklogItemFilter{
						ProjectID: projectID,
						Sorting:   reportSorting(),
						Paging:    reportPaging(cursor),
					})
					return set, f.NextPage, err
				},
				func(r *types.ProjectBacklogItem) reportSample {
					return reportSample{createdAt: r.CreatedAt, dateDue: r.DateDue, dims: map[string]any{
						"status": r.Status, "priority": r.Priority, "category": r.Category,
					}}
				},
			)
		},
	},
}

func ProjectReport() *projectReport {
	return &projectReport{
		incident:    DefaultProjectIncident,
		task:        DefaultProjectTask,
		feature:     DefaultProjectFeature,
		privacy:     DefaultProjectPrivacy,
		review:      DefaultProjectReview,
		backlogItem: DefaultProjectBacklogItem,
	}
}

// Report validates the request, loads the project-scoped rows and returns the
// grouped result. An empty Dimensions set yields a single grand-total row.
func (svc *projectReport) Report(ctx context.Context, rr *types.ProjectReportRequest) (*types.ProjectReportResult, error) {
	if rr.ProjectID == 0 {
		return nil, errors.InvalidData("report requires a project scope")
	}

	if rr.FromTimestamp != nil && rr.ToTimestamp != nil && rr.ToTimestamp.Before(*rr.FromTimestamp) {
		return nil, errors.InvalidData("report time range is inverted")
	}

	src, ok := projectReportSources[rr.Resource]
	if !ok {
		return nil, errors.InvalidData("unknown report resource %q (available: %s)", rr.Resource, sortedKeys(projectReportSourceKeys()))
	}

	if err := normalizeProjectReport(rr, src); err != nil {
		return nil, errors.InvalidData("%v", err)
	}

	samples, err := src.load(ctx, svc, rr.ProjectID)
	if err != nil {
		return nil, err
	}

	return &types.ProjectReportResult{
		Resource:   rr.Resource,
		Dimensions: rr.Dimensions,
		Metrics:    rr.Metrics,
		Set:        aggregateProjectReport(samples, rr),
	}, nil
}

// collectReportSamples pages a category's Search to exhaustion, projecting each
// row to a reportSample. Cursor paging needs a stable sort, supplied by the
// caller via reportSorting().
func collectReportSamples[S ~[]E, E any](
	search func(cursor *filter.PagingCursor) (S, *filter.PagingCursor, error),
	project func(E) reportSample,
) ([]reportSample, error) {
	var (
		out    []reportSample
		cursor *filter.PagingCursor
	)

	for {
		set, next, err := search(cursor)
		if err != nil {
			return nil, err
		}

		for _, e := range set {
			out = append(out, project(e))
		}

		if next == nil {
			break
		}
		cursor = next
	}

	return out, nil
}

func reportPaging(cursor *filter.PagingCursor) filter.Paging {
	return filter.Paging{Limit: projectReportPageSize, PageCursor: cursor}
}

// reportSorting gives cursor paging a stable order; id is the primary index on
// every category.
func reportSorting() filter.Sorting {
	s, _ := filter.NewSorting("id")
	return s
}

// aggregateProjectReport groups samples by the requested dimensions and counts
// each group. Windowing on created-at is applied here.
func aggregateProjectReport(samples []reportSample, rr *types.ProjectReportRequest) []*types.ProjectReportRow {
	type bucket struct {
		dims    map[string]any
		count   float64
		open    float64
		overdue float64
	}

	var (
		buckets = map[string]*bucket{}
		order   = make([]string, 0)
		now     = projectReportNowFn()
	)

	for _, s := range samples {
		if rr.FromTimestamp != nil && s.createdAt.Before(*rr.FromTimestamp) {
			continue
		}
		if rr.ToTimestamp != nil && s.createdAt.After(*rr.ToTimestamp) {
			continue
		}

		dims := make(map[string]any, len(rr.Dimensions))
		keyParts := make([]string, len(rr.Dimensions))
		for i, d := range rr.Dimensions {
			v := reportDimValue(s, d)
			dims[d] = v
			keyParts[i] = fmt.Sprintf("%v", v)
		}

		key := strings.Join(keyParts, "\x00")
		b := buckets[key]
		if b == nil {
			b = &bucket{dims: dims}
			buckets[key] = b
			order = append(order, key)
		}
		b.count++

		status, _ := s.dims["status"].(string)
		if isOpenReportStatus(status) {
			b.open++
			if isReportDueOverdue(s.dateDue, now) {
				b.overdue++
			}
		}
	}

	sort.Strings(order)

	rows := make([]*types.ProjectReportRow, 0, len(order))
	for _, key := range order {
		b := buckets[key]

		metrics := make(map[string]float64, len(rr.Metrics))
		for _, m := range rr.Metrics {
			switch m {
			case "count":
				metrics["count"] = b.count
			case "open":
				metrics["open"] = b.open
			case "overdue":
				metrics["overdue"] = b.overdue
			}
		}

		rows = append(rows, &types.ProjectReportRow{Dimensions: b.dims, Metrics: metrics})
	}

	return rows
}

func reportDimValue(s reportSample, dim string) any {
	if dim == "day" {
		return s.createdAt.Format("2006-01-02")
	}
	return s.dims[dim]
}

// normalizeProjectReport dedupes dimension/metric keys, applies the default
// metric and rejects keys the resource does not support.
func normalizeProjectReport(rr *types.ProjectReportRequest, src projectReportSource) error {
	rr.Dimensions = dedupeReportKeys(rr.Dimensions)
	for _, d := range rr.Dimensions {
		if !src.dimensions[d] {
			return fmt.Errorf("unknown report dimension %q for resource %q (available: %s)", d, rr.Resource, sortedKeys(src.dimensions))
		}
	}

	if len(rr.Metrics) == 0 {
		rr.Metrics = []string{"count"}
	}

	rr.Metrics = dedupeReportKeys(rr.Metrics)
	for _, m := range rr.Metrics {
		if !src.metrics[m] {
			return fmt.Errorf("unknown report metric %q (available: %s)", m, sortedKeys(src.metrics))
		}
	}

	return nil
}

func dedupeReportKeys(kk []string) []string {
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

func sortedKeys(m map[string]bool) string {
	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	return strings.Join(kk, ", ")
}

func projectReportSourceKeys() map[string]bool {
	kk := make(map[string]bool, len(projectReportSources))
	for k := range projectReportSources {
		kk[k] = true
	}
	return kk
}
