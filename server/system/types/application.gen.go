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
	Application struct {
		ID         uint64                           `json:"applicationID,string"`
		TenantID   uint64                           `json:"tenantID,string,omitempty"`
		ProjectID  uint64                           `json:"projectID,string,omitempty"`
		Name       string                           `json:"name"`
		Enabled    bool                             `json:"enabled"`
		Weight     int                              `json:"weight"`
		Meta       *ApplicationMeta                 `json:"meta,omitempty"`
		Unify      *ApplicationUnify                `json:"unify,omitempty"`
		Source     string                           `json:"-"`
		SourceMeta *ApplicationSourceMeta           `json:"sourceMeta,omitempty"`
		OwnerID    uint64                           `json:"ownerID,string"`
		CreatedAt  time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Flags      []string                         `json:"flags,omitempty"`
	}

	ApplicationMeta struct {
		Description string `json:"description,omitempty"`
	}

	ApplicationUnify struct {
		Name   string `json:"name,omitempty"`
		Listed bool   `json:"listed"`
		Url    string `json:"url"`
		Config string `json:"config"`
		Icon   string `json:"icon,omitempty"`
		IconID uint64 `json:"iconID,string"`
		Logo   string `json:"logo,omitempty"`
		LogoID uint64 `json:"logoID,string"`
		Kind   string `json:"kind,omitempty"`
		Home   bool   `json:"home,omitempty"`
	}

	ApplicationSourceMeta struct {
		Hash        string            `json:"hash,omitempty"`
		Size        int               `json:"size"`
		Namespace   string            `json:"namespace,omitempty"`
		Modules     []string          `json:"modules,omitempty"`
		Writes      []string          `json:"writes,omitempty"`
		NamespaceID uint64            `json:"namespaceID,string,omitempty"`
		ModuleIDs   map[string]string `json:"moduleIDs,omitempty"`
		UpdatedAt   *time.Time        `json:"updatedAt,omitempty"`
		UpdatedBy   uint64            `json:"updatedBy,string,omitempty"`
	}
)

