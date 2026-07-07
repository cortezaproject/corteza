package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/ast"
	"github.com/crusttech/human/server/pkg/expr"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	NgAutomation struct {
		ID        uint64                           `json:"automationID,string"`
		TenantID  uint64                           `json:"tenantID,string,omitempty"`
		ProjectID uint64                           `json:"projectID,string,omitempty"`
		Handle    string                           `json:"handle"`
		Meta      *NgAutomationMeta                `json:"meta,omitempty"`
		Enabled   bool                             `json:"enabled"`
		Scope     *expr.Vars                       `json:"scope"`
		Triggers  NgAutomationTriggerSet           `json:"triggers"`
		Steps     NgAutomationStepSet              `json:"steps"`
		Paths     NgAutomationPathSet              `json:"paths"`
		Issues    NgAutomationIssueSet             `json:"issues,omitempty"`
		RunAs     uint64                           `json:"runAs,string"`
		OwnedBy   uint64                           `json:"ownedBy,string"`
		CreatedAt time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy uint64                           `json:"createdBy,string"`
		UpdatedBy uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy uint64                           `json:"deletedBy,string,omitempty"`
		Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	NgAutomationMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
	}

	NgAutomationVisual struct {
	}

	NgAutomationIcon struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	NgAutomationTrigger struct {
		ID           uint64                           `json:"triggerID,string"`
		Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Handle       string                           `json:"handle"`
		Meta         *NgTriggerMeta                   `json:"meta,omitempty"`
		Enabled      bool                             `json:"enabled"`
		ResourceType string                           `json:"resourceType"`
		EventType    string                           `json:"eventType"`
		Constraints  []NgTriggerConstraint            `json:"constraints"`
		Input        *expr.Vars                       `json:"input"`
		InputSchema  NgAutomationTriggerSchema        `json:"inputSchema,omitempty"`
	}

	NgAutomationTriggerParam struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Required    bool   `json:"required"`
		Description string `json:"description"`
	}

	NgTriggerConstraint struct {
		Name   string                     `json:"name"`
		Op     string                     `json:"op,omitempty"`
		Values []NgTriggerConstraintValue `json:"values,omitempty"`
	}

	NgTriggerConstraintValue struct {
		Type  string `json:"@type"`
		Value string `json:"@value"`
	}

	NgTriggerMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
	}

	NgAutomationStep struct {
		ID          uint64               `json:"stepID,string"`
		Handle      string               `json:"handle"`
		Meta        NgAutomationStepMeta `json:"meta"`
		Kind        string               `json:"kind"`
		Ref         string               `json:"ref"`
		Arguments   []*Expr              `json:"arguments"`
		Results     []*Expr              `json:"results"`
		Recoverable bool                 `json:"recoverable,omitempty"`
		MaxRetries  int                  `json:"maxRetries,omitempty"`
	}

	NgAutomationStepMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
		Extra       map[string]any     `json:"extra,omitempty"`
	}

	NgAutomationPath struct {
		ParentID  uint64               `json:"parentID,string"`
		ChildID   uint64               `json:"childID,string"`
		Condition *ast.ASTNode         `json:"condition,omitempty"`
		Kind      string               `json:"kind,omitempty"`
		Meta      NgAutomationPathMeta `json:"meta"`
	}

	NgAutomationPathMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
	}

	NgAutomationIssue struct {
		Code     string                     `json:"code"`
		Severity string                     `json:"severity"`
		Message  string                     `json:"message"`
		Details  []*NgAutomationIssueDetail `json:"details,omitempty"`
	}

	NgAutomationIssueDetail struct {
		MissingReference *DetailMissingReference    `json:"missingReference"`
		InvalidType      *DetailInvalidType         `json:"invalidType"`
		DuplicateID      *DetailDuplicateID         `json:"duplicateID"`
		Cycle            *DetailCycle               `json:"cycle"`
		ResourceRef      *DetailResourceRef         `json:"resourceRef"`
		GatewayPaths     *DetailGatewayPaths        `json:"gatewayPaths"`
		EmptyField       *DetailEmptyField          `json:"emptyField"`
		Details          []*NgAutomationIssueDetail `json:"details"`
	}

	DetailMissingReference struct {
		RefKind    string `json:"refKind"`
		Ref        string `json:"ref"`
		StepID     uint64 `json:"stepID,string,omitempty"`
		Field      string `json:"field,omitempty"`
		FieldIndex int    `json:"fieldIndex,omitempty"`
	}

	DetailInvalidType struct {
		StepID     uint64 `json:"stepID,string"`
		Field      string `json:"field,omitempty"`
		FieldIndex int    `json:"fieldIndex,omitempty"`
		Target     string `json:"target,omitempty"`
		Expected   string `json:"expected"`
		Actual     string `json:"actual"`
	}

	DetailDuplicateID struct {
		Resource string `json:"resource"`
		ID       uint64 `json:"id,string"`
		Indices  []int  `json:"indices,omitempty"`
	}

	DetailCycle struct {
		StepIDs IssueIDs `json:"stepIDs"`
	}

	DetailResourceRef struct {
		Resource   string `json:"resource"`
		ID         uint64 `json:"id,string,omitempty"`
		Index      int    `json:"index,omitempty"`
		Field      string `json:"field,omitempty"`
		FieldIndex int    `json:"fieldIndex,omitempty"`
	}

	DetailGatewayPaths struct {
		StepID    uint64 `json:"stepID,string"`
		Violation string `json:"violation"`
		Got       int    `json:"got"`
		Want      int    `json:"want,omitempty"`
	}

	DetailEmptyField struct {
		Resource string `json:"resource"`
		Index    int    `json:"index"`
		Field    string `json:"field,omitempty"`
	}
)

