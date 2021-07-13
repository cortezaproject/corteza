package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
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

		// Stacktrace that gets stored (if/when configured)
		Stacktrace Stacktrace `json:"stacktrace"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		PurgeAt   *time.Time `json:"purgeAt,omitempty"`

		// here we join suspended & prompted state;
		// we treat both states as suspended
		SuspendedAt *time.Time `json:"suspendedAt,omitempty"`
		CompletedAt *time.Time `json:"completedAt,omitempty"`
		Error       string     `json:"error,omitempty"`

		session *wfexec.Session

		// For keeping runtime stacktrace,
		// even if we do not want to store it on every update
		//
		// This will aid us when session fails and we can access
		// the whole stacktrace
		RuntimeStacktrace Stacktrace `json:"-"`

		l sync.RWMutex
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
		CreatedBy    []uint64 `json:"createdBy"`
		EventType    string   `json:"eventType"`
		ResourceType string   `json:"resourceType"`

		Completed filter.State `json:"deleted"`
		Status    []uint       `json:"status"`

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

	SessionStatus uint
)

const (
	SessionStarted SessionStatus = iota
	SessionPrompted
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

func (s Session) Resume(ctx context.Context, stateID uint64, input *expr.Vars) (*wfexec.ResumedPrompt, error) {
	return s.session.Resume(ctx, stateID, input)
}

func (s Session) PendingPrompts(ownerId uint64) []*wfexec.PendingPrompt {
	return s.session.UserPendingPrompts(ownerId)
}

func (s *Session) GC() bool {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.CompletedAt != nil || s.session.Error() != nil
}

// WaitResults wait blocks until workflow session is completed or fails (or context is canceled) and returns resuts
func (s *Session) WaitResults(ctx context.Context) (*expr.Vars, wfexec.SessionStatus, Stacktrace, error) {
	s.l.RLock()
	defer s.l.RUnlock()

	if err := s.session.WaitUntil(ctx, wfexec.SessionFailed, wfexec.SessionCompleted); err != nil {
		return nil, -1, s.Stacktrace, err
	}

	return s.session.Result(), s.session.Status(), s.Stacktrace, nil
}

func (s *Session) Apply(ssp SessionStartParams) {
	s.l.Lock()
	defer s.l.Unlock()

	s.WorkflowID = ssp.WorkflowID
	s.EventType = ssp.EventType
	s.ResourceType = ssp.ResourceType
	s.Input = ssp.Input

	if ssp.KeepFor > 0 {
		at := time.Now().Add(time.Duration(ssp.KeepFor) * time.Second)
		s.PurgeAt = &at
	}

	if ssp.Trace {
		// set Stacktrace prop to signal status handler
		// that we're interested in storing stacktrace
		s.Stacktrace = Stacktrace{}
	}
}

func (s *Session) CopyRuntimeStacktrace() {
	s.l.Lock()
	defer s.l.Unlock()

	if s.Stacktrace != nil || s.Error != "" {
		// Save stacktrace when we know we're tracing workflows OR whenever there is an error...
		s.Stacktrace = s.RuntimeStacktrace
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
			return fmt.Errorf("cannot scan '%v' into Stacktrace: %w", string(b), err)
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
