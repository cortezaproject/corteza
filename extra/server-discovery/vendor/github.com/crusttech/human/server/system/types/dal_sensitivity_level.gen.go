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
	DalSensitivityLevel struct {
		ID        uint64                  `json:"sensitivityLevelID,string"`
		TenantID  uint64                  `json:"tenantID,string,omitempty"`
		ProjectID uint64                  `json:"projectID,string,omitempty"`
		Handle    string                  `json:"handle"`
		Level     int                     `json:"level"`
		Meta      DalSensitivityLevelMeta `json:"meta"`
		Labels    map[string]string       `json:"labels,omitempty"`
		CreatedAt time.Time               `json:"createdAt,omitempty"`
		UpdatedAt *time.Time              `json:"updatedAt,omitempty"`
		DeletedAt *time.Time              `json:"deletedAt,omitempty"`
		CreatedBy uint64                  `json:"createdBy,string"`
		UpdatedBy uint64                  `json:"updatedBy,string,omitempty"`
		DeletedBy uint64                  `json:"deletedBy,string,omitempty"`
	}

	DalSensitivityLevelMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func (r DalSensitivityLevel) Clone() *DalSensitivityLevel {
	dup := r
	dup.Meta = *r.Meta.Clone()

	if r.Labels != nil {
		dup.Labels = make(map[string]string, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}

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

func (r DalSensitivityLevel) Diff(cmp *DalSensitivityLevel) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalSensitivityLevel{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "sensitivityLevelID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Level != cmp.Level {
		out = append(out, &revisions.Change{Key: "level", Old: []any{cmp.Level}, New: []any{r.Level}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
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

func (r DalSensitivityLevelMeta) Clone() *DalSensitivityLevelMeta {
	dup := r
	return &dup
}

func (r DalSensitivityLevelMeta) Diff(cmp *DalSensitivityLevelMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DalSensitivityLevelMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *DalSensitivityLevelMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DalSensitivityLevelMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseDalSensitivityLevelMeta(ss []string) (p DalSensitivityLevelMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