func (r NgAutomation) Clone() *NgAutomation {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.Scope != nil {
		v := *r.Scope
		dup.Scope = &v
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

func (r NgAutomation) Diff(cmp *NgAutomation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomation{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "automationID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if !reflect.DeepEqual(r.Scope, cmp.Scope) {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
	}

	if !reflect.DeepEqual(r.Triggers, cmp.Triggers) {
		out = append(out, &revisions.Change{Key: "triggers", Old: []any{cmp.Triggers}, New: []any{r.Triggers}})
	}

	if !reflect.DeepEqual(r.Steps, cmp.Steps) {
		out = append(out, &revisions.Change{Key: "steps", Old: []any{cmp.Steps}, New: []any{r.Steps}})
	}

	if !reflect.DeepEqual(r.Paths, cmp.Paths) {
		out = append(out, &revisions.Change{Key: "paths", Old: []any{cmp.Paths}, New: []any{r.Paths}})
	}

	if !reflect.DeepEqual(r.Issues, cmp.Issues) {
		out = append(out, &revisions.Change{Key: "issues", Old: []any{cmp.Issues}, New: []any{r.Issues}})
	}

	if r.RunAs != cmp.RunAs {
		out = append(out, &revisions.Change{Key: "runAs", Old: []any{cmp.RunAs}, New: []any{r.RunAs}})
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

func (r NgAutomationMeta) Clone() *NgAutomationMeta {
	dup := r
	dup.Visual = *r.Visual.Clone()

	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	return &dup
}

func (r NgAutomationMeta) Diff(cmp *NgAutomationMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	for _, c := range r.Visual.Diff(&cmp.Visual) {
		c.Key = "visual." + c.Key
		out = append(out, c)
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

func (r *NgAutomationMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationVisual) Clone() *NgAutomationVisual {
	dup := r
	return &dup
}

func (r NgAutomationVisual) Diff(cmp *NgAutomationVisual) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationVisual{}
	}
	return out
}

func (r *NgAutomationVisual) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationVisual) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationIcon) Clone() *NgAutomationIcon {
	dup := r
	return &dup
}

func (r NgAutomationIcon) Diff(cmp *NgAutomationIcon) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationIcon{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Value != cmp.Value {
		out = append(out, &revisions.Change{Key: "value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	return out
}

func (r NgAutomationTrigger) Clone() *NgAutomationTrigger {
	dup := r
	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}

	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.Constraints != nil {
		dup.Constraints = make([]NgTriggerConstraint, len(r.Constraints))
		for i := range r.Constraints {
			dup.Constraints[i] = *r.Constraints[i].Clone()
		}
	}

	if r.Input != nil {
		v := *r.Input
		dup.Input = &v
	}

	return &dup
}

func (r NgAutomationTrigger) Diff(cmp *NgAutomationTrigger) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationTrigger{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "triggerID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
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

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.ResourceType != cmp.ResourceType {
		out = append(out, &revisions.Change{Key: "resourceType", Old: []any{cmp.ResourceType}, New: []any{r.ResourceType}})
	}

	if r.EventType != cmp.EventType {
		out = append(out, &revisions.Change{Key: "eventType", Old: []any{cmp.EventType}, New: []any{r.EventType}})
	}

	if !reflect.DeepEqual(r.Constraints, cmp.Constraints) {
		out = append(out, &revisions.Change{Key: "constraints", Old: []any{cmp.Constraints}, New: []any{r.Constraints}})
	}

	if !reflect.DeepEqual(r.Input, cmp.Input) {
		out = append(out, &revisions.Change{Key: "input", Old: []any{cmp.Input}, New: []any{r.Input}})
	}

	if !reflect.DeepEqual(r.InputSchema, cmp.InputSchema) {
		out = append(out, &revisions.Change{Key: "inputSchema", Old: []any{cmp.InputSchema}, New: []any{r.InputSchema}})
	}

	return out
}

func (r *NgAutomationTrigger) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationTrigger) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationTriggerParam) Clone() *NgAutomationTriggerParam {
	dup := r
	return &dup
}

func (r NgAutomationTriggerParam) Diff(cmp *NgAutomationTriggerParam) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationTriggerParam{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Required != cmp.Required {
		out = append(out, &revisions.Change{Key: "required", Old: []any{cmp.Required}, New: []any{r.Required}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *NgAutomationTriggerParam) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationTriggerParam) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgTriggerConstraint) Clone() *NgTriggerConstraint {
	dup := r
	if r.Values != nil {
		dup.Values = make([]NgTriggerConstraintValue, len(r.Values))
		for i := range r.Values {
			dup.Values[i] = *r.Values[i].Clone()
		}
	}

	return &dup
}

func (r NgTriggerConstraint) Diff(cmp *NgTriggerConstraint) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgTriggerConstraint{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Op != cmp.Op {
		out = append(out, &revisions.Change{Key: "op", Old: []any{cmp.Op}, New: []any{r.Op}})
	}

	if !reflect.DeepEqual(r.Values, cmp.Values) {
		out = append(out, &revisions.Change{Key: "values", Old: []any{cmp.Values}, New: []any{r.Values}})
	}

	return out
}

func (r *NgTriggerConstraint) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgTriggerConstraint) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgTriggerConstraintValue) Clone() *NgTriggerConstraintValue {
	dup := r
	return &dup
}

func (r NgTriggerConstraintValue) Diff(cmp *NgTriggerConstraintValue) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgTriggerConstraintValue{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "@type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Value != cmp.Value {
		out = append(out, &revisions.Change{Key: "@value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	return out
}

func (r NgTriggerMeta) Clone() *NgTriggerMeta {
	dup := r
	dup.Visual = *r.Visual.Clone()

	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	return &dup
}

func (r NgTriggerMeta) Diff(cmp *NgTriggerMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgTriggerMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	for _, c := range r.Visual.Diff(&cmp.Visual) {
		c.Key = "visual." + c.Key
		out = append(out, c)
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

func (r *NgTriggerMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgTriggerMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationStep) Clone() *NgAutomationStep {
	dup := r
	dup.Meta = *r.Meta.Clone()

	if r.Arguments != nil {
		dup.Arguments = make([]*Expr, len(r.Arguments))
		copy(dup.Arguments, r.Arguments)
	}

	if r.Results != nil {
		dup.Results = make([]*Expr, len(r.Results))
		copy(dup.Results, r.Results)
	}

	return &dup
}

func (r NgAutomationStep) Diff(cmp *NgAutomationStep) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationStep{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if r.Ref != cmp.Ref {
		out = append(out, &revisions.Change{Key: "ref", Old: []any{cmp.Ref}, New: []any{r.Ref}})
	}

	if !reflect.DeepEqual(r.Arguments, cmp.Arguments) {
		out = append(out, &revisions.Change{Key: "arguments", Old: []any{cmp.Arguments}, New: []any{r.Arguments}})
	}

	if !reflect.DeepEqual(r.Results, cmp.Results) {
		out = append(out, &revisions.Change{Key: "results", Old: []any{cmp.Results}, New: []any{r.Results}})
	}

	if r.Recoverable != cmp.Recoverable {
		out = append(out, &revisions.Change{Key: "recoverable", Old: []any{cmp.Recoverable}, New: []any{r.Recoverable}})
	}

	if r.MaxRetries != cmp.MaxRetries {
		out = append(out, &revisions.Change{Key: "maxRetries", Old: []any{cmp.MaxRetries}, New: []any{r.MaxRetries}})
	}

	return out
}

func (r *NgAutomationStep) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationStep) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationStepMeta) Clone() *NgAutomationStepMeta {
	dup := r
	dup.Visual = *r.Visual.Clone()

	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	if r.Extra != nil {
		dup.Extra = make(map[string]any, len(r.Extra))
		for k, v := range r.Extra {
			dup.Extra[k] = v
		}
	}

	return &dup
}

func (r NgAutomationStepMeta) Diff(cmp *NgAutomationStepMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationStepMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	for _, c := range r.Visual.Diff(&cmp.Visual) {
		c.Key = "visual." + c.Key
		out = append(out, c)
	}

	if (r.Icon == nil) != (cmp.Icon == nil) {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	} else if r.Icon != nil {
		for _, c := range r.Icon.Diff(cmp.Icon) {
			c.Key = "icon." + c.Key
			out = append(out, c)
		}
	}

	if !reflect.DeepEqual(r.Extra, cmp.Extra) {
		out = append(out, &revisions.Change{Key: "extra", Old: []any{cmp.Extra}, New: []any{r.Extra}})
	}

	return out
}

func (r *NgAutomationStepMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationStepMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationPath) Clone() *NgAutomationPath {
	dup := r
	if r.Condition != nil {
		v := *r.Condition
		dup.Condition = &v
	}

	dup.Meta = *r.Meta.Clone()

	return &dup
}

func (r NgAutomationPath) Diff(cmp *NgAutomationPath) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationPath{}
	}
	if r.ParentID != cmp.ParentID {
		out = append(out, &revisions.Change{Key: "parentID", Old: []any{cmp.ParentID}, New: []any{r.ParentID}})
	}

	if r.ChildID != cmp.ChildID {
		out = append(out, &revisions.Change{Key: "childID", Old: []any{cmp.ChildID}, New: []any{r.ChildID}})
	}

	if !reflect.DeepEqual(r.Condition, cmp.Condition) {
		out = append(out, &revisions.Change{Key: "condition", Old: []any{cmp.Condition}, New: []any{r.Condition}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *NgAutomationPath) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationPath) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationPathMeta) Clone() *NgAutomationPathMeta {
	dup := r
	dup.Visual = *r.Visual.Clone()

	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	return &dup
}

func (r NgAutomationPathMeta) Diff(cmp *NgAutomationPathMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationPathMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	for _, c := range r.Visual.Diff(&cmp.Visual) {
		c.Key = "visual." + c.Key
		out = append(out, c)
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

func (r *NgAutomationPathMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationPathMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationIssue) Clone() *NgAutomationIssue {
	dup := r
	if r.Details != nil {
		dup.Details = make([]*NgAutomationIssueDetail, len(r.Details))
		copy(dup.Details, r.Details)
	}

	return &dup
}

func (r NgAutomationIssue) Diff(cmp *NgAutomationIssue) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationIssue{}
	}
	if r.Code != cmp.Code {
		out = append(out, &revisions.Change{Key: "code", Old: []any{cmp.Code}, New: []any{r.Code}})
	}

	if r.Severity != cmp.Severity {
		out = append(out, &revisions.Change{Key: "severity", Old: []any{cmp.Severity}, New: []any{r.Severity}})
	}

	if r.Message != cmp.Message {
		out = append(out, &revisions.Change{Key: "message", Old: []any{cmp.Message}, New: []any{r.Message}})
	}

	if !reflect.DeepEqual(r.Details, cmp.Details) {
		out = append(out, &revisions.Change{Key: "details", Old: []any{cmp.Details}, New: []any{r.Details}})
	}

	return out
}

func (r *NgAutomationIssue) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationIssue) Value() (driver.Value, error) { return json.Marshal(r) }

func (r NgAutomationIssueDetail) Clone() *NgAutomationIssueDetail {
	dup := r
	if r.MissingReference != nil {
		dup.MissingReference = r.MissingReference.Clone()
	}

	if r.InvalidType != nil {
		dup.InvalidType = r.InvalidType.Clone()
	}

	if r.DuplicateID != nil {
		dup.DuplicateID = r.DuplicateID.Clone()
	}

	if r.Cycle != nil {
		dup.Cycle = r.Cycle.Clone()
	}

	if r.ResourceRef != nil {
		dup.ResourceRef = r.ResourceRef.Clone()
	}

	if r.GatewayPaths != nil {
		dup.GatewayPaths = r.GatewayPaths.Clone()
	}

	if r.EmptyField != nil {
		dup.EmptyField = r.EmptyField.Clone()
	}

	if r.Details != nil {
		dup.Details = make([]*NgAutomationIssueDetail, len(r.Details))
		copy(dup.Details, r.Details)
	}

	return &dup
}

func (r NgAutomationIssueDetail) Diff(cmp *NgAutomationIssueDetail) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NgAutomationIssueDetail{}
	}
	if (r.MissingReference == nil) != (cmp.MissingReference == nil) {
		out = append(out, &revisions.Change{Key: "missingReference", Old: []any{cmp.MissingReference}, New: []any{r.MissingReference}})
	} else if r.MissingReference != nil {
		for _, c := range r.MissingReference.Diff(cmp.MissingReference) {
			c.Key = "missingReference." + c.Key
			out = append(out, c)
		}
	}

	if (r.InvalidType == nil) != (cmp.InvalidType == nil) {
		out = append(out, &revisions.Change{Key: "invalidType", Old: []any{cmp.InvalidType}, New: []any{r.InvalidType}})
	} else if r.InvalidType != nil {
		for _, c := range r.InvalidType.Diff(cmp.InvalidType) {
			c.Key = "invalidType." + c.Key
			out = append(out, c)
		}
	}

	if (r.DuplicateID == nil) != (cmp.DuplicateID == nil) {
		out = append(out, &revisions.Change{Key: "duplicateID", Old: []any{cmp.DuplicateID}, New: []any{r.DuplicateID}})
	} else if r.DuplicateID != nil {
		for _, c := range r.DuplicateID.Diff(cmp.DuplicateID) {
			c.Key = "duplicateID." + c.Key
			out = append(out, c)
		}
	}

	if (r.Cycle == nil) != (cmp.Cycle == nil) {
		out = append(out, &revisions.Change{Key: "cycle", Old: []any{cmp.Cycle}, New: []any{r.Cycle}})
	} else if r.Cycle != nil {
		for _, c := range r.Cycle.Diff(cmp.Cycle) {
			c.Key = "cycle." + c.Key
			out = append(out, c)
		}
	}

	if (r.ResourceRef == nil) != (cmp.ResourceRef == nil) {
		out = append(out, &revisions.Change{Key: "resourceRef", Old: []any{cmp.ResourceRef}, New: []any{r.ResourceRef}})
	} else if r.ResourceRef != nil {
		for _, c := range r.ResourceRef.Diff(cmp.ResourceRef) {
			c.Key = "resourceRef." + c.Key
			out = append(out, c)
		}
	}

	if (r.GatewayPaths == nil) != (cmp.GatewayPaths == nil) {
		out = append(out, &revisions.Change{Key: "gatewayPaths", Old: []any{cmp.GatewayPaths}, New: []any{r.GatewayPaths}})
	} else if r.GatewayPaths != nil {
		for _, c := range r.GatewayPaths.Diff(cmp.GatewayPaths) {
			c.Key = "gatewayPaths." + c.Key
			out = append(out, c)
		}
	}

	if (r.EmptyField == nil) != (cmp.EmptyField == nil) {
		out = append(out, &revisions.Change{Key: "emptyField", Old: []any{cmp.EmptyField}, New: []any{r.EmptyField}})
	} else if r.EmptyField != nil {
		for _, c := range r.EmptyField.Diff(cmp.EmptyField) {
			c.Key = "emptyField." + c.Key
			out = append(out, c)
		}
	}

	if !reflect.DeepEqual(r.Details, cmp.Details) {
		out = append(out, &revisions.Change{Key: "details", Old: []any{cmp.Details}, New: []any{r.Details}})
	}

	return out
}

