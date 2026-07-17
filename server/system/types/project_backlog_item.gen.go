package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
	"time"
)

type (
	ProjectBacklogItem struct {
		ID          uint64     `json:"backlogItemID,string"`
		TenantID    uint64     `json:"tenantID,string,omitempty"`
		ProjectID   uint64     `json:"projectID,string,omitempty"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Category    string     `json:"category"`
		EventID     uint64     `json:"eventID,string,omitempty"`
		Assignee    uint64     `json:"assignee,string,omitempty"`
		Priority    string     `json:"priority"`
		Status      string     `json:"status"`
		DateDue     string     `json:"dateDue,omitempty"`
		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty"`
		CreatedBy   uint64     `json:"createdBy,string"`
		UpdatedBy   uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy   uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r ProjectBacklogItem) Clone() *ProjectBacklogItem {
	dup := r
	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r ProjectBacklogItem) Diff(cmp *ProjectBacklogItem) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectBacklogItem{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "backlogItemID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Category != cmp.Category {
		out = append(out, &revisions.Change{Key: "category", Old: []any{cmp.Category}, New: []any{r.Category}})
	}

	if r.EventID != cmp.EventID {
		out = append(out, &revisions.Change{Key: "eventID", Old: []any{cmp.EventID}, New: []any{r.EventID}})
	}

	if r.Assignee != cmp.Assignee {
		out = append(out, &revisions.Change{Key: "assignee", Old: []any{cmp.Assignee}, New: []any{r.Assignee}})
	}

	if r.Priority != cmp.Priority {
		out = append(out, &revisions.Change{Key: "priority", Old: []any{cmp.Priority}, New: []any{r.Priority}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.DateDue != cmp.DateDue {
		out = append(out, &revisions.Change{Key: "dateDue", Old: []any{cmp.DateDue}, New: []any{r.DateDue}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	return out
}
