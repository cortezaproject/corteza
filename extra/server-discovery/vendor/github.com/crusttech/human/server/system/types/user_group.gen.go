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
	UserGroup struct {
		ID         uint64                           `json:"userGroupID,string"`
		TenantID   uint64                           `json:"tenantID,string,omitempty"`
		ProjectID  uint64                           `json:"projectID,string,omitempty"`
		Handle     string                           `json:"handle"`
		Meta       *UserGroupMeta                   `json:"meta"`
		Config     *UserGroupConfig                 `json:"config"`
		IsRoot     bool                             `json:"isRoot"`
		ArchivedAt *time.Time                       `json:"archivedAt,omitempty"`
		CreatedAt  time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	UserGroupConfig struct {
		Paths []UserGroupPath `json:"path"`
	}

	UserGroupMeta struct {
		Description string `json:"description"`
		Short       string `json:"short"`
	}

	UserGroupPath struct {
		SelfID uint64 `json:"selfID,string"`
		Name   string `json:"name"`
	}
)

func (r UserGroup) Clone() *UserGroup {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.Config != nil {
		dup.Config = r.Config.Clone()
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

func (r UserGroup) Diff(cmp *UserGroup) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &UserGroup{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "userGroupID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	if (r.Config == nil) != (cmp.Config == nil) {
		out = append(out, &revisions.Change{Key: "config", Old: []any{cmp.Config}, New: []any{r.Config}})
	} else if r.Config != nil {
		for _, c := range r.Config.Diff(cmp.Config) {
			c.Key = "config." + c.Key
			out = append(out, c)
		}
	}

	if r.IsRoot != cmp.IsRoot {
		out = append(out, &revisions.Change{Key: "isRoot", Old: []any{cmp.IsRoot}, New: []any{r.IsRoot}})
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

func (r UserGroupConfig) Clone() *UserGroupConfig {
	dup := r
	if r.Paths != nil {
		dup.Paths = make([]UserGroupPath, len(r.Paths))
		for i := range r.Paths {
			dup.Paths[i] = *r.Paths[i].Clone()
		}
	}

	return &dup
}

func (r UserGroupConfig) Diff(cmp *UserGroupConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &UserGroupConfig{}
	}
	if !reflect.DeepEqual(r.Paths, cmp.Paths) {
		out = append(out, &revisions.Change{Key: "path", Old: []any{cmp.Paths}, New: []any{r.Paths}})
	}

	return out
}

func (r *UserGroupConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r UserGroupConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r UserGroupMeta) Clone() *UserGroupMeta {
	dup := r
	return &dup
}

func (r UserGroupMeta) Diff(cmp *UserGroupMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &UserGroupMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	return out
}

func (r *UserGroupMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r UserGroupMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r UserGroupPath) Clone() *UserGroupPath {
	dup := r
	return &dup
}

func (r UserGroupPath) Diff(cmp *UserGroupPath) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &UserGroupPath{}
	}
	if r.SelfID != cmp.SelfID {
		out = append(out, &revisions.Change{Key: "selfID", Old: []any{cmp.SelfID}, New: []any{r.SelfID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	return out
}

func (r *UserGroupPath) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r UserGroupPath) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseUserGroupMeta(ss []string) (p *UserGroupMeta, err error) {
	p = &UserGroupMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func ParseUserGroupConfig(ss []string) (p *UserGroupConfig, err error) {
	p = &UserGroupConfig{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
