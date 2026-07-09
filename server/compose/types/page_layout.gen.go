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
	PageLayout struct {
		ID             uint64                           `json:"pageLayoutID,string"`
		TenantID       uint64                           `json:"tenantID,string,omitempty"`
		ProjectID      uint64                           `json:"projectID,string,omitempty"`
		Handle         string                           `json:"handle"`
		Primary        bool                             `json:"primary"`
		PageID         uint64                           `json:"pageID,string"`
		ParentID       uint64                           `json:"parentID,string"`
		NamespaceID    uint64                           `json:"namespaceID,string"`
		Weight         int                              `json:"weight"`
		Meta           PageLayoutMeta                   `json:"meta,omitempty"`
		Config         PageLayoutConfig                 `json:"config"`
		Blocks         PageLayoutBlocks                 `json:"blocks,omitempty"`
		OwnedBy        uint64                           `json:"ownedBy,string"`
		CreatedAt      time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
		CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
		Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	PageLayoutMeta struct {
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Style       map[string]any `json:"style,omitempty"`
	}

	PageLayoutConfig struct {
		Visibility PageLayoutVisibility   `json:"visibility"`
		Buttons    PageLayoutButtonConfig `json:"buttons"`
		Actions    []PageLayoutAction     `json:"actions,omitempty"`
		Validation PageLayoutValidation   `json:"validation"`
		UseTitle   bool                   `json:"useTitle"`
	}

	PageLayoutVisibility struct {
		Expression string   `json:"expression"`
		Roles      []string `json:"roles,omitempty"`
	}

	PageLayoutButton struct {
		Enabled bool   `json:"enabled"`
		Label   string `json:"label"`
	}

	PageLayoutAction struct {
		ActionID  uint64               `json:"actionID,string"`
		Placement string               `json:"placement"`
		Meta      PageLayoutActionMeta `json:"meta"`
		Enabled   bool                 `json:"enabled"`
		Kind      string               `json:"kind"`
		Params    any                  `json:"params"`
	}

	PageLayoutActionMeta struct {
		Label string         `json:"label"`
		Style map[string]any `json:"style,omitempty"`
	}

	PageLayoutValidation struct {
		RequiredFields []PageLayoutRequiredField `json:"requiredFields,omitempty"`
	}

	PageLayoutRequiredField struct {
		Field     string `json:"field"`
		Condition string `json:"condition"`
	}

	PageLayoutBlock struct {
		BlockID uint64         `json:"blockID,string,omitempty"`
		XYWH    [4]int         `json:"xywh"`
		Meta    map[string]any `json:"meta,omitempty"`
	}
)

func (r PageLayout) Clone() *PageLayout {
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

func (r PageLayout) Diff(cmp *PageLayout) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayout{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "pageLayoutID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Primary != cmp.Primary {
		out = append(out, &revisions.Change{Key: "primary", Old: []any{cmp.Primary}, New: []any{r.Primary}})
	}

	if r.PageID != cmp.PageID {
		out = append(out, &revisions.Change{Key: "pageID", Old: []any{cmp.PageID}, New: []any{r.PageID}})
	}

	if r.ParentID != cmp.ParentID {
		out = append(out, &revisions.Change{Key: "parentID", Old: []any{cmp.ParentID}, New: []any{r.ParentID}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if r.Weight != cmp.Weight {
		out = append(out, &revisions.Change{Key: "weight", Old: []any{cmp.Weight}, New: []any{r.Weight}})
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

	if r.OwnedBy != cmp.OwnedBy {
		out = append(out, &revisions.Change{Key: "ownedBy", Old: []any{cmp.OwnedBy}, New: []any{r.OwnedBy}})
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

func (r PageLayoutMeta) Clone() *PageLayoutMeta {
	dup := r
	if r.Style != nil {
		dup.Style = make(map[string]any, len(r.Style))
		for k, v := range r.Style {
			dup.Style[k] = v
		}
	}

	return &dup
}

func (r PageLayoutMeta) Diff(cmp *PageLayoutMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutMeta{}
	}
	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.Style, cmp.Style) {
		out = append(out, &revisions.Change{Key: "style", Old: []any{cmp.Style}, New: []any{r.Style}})
	}

	return out
}

func (r *PageLayoutMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutConfig) Clone() *PageLayoutConfig {
	dup := r
	dup.Visibility = *r.Visibility.Clone()

	if r.Actions != nil {
		dup.Actions = make([]PageLayoutAction, len(r.Actions))
		for i := range r.Actions {
			dup.Actions[i] = *r.Actions[i].Clone()
		}
	}

	dup.Validation = *r.Validation.Clone()

	return &dup
}

func (r PageLayoutConfig) Diff(cmp *PageLayoutConfig) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutConfig{}
	}
	for _, c := range r.Visibility.Diff(&cmp.Visibility) {
		c.Key = "visibility." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Buttons, cmp.Buttons) {
		out = append(out, &revisions.Change{Key: "buttons", Old: []any{cmp.Buttons}, New: []any{r.Buttons}})
	}

	if !reflect.DeepEqual(r.Actions, cmp.Actions) {
		out = append(out, &revisions.Change{Key: "actions", Old: []any{cmp.Actions}, New: []any{r.Actions}})
	}

	for _, c := range r.Validation.Diff(&cmp.Validation) {
		c.Key = "validation." + c.Key
		out = append(out, c)
	}

	if r.UseTitle != cmp.UseTitle {
		out = append(out, &revisions.Change{Key: "useTitle", Old: []any{cmp.UseTitle}, New: []any{r.UseTitle}})
	}

	return out
}

func (r *PageLayoutConfig) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutConfig) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutVisibility) Clone() *PageLayoutVisibility {
	dup := r
	if r.Roles != nil {
		dup.Roles = make([]string, len(r.Roles))
		copy(dup.Roles, r.Roles)
	}

	return &dup
}

