package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/automation/types"
	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	execTypes "github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/davecgh/go-spew/spew"
)

// ConvertNgAutomation converts service lvl structs into pkg/automation_exec
//
// @todo report issues so we don't brick the system in case of semantic errors
func ConvertNgAutomation(ctx context.Context, svc *ngAutomation, a *automationTypes.NgAutomation) (execTypes.Executable, automationTypes.NgAutomationIssueSet) {
	var issues automationTypes.NgAutomationIssueSet
	var err error

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

		aux := execTypes.Step{
			ID:     sid,
			Kind:   ui.Kind,
			Handle: ui.Handle,
		}

		for _, e := range ui.Arguments {
			aux.Arguments = append(aux.Arguments, execTypes.StepArg{
				Expr: &execTypes.Expr{
					Context:    e.Context,
					Target:     e.Target,
					Source:     e.Source,
					Expression: e.Expr,
					Value:      e.Value,
					Type:       e.Type,
				},
			})
		}

		// @todo fugly
		if ui.Kind == "function" {
			reg := Registry()
			def := reg.Function(ui.Ref)
			if def == nil {
				err = errors.Internal("unknown function %q", ui.Ref)
				panic(err)
			}

			ui.Results = []*automationTypes.Expr{}
			for _, r := range def.Results {
				ui.Results = append(ui.Results, &automationTypes.Expr{
					Target: r.Name,
					Type:   r.Types[0],
					Source: r.Name,
				})
			}
		}

		for _, e := range ui.Results {
			aux.Results = append(aux.Results, execTypes.StepRst{
				Name: e.Target,
			})
		}

		aux.Handler, err = stepConv(svc, ui)
		// @todo err handling...
		if err != nil {
			spew.Dump(err)
			err = nil
			continue
		}

		exSteps = append(exSteps, aux)
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

func stepConv(svc *ngAutomation, step *automationTypes.NgAutomationStep) (out execTypes.StepHandler, err error) {
	if err := parseExpressions(svc, step.Arguments...); err != nil {
		return nil, errors.Internal("failed to parse step arguments expressions for %s: %s", step.Kind, err).Wrap(err)
	}

	if err := parseExpressions(svc, step.Results...); err != nil {
		return nil, errors.Internal("failed to parse step results expressions for %s: %s", step.Kind, err).Wrap(err)
	}

	out, err = func() (out execTypes.StepHandler, err error) {
		switch step.Kind {
		case "function":
			return stepConvFunction(step)

		default:
			return nil, errors.Internal("unsupported step kind %q", step.Kind)
		}
	}()

	if err != nil {
		return nil, err
	}

	// @todo skip visual stuff
	if out == nil {
		panic("not supported")
	}

	return out, nil
}

func parseExpressions(svc *ngAutomation, ee ...*types.Expr) (err error) {
	for _, e := range ee {

		if len(strings.TrimSpace(e.Expr)) > 0 {
			if err = svc.parser.ParseEvaluators(e); err != nil {
				return
			}
		}

		if err = e.SetType(exprTypeSetter(Registry(), e)); err != nil {
			return err
		}

		for _, t := range e.Tests {
			if err = svc.parser.ParseEvaluators(t); err != nil {
				return
			}
		}
	}

	return nil
}

func stepConvFunction(step *automationTypes.NgAutomationStep) (out execTypes.StepHandler, err error) {
	reg := Registry()

	if def := reg.Function(step.Ref); def == nil {
		return nil, errors.Internal("unknown function %q", step.Ref)
	} else {
		if def.Kind != string(step.Kind) {
			return nil, fmt.Errorf("unexpected %s on %s step", def.Kind, step.Kind)
		}

		var (
			err        error
			isIterator = def.Kind == types.FunctionKindIterator
		)

		if isIterator {
			if def.Iterator == nil {
				return nil, errors.Internal("iterator handler for %q not set", step.Ref)
			}
		} else {
			if def.Handler == nil {
				return nil, errors.Internal("function handler for %q not set", step.Ref)
			}
		}

		if err = def.Parameters.VerifyArguments(step.Arguments); err != nil {
			return nil, errors.Internal("failed to verify argument expressions for %s %s: %s", step.Kind, step.Ref, err).Wrap(err)
		}

		if err = def.Results.VerifyResults(step.Results); err != nil {
			return nil, errors.Internal("failed to verify result expressions for %s %s: %s", step.Kind, step.Ref, err).Wrap(err)
		}

		// if isIterator {
		// 	if len(out) != 2 {
		// 		return nil, fmt.Errorf("expecting exactly 2 outbound paths for iterator")
		// 	}

		// 	var (
		// 		next = g.StepByID(out[0].ChildID)
		// 		exit = g.StepByID(out[1].ChildID)
		// 	)

		// 	if next == nil || exit == nil {
		// 		// wait for steps to be resolved
		// 		return nil, nil
		// 	}

		// 	return types.IteratorStep(def, step.Arguments, step.Results, next, exit)

		// } else {
		// }
		return types.FunctionStep(def, step.Arguments, step.Results)
	}
}
