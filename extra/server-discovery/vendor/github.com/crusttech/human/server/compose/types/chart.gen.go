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
	Chart struct {
		ID          uint64                           `json:"chartID,string"`
		Handle      string                           `json:"handle"`
		TenantID    uint64                           `json:"tenantID,string,omitempty"`
		ProjectID   uint64                           `json:"projectID,string,omitempty"`
		NamespaceID uint64                           `json:"namespaceID,string"`
		Name        string                           `json:"name"`
		Meta        *ChartMeta                       `json:"meta,omitempty"`
		Config      ChartConfig                      `json:"config"`
		CreatedAt   time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
		Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ChartMeta struct {
		Description string `json:"description,omitempty"`
	}

	ChartConfig struct {
		Reports     []*ChartConfigReport   `json:"reports,omitempty"`
		ColorScheme string                 `json:"colorScheme,omitempty"`
		NoAnimation bool                   `json:"noAnimation,omitempty"`
		Toolbox     map[string]interface{} `json:"toolbox,omitempty"`
	}
)

func (r Chart) Clone() *Chart {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	dup.Config = *r.Config.Clone()

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

func (r Chart) Diff(cmp *Chart) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Chart{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "chartID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
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

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r ChartMeta) Clone() *ChartMeta {
	dup := r
	return &dup
}

func (r ChartMeta) Diff(cmp *ChartMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChartMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *ChartMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChartMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChartConfig) Clone() *ChartConfig {
	dup := r
	if r.Reports != nil {
		dup.Reports = make([]*ChartConfigReport, len(r.Reports))
		copy(dup.Reports, r.Reports)
	}

	if r.Toolbox != nil {
		dup.Toolbox = make(map[string]interface{}, len(r.Toolbox))
		for k, v := range r.Toolbox {
			dup.Toolbox[k] = v
		}
	}

	return &dup
}

func (r ChartConfig) Diff(cmp *ChartConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChartConfig{}
	}
	if !reflect.DeepEqual(r.Reports, cmp.Reports) {
		out = append(out, &revisions.Change{Key: "reports", Old: []any{cmp.Reports}, New: []any{r.Reports}})
	}

	if r.ColorScheme != cmp.ColorScheme {
		out = append(out, &revisions.Change{Key: "colorScheme", Old: []any{cmp.ColorScheme}, New: []any{r.ColorScheme}})
	}

	if r.NoAnimation != cmp.NoAnimation {
		out = append(out, &revisions.Change{Key: "noAnimation", Old: []any{cmp.NoAnimation}, New: []any{r.NoAnimation}})
	}

	if !reflect.DeepEqual(r.Toolbox, cmp.Toolbox) {
		out = append(out, &revisions.Change{Key: "toolbox", Old: []any{cmp.Toolbox}, New: []any{r.Toolbox}})
	}

	return out
}

func (r *ChartConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChartConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseChartMeta(ss []string) (p ChartMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
