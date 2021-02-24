package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"go.uber.org/zap"
	"strings"
)

type (
	workflowConverter struct {
		// workflow function registry
		reg    *registry
		parser expr.Parsable
		log    *zap.Logger
	}
)

func Convert(wfService *workflow, wf *types.Workflow) (*wfexec.Graph, types.WorkflowIssueSet) {
	conv := &workflowConverter{
		reg:    wfService.reg,
		parser: wfService.parser,
		log:    wfService.log,
	}

	return conv.makeGraph(wf)
}

// Converts workflow definition to wf execution graph
func (svc workflowConverter) makeGraph(def *types.Workflow) (*wfexec.Graph, types.WorkflowIssueSet) {
	var (
		g           = wfexec.NewGraph()
		wfii        = types.WorkflowIssueSet{}
		IDs         = make(map[uint64]int)
		lastResStep *types.WorkflowStep
	)

	// Basic step verification
	for i, s := range def.Steps {
		if _, has := IDs[s.ID]; has {
			wfii = wfii.Append(
				fmt.Errorf("step ID not unique"),
				map[string]int{"step": i, "duplicate": IDs[s.ID]},
			)
		} else {
			IDs[s.ID] = i
		}
	}

	// if we have one or more duplicated-id error we need to abort right away
	// because all further step positions in issue culprit will be invalid
	if len(wfii) > 0 {
		return nil, wfii
	}

	for g.Len() < len(def.Steps) {
		progress := false
		lastResStep = nil
		for _, step := range def.Steps {
			lastResStep = step
			if g.StepByID(step.ID) != nil {
				// resolved
				continue
			}

			// Collect all incoming and outgoing paths
			inPaths := make([]*types.WorkflowPath, 0, 8)
			outPaths := make([]*types.WorkflowPath, 0, 8)
			for _, path := range def.Paths {
				if path.ChildID == step.ID {
					inPaths = append(inPaths, path)
				} else if path.ParentID == step.ID {
					outPaths = append(outPaths, path)
				}
			}

			stepIssues := verifyStep(step, inPaths, outPaths)

			if step.Kind == types.WorkflowStepKindVisual {
				// make sure visual steps are skipped
				continue
			}

			if resolved, err := svc.workflowStepDefConv(g, step, inPaths, outPaths); err != nil {
				switch aux := err.(type) {
				case types.WorkflowIssueSet:
					stepIssues = append(stepIssues, aux...)
					continue
				case error:
					stepIssues = stepIssues.Append(err, nil)
				}
			} else if resolved {
				progress = true
			}

			wfii = append(wfii, stepIssues.SetCulprit("step", IDs[step.ID])...)
		}

		if !progress {
			var culprit = make(map[string]int)
			if lastResStep != nil {
				culprit = map[string]int{"step": IDs[lastResStep.ID]}
			}

			// nothing resolved for 1 cycle
			wfii = wfii.Append(fmt.Errorf("failed to resolve workflow step dependencies"), culprit)
			break
		}
	}

	for pos, path := range def.Paths {
		if g.StepByID(path.ChildID) == nil {
			wfii = wfii.Append(fmt.Errorf("failed to resolve step with ID %d", path.ChildID), map[string]int{"path": pos})
			continue
		}

		if g.StepByID(path.ParentID) == nil {
			wfii = wfii.Append(fmt.Errorf("failed to resolve step with ID %d", path.ParentID), map[string]int{"path": pos})
			continue
		}

		if len(wfii) > 0 {
			// pointless to fill the if there are errors
			continue
		}

		g.AddParent(
			g.StepByID(path.ChildID),
			g.StepByID(path.ParentID),
		)
	}

	if len(wfii) > 0 {
		return nil, wfii
	}

	return g, nil
}

