package service

import (
	"context"
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type (
	workflowService struct {
		eventbus  workflowEventTriggerHandler
		store     store.Storer
		actionlog actionlog.Recorder
		ac        workflowAccessController

		log *zap.Logger

		// maps registered triggers (value, uintptr) to trigger ID (key, uint64)
		// this will keep track of all our trigger registrations and help us do a cleanup on
		// trigger update.
		triggers map[uint64]uintptr

		sessions map[uint64]*workflow.Session

		// maps resolved workflow graphs to workflow ID (key, uint64)
		wfgs map[uint64]*workflow.Graph

		mux *sync.RWMutex
	}

	workflowAccessController interface {
		CanCreateWorkflow(context.Context) bool
		CanReadWorkflow(context.Context, *types.Workflow) bool
		CanUpdateWorkflow(context.Context, *types.Workflow) bool
		CanDeleteWorkflow(context.Context, *types.Workflow) bool

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
)

var (
	workflowFuncRegistry map[string]*types.WorkflowFunction
)

func init() {
	workflowFuncRegistry = make(map[string]*types.WorkflowFunction)
	workflowFuncRegistry["HelloWorld"] = &types.WorkflowFunction{
		Ref: "HelloWorld",
		Handler: func(ctx context.Context, variables workflow.Variables) (workflow.Variables, error) {
			println("hello from workflow function;")
			return nil, nil
		},
	}
}

func WorkflowService(log *zap.Logger, _ store.Storer, _ actionlog.Recorder, eb workflowEventTriggerHandler) *workflowService {
	return &workflowService{
		eventbus:  eb,
		actionlog: DefaultActionlog,
		wfgs:      make(map[uint64]*workflow.Graph),
		triggers:  make(map[uint64]uintptr),
		log:       log.Named("WORKFLOW"),
		mux:       &sync.RWMutex{},
	}
}

func (svc *workflowService) Find(ctx context.Context, filter types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error) {
	var (
		wap = &workflowServiceActionProps{filter: &filter}
	)

	err := errors.Internal("pending implementation")

	return nil, filter, svc.recordAction(ctx, wap, WorkflowServiceActionLookup, err)
}

func (svc *workflowService) FindByID(ctx context.Context, workflowID uint64) (wf *types.Workflow, err error) {
	var (
		wap = &workflowServiceActionProps{workflow: &types.Workflow{ID: workflowID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if wf, err = loadWorkflow(ctx, s, workflowID); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, wf); err != nil {
			return err
		}

		return nil
	})

	return wf, svc.recordAction(ctx, wap, WorkflowServiceActionLookup, err)
}

// Create adds new workflow resource and saves it into store
// It updates service's cache
func (svc *workflowService) Create(ctx context.Context, new *types.Workflow) (wf *types.Workflow, err error) {
	var (
		wap   = &workflowServiceActionProps{new: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		wf = &types.Workflow{
			ID:           nextID(),
			Handle:       new.Handle,
			Labels:       new.Labels,
			Meta:         new.Meta,
			Enabled:      new.Enabled,
			Trace:        new.Trace,
			KeepSessions: new.KeepSessions,

			Scope:    new.Scope,
			Steps:    new.Steps,
			Paths:    new.Paths,
			Triggers: new.Triggers,

			// @todo need to check against access control if current user can modify security descriptor
			RunAs:     new.RunAs,
			OwnedBy:   cUser,
			CreatedAt: *now(),
			CreatedBy: cUser,
		}

		return store.CreateWorkflow(ctx, svc.store, wf)
	})

	return wf, svc.recordAction(ctx, wap, WorkflowServiceActionCreate, err)
}

// Update modifies existing workflow resource in the store
func (svc *workflowService) Update(ctx context.Context, upd *types.Workflow) (wf *types.Workflow, err error) {
	return svc.updater(ctx, upd.ID, WorkflowServiceActionUpdate, svc.handleUpdate(upd))
}

func (svc *workflowService) DeleteByID(ctx context.Context, workflowID uint64) (wf *types.Workflow, err error) {
	return svc.updater(ctx, workflowID, WorkflowServiceActionDelete, svc.handleDelete)
}

func (svc *workflowService) UndeleteByID(ctx context.Context, workflowID uint64) (wf *types.Workflow, err error) {
	return svc.updater(ctx, workflowID, WorkflowServiceActionUndelete, svc.handleUndelete)
}

// Start runs a new workflow
//
// Workflow execution is asynchronous operation.
func (svc *workflowService) Start(ctx context.Context, workflowID uint64, scope workflow.Variables) error {
	defer svc.mux.Unlock()
	svc.mux.Lock()
	return errors.Internal("pending implementation")
}

// Resume resumes suspended session/state
//
// Session can only be resumed by knowing session and state ID. Resume is an asynchronous operation
func (svc *workflowService) Resume(ctx context.Context, sessionID, stateID uint64, scope workflow.Variables) error {
	defer svc.mux.Unlock()
	svc.mux.Lock()
	return errors.Internal("pending implementation")
}

func (svc *workflowService) TEMP() {
	wfID := nextID()
	stepID := nextID()
	wfTmp := &types.Workflow{
		ID:      wfID,
		Enabled: true,
		Steps: types.WorkflowStepSet{
			{ID: stepID, Kind: types.WorkflowStepKindFunction, Ref: "HelloWorld"},
		},
		Triggers: types.WorkflowTriggerSet{
			{
				ID:           nextID(),
				WorkflowID:   wfID,
				StepID:       stepID,
				Enabled:      true,
				ResourceType: "system:sink",
				EventType:    "onRequest",
			},
		},
	}

	if g, err := workflowDefToGraph(expr.Parser(), wfTmp); err != nil {
		panic(err)
	} else {
		svc.wfgs[wfTmp.ID] = g
	}

	svc.registerTriggers(context.TODO(), wfTmp)
}

func (svc workflowService) uniqueCheck(ctx context.Context, wf *types.Workflow) (err error) {
	if wf.Handle != "" {
		if e, _ := store.LookupWorkflowByHandle(ctx, svc.store, wf.Handle); e != nil && e.ID != wf.ID {
			return WorkflowServiceErrHandleNotUnique()
		}
	}

	return nil
}

func (svc workflowService) updater(ctx context.Context, workflowID uint64, action func(...*workflowServiceActionProps) *workflowServiceAction, fn workflowUpdateHandler) (*types.Workflow, error) {
	var (
		changes workflowChanges
		wf      *types.Workflow
		aProps  = &workflowServiceActionProps{workflow: &types.Workflow{ID: workflowID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		wf, err = loadWorkflow(ctx, s, workflowID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, wf); err != nil {
			return err
		}

		aProps.setWorkflow(wf)
		aProps.setUpdate(wf)

		if changes, err = fn(ctx, wf); err != nil {
			return err
		}

		if changes&workflowChanged > 0 {
			if err = store.UpdateWorkflow(ctx, svc.store, wf); err != nil {
				return err
			}
		}

		if changes&workflowLabelsChanged > 0 {
			if err = label.Update(ctx, s, wf); err != nil {
				return
			}
		}

		return err
	})

	return wf, svc.recordAction(ctx, aProps, action, err)
}

func (svc workflowService) handleUpdate(upd *types.Workflow) workflowUpdateHandler {
	return func(ctx context.Context, res *types.Workflow) (changes workflowChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return workflowUnchanged, WorkflowServiceErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return workflowUnchanged, WorkflowServiceErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return workflowUnchanged, err
		}

		if !svc.ac.CanUpdateWorkflow(ctx, res) {
			return workflowUnchanged, WorkflowServiceErrNotAllowedToUpdate()
		}

		if res.Handle != upd.Handle {
			changes |= workflowChanged
			res.Handle = upd.Handle
		}

		if res.Enabled != upd.Enabled {
			changes |= workflowChanged
			res.Enabled = upd.Enabled
		}

		if !reflect.DeepEqual(upd.Meta, res.Meta) {
			changes |= workflowChanged
			res.Meta = upd.Meta
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= workflowLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if res.Trace != upd.Trace {
			changes |= workflowChanged
			res.Trace = upd.Trace
		}

		if res.KeepSessions != upd.KeepSessions {
			changes |= workflowChanged
			res.KeepSessions = upd.KeepSessions
		}

		if !reflect.DeepEqual(upd.Meta, res.Meta) {
			changes |= workflowChanged
			res.Meta = upd.Meta
		}

		if !reflect.DeepEqual(upd.Scope, res.Scope) {
			changes |= workflowChanged
			res.Scope = upd.Scope
		}

		if !reflect.DeepEqual(upd.Steps, res.Steps) {
			changes |= workflowChanged
			res.Steps = upd.Steps
		}

		if !reflect.DeepEqual(upd.Paths, res.Paths) {
			changes |= workflowChanged
			res.Paths = upd.Paths
		}

		if !reflect.DeepEqual(upd.Triggers, res.Triggers) {
			changes |= workflowChanged
			res.Triggers = upd.Triggers
		}

		if res.RunAs != upd.RunAs {
			// @todo need to check against access control if current user can modify security descriptor
			changes |= workflowChanged
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

func (svc workflowService) handleDelete(ctx context.Context, wf *types.Workflow) (workflowChanges, error) {
	if !svc.ac.CanDeleteWorkflow(ctx, wf) {
		return workflowUnchanged, WorkflowServiceErrNotAllowedToDelete()
	}

	if wf.DeletedAt != nil {
		// workflow already deleted
		return workflowUnchanged, nil
	}

	wf.DeletedAt = now()
	return workflowChanged, nil
}

func (svc workflowService) handleUndelete(ctx context.Context, wf *types.Workflow) (workflowChanges, error) {
	if !svc.ac.CanDeleteWorkflow(ctx, wf) {
		return workflowUnchanged, WorkflowServiceErrNotAllowedToUndelete()
	}

	if wf.DeletedAt == nil {
		// workflow not deleted
		return workflowUnchanged, nil
	}

	wf.DeletedAt = nil
	return workflowChanged, nil
}

// registerTriggers registeres workflows triggers to eventbus
//
// It preloads run-as identity and finds a starting step for each trigger
func (svc *workflowService) registerTriggers(ctx context.Context, wf *types.Workflow) {
	svc.unregisterTriggers(wf.Triggers...)

	defer svc.mux.Unlock()
	svc.mux.Lock()

	var (
		wfLog = svc.log.Named("trigger-registration").With(
			zap.Uint64("workflowID", wf.ID),
			zap.String("workflow", wf.Handle),
		)

		g = svc.wfgs[wf.ID]

		runAs intAuth.Identifiable
	)

	if !wf.Enabled {
		wfLog.Debug("workflow disabled")
		return
	}

	if wf.RunAs > 0 {
		if u, err := DefaultUser.FindByID(wf.RunAs); err != nil {
			wfLog.Error("failed to load run-as user", zap.Error(err))
			return
		} else if !u.Valid() {
			wfLog.Error("invalid user used for workflow run-as",
				zap.Uint64("userID", u.ID),
				zap.String("email", u.Email),
			)
		} else {
			runAs = u
		}
	}

	for _, t := range wf.Triggers {
		log := wfLog.With(
			zap.Uint64("triggerID", t.ID),
			zap.Uint64("workflowID", wf.ID),
		)

		if !t.Enabled {
			log.Debug("skipping disabled trigger")
			continue
		}

		var (
			start workflow.Step
		)

		if t.StepID == 0 {
			// starting step is not explicitly workflows on trigger, find orphan step
			switch oo := g.Orphans(); len(oo) {
			case 1:
				start = oo[0]
			case 0:
				log.Error("could not find step without parents")
				continue
			default:
				log.Error("multiple steps without parents")
				continue
			}
		} else if start = g.GetStepByIdentifier(t.StepID); start == nil {
			log.Error("trigger staring step references nonexisting step")
			continue
		}

		var (
			handler = func(handleCtx context.Context, ev eventbus.Event) error {
				println("workflow trigger handler; handling event")
				// @todo how do we find out where the workflow should start??

				// create session scope from predefined workflow scope and trigger input
				var (
					scope = wf.Scope.Merge(t.Input)
				)
				scope["event"] = ev

				if runAs != nil {
					// @todo can we pluck alternative identity from Event?
					//       for example:
					//         - use http auth header and get username
					//         - use from/to/replyTo and use that as an identifier
					handleCtx = intAuth.SetIdentityToContext(handleCtx, runAs)
				}

				session := workflow.NewSession(ctx, g)
				svc.mux.Unlock()
				svc.sessions[session.ID()] = session
				svc.mux.Lock()

				return session.Exec(handleCtx, start, scope)
			}

			ops   = make([]eventbus.HandlerRegOp, 0, len(t.Constraints)+2)
			cnstr eventbus.ConstraintMatcher
			err   error
		)

		ops = append(
			ops,
			eventbus.On(t.EventType),
			eventbus.For(t.ResourceType),
		)

		for _, c := range t.Constraints {
			if cnstr, err = eventbus.ConstraintMaker(c.Name, c.Op, c.Values...); err != nil {
				log.Debug(
					"failed to make constraint for workflow trigger",
					zap.Any("constraint", c),
					zap.Error(err),
				)
			} else {
				ops = append(ops, eventbus.Constraint(cnstr))
			}
		}

		svc.triggers[t.ID] = svc.eventbus.Register(handler, ops...)

		log.Debug("trigger registered",
			zap.String("eventType", t.EventType),
			zap.String("resourceType", t.ResourceType),
			zap.Any("constraints", t.Constraints),
		)
	}
}

func (svc *workflowService) unregisterTriggers(tt ...*types.WorkflowTrigger) {
	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, t := range tt {
		if ptr, has := svc.triggers[t.ID]; has {
			svc.eventbus.Unregister(ptr)
			svc.log.Debug("trigger unregistered", zap.Uint64("triggerID", t.ID), zap.Uint64("workflowID", t.WorkflowID))
		}
	}
}

func loadWorkflow(ctx context.Context, s store.Storer, workflowID uint64) (wf *types.Workflow, err error) {
	if workflowID == 0 {
		return nil, WorkflowServiceErrInvalidID()
	}

	if wf, err = store.LookupWorkflowByID(ctx, s, workflowID); errors.IsNotFound(err) {
		return nil, WorkflowServiceErrNotFound()
	}

	return
}

func workflowDefToGraph(lang gval.Language, def *types.Workflow) (*workflow.Graph, error) {
	var (
		g = workflow.NewGraph()
	)

	for g.Len() < len(def.Steps) {
		progress := false
		for _, step := range def.Steps {
			if g.GetStepByIdentifier(step.ID) != nil {
				// resolved
				continue
			}

			inPaths := make([]*types.WorkflowPath, 0, 8)
			outPaths := make([]*types.WorkflowPath, 0, 8)
			for _, path := range def.Paths {
				if path.ChildID == step.ID {
					outPaths = append(inPaths, path)
				} else if path.ParentID == step.ID {
					inPaths = append(inPaths, path)
				}
			}

			if resolved, err := workflowStepDefConv(g, lang, step, inPaths, outPaths); err != nil {
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
		if g.GetStepByIdentifier(path.ChildID) == nil {
			return nil, errors.Internal("failed to resolve paths for %d", path.ChildID)
		}

		if g.GetStepByIdentifier(path.ParentID) == nil {
			return nil, errors.Internal("failed to resolve paths for %d", path.ParentID)
		}

		g.AddParent(
			g.GetStepByIdentifier(path.ChildID),
			g.GetStepByIdentifier(path.ParentID),
		)
	}

	return g, nil
}

// converts all step definitions into workflow.Step instances
//
// if this func returns nil for step and error, assume unresolved dependencies
func workflowStepDefConv(g *workflow.Graph, lang gval.Language, s *types.WorkflowStep, in, out []*types.WorkflowPath) (bool, error) {
	conv, err := func() (workflow.Step, error) {
		switch s.Kind {
		case types.WorkflowStepKindExpressions:
			return workflowExprDefConv(lang, s.Arguments...)

		case types.WorkflowStepKindGatewayFork:
			return workflow.ForkGateway(), nil

		case types.WorkflowStepKindGatewayJoin:
			var (
				ss []workflow.Step
			)
			for _, p := range in {
				if parent := g.GetStepByIdentifier(p.ParentID); parent != nil {
					ss = append(ss, parent)
				} else {
					// unresolved parent, come back later.
					return nil, nil
				}
			}

			return workflow.JoinGateway(ss...), nil

		case types.WorkflowStepKindGatewayIncl, types.WorkflowStepKindGatewayExcl:
			var (
				pp []*workflow.GatewayPath
			)

			for _, p := range in {
				child := g.GetStepByIdentifier(p.ChildID)
				if child == nil {
					return nil, nil
				}

				p, err := workflow.NewGatewayPath(child, p.Test)
				if err != nil {
					return nil, err
				} else {
					pp = append(pp, p)
				}
			}

			if s.Kind == types.WorkflowStepKindGatewayExcl {
				return workflow.ExclGateway(pp...)
			} else {
				return workflow.InclGateway(pp...)
			}

		case types.WorkflowStepKindFunction:
			if s.Ref == "" {
				return nil, errors.Internal("function reference missing")
			}

			if fn, has := workflowFuncRegistry[s.Ref]; !has {
				return nil, errors.Internal("function reference missing")
			} else {
				var (
					err    error
					aa, rr *workflow.Expressions
				)

				if aa, err = workflowExprDefConv(lang, s.Arguments...); err != nil {
					return nil, errors.Internal("failed to convert function arguments: %w", err)
				}

				if rr, err = workflowExprDefConv(lang, s.Results...); err != nil {
					return nil, errors.Internal("failed to convert function arguments: %w", err)
				}

				return workflow.Activity(fn.Handler, aa, rr), nil
			}

		case types.WorkflowStepKindSubprocess:
			return nil, errors.Internal("pending implementation")

		default:
			return nil, errors.Internal("invalid step kind %q", s.Kind)
		}
	}()

	if err != nil {
		return false, err
	} else if conv != nil {
		g.AddStep(conv)
		g.SetStepIdentifier(conv, s.ID)
		return true, err
	} else {
		// unresolved
		return false, nil
	}
}

func workflowExprDefConv(lang gval.Language, ee ...*types.WorkflowExpression) (*workflow.Expressions, error) {
	var (
		set = workflow.NewExpressions(lang)
	)

	for _, e := range ee {
		if err := set.Set(e.Name, e.Expr); err != nil {
			return nil, err
		}
	}

	return set, nil
}
