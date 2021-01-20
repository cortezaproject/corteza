package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
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

		Input  *expr.Vars `json:"input"`
		Output *expr.Vars `json:"output"`

		Stacktrace Stacktrace `json:"stacktrace"`

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
		Input        *expr.Vars
		StepID       uint64
		EventType    string
		ResourceType string
	}

	SessionFilter struct {
		SessionID    []uint64 `json:"sessionID"`
		WorkflowID   []uint64 `json:"workflowID"`
		EventType    string   `json:"eventType"`
		ResourceType string   `json:"resourceType"`

		Completed filter.State `json:"deleted"`
		Suspended filter.State `json:"disabled"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Session) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	Stacktrace []*wfexec.Frame

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

func (s Session) Exec(ctx context.Context, step wfexec.Step, input *expr.Vars) error {
	return s.session.Exec(ctx, step, input)
}

func (s Session) Resume(ctx context.Context, stateID uint64, input *expr.Vars) error {
	return s.session.Resume(ctx, stateID, input)
}

func (s Session) Wait(ctx context.Context) error {
	return s.session.Wait(ctx)
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

	if ssp.Trace {
		// set prop
		s.Stacktrace = Stacktrace{}
	}
}

func (set *Stacktrace) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*set = Stacktrace{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, set); err != nil {
			return fmt.Errorf("can not scan '%v' into Stacktrace: %w", string(b), err)
		}
	}

	return nil
}

// Scan on WorkflowStepSet gracefully handles conversion from NULL
func (set Stacktrace) Value() (driver.Value, error) {
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
