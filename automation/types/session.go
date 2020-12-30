package types

import (
	"time"
)

type (
	// Instance of single workflow execution
	Session struct {
		ID         uint64 `json:"sessionID,string"`
		WorkflowID uint64 `json:"workflowID,string"`

		EventType    string `json:"eventType,string"`
		ResourceType string `json:"resourceType,string"`

		ExecutedAs uint64 `json:"executedBy,string"` // runner (might be different then creator)
		WallTime   int    `json:"wallTime"`          // how long did it take (ms) to run it (inc all suspension)
		UserTime   int    `json:"userTime"`          // how long did it take (ms) to run it (sum of all time spent in each step)

		Input  Variables `json:"input"`
		Output Variables `json:"output"`

		Trace []SessionTraceStep `json:"trace"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		PurgeAt   *time.Time `json:"purgeAt,omitempty"`
	}

	SessionFilter struct {
		TriggerID  []uint64 `json:"triggerID"`
		WorkflowID []uint64 `json:"workflowID"`

		EventType    string `json:"eventType"`
		ResourceType string `json:"resourceType"`
	}

	// WorkflowSessionTraceStep stores info and instrumentation on visited workflow steps
	SessionTraceStep struct {
		ID         uint64    `json:"traceStepID,string"`
		CallerStep uint64    `json:"traceCallerStepID,string"`
		WorkflowID uint64    `json:"workflowID,string"`
		StateID    uint64    `json:"stateID,string"`
		SessionID  uint64    `json:"sessionID,string"`
		CallerID   uint64    `json:"callerID,string"`
		StepID     uint64    `json:"stepID,string"`
		Depth      uint64    `json:"depth,string"`
		Scope      Variables `json:"scope"`
		Duration   int       `json:"duration"` // in ms
	}
)
