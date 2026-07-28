package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/types"
)

// projectBoard serves the kanban board (BoardPanel.vue): the four shared
// status columns (Open / In Progress / Ready to Test / Completed), each
// carrying every work item across all six resources (the five event
// categories plus backlog items) that currently sits in that status.
//
// Reuses the exact same per-resource Search access (RBAC + project scope
// already enforced there — see project_report.go's own doc comment) and the
// same searcher interfaces project_report.go declares
// (projectReportIncidentSearcher etc.), rather than a parallel registry.
//
// Where this deliberately differs from project_report.go: that endpoint
// pages every matching row into memory to aggregate (fine for counts, where
// every row must be visited once). A board column only ever needs its next
// page of cards plus a true total, so loadColumn below never reads more rows
// per resource than a page needs.
type projectBoard struct {
	incident    projectReportIncidentSearcher
	task        projectReportTaskSearcher
	feature     projectReportFeatureSearcher
	privacy     projectReportPrivacySearcher
	review      projectReportReviewSearcher
	backlogItem projectReportBacklogItemSearcher
}

func ProjectBoard() *projectBoard {
	return &projectBoard{
		incident:    DefaultProjectIncident,
		task:        DefaultProjectTask,
		feature:     DefaultProjectFeature,
		privacy:     DefaultProjectPrivacy,
		review:      DefaultProjectReview,
		backlogItem: DefaultProjectBacklogItem,
	}
}

const projectBoardDefaultLimit = 20

// projectBoardTotalProbeLimit bounds the per-source fetch used solely to
// learn a column's true total (see loadColumn) — fixed, independent of the
// caller's requested item-page limit. See loadColumn's doc comment for why
// that decoupling matters.
const projectBoardTotalProbeLimit = 20

// projectBoardStatuses is the board's fixed column set — mirrors the
// frontend's EVENT_STATUS (client/web/unify/src/sections/project/config/
// eventForm.js), which every work item's free-text Status field is expected
// to hold one of. Keep the two lists in sync by hand; same cross-reference
// idiom as projectReportCompletedStatus in project_report.go.
var projectBoardStatuses = []string{"Open", "In Progress", "Ready to Test", "Completed"}

// projectBoardSources declares, in a fixed order, how to page one resource's
// rows for a single (project, revision, status) scope and map them to the
// board's normalized card shape. The order is CATEGORY_ORDER (config/
// categories.js) plus backlog last — the same concatenation order
// BoardPanel.vue's own boardItems computed already produces today
// ([...evItems, ...blItems], with evItems itself in CATS' key order:
// incident, feature, privacy, task, review). A column fills from these
// sources in order (see loadColumn): all of source N's items in a status
// before any of source N+1's. This is a deliberate simplification over a
// true cross-type chronological merge — cheap to page and to resume
// correctly (see boardCursorArb), at the cost of a column reading as
// "all incidents, then all features, ..." rather than strictly newest-first
// across types. Acceptable for v1; revisit if that ordering reads oddly once
// the board is under real load.
type projectBoardSource struct {
	key string
	// search runs one bounded Search call against the resource for the given
	// scope/status/paging and returns the mapped cards plus that call's
	// NextPage/Total (Total is only meaningful when paging.IncTotal was set).
	search func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error)
}

