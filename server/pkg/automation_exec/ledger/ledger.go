package ledger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	ledger struct {
		mu    sync.RWMutex
		store map[id.ID]*types.Execution
	}
)

// The Ledger is the source of truth for what happened.
//
// Records executions
// Records step events
// Tracks current and terminal state
// Answers inspection questions (“what’s running?”, “what finished?”, “what step failed?”)
func Ledger() *ledger {
	return &ledger{
		store: make(map[id.ID]*types.Execution),
	}
}

func (l *ledger) Init(
	ctx context.Context,
	execID, exeID id.ID,
	rev uint32,
	params map[string]any,
) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.store[execID]; exists {
		return fmt.Errorf("execution %d already exists", execID)
	}

	l.store[execID] = &types.Execution{
		ID:           execID,
		ExecutableID: exeID,
		Revision:     rev,
		Status:       types.StatusCreated,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Variables:    params,
		Events:       make([]types.StepEvent, 0),
	}

	return nil
}

func (l *ledger) ExecutionCompleted(ctx context.Context, execID id.ID) error {
	return l.transition(execID, types.StatusCompleted)
}

func (l *ledger) ExecutionFailed(ctx context.Context, execID id.ID, err error) error {
	return l.transition(execID, types.StatusFailed)
}

func (l *ledger) transition(execID id.ID, status types.Status) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[execID]
	if !ok {
		return fmt.Errorf("execution %d not found", execID)
	}

	ex.Status = status
	now := time.Now()
	ex.UpdatedAt = now

	if isTerminal(status) {
		ex.EndedAt = &now
	}

	return nil
}

func (l *ledger) StepStarted(ctx context.Context, execID, stepID id.ID) error {
	return l.logStep(execID, types.StepEvent{
		Type:   types.EventStepStarted,
		StepID: stepID,
	})
}

func (l *ledger) StepCompleted(
	ctx context.Context,
	execID, stepID id.ID,
	result map[string]any,
) error {
	return l.logStep(execID, types.StepEvent{
		Type:   types.EventStepCompleted,
		StepID: stepID,
	})
}

func (l *ledger) StepFailed(
	ctx context.Context,
	execID, stepID id.ID,
	err error,
) error {
	return l.logStep(execID, types.StepEvent{
		Type:   types.EventStepFailed,
		StepID: stepID,
		Error:  err,
	})
}

func (l *ledger) logStep(execID id.ID, event types.StepEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[execID]
	if !ok {
		return fmt.Errorf("execution %d not found", execID)
	}

	event.Timestamp = time.Now()
	ex.Events = append(ex.Events, event)
	ex.UpdatedAt = event.Timestamp

	return nil
}

func (l *ledger) IsExecutableInUse(ctx context.Context, exeID id.ID) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, ex := range l.store {
		if ex.ExecutableID.Equal(exeID) && !isTerminal(ex.Status) {
			return true, nil
		}
	}
	return false, nil
}

func (l *ledger) IsStepInUse(ctx context.Context, exeID, stepID id.ID) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, ex := range l.store {
		if !ex.ExecutableID.Equal(exeID) || isTerminal(ex.Status) {
			continue
		}

		for i := len(ex.Events) - 1; i >= 0; i-- {
			ev := ex.Events[i]
			if ev.StepID.Equal(stepID) {
				return ev.Type == types.EventStepStarted, nil
			}
		}
	}

	return false, nil
}

func (l *ledger) GetExecution(ctx context.Context, execID id.ID) (*types.Execution, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ex, ok := l.store[execID]
	if !ok {
		return nil, fmt.Errorf("execution %d not found", execID)
	}

	return ex, nil
}

func (l *ledger) ListExecutions(ctx context.Context) ([]*types.Execution, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]*types.Execution, 0, len(l.store))
	for _, ex := range l.store {
		out = append(out, ex)
	}

	return out, nil
}

func isTerminal(s types.Status) bool {
	return s == types.StatusCompleted ||
		s == types.StatusFailed ||
		s == types.StatusCancelled
}
