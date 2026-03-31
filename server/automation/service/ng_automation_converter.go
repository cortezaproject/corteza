package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/automation/types"
	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/ast"
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

	if a == nil {
		return execTypes.Executable{}, automationTypes.NgAutomationIssueSet{
			issue("nil automation", nil),
		}
	}

	triggerIdx, triggerPathCount, triggerIssues := indexTriggers(a.Triggers)
	issues = append(issues, triggerIssues...)

	stepIdx, stepIssues := indexSteps(a.Steps, triggerIdx)
	issues = append(issues, stepIssues...)

	if len(stepIdx) == 0 {
		return execTypes.Executable{}, append(issues,
			issue("no valid steps defined", nil),
		)
	}

	exSteps, idMap, buildIssues := buildExecSteps(svc, stepIdx)
	issues = append(issues, buildIssues...)

	exByID := make(map[id.ID]*execTypes.Step, len(exSteps))
	for i := range exSteps {
		exByID[exSteps[i].ID] = &exSteps[i]
	}

	wireIssues := wirePaths(a.Paths, stepIdx, triggerIdx, triggerPathCount, idMap, exByID)
	issues = append(issues, wireIssues...)

	if !hasEntry(exByID) {
		issues = append(issues, issue("no entry steps (every step has at least one parent)", nil))
	}

	if err := detectCycle(exByID); err != nil {
		issues = append(issues, issue(err.Error(), nil))
	}

	return execTypes.Executable{
		ID:     id.MustNumID(a.ID),
		Handle: a.Handle,
		Steps:  exSteps,
	}, issues
}

func indexTriggers(triggers automationTypes.NgAutomationTriggerSet) (
	idx map[uint64]*automationTypes.NgAutomationTrigger,
	pathCount map[uint64]int,
	issues automationTypes.NgAutomationIssueSet,
) {
	idx = make(map[uint64]*automationTypes.NgAutomationTrigger, len(triggers))
	pathCount = make(map[uint64]int, len(triggers))

	for i := range triggers {
		t := triggers[i]

		if t.ID == 0 {
			issues = append(issues, issue("trigger has empty ID", map[string]int{"trigger": i}))
			continue
		}

		if _, exists := idx[t.ID]; exists {
			issues = append(issues, issue(fmt.Sprintf("duplicate trigger ID %d", t.ID), map[string]int{"trigger": i}))
			continue
		}

		idx[t.ID] = t
		pathCount[t.ID] = 0
	}

	return
}

func indexSteps(
	steps automationTypes.NgAutomationStepSet,
	triggerIdx map[uint64]*automationTypes.NgAutomationTrigger,
) (
	idx map[uint64]*automationTypes.NgAutomationStep,
	issues automationTypes.NgAutomationIssueSet,
) {
	idx = make(map[uint64]*automationTypes.NgAutomationStep, len(steps))

	for i := range steps {
		s := steps[i]

		if s.ID == 0 {
			issues = append(issues, issue("step has empty ID", map[string]int{"step": i}))
			continue
		}

		if _, exists := idx[s.ID]; exists {
			issues = append(issues, issue(fmt.Sprintf("duplicate step ID %d", s.ID), map[string]int{"step": i}))
			continue
		}

		if _, exists := triggerIdx[s.ID]; exists {
			issues = append(issues, issue(fmt.Sprintf("step ID %d collides with trigger ID", s.ID), map[string]int{"step": i}))
			continue
		}

		idx[s.ID] = s
	}

	return
}

func buildExecSteps(
	svc *ngAutomation,
	stepIdx map[uint64]*automationTypes.NgAutomationStep,
) (
	exSteps []execTypes.Step,
	idMap map[uint64]id.ID,
	issues automationTypes.NgAutomationIssueSet,
) {
	exSteps = make([]execTypes.Step, 0, len(stepIdx))
	idMap = make(map[uint64]id.ID, len(stepIdx))

	for uiID, step := range stepIdx {
		stepID := id.MustNumID(uiID)
		idMap[uiID] = stepID

		aux := execTypes.Step{
			ID:          stepID,
			Kind:        step.Kind,
			Handle:      step.Handle,
			Recoverable: step.Recoverable,
			MaxRetries:  step.MaxRetries,
		}

		for _, e := range step.Arguments {
			aux.Arguments = append(aux.Arguments, execTypes.StepArg{
				Expr: &execTypes.Expr{
					ArgumentName: e.ArgumentName,
					Scope:        e.Scope,
					Target:       e.Target,
					Source:       e.Source,
					Expression:   e.Expr,
					Value:        e.Value,
					Type:         e.Type,
				},
			})
		}

		// @todo fugly
		if step.Kind == "function" {
			reg := ConstructLibrary()
			def, ok := reg.Function(step.Ref)
			if !ok {
				issues = append(issues, issue(fmt.Sprintf("unknown function %q", step.Ref), map[string]int{}))
				continue
			}

			step.Results = []*automationTypes.Expr{}
			for _, r := range def.Results {
				step.Results = append(step.Results, &automationTypes.Expr{
					ArgumentName: r.ArgumentName,
					Target:       r.Name,
					Type:         r.Types[0],
					Source:       r.Name,
				})
			}
		}

		for _, e := range step.Results {
			aux.Results = append(aux.Results, execTypes.StepRst{
				ArgumentName: e.ArgumentName,
			})
		}

		var err error
		aux.Handler, err = stepConv(svc, step)
		if err != nil {
			// @todo err handling...
			spew.Dump(err)
			continue
		}

		exSteps = append(exSteps, aux)
	}

	return
}

