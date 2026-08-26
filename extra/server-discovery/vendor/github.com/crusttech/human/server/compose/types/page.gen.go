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
	Page struct {
		ID             uint64                           `json:"pageID,string"`
		TenantID       uint64                           `json:"tenantID,string,omitempty"`
		ProjectID      uint64                           `json:"projectID,string,omitempty"`
		Title          string                           `json:"title"`
		Handle         string                           `json:"handle"`
		SelfID         uint64                           `json:"selfID,string"`
		ModuleID       uint64                           `json:"moduleID,string"`
		NamespaceID    uint64                           `json:"namespaceID,string"`
		Meta           PageMeta                         `json:"meta"`
		Config         PageConfig                       `json:"config"`
		Blocks         PageBlocks                       `json:"blocks"`
		Children       PageSet                          `json:"children,omitempty"`
		Visible        bool                             `json:"visible"`
		Weight         int                              `json:"weight"`
		Description    string                           `json:"description"`
		CreatedAt      time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
		CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
		Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	PageBlock struct {
		BlockID     uint64                 `json:"blockID,string,omitempty"`
		Options     map[string]interface{} `json:"options,omitempty"`
		Style       PageBlockStyle         `json:"style,omitempty"`
		Kind        string                 `json:"kind"`
		XYWH        [4]int                 `json:"xywh"`
		Meta        map[string]any         `json:"meta,omitempty"`
		Title       string                 `json:"title,omitempty"`
		Description string                 `json:"description,omitempty"`
	}

	PageBlockStyle struct {
		Variants map[string]string      `json:"variants,omitempty"`
		Wrap     map[string]string      `json:"wrap,omitempty"`
		Border   map[string]interface{} `json:"border,omitempty"`
	}

	PageMeta struct {
		AllowPersonalLayouts bool           `json:"allowPersonalLayouts"`
		Notifications        map[string]any `json:"notifications,omitempty"`
	}

	PageConfig struct {
		NavItem PageConfigNavItem `json:"navItem"`
	}

	PageConfigNavItem struct {
		Expanded bool            `json:"expanded"`
		Icon     *PageConfigIcon `json:"icon,omitempty"`
	}

	PageConfigIcon struct {
		Type  IconType          `json:"type,omitempty"`
		Src   string            `json:"src"`
		Style map[string]string `json:"style,omitempty"`
	}

	// how child pages are handled on delete
	PageChildrenDeleteStrategy string
)

