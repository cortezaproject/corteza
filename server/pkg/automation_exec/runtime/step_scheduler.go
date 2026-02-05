package runtime

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

var (
	ErrNoMoreSteps    = fmt.Errorf("scheduler: no runnable steps")
	ErrStepNotFound   = fmt.Errorf("scheduler: step not found")
	ErrOutputNotFound = fmt.Errorf("scheduler: output not found")
)

const (
	frameTypeStep frameType = iota
	frameTypeIterator
	frameTypeBranch
)

type (
	scheduler struct {
		executable types.Executable
		steps      map[id.ID]*types.Step

		// Execution stack
		stack []*frame

		// Completed step outputs indexed by handle
		completedOutputs map[string]map[string]*expr.Vars
	}

	IteratorHandler interface {
		Start(context.Context, *expr.Vars) error
		More(context.Context, *expr.Vars) (bool, error)
		Next(context.Context, *expr.Vars) (*expr.Vars, error)
	}

	BranchHandler interface {
		Evaluate(context.Context, *expr.Vars) (bool, error)
	}

	frame struct {
		typ frameType

		id       id.ID
		parentID id.ID

		stepID id.ID
		step   *types.Step

		// Output tracking
		handle  string
		outputs map[string]*expr.Vars

		// For iterators
		iterHandler IteratorHandler
		iterStarted bool

		// For branches
		branchTaken bool
	}

	frameType uint8
)

func newScheduler(exe types.Executable) *scheduler {
	ss := &scheduler{
		executable:       exe,
		steps:            make(map[id.ID]*types.Step, len(exe.Steps)),
		stack:            make([]*frame, 0, 16),
		completedOutputs: make(map[string]map[string]*expr.Vars),
	}

	return ss.init(exe)
}

// Next returns the next step to execute
// Returns (step, hasNext, error)
// Note: For frameTypeStep, the frame is NOT popped here - it remains on stack
// until PopFrame() is called after StoreOutputs
func (ss *scheduler) Next(ctx context.Context) (*types.Step, id.ID, id.ID, bool, error) {
	for len(ss.stack) > 0 {
		current := ss.stack[len(ss.stack)-1]

		if current.id.IsZero() {
			current.id = id.MustNumID(id.Next())
		}

		switch current.typ {
		case frameTypeIterator:
			return ss.handleIterator(ctx, current)

		case frameTypeBranch:
			return ss.handleBranch(ctx, current)

		case frameTypeStep:
			// Don't pop yet - wait for StoreOutputs to be called first
			return current.step, current.id, current.parentID, true, nil
		}
	}

	return nil, id.Zero(), id.Zero(), false, nil
}

func (ss *scheduler) FindOutput(handle string) (*expr.Vars, error) {
	// First check completed outputs
	if outputs, ok := ss.completedOutputs[handle]; ok {
		return expr.NewVars(outputs)
	}

	// Then check stack for in-progress frames
	for i := len(ss.stack) - 1; i >= 0; i-- {
		f := ss.stack[i]

		if f.handle == handle && f.outputs != nil {
			return expr.NewVars(f.outputs)
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrOutputNotFound, handle)
}

// StoreOutputs saves step execution results and moves them to completedOutputs
func (ss *scheduler) StoreOutputs(stepID id.ID, outputs map[string]*expr.Vars) error {
	// Find the frame for this step
	for i := len(ss.stack) - 1; i >= 0; i-- {
		f := ss.stack[i]
		if f.stepID.Equal(stepID) {
			handle := f.handle
			// Store outputs in completed map for future lookups
			if handle != "" {
				ss.completedOutputs[handle] = outputs
			}
			f.outputs = outputs
			return nil
		}
	}

	return fmt.Errorf("step %v not found on stack", stepID)
}

// OnStepComplete notifies scheduler that a step finished
// Pushes appropriate child frames based on step kind
func (ss *scheduler) OnStepComplete(ctx context.Context, stepID id.ID) error {
	step, ok := ss.steps[stepID]
	if !ok {
		return ErrStepNotFound
	}

	// We finished a step, so the top of the stack should be this step's frame.
	// We pop it now.
	var currentFrame *frame
	if len(ss.stack) > 0 {
		top := ss.stack[len(ss.stack)-1]
		if top.stepID.Equal(stepID) && top.typ == frameTypeStep {
			currentFrame = top
			ss.stack = ss.stack[:len(ss.stack)-1]
		}
	}

	// If the stack is now empty, we've finished a root sequence.
	// That's fine, no error needed.
	var top *frame
	if len(ss.stack) > 0 {
		top = ss.stack[len(ss.stack)-1]
	}

	// Check if the current frame was inside an iterator loop
	// (Note: frameTypeIterator doesn't represent the step itself, but the loop management)
	if top != nil && top.typ == frameTypeIterator && top.stepID == stepID {
		// Iterator frame handles its own continuation
		return nil
	}

	// If we don't have a current frame (already popped or somehow missing),
	// we can't reliably push children with correct parent ID.
	// But we should at least try to push them if they exist.
	var currentFrameID id.ID
	if currentFrame != nil {
		currentFrameID = currentFrame.id
	}

	switch step.Kind {
	case "iterator":
		handler, ok := step.Handler.(IteratorHandler)
		if !ok {
			return fmt.Errorf("step %v: handler is not IteratorHandler", stepID)
		}

		var parentID id.ID
		if top != nil {
			parentID = top.id
		}

		ss.stack = append(ss.stack, &frame{
			typ:         frameTypeIterator,
			id:          id.MustNumID(id.Next()),
			parentID:    parentID,
			stepID:      stepID,
			step:        step,
			handle:      step.Handle,
			iterHandler: handler,
		})

	case "branch":
		var parentID id.ID
		if top != nil {
			parentID = top.id
		}

		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeBranch,
			id:       id.MustNumID(id.Next()),
			parentID: parentID,
			stepID:   stepID,
			step:     step,
			handle:   step.Handle,
		})

	default:
		// Expression/Function: push children
		// Children of this step should have this frame as parent
		ss.pushChildren(step, currentFrameID)
	}

	return nil
}