// gatewayPath pairs a raw path with its slice index for ordered gateway wiring.
type gatewayPath struct {
	path  *automationTypes.NgAutomationPath
	index int
}

func isGatewayKind(kind string) bool {
	return kind == "gatewayExclusive" || kind == "gatewayInclusive"
}

func wirePaths(
	paths automationTypes.NgAutomationPathSet,
	stepIdx map[uint64]*automationTypes.NgAutomationStep,
	triggerIdx map[uint64]*automationTypes.NgAutomationTrigger,
	triggerPathCount map[uint64]int,
	idMap map[uint64]id.ID,
	exByID map[id.ID]*execTypes.Step,
) (issues automationTypes.NgAutomationIssueSet) {
	stepsWithParents := make(map[uint64]bool)

	// Collect non-gateway, non-trigger step paths grouped by parent (in original order).
	// Index 0 = try (normal child), index 1 = catch (ErrHandlerStepID).
	type stepPath struct {
		path  *automationTypes.NgAutomationPath
		pathI int
	}
	stepPaths := make(map[uint64][]stepPath)
	// Collect gateway paths for deferred ordered wiring.
	gwPaths := make(map[uint64][]gatewayPath)
	for i := range paths {
		p := paths[i]
		if _, isTrigger := triggerIdx[p.ParentID]; isTrigger {
			continue
		}
		if s, ok := stepIdx[p.ParentID]; ok && isGatewayKind(s.Kind) {
			gwPaths[p.ParentID] = append(gwPaths[p.ParentID], gatewayPath{path: p, index: i})
			continue
		}
		if _, ok := stepIdx[p.ParentID]; !ok {
			continue
		}
		stepPaths[p.ParentID] = append(stepPaths[p.ParentID], stepPath{path: p, pathI: i})
	}

	issues = append(issues, validateGatewayPaths(gwPaths)...)

	// Wire non-gateway paths by position.
	for i := range paths {
		p := paths[i]

		if p.ParentID == 0 || p.ChildID == 0 {
			issues = append(issues, issue("path has empty parent or child", map[string]int{"path": i}))
			continue
		}

		if p.ParentID == p.ChildID {
			issues = append(issues, issue("path is a self-loop", map[string]int{"path": i}))
			continue
		}

		if _, isTrigger := triggerIdx[p.ParentID]; isTrigger {
			childExecID, ok := idMap[p.ChildID]
			if !ok {
				issues = append(issues, issue(fmt.Sprintf("unknown child step %d for trigger path", p.ChildID), map[string]int{"path": i}))
				continue
			}

			triggerPathCount[p.ParentID]++
			if triggerPathCount[p.ParentID] > 1 {
				issues = append(issues, issue(fmt.Sprintf("trigger %d has multiple outbound paths (only one allowed)", p.ParentID), map[string]int{"path": i}))
				continue
			}

			stepsWithParents[p.ChildID] = true
			_ = childExecID
			continue
		}

		parentExecID, ok := idMap[p.ParentID]
		if !ok {
			issues = append(issues, issue(fmt.Sprintf("unknown parent step %d", p.ParentID), map[string]int{"path": i}))
			continue
		}

		childExecID, ok := idMap[p.ChildID]
		if !ok {
			issues = append(issues, issue(fmt.Sprintf("unknown child step %d", p.ChildID), map[string]int{"path": i}))
			continue
		}

		parent := exByID[parentExecID]
		child := exByID[childExecID]
		if parent == nil || child == nil {
			issues = append(issues, issue("internal step index missing", map[string]int{"path": i}))
			continue
		}

		stepsWithParents[p.ChildID] = true

		// Gateway paths are wired below in sorted order.
		if isGatewayKind(stepIdx[p.ParentID].Kind) {
			continue
		}

		// Determine this path's position among all outbound paths from this parent.
		// Position 0 = try (normal child); position 1 = catch (ErrHandlerStepID).
		pos := -1
		for j, sp := range stepPaths[p.ParentID] {
			if sp.pathI == i {
				pos = j
				break
			}
		}

		if pos == 1 {
			// Second outbound path is the catch handler.
			parent.ErrHandlerStepID = childExecID
		} else {
			// First (or only) outbound path is the normal child edge.
			parent.Children = append(parent.Children, *child)
			child.Parents = append(child.Parents, *parent)
		}
	}

	// Wire gateway paths in sorted order (conditions first, else last).
	for parentID, gps := range gwPaths {
		sorted := sortGatewayPaths(gps)

		parentExecID, ok := idMap[parentID]
		if !ok {
			continue
		}

		parent := exByID[parentExecID]
		if parent == nil {
			continue
		}

		nodeConds := make([]*ast.ASTNode, 0, len(sorted))
		for _, gp := range sorted {
			nodeConds = append(nodeConds, gp.path.Condition)

			childExecID, ok := idMap[gp.path.ChildID]
			if !ok {
				continue
			}

			child := exByID[childExecID]
			if child == nil {
				continue
			}

			parent.Children = append(parent.Children, *child)
			child.Parents = append(child.Parents, *parent)
		}

		if cs, ok := parent.Handler.(automationTypes.GatewayConditionsSetter); ok {
			cs.SetConditions(nodeConds)
		}
	}

	// Infer trigger-to-step path when there is exactly one orphan of each.
	var orphanSteps []uint64
	for stepID := range stepIdx {
		if !stepsWithParents[stepID] {
			orphanSteps = append(orphanSteps, stepID)
		}
	}

	var orphanTriggers []uint64
	for triggerID := range triggerIdx {
		if triggerPathCount[triggerID] == 0 {
			orphanTriggers = append(orphanTriggers, triggerID)
		}
	}

	if len(orphanTriggers) > 0 && len(orphanSteps) > 0 {
		if len(orphanTriggers) == 1 && len(orphanSteps) == 1 {
			stepsWithParents[orphanSteps[0]] = true
			triggerPathCount[orphanTriggers[0]] = 1
		} else {
			issues = append(issues, issue(
				fmt.Sprintf("cannot infer trigger-to-step paths: %d triggers and %d entry steps without explicit paths",
					len(orphanTriggers), len(orphanSteps)),
				nil,
			))
		}
	}

	return
}