func (r PageLayoutVisibility) Diff(cmp *PageLayoutVisibility) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutVisibility{}
	}
	if r.Expression != cmp.Expression {
		out = append(out, &revisions.Change{Key: "expression", Old: []any{cmp.Expression}, New: []any{r.Expression}})
	}

	if !reflect.DeepEqual(r.Roles, cmp.Roles) {
		out = append(out, &revisions.Change{Key: "roles", Old: []any{cmp.Roles}, New: []any{r.Roles}})
	}

	return out
}

func (r *PageLayoutVisibility) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutVisibility) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutButton) Clone() *PageLayoutButton {
	dup := r
	return &dup
}

func (r PageLayoutButton) Diff(cmp *PageLayoutButton) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutButton{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	return out
}

func (r *PageLayoutButton) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutButton) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutAction) Clone() *PageLayoutAction {
	dup := r
	dup.Meta = *r.Meta.Clone()

	return &dup
}

func (r PageLayoutAction) Diff(cmp *PageLayoutAction) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutAction{}
	}
	if r.ActionID != cmp.ActionID {
		out = append(out, &revisions.Change{Key: "actionID", Old: []any{cmp.ActionID}, New: []any{r.ActionID}})
	}

	if r.Placement != cmp.Placement {
		out = append(out, &revisions.Change{Key: "placement", Old: []any{cmp.Placement}, New: []any{r.Placement}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	return out
}

func (r *PageLayoutAction) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutAction) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutActionMeta) Clone() *PageLayoutActionMeta {
	dup := r
	if r.Style != nil {
		dup.Style = make(map[string]any, len(r.Style))
		for k, v := range r.Style {
			dup.Style[k] = v
		}
	}

	return &dup
}

func (r PageLayoutActionMeta) Diff(cmp *PageLayoutActionMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutActionMeta{}
	}
	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if !reflect.DeepEqual(r.Style, cmp.Style) {
		out = append(out, &revisions.Change{Key: "style", Old: []any{cmp.Style}, New: []any{r.Style}})
	}

	return out
}

func (r *PageLayoutActionMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutActionMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutValidation) Clone() *PageLayoutValidation {
	dup := r
	if r.RequiredFields != nil {
		dup.RequiredFields = make([]PageLayoutRequiredField, len(r.RequiredFields))
		for i := range r.RequiredFields {
			dup.RequiredFields[i] = *r.RequiredFields[i].Clone()
		}
	}

	return &dup
}

func (r PageLayoutValidation) Diff(cmp *PageLayoutValidation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutValidation{}
	}
	if !reflect.DeepEqual(r.RequiredFields, cmp.RequiredFields) {
		out = append(out, &revisions.Change{Key: "requiredFields", Old: []any{cmp.RequiredFields}, New: []any{r.RequiredFields}})
	}

	return out
}

func (r *PageLayoutValidation) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutValidation) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutRequiredField) Clone() *PageLayoutRequiredField {
	dup := r
	return &dup
}

func (r PageLayoutRequiredField) Diff(cmp *PageLayoutRequiredField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutRequiredField{}
	}
	if r.Field != cmp.Field {
		out = append(out, &revisions.Change{Key: "field", Old: []any{cmp.Field}, New: []any{r.Field}})
	}

	if r.Condition != cmp.Condition {
		out = append(out, &revisions.Change{Key: "condition", Old: []any{cmp.Condition}, New: []any{r.Condition}})
	}

	return out
}

func (r *PageLayoutRequiredField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutRequiredField) Value() (driver.Value, error) { return json.Marshal(r) }

func (r PageLayoutBlock) Clone() *PageLayoutBlock {
	dup := r
	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r PageLayoutBlock) Diff(cmp *PageLayoutBlock) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &PageLayoutBlock{}
	}
	if r.BlockID != cmp.BlockID {
		out = append(out, &revisions.Change{Key: "blockID", Old: []any{cmp.BlockID}, New: []any{r.BlockID}})
	}

	if r.XYWH != cmp.XYWH {
		out = append(out, &revisions.Change{Key: "xywh", Old: []any{cmp.XYWH}, New: []any{r.XYWH}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	return out
}

func (r *PageLayoutBlock) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r PageLayoutBlock) Value() (driver.Value, error) { return json.Marshal(r) }

func ParsePageLayoutMeta(ss []string) (p PageLayoutMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParsePageLayoutConfig(ss []string) (p PageLayoutConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageLayoutBlocks) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageLayoutBlocks) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageLayoutBlocks(ss []string) (p PageLayoutBlocks, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
