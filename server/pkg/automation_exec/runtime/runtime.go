package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

var (
	ErrExecutionStopped = errors.New("runtime: execution stopped")
	ErrStepFailed       = errors.New("runtime: step failed")
)

type ExecutionGate interface {
	Request(ops int) (<-chan struct{}, error)
}

type StateReporter interface {
	ReportStepStart(stepID id.ID) error
	ReportStepComplete(stepID id.ID, result StepResult) error
	ReportStepFailed(stepID id.ID, err error) error
	ReportExecutionComplete() error
	ReportExecutionFailed(err error) error
}

type StepResult struct {
	StepID      id.ID
	StartedAt   time.Time
	CompletedAt time.Time
	Output      map[string]any
	Error       *string
}

type ExecutionState struct {
	CurrentStep     *id.ID
	Variables       map[string]any
	CompletedSteps  map[id.ID]StepResult
	InProgressSteps map[id.ID]bool
	mu              sync.RWMutex
}

type runtime struct {
	exec      types.Executable
	scheduler *scheduler

	gate     ExecutionGate
	reporter StateReporter

	state *ExecutionState

	stopped atomic.Bool
	blocked atomic.Bool

	stopCh   chan struct{}
	resumeCh chan struct{}
}

// Runtime runs the executable
func Runtime(exec types.Executable, gate ExecutionGate, reporter StateReporter) *runtime {
	state := &ExecutionState{
		Variables:       make(map[string]any),
		CompletedSteps:  make(map[id.ID]StepResult),
		InProgressSteps: make(map[id.ID]bool),
	}

	return &runtime{
		exec:      exec,
		gate:      gate,
		reporter:  reporter,
		state:     state,
		scheduler: newScheduler(exec),
		stopCh:    make(chan struct{}),
		resumeCh:  make(chan struct{}, 1),
	}
}

func (r *runtime) Start(ctx context.Context, scope map[string]any) error {
	r.state.mu.Lock()
	r.state.Variables = scope
	r.state.mu.Unlock()

	for {
		select {
		case <-r.stopCh:
			return ErrExecutionStopped
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if r.blocked.Load() {
			select {
			case <-r.resumeCh:
			case <-r.stopCh:
				return ErrExecutionStopped
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		step, more, err := r.scheduler.Next()
		if err != nil {
			return err
		}
		if !more {
			return r.complete()
		}

		ops := r.getStepOps(step)
		if err := r.waitForPermission(ctx, ops); err != nil {
			return err
		}

		if err := r.executeStep(ctx, step); err != nil {
			return r.fail(err)
		}

		r.scheduler.OnStepComplete(step.ID)
	}
}

func (r *runtime) Stop() {
	if r.stopped.CompareAndSwap(false, true) {
		close(r.stopCh)
	}
}

func (r *runtime) Block() {
	r.blocked.Store(true)
}

func (r *runtime) Resume() {
	if r.blocked.CompareAndSwap(true, false) {
		select {
		case r.resumeCh <- struct{}{}:
		default:
		}
	}
}

func (r *runtime) IsBlocked() bool { return r.blocked.Load() }
func (r *runtime) IsStopped() bool { return r.stopped.Load() }

func (r *runtime) executeStep(ctx context.Context, step *types.Step) error {
	r.state.mu.Lock()
	r.state.CurrentStep = &step.ID
	r.state.InProgressSteps[step.ID] = true
	r.state.mu.Unlock()

	if err := r.reporter.ReportStepStart(step.ID); err != nil {
		return fmt.Errorf("report step start: %w", err)
	}

	r.state.mu.RLock()
	scope := r.copyVariablesLocked()
	r.state.mu.RUnlock()

	result := StepResult{
		StepID: step.ID,
		Output: make(map[string]any),
	}

	result.StartedAt = time.Now()
	output, err := step.Handler.Execute(ctx, scope)
	result.CompletedAt = time.Now()

	if err != nil {
		errMsg := err.Error()
		result.Error = &errMsg

		r.state.mu.Lock()
		delete(r.state.InProgressSteps, step.ID)
		r.state.CompletedSteps[step.ID] = result
		r.state.CurrentStep = nil
		r.state.mu.Unlock()

		_ = r.reporter.ReportStepFailed(step.ID, err)
		return fmt.Errorf("%w: %s", ErrStepFailed, step.ID)
	}

	result.Output = output

	r.state.mu.Lock()
	r.state.Variables = r.mergeScope(r.state.Variables, output)
	delete(r.state.InProgressSteps, step.ID)
	r.state.CompletedSteps[step.ID] = result
	r.state.CurrentStep = nil
	r.state.mu.Unlock()

	if err := r.reporter.ReportStepComplete(step.ID, result); err != nil {
		return fmt.Errorf("report step complete: %w", err)
	}

	return nil
}

func (r *runtime) mergeScope(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

func (r *runtime) waitForPermission(ctx context.Context, ops int) error {
	permCh, err := r.gate.Request(ops)
	if err != nil {
		return err
	}

	select {
	case <-permCh:
		return nil
	case <-r.stopCh:
		return ErrExecutionStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *runtime) complete() error {
	if err := r.reporter.ReportExecutionComplete(); err != nil {
		return fmt.Errorf("report execution complete: %w", err)
	}
	return nil
}

func (r *runtime) fail(err error) error {
	if repErr := r.reporter.ReportExecutionFailed(err); repErr != nil {
		return fmt.Errorf("execution failed: %w (report error: %v)", err, repErr)
	}
	return err
}

func (r *runtime) getStepOps(step *types.Step) int {
	if v, ok := step.Config["ops"].(int); ok && v > 0 {
		return v
	}
	return 1
}

func (r *runtime) copyVariablesLocked() map[string]any {
	vars := make(map[string]any, len(r.state.Variables))
	for k, v := range r.state.Variables {
		vars[k] = v
	}
	return vars
}

func (r *runtime) readOnlyStateLocked() map[string]any {
	return map[string]any{
		"current_step":    r.state.CurrentStep,
		"completed_steps": len(r.state.CompletedSteps),
	}
}
