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
	frameTypeGateway
)

type (
	scheduler struct {
		executable types.Executable
		steps      map[id.ID]*types.Step

		// Execution stack
		stack []*frame

		// Completed step outputs indexed by handle
		completedOutputs map[string]map[string]expr.TypedValue

		// Branch tracking for termination
		activeBranches map[id.ID]bool
	}

	IteratorHandler interface {
		Start(context.Context, *expr.Vars) error
		More(context.Context, *expr.Vars) (bool, error)
		Next(context.Context, *expr.Vars) (*expr.Vars, error)
	}

	GatewayHandler interface {
		Select(context.Context, map[string]*expr.Vars) ([]int, error)
	}

	frame struct {
		typ frameType

		id       id.ID
		parentID id.ID

		stepID id.ID
		step   *types.Step

		branchID id.ID

		// Output tracking
		handle  string
		outputs map[string]expr.TypedValue

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
		completedOutputs: make(map[string]map[string]expr.TypedValue),
		activeBranches:   make(map[id.ID]bool),
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

		case frameTypeGateway:
			return ss.handleGateway(ctx, current)

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
func (ss *scheduler) StoreOutputs(stepID id.ID, outputs map[string]expr.TypedValue) error {
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
		var bid id.ID
		if top != nil {
			parentID = top.id
			bid = top.branchID
		}

		ss.stack = append(ss.stack, &frame{
			typ:         frameTypeIterator,
			id:          id.MustNumID(id.Next()),
			parentID:    parentID,
			stepID:      stepID,
			step:        step,
			handle:      step.Handle,
			iterHandler: handler,
			branchID:    bid,
		})

	case "gatewayExclusive", "gatewayInclusive":
		var parentID id.ID
		var bid id.ID
		if currentFrame != nil {
			bid = currentFrame.branchID
		} else if top != nil {
			bid = top.branchID
		}
		if top != nil {
			parentID = top.id
		}

		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeGateway,
			id:       id.MustNumID(id.Next()),
			parentID: parentID,
			stepID:   stepID,
			step:     step,
			handle:   step.Handle,
			branchID: bid,
		})

	default:
		// Expression/Function: push children
		// Children of this step should have this frame as parent
		var bid id.ID
		if currentFrame != nil {
			bid = currentFrame.branchID
		} else if top != nil {
			bid = top.branchID
		}

		ss.pushChildren(step, currentFrameID, bid)
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
				branchID: f.branchID,
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
	iterOutputs := make(map[string]expr.TypedValue)
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
		branchID: f.branchID,
	})

	return ss.Next(ctx)
}

func (ss *scheduler) handleGateway(ctx context.Context, f *frame) (*types.Step, id.ID, id.ID, bool, error) {
	if f.branchTaken {
		ss.stack = ss.stack[:len(ss.stack)-1]
		return ss.Next(ctx)
	}

	handler, ok := f.step.Handler.(GatewayHandler)
	if !ok {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("step %v: handler is not GatewayHandler", f.stepID)
	}

	// Build scope map from frame outputs. Gateway conditions reference variables
	// via Meta["scope"]; without explicit scope they use "global".
	globalVars := &expr.Vars{}
	if f.outputs != nil {
		for k, v := range f.outputs {
			globalVars.Set(k, v)
		}
	}
	condScope := map[string]*expr.Vars{"global": globalVars}

	indices, err := handler.Select(ctx, condScope)
	if err != nil {
		return nil, id.Zero(), id.Zero(), false, fmt.Errorf("gateway selection failed: %w", err)
	}

	f.branchTaken = true

	// Push selected children for execution
	// Push in reverse as this is a stack
	for i := len(indices) - 1; i >= 0; i-- {
		idx := indices[i]
		if idx < 0 || idx >= len(f.step.Children) {
			continue
		}
		nextStep := &f.step.Children[idx]
		bid := f.branchID
		if len(indices) > 1 {
			bid = id.MustNumID(id.Next())
			ss.activeBranches[bid] = true
		}
		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeStep,
			id:       id.MustNumID(id.Next()),
			parentID: f.id,
			stepID:   nextStep.ID,
			step:     nextStep,
			handle:   nextStep.Handle,
			branchID: bid,
		})
	}

	return ss.Next(ctx)
}

func (ss *scheduler) pushChildren(step *types.Step, parentID, branchID id.ID) {
	// Push children in reverse order so first child executes first

	if len(step.Children) == 1 {
		child := &step.Children[0]
		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeStep,
			id:       id.MustNumID(id.Next()),
			parentID: parentID,
			stepID:   child.ID,
			step:     child,
			handle:   child.Handle,
			branchID: branchID,
		})
		return
	}

	// Multiple children spawn new branches
	delete(ss.activeBranches, branchID)

	for i := len(step.Children) - 1; i >= 0; i-- {
		child := &step.Children[i]
		bid := id.MustNumID(id.Next())
		ss.activeBranches[bid] = true

		ss.stack = append(ss.stack, &frame{
			typ:      frameTypeStep,
			id:       id.MustNumID(id.Next()),
			parentID: parentID,
			stepID:   child.ID,
			step:     child,
			handle:   child.Handle,
			branchID: bid,
		})
	}
}

// TerminateBranch marks the current branch as terminated
// Returns true if all branches are terminated (automation should end)
func (ss *scheduler) TerminateBranch(stepID id.ID) (bool, error) {
	// Find the current branch ID from the frame associated with this stepID
	var branchID id.ID
	for i := len(ss.stack) - 1; i >= 0; i-- {
		if ss.stack[i].stepID.Equal(stepID) {
			branchID = ss.stack[i].branchID
			break
		}
	}

	if branchID.IsZero() {
		// Should not happen if step is on stack
		return len(ss.activeBranches) == 0, nil
	}

	delete(ss.activeBranches, branchID)
	ss.clearBranchFrames(branchID)

	return len(ss.activeBranches) == 0, nil
}

func (ss *scheduler) clearBranchFrames(branchID id.ID) {
	newStack := make([]*frame, 0, len(ss.stack))
	for _, f := range ss.stack {
		if !f.branchID.Equal(branchID) {
			newStack = append(newStack, f)
		}
	}
	ss.stack = newStack
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
			bid := id.MustNumID(id.Next())
			ss.activeBranches[bid] = true

			ss.stack = append(ss.stack, &frame{
				typ:      frameTypeStep,
				id:       id.MustNumID(id.Next()),
				stepID:   step.ID,
				step:     step,
				handle:   step.Handle,
				branchID: bid,
			})
		}
	}

	return ss
}