func validateGatewayPaths(gwPaths map[uint64][]gatewayPath) (issues automationTypes.NgAutomationIssueSet) {
	for parentID, gps := range gwPaths {
		if len(gps) < 2 {
			issues = append(issues, issue(fmt.Sprintf("gateway step %d must have at least 2 outbound paths, got %d", parentID, len(gps)), nil))
		}

		nilCount := 0
		for _, gp := range gps {
			if gp.path.Condition == nil {
				nilCount++
			}
		}

		if nilCount > 1 {
			issues = append(issues, issue(fmt.Sprintf("gateway step %d has %d else paths (exactly one allowed)", parentID, nilCount), nil))
		} else if nilCount == 0 {
			issues = append(issues, issue(fmt.Sprintf("gateway step %d has no else path", parentID), nil))
		}
	}

	return
}

// sortGatewayPaths returns gps with non-nil conditions first, nil (else) last.
func sortGatewayPaths(gps []gatewayPath) []gatewayPath {
	sorted := make([]gatewayPath, 0, len(gps))
	var elsePath *gatewayPath

	for i := range gps {
		if gps[i].path.Condition == nil {
			elsePath = &gps[i]
		} else {
			sorted = append(sorted, gps[i])
		}
	}

	if elsePath != nil {
		sorted = append(sorted, *elsePath)
	}

	return sorted
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

// normalizeExprType maps common type aliases to canonical automation expr registry names.
func normalizeExprType(t string) string {
	switch strings.ToLower(t) {
	case "number":
		return "Integer"
	}
	return t
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

		case "termination":
			return stepConvTermination(step)

		case "gatewayExclusive":
			return stepConvGateway(step, false)

		case "gatewayInclusive":
			return stepConvGateway(step, true)

		case "error":
			return stepConvError(step)

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

		e.Type = normalizeExprType(e.Type)
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

		ArgsMerger: def.ArgsMerger,
		Handler:    def.Handler,
		Iterator:   def.Iterator,

		Labels:   def.Labels,
		Disabled: def.Disabled,
	}, step.Arguments, step.Results)
}

func stepConvTermination(step *automationTypes.NgAutomationStep) (out execTypes.StepHandler, err error) {
	return automationTypes.TerminationStep()
}

// stepConvGateway creates a gateway handler.
// Conditions are populated in Phase 3 via GatewayConditionsSetter.SetConditions.
func stepConvGateway(step *automationTypes.NgAutomationStep, inclusive bool) (out execTypes.StepHandler, err error) {
	if inclusive {
		return automationTypes.InclusiveGatewayStep(nil), nil
	}
	return automationTypes.ExclusiveGatewayStep(nil), nil
}

// stepConvError creates an error step handler.
// It reads "message" (*ast.ASTNode) and "recoverable" (bool) from step.Meta.
func stepConvError(step *automationTypes.NgAutomationStep) (out execTypes.StepHandler, err error) {
	var message *ast.ASTNode
	if v, ok := step.Meta.Extra["message"]; ok {
		if node, ok := v.(*ast.ASTNode); ok {
			message = node
		}
	}

	recoverable := false
	if v, ok := step.Meta.Extra["recoverable"]; ok {
		if b, ok := v.(bool); ok {
			recoverable = b
		}
	}

	return automationTypes.ErrorStep(message, recoverable), nil
}
