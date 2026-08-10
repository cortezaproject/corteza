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
	Workflow struct {
		ID           uint64                           `json:"workflowID,string"`
		TenantID     uint64                           `json:"tenantID,string,omitempty"`
		ProjectID    uint64                           `json:"projectID,string,omitempty"`
		Handle       string                           `json:"handle"`
		Meta         *WorkflowMeta                    `json:"meta,omitempty"`
		Enabled      bool                             `json:"enabled"`
		Trace        bool                             `json:"trace"`
		KeepSessions int                              `json:"keepSessions"`
		Scope        *expr.Vars                       `json:"scope"`
		Steps        WorkflowStepSet                  `json:"steps"`
		Paths        WorkflowPathSet                  `json:"paths"`
		Issues       WorkflowIssueSet                 `json:"issues,omitempty"`
		RunAs        uint64                           `json:"runAs,string"`
		OwnedBy      uint64                           `json:"ownedBy,string"`
		CreatedAt    time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy    uint64                           `json:"createdBy,string"`
		UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
		Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	WorkflowMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
		SubWorkflow bool                   `json:"subWorkflow,omitempty"`
		Input       []WorkflowIODef        `json:"input,omitempty"`
		Output      []WorkflowIODef        `json:"output,omitempty"`
	}

	WorkflowStep struct {
		ID        uint64            `json:"stepID,string"`
		Kind      WorkflowStepKind  `json:"kind"`
		Ref       string            `json:"ref"`
		Arguments []*Expr           `json:"arguments"`
		Results   []*Expr           `json:"results"`
		Meta      WorkflowStepMeta  `json:"meta,omitempty"`
		Labels    map[string]string `json:"labels,omitempty"`
	}

	WorkflowStepMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowPathMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowIssue struct {
		Culprit     map[string]int `json:"culprit"`
		Description string         `json:"description"`
	}

	WorkflowExecParams struct {
		CallerWorkflowID uint64     `json:"callerWorkflowID"`
		CallerSessionID  uint64     `json:"callerSessionID"`
		CallerStepID     uint64     `json:"callerStepID"`
		StepID           uint64     `json:"stepID"`
		EventType        string     `json:"eventType"`
		ResourceType     string     `json:"resourceType"`
		Trace            bool       `json:"trace"`
		Async            bool       `json:"async"`
		Wait             bool       `json:"wait"`
		Input            *expr.Vars `json:"input"`
	}

	WorkflowStepKind string
)

func (r Workflow) Clone() *Workflow {
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

func (r Workflow) Diff(cmp *Workflow) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Workflow{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "workflowID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if r.Trace != cmp.Trace {
		out = append(out, &revisions.Change{Key: "trace", Old: []any{cmp.Trace}, New: []any{r.Trace}})
	}

	if r.KeepSessions != cmp.KeepSessions {
		out = append(out, &revisions.Change{Key: "keepSessions", Old: []any{cmp.KeepSessions}, New: []any{r.KeepSessions}})
	}

	if !reflect.DeepEqual(r.Scope, cmp.Scope) {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
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

func (r WorkflowMeta) Clone() *WorkflowMeta {
	dup := r
	if r.Visual != nil {
		dup.Visual = make(map[string]interface{}, len(r.Visual))
		for k, v := range r.Visual {
			dup.Visual[k] = v
		}
	}

	if r.Input != nil {
		dup.Input = make([]WorkflowIODef, len(r.Input))
		copy(dup.Input, r.Input)
	}

	if r.Output != nil {
		dup.Output = make([]WorkflowIODef, len(r.Output))
		copy(dup.Output, r.Output)
	}

	return &dup
}

func (r WorkflowMeta) Diff(cmp *WorkflowMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.Visual, cmp.Visual) {
		out = append(out, &revisions.Change{Key: "visual", Old: []any{cmp.Visual}, New: []any{r.Visual}})
	}

	if r.SubWorkflow != cmp.SubWorkflow {
		out = append(out, &revisions.Change{Key: "subWorkflow", Old: []any{cmp.SubWorkflow}, New: []any{r.SubWorkflow}})
	}

	if !reflect.DeepEqual(r.Input, cmp.Input) {
		out = append(out, &revisions.Change{Key: "input", Old: []any{cmp.Input}, New: []any{r.Input}})
	}

	if !reflect.DeepEqual(r.Output, cmp.Output) {
		out = append(out, &revisions.Change{Key: "output", Old: []any{cmp.Output}, New: []any{r.Output}})
	}

	return out
}

func (r *WorkflowMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r WorkflowStep) Clone() *WorkflowStep {
	dup := r
	if r.Arguments != nil {
		dup.Arguments = make([]*Expr, len(r.Arguments))
		copy(dup.Arguments, r.Arguments)
	}

	if r.Results != nil {
		dup.Results = make([]*Expr, len(r.Results))
		copy(dup.Results, r.Results)
	}

	dup.Meta = *r.Meta.Clone()

	if r.Labels != nil {
		dup.Labels = make(map[string]string, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}

	return &dup
}

func (r WorkflowStep) Diff(cmp *WorkflowStep) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowStep{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if !reflect.DeepEqual(r.Kind, cmp.Kind) {
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

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}

	return out
}

func (r *WorkflowStep) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowStep) Value() (driver.Value, error) { return json.Marshal(r) }

func (r WorkflowStepMeta) Clone() *WorkflowStepMeta {
	dup := r
	if r.Visual != nil {
		dup.Visual = make(map[string]interface{}, len(r.Visual))
		for k, v := range r.Visual {
			dup.Visual[k] = v
		}
	}

	return &dup
}

func (r WorkflowStepMeta) Diff(cmp *WorkflowStepMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowStepMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.Visual, cmp.Visual) {
		out = append(out, &revisions.Change{Key: "visual", Old: []any{cmp.Visual}, New: []any{r.Visual}})
	}

	return out
}

