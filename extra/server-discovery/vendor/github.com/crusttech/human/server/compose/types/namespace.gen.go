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
	Namespace struct {
		ID             uint64                           `json:"namespaceID,string"`
		TenantID       uint64                           `json:"tenantID,string,omitempty"`
		ProjectID      uint64                           `json:"projectID,string,omitempty"`
		Slug           string                           `json:"slug"`
		Enabled        bool                             `json:"enabled"`
		Meta           NamespaceMeta                    `json:"meta"`
		Name           string                           `json:"name"`
		CreatedAt      time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy      uint64                           `json:"createdBy,string"`
		CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
		Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	NamespaceMeta struct {
		Icon        string `json:"icon,omitempty"`
		IconID      uint64 `json:"iconID,string"`
		Logo        string `json:"logo,omitempty"`
		LogoID      uint64 `json:"logoID,string"`
		LogoEnabled bool   `json:"logoEnabled,omitempty"`
		HideSidebar bool   `json:"hideSidebar"`
		Subtitle    string `json:"subtitle,omitempty"`
		Description string `json:"description,omitempty"`
	}
)

func (r Namespace) Clone() *Namespace {
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

func (r Namespace) Diff(cmp *Namespace) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Namespace{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Slug != cmp.Slug {
		out = append(out, &revisions.Change{Key: "slug", Old: []any{cmp.Slug}, New: []any{r.Slug}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
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

	if r.CreatedByAgent != cmp.CreatedByAgent {
		out = append(out, &revisions.Change{Key: "createdByAgent", Old: []any{cmp.CreatedByAgent}, New: []any{r.CreatedByAgent}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r NamespaceMeta) Clone() *NamespaceMeta {
	dup := r
	return &dup
}

func (r NamespaceMeta) Diff(cmp *NamespaceMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NamespaceMeta{}
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

	if r.LogoEnabled != cmp.LogoEnabled {
		out = append(out, &revisions.Change{Key: "logoEnabled", Old: []any{cmp.LogoEnabled}, New: []any{r.LogoEnabled}})
	}

	if r.HideSidebar != cmp.HideSidebar {
		out = append(out, &revisions.Change{Key: "hideSidebar", Old: []any{cmp.HideSidebar}, New: []any{r.HideSidebar}})
	}

	if r.Subtitle != cmp.Subtitle {
		out = append(out, &revisions.Change{Key: "subtitle", Old: []any{cmp.Subtitle}, New: []any{r.Subtitle}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *NamespaceMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NamespaceMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseNamespaceMeta(ss []string) (p NamespaceMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