var projectBoardSources = []projectBoardSource{
	{
		key: "incident",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.incident.Search(ctx, types.ProjectIncidentFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "incident", Title: r.Title, Status: r.Status,
					Severity: r.Severity, Owner: r.IssueOwner, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
	{
		key: "feature",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.feature.Search(ctx, types.ProjectFeatureFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "feature", Title: r.Title, Status: r.Status,
					Severity: r.Severity, Owner: r.FeatureOwner, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
	{
		key: "privacy",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.privacy.Search(ctx, types.ProjectPrivacyFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "privacy", Title: r.Title, Status: r.Status,
					Severity: r.Severity, Owner: r.RequestOwner, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
	{
		key: "task",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.task.Search(ctx, types.ProjectTaskFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "task", Title: r.Title, Status: r.Status,
					Severity: r.Severity, Owner: r.Owner, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
	{
		key: "review",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.review.Search(ctx, types.ProjectReviewFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				// review has no severity — left blank, same as BoardPanel's
				// current OWNER_KEY/severity handling for this category.
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "review", Title: r.Title, Status: r.Status,
					Owner: r.Reviewer, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
	{
		key: "backlog",
		search: func(ctx context.Context, svc *projectBoard, projectID, revisionID uint64, status string, paging filter.Paging) ([]*types.ProjectBoardItem, *filter.PagingCursor, uint, error) {
			set, f, err := svc.backlogItem.Search(ctx, types.ProjectBacklogItemFilter{
				ProjectID:  projectID,
				RevisionID: revisionID,
				Status:     status,
				Sorting:    boardSorting(),
				Paging:     paging,
			})
			if err != nil {
				return nil, nil, 0, err
			}
			items := make([]*types.ProjectBoardItem, len(set))
			for i, r := range set {
				items[i] = &types.ProjectBoardItem{
					ID: r.ID, ItemType: "backlog", LinkedCategory: r.Category, Title: r.Title, Status: r.Status,
					Priority: r.Priority, Owner: r.Assignee, DueDate: r.DateDue, RevisionID: r.RevisionID,
				}
			}
			return items, f.NextPage, f.Total, nil
		},
	},
}

// boardSorting gives every source's paging a stable, cursor-compatible order:
// newest first by primary key. Fixed and identical across all six resources
// (rather than e.g. dateDue, a free-text/unvalidated field — see
// project_report.go's reportDueDateFormats) so a single cursor shape works
// for every source in projectBoardSources.
func boardSorting() filter.Sorting {
	s, _ := filter.NewSorting("id DESC")
	return s
}

// boardCursorArb is the extra bookkeeping stashed in a column cursor's Arb
// (pkg/filter/pagination.go's PagingCursor.Arb — "arbitrary per-consumer
// cursor state") alongside whichever source's own seek cursor is active:
// which entry of projectBoardSources to resume from. Same "cursor + side-
// channel state" idiom service/connection.go's catalogArb uses.
type boardCursorArb struct {
	ResourceIdx int `json:"i"`
}

// Board validates the request and returns either every column (Status
// unset — the board's initial load) or a single column's next page
// (Status + PageCursor set — a column's "load more").
func (svc *projectBoard) Board(ctx context.Context, rr *types.ProjectBoardRequest) (*types.ProjectBoardResult, error) {
	if rr.ProjectID == 0 {
		return nil, errors.InvalidData("board requires a project scope")
	}

	if rr.Status == "" && rr.PageCursor != "" {
		return nil, errors.InvalidData("a page cursor requires a single status column")
	}

	limit := rr.Limit
	if limit == 0 {
		limit = projectBoardDefaultLimit
	}

	statuses := projectBoardStatuses
	if rr.Status != "" {
		statuses = []string{rr.Status}
	}

	var cursor *filter.PagingCursor
	if rr.PageCursor != "" {
		p, err := filter.NewPaging(limit, rr.PageCursor)
		if err != nil {
			return nil, errors.InvalidData("invalid page cursor: %v", err)
		}
		cursor = p.PageCursor
	}

	columns := make([]*types.ProjectBoardColumn, 0, len(statuses))
	for _, status := range statuses {
		col, err := svc.loadColumn(ctx, rr.ProjectID, rr.RevisionID, status, cursor, limit)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return &types.ProjectBoardResult{Columns: columns}, nil
}

// loadColumn fetches one status column's true total (fresh loads only —
// see below) and its next page of cards, resuming from cursor when set.
//
// Total: only computed when cursor is nil (a fresh load, never a
// continuation) — the store forbids combining IncTotal with a page cursor
// (see rdbms SearchProjectIncidents's own check), and the frontend already
// has the total from the column's first page, so a "load more" doesn't need
// it recomputed. Costs one bounded Search call per source (6 calls per
// column, each capped at projectBoardTotalProbeLimit rows regardless of the
// caller's requested item-page `limit` — deliberately decoupled, so a caller
// that only wants totals (e.g. RevisionCompletenessBar.vue, which never
// renders a single card) can pass a small `limit` without shrinking the
// probe and thereby making the store's own IncTotal fallback trigger more
// often. A source's own IncTotal can still fall back to an unbounded count
// query when that one source/status bucket alone holds more than the probe
// size (store behaviour, not introduced here — the same fallback every
// existing `incTotal` list param already relies on).
//
// Items: walks projectBoardSources in order starting at cursor's resource
// (or the first source, fresh), pulling exactly as many rows from each as
// still fit in the page — never a whole source's rowset. Stops the moment
// the page is full, tagging the returned cursor with which source to resume
// from next time (boardCursorArb) so a later "load more" call picks up
// mid-source or at the next source boundary without re-reading anything
// already returned.
func (svc *projectBoard) loadColumn(ctx context.Context, projectID, revisionID uint64, status string, cursor *filter.PagingCursor, limit uint) (*types.ProjectBoardColumn, error) {
	col := &types.ProjectBoardColumn{Status: status}

	if cursor == nil {
		var total uint64
		for _, src := range projectBoardSources {
			_, _, t, err := src.search(ctx, svc, projectID, revisionID, status, filter.Paging{Limit: projectBoardTotalProbeLimit, IncTotal: true})
			if err != nil {
				return nil, err
			}
			total += uint64(t)
		}
		col.Total = total
	}

	startIdx, innerCursor, err := decodeBoardCursor(cursor)
	if err != nil {
		return nil, err
	}

	var items []*types.ProjectBoardItem
	idx := startIdx
	for idx < len(projectBoardSources) && uint(len(items)) < limit {
		src := projectBoardSources[idx]
		remaining := limit - uint(len(items))

		page, next, _, err := src.search(ctx, svc, projectID, revisionID, status, filter.Paging{Limit: remaining, PageCursor: innerCursor})
		if err != nil {
			return nil, err
		}
		items = append(items, page...)
		// Only the resumed source (idx == startIdx) ever carries an inbound
		// cursor; every source visited after it in this same call starts fresh.
		innerCursor = nil

		if next != nil {
			// This source isn't exhausted — resume here next time. By
			// construction (the store only sets NextPage when a full page of
			// `remaining` rows came back), items is now exactly at `limit`.
			_ = next.SetArb(boardCursorArb{ResourceIdx: idx})
			col.NextPage = next
			break
		}
		idx++
	}

	if col.NextPage == nil && idx < len(projectBoardSources) {
		// The page filled exactly at a source boundary — more sources remain
		// but none has been queried yet, so there's no real seek cursor to
		// carry, only which source to resume at.
		next := &filter.PagingCursor{}
		_ = next.SetArb(boardCursorArb{ResourceIdx: idx})
		col.NextPage = next
	}

	col.Items = items
	return col, nil
}

// decodeBoardCursor reads a column cursor back into (which source to resume,
// that source's own seek cursor or nil to start it fresh). A nil cursor
// (the column's first page) resumes at source 0 with no seek state.
func decodeBoardCursor(cursor *filter.PagingCursor) (idx int, inner *filter.PagingCursor, err error) {
	if cursor == nil {
		return 0, nil, nil
	}

	var arb boardCursorArb
	if err = cursor.GetArb(&arb); err != nil {
		return 0, nil, errors.InvalidData("invalid page cursor: %v", err)
	}
	if arb.ResourceIdx < 0 || arb.ResourceIdx >= len(projectBoardSources) {
		return 0, nil, errors.InvalidData("invalid page cursor")
	}

	if len(cursor.Keys()) == 0 {
		// A source-boundary marker (see loadColumn) — no seek state for this
		// source yet, start it fresh.
		return arb.ResourceIdx, nil, nil
	}
	return arb.ResourceIdx, cursor, nil
}
