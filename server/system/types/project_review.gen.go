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
	ProjectReview struct {
		ID              uint64     `json:"reviewID,string"`
		TenantID        uint64     `json:"tenantID,string,omitempty"`
		ProjectID       uint64     `json:"projectID,string,omitempty"`
		RevisionID      uint64     `json:"revisionID,string,omitempty"`
		Title           string     `json:"title"`
		Description     string     `json:"description"`
		ReviewType      string     `json:"reviewType,omitempty"`
		ReviewFrequency string     `json:"reviewFrequency,omitempty"`
		Scope           string     `json:"scope"`
		Reviewer        uint64     `json:"reviewer,string,omitempty"`
		ApprovedBy      uint64     `json:"approvedBy,string,omitempty"`
		Status          string     `json:"status"`
		DateDue         string     `json:"dateDue,omitempty"`
		CreatedAt       time.Time  `json:"createdAt,omitempty"`
		UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
		DeletedAt       *time.Time `json:"deletedAt,omitempty"`
		CreatedBy       uint64     `json:"createdBy,string"`
		UpdatedBy       uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy       uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r ProjectReview) Clone() *ProjectReview {
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

func (r ProjectReview) Diff(cmp *ProjectReview) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectReview{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "reviewID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.RevisionID != cmp.RevisionID {
		out = append(out, &revisions.Change{Key: "revisionID", Old: []any{cmp.RevisionID}, New: []any{r.RevisionID}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.ReviewType != cmp.ReviewType {
		out = append(out, &revisions.Change{Key: "reviewType", Old: []any{cmp.ReviewType}, New: []any{r.ReviewType}})
	}

	if r.ReviewFrequency != cmp.ReviewFrequency {
		out = append(out, &revisions.Change{Key: "reviewFrequency", Old: []any{cmp.ReviewFrequency}, New: []any{r.ReviewFrequency}})
	}

	if r.Scope != cmp.Scope {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
	}

	if r.Reviewer != cmp.Reviewer {
		out = append(out, &revisions.Change{Key: "reviewer", Old: []any{cmp.Reviewer}, New: []any{r.Reviewer}})
	}

	if r.ApprovedBy != cmp.ApprovedBy {
		out = append(out, &revisions.Change{Key: "approvedBy", Old: []any{cmp.ApprovedBy}, New: []any{r.ApprovedBy}})
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
