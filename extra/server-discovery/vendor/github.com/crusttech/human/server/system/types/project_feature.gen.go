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
	ProjectFeature struct {
		ID               uint64     `json:"featureID,string"`
		TenantID         uint64     `json:"tenantID,string,omitempty"`
		ProjectID        uint64     `json:"projectID,string,omitempty"`
		RevisionID       uint64     `json:"revisionID,string,omitempty"`
		Title            string     `json:"title"`
		Description      string     `json:"description"`
		FeatureType      string     `json:"featureType,omitempty"`
		Status           string     `json:"status"`
		Severity         string     `json:"severity"`
		Risk             string     `json:"risk"`
		FeatureOwner     uint64     `json:"featureOwner,string,omitempty"`
		ChangeOwner      uint64     `json:"changeOwner,string,omitempty"`
		ChangeApprovedBy uint64     `json:"changeApprovedBy,string,omitempty"`
		RiskFeature      string     `json:"riskFeature,omitempty"`
		ChangeRequired   string     `json:"changeRequired,omitempty"`
		RiskChange       string     `json:"riskChange,omitempty"`
		DateDue          string     `json:"dateDue,omitempty"`
		CreatedAt        time.Time  `json:"createdAt,omitempty"`
		UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
		DeletedAt        *time.Time `json:"deletedAt,omitempty"`
		CreatedBy        uint64     `json:"createdBy,string"`
		UpdatedBy        uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy        uint64     `json:"deletedBy,string,omitempty"`
	}
)

func (r ProjectFeature) Clone() *ProjectFeature {
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

func (r ProjectFeature) Diff(cmp *ProjectFeature) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectFeature{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "featureID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.FeatureType != cmp.FeatureType {
		out = append(out, &revisions.Change{Key: "featureType", Old: []any{cmp.FeatureType}, New: []any{r.FeatureType}})
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

	if r.FeatureOwner != cmp.FeatureOwner {
		out = append(out, &revisions.Change{Key: "featureOwner", Old: []any{cmp.FeatureOwner}, New: []any{r.FeatureOwner}})
	}

	if r.ChangeOwner != cmp.ChangeOwner {
		out = append(out, &revisions.Change{Key: "changeOwner", Old: []any{cmp.ChangeOwner}, New: []any{r.ChangeOwner}})
	}

	if r.ChangeApprovedBy != cmp.ChangeApprovedBy {
		out = append(out, &revisions.Change{Key: "changeApprovedBy", Old: []any{cmp.ChangeApprovedBy}, New: []any{r.ChangeApprovedBy}})
	}

	if r.RiskFeature != cmp.RiskFeature {
		out = append(out, &revisions.Change{Key: "riskFeature", Old: []any{cmp.RiskFeature}, New: []any{r.RiskFeature}})
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
