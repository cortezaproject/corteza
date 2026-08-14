package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/automation/types"
	automationTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/ast"
	execTypes "github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/id"
)

// ConvertNgAutomation converts service lvl structs into pkg/automation_exec
//
// @todo report issues so we don't brick the system in case of semantic errors
func ConvertNgAutomation(ctx context.Context, svc *ngAutomation, a *automationTypes.NgAutomation) (execTypes.Executable, automationTypes.NgAutomationIssueSet) {
	var issues automationTypes.NgAutomationIssueSet

	if a == nil {
		return execTypes.Executable{}, automationTypes.NgAutomationIssueSet{
			issue(automationTypes.IssueCodeAutomationNil, automationTypes.NgAutomationSeverityError),
		}
	}

	triggerIdx, triggerPathCount, triggerIssues := indexTriggers(a.Triggers)
	issues = append(issues, triggerIssues...)

	stepIdx, stepIssues := indexSteps(a.Steps, triggerIdx)
	issues = append(issues, stepIssues...)

	issues = append(issues, validateScopeRefs(a.Steps, a.Triggers)...)

	exSteps, idMap, buildIssues := buildExecSteps(svc, stepIdx)
	issues = append(issues, buildIssues...)

	exByID := make(map[id.ID]*execTypes.Step, len(exSteps))
	for i := range exSteps {
		exByID[exSteps[i].ID] = &exSteps[i]
	}

	wireIssues := wirePaths(a.Paths, stepIdx, triggerIdx, triggerPathCount, idMap, exByID)
	issues = append(issues, wireIssues...)

	// Auto-inject a synthetic termination step for any leaf steps that have no
	// children and aren't already a termination step. This makes the termination
	// step optional in the automation config.
	exSteps = ensureTermination(exSteps, exByID)

	if !hasEntry(exByID) {
		issues = append(issues, issue(automationTypes.IssueCodeGraphNoEntry, automationTypes.NgAutomationSeverityError))
	}

	if ring, err := detectCycle(exByID); err != nil {
		ids := make([]uint64, len(ring))
		for i, n := range ring {
			ids[i] = n.Num()
		}
		issues = append(issues, issue(automationTypes.IssueCodeGraphCycle, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailCycle(ids...)))
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
	firstIdx := make(map[uint64]int, len(triggers))

	for i := range triggers {
		t := triggers[i]

		if t.ID == 0 {
			issues = append(issues, issue(automationTypes.IssueCodeTriggerEmptyID, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailEmptyField("trigger", i, "")))
			continue
		}

		if _, exists := idx[t.ID]; exists {
			issues = append(issues, issue(automationTypes.IssueCodeTriggerDuplicateID, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailDuplicateID("trigger", t.ID, firstIdx[t.ID], i)))
			continue
		}

		idx[t.ID] = t
		firstIdx[t.ID] = i
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
	firstIdx := make(map[uint64]int, len(steps))

	for i := range steps {
		s := steps[i]

		if s.ID == 0 {
			issues = append(issues, issue(automationTypes.IssueCodeStepEmptyID, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailEmptyField("step", i, "")))
			continue
		}

		if _, exists := idx[s.ID]; exists {
			issues = append(issues, issue(automationTypes.IssueCodeStepDuplicateID, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailDuplicateID("step", s.ID, firstIdx[s.ID], i)))
			continue
		}

		if _, exists := triggerIdx[s.ID]; exists {
			issues = append(issues, issue(automationTypes.IssueCodeStepIDCollision, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailResourceRef("step", s.ID, i)))
			continue
		}

		idx[s.ID] = s
		firstIdx[s.ID] = i
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
					// ArgKey, not ArgumentName — same rule validation uses.
					ArgumentName: e.ArgKey(),
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
		if step.Kind == "function" || step.Kind == "iterator" {
			reg := ConstructLibrary()
			def, ok := reg.Function(step.Ref)
			if !ok {
				issues = append(issues, issue(automationTypes.IssueCodeFunctionUnknown, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailMissingReference("function", step.Ref, uiID, "", 0)))
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
			// Report each conversion failure as an issue. Dropping the step instead
			// stores it cleanly, returns no issues, registers an executable with a hole
			// in it, and executes as "completed" having run nothing — for an unsupported
			// kind, an unknown function ref, a missing or misnamed argument, a wrong type
			// or an unparseable expression alike.
			// -1: stepIdx is a map, so there is no positional index to report.
			issues = append(issues, automationTypes.NewStepConversionIssue(uiID, -1, err))
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
			issues = append(issues, issue(automationTypes.IssueCodePathEmpty, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailResourceRef("path", 0, i)))
			continue
		}

		if p.ParentID == p.ChildID {
			issues = append(issues, issue(automationTypes.IssueCodePathSelfLoop, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailResourceRef("path", 0, i)))
			continue
		}

		if _, isTrigger := triggerIdx[p.ParentID]; isTrigger {
			childExecID, ok := idMap[p.ChildID]
			if !ok {
				issues = append(issues, issue(automationTypes.IssueCodePathUnknownChild, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailMissingReference("step", strconv.FormatUint(p.ChildID, 10), 0, "", 0)))
				continue
			}

			triggerPathCount[p.ParentID]++
			if triggerPathCount[p.ParentID] > 1 {
				issues = append(issues, issue(automationTypes.IssueCodeTriggerMultiPaths, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailResourceRef("path", 0, i)))
				continue
			}

			stepsWithParents[p.ChildID] = true
			_ = childExecID
			continue
		}

		parentExecID, ok := idMap[p.ParentID]
		if !ok {
			issues = append(issues, issue(automationTypes.IssueCodePathUnknownParent, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailMissingReference("step", strconv.FormatUint(p.ParentID, 10), 0, "", 0)))
			continue
		}

		childExecID, ok := idMap[p.ChildID]
		if !ok {
			issues = append(issues, issue(automationTypes.IssueCodePathUnknownChild, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailMissingReference("step", strconv.FormatUint(p.ChildID, 10), 0, "", 0)))
			continue
		}

		parent := exByID[parentExecID]
		child := exByID[childExecID]
		if parent == nil || child == nil {
			issues = append(issues, issue(automationTypes.IssueCodeInternal, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailResourceRef("path", 0, i)))
			continue
		}

		stepsWithParents[p.ChildID] = true

		// Gateway paths are wired below in sorted order.
		parentKind := stepIdx[p.ParentID].Kind
		if isGatewayKind(parentKind) {
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

		// Iterators must map their multiple out-paths to Children sequence (try/exit), not ErrHandlers.
		if pos == 1 && parentKind != "iterator" {
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
				automationTypes.IssueCodeGraphAmbiguousEntry,
				automationTypes.NgAutomationSeverityError,
			))
		}
	}

	return
}

// validateScopeRefs checks that every expression's Scope points to an existing
// step/trigger handle (or the global scope). A dangling Scope means the
// referenced step was removed while a downstream step still pulls from it.
func validateScopeRefs(
	steps automationTypes.NgAutomationStepSet,
	triggers automationTypes.NgAutomationTriggerSet,
) (issues automationTypes.NgAutomationIssueSet) {
	validScopes := map[string]bool{
		"":       true,
		"global": true,
	}
	for _, s := range steps {
		if s.Handle != "" {
			validScopes[s.Handle] = true
		}
	}
	for _, t := range triggers {
		if t.Handle != "" {
			validScopes[t.Handle] = true
		}
	}

	check := func(stepI int, field string, ee []*automationTypes.Expr) {
		for exprI, e := range ee {
			if e == nil || validScopes[e.Scope] {
				continue
			}
			issues = append(issues, issue(
				automationTypes.IssueCodeScopeUnknown,
				automationTypes.NgAutomationSeverityError,
				automationTypes.NewDetailMissingReference("scope", e.Scope, steps[stepI].ID, field, exprI),
			))
		}
	}

	for i := range steps {
		check(i, "arguments", steps[i].Arguments)
		check(i, "results", steps[i].Results)
	}

	return
}

func validateGatewayPaths(gwPaths map[uint64][]gatewayPath) (issues automationTypes.NgAutomationIssueSet) {
	for parentID, gps := range gwPaths {
		if len(gps) < 2 {
			issues = append(issues, issue(automationTypes.IssueCodeGatewayTooFewPaths, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailGatewayPaths(parentID, "tooFew", len(gps), 2)))
		}

		nilCount := 0
		for _, gp := range gps {
			if gp.path.Condition == nil {
				nilCount++
			}
		}

		if nilCount > 1 {
			issues = append(issues, issue(automationTypes.IssueCodeGatewayMultiElse, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailGatewayPaths(parentID, "multipleElse", nilCount, 1)))
		} else if nilCount == 0 {
			issues = append(issues, issue(automationTypes.IssueCodeGatewayNoElse, automationTypes.NgAutomationSeverityError, automationTypes.NewDetailGatewayPaths(parentID, "noElse", 0, 1)))
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

// detectCycle runs a DFS cycle check using only IDs and, on detection, returns
// the loop ring (the gray-path slice from the re-entered node to the current
// node). Uses step.Children slice (value-copies) but depends only on child IDs.
func detectCycle(exByID map[id.ID]*execTypes.Step) (ring []id.ID, err error) {
	const (
		white = 0
		gray  = 1
		black = 2
	)

	color := make(map[id.ID]uint8, len(exByID))
	var path []id.ID

	var visit func(n id.ID) (id.ID, bool)
	visit = func(n id.ID) (id.ID, bool) {
		switch color[n] {
		case gray:
			return n, true
		case black:
			return id.Zero(), false
		}

		color[n] = gray
		path = append(path, n)

		if s := exByID[n]; s != nil {
			for _, c := range s.Children {
				if hit, found := visit(c.ID); found {
					return hit, true
				}
			}
		}

		color[n] = black
		path = path[:len(path)-1]
		return id.Zero(), false
	}

	for sid := range exByID {
		if color[sid] != white {
			continue
		}
		if hit, found := visit(sid); found {
			start := 0
			for i, n := range path {
				if n == hit {
					start = i
					break
				}
			}
			ring = append(ring, path[start:]...)
			return ring, fmt.Errorf("cycle detected at step %v", hit)
		}
	}

	return nil, nil
}

func issue(code, severity string, details ...*automationTypes.NgAutomationIssueDetail) *automationTypes.NgAutomationIssue {
	return automationTypes.NewIssue(code, severity, details...)
}

// ensureTermination auto-wires a single synthetic termination step to all leaf
// steps (steps with no children) that are not already termination steps.
// This makes the explicit termination step optional in the automation config.
func ensureTermination(exSteps []execTypes.Step, exByID map[id.ID]*execTypes.Step) []execTypes.Step {
	var leafIDs []id.ID
	for _, s := range exSteps {
		if s.Kind == "termination" {
			return exSteps // already has one, nothing to do
		}
		if len(s.Children) == 0 {
			leafIDs = append(leafIDs, s.ID)
		}
	}

	if len(leafIDs) == 0 {
		return exSteps
	}

	termHandler, _ := automationTypes.TerminationStep()
	termStep := execTypes.Step{
		ID:      id.MustNumID(id.Next()),
		Handle:  "_auto_termination",
		Kind:    "termination",
		Handler: termHandler,
	}

	for _, leafID := range leafIDs {
		if s, ok := exByID[leafID]; ok {
			s.Children = append(s.Children, termStep)
			termStep.Parents = append(termStep.Parents, *s)
		}
	}

	exByID[termStep.ID] = &termStep
	return append(exSteps, termStep)
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
		case "function", "iterator":
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
			if err = svc.services.parser.ParseEvaluators(e); err != nil {
				return
			}
		}

		e.Type = normalizeExprType(e.Type)
		if err = e.SetType(exprTypeSetter(Registry(), e)); err != nil {
			return err
		}

		for _, t := range e.Tests {
			if err = svc.services.parser.ParseEvaluators(t); err != nil {
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

	// functionStep.ExecN groups evaluated arguments by ArgumentName, and the
	// construct library's parameters are matched against those keys. Validation
	// resolves an argument's name with a fallback to Target (Expr.ArgKey), so
	// anything that binds has to resolve it the same way — otherwise a
	// target-only argument passes VerifyArguments, which even names the
	// parameter it matched, and then binds to the empty key at run time: the
	// step executes with an empty args map and reports "completed" having done
	// nothing.
	//
	// Copies, not in-place: step is the object that gets persisted, and
	// rewriting the caller's payload is a separate concern from binding it.
	// Safe to copy here because parseExpressions has already populated the
	// unexported eval/typ fields, which the shallow copy carries over.
	arguments := resolveArgNames(step.Arguments)

	handler := def.Handler
	kind := types.FunctionKindFunction
	if isIterator {
		// For iterators, we create a handler that implements both ExecN
		// (for the runtime's step execution) and the scheduler's
		// IteratorHandler interface (Start/More/Next).
		iterFn := def.Iterator
		return &ngIteratorStep{
			def:       &def,
			iterFn:    iterFn,
			arguments: arguments,
			results:   step.Results,
		}, nil
	}

	return types.FunctionStep(&automationTypes.Function{
		Ref:  def.Ref,
		Kind: kind,
		Meta: &automationTypes.FunctionMeta{
			Short:       def.Meta.Short,
			Description: def.Meta.Description,
		},
		Parameters: def.Parameters,
		Results:    def.Results,

		ArgsMerger: def.ArgsMerger,
		Handler:    handler,

		Labels:   def.Labels,
		Disabled: def.Disabled,
	}, arguments, step.Results)
}

// resolveArgNames returns a copy of the set with each expression's ArgumentName
// filled in from Expr.ArgKey, so binding keys off the same name validation did.
func resolveArgNames(in automationTypes.ExprSet) automationTypes.ExprSet {
	out := make(automationTypes.ExprSet, len(in))
	for i, e := range in {
		aux := *e
		aux.ArgumentName = e.ArgKey()
		out[i] = &aux
	}

	return out
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

// ngIteratorStep implements execTypes.StepHandler and the scheduler's
// IteratorHandler interface for iterator steps in the ng runtime.
type ngIteratorStep struct {
	def       *automationTypes.ConstructFunction
	iterFn    automationTypes.IteratorHandler
	arguments automationTypes.ExprSet
	results   automationTypes.ExprSet
}

// ExecN evaluates arguments and initializes the underlying wfexec iterator.
// The scheduler calls Start/More/Next after this.
func (s *ngIteratorStep) ExecN(ctx context.Context, r *execTypes.ExecRequest) (execTypes.ExecResponse, error) {
	var args *expr.Vars

	if len(s.arguments) > 0 {
		evaled, err := s.arguments.EvalN(ctx, r.Scope)
		if err != nil {
			return nil, err
		}

		if s.def.ArgsMerger == nil {
			grouped := make(map[string][]expr.TypedValue)
			for i, e := range s.arguments {
				grouped[e.ArgumentName] = append(grouped[e.ArgumentName], evaled[i])
			}

			final := make(map[string]expr.TypedValue)
			for _, p := range s.def.Parameters {
				vals, exists := grouped[p.ArgumentName]
				if !exists {
					continue
				}
				if p.Aggregate {
					var err error
					final[p.ArgumentName], err = expr.NewArray(vals)
					if err != nil {
						return nil, err
					}
				} else {
					final[p.ArgumentName] = vals[0]
				}
			}

			var err error
			args, err = expr.NewVars(final)
			if err != nil {
				return nil, err
			}
		} else {
			var err error
			args, err = s.def.ArgsMerger(ctx, s.arguments, evaled)
			if err != nil {
				return nil, err
			}
		}

		if args != nil {
			if r.Arguments == nil {
				r.Arguments = make(map[string]expr.TypedValue)
			}
			for k, v := range args.GetValue() {
				r.Arguments[k] = v
			}
		}
	}

	ih, err := s.iterFn(ctx, args)
	if err != nil {
		return nil, err
	}

	out := &expr.Vars{}
	err = out.Set("_iter_handler", expr.Must(expr.NewAny(ih)))
	if err != nil {
		return nil, err
	}

	return out, nil
}