func (r Page) Clone() *Page {
	dup := r
	dup.Meta = *r.Meta.Clone()

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

func (r Page) Diff(cmp *Page) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Page{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "pageID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.SelfID != cmp.SelfID {
		out = append(out, &revisions.Change{Key: "selfID", Old: []any{cmp.SelfID}, New: []any{r.SelfID}})
	}

	if r.ModuleID != cmp.ModuleID {
		out = append(out, &revisions.Change{Key: "moduleID", Old: []any{cmp.ModuleID}, New: []any{r.ModuleID}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Config.Diff(&cmp.Config) {
		c.Key = "config." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Blocks, cmp.Blocks) {
		out = append(out, &revisions.Change{Key: "blocks", Old: []any{cmp.Blocks}, New: []any{r.Blocks}})
	}

	if !reflect.DeepEqual(r.Children, cmp.Children) {
		out = append(out, &revisions.Change{Key: "children", Old: []any{cmp.Children}, New: []any{r.Children}})
	}

	if r.Visible != cmp.Visible {
		out = append(out, &revisions.Change{Key: "visible", Old: []any{cmp.Visible}, New: []any{r.Visible}})
	}

	if r.Weight != cmp.Weight {
		out = append(out, &revisions.Change{Key: "weight", Old: []any{cmp.Weight}, New: []any{r.Weight}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
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

	if r.CreatedByAgent != cmp.CreatedByAgent {
		out = append(out, &revisions.Change{Key: "createdByAgent", Old: []any{cmp.CreatedByAgent}, New: []any{r.CreatedByAgent}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r PageBlock) Clone() *PageBlock {
	dup := r
	if r.Options != nil {
		dup.Options = make(map[string]interface{}, len(r.Options))
		for k, v := range r.Options {
			dup.Options[k] = v
		}
	}

	dup.Style = *r.Style.Clone()

	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r PageBlock) Diff(cmp *PageBlock) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageBlock{}
	}
	if r.BlockID != cmp.BlockID {
		out = append(out, &revisions.Change{Key: "blockID", Old: []any{cmp.BlockID}, New: []any{r.BlockID}})
	}

	if !reflect.DeepEqual(r.Options, cmp.Options) {
		out = append(out, &revisions.Change{Key: "options", Old: []any{cmp.Options}, New: []any{r.Options}})
	}

	for _, c := range r.Style.Diff(&cmp.Style) {
		c.Key = "style." + c.Key
		out = append(out, c)
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if r.XYWH != cmp.XYWH {
		out = append(out, &revisions.Change{Key: "xywh", Old: []any{cmp.XYWH}, New: []any{r.XYWH}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *PageBlock) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageBlock) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageBlockStyle) Clone() *PageBlockStyle {
	dup := r
	if r.Variants != nil {
		dup.Variants = make(map[string]string, len(r.Variants))
		for k, v := range r.Variants {
			dup.Variants[k] = v
		}
	}

	if r.Wrap != nil {
		dup.Wrap = make(map[string]string, len(r.Wrap))
		for k, v := range r.Wrap {
			dup.Wrap[k] = v
		}
	}

	if r.Border != nil {
		dup.Border = make(map[string]interface{}, len(r.Border))
		for k, v := range r.Border {
			dup.Border[k] = v
		}
	}

	return &dup
}

func (r PageBlockStyle) Diff(cmp *PageBlockStyle) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageBlockStyle{}
	}
	if !reflect.DeepEqual(r.Variants, cmp.Variants) {
		out = append(out, &revisions.Change{Key: "variants", Old: []any{cmp.Variants}, New: []any{r.Variants}})
	}

	if !reflect.DeepEqual(r.Wrap, cmp.Wrap) {
		out = append(out, &revisions.Change{Key: "wrap", Old: []any{cmp.Wrap}, New: []any{r.Wrap}})
	}

	if !reflect.DeepEqual(r.Border, cmp.Border) {
		out = append(out, &revisions.Change{Key: "border", Old: []any{cmp.Border}, New: []any{r.Border}})
	}

	return out
}

func (r *PageBlockStyle) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageBlockStyle) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageMeta) Clone() *PageMeta {
	dup := r
	if r.Notifications != nil {
		dup.Notifications = make(map[string]any, len(r.Notifications))
		for k, v := range r.Notifications {
			dup.Notifications[k] = v
		}
	}

	return &dup
}

func (r PageMeta) Diff(cmp *PageMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageMeta{}
	}
	if r.AllowPersonalLayouts != cmp.AllowPersonalLayouts {
		out = append(out, &revisions.Change{Key: "allowPersonalLayouts", Old: []any{cmp.AllowPersonalLayouts}, New: []any{r.AllowPersonalLayouts}})
	}

	if !reflect.DeepEqual(r.Notifications, cmp.Notifications) {
		out = append(out, &revisions.Change{Key: "notifications", Old: []any{cmp.Notifications}, New: []any{r.Notifications}})
	}

	return out
}

func (r *PageMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageConfig) Clone() *PageConfig {
	dup := r
	dup.NavItem = *r.NavItem.Clone()

	return &dup
}

func (r PageConfig) Diff(cmp *PageConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageConfig{}
	}
	for _, c := range r.NavItem.Diff(&cmp.NavItem) {
		c.Key = "navItem." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *PageConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageConfigNavItem) Clone() *PageConfigNavItem {
	dup := r
	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	return &dup
}

func (r PageConfigNavItem) Diff(cmp *PageConfigNavItem) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageConfigNavItem{}
	}
	if r.Expanded != cmp.Expanded {
		out = append(out, &revisions.Change{Key: "expanded", Old: []any{cmp.Expanded}, New: []any{r.Expanded}})
	}

	if (r.Icon == nil) != (cmp.Icon == nil) {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	} else if r.Icon != nil {
		for _, c := range r.Icon.Diff(cmp.Icon) {
			c.Key = "icon." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *PageConfigNavItem) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageConfigNavItem) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageConfigIcon) Clone() *PageConfigIcon {
	dup := r
	if r.Style != nil {
		dup.Style = make(map[string]string, len(r.Style))
		for k, v := range r.Style {
			dup.Style[k] = v
		}
	}

	return &dup
}

func (r PageConfigIcon) Diff(cmp *PageConfigIcon) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageConfigIcon{}
	}
	if !reflect.DeepEqual(r.Type, cmp.Type) {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Src != cmp.Src {
		out = append(out, &revisions.Change{Key: "src", Old: []any{cmp.Src}, New: []any{r.Src}})
	}

	if !reflect.DeepEqual(r.Style, cmp.Style) {
		out = append(out, &revisions.Change{Key: "style", Old: []any{cmp.Style}, New: []any{r.Style}})
	}

	return out
}

func (r *PageConfigIcon) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageConfigIcon) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	PageChildrenDeleteStrategyAbort PageChildrenDeleteStrategy = "abort"
	// reattach children to parent
	PageChildrenDeleteStrategyRebase PageChildrenDeleteStrategy = "rebase"
)

func ParsePageMeta(ss []string) (p PageMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParsePageConfig(ss []string) (p PageConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageBlocks) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageBlocks) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageBlocks(ss []string) (p PageBlocks, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