func (ss *scheduler) handleIterator(ctx context.Context, f *frame) (*types.Step, id.ID, id.ID, bool, error) {
	if !f.iterStarted {
		// Get initial vars from frame outputs (set by runtime after step exec)
		initialVars := &expr.Vars{}
		if f.outputs != nil {
			for k, v := range f.outputs {
				initialVars.Set(k, v)
			}
		}

		if err := f.iterHandler.Start(ctx, initialVars); err != nil {
			return nil, id.Zero(), id.Zero(), false, fmt.Errorf("iterator start failed: %w", err)
		}
		f.iterStarted = true
	}

	currentVars := &expr.Vars{}
	if f.outputs != nil {
		for k, v := range f.outputs {
			currentVars.Set(k, v)
		}
	}

	more, err := f.iterHandler.More(ctx, currentVars)
	if err != nil {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("iterator more check failed: %w", err)
	}

	if !more {
		// Iterator done, pop frame and push exit edge
		ss.stack = ss.stack[:len(ss.stack)-1]

		// Second child is the exit path (if exists)
		if len(f.step.Children) > 1 {
			exitStep := &f.step.Children[1]
			ss.stack = append(ss.stack, &frame{
				typ:      frameTypeStep,
				id:       id.MustNumID(id.Next()),
				parentID: f.parentID, // sibling of iterator
				stepID:   exitStep.ID,
				step:     exitStep,
				handle:   exitStep.Handle,
			})
		}

		return ss.Next(ctx)
	}

	// Get next iteration vars
	iterVars, err := f.iterHandler.Next(ctx, currentVars)
	if err != nil {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("iterator next failed: %w", err)
	}

	// Convert expr.Vars to map for frame outputs
	iterOutputs := make(map[string]*expr.Vars)
	if iterVars != nil {
		// TODO: implement vars to map conversion based on expr.Vars API
	}

	// First child is the body path
	if len(f.step.Children) == 0 {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("iterator missing body edge")
	}

	bodyStep := &f.step.Children[0]

	ss.stack = append(ss.stack, &frame{
		typ:      frameTypeStep,
		id:       id.MustNumID(id.Next()),
		parentID: f.id,
		stepID:   bodyStep.ID,
		step:     bodyStep,
		handle:   bodyStep.Handle,
		outputs:  iterOutputs,
	})

	return ss.Next(ctx)
}

func (ss *scheduler) handleBranch(ctx context.Context, f *frame) (*types.Step, id.ID, id.ID, bool, error) {
	if f.branchTaken {
		ss.stack = ss.stack[:len(ss.stack)-1]
		return ss.Next(ctx)
	}

	handler, ok := f.step.Handler.(BranchHandler)
	if !ok {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("step %v: handler is not BranchHandler", f.stepID)
	}

	// Build vars from frame outputs
	condVars := &expr.Vars{}
	if f.outputs != nil {
		for k, v := range f.outputs {
			condVars.Set(k, v)
		}
	}

	condition, err := handler.Evaluate(ctx, condVars)
	if err != nil {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("branch evaluation failed: %w", err)
	}

	f.branchTaken = true

	// First child is true path, second child is false path
	childIndex := 1 // false
	if condition {
		childIndex = 0 // true
	}

	if len(f.step.Children) > childIndex {
		nextStep := &f.step.Children[childIndex]
		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeStep,
			id:       id.MustNumID(id.Next()),
			parentID: f.id,
			stepID:   nextStep.ID,
			step:     nextStep,
			handle:   nextStep.Handle,
		})
	}

	return ss.Next(ctx)
}

func (ss *scheduler) pushChildren(step *types.Step, parentID id.ID) {
	// Push children in reverse order so first child executes first
	for i := len(step.Children) - 1; i >= 0; i-- {
		child := &step.Children[i]
		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeStep,
			id:       id.MustNumID(id.Next()),
			parentID: parentID,
			stepID:   child.ID,
			step:     child,
			handle:   child.Handle,
		})
	}
}

func (ss *scheduler) init(exe types.Executable) *scheduler {
	// Index steps
	for i := range exe.Steps {
		step := &exe.Steps[i]
		ss.steps[step.ID] = step
	}

	// Find entry points (steps with no parents)
	for i := range exe.Steps {
		step := &exe.Steps[i]
		if len(step.Parents) == 0 {
			ss.stack = append(ss.stack, &frame{
				typ:    frameTypeStep,
				id:     id.MustNumID(id.Next()),
				stepID: step.ID,
				step:   step,
				handle: step.Handle,
			})
		}
	}

	return ss
}
