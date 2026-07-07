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
	ApigwFilter struct {
		ID        uint64            `json:"filterID,string"`
		TenantID  uint64            `json:"tenantID,string,omitempty"`
		ProjectID uint64            `json:"projectID,string,omitempty"`
		Route     uint64            `json:"routeID,string"`
		Weight    uint64            `json:"weight,string"`
		Kind      string            `json:"kind,omitempty"`
		Ref       string            `json:"ref,omitempty"`
		Enabled   bool              `json:"enabled,omitempty"`
		Params    ApigwFilterParams `json:"params"`
		CreatedAt time.Time         `json:"createdAt,omitempty"`
		UpdatedAt *time.Time        `json:"updatedAt,omitempty"`
		DeletedAt *time.Time        `json:"deletedAt,omitempty"`
		CreatedBy uint64            `json:"createdBy,string"`
		UpdatedBy uint64            `json:"updatedBy,string,omitempty"`
		DeletedBy uint64            `json:"deletedBy,string,omitempty"`
	}
)

func (r ApigwFilter) Clone() *ApigwFilter {
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

func (r ApigwFilter) Diff(cmp *ApigwFilter) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApigwFilter{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "filterID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Route != cmp.Route {
		out = append(out, &revisions.Change{Key: "routeID", Old: []any{cmp.Route}, New: []any{r.Route}})
	}

	if r.Weight != cmp.Weight {
		out = append(out, &revisions.Change{Key: "weight", Old: []any{cmp.Weight}, New: []any{r.Weight}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if r.Ref != cmp.Ref {
		out = append(out, &revisions.Change{Key: "ref", Old: []any{cmp.Ref}, New: []any{r.Ref}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
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

func (m *ApigwFilterParams) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ApigwFilterParams) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseApigwFilterParams(ss []string) (p ApigwFilterParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
