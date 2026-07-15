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
		incident projectReportIncidentSearcher
		task     projectReportTaskSearcher
		feature  projectReportFeatureSearcher
		privacy  projectReportPrivacySearcher
		review   projectReportReviewSearcher
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

	// reportSample is the neutral projected form the aggregator groups over:
	// the created-at timestamp (for the "day" dimension and windowing) plus the
	// pre-extracted column dimension values.
	reportSample struct {
		createdAt time.Time
		dims      map[string]any
	}

	// projectReportSource declares the dimensions/metrics a category supports
	// and how to load its project-scoped rows as samples.
	projectReportSource struct {
		dimensions map[string]bool
		metrics    map[string]bool
		load       func(ctx context.Context, svc *projectReport, projectID uint64) ([]reportSample, error)
	}
)

// v1 metrics are the same for every category.
var projectReportMetrics = map[string]bool{"count": true}

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
					return reportSample{createdAt: r.CreatedAt, dims: map[string]any{
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
					return reportSample{createdAt: r.CreatedAt, dims: map[string]any{
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
					return reportSample{createdAt: r.CreatedAt, dims: map[string]any{
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
					return reportSample{createdAt: r.CreatedAt, dims: map[string]any{
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
					return reportSample{createdAt: r.CreatedAt, dims: map[string]any{
						"status": r.Status, "type": r.ReviewType, "scope": r.Scope,
					}}
				},
			)
		},
	},
}

func ProjectReport() *projectReport {
	return &projectReport{
		incident: DefaultProjectIncident,
		task:     DefaultProjectTask,
		feature:  DefaultProjectFeature,
		privacy:  DefaultProjectPrivacy,
		review:   DefaultProjectReview,
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
		dims  map[string]any
		count float64
	}

	var (
		buckets = map[string]*bucket{}
		order   = make([]string, 0)
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
