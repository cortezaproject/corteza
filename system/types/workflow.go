package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"time"
)

type (
	// Workflow represents entire workflow definition
	Workflow struct {
		ID      uint64            `json:"workflowID,string"`
		Handle  string            `json:"handle"`
		Labels  map[string]string `json:"labels,omitempty"`
		Meta    *WorkflowMeta     `json:"meta"`
		Enabled bool              `json:"enabled"`

		Trace bool `json:"trace"`

		// how much time do we keep completed sessions (in sec)
		KeepSessions int `json:"keepSessions"`

		// Initial input scope
		Scope wfexec.Variables `json:"scope"`

		Steps    WorkflowStepSet
		Paths    WorkflowPathSet
		Triggers WorkflowTriggerSet

		RunAs uint64 `json:"runAs,string"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	WorkflowFilter struct {
		WorkflowID []uint64 `json:"workflowID"`

		Query string `json:"query"`

		Deleted  filter.State `json:"deleted"`
		Archived filter.State `json:"archived"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Workflow) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	WorkflowMeta struct {
		Name        string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowTrigger struct {
		ID         uint64 `json:"triggerID,string"`
		WorkflowID uint64 `json:"workflowID,string"`
		Enabled    bool   `json:"enabled"`

		// Start workflow on this step. If 0, find first (only) orphan
		StepID uint64

		// Resource type that can trigger the workflow
		ResourceType string

		// Event type that can trigger the workflow
		EventType string

		// Trigger constraints
		Constraints []WorkflowTriggerConstraint

		Meta WorkflowTriggerMeta `json:"meta"`

		// Initial input scope,
		// will be merged merged with workflow variables
		Input wfexec.Variables

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	WorkflowTriggerConstraint struct {
		Name   string   `json:"name"`
		Op     string   `json:"op,omitempty"`
		Values []string `json:"values,omitempty"`
	}

	WorkflowTriggerFilter struct{}

	WorkflowTriggerMeta struct {
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// WorkflowStep describes one workflow step
	WorkflowStep struct {
		ID   uint64           `json:"stepID,string"`
		Kind WorkflowStepKind `json:"kind"`

		// reference to function or subprocess (workflow)
		Ref string `json:"ref"`

		// set of expressions to evaluate or pass to function
		// invalid for for kind=~gateway:*
		Arguments []*WorkflowExpression `json:"arguments,string"`

		// only valid when kind=function
		Results []*WorkflowExpression `json:"results,string"`

		Meta WorkflowStepMeta `json:"meta"`
	}

	WorkflowStepMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// WorkflowPath defines connection between two workflow steps
	WorkflowPath struct {
		ParentID uint64 `json:"parentID,string"`
		ChildID  uint64 `json:"childID,string"`

		// test expression for gateway paths
		Test string           `json:"test,string"`
		Meta WorkflowPathMeta `json:"meta"`
	}

	WorkflowPathMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// Instance of single workflow execution
	WorkflowSession struct {
		ID         uint64 `json:"sessionID,string"`
		WorkflowID uint64 `json:"workflowID,string"`

		EventType       string `json:"eventType,string"`       // event name
		EventResourceID string `json:"eventResourceID,string"` // resource ID
		ExecutedAs      uint64 `json:"executedBy,string"`      // runner (might be different then creator)
		WallTime        int    `json:"wallTime"`               // how long did it take (ms) to run it (inc all suspension)
		UserTime        int    `json:"userTime"`               // how long did it take (ms) to run it (sum of all time spent in each step)

		Input  wfexec.Variables `json:"input"`
		Output wfexec.Variables `json:"output"`

		Trace []WorkflowSessionTraceStep `json:"trace"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
		PurgeAt   *time.Time `json:"purgeAt,omitempty"`
	}

	WorkflowSessionFilter struct{}

	// WorkflowSessionTraceStep stores info and instrumentation on visited workflow steps
	WorkflowSessionTraceStep struct {
		ID         uint64           `json:"traceStepID,string"`
		CallerStep uint64           `json:"traceCallerStepID,string"`
		WorkflowID uint64           `json:"workflowID,string"`
		StateID    uint64           `json:"stateID,string"`
		SessionID  uint64           `json:"sessionID,string"`
		CallerID   uint64           `json:"callerID,string"`
		StepID     uint64           `json:"stepID,string"`
		Depth      uint64           `json:"depth,string"`
		Scope      wfexec.Variables `json:"scope"`
		Duration   int              `json:"duration"` // in ms
	}

	// WorkflowState tracks suspended sessions
	// Session can have more than one state
	WorkflowState struct {
		ID        uint64 `json:"stateID,string"`
		SessionID uint64 `json:"sessionID,string"`

		ResumeAt        *time.Time `json:"resumeAt"`
		WaitingForInput bool       `json:"waitingForInput"`

		CreatedAt time.Time `json:"createdAt,omitempty"`
		CreatedBy uint64    `json:"createdBy,string"`

		CallerID uint64           `json:"callerID,string"`
		StepID   uint64           `json:"stepID,string"`
		Scope    wfexec.Variables `json:"scope"`
	}

	// workflow functions are defined in the core code and through plugins
	WorkflowFunction struct {
		Ref        string                 `json:"ref"`
		Meta       WorkflowFunctionMeta   `json:"meta"`
		Handler    wfexec.ActivityHandler `json:"-"`
		Parameters []*WorkflowParameter   `json:"parameters"`
		Results    []*WorkflowParameter   `json:"results"`
	}

	WorkflowFunctionMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowParameter struct {
		Name string                `json:"name"`
		Type string                `json:"type"`
		Meta WorkflowParameterMeta `json:"meta"`
	}

	WorkflowParameterMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// Used for expression steps, arguments and results mapping
	WorkflowExpression struct {
		Name string `json:"name"`
		Expr string `json:"expr"`
	}

	WorkflowStepKind string
)

const (
	WorkflowStepKindExpressions WorkflowStepKind = "expressions"
	WorkflowStepKindGatewayIncl WorkflowStepKind = "gateway:incl"
	WorkflowStepKindGatewayExcl WorkflowStepKind = "gateway:excl"
	WorkflowStepKindGatewayFork WorkflowStepKind = "gateway:fork"
	WorkflowStepKindGatewayJoin WorkflowStepKind = "gateway:join"
	WorkflowStepKindFunction    WorkflowStepKind = "function"
	WorkflowStepKindSubprocess  WorkflowStepKind = "subprocess"
	//WorkflowStepKindInput       WorkflowStepKind = "input" // ref = frontend function
	//WorkflowStepKindEvent       WorkflowStepKind = "event" // ref = ??
	//WorkflowStepKindAlert       WorkflowStepKind = "alert" // ref = error, warning, info
)

// Resource returns a resource ID for this type
func (r *Workflow) RBACResource() rbac.Resource {
	return WorkflowRBACResource.AppendID(r.ID)
}

func ParseWorkflowMeta(ss []string) (p *WorkflowMeta, err error) {
	p = &WorkflowMeta{}
	return p, parseStringsInput(ss, p)
}

func ParseWorkflowStepSet(ss []string) (p WorkflowStepSet, err error) {
	p = WorkflowStepSet{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowPathSet(ss []string) (p WorkflowPathSet, err error) {
	p = WorkflowPathSet{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowTriggerSet(ss []string) (p WorkflowTriggerSet, err error) {
	p = WorkflowTriggerSet{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowVariables(ss []string) (p wfexec.Variables, err error) {
	p = wfexec.Variables{}
	return p, parseStringsInput(ss, &p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
}

func (vv *WorkflowMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = WorkflowMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into WorkflowMeta: %w", string(b), err)
		}
	}

	return nil
}

// Scan on WorkflowMeta gracefully handles conversion from NULL
func (vv *WorkflowMeta) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}