// converts all step definitions into workflow.Step instances
//
// if this func returns nil for step and error, assume unresolved dependencies
func (svc workflowConverter) workflowStepDefConv(g *wfexec.Graph, s *types.WorkflowStep, in, out []*types.WorkflowPath) (bool, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return false, errors.Internal("failed to parse step arguments expressions for %s: %s", s.Kind, err).Wrap(err)
	}

	if err := svc.parseExpressions(s.Results...); err != nil {
		return false, errors.Internal("failed to parse step results expressions for %s: %s", s.Kind, err).Wrap(err)
	}

	conv, err := func() (wfexec.Step, error) {
		switch s.Kind {
		case types.WorkflowStepKindVisual:
			return nil, nil

		case types.WorkflowStepKindDebug:
			return svc.convDebugStep(s)

		case types.WorkflowStepKindExpressions:
			return svc.convExpressionStep(s)

		case types.WorkflowStepKindGateway:
			return svc.convGateway(g, s, in, out)

		case types.WorkflowStepKindFunction, types.WorkflowStepKindIterator:
			return svc.convFunctionStep(g, s, out)

		case types.WorkflowStepKindError:
			return svc.convErrorStep(s)

		case types.WorkflowStepKindTermination:
			return svc.convTerminationStep()

		case types.WorkflowStepKindPrompt:
			return svc.convPromptStep(s)

		case types.WorkflowStepKindDelay:
			return svc.convDelayStep(s)

		case types.WorkflowStepKindErrHandler:
			return svc.convErrorHandlerStep(g, out)

		case types.WorkflowStepKindBreak:
			return svc.convBreakStep()

		case types.WorkflowStepKindContinue:
			return svc.convContinueStep()

		default:
			return nil, errors.Internal("unsupported step kind %q", s.Kind)
		}
	}()

	if err != nil {
		return false, err
	} else if conv != nil {
		conv.SetID(s.ID)
		g.AddStep(conv)
		return true, err
	} else {
		// signal caller that we were unable to
		// resolve definition at the moment
		return false, nil
	}
}

func (svc workflowConverter) convGateway(g *wfexec.Graph, s *types.WorkflowStep, in, out []*types.WorkflowPath) (wfexec.Step, error) {
	switch s.Ref {
	case "fork":
		return wfexec.ForkGateway(), nil

	case "join":
		var (
			ss []wfexec.Step
		)
		for _, p := range in {
			if parent := g.StepByID(p.ParentID); parent != nil {
				ss = append(ss, parent)
			} else {
				// unresolved parent, come back later.
				return nil, nil
			}
		}

		return wfexec.JoinGateway(ss...), nil

	case "incl", "excl":
		var (
			pp []*wfexec.GatewayPath
		)

		for _, c := range out {
			child := g.StepByID(c.ChildID)
			if child == nil {
				return nil, nil
			}

			if len(c.Expr) > 0 {
				if err := svc.parser.ParseEvaluators(c); err != nil {
					return nil, err
				}
			}

			// wrapping with fn to make sure that we're dealing with the right wf path inside gw-path tester fn
			err := func(c types.WorkflowPath) error {
				p, err := wfexec.NewGatewayPath(child, func(ctx context.Context, scope *expr.Vars) (bool, error) {
					if len(c.Expr) == 0 {
						return true, nil
					}

					return c.Test(ctx, scope)
				})

				if err != nil {
					return err
				}

				pp = append(pp, p)
				return nil
			}(*c)

			if err != nil {
				return nil, err
			}
		}

		if s.Ref == "excl" {
			return wfexec.ExclGateway(pp...)
		} else {
			return wfexec.InclGateway(pp...)
		}
	}

	return nil, errors.Internal("unexpected workflow configuration")
}

func (svc workflowConverter) convErrorHandlerStep(g *wfexec.Graph, out []*types.WorkflowPath) (wfexec.Step, error) {
	switch len(out) {
	case 0:
		return nil, fmt.Errorf("expecting at least one path out of error handling step")
	case 1:
		// remove error handler
		return types.ErrorHandlerStep(nil), nil
	case 2:
		errorHandler := g.StepByID(out[1].ChildID)
		if errorHandler == nil {
			// wait for it to be resolved
			return nil, nil
		}

		return types.ErrorHandlerStep(errorHandler), nil

	default:
		// this might be extended in the future to allow different paths using expression
		// but then again, this can be solved by gateway path following the error handling step
		return nil, fmt.Errorf("max 2 paths out of error handling step")
	}
}

func (svc workflowConverter) convExpressionStep(s *types.WorkflowStep) (wfexec.Step, error) {
	return types.ExpressionsStep(s.Arguments...), nil
}

// internal debug step that can log entire
func (svc workflowConverter) convDebugStep(s *types.WorkflowStep) (wfexec.Step, error) {
	return types.DebugStep(svc.log), nil
}

