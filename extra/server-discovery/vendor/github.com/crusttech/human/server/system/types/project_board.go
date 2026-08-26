package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

// Project board DTOs.
//
// The kanban board (BoardPanel.vue) needs every work item across all six
// resources (the five event categories plus backlog items) bucketed into the
// four shared status columns, with a TRUE per-column total and only the
// first page of cards — not the whole project's rowset. Today the frontend
// loads each of the six resources capped at 200 rows (stores/events.js,
// stores/backlogItems.js) and groups client-side; past 200 rows per
// resource, cards silently stop rendering and the counts go wrong. This
// endpoint replaces that for the board specifically.
//
// Unlike the report endpoint (project_report.go), which pages every matching
// row into memory to aggregate, this endpoint only ever fetches as many rows
// per resource as a column page needs — see service/project_board.go's
// loadColumn.
//
// The five category lists and the backlog list (CategoryPanel.vue,
// BacklogView.vue) are NOT served by this endpoint: they already read full,
// type-specific rows (every field the edit dialog/drawer need, not just the
// board's slim card projection), and the existing per-resource list
// endpoints (e.g. GET /project-incidents/) already support cursor paging and
// a true total via `incTotal` — the fix there is rewiring the frontend off
// its 200-row cap onto that existing capability, not new backend surface.
type (
	ProjectBoardRequest struct {
		// Chain-root project scope — required.
		ProjectID uint64

		// Optional revision scope. Zero (default) is chain-wide, including
		// unassigned (rel_revision = 0) rows; non-zero scopes every
		// underlying query to that one revision, excluding unassigned rows —
		// same contract as ProjectReportRequest.RevisionID.
		RevisionID uint64

		// Status scopes the response to a single column. Required for
		// paging (PageCursor is meaningless without it); unset, every column
		// is returned — the board's initial load.
		Status string

		// PageCursor continues a single column's paging from where a
		// previous page (for the same Status) left off; opaque, round-tripped
		// from that page's Column.NextPage. Ignored/invalid when Status is
		// unset — a single cursor can't resume four independent columns.
		PageCursor string

		// Cards per column page; the service applies a default when zero.
		Limit uint
	}

	// ProjectBoardItem is the card shape BoardCard.vue renders, normalized
	// across all six source resources so the frontend never special-cases a
	// field name per category (mirrors, and replaces, BoardPanel.vue's own
	// client-side OWNER_KEY normalization) — see service/project_board.go's
	// projectBoardSources registry for the per-resource mapping.
	ProjectBoardItem struct {
		ID       uint64 `json:"id,string"`
		ItemType string `json:"itemType"`

		// Set only for backlog items — the category of the event they're
		// linked to (drives the card's type badge; mirrors
		// BoardCard.vue's isBacklog branch).
		LinkedCategory string `json:"linkedCategory,omitempty"`

		Title  string `json:"title"`
		Status string `json:"status"`

		// Severity is blank for review and backlog items; Priority is set
		// only for backlog items — mirrors the underlying resource shapes
		// exactly rather than overloading one field to mean different
		// things for different item types.
		Severity string `json:"severity,omitempty"`
		Priority string `json:"priority,omitempty"`

		// Raw owner/assignee user reference (0 = unassigned), normalized off
		// whichever per-category field carries it (issueOwner, featureOwner,
		// requestOwner, owner, reviewer, or backlog's assignee). Left as an
		// ID rather than resolved to a display name: name resolution is a
		// project-user-directory concern the frontend already owns
		// (stores/users.js) and is orthogonal to the board-loading problem
		// this endpoint fixes.
		Owner uint64 `json:"owner,string,omitempty"`

		DueDate    string `json:"dueDate,omitempty"`
		RevisionID uint64 `json:"revisionID,string,omitempty"`
	}

	// ProjectBoardColumn is one status column. Total is the TRUE count across
	// all six resources for this project/revision/status; set only on a
	// fresh load (PageCursor unset) — a continuation page omits it, since the
	// store's IncTotal can't be combined with a page cursor (see
	// service/project_board.go's loadColumn) and the frontend already has
	// the total from the column's first page.
	ProjectBoardColumn struct {
		Status   string               `json:"status"`
		Total    uint64               `json:"total,omitempty"`
		Items    []*ProjectBoardItem  `json:"items"`
		NextPage *filter.PagingCursor `json:"nextPage,omitempty"`
	}

	ProjectBoardResult struct {
		Columns []*ProjectBoardColumn `json:"columns"`
	}
)