func (r *NgAutomationIssueDetail) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r NgAutomationIssueDetail) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailMissingReference) Clone() *DetailMissingReference {
	dup := r
	return &dup
}

func (r DetailMissingReference) Diff(cmp *DetailMissingReference) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailMissingReference{}
	}
	if r.RefKind != cmp.RefKind {
		out = append(out, &revisions.Change{Key: "refKind", Old: []any{cmp.RefKind}, New: []any{r.RefKind}})
	}

	if r.Ref != cmp.Ref {
		out = append(out, &revisions.Change{Key: "ref", Old: []any{cmp.Ref}, New: []any{r.Ref}})
	}

	if r.StepID != cmp.StepID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.StepID}, New: []any{r.StepID}})
	}

	if r.Field != cmp.Field {
		out = append(out, &revisions.Change{Key: "field", Old: []any{cmp.Field}, New: []any{r.Field}})
	}

	if r.FieldIndex != cmp.FieldIndex {
		out = append(out, &revisions.Change{Key: "fieldIndex", Old: []any{cmp.FieldIndex}, New: []any{r.FieldIndex}})
	}

	return out
}

func (r *DetailMissingReference) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailMissingReference) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailInvalidType) Clone() *DetailInvalidType {
	dup := r
	return &dup
}

