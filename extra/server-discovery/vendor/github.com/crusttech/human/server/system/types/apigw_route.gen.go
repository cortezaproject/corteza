package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	ApigwRoute struct {
		ID        uint64                           `json:"routeID,string"`
		TenantID  uint64                           `json:"tenantID,string,omitempty"`
		ProjectID uint64                           `json:"projectID,string,omitempty"`
		Endpoint  string                           `json:"endpoint"`
		Method    string                           `json:"method"`
		Enabled   bool                             `json:"enabled"`
		Meta      ApigwRouteMeta                   `json:"meta"`
		Group     uint64                           `json:"group,string"`
		CreatedAt time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy uint64                           `json:"createdBy,string"`
		UpdatedBy uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy uint64                           `json:"deletedBy,string,omitempty"`
		Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ApigwRouteMeta struct {
		Debug  bool                             `json:"debug"`
		Async  bool                             `json:"async"`
		Desc   string                           `json:"description"`
		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}
)

func (r ApigwRoute) Clone() *ApigwRoute {
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

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r ApigwRoute) Diff(cmp *ApigwRoute) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApigwRoute{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "routeID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Endpoint != cmp.Endpoint {
		out = append(out, &revisions.Change{Key: "endpoint", Old: []any{cmp.Endpoint}, New: []any{r.Endpoint}})
	}

	if r.Method != cmp.Method {
		out = append(out, &revisions.Change{Key: "method", Old: []any{cmp.Method}, New: []any{r.Method}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if r.Group != cmp.Group {
		out = append(out, &revisions.Change{Key: "group", Old: []any{cmp.Group}, New: []any{r.Group}})
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

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r ApigwRouteMeta) Clone() *ApigwRouteMeta {
	dup := r
	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}

	return &dup
}

func (r ApigwRouteMeta) Diff(cmp *ApigwRouteMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApigwRouteMeta{}
	}
	if r.Debug != cmp.Debug {
		out = append(out, &revisions.Change{Key: "debug", Old: []any{cmp.Debug}, New: []any{r.Debug}})
	}

	if r.Async != cmp.Async {
		out = append(out, &revisions.Change{Key: "async", Old: []any{cmp.Async}, New: []any{r.Async}})
	}

	if r.Desc != cmp.Desc {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Desc}, New: []any{r.Desc}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}

	return out
}

func (r *ApigwRouteMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ApigwRouteMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseApigwRouteMeta(ss []string) (p ApigwRouteMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
