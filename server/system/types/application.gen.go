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
		ID        uint64                           `json:"applicationID,string"`
		TenantID  uint64                           `json:"tenantID,string,omitempty"`
		ProjectID uint64                           `json:"projectID,string,omitempty"`
		Name      string                           `json:"name"`
		Enabled   bool                             `json:"enabled"`
		Weight    int                              `json:"weight"`
		Unify     *ApplicationUnify                `json:"unify,omitempty"`
		OwnerID   uint64                           `json:"ownerID"`
		CreatedAt time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
		Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Flags     []string                         `json:"flags,omitempty"`
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
	}
)

func (r Application) Clone() *Application {
	dup := r
	if r.Unify != nil {
		dup.Unify = r.Unify.Clone()
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

	if (r.Unify == nil) != (cmp.Unify == nil) {
		out = append(out, &revisions.Change{Key: "unify", Old: []any{cmp.Unify}, New: []any{r.Unify}})
	} else if r.Unify != nil {
		for _, c := range r.Unify.Diff(cmp.Unify) {
			c.Key = "unify." + c.Key
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

	return out
}

func (r *ApplicationUnify) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ApplicationUnify) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseApplicationUnify(ss []string) (p ApplicationUnify, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
