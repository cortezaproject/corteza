package types

import (
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"time"
)

type (
	// Workflow represents entire workflow definition
	Workflow struct {
		ID      uint64            `json:"workflowID,string"`
		Handle  string            `json:"handle"`
		Labels  map[string]string `json:"labels,omitempty"`
		Meta    WorkflowMeta      `json:"meta"`
		Enabled bool              `json:"enabled"`

		Trace        bool          `json:"trace"`
		KeepSessions time.Duration `json:"keepSessions"`

		// Initial input scope
		Scope workflow.Variables `json:"scope"`

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

	WorkflowFilter struct{}

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
		Input workflow.Variables

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

		Triggered   string        `json:"triggered,string"`   // event name
		TriggeredBy string        `json:"triggeredBy,string"` // resource ID
		ExecutedAs  uint64        `json:"executedBy,string"`  // runner (might be different then creator)
		WallTime    time.Duration `json:"wallTime"`           // how long did it take to run it (inc all suspension)
		UserTime    time.Duration `json:"userTime"`           // how long did it take to run it (sum of all time spent in each step)

		Input  workflow.Variables `json:"input"`
		Output workflow.Variables `json:"output"`

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
		ID         uint64             `json:"traceStepID,string"`
		CallerStep uint64             `json:"traceCallerStepID,string"`
		WorkflowID uint64             `json:"workflowID,string"`
		StateID    uint64             `json:"stateID,string"`
		SessionID  uint64             `json:"sessionID,string"`
		CallerID   uint64             `json:"callerID,string"`
		StepID     uint64             `json:"stepID,string"`
		Depth      uint64             `json:"depth,string"`
		Scope      workflow.Variables `json:"scope"`
		Duration   time.Duration      `json:"duration"`
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

		CallerID uint64             `json:"callerID,string"`
		StepID   uint64             `json:"stepID,string"`
		Scope    workflow.Variables `json:"scope"`
	}

	// workflow functions are defined in the core code and through plugins
	WorkflowFunction struct {
		Ref        string                   `json:"ref"`
		Meta       WorkflowFunctionMeta     `json:"meta"`
		Handler    workflow.ActivityHandler `json:"-"`
		Parameters []*WorkflowParameter     `json:"parameters"`
		Results    []*WorkflowParameter     `json:"results"`
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
)