func (svc workflowConverter) convFunctionStep(g *wfexec.Graph, s *types.WorkflowStep, out []*types.WorkflowPath) (wfexec.Step, error) {
	reg := Registry()

	if def := reg.Function(s.Ref); def == nil {
		return nil, errors.Internal("unknown function %q", s.Ref)
	} else {
		if def.Kind != string(s.Kind) {
			return nil, fmt.Errorf("unexpected %s on %s step", def.Kind, s.Kind)
		}

		var (
			err        error
			isIterator = def.Kind == types.FunctionKindIterator
		)

		if isIterator {
			if def.Iterator == nil {
				return nil, errors.Internal("iterator handler for %q not set", s.Ref)
			}
		} else {
			if def.Handler == nil {
				return nil, errors.Internal("function handler for %q not set", s.Ref)
			}
		}

		if err = def.Parameters.VerifyArguments(s.Arguments); err != nil {
			return nil, errors.Internal("failed to verify argument expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		}

		if err = def.Results.VerifyResults(s.Results); err != nil {
			return nil, errors.Internal("failed to verify result expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		}

		if isIterator {
			if len(out) != 2 {
				return nil, fmt.Errorf("expecting exactly 2 outbound paths for iterator")
			}

			var (
				next = g.StepByID(out[0].ChildID)
				exit = g.StepByID(out[1].ChildID)
			)

			if next == nil || exit == nil {
				// wait for steps to be resolved
				return nil, nil
			}

			return types.IteratorStep(def, s.Arguments, s.Results, next, exit)

		} else {
			return types.FunctionStep(def, s.Arguments, s.Results)
		}
	}
}

// creates error step
//
// Expects ZERO outgoing paths and
func (svc workflowConverter) convErrorStep(s *types.WorkflowStep) (wfexec.Step, error) {
	const (
		argName = "message"
	)

	var (
		args = types.ExprSet(s.Arguments)
	)

	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		var (
			msg         string
			result, err = args.Eval(ctx, r.Scope)
		)
		if err != nil {
			return nil, err
		}

		if result.Has(argName) {
			str, _ := expr.NewString(expr.Must(result.Select(argName)))
			msg = str.GetValue()
		} else {
			if aux, is := args.GetByTarget(argName).Value.(string); is {
				msg = aux
			} else {
				msg = "ERROR"
			}
		}

		return nil, errors.Automation(msg)
	}), nil
}

// converts termination definition to wfexec.Step
func (svc workflowConverter) convTerminationStep() (wfexec.Step, error) {
	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		return wfexec.Termination(), nil
	}), nil
}

// converts prompt definition to wfexec.Step
//
// If input is not passed (session was not resumed from suspension after prompt), we send wait-for-input
// result that signals session it needs to suspend and wait for the input.
// At this point, we take the current scope run it through step arguments and store the output with the suspended state
//
// After session is resumed, input should be set (not nil) and evaluated through results.
func (svc workflowConverter) convPromptStep(s *types.WorkflowStep) (wfexec.Step, error) {
	// Use expression step as base for prompt step
	var (
		args = types.ExprSet(s.Arguments)
		res  = types.ExprSet(s.Results)
	)

	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		if r.Input == nil {
			// input is only set (not nil) when session is resumed on prompt step
			// suspend the session and wait for input

			// @todo take scope, eval with arguments and add it to suspended state
			payload, err := args.Eval(ctx, r.Scope)
			if err != nil {
				return nil, err
			}

			var ownerId uint64 = 0
			if i := auth.GetIdentityFromContext(ctx); i != nil {
				ownerId = i.Identity()
			}

			return wfexec.Prompt(ownerId, s.Ref, payload), nil
		}

		results, err := res.Eval(ctx, r.Scope.Merge(r.Input))
		if err != nil {
			return nil, err
		}

		return results, nil
	}), nil
}

// converts delay definition to wfexec.Step
func (svc workflowConverter) convDelayStep(s *types.WorkflowStep) (wfexec.Step, error) {
	return types.DelayStep(s.Arguments), nil
}

func (svc workflowConverter) convBreakStep() (wfexec.Step, error) {
	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		return wfexec.LoopBreak(), nil
	}), nil

}

func (svc workflowConverter) convContinueStep() (wfexec.Step, error) {
	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		return wfexec.LoopContinue(), nil
	}), nil

}

