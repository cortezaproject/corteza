package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/expr"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	Trigger struct {
		ID           uint64                           `json:"triggerID,string"`
		TenantID     uint64                           `json:"tenantID,string,omitempty"`
		ProjectID    uint64                           `json:"projectID,string,omitempty"`
		WorkflowID   uint64                           `json:"workflowID,string"`
		StepID       uint64                           `json:"stepID,string"`
		Enabled      bool                             `json:"enabled"`
		Meta         *TriggerMeta                     `json:"meta,omitempty"`
		ResourceType string                           `json:"resourceType"`
		EventType    string                           `json:"eventType"`
		Constraints  TriggerConstraintSet             `json:"constraints"`
		Input        *expr.Vars                       `json:"input"`
		OwnedBy      uint64                           `json:"ownedBy,string"`
		CreatedAt    time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy    uint64                           `json:"createdBy,string"`
		UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
		Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	TriggerMeta struct {
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	TriggerConstraint struct {
		Name   string   `json:"name"`
		Op     string   `json:"op,omitempty"`
		Values []string `json:"values,omitempty"`
	}
)

func (r Trigger) Clone() *Trigger {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.Input != nil {
		v := *r.Input
		dup.Input = &v
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

func (r Trigger) Diff(cmp *Trigger) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Trigger{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "triggerID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.WorkflowID != cmp.WorkflowID {
		out = append(out, &revisions.Change{Key: "workflowID", Old: []any{cmp.WorkflowID}, New: []any{r.WorkflowID}})
	}

	if r.StepID != cmp.StepID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.StepID}, New: []any{r.StepID}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
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

func (r TriggerMeta) Clone() *TriggerMeta {
	dup := r
	if r.Visual != nil {
		dup.Visual = make(map[string]interface{}, len(r.Visual))
		for k, v := range r.Visual {
			dup.Visual[k] = v
		}
	}

	return &dup
}

func (r TriggerMeta) Diff(cmp *TriggerMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TriggerMeta{}
	}
	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.Visual, cmp.Visual) {
		out = append(out, &revisions.Change{Key: "visual", Old: []any{cmp.Visual}, New: []any{r.Visual}})
	}

	return out
}

func (r *TriggerMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TriggerMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r TriggerConstraint) Clone() *TriggerConstraint {
	dup := r
	if r.Values != nil {
		dup.Values = make([]string, len(r.Values))
		copy(dup.Values, r.Values)
	}

	return &dup
}

func (r TriggerConstraint) Diff(cmp *TriggerConstraint) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TriggerConstraint{}
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

func (r *TriggerConstraint) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r TriggerConstraint) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseTriggerMeta(ss []string) (p *TriggerMeta, err error) {
	p = &TriggerMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *TriggerConstraintSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TriggerConstraintSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTriggerConstraintSet(ss []string) (p TriggerConstraintSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
