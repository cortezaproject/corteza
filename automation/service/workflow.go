package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"sync"
)

type (
	workflow struct {
		eventbus  workflowEventTriggerHandler
		store     store.Storer
		actionlog actionlog.Recorder
		ac        workflowAccessController
		triggers  *trigger

		log *zap.Logger

		// maps resolved workflow graphs to workflow ID (key, uint64)
		wfgs map[uint64]*wfexec.Graph

		// workflow function registry
		reg *registry

		mux    *sync.RWMutex
		parser expr.Parsable
	}

	workflowAccessController interface {
		CanAccess(context.Context) bool

		CanCreateWorkflow(context.Context) bool
		CanReadWorkflow(context.Context, *types.Workflow) bool
		CanUpdateWorkflow(context.Context, *types.Workflow) bool
		CanDeleteWorkflow(context.Context, *types.Workflow) bool
		CanUndeleteWorkflow(context.Context, *types.Workflow) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}

	workflowEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	workflowUpdateHandler func(ctx context.Context, ns *types.Workflow) (workflowChanges, error)
	workflowChanges       uint8
)

const (
	workflowUnchanged     workflowChanges = 0
	workflowChanged       workflowChanges = 1
	workflowLabelsChanged workflowChanges = 2
	workflowDefChanged    workflowChanges = 4
)

func Workflow(log *zap.Logger) *workflow {
	return &workflow{
		log:       log,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		triggers:  DefaultTrigger,
		eventbus:  eventbus.Service(),
		wfgs:      make(map[uint64]*wfexec.Graph),
		mux:       &sync.RWMutex{},
		parser:    expr.NewParser(),
		reg:       Registry(),
	}
}

