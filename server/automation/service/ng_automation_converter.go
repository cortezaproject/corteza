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

	// --- Phase 1a: validate & index triggers ---
	triggerIdx := make(map[uint64]*automationTypes.NgAutomationTrigger, len(a.Triggers))
	triggerPathCount := make(map[uint64]int, len(a.Triggers))
	for i := range a.Triggers {
		t := a.Triggers[i]

		if t.ID == 0 {
			issues = append(issues, issue(
				"trigger has empty ID",
				map[string]int{"trigger": i},
			))
			continue
		}

		if _, exists := triggerIdx[t.ID]; exists {
			issues = append(issues, issue(
				fmt.Sprintf("duplicate trigger ID %d", t.ID),
				map[string]int{"trigger": i},
			))
			continue
		}

		triggerIdx[t.ID] = t
		triggerPathCount[t.ID] = 0
	}

	// --- Phase 1b: validate & index steps ---
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

		// Check for ID collision with triggers
		if _, exists := triggerIdx[s.ID]; exists {
			issues = append(issues, issue(
				fmt.Sprintf("step ID %d collides with trigger ID", s.ID),
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
			reg := ConstructLibrary()
			def, ok := reg.Function(ui.Ref)
			if !ok {
				issues = append(issues, issue(
					fmt.Sprintf("unknown function %q", ui.Ref),
					map[string]int{},
				))
				continue
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

	// Track which steps have parents (for trigger path inference)
	stepsWithParents := make(map[uint64]bool, len(stepIdx))

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

		// Check if parent is a trigger
		if _, isTrigger := triggerIdx[p.ParentID]; isTrigger {
			// Trigger → Step path
			childExecID, ok := idMap[p.ChildID]
			if !ok {
				issues = append(issues, issue(
					fmt.Sprintf("unknown child step %d for trigger path", p.ChildID),
					map[string]int{"path": i},
				))
				continue
			}

			// Validate one trigger can only have a single path
			triggerPathCount[p.ParentID]++
			if triggerPathCount[p.ParentID] > 1 {
				issues = append(issues, issue(
					fmt.Sprintf("trigger %d has multiple outbound paths (only one allowed)", p.ParentID),
					map[string]int{"path": i},
				))
				continue
			}

			// Mark step as having a parent (from trigger)
			stepsWithParents[p.ChildID] = true

			// For trigger paths, the child becomes an entry point
			// (no parent Step added, but we track it has a trigger parent)
			_ = childExecID
			continue
		}

		// Step → Step path
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

		// Mark step as having a parent
		stepsWithParents[p.ChildID] = true

		parent.Children = append(parent.Children, *child)
		child.Parents = append(child.Parents, *parent)
	}

	// --- Phase 3b: infer trigger-to-step paths if needed ---
	// Find steps without parents (entry steps)
	var orphanSteps []uint64
	for stepID := range stepIdx {
		if !stepsWithParents[stepID] {
			orphanSteps = append(orphanSteps, stepID)
		}
	}

	// Find triggers without paths
	var orphanTriggers []uint64
	for triggerID := range triggerIdx {
		if triggerPathCount[triggerID] == 0 {
			orphanTriggers = append(orphanTriggers, triggerID)
		}
	}

	// Infer path: exactly one orphan trigger and exactly one orphan step
	if len(orphanTriggers) > 0 && len(orphanSteps) > 0 {
		if len(orphanTriggers) == 1 && len(orphanSteps) == 1 {
			// Auto-infer the path
			stepsWithParents[orphanSteps[0]] = true
			triggerPathCount[orphanTriggers[0]] = 1
			// The step becomes an entry point (triggered by the single trigger)
		} else if len(orphanTriggers) > 1 || len(orphanSteps) > 1 {
			// Cannot infer - ambiguous
			issues = append(issues, issue(
				fmt.Sprintf("cannot infer trigger-to-step paths: %d triggers and %d entry steps without explicit paths",
					len(orphanTriggers), len(orphanSteps)),
				nil,
			))
		}
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
	reg := ConstructLibrary()

	def, ok := reg.Function(step.Ref)
	if !ok {
		return nil, errors.Internal("unknown function %q", step.Ref)
	}

	if def.Kind != string(step.Kind) {
		return nil, fmt.Errorf("unexpected %s on %s step", def.Kind, step.Kind)
	}

	var (
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

	return types.FunctionStep(&automationTypes.Function{
		Ref:  def.Ref,
		Kind: def.Kind,
		Meta: &automationTypes.FunctionMeta{
			Short:       def.Meta.Short,
			Description: def.Meta.Description,
		},
		Parameters: def.Parameters,
		Results:    def.Results,

		Handler:  def.Handler,
		Iterator: def.Iterator,

		Labels:   def.Labels,
		Disabled: def.Disabled,
	}, step.Arguments, step.Results)
}
