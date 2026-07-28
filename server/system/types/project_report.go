package types

import (
	"time"
)

// Project category reporting DTOs.
//
// Categories (incident, task, feature, privacy, review) live on the generic
// codegen store and expose only Search. This report reads the project-scoped
// rowset and aggregates it in-memory into grouped rows — enough for the
// dashboard overview and per-category charts without SQL/DAL pushdown.
//
// The dimension/metric registry and row projection live in the service layer
// (service/project_report.go) because they reference the concrete category
// types; this file only carries the transport shapes.
type (
	ProjectReportRequest struct {
		// Category to report on (incident, task, feature, privacy, review).
		Resource string

		// Group rows by one or more dimensions (e.g. status, severity, day).
		Dimensions []string

		// Metrics to compute per group (v1: count, open, overdue).
		Metrics []string

		// Project scope — required; a report never spans projects.
		ProjectID uint64

		// Optional revision scope. Zero (default) aggregates chain-wide,
		// across every revision plus unassigned (rel_revision = 0) rows,
		// same as before this field existed. Non-zero constrains every
		// underlying query to rel_revision = RevisionID, which — same as
		// the existing revisionID list/filter param on the six category
		// resources — excludes unassigned rows rather than folding them in.
		RevisionID uint64

		// Optional created-at window; applied in-memory over each row.
		FromTimestamp *time.Time
		ToTimestamp   *time.Time
	}

	ProjectReportRow struct {
		Dimensions map[string]any     `json:"dimensions"`
		Metrics    map[string]float64 `json:"metrics"`
	}

	ProjectReportResult struct {
		Resource   string              `json:"resource"`
		Dimensions []string            `json:"dimensions"`
		Metrics    []string            `json:"metrics"`
		Set        []*ProjectReportRow `json:"set"`
	}
)
