package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	ProjectFriaScenario struct {
		ID         uint64                  `json:"projectFriaScenarioID,string"`
		TenantID   uint64                  `json:"tenantID,string,omitempty"`
		ProjectID  uint64                  `json:"projectID,string"`
		AiSystemID uint64                  `json:"aiSystemID,string"`
		Title      string                  `json:"title"`
		Severity   string                  `json:"severity"`
		Meta       ProjectFriaScenarioMeta `json:"meta"`
		CreatedAt  time.Time               `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time              `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time              `json:"deletedAt,omitempty"`
		CreatedBy  uint64                  `json:"createdBy,string"`
		UpdatedBy  uint64                  `json:"updatedBy,string,omitempty"`
		DeletedBy  uint64                  `json:"deletedBy,string,omitempty"`
	}

	ProjectFriaScenarioMeta struct {
		Description            string   `json:"description,omitempty"`
		TriggerTypes           []string `json:"triggerTypes,omitempty"`
		TriggerDescription     string   `json:"triggerDescription,omitempty"`
		ImpactedParties        []string `json:"impactedParties,omitempty"`
		VulnerableGroups       []string `json:"vulnerableGroups,omitempty"`
		VulnerableGroupsNotes  string   `json:"vulnerableGroupsNotes,omitempty"`
		Rights                 []string `json:"rights,omitempty"`
		HarmVectors            []string `json:"harmVectors,omitempty"`
		HarmVectorsDescription string   `json:"harmVectorsDescription,omitempty"`
	}
)

func (r ProjectFriaScenario) Clone() *ProjectFriaScenario {
	dup := r
	dup.Meta = *r.Meta.Clone()

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

func (r ProjectFriaScenario) Diff(cmp *ProjectFriaScenario) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectFriaScenario{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "projectFriaScenarioID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.AiSystemID != cmp.AiSystemID {
		out = append(out, &revisions.Change{Key: "aiSystemID", Old: []any{cmp.AiSystemID}, New: []any{r.AiSystemID}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Severity != cmp.Severity {
		out = append(out, &revisions.Change{Key: "severity", Old: []any{cmp.Severity}, New: []any{r.Severity}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
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

func (r ProjectFriaScenarioMeta) Clone() *ProjectFriaScenarioMeta {
	dup := r
	if r.TriggerTypes != nil {
		dup.TriggerTypes = make([]string, len(r.TriggerTypes))
		copy(dup.TriggerTypes, r.TriggerTypes)
	}

	if r.ImpactedParties != nil {
		dup.ImpactedParties = make([]string, len(r.ImpactedParties))
		copy(dup.ImpactedParties, r.ImpactedParties)
	}

	if r.VulnerableGroups != nil {
		dup.VulnerableGroups = make([]string, len(r.VulnerableGroups))
		copy(dup.VulnerableGroups, r.VulnerableGroups)
	}

	if r.Rights != nil {
		dup.Rights = make([]string, len(r.Rights))
		copy(dup.Rights, r.Rights)
	}

	if r.HarmVectors != nil {
		dup.HarmVectors = make([]string, len(r.HarmVectors))
		copy(dup.HarmVectors, r.HarmVectors)
	}

	return &dup
}

func (r ProjectFriaScenarioMeta) Diff(cmp *ProjectFriaScenarioMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectFriaScenarioMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.TriggerTypes, cmp.TriggerTypes) {
		out = append(out, &revisions.Change{Key: "triggerTypes", Old: []any{cmp.TriggerTypes}, New: []any{r.TriggerTypes}})
	}

	if r.TriggerDescription != cmp.TriggerDescription {
		out = append(out, &revisions.Change{Key: "triggerDescription", Old: []any{cmp.TriggerDescription}, New: []any{r.TriggerDescription}})
	}

	if !reflect.DeepEqual(r.ImpactedParties, cmp.ImpactedParties) {
		out = append(out, &revisions.Change{Key: "impactedParties", Old: []any{cmp.ImpactedParties}, New: []any{r.ImpactedParties}})
	}

	if !reflect.DeepEqual(r.VulnerableGroups, cmp.VulnerableGroups) {
		out = append(out, &revisions.Change{Key: "vulnerableGroups", Old: []any{cmp.VulnerableGroups}, New: []any{r.VulnerableGroups}})
	}

	if r.VulnerableGroupsNotes != cmp.VulnerableGroupsNotes {
		out = append(out, &revisions.Change{Key: "vulnerableGroupsNotes", Old: []any{cmp.VulnerableGroupsNotes}, New: []any{r.VulnerableGroupsNotes}})
	}

	if !reflect.DeepEqual(r.Rights, cmp.Rights) {
		out = append(out, &revisions.Change{Key: "rights", Old: []any{cmp.Rights}, New: []any{r.Rights}})
	}

	if !reflect.DeepEqual(r.HarmVectors, cmp.HarmVectors) {
		out = append(out, &revisions.Change{Key: "harmVectors", Old: []any{cmp.HarmVectors}, New: []any{r.HarmVectors}})
	}

	if r.HarmVectorsDescription != cmp.HarmVectorsDescription {
		out = append(out, &revisions.Change{Key: "harmVectorsDescription", Old: []any{cmp.HarmVectorsDescription}, New: []any{r.HarmVectorsDescription}})
	}

	return out
}

func (r *ProjectFriaScenarioMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectFriaScenarioMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseProjectFriaScenarioMeta(ss []string) (p ProjectFriaScenarioMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