func (r DetailInvalidType) Diff(cmp *DetailInvalidType) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailInvalidType{}
	}
	if r.StepID != cmp.StepID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.StepID}, New: []any{r.StepID}})
	}

	if r.Field != cmp.Field {
		out = append(out, &revisions.Change{Key: "field", Old: []any{cmp.Field}, New: []any{r.Field}})
	}

	if r.FieldIndex != cmp.FieldIndex {
		out = append(out, &revisions.Change{Key: "fieldIndex", Old: []any{cmp.FieldIndex}, New: []any{r.FieldIndex}})
	}

	if r.Target != cmp.Target {
		out = append(out, &revisions.Change{Key: "target", Old: []any{cmp.Target}, New: []any{r.Target}})
	}

	if r.Expected != cmp.Expected {
		out = append(out, &revisions.Change{Key: "expected", Old: []any{cmp.Expected}, New: []any{r.Expected}})
	}

	if r.Actual != cmp.Actual {
		out = append(out, &revisions.Change{Key: "actual", Old: []any{cmp.Actual}, New: []any{r.Actual}})
	}

	return out
}

func (r *DetailInvalidType) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailInvalidType) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailDuplicateID) Clone() *DetailDuplicateID {
	dup := r
	if r.Indices != nil {
		dup.Indices = make([]int, len(r.Indices))
		copy(dup.Indices, r.Indices)
	}

	return &dup
}

