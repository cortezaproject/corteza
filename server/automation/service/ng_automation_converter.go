package service

import (
	"context"
	"fmt"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	execTypes "github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/davecgh/go-spew/spew"
)

// ConvertNgAutomation converts service lvl structs into pkg/automation_exec
//
// @todo report issues so we don't brick the system in case of semantic errors
func ConvertNgAutomation(ctx context.Context, a *automationTypes.NgAutomation) (execTypes.Executable, error) {
	if a == nil {
		return execTypes.Executable{}, fmt.Errorf("nil automation")
	}

	// @todo when we do revisions
	revision := 0

	// Basic validation
	stepIdx := make(map[uint64]*automationTypes.NgAutomationStep, len(a.Steps))
	for i := range a.Steps {
		s := a.Steps[i]
		if s.ID == 0 {
			return execTypes.Executable{}, fmt.Errorf("step at index %d has empty ID", i)
		}
		if _, exists := stepIdx[s.ID]; exists {
			return execTypes.Executable{}, fmt.Errorf("duplicate step ID %d", s.ID)
		}
		stepIdx[s.ID] = s
	}

	// Build up steps
	exSteps := make([]execTypes.Step, 0, len(a.Steps))
	idMap := make(map[uint64]id.ID, len(a.Steps)) // uiStepID -> execStepID

	for _, ui := range stepIdx {
		sid := id.MustNumID(ui.ID)
		idMap[ui.ID] = sid

		exSteps = append(exSteps, execTypes.Step{
			ID:      sid,
			Kind:    ui.Kind,
			Handler: noopHandler{},
			// Parents/Children wired in Phase 3
		})
	}

	// Index execution steps by id.ID for wiring
	exByID := make(map[id.ID]*execTypes.Step, len(exSteps))
	for i := range exSteps {
		exByID[exSteps[i].ID] = &exSteps[i]
	}

	// Connect bits
	for i := range a.Paths {
		p := a.Paths[i]

		if p.ParentID == 0 || p.ChildID == 0 {
			return execTypes.Executable{}, fmt.Errorf("path at index %d has empty parent/child", i)
		}
		if p.ParentID == p.ChildID {
			return execTypes.Executable{}, fmt.Errorf("path at index %d is a self-loop (%d)", i, p.ParentID)
		}

		parentExecID, ok := idMap[p.ParentID]
		if !ok {
			return execTypes.Executable{}, fmt.Errorf("path at index %d references unknown parent step %d", i, p.ParentID)
		}
		childExecID, ok := idMap[p.ChildID]
		if !ok {
			return execTypes.Executable{}, fmt.Errorf("path at index %d references unknown child step %d", i, p.ChildID)
		}

		parent := exByID[parentExecID]
		child := exByID[childExecID]
		if parent == nil || child == nil {
			return execTypes.Executable{}, fmt.Errorf("internal step index missing for path at index %d", i)
		}

		// NOTE: types.Step uses value slices for Parents/Children, so this copies the structs.
		// This is fine for now (scheduler only needs IDs), but it is structurally lossy.
		parent.Children = append(parent.Children, *child)
		child.Parents = append(child.Parents, *parent)
	}

	// Graph validation
	if !hasEntry(exByID) {
		return execTypes.Executable{}, fmt.Errorf("no entry steps (every step has at least one parent)")
	}
	if err := detectCycle(exByID); err != nil {
		return execTypes.Executable{}, err
	}

	return execTypes.Executable{
		ID:       id.MustNumID(a.ID),
		Revision: revision,
		Handle:   a.Handle,
		Steps:    exSteps,
	}, nil
}

// noopHandler is the temporary handler used for all steps.
// It returns the scope unchanged and never errors.
type noopHandler struct{}

func (noopHandler) Exec(ctx context.Context, r *execTypes.ExecRequest) (execTypes.ExecResponse, error) {
	spew.Dump("noop exec")

	if r == nil {
		return nil, nil
	}
	return r.Scope, nil
}

func hasEntry(exByID map[id.ID]*execTypes.Step) bool {
	if len(exByID) == 0 {
		return true
	}

	for _, s := range exByID {
		if s == nil {
			continue
		}
		if len(s.Parents) == 0 {
			return true
		}
	}
	return false
}

// detectCycle runs a DFS cycle check using only IDs.
// Uses step.Children slice (value-copies) but depends only on child IDs.
func detectCycle(exByID map[id.ID]*execTypes.Step) error {
	const (
		white = 0
		gray  = 1
		black = 2
	)

	color := make(map[id.ID]uint8, len(exByID))

	var visit func(n id.ID) error
	visit = func(n id.ID) error {
		switch color[n] {
		case gray:
			return fmt.Errorf("cycle detected at step %v", n)
		case black:
			return nil
		}

		color[n] = gray
		s := exByID[n]
		if s != nil {
			for _, c := range s.Children {
				if err := visit(c.ID); err != nil {
					return err
				}
			}
		}
		color[n] = black
		return nil
	}

	for sid := range exByID {
		if color[sid] == white {
			if err := visit(sid); err != nil {
				return err
			}
		}
	}

	return nil
}
