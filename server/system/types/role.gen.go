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
	Role struct {
		ID         uint64                           `json:"roleID,string"`
		TenantID   uint64                           `json:"tenantID,string,omitempty"`
		ProjectID  uint64                           `json:"projectID,string,omitempty"`
		Name       string                           `json:"name"`
		Handle     string                           `json:"handle"`
		Meta       *RoleMeta                        `json:"meta"`
		ArchivedAt *time.Time                       `json:"archivedAt,omitempty"`
		CreatedAt  time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	RoleMeta struct {
		Description string       `json:"description,omitempty"`
		Context     *RoleContext `json:"context,omitempty"`
	}

	RoleContext struct {
		Resource []string `json:"resourceTypes,omitempty"`
		Expr     string   `json:"expr,omitempty"`
	}
)

func (r Role) Clone() *Role {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.ArchivedAt != nil {
		v := *r.ArchivedAt
		dup.ArchivedAt = &v
	}

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

func (r Role) Diff(cmp *Role) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Role{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "roleID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	if !reflect.DeepEqual(r.ArchivedAt, cmp.ArchivedAt) {
		out = append(out, &revisions.Change{Key: "archivedAt", Old: []any{cmp.ArchivedAt}, New: []any{r.ArchivedAt}})
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

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r RoleMeta) Clone() *RoleMeta {
	dup := r
	if r.Context != nil {
		dup.Context = r.Context.Clone()
	}

	return &dup
}

func (r RoleMeta) Diff(cmp *RoleMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &RoleMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if (r.Context == nil) != (cmp.Context == nil) {
		out = append(out, &revisions.Change{Key: "context", Old: []any{cmp.Context}, New: []any{r.Context}})
	} else if r.Context != nil {
		for _, c := range r.Context.Diff(cmp.Context) {
			c.Key = "context." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *RoleMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r RoleMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r RoleContext) Clone() *RoleContext {
	dup := r
	if r.Resource != nil {
		dup.Resource = make([]string, len(r.Resource))
		copy(dup.Resource, r.Resource)
	}

	return &dup
}

func (r RoleContext) Diff(cmp *RoleContext) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &RoleContext{}
	}
	if !reflect.DeepEqual(r.Resource, cmp.Resource) {
		out = append(out, &revisions.Change{Key: "resourceTypes", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.Expr != cmp.Expr {
		out = append(out, &revisions.Change{Key: "expr", Old: []any{cmp.Expr}, New: []any{r.Expr}})
	}

	return out
}

func (r *RoleContext) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r RoleContext) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseRoleMeta(ss []string) (p *RoleMeta, err error) {
	p = &RoleMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
