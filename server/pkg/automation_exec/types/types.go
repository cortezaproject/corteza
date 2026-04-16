package types

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	Executable struct {
		ID       id.ID
		Revision int
		Handle   string

		RunAs auth.Identifiable

		Steps []Step

		// @todo
		// Limits        ExecutableLimits
	}

	// Expr represents an expression mapping for arguments/results
	// Cloned from automation/types.Expr to avoid circular dependencies
	Expr struct {
		ArgumentName string `json:"argumentName"`

		// Scope defines where evaluation context of this expression
		//
		// Leave empty for global context
		Scope string `json:"scope,omitempty"`

		// Variable name to set results of the expression to
		Target string `json:"target"`

		// Source of the value / name of the variable from scope
		// Takes precedence over Value
		Source string `json:"source,omitempty"`

		// Expression to evaluate over the input variables
		Expression string `json:"expr,omitempty"`

		// Raw value to be set to target
		// If expression is set and fails, evaluation defaults to value
		Value interface{} `json:"value,omitempty"`

		// Expected type of the input value
		Type string `json:"type,omitempty"`
	}

	StepArg struct {
		*Expr
	}

	StepRst struct {
		ArgumentName string
	}

	// Step represents a single step in an executable workflow
	Step struct {
		ID     id.ID
		Handle string
		Kind   string

		Children []Step
		Parents  []Step

		Arguments []StepArg
		Results   []StepRst

		Handler StepHandler

		// Recoverable enables pause-and-retry behaviour when a phase returns a RecoverableError.
		// Both the step flag and the error type must be set for a pause to trigger.
		Recoverable bool

		// MaxRetries is the per-phase retry limit.
		// 0 = pause indefinitely (user must call ResumeStep or TerminateExecution).
		// >0 = auto-escalate to the catch handler after N failed retries.
		MaxRetries int

		// ErrHandlerStepID is the ID of the catch step registered for this step.
		// Copied into frame.errHandlerStepID when the frame is first pushed.
		// Zero means no catch handler.
		ErrHandlerStepID id.ID
	}

	ExecRequest struct {
		// Current scope
		Scope map[string]*expr.Vars

		// Evaluated arguments passed into this step, populated by the step's handler
		Arguments map[string]expr.TypedValue
	}

	ExecResponse any

	ExecutableLimits struct{}

	// Executable represents a complete workflow definition

	Execution struct {
		ID           id.ID
		ExecutableID id.ID
		Revision     int
		Status       Status
		CreatedAt    time.Time
		UpdatedAt    time.Time
		EndedAt      *time.Time

		EventType    string
		ResourceType string

		Error error

		Trace []StackFrame

		Events []StepEvent
	}

	StepEvent struct {
		Type      EventType
		StepID    id.ID
		Timestamp time.Time
		Payload   any
		Error     error
	}

	StepHandler interface {
		ExecN(context.Context, *ExecRequest) (ExecResponse, error)
	}

	Budget struct {
		MaxOps int
		Window time.Duration // 0 = hard cap
	}

	RateLimit struct {
		MaxOps int
		Window time.Duration
	}

	EventType        string
	Status           string
	ExecutableStatus int

	StackFrame struct {
		ID        id.ID      `json:"id"`
		StepID    id.ID      `json:"stepID"`
		ParentID  id.ID      `json:"parentID,omitempty"`
		Handle    string     `json:"handle,omitempty"`
		Kind      string     `json:"kind,omitempty"`
		Input     any        `json:"input,omitempty"`
		Output    any        `json:"output,omitempty"`
		Args      any        `json:"args,omitempty"`
		StartedAt time.Time  `json:"startedAt"`
		EndedAt   *time.Time `json:"endedAt,omitempty"`
		Error     error      `json:"error,omitempty"`
	}

	ExecutionParams struct {
		EntryPoint   string
		Input        *expr.Vars
		EventType    string
		ResourceType string
	}

	ExecutionFilter struct {
		AutomationID []string
		EventType    string
		ResourceType string
		Status       []string
	}

	ExecutionResult struct {
		ExecutionID  id.ID      `json:"executionID"`
		ExecutableID id.ID      `json:"executableID"`
		Revision     int        `json:"revision"`
		Status       Status     `json:"status"`
		Error        string     `json:"error,omitempty"`
		StartedAt    time.Time  `json:"startedAt"`
		EndedAt      *time.Time `json:"endedAt,omitempty"`
		Duration     string     `json:"duration"`
		EventType    string     `json:"eventType,omitempty"`
		ResourceType string     `json:"resourceType,omitempty"`
	}
)

const (
	StatusActive ExecutableStatus = iota
	StatusDeprecated
)

const (
	MaxIteratorFrames = 1000
)

const (
	StatusCreated   Status = "created"
	StatusRunning   Status = "running"
	StatusPaused    Status = "paused"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

const (
	EventStepStarted   EventType = "step_started"
	EventStepCompleted EventType = "step_completed"
	EventStepFailed    EventType = "step_failed"
	EventStepPaused    EventType = "step_paused"
)

func (s ExecutableStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusDeprecated:
		return "deprecated"
	default:
		return "unknown"
	}
}
