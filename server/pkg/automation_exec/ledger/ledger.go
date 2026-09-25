package ledger

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
	"go.uber.org/zap"
)

type (
	ledger struct {
		mu sync.RWMutex
		// ExecutableID -> ExecutionID -> Revision -> Execution
		store map[id.ID]map[id.ID]map[int]*types.Execution
		log   *zap.Logger

		// sink receives a copy of every execution that reaches a terminal
		// status; the ledger itself is in-memory, so this is the one place a
		// run can be made durable
		sink ExecutionSink
	}

	// ExecutionSink is called with a snapshot of an execution once it has
	// completed, failed or been cancelled
	ExecutionSink func(ctx context.Context, ex types.Execution)

	Option func(*ledger)
)

func Ledger(log *zap.Logger, opts ...Option) *ledger {
	l := &ledger{
		store: make(map[id.ID]map[id.ID]map[int]*types.Execution),
		log:   log,
	}

	for _, o := range opts {
		o(l)
	}

	return l
}

// WithTerminalSink registers the function that receives every terminal execution
func WithTerminalSink(fn ExecutionSink) Option {
	return func(l *ledger) { l.sink = fn }
}

func (l *ledger) RegisterExecution(ctx context.Context, executableID, executionID id.ID, revision int, params types.ExecutionParams) error {
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
		EventType:    params.EventType,
		ResourceType: params.ResourceType,
		Events:       make([]types.StepEvent, 0),
	}

	return nil
}

func (l *ledger) ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, revision int) error {
	return l.transition(ctx, executableID, executionID, revision, types.StatusCompleted, nil)
}

func (l *ledger) ExecutionFailed(ctx context.Context, executableID, executionID id.ID, revision int, err error) error {
	return l.transition(ctx, executableID, executionID, revision, types.StatusFailed, err)
}

func (l *ledger) transition(ctx context.Context, executableID, executionID id.ID, revision int, status types.Status, err error) error {
	snapshot, terminal, trErr := l.apply(executableID, executionID, revision, status, err)
	if trErr != nil {
		return trErr
	}

	// the sink runs outside the lock so it may read the ledger freely
	if terminal && l.sink != nil {
		l.sink(ctx, snapshot)
	}

	return nil
}

func (l *ledger) apply(executableID, executionID id.ID, revision int, status types.Status, err error) (snapshot types.Execution, terminal bool, _ error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return snapshot, false, fmt.Errorf("execution not found")
	}

	// the runtime and its manager both report the end of a run; only the
	// first report ends it
	wasTerminal := isTerminal(ex.Status)

	now := time.Now()
	ex.Status = status
	ex.UpdatedAt = now
	ex.Error = err

	if isTerminal(status) {
		ex.EndedAt = &now
	}

	// trace and events are not copied: the sink records the outcome, not the steps
	snapshot = *ex
	snapshot.Trace = nil
	snapshot.Events = nil

	return snapshot, isTerminal(status) && !wasTerminal, nil
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

func (l *ledger) ExecutionPaused(ctx context.Context, executableID, executionID, stepID id.ID, phaseIndex, revision int, err error) error {
	if logErr := l.logStep(executableID, executionID, revision, types.StepEvent{
		Type:    types.EventStepPaused,
		StepID:  stepID,
		Payload: phaseIndex,
		Error:   err,
	}); logErr != nil {
		return logErr
	}

	return l.transition(ctx, executableID, executionID, revision, types.StatusPaused, nil)
}

func (l *ledger) RecordFrame(ctx context.Context, executableID, executionID id.ID, revision int, frame types.StackFrame) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	ex, ok := l.store[executableID][executionID][revision]
	if !ok {
		return fmt.Errorf("execution not found")
	}

	// Truncation: keep the last MaxIteratorFrames siblings per parent (sliding window).
	// When the cap is reached, drop the oldest sibling so the count stays stable.
	if !frame.ParentID.IsZero() {
		siblingIndices := make([]int, 0)
		for i, f := range ex.Trace {
			if f.ParentID == frame.ParentID {
				siblingIndices = append(siblingIndices, i)
			}
		}

		if len(siblingIndices) >= types.MaxIteratorFrames {
			removeIndex := siblingIndices[0] // drop the oldest
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

	// A copy: the stored trace stays live for the execution's lifetime, and a
	// caller that filters or sorts the returned slice in place would be
	// rewriting it.
	out := make([]types.StackFrame, len(ex.Trace))
	copy(out, ex.Trace)

	return out, nil
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

	sortExecutions(out)
	return out, nil
}

func (l *ledger) ListExecutionsByExecutable(ctx context.Context, executableID id.ID) ([]*types.Execution, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]*types.Execution, 0)
	execs, ok := l.store[executableID]
	if !ok {
		return out, nil
	}

	for _, revs := range execs {
		for _, ex := range revs {
			out = append(out, ex)
		}
	}

	sortExecutions(out)
	return out, nil
}

// sortExecutions orders executions newest first. The store is a map, so
// without this the caller gets Go's randomised iteration order and the same
// list comes back differently on consecutive calls — which makes out[0] the
// latest run only by chance, and sends anyone reading a trace to an arbitrary
// execution. ID descending breaks ties, since IDs are monotonic and two
// executions can share a timestamp.
func sortExecutions(ee []*types.Execution) {
	sort.Slice(ee, func(i, j int) bool {
		if !ee[i].CreatedAt.Equal(ee[j].CreatedAt) {
			return ee[i].CreatedAt.After(ee[j].CreatedAt)
		}
		return ee[i].ID.Num() > ee[j].ID.Num()
	})
}

func isTerminal(s types.Status) bool {
	return s == types.StatusCompleted ||
		s == types.StatusFailed ||
		s == types.StatusCancelled
}
