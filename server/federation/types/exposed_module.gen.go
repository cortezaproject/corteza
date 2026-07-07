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
	ExposedModule struct {
		ID                 uint64         `json:"moduleID,string"`
		TenantID           uint64         `json:"tenantID,string,omitempty"`
		ProjectID          uint64         `json:"projectID,string,omitempty"`
		Handle             string         `json:"handle"`
		Name               string         `json:"name"`
		NodeID             uint64         `json:"nodeID,string"`
		ComposeModuleID    uint64         `json:"composeModuleID,string"`
		ComposeNamespaceID uint64         `json:"composeNamespaceID,string"`
		Fields             ModuleFieldSet `json:"fields"`
		CreatedAt          time.Time      `json:"createdAt,omitempty"`
		UpdatedAt          *time.Time     `json:"updatedAt,omitempty"`
		DeletedAt          *time.Time     `json:"deletedAt,omitempty"`
		CreatedBy          uint64         `json:"createdBy,string"`
		UpdatedBy          uint64         `json:"updatedBy,string,omitempty"`
		DeletedBy          uint64         `json:"deletedBy,string,omitempty"`
	}

	ModuleField struct {
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Label   string `json:"label"`
		IsMulti bool   `json:"isMulti"`
	}
)

func (r ExposedModule) Clone() *ExposedModule {
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

func (r ExposedModule) Diff(cmp *ExposedModule) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ExposedModule{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "moduleID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.NodeID != cmp.NodeID {
		out = append(out, &revisions.Change{Key: "nodeID", Old: []any{cmp.NodeID}, New: []any{r.NodeID}})
	}

	if r.ComposeModuleID != cmp.ComposeModuleID {
		out = append(out, &revisions.Change{Key: "composeModuleID", Old: []any{cmp.ComposeModuleID}, New: []any{r.ComposeModuleID}})
	}

	if r.ComposeNamespaceID != cmp.ComposeNamespaceID {
		out = append(out, &revisions.Change{Key: "composeNamespaceID", Old: []any{cmp.ComposeNamespaceID}, New: []any{r.ComposeNamespaceID}})
	}

	if !reflect.DeepEqual(r.Fields, cmp.Fields) {
		out = append(out, &revisions.Change{Key: "fields", Old: []any{cmp.Fields}, New: []any{r.Fields}})
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

func (r ModuleField) Clone() *ModuleField {
	dup := r
	return &dup
}

func (r ModuleField) Diff(cmp *ModuleField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ModuleField{}
	}
	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if r.IsMulti != cmp.IsMulti {
		out = append(out, &revisions.Change{Key: "isMulti", Old: []any{cmp.IsMulti}, New: []any{r.IsMulti}})
	}

	return out
}

func (r *ModuleField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ModuleField) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *ModuleFieldSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ModuleFieldSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseModuleFieldSet(ss []string) (p ModuleFieldSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