func (r DetailDuplicateID) Diff(cmp *DetailDuplicateID) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailDuplicateID{}
	}
	if r.Resource != cmp.Resource {
		out = append(out, &revisions.Change{Key: "resource", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "id", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if !reflect.DeepEqual(r.Indices, cmp.Indices) {
		out = append(out, &revisions.Change{Key: "indices", Old: []any{cmp.Indices}, New: []any{r.Indices}})
	}

	return out
}

func (r *DetailDuplicateID) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailDuplicateID) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailCycle) Clone() *DetailCycle {
	dup := r
	return &dup
}

func (r DetailCycle) Diff(cmp *DetailCycle) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailCycle{}
	}
	if !reflect.DeepEqual(r.StepIDs, cmp.StepIDs) {
		out = append(out, &revisions.Change{Key: "stepIDs", Old: []any{cmp.StepIDs}, New: []any{r.StepIDs}})
	}

	return out
}

func (r *DetailCycle) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailCycle) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailResourceRef) Clone() *DetailResourceRef {
	dup := r
	return &dup
}

func (r DetailResourceRef) Diff(cmp *DetailResourceRef) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailResourceRef{}
	}
	if r.Resource != cmp.Resource {
		out = append(out, &revisions.Change{Key: "resource", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "id", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Index != cmp.Index {
		out = append(out, &revisions.Change{Key: "index", Old: []any{cmp.Index}, New: []any{r.Index}})
	}

	if r.Field != cmp.Field {
		out = append(out, &revisions.Change{Key: "field", Old: []any{cmp.Field}, New: []any{r.Field}})
	}

	if r.FieldIndex != cmp.FieldIndex {
		out = append(out, &revisions.Change{Key: "fieldIndex", Old: []any{cmp.FieldIndex}, New: []any{r.FieldIndex}})
	}

	return out
}

func (r *DetailResourceRef) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailResourceRef) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailGatewayPaths) Clone() *DetailGatewayPaths {
	dup := r
	return &dup
}

