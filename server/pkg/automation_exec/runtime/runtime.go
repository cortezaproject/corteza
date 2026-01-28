package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

var (
	ErrExecutionStopped = errors.New("runtime: execution stopped")
	ErrStepFailed       = errors.New("runtime: step failed")
)

type executionGate interface {
	Request(execID id.ID, ops int) (<-chan struct{}, error)
}

type stateReporter interface {
	StepStarted(ctx context.Context, executableID, executionID, stepID id.ID, rev int) error
	StepCompleted(ctx context.Context, executableID, executionID, stepID id.ID, rev int, out any) error
	StepFailed(ctx context.Context, executableID, executionID, stepID id.ID, rev int, err error) error

	ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, rev int) error
	ExecutionFailed(ctx context.Context, executableID, executionID id.ID, rev int, err error) error
}

type StepResult struct {
	StepID      id.ID
	StartedAt   time.Time
	CompletedAt time.Time
	Output      types.ExecResponse
	Error       error
}

type executionState struct {
	CurrentStep     *id.ID
	Variables       *expr.Vars
	CompletedSteps  map[id.ID]StepResult
	InProgressSteps map[id.ID]bool
	mux             sync.RWMutex
}

type runtime struct {
	execID    id.ID
	exec      types.Executable
	scheduler *scheduler

	gate     executionGate
	reporter stateReporter

	state *executionState

	stopped atomic.Bool
	blocked atomic.Bool

	stopCh   chan struct{}
	resumeCh chan struct{}
}

func Runtime(
	executionID id.ID,
	exec types.Executable,
	gate executionGate,
	reporter stateReporter,
) *runtime {
	state := &executionState{
		CompletedSteps:  make(map[id.ID]StepResult),
		InProgressSteps: make(map[id.ID]bool),
	}

	return &runtime{
		execID:    executionID,
		exec:      exec,
		gate:      gate,
		reporter:  reporter,
		state:     state,
		scheduler: newScheduler(exec),
		stopCh:    make(chan struct{}),
		resumeCh:  make(chan struct{}, 1),
	}
}

func (r *runtime) Start(ctx context.Context, scope *expr.Vars) error {
	r.state.mux.Lock()
	r.state.Variables = scope
	r.state.mux.Unlock()

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
			return r.complete(ctx)
		}

		if err := r.waitForPermission(ctx, r.getStepOps(step)); err != nil {
			return err
		}

		if err := r.executeStep(ctx, step); err != nil {
			return r.fail(ctx, err)
		}

		r.scheduler.OnStepComplete(step.ID)
	}
}

func (r *runtime) Stop() {
	if r.stopped.CompareAndSwap(false, true) {
		close(r.stopCh)
	}
}

func (r *runtime) Block() { r.blocked.Store(true) }
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
	r.state.mux.Lock()
	r.state.CurrentStep = &step.ID
	r.state.InProgressSteps[step.ID] = true
	r.state.mux.Unlock()

	if err := r.reporter.StepStarted(ctx, r.exec.ID, r.execID, step.ID, r.exec.Revision); err != nil {
		return fmt.Errorf("step started: %w", err)
	}

	r.state.mux.RLock()
	scope := (&expr.Vars{}).MustMerge(r.state.Variables)
	r.state.mux.RUnlock()

	result := StepResult{
		StepID:    step.ID,
		StartedAt: time.Now(),
	}

	output, err := step.Handler.Exec(ctx, &types.ExecRequest{
		Scope: scope,
	})
	result.CompletedAt = time.Now()

	if err != nil {
		result.Error = err

		r.state.mux.Lock()
		delete(r.state.InProgressSteps, step.ID)
		r.state.CompletedSteps[step.ID] = result
		r.state.CurrentStep = nil
		r.state.mux.Unlock()

		_ = r.reporter.StepFailed(ctx, r.exec.ID, r.execID, step.ID, r.exec.Revision, err)
		return fmt.Errorf("%w: %s", ErrStepFailed, step.ID)
	}

	var newScope *expr.Vars
	if rs, ok := output.(*expr.Vars); ok {
		newScope = rs
	}

	result.Output = output

	r.state.mux.Lock()
	if newScope != nil {
		r.state.Variables = r.state.Variables.MustMerge(newScope)
	}
	delete(r.state.InProgressSteps, step.ID)
	r.state.CompletedSteps[step.ID] = result
	r.state.CurrentStep = nil
	r.state.mux.Unlock()

	if err := r.reporter.StepCompleted(ctx, r.exec.ID, r.execID, step.ID, r.exec.Revision, output); err != nil {
		return fmt.Errorf("step completed: %w", err)
	}

	return nil
}

func (r *runtime) waitForPermission(ctx context.Context, ops int) error {
	permCh, err := r.gate.Request(r.execID, ops)
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

func (r *runtime) complete(ctx context.Context) error {
	if err := r.reporter.ExecutionCompleted(ctx, r.exec.ID, r.execID, r.exec.Revision); err != nil {
		return fmt.Errorf("execution complete: %w", err)
	}
	return nil
}

func (r *runtime) fail(ctx context.Context, err error) error {
	if repErr := r.reporter.ExecutionFailed(ctx, r.exec.ID, r.execID, r.exec.Revision, err); repErr != nil {
		return fmt.Errorf("execution failed: %w (report error: %v)", err, repErr)
	}
	return err
}

func (r *runtime) getStepOps(step *types.Step) int {
	return 1
}
