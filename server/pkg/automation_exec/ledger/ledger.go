package ledger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"go.uber.org/zap"
)

type ledger struct {
	mu sync.RWMutex
	// ExecutableID -> ExecutionID -> Revision -> Execution
	store map[id.ID]map[id.ID]map[int]*types.Execution
	log   *zap.Logger
}

func Ledger(log *zap.Logger) *ledger {
	return &ledger{
		store: make(map[id.ID]map[id.ID]map[int]*types.Execution),
		log:   log,
	}
}

func (l *ledger) RegisterExecution(ctx context.Context, executableID, executionID id.ID, revision int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.store[executableID]; !ok {
		l.store[executableID] = make(map[id.ID]map[int]*types.Execution)
	}
	if _, ok := l.store[executableID][executionID]; !ok {
		l.store[executableID][executionID] = make(map[int]*types.Execution)
	}

	now := time.Now()
	l.store[executableID][executionID][revision] = &types.Execution{
		ID:           executionID,
		ExecutableID: executableID,
		Revision:     revision,
		Status:       types.StatusCreated,
		CreatedAt:    now,
		UpdatedAt:    now,
		Events:       make([]types.StepEvent, 0),
	}

	return nil
}

func (l *ledger) ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, revision int) error {
	return l.transition(executableID, executionID, revision, types.StatusCompleted, nil)
}

func (l *ledger) ExecutionFailed(ctx context.Context, executableID, executionID id.ID, revision int, err error) error {
	return l.transition(executableID, executionID, revision, types.StatusFailed, err)
}

func (l *ledger) transition(executableID, executionID id.ID, revision int, status types.Status, err error) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return fmt.Errorf("execution not found")
	}

	now := time.Now()
	ex.Status = status
	ex.UpdatedAt = now
	ex.Error = err

	if isTerminal(status) {
		ex.EndedAt = &now
	}

	return nil
}

func (l *ledger) StepStarted(ctx context.Context, executableID, executionID, stepID id.ID, revision int) error {
	return l.logStep(executableID, executionID, revision, types.StepEvent{
		Type:   types.EventStepStarted,
		StepID: stepID,
	})
}

func (l *ledger) StepCompleted(ctx context.Context, executableID, executionID, stepID id.ID, revision int, payload any) error {
	return l.logStep(executableID, executionID, revision, types.StepEvent{
		Type:    types.EventStepCompleted,
		StepID:  stepID,
		Payload: payload,
	})
}

func (l *ledger) StepFailed(ctx context.Context, executableID, executionID, stepID id.ID, revision int, err error) error {
	return l.logStep(executableID, executionID, revision, types.StepEvent{
		Type:   types.EventStepFailed,
		StepID: stepID,
		Error:  err,
	})
}

func (l *ledger) RecordFrame(ctx context.Context, executableID, executionID id.ID, revision int, frame types.StackFrame) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return fmt.Errorf("execution not found")
	}

	// Truncation logic: Per-iterator truncation
	if !frame.ParentID.IsZero() {
		siblingIndices := make([]int, 0)
		for i, f := range ex.Trace {
			if f.ParentID == frame.ParentID {
				siblingIndices = append(siblingIndices, i)
			}
		}

		if len(siblingIndices) >= types.MaxIteratorFrames {
			// Keep first 100 iterations, truncate from the 101st
			// removeIndex is the index in the global Trace slice
			removeIndex := siblingIndices[100]
			ex.Trace = append(ex.Trace[:removeIndex], ex.Trace[removeIndex+1:]...)
		}
	}

	ex.Trace = append(ex.Trace, frame)
	ex.UpdatedAt = time.Now()

	return nil
}

func (l *ledger) GetTrace(ctx context.Context, executableID, executionID id.ID, revision int) ([]types.StackFrame, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return nil, fmt.Errorf("execution not found")
	}

	return ex.Trace, nil
}

func (l *ledger) logStep(executableID, executionID id.ID, revision int, event types.StepEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return fmt.Errorf("execution not found")
	}

	event.Timestamp = time.Now()
	ex.Events = append(ex.Events, event)
	ex.UpdatedAt = event.Timestamp

	return nil
}

func (l *ledger) IsExecutableInUse(ctx context.Context, executableID id.ID, revision int) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	execs, ok := l.store[executableID]
	if !ok {
		return false, nil
	}

	for _, revs := range execs {
		ex, ok := revs[revision]
		if !ok {
			continue
		}

		if !isTerminal(ex.Status) {
			return true, nil
		}
	}

	return false, nil
}

func (l *ledger) GetExecution(ctx context.Context, executableID, executionID id.ID, revision int) (*types.Execution, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return nil, fmt.Errorf("execution not found")
	}

	return ex, nil
}

func (l *ledger) ListExecutions(ctx context.Context) ([]*types.Execution, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]*types.Execution, 0)
	for _, execs := range l.store {
		for _, revs := range execs {
			for _, ex := range revs {
				out = append(out, ex)
			}
		}
	}

	return out, nil
}

func isTerminal(s types.Status) bool {
	return s == types.StatusCompleted ||
		s == types.StatusFailed ||
		s == types.StatusCancelled
}