func (r DetailGatewayPaths) Diff(cmp *DetailGatewayPaths) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailGatewayPaths{}
	}
	if r.StepID != cmp.StepID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.StepID}, New: []any{r.StepID}})
	}

	if r.Violation != cmp.Violation {
		out = append(out, &revisions.Change{Key: "violation", Old: []any{cmp.Violation}, New: []any{r.Violation}})
	}

	if r.Got != cmp.Got {
		out = append(out, &revisions.Change{Key: "got", Old: []any{cmp.Got}, New: []any{r.Got}})
	}

	if r.Want != cmp.Want {
		out = append(out, &revisions.Change{Key: "want", Old: []any{cmp.Want}, New: []any{r.Want}})
	}

	return out
}

func (r *DetailGatewayPaths) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailGatewayPaths) Value() (driver.Value, error) { return json.Marshal(r) }

func (r DetailEmptyField) Clone() *DetailEmptyField {
	dup := r
	return &dup
}

func (r DetailEmptyField) Diff(cmp *DetailEmptyField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &DetailEmptyField{}
	}
	if r.Resource != cmp.Resource {
		out = append(out, &revisions.Change{Key: "resource", Old: []any{cmp.Resource}, New: []any{r.Resource}})
	}

	if r.Index != cmp.Index {
		out = append(out, &revisions.Change{Key: "index", Old: []any{cmp.Index}, New: []any{r.Index}})
	}

	if r.Field != cmp.Field {
		out = append(out, &revisions.Change{Key: "field", Old: []any{cmp.Field}, New: []any{r.Field}})
	}

	return out
}

func (r *DetailEmptyField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r DetailEmptyField) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseNgAutomationMeta(ss []string) (p *NgAutomationMeta, err error) {
	p = &NgAutomationMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *NgAutomationTriggerSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationTriggerSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationTriggerSet(ss []string) (p NgAutomationTriggerSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationStepSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationStepSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationStepSet(ss []string) (p NgAutomationStepSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationPathSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationPathSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationPathSet(ss []string) (p NgAutomationPathSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationIssueSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationIssueSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationIssueSet(ss []string) (p NgAutomationIssueSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
