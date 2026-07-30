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
	ProjectAiSystem struct {
		ID        uint64              `json:"projectAiSystemID,string"`
		TenantID  uint64              `json:"tenantID,string,omitempty"`
		ProjectID uint64              `json:"projectID,string"`
		Handle    string              `json:"handle"`
		RiskClass string              `json:"riskClass"`
		Meta      ProjectAiSystemMeta `json:"meta"`
		CreatedAt time.Time           `json:"createdAt,omitempty"`
		UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
		DeletedAt *time.Time          `json:"deletedAt,omitempty"`
	}

	ProjectAiSystemMeta struct {
		Short           string `json:"short"`
		Description     string `json:"description,omitempty"`
		IntendedPurpose string `json:"intendedPurpose,omitempty"`
	}
)

func (r ProjectAiSystem) Clone() *ProjectAiSystem {
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

func (r ProjectAiSystem) Diff(cmp *ProjectAiSystem) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectAiSystem{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "projectAiSystemID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.RiskClass != cmp.RiskClass {
		out = append(out, &revisions.Change{Key: "riskClass", Old: []any{cmp.RiskClass}, New: []any{r.RiskClass}})
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

	return out
}

func (r ProjectAiSystemMeta) Clone() *ProjectAiSystemMeta {
	dup := r
	return &dup
}

func (r ProjectAiSystemMeta) Diff(cmp *ProjectAiSystemMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectAiSystemMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.IntendedPurpose != cmp.IntendedPurpose {
		out = append(out, &revisions.Change{Key: "intendedPurpose", Old: []any{cmp.IntendedPurpose}, New: []any{r.IntendedPurpose}})
	}

	return out
}

func (r *ProjectAiSystemMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ProjectAiSystemMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseProjectAiSystemMeta(ss []string) (p ProjectAiSystemMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