func (svc workflowConverter) parseExpressions(ee ...*types.Expr) (err error) {
	for _, e := range ee {

		if len(strings.TrimSpace(e.Expr)) > 0 {
			if err = svc.parser.ParseEvaluators(e); err != nil {
				return
			}
		}

		if err = e.SetType(exprTypeSetter(svc.reg, e)); err != nil {
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

func verifyStep(s *types.WorkflowStep, in, out types.WorkflowPathSet) types.WorkflowIssueSet {
	const (
		arguments = "argument"
		results   = "result"
		outbound  = "outbound path"
		inbound   = "inbound path"
	)

	var (
		ii = types.WorkflowIssueSet{}

		// checks if count of arguments, results, outbound or inbound paths is between (incl) min and max
		count = func(min, max int, typ string) func() error {
			return func() error {
				var (
					l int
				)

				switch typ {
				case arguments:
					l = len(s.Arguments)
				case results:
					l = len(s.Results)
				case outbound:
					l = len(out)
				case inbound:
					l = len(in)
				}

				switch {
				case max == 0 && min == max && l != min:
					return errors.Internal("%s step does not expect any %ss", s.Kind, typ)

				case max > 0 && min == max && l != min:
					return errors.Internal("%s step expects exactly %d %s(s)", s.Kind, min, typ)

				case l < min:
					return errors.Internal("%s step expects at least %d %s(s)", s.Kind, min, typ)

				case max > 0 && l > max:
					return errors.Internal("%s step expects no more than %d %s(s)", s.Kind, max, typ)

				}

				return nil
			}
		}

		// check if reference is set on the step
		requiredRef = func() error {
			if s.Ref == "" {
				return errors.Internal("%s step expects reference", s.Kind)
			}

			return nil
		}

		// reference should not be set on the step
		noRef = func() error {
			if s.Ref == "" {
				return errors.Internal("%s step expects reference", s.Kind)
			}

			return nil
		}

		// checks if argument is present
		checkArg = func(argName string, typ expr.Type) func() error {
			return func() error {
				msgArg := types.ExprSet(s.Arguments).GetByTarget(argName)
				if msgArg != nil && msgArg.Type != typ.Type() {
					return errors.Internal("%s argument on %s step must be %s, got type '%s'", argName, s.Kind, typ.Type(), msgArg.Type)
				}

				return nil
			}
		}

		// checks if argument is present
		requiredArg = func(argName string, typ expr.Type) func() error {
			return func() error {
				if msgArg := types.ExprSet(s.Arguments).GetByTarget(argName); msgArg == nil {
					return errors.Internal("%s step expects to have '%s' argument", s.Kind, argName)
				}

				return checkArg(argName, typ)()
			}
		}

		// should not have any arguments, results, outbound or inbound paths
		zero = func(typ string) func() error { return count(0, 0, typ) }

		// step should not have any outbound paths
		last = func() error { return count(0, 0, outbound)() }

		// gw sanity check
		gatewayCheck = func(checks ...func() error) []func() error {
			switch s.Ref {
			case "join":
				return append(checks, count(1, -1, inbound))

			case "fork":
			case "incl", "excl":
				return append(checks, count(1, -1, outbound))
			}

			return append(checks, func() error { return fmt.Errorf("unknown gateway type") })
		}

		checks = make([]func() error, 0)
	)

	switch s.Kind {
	case types.WorkflowStepKindErrHandler:
		checks = append(checks,
			noRef,
			zero(arguments),
			zero(results),
			count(1, 2, outbound),
		)

	case types.WorkflowStepKindDebug:
		checks = append(checks,
			noRef,
			zero(results),
			count(0, 1, outbound),
		)

	case types.WorkflowStepKindVisual:
		checks = append(checks,
			zero(results),
			zero(arguments),
		)

	case types.WorkflowStepKindExpressions:
		checks = append(checks,
			noRef,
			zero(results),
			count(1, -1, arguments),
			count(0, 1, outbound),
		)

	case types.WorkflowStepKindGateway:
		checks = append(checks, gatewayCheck(zero(arguments), zero(results))...)

	case types.WorkflowStepKindError:
		checks = append(checks,
			requiredArg("message", expr.String{}),
			count(0, 1, arguments),
			zero(results),
			last,
		)

	case types.WorkflowStepKindTermination:
		checks = append(checks,
			noRef,
			zero(arguments),
			zero(results),
			last,
		)

	case types.WorkflowStepKindFunction:
		checks = append(checks,
			requiredRef,
			count(0, 1, outbound),
		)

	case types.WorkflowStepKindIterator:
		checks = append(checks,
			requiredRef,
			count(2, 2, outbound),
		)

	case types.WorkflowStepKindPrompt:
		checks = append(checks,
			requiredRef,
			count(0, 1, outbound),
		)

	case types.WorkflowStepKindDelay:
		checks = append(checks,
			checkArg("timestamp", expr.DateTime{}),
			checkArg("offset", expr.Duration{}),
			count(1, 1, arguments),
			zero(results),
			count(0, 1, outbound),
		)

	case types.WorkflowStepKindBreak:
		checks = append(checks,
			noRef,
			zero(arguments),
			zero(results),
			count(0, 1, outbound),
			last,
		)

	case types.WorkflowStepKindContinue:
		checks = append(checks,
			noRef,
			zero(arguments),
			zero(results),
			count(0, 1, outbound),
			last,
		)

	case "":
		return ii.Append(fmt.Errorf("missing step kind"), nil)

	default:
		return ii.Append(fmt.Errorf("unknown step kind '%s'", s.Kind), nil)

	}

	for _, check := range checks {
		if err := check(); err != nil {
			ii = ii.Append(err, nil)
		}
	}

	return ii
}
