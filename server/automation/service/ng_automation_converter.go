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
func ConvertNgAutomation(ctx context.Context, a *automationTypes.NgAutomation) (execTypes.Executable, automationTypes.NgAutomationIssueSet) {
	var issues automationTypes.NgAutomationIssueSet

	if a == nil {
		return execTypes.Executable{}, automationTypes.NgAutomationIssueSet{
			issue("nil automation", nil),
		}
	}

	revision := 0

	// --- Phase 1: validate & index steps ---
	stepIdx := make(map[uint64]*automationTypes.NgAutomationStep, len(a.Steps))
	for i := range a.Steps {
		s := a.Steps[i]

		if s.ID == 0 {
			issues = append(issues, issue(
				"step has empty ID",
				map[string]int{"step": i},
			))
			continue
		}

		if _, exists := stepIdx[s.ID]; exists {
			issues = append(issues, issue(
				fmt.Sprintf("duplicate step ID %d", s.ID),
				map[string]int{"step": i},
			))
			continue
		}

		stepIdx[s.ID] = s
	}

	if len(stepIdx) == 0 {
		return execTypes.Executable{}, append(issues,
			issue("no valid steps defined", nil),
		)
	}

	// --- Phase 2: build execution steps ---
	exSteps := make([]execTypes.Step, 0, len(stepIdx))
	idMap := make(map[uint64]id.ID, len(stepIdx))

	for uiID, ui := range stepIdx {
		sid := id.MustNumID(uiID)
		idMap[uiID] = sid

		exSteps = append(exSteps, execTypes.Step{
			ID:      sid,
			Kind:    ui.Kind,
			Handler: noopHandler{},
		})
	}

	exByID := make(map[id.ID]*execTypes.Step, len(exSteps))
	for i := range exSteps {
		exByID[exSteps[i].ID] = &exSteps[i]
	}

	// --- Phase 3: wire paths ---
	for i := range a.Paths {
		p := a.Paths[i]

		if p.ParentID == 0 || p.ChildID == 0 {
			issues = append(issues, issue(
				"path has empty parent or child",
				map[string]int{"path": i},
			))
			continue
		}

		if p.ParentID == p.ChildID {
			issues = append(issues, issue(
				"path is a self-loop",
				map[string]int{"path": i},
			))
			continue
		}

		parentExecID, ok := idMap[p.ParentID]
		if !ok {
			issues = append(issues, issue(
				fmt.Sprintf("unknown parent step %d", p.ParentID),
				map[string]int{"path": i},
			))
			continue
		}

		childExecID, ok := idMap[p.ChildID]
		if !ok {
			issues = append(issues, issue(
				fmt.Sprintf("unknown child step %d", p.ChildID),
				map[string]int{"path": i},
			))
			continue
		}

		parent := exByID[parentExecID]
		child := exByID[childExecID]
		if parent == nil || child == nil {
			issues = append(issues, issue(
				"internal step index missing",
				map[string]int{"path": i},
			))
			continue
		}

		parent.Children = append(parent.Children, *child)
		child.Parents = append(child.Parents, *parent)
	}

	// --- Phase 4: graph validation ---
	if !hasEntry(exByID) {
		issues = append(issues, issue(
			"no entry steps (every step has at least one parent)",
			nil,
		))
	}

	if err := detectCycle(exByID); err != nil {
		issues = append(issues, issue(
			err.Error(),
			nil,
		))
	}

	// Return executable even if issues exist
	return execTypes.Executable{
		ID:       id.MustNumID(a.ID),
		Revision: revision,
		Handle:   a.Handle,
		Steps:    exSteps,
	}, issues
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

func issue(desc string, culprit map[string]int) *automationTypes.NgAutomationIssue {
	return &automationTypes.NgAutomationIssue{
		Description: desc,
		Culprit:     culprit,
	}
}
