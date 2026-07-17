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
	ProjectIncident struct {
		ID               uint64     `json:"incidentID,string"`
		TenantID         uint64     `json:"tenantID,string,omitempty"`
		ProjectID        uint64     `json:"projectID,string,omitempty"`
		Title            string     `json:"title"`
		Description      string     `json:"description"`
		IncidentType     string     `json:"incidentType,omitempty"`
		GroupSystem      string     `json:"groupSystem,omitempty"`
		Status           string     `json:"status"`
		Severity         string     `json:"severity"`
		Risk             string     `json:"risk"`
		IssueOwner       uint64     `json:"issueOwner,string,omitempty"`
		ChangeOwner      uint64     `json:"changeOwner,string,omitempty"`
		ChangeApprovedBy uint64     `json:"changeApprovedBy,string,omitempty"`
		RiskIssue        string     `json:"riskIssue,omitempty"`
		ChangeRequired   string     `json:"changeRequired,omitempty"`
		RiskChange       string     `json:"riskChange,omitempty"`
		DateDue          string     `json:"dateDue,omitempty"`
		CompletedDate    string     `json:"completedDate,omitempty"`
		CreatedAt        time.Time  `json:"createdAt,omitempty"`
		UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
		DeletedAt        *time.Time `json:"deletedAt,omitempty"`
		CreatedBy        uint64     `json:"createdBy,string"`
		UpdatedBy        uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy        uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r ProjectIncident) Clone() *ProjectIncident {
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

func (r ProjectIncident) Diff(cmp *ProjectIncident) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectIncident{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "incidentID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.IncidentType != cmp.IncidentType {
		out = append(out, &revisions.Change{Key: "incidentType", Old: []any{cmp.IncidentType}, New: []any{r.IncidentType}})
	}

	if r.GroupSystem != cmp.GroupSystem {
		out = append(out, &revisions.Change{Key: "groupSystem", Old: []any{cmp.GroupSystem}, New: []any{r.GroupSystem}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Severity != cmp.Severity {
		out = append(out, &revisions.Change{Key: "severity", Old: []any{cmp.Severity}, New: []any{r.Severity}})
	}

	if r.Risk != cmp.Risk {
		out = append(out, &revisions.Change{Key: "risk", Old: []any{cmp.Risk}, New: []any{r.Risk}})
	}

	if r.IssueOwner != cmp.IssueOwner {
		out = append(out, &revisions.Change{Key: "issueOwner", Old: []any{cmp.IssueOwner}, New: []any{r.IssueOwner}})
	}

	if r.ChangeOwner != cmp.ChangeOwner {
		out = append(out, &revisions.Change{Key: "changeOwner", Old: []any{cmp.ChangeOwner}, New: []any{r.ChangeOwner}})
	}

	if r.ChangeApprovedBy != cmp.ChangeApprovedBy {
		out = append(out, &revisions.Change{Key: "changeApprovedBy", Old: []any{cmp.ChangeApprovedBy}, New: []any{r.ChangeApprovedBy}})
	}

	if r.RiskIssue != cmp.RiskIssue {
		out = append(out, &revisions.Change{Key: "riskIssue", Old: []any{cmp.RiskIssue}, New: []any{r.RiskIssue}})
	}

	if r.ChangeRequired != cmp.ChangeRequired {
		out = append(out, &revisions.Change{Key: "changeRequired", Old: []any{cmp.ChangeRequired}, New: []any{r.ChangeRequired}})
	}

	if r.RiskChange != cmp.RiskChange {
		out = append(out, &revisions.Change{Key: "riskChange", Old: []any{cmp.RiskChange}, New: []any{r.RiskChange}})
	}

	if r.DateDue != cmp.DateDue {
		out = append(out, &revisions.Change{Key: "dateDue", Old: []any{cmp.DateDue}, New: []any{r.DateDue}})
	}

	if r.CompletedDate != cmp.CompletedDate {
		out = append(out, &revisions.Change{Key: "completedDate", Old: []any{cmp.CompletedDate}, New: []any{r.CompletedDate}})
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
