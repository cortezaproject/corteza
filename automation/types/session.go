package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"time"
)

type (
	// Instance of single workflow execution
	Session struct {
		ID         uint64 `json:"sessionID,string"`
		WorkflowID uint64 `json:"workflowID,string"`

		Status SessionStatus `json:"status,string"`

		EventType    string `json:"eventType"`
		ResourceType string `json:"resourceType"`

		WallTime int `json:"wallTime"` // how long did it take (ms) to run it (inc all suspension)
		UserTime int `json:"userTime"` // how long did it take (ms) to run it (sum of all time spent in each step)

		Input  Variables `json:"input"`
		Output Variables `json:"output"`

		Trace SessionTraceStepSet `json:"trace"`

		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		CreatedBy   uint64     `json:"createdBy,string"`
		PurgeAt     *time.Time `json:"purgeAt,omitempty"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty"`
		CompletedAt *time.Time `json:"completedAt,omitempty"`
		Error       string     `json:"error,omitempty"`

		session *wfexec.Session
	}

	SessionStartParams struct {
		WorkflowID   uint64
		KeepFor      int
		Trace        bool
		Input        Variables
		StepID       uint64
		EventType    string
		ResourceType string
	}

	SessionFilter struct {
		SessionID    []uint64 `json:"sessionID"`
		WorkflowID   []uint64 `json:"workflowID"`
		EventType    string   `json:"eventType"`
		ResourceType string   `json:"resourceType"`
	}

	// WorkflowSessionTraceStep stores info and instrumentation on visited workflow steps
	SessionTraceStep struct {
		ID         uint64    `json:"traceStepID,string"`
		CallerStep uint64    `json:"traceCallerStepID,string"`
		StateID    uint64    `json:"stateID,string"`
		CallerID   uint64    `json:"callerID,string"`
		StepID     uint64    `json:"stepID,string"`
		Depth      uint64    `json:"depth,string"`
		Scope      Variables `json:"scope"`
		Duration   int       `json:"duration"` // in ms
	}

	SessionStatus int
)

const (
	SessionStarted SessionStatus = iota
	SessionSuspended
	SessionFailed
	SessionCompleted
)

func NewSession(s *wfexec.Session) *Session {
	return &Session{
		ID:      s.ID(),
		session: s,
	}
}

func (s Session) Exec(ctx context.Context, step wfexec.Step, input Variables) error {
	return s.session.Exec(ctx, step, wfexec.Variables(input))
}

func (s Session) Resume(ctx context.Context, stateID uint64, input Variables) error {
	return s.session.Resume(ctx, stateID, wfexec.Variables(input))
}

func (s *Session) Apply(ssp SessionStartParams) {
	s.WorkflowID = ssp.WorkflowID
	s.EventType = ssp.EventType
	s.ResourceType = ssp.ResourceType
	s.Input = ssp.Input

	if ssp.KeepFor > 0 {
		at := time.Now().Add(time.Duration(ssp.KeepFor) * time.Second)
		s.PurgeAt = &at
	}
}

func (set *SessionTraceStepSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*set = SessionTraceStepSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, set); err != nil {
			return fmt.Errorf("can not scan '%v' into SessionTraceStepSet: %w", string(b), err)
		}
	}

	return nil
}

// Scan on WorkflowStepSet gracefully handles conversion from NULL
func (set SessionTraceStepSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (s SessionStatus) String() string {
	switch s {
	case SessionStarted:
		return "started"
	case SessionSuspended:
		return "suspended"
	case SessionFailed:
		return "failed"
	case SessionCompleted:
		return "completed"
	}

	return "unknown"
}

func (s SessionStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