func (r *WorkflowStepMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowStepMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r WorkflowPathMeta) Clone() *WorkflowPathMeta {
	dup := r
	if r.Visual != nil {
		dup.Visual = make(map[string]interface{}, len(r.Visual))
		for k, v := range r.Visual {
			dup.Visual[k] = v
		}
	}

	return &dup
}

func (r WorkflowPathMeta) Diff(cmp *WorkflowPathMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowPathMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.Visual, cmp.Visual) {
		out = append(out, &revisions.Change{Key: "visual", Old: []any{cmp.Visual}, New: []any{r.Visual}})
	}

	return out
}

func (r *WorkflowPathMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowPathMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r WorkflowIssue) Clone() *WorkflowIssue {
	dup := r
	if r.Culprit != nil {
		dup.Culprit = make(map[string]int, len(r.Culprit))
		for k, v := range r.Culprit {
			dup.Culprit[k] = v
		}
	}

	return &dup
}

func (r WorkflowIssue) Diff(cmp *WorkflowIssue) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowIssue{}
	}
	if !reflect.DeepEqual(r.Culprit, cmp.Culprit) {
		out = append(out, &revisions.Change{Key: "culprit", Old: []any{cmp.Culprit}, New: []any{r.Culprit}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *WorkflowIssue) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowIssue) Value() (driver.Value, error) { return json.Marshal(r) }

func (r WorkflowExecParams) Clone() *WorkflowExecParams {
	dup := r
	if r.Input != nil {
		v := *r.Input
		dup.Input = &v
	}

	return &dup
}

func (r WorkflowExecParams) Diff(cmp *WorkflowExecParams) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &WorkflowExecParams{}
	}
	if r.CallerWorkflowID != cmp.CallerWorkflowID {
		out = append(out, &revisions.Change{Key: "callerWorkflowID", Old: []any{cmp.CallerWorkflowID}, New: []any{r.CallerWorkflowID}})
	}

	if r.CallerSessionID != cmp.CallerSessionID {
		out = append(out, &revisions.Change{Key: "callerSessionID", Old: []any{cmp.CallerSessionID}, New: []any{r.CallerSessionID}})
	}

	if r.CallerStepID != cmp.CallerStepID {
		out = append(out, &revisions.Change{Key: "callerStepID", Old: []any{cmp.CallerStepID}, New: []any{r.CallerStepID}})
	}

	if r.StepID != cmp.StepID {
		out = append(out, &revisions.Change{Key: "stepID", Old: []any{cmp.StepID}, New: []any{r.StepID}})
	}

	if r.EventType != cmp.EventType {
		out = append(out, &revisions.Change{Key: "eventType", Old: []any{cmp.EventType}, New: []any{r.EventType}})
	}

	if r.ResourceType != cmp.ResourceType {
		out = append(out, &revisions.Change{Key: "resourceType", Old: []any{cmp.ResourceType}, New: []any{r.ResourceType}})
	}

	if r.Trace != cmp.Trace {
		out = append(out, &revisions.Change{Key: "trace", Old: []any{cmp.Trace}, New: []any{r.Trace}})
	}

	if r.Async != cmp.Async {
		out = append(out, &revisions.Change{Key: "async", Old: []any{cmp.Async}, New: []any{r.Async}})
	}

	if r.Wait != cmp.Wait {
		out = append(out, &revisions.Change{Key: "wait", Old: []any{cmp.Wait}, New: []any{r.Wait}})
	}

	if !reflect.DeepEqual(r.Input, cmp.Input) {
		out = append(out, &revisions.Change{Key: "input", Old: []any{cmp.Input}, New: []any{r.Input}})
	}

	return out
}

func (r *WorkflowExecParams) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r WorkflowExecParams) Value() (driver.Value, error) { return json.Marshal(r) }

const (
	WorkflowStepKindExpressions  WorkflowStepKind = "expressions"
	WorkflowStepKindGateway      WorkflowStepKind = "gateway"
	WorkflowStepKindFunction     WorkflowStepKind = "function"
	WorkflowStepKindIterator     WorkflowStepKind = "iterator"
	WorkflowStepKindError        WorkflowStepKind = "error"
	WorkflowStepKindTermination  WorkflowStepKind = "termination"
	WorkflowStepKindPrompt       WorkflowStepKind = "prompt"
	WorkflowStepKindDelay        WorkflowStepKind = "delay"
	WorkflowStepKindErrHandler   WorkflowStepKind = "error-handler"
	WorkflowStepKindVisual       WorkflowStepKind = "visual"
	WorkflowStepKindDebug        WorkflowStepKind = "debug"
	WorkflowStepKindBreak        WorkflowStepKind = "break"
	WorkflowStepKindContinue     WorkflowStepKind = "continue"
	WorkflowStepKindExecWorkflow WorkflowStepKind = "exec-workflow"
)

func ParseWorkflowMeta(ss []string) (p *WorkflowMeta, err error) {
	p = &WorkflowMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *WorkflowStepSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowStepSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowStepSet(ss []string) (p WorkflowStepSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *WorkflowPathSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowPathSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowPathSet(ss []string) (p WorkflowPathSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *WorkflowIssueSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowIssueSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowIssueSet(ss []string) (p WorkflowIssueSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