func (svc *workflow) Search(ctx context.Context, filter types.WorkflowFilter) (rr types.WorkflowSet, f types.WorkflowFilter, err error) {
	var (
		wap = &workflowActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Workflow) (bool, error) {
		if !svc.ac.CanReadWorkflow(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() (err error) {
		if filter.Deleted > 0 {
			// If list with deleted or suspended users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted or archived workflows
			if !svc.ac.CanAccess(ctx) {
				return WorkflowErrNotAllowedToSearch()
			}
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Workflow{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if rr, f, err = store.SearchAutomationWorkflows(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledWorkflows(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, filter, svc.recordAction(ctx, wap, WorkflowActionSearch, err)
}

func (svc *workflow) LookupByID(ctx context.Context, workflowID uint64) (wf *types.Workflow, err error) {
	var (
		wap = &workflowActionProps{workflow: &types.Workflow{ID: workflowID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if wf, err = loadWorkflow(ctx, s, workflowID); err != nil {
			return err
		}

		if !svc.ac.CanReadWorkflow(ctx, wf) {
			return WorkflowErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, wf); err != nil {
			return err
		}

		return nil
	})

	return wf, svc.recordAction(ctx, wap, WorkflowActionLookup, err)
}

// Create adds new workflow resource and saves it into store
// It updates service's cache
func (svc *workflow) Create(ctx context.Context, new *types.Workflow) (wf *types.Workflow, err error) {
	var (
		wap   = &workflowActionProps{new: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanCreateWorkflow(ctx) {
			return WorkflowErrNotAllowedToCreate()
		}

		if !handle.IsValid(new.Handle) {
			return WorkflowErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		if err = validateSteps(new.Steps...); err != nil {
			return
		}

		wf = &types.Workflow{
			ID:           nextID(),
			Handle:       new.Handle,
			Labels:       new.Labels,
			Meta:         new.Meta,
			Enabled:      new.Enabled,
			Trace:        new.Trace,
			KeepSessions: new.KeepSessions,

			Scope: new.Scope,
			Steps: new.Steps,
			Paths: new.Paths,

			// @todo need to check against access control if current user can modify security descriptor
			RunAs:     new.RunAs,
			OwnedBy:   cUser,
			CreatedAt: *now(),
			CreatedBy: cUser,
		}

		if err = store.CreateAutomationWorkflow(ctx, s, wf); err != nil {
			return
		}

		if err = label.Create(ctx, s, wf); err != nil {
			return
		}

		return
	})

	return wf, svc.recordAction(ctx, wap, WorkflowActionCreate, err)
}

// Update modifies existing workflow resource in the store
func (svc *workflow) Update(ctx context.Context, upd *types.Workflow) (*types.Workflow, error) {
	return svc.updater(ctx, upd.ID, WorkflowActionUpdate, func(ctx context.Context, res *types.Workflow) (workflowChanges, error) {
		if !svc.ac.CanUpdateWorkflow(ctx, res) {
			return workflowUnchanged, WorkflowErrNotAllowedToUpdate()
		}

		handler := svc.handleUpdate(upd)
		return handler(ctx, res)
	})
}

func (svc *workflow) DeleteByID(ctx context.Context, workflowID uint64) error {
	return trim1st(svc.updater(ctx, workflowID, WorkflowActionDelete, svc.handleDelete))
}

func (svc *workflow) UndeleteByID(ctx context.Context, workflowID uint64) error {
	return trim1st(svc.updater(ctx, workflowID, WorkflowActionUndelete, svc.handleUndelete))
}

// Start runs a new workflow
//
// Workflow execution is asynchronous operation.
func (svc *workflow) Start(ctx context.Context, workflowID uint64, scope *expr.Vars) error {
	defer svc.mux.Unlock()
	svc.mux.Lock()
	return errors.Internal("pending implementation")
}

func (svc workflow) uniqueCheck(ctx context.Context, res *types.Workflow) (err error) {
	if res.Handle != "" {
		if e, _ := store.LookupAutomationWorkflowByHandle(ctx, svc.store, res.Handle); e != nil && e.ID != res.ID {
			return WorkflowErrHandleNotUnique()
		}
	}

	return nil
}

func (svc workflow) updater(ctx context.Context, workflowID uint64, action func(...*workflowActionProps) *workflowAction, fn workflowUpdateHandler) (*types.Workflow, error) {
	var (
		changes workflowChanges
		res     *types.Workflow
		aProps  = &workflowActionProps{workflow: &types.Workflow{ID: workflowID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		res, err = loadWorkflow(ctx, s, workflowID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setWorkflow(res)
		aProps.setUpdate(res)

		if changes, err = fn(ctx, res); err != nil {
			return err
		}

		if changes&workflowChanged > 0 {
			if err = store.UpdateAutomationWorkflow(ctx, svc.store, res); err != nil {
				return err
			}

		}

		if changes&workflowDefChanged > 0 {
			if err = svc.triggers.registerWorkflows(ctx, res); err != nil {
				return err
			}
		}

		if changes&workflowLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		return err
	})

	return res, svc.recordAction(ctx, aProps, action, err)
}

func (svc workflow) handleUpdate(upd *types.Workflow) workflowUpdateHandler {
	return func(ctx context.Context, res *types.Workflow) (changes workflowChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return workflowUnchanged, WorkflowErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return workflowUnchanged, WorkflowErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return workflowUnchanged, err
		}

		if !svc.ac.CanUpdateWorkflow(ctx, res) {
			return workflowUnchanged, WorkflowErrNotAllowedToUpdate()
		}

		if res.Handle != upd.Handle {
			changes |= workflowChanged
			res.Handle = upd.Handle
		}

		if res.Enabled != upd.Enabled {
			changes |= workflowChanged | workflowDefChanged
			res.Enabled = upd.Enabled
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= workflowLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if res.Trace != upd.Trace {
			changes |= workflowChanged | workflowDefChanged
			res.Trace = upd.Trace
		}

		if res.KeepSessions != upd.KeepSessions {
			changes |= workflowChanged | workflowDefChanged
			res.KeepSessions = upd.KeepSessions
		}

		if upd.Meta != nil {
			if !reflect.DeepEqual(upd.Meta, res.Meta) {
				changes |= workflowChanged
				res.Meta = upd.Meta
			}
		}

		if upd.Scope != nil {
			if !reflect.DeepEqual(upd.Scope, res.Scope) {
				changes |= workflowChanged | workflowDefChanged
				res.Scope = upd.Scope
			}
		}

		if upd.Steps != nil {
			if !reflect.DeepEqual(upd.Steps, res.Steps) {
				if err = validateSteps(upd.Steps...); err != nil {
					return
				}

				changes |= workflowChanged | workflowDefChanged
				res.Steps = upd.Steps
			}
		}

		if upd.Paths != nil {
			if !reflect.DeepEqual(upd.Paths, res.Paths) {
				changes |= workflowChanged | workflowDefChanged
				res.Paths = upd.Paths
			}
		}

		if res.RunAs != upd.RunAs {
			// @todo need to check against access control if current user can modify security descriptor
			changes |= workflowChanged | workflowDefChanged
			res.RunAs = upd.RunAs
		}

		if res.OwnedBy != upd.OwnedBy {
			// @todo need to check against access control if current user can modify owner
			changes |= workflowChanged
			res.OwnedBy = upd.OwnedBy
		}

		if changes&workflowChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc workflow) handleDelete(ctx context.Context, res *types.Workflow) (workflowChanges, error) {
	if !svc.ac.CanDeleteWorkflow(ctx, res) {
		return workflowUnchanged, WorkflowErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		// workflow already deleted
		return workflowUnchanged, nil
	}

	res.DeletedAt = now()
	return workflowChanged, nil
}

func (svc workflow) handleUndelete(ctx context.Context, res *types.Workflow) (workflowChanges, error) {
	if !svc.ac.CanDeleteWorkflow(ctx, res) {
		return workflowUnchanged, WorkflowErrNotAllowedToUndelete()
	}

	if res.DeletedAt == nil {
		// workflow not deleted
		return workflowUnchanged, nil
	}

	res.DeletedAt = nil
	return workflowChanged, nil
}

func (svc *workflow) Load(ctx context.Context) error {
	wwf, _, err := store.SearchAutomationWorkflows(ctx, svc.store, types.WorkflowFilter{
		Deleted:  filter.StateInclusive,
		Disabled: filter.StateExcluded,
	})

	if err != nil {
		return err
	}

	return svc.triggers.registerWorkflows(ctx, wwf...)
}

// Converts workflow definition to wf execution graph
func (svc *workflow) toGraph(def *types.Workflow) (*wfexec.Graph, error) {
	var (
		g = wfexec.NewGraph()
	)

	for g.Len() < len(def.Steps) {
		progress := false
		for _, step := range def.Steps {
			if g.StepByID(step.ID) != nil {
				// resolved
				continue
			}

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
				return nil, err
			} else if resolved {
				progress = true
			}
		}

		if !progress {
			// nothing resolved
			return nil, errors.Internal("failed to resolve workflow step dependencies")
		}
	}

	for _, path := range def.Paths {
		if g.StepByID(path.ChildID) == nil {
			return nil, errors.Internal("failed to resolve step with ID %d", path.ChildID)
		}

		if g.StepByID(path.ParentID) == nil {
			return nil, errors.Internal("failed to resolve step with ID %d", path.ParentID)
		}

		g.AddParent(
			g.StepByID(path.ChildID),
			g.StepByID(path.ParentID),
		)
	}

	return g, nil
}

// converts all step definitions into workflow.Step instances
//
// if this func returns nil for step and error, assume unresolved dependencies
func (svc *workflow) workflowStepDefConv(g *wfexec.Graph, s *types.WorkflowStep, in, out []*types.WorkflowPath) (bool, error) {
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

		//case types.WorkflowStepKindMessage:
		//	return svc.convMessageStep(s)

		case types.WorkflowStepKindPrompt:
			return svc.convPromptStep(s)

		case types.WorkflowStepKindErrHandler:
			return svc.convErrorHandlerStep(g, out)

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

func (svc *workflow) convGateway(g *wfexec.Graph, s *types.WorkflowStep, in, out []*types.WorkflowPath) (wfexec.Step, error) {
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

func (svc *workflow) convErrorHandlerStep(g *wfexec.Graph, out []*types.WorkflowPath) (wfexec.Step, error) {
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

func (svc *workflow) convExpressionStep(s *types.WorkflowStep) (wfexec.Step, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	return types.ExpressionsStep(s.Arguments...), nil
}

// internal debug step that can log entire
func (svc *workflow) convDebugStep(s *types.WorkflowStep) (wfexec.Step, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	return types.DebugStep(svc.log), nil
}

func (svc *workflow) convFunctionStep(g *wfexec.Graph, s *types.WorkflowStep, out []*types.WorkflowPath) (wfexec.Step, error) {
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

// converts prompt definition to wfexec.Step
func (svc *workflow) convPromptStep(s *types.WorkflowStep) (wfexec.Step, error) {
	if err := svc.parseExpressions(s.Arguments...); err != nil {
		return nil, err
	}

	// Use expression step as base for prompt step
	return types.PromptStep(s.Ref, types.ExpressionsStep(s.Arguments...)), nil
}

func (svc *workflow) parseExpressions(ee ...*types.Expr) (err error) {
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

func validateSteps(ss ...*types.WorkflowStep) error {
	var (
		IDs = make(map[uint64]int)

		noArgs = func(i int, s *types.WorkflowStep) error {
			if len(s.Arguments) > 0 {
				return errors.Internal("%s step (ID=%d, position=%d) does not accept arguments", s.Kind, s.ID, i)
			}

			return nil
		}

		reqArgs = func(i int, s *types.WorkflowStep) error {
			if len(s.Arguments) == 0 {
				return errors.Internal("%s step (ID=%d, position=%d) require defined arguments", s.Kind, s.ID, i)
			}

			return nil
		}

		noResults = func(i int, s *types.WorkflowStep) error {
			if len(s.Results) > 0 {
				return errors.Internal("%s step (ID=%d, position=%d) does not accept results", s.Kind, s.ID, i)
			}

			return nil
		}

		checks = make([]func(i int, s *types.WorkflowStep) error, 0)
	)

	for i, s := range ss {
		if p, has := IDs[s.ID]; has {
			return fmt.Errorf("duplicate step ID (%d) used for steps on positions %d and %d", s.ID, p, i)
		}

		IDs[s.ID] = i
		checks = nil

		switch s.Kind {
		case types.WorkflowStepKindErrHandler:

		case types.WorkflowStepKindDebug:
			checks = append(checks, noResults)

		case types.WorkflowStepKindVisual:
			checks = append(checks, noArgs, noResults)

		case types.WorkflowStepKindExpressions:
			checks = append(checks, reqArgs, noResults)

		case types.WorkflowStepKindGateway:
			checks = append(checks, noArgs, noResults)

		case types.WorkflowStepKindFunction, types.WorkflowStepKindIterator:

		case types.WorkflowStepKindPrompt:
			checks = append(checks, noResults)

		default:
			return errors.Internal("unknown step kind (ID=%d, position=%d)", s.ID, i)
		}

		for _, check := range checks {
			if err := check(i, s); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadWorkflow(ctx context.Context, s store.Storer, workflowID uint64) (res *types.Workflow, err error) {
	if workflowID == 0 {
		return nil, WorkflowErrInvalidID()
	}

	if res, err = store.LookupAutomationWorkflowByID(ctx, s, workflowID); errors.IsNotFound(err) {
		return nil, WorkflowErrNotFound()
	}

	return
}

// toLabeledWorkflows converts to []label.LabeledResource
func toLabeledWorkflows(set []*types.Workflow) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func exprTypeSetter(reg *registry, e *types.Expr) func(string) (expr.Type, error) {
	return func(name string) (expr.Type, error) {
		if name == "" {
			name = "Any"
		}

		if typ := reg.Type(name); typ != nil {
			return typ, nil
		} else {
			return nil, errors.NotFound(
				"unknown or unregistered type %q used for expression %q on %q",
				name,
				e.Expr,
				e.Target,
			)
		}
	}
}
