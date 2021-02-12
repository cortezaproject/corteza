package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
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

			stepIssues := verifyStep(step)

			if step.Kind == types.WorkflowStepKindVisual {
				// make sure visual steps are skipped
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
			return svc.convErrorStep(s, out)

		case types.WorkflowStepKindTermination:
			return svc.convTerminationStep(out)

		case types.WorkflowStepKindPrompt:
			return svc.convPromptStep(s)

		case types.WorkflowStepKindErrHandler:
			return svc.convErrorHandlerStep(g, out)

		case types.WorkflowStepKindBreak:
			return svc.convBreakStep(out)
		case types.WorkflowStepKindContinue:
			return svc.convContinueStep(out)

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

	return nil, fmt.Errorf("unknown gateway type")
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
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	return types.ExpressionsStep(s.Arguments...), nil
}

// internal debug step that can log entire
func (svc workflowConverter) convDebugStep(s *types.WorkflowStep) (wfexec.Step, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	return types.DebugStep(svc.log), nil
}

func (svc workflowConverter) convFunctionStep(g *wfexec.Graph, s *types.WorkflowStep, out []*types.WorkflowPath) (wfexec.Step, error) {
	if s.Ref == "" {
		return nil, errors.Internal("function reference missing")
	}

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
			if len(out) != 2 {
				return nil, fmt.Errorf("expecting exactly two paths (next, exit) out of iterator function step")
			}

			if def.Iterator == nil {
				return nil, errors.Internal("iterator handler for %q not set", s.Ref)
			}
		} else {
			if def.Handler == nil {
				return nil, errors.Internal("function handler for %q not set", s.Ref)
			}
		}

		if err = svc.parseExpressions(s.Arguments...); err != nil {
			return nil, errors.Internal("failed to parse argument expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		} else if err = def.Parameters.VerifyArguments(s.Arguments); err != nil {
			return nil, errors.Internal("failed to verify argument expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		}

		if err = svc.parseExpressions(s.Results...); err != nil {
			return nil, errors.Internal("failed to parse result expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		} else if err = def.Results.VerifyResults(s.Results); err != nil {
			return nil, errors.Internal("failed to verify result expressions for %s %s: %s", s.Kind, s.Ref, err).Wrap(err)
		}

		if isIterator {
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
func (svc workflowConverter) convErrorStep(s *types.WorkflowStep, out types.WorkflowPathSet) (wfexec.Step, error) {
	const (
		argName = "message"
	)

	if len(out) > 0 {
		return nil, errors.Internal("error step must be last step in branch")
	}

	var (
		args = types.ExprSet(s.Arguments)
	)

	if msgArg := args.GetByTarget(argName); msgArg == nil {
		return nil, errors.Internal("error step must have %s argument", argName)
	} else if msgArg.Type != (expr.String{}).Type() {
		return nil, errors.Internal("%s argument on error step must be string, got type '%s'", argName, msgArg.Type)
	} else if len(args) > 1 {
		return nil, errors.Internal("too many arguments on error step")
	}

	if err := svc.parseExpressions(args...); err != nil {
		return nil, err
	}

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

// converts prompt definition to wfexec.Step
func (svc workflowConverter) convTerminationStep(out types.WorkflowPathSet) (wfexec.Step, error) {
	if len(out) > 0 {
		return nil, errors.Internal("termination step must be last step in branch")
	}

	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		return wfexec.Termination(), nil
	}), nil
}

// converts prompt definition to wfexec.Step
func (svc workflowConverter) convPromptStep(s *types.WorkflowStep) (wfexec.Step, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	// Use expression step as base for prompt step
	return types.PromptStep(s.Ref, types.ExpressionsStep(s.Arguments...)), nil
}

func (svc workflowConverter) convBreakStep(out types.WorkflowPathSet) (wfexec.Step, error) {
	if len(out) > 0 {
		return nil, errors.Internal("break step must be last step in branch")
	}

	return wfexec.NewGenericStep(func(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
		return wfexec.LoopBreak(), nil
	}), nil

}

func (svc workflowConverter) convContinueStep(out types.WorkflowPathSet) (wfexec.Step, error) {
	if len(out) > 0 {
		return nil, errors.Internal("continue step must be last step in branch")
	}

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

func verifyStep(step *types.WorkflowStep) types.WorkflowIssueSet {
	var (
		ii = types.WorkflowIssueSet{}

		noArgs = func(s *types.WorkflowStep) error {
			if len(s.Arguments) > 0 {
				return errors.Internal("%s step does not accept arguments", s.Kind)
			}

			return nil
		}

		noResults = func(s *types.WorkflowStep) error {
			if len(s.Results) > 0 {
				return errors.Internal("%s step does not accept results", s.Kind)
			}

			return nil
		}

		checks = make([]func(s *types.WorkflowStep) error, 0)
	)

	switch step.Kind {
	case types.WorkflowStepKindErrHandler:
		checks = append(checks, noArgs, noResults)

	case types.WorkflowStepKindDebug:
		checks = append(checks, noResults)

	case types.WorkflowStepKindVisual:
		checks = append(checks, noArgs, noResults)

	case types.WorkflowStepKindExpressions:
		checks = append(checks, noResults, func(s *types.WorkflowStep) error {
			if len(s.Arguments) == 0 {
				return errors.Internal("%s step require at least one argument", s.Kind)
			}

			return nil
		})

	case types.WorkflowStepKindGateway:
		checks = append(checks, noArgs, noResults)

	case types.WorkflowStepKindError:
		checks = append(checks, noResults)

	case types.WorkflowStepKindTermination:
		checks = append(checks, noArgs, noResults)

	case types.WorkflowStepKindFunction, types.WorkflowStepKindIterator:

	case types.WorkflowStepKindPrompt:
		checks = append(checks, noResults)

	case types.WorkflowStepKindBreak:
		checks = append(checks, noArgs, noResults)

	case types.WorkflowStepKindContinue:
		checks = append(checks, noArgs, noResults)

	default:
		return ii.Append(fmt.Errorf("unknown step kind"), nil)
	}

	for _, check := range checks {
		if err := check(step); err != nil {
			ii = ii.Append(err, nil)
		}
	}

	return ii
}
