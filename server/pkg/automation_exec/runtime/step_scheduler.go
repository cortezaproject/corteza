package runtime

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

var (
	ErrNoMoreSteps  = fmt.Errorf("scheduler: no runnable steps")
	ErrStepNotFound = fmt.Errorf("scheduler: step not found")
)

type scheduler struct {
	executable types.Executable

	// step index
	steps map[id.ID]*types.Step

	// remaining unmet dependencies per step
	remainingDeps map[id.ID]int

	// runnable queue
	nextSteps []id.ID
	queued    map[id.ID]bool
}

func newScheduler(exe types.Executable) *scheduler {
	ss := &scheduler{
		executable:    exe,
		steps:         make(map[id.ID]*types.Step, len(exe.Steps)),
		remainingDeps: make(map[id.ID]int, len(exe.Steps)),
		queued:        make(map[id.ID]bool, len(exe.Steps)),
	}

	return ss.init(exe)
}

func (ss *scheduler) Next() (*types.Step, bool, error) {
	if len(ss.nextSteps) == 0 {
		return nil, false, nil
	}

	// @todo this isn't the most efficient way to do it
	// we can use a ring buffer or something
	stepID := ss.nextSteps[0]
	ss.nextSteps = ss.nextSteps[1:]
	delete(ss.queued, stepID)

	step, ok := ss.steps[stepID]
	if !ok {
		return nil, false, ErrStepNotFound
	}

	return step, true, nil
}

func (ss *scheduler) OnStepComplete(stepID id.ID) {
	step, ok := ss.steps[stepID]
	if !ok {
		return
	}

	for _, child := range step.Children {
		ss.remainingDeps[child]--
		if ss.remainingDeps[child] == 0 {
			ss.enqueue(child)
		}
	}
}

func (ss *scheduler) enqueue(stepID id.ID) {
	if ss.queued[stepID] {
		return
	}
	ss.nextSteps = append(ss.nextSteps, stepID)
	ss.queued[stepID] = true
}

func (ss *scheduler) init(exe types.Executable) *scheduler {
	// index steps
	for i := range exe.Steps {
		step := &exe.Steps[i]
		ss.steps[step.ID] = step
	}

	// compute dependency counts
	for _, step := range exe.Steps {
		for _, child := range step.Children {
			ss.remainingDeps[child]++
		}
	}

	// enqueue entry steps (zero dependencies)
	for stepID := range ss.steps {
		if ss.remainingDeps[stepID] == 0 {
			ss.enqueue(stepID)
		}
	}

	return ss
}