func (r Application) Clone() *Application {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.Unify != nil {
		dup.Unify = r.Unify.Clone()
	}

	if r.SourceMeta != nil {
		dup.SourceMeta = r.SourceMeta.Clone()
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
	if r.Flags != nil {
		dup.Flags = make([]string, len(r.Flags))
		copy(dup.Flags, r.Flags)
	}
	return &dup
}

func (r Application) Diff(cmp *Application) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Application{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "applicationID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Weight != cmp.Weight {
		out = append(out, &revisions.Change{Key: "weight", Old: []any{cmp.Weight}, New: []any{r.Weight}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	if (r.Unify == nil) != (cmp.Unify == nil) {
		out = append(out, &revisions.Change{Key: "unify", Old: []any{cmp.Unify}, New: []any{r.Unify}})
	} else if r.Unify != nil {
		for _, c := range r.Unify.Diff(cmp.Unify) {
			c.Key = "unify." + c.Key
			out = append(out, c)
		}
	}

	if (r.SourceMeta == nil) != (cmp.SourceMeta == nil) {
		out = append(out, &revisions.Change{Key: "sourceMeta", Old: []any{cmp.SourceMeta}, New: []any{r.SourceMeta}})
	} else if r.SourceMeta != nil {
		for _, c := range r.SourceMeta.Diff(cmp.SourceMeta) {
			c.Key = "sourceMeta." + c.Key
			out = append(out, c)
		}
	}

	if r.OwnerID != cmp.OwnerID {
		out = append(out, &revisions.Change{Key: "ownerID", Old: []any{cmp.OwnerID}, New: []any{r.OwnerID}})
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
	if !reflect.DeepEqual(r.Flags, cmp.Flags) {
		out = append(out, &revisions.Change{Key: "flags", Old: []any{cmp.Flags}, New: []any{r.Flags}})
	}
	return out
}

func (r ApplicationMeta) Clone() *ApplicationMeta {
	dup := r
	return &dup
}

func (r ApplicationMeta) Diff(cmp *ApplicationMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApplicationMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *ApplicationMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ApplicationMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ApplicationUnify) Clone() *ApplicationUnify {
	dup := r
	return &dup
}

func (r ApplicationUnify) Diff(cmp *ApplicationUnify) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApplicationUnify{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Listed != cmp.Listed {
		out = append(out, &revisions.Change{Key: "listed", Old: []any{cmp.Listed}, New: []any{r.Listed}})
	}

	if r.Url != cmp.Url {
		out = append(out, &revisions.Change{Key: "url", Old: []any{cmp.Url}, New: []any{r.Url}})
	}

	if r.Config != cmp.Config {
		out = append(out, &revisions.Change{Key: "config", Old: []any{cmp.Config}, New: []any{r.Config}})
	}

	if r.Icon != cmp.Icon {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	}

	if r.IconID != cmp.IconID {
		out = append(out, &revisions.Change{Key: "iconID", Old: []any{cmp.IconID}, New: []any{r.IconID}})
	}

	if r.Logo != cmp.Logo {
		out = append(out, &revisions.Change{Key: "logo", Old: []any{cmp.Logo}, New: []any{r.Logo}})
	}

	if r.LogoID != cmp.LogoID {
		out = append(out, &revisions.Change{Key: "logoID", Old: []any{cmp.LogoID}, New: []any{r.LogoID}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if r.Home != cmp.Home {
		out = append(out, &revisions.Change{Key: "home", Old: []any{cmp.Home}, New: []any{r.Home}})
	}

	return out
}

func (r *ApplicationUnify) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ApplicationUnify) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ApplicationSourceMeta) Clone() *ApplicationSourceMeta {
	dup := r
	if r.Modules != nil {
		dup.Modules = make([]string, len(r.Modules))
		copy(dup.Modules, r.Modules)
	}

	if r.Writes != nil {
		dup.Writes = make([]string, len(r.Writes))
		copy(dup.Writes, r.Writes)
	}

	if r.ModuleIDs != nil {
		dup.ModuleIDs = make(map[string]string, len(r.ModuleIDs))
		for k, v := range r.ModuleIDs {
			dup.ModuleIDs[k] = v
		}
	}

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	return &dup
}

func (r ApplicationSourceMeta) Diff(cmp *ApplicationSourceMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ApplicationSourceMeta{}
	}
	if r.Hash != cmp.Hash {
		out = append(out, &revisions.Change{Key: "hash", Old: []any{cmp.Hash}, New: []any{r.Hash}})
	}

	if r.Size != cmp.Size {
		out = append(out, &revisions.Change{Key: "size", Old: []any{cmp.Size}, New: []any{r.Size}})
	}

	if r.Namespace != cmp.Namespace {
		out = append(out, &revisions.Change{Key: "namespace", Old: []any{cmp.Namespace}, New: []any{r.Namespace}})
	}

	if !reflect.DeepEqual(r.Modules, cmp.Modules) {
		out = append(out, &revisions.Change{Key: "modules", Old: []any{cmp.Modules}, New: []any{r.Modules}})
	}

	if !reflect.DeepEqual(r.Writes, cmp.Writes) {
		out = append(out, &revisions.Change{Key: "writes", Old: []any{cmp.Writes}, New: []any{r.Writes}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if !reflect.DeepEqual(r.ModuleIDs, cmp.ModuleIDs) {
		out = append(out, &revisions.Change{Key: "moduleIDs", Old: []any{cmp.ModuleIDs}, New: []any{r.ModuleIDs}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	return out
}

func (r *ApplicationSourceMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ApplicationSourceMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseApplicationMeta(ss []string) (p ApplicationMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseApplicationUnify(ss []string) (p ApplicationUnify, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseApplicationSourceMeta(ss []string) (p ApplicationSourceMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
