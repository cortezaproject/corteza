package service

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/crusttech/human/server/automation/types"
	intAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/wfexec"
	"github.com/crusttech/human/server/store"
	"go.uber.org/zap"
)

type (
	workflowServices struct {
		eventbus workflowEventTriggerHandler
		triggers *trigger
		session  *session
		opt      options.WorkflowOpt
		log      *zap.Logger
		// cache of workflows, graphs to workflow ID (key, uint64)
		cache    map[uint64]*wfCacheItem
		muxCache *sync.RWMutex
		// handle to workflow index
		wIndex    map[string]uint64
		muxWIndex *sync.RWMutex
		// workflow function registry
		reg         *registry
		corredorOpt options.CorredorOpt
		parser      expr.Parsable
	}

	wfCacheItem struct {
		wf *types.Workflow

		// caching exec graph
		g *wfexec.Graph

		// caching user we'll executing workflow with
		runAs intAuth.Identifiable
	}

	workflowAccessController interface {
		CanSearchWorkflows(context.Context) bool
		CanCreateWorkflow(context.Context) bool
		CanReadWorkflow(context.Context, *types.Workflow) bool
		CanUpdateWorkflow(context.Context, *types.Workflow) bool
		CanDeleteWorkflow(context.Context, *types.Workflow) bool
		CanUndeleteWorkflow(context.Context, *types.Workflow) bool
		CanManageSessionsOnWorkflow(context.Context, *types.Workflow) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error

		workflowExecController
	}

	workflowExecController interface {
		CanExecuteWorkflow(context.Context, *types.Workflow) bool
	}

	workflowEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	workflowInvokerCtxKey struct{}
)

func Workflow(log *zap.Logger, corredorOpt options.CorredorOpt, opt options.WorkflowOpt) *workflow {
	return &workflow{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		services: &workflowServices{
			log:         log,
			opt:         opt,
			triggers:    DefaultTrigger,
			session:     DefaultSession,
			eventbus:    eventbus.Service(),
			cache:       make(map[uint64]*wfCacheItem),
			wIndex:      make(map[string]uint64),
			muxCache:    &sync.RWMutex{},
			muxWIndex:   &sync.RWMutex{},
			parser:      expr.NewParser(),
			reg:         Registry(),
			corredorOpt: corredorOpt,
		},
	}
}

func (svc *workflow) FindByID(ctx context.Context, workflowID uint64) (wf *types.Workflow, err error) {
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

func (svc *workflow) LookupByHandle(ctx context.Context, handle string) (*types.Workflow, error) {
	rr, _, err := svc.Search(ctx, types.WorkflowFilter{Handle: handle})
	if err != nil {
		return nil, err
	}
	if len(rr) == 0 {
		return nil, WorkflowErrNotFound()
	}
	return rr[0], nil
}

// Create adds new workflow resource and saves it into store
// It updates service's cache
func (svc *workflow) guard(ctx context.Context, res *types.Workflow) error {
	return guardProjectWritable(ctx, svc.store, res.ProjectID)
}

func (svc *workflow) Create(ctx context.Context, new *types.Workflow) (wf *types.Workflow, err error) {
	var (
		wap   = &workflowActionProps{workflow: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
		g     *wfexec.Graph
		runAs intAuth.Identifiable
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanCreateWorkflow(ctx) {
			return WorkflowErrNotAllowedToCreate()
		}

		if new.Meta.Name == "" {
			return WorkflowErrMissingName()
		}

		if !handle.IsValid(new.Handle) {
			return WorkflowErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
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

		wap.workflow = wf

		if g, runAs, err = svc.validateWorkflow(ctx, wf); err != nil {
			return
		}

		svc.updateCache(wf, runAs, g)

		if len(wf.Issues) == 0 {
			if err = svc.services.triggers.registerWorkflows(ctx, wf); err != nil {
				return err
			}
		}

		if err = store.CreateAutomationWorkflow(ctx, s, wf); err != nil {
			return
		}

		if err = label.Create(ctx, s, wf); err != nil {
			return
		}

		wap.setNew(wf)

		return
	})

	return wf, svc.recordAction(ctx, wap, WorkflowActionCreate, err)
}

// Update modifies existing workflow resource in the store
func (svc *workflow) Update(ctx context.Context, upd *types.Workflow) (res *types.Workflow, err error) {
	var (
		old    *types.Workflow
		aProps = &workflowActionProps{workflow: &types.Workflow{ID: upd.ID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadWorkflow(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		old = res.Clone()
		aProps.setWorkflow(res)
		aProps.setUpdate(res)

		return svc.onUpdate(ctx, s, upd, res, aProps, nil, nil)
	})

	return res, svc.recordAction(ctx, aProps, WorkflowActionUpdate, err, old, res)
}

func (svc *workflow) DeleteByID(ctx context.Context, workflowID uint64) (err error) {
	var (
		aProps = &workflowActionProps{workflow: &types.Workflow{ID: workflowID}}
		res    *types.Workflow
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadWorkflow(ctx, s, workflowID); err != nil {
			return
		}

		aProps.setWorkflow(res)
		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, WorkflowActionDelete, err)
}

func (svc *workflow) UndeleteByID(ctx context.Context, workflowID uint64) (err error) {
	var (
		aProps = &workflowActionProps{workflow: &types.Workflow{ID: workflowID}}
		res    *types.Workflow
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadWorkflow(ctx, s, workflowID); err != nil {
			return
		}

		aProps.setWorkflow(res)
		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, WorkflowActionUndelete, err)
}

func (svc workflow) uniqueCheck(ctx context.Context, res *types.Workflow) (err error) {
	if res.Handle != "" {
		if e, _ := store.LookupAutomationWorkflowByHandle(ctx, svc.store, res.Handle); e != nil && e.ID != res.ID {
			return WorkflowErrHandleNotUnique()
		}
	}

	return nil
}

// onUpdate applies field changes from upd onto res, validates, updates the cache
// and trigger registrations, then persists.
func (svc *workflow) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Workflow, aProps *workflowActionProps, _ func() error, _ func() error) error {
	if upd.Meta.Name == "" {
		return WorkflowErrMissingName()
	}

	if !svc.ac.CanUpdateWorkflow(ctx, res) {
		return WorkflowErrNotAllowedToUpdate()
	}

	if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
		return WorkflowErrStaleData()
	}

	// res is loaded from the store, so res.Issues holds the persisted issue set.
	// Capture it before validateWorkflow recomputes it so we can detect issue
	// transitions (e.g. issues resolved) and re-persist even when no other field
	// changed — otherwise the read endpoint keeps serving stale issues.
	priorIssues := res.Issues

	if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
		return WorkflowErrInvalidHandle()
	}

	if err := svc.uniqueCheck(ctx, upd); err != nil {
		return err
	}

	changed := false
	labelsChanged := false

	if res.Handle != upd.Handle {
		changed = true
		res.Handle = upd.Handle
	}
	if res.Enabled != upd.Enabled {
		changed = true
		res.Enabled = upd.Enabled
	}
	if upd.Labels != nil && label.Changed(res.Labels, upd.Labels) {
		labelsChanged = true
		res.Labels = upd.Labels
	}
	if res.Trace != upd.Trace {
		changed = true
		res.Trace = upd.Trace
	}
	if res.KeepSessions != upd.KeepSessions {
		changed = true
		res.KeepSessions = upd.KeepSessions
	}
	if upd.Meta != nil && !reflect.DeepEqual(upd.Meta, res.Meta) {
		changed = true
		res.Meta = upd.Meta
	}
	if upd.Scope != nil && !reflect.DeepEqual(upd.Scope, res.Scope) {
		changed = true
		res.Scope = upd.Scope
	}
	if upd.Steps != nil && !reflect.DeepEqual(upd.Steps, res.Steps) {
		changed = true
		res.Steps = upd.Steps
	}
	if upd.Paths != nil && !reflect.DeepEqual(upd.Paths, res.Paths) {
		changed = true
		res.Paths = upd.Paths
	}
	if res.RunAs != upd.RunAs {
		// @todo need to check against access control if current user can modify security descriptor
		changed = true
		res.RunAs = upd.RunAs
	}
	if res.OwnedBy != upd.OwnedBy {
		// @todo need to check against access control if current user can modify owner
		changed = true
		res.OwnedBy = upd.OwnedBy
	}
	if changed {
		res.UpdatedAt = now()
	}

	g, runAs, err := svc.validateWorkflow(ctx, res)
	if err != nil {
		return err
	}

	svc.updateCache(res, runAs, g)

	if len(res.Issues) == 0 {
		if err = svc.services.triggers.registerWorkflows(ctx, res); err != nil {
			return err
		}
	}

	if changed || !reflect.DeepEqual(priorIssues, res.Issues) {
		if err = store.UpdateAutomationWorkflow(ctx, s, res); err != nil {
			return err
		}
	}

	if labelsChanged {
		if err = label.Update(ctx, s, res); err != nil {
			return err
		}
	}

	return nil
}

// onDelete soft-deletes the workflow, updates cache, and persists.
func (svc *workflow) onDelete(ctx context.Context, s store.Storer, res *types.Workflow, aProps *workflowActionProps) error {
	if !svc.ac.CanDeleteWorkflow(ctx, res) {
		return WorkflowErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		// already deleted
		return nil
	}

	res.DeletedAt = now()

	g, runAs, err := svc.validateWorkflow(ctx, res)
	if err != nil {
		return err
	}
	svc.updateCache(res, runAs, g)

	return store.UpdateAutomationWorkflow(ctx, s, res)
}

// onUndelete reverses a soft-delete, updates cache, and persists.
func (svc *workflow) onUndelete(ctx context.Context, s store.Storer, res *types.Workflow, aProps *workflowActionProps) error {
	if !svc.ac.CanUndeleteWorkflow(ctx, res) {
		return WorkflowErrNotAllowedToUndelete()
	}

	if res.DeletedAt == nil {
		// not deleted
		return nil
	}

	res.DeletedAt = nil

	g, runAs, err := svc.validateWorkflow(ctx, res)
	if err != nil {
		return err
	}
	svc.updateCache(res, runAs, g)

	if len(res.Issues) == 0 {
		if err = svc.services.triggers.registerWorkflows(ctx, res); err != nil {
			return err
		}
	}

	return store.UpdateAutomationWorkflow(ctx, s, res)
}

func (svc *workflow) Load(ctx context.Context) error {
	var (
		set, _, err = store.SearchAutomationWorkflows(ctx, svc.store, types.WorkflowFilter{
			Deleted:     filter.StateInclusive,
			Disabled:    filter.StateExcluded,
			SubWorkflow: filter.StateInclusive,
		})
		g     *wfexec.Graph
		runAs intAuth.Identifiable
	)
	if err != nil {
		return err
	}

	svc.services.muxWIndex.Lock()
	defer svc.services.muxWIndex.Unlock()

	for _, wf := range set {
		svc.services.wIndex[wf.Handle] = wf.ID

		if g, runAs, err = svc.validateWorkflow(ctx, wf); err != nil {
			continue
		}

		svc.updateCache(wf, runAs, g)
	}

	return svc.services.triggers.registerWorkflows(ctx, set...)
}

// updateCache
func (svc *workflow) updateCache(wf *types.Workflow, runAs intAuth.Identifiable, g *wfexec.Graph) {
	defer svc.services.muxCache.Unlock()
	svc.services.muxCache.Lock()

	if wf.Executable() {
		svc.services.wIndex[wf.Handle] = wf.ID
		svc.services.cache[wf.ID] = &wfCacheItem{g: g, wf: wf, runAs: runAs}
	} else {
		// remove deleted
		delete(svc.services.cache, wf.ID)
		delete(svc.services.wIndex, wf.Handle)
	}

	return
}

func (svc *workflow) onExec(ctx context.Context, aProps *workflowActionProps, workflowID uint64, p types.WorkflowExecParams) (results *expr.Vars, sessionID uint64, stacktrace types.Stacktrace, err error) {
	var (
		t    *types.Trigger
		wait WaitFn
	)

	svc.services.muxCache.Lock()
	if nil == svc.services.cache[workflowID] || nil == svc.services.cache[workflowID].wf {
		svc.services.muxCache.Unlock()
		return nil, 0, nil, WorkflowErrNotFound()
	}

	wf := svc.services.cache[workflowID].wf
	svc.services.muxCache.Unlock()

	aProps.setWorkflow(wf)

	if !svc.ac.CanExecuteWorkflow(ctx, wf) {
		return nil, 0, nil, WorkflowErrNotAllowedToExecute()
	}

	if !wf.Enabled && !p.Trace {
		return nil, 0, nil, WorkflowErrDisabled()
	}

	// Find the trigger.
	// @todo can we cache this as well?
	t, err = func() (*types.Trigger, error) {
		if p.CallerWorkflowID > 0 {
			// skip triggers checking when executed as sub-workflow
			// @todo be more strict and allow this ONLY when workflow is flagged as a sub-workflow
			return nil, nil
		}

		var tt types.TriggerSet
		// Load triggers directly from the store. At this point we do not care
		// about trigger search or read permissions
		tt, err = loadWorkflowTriggers(ctx, svc.store, workflowID)
		if err != nil {
			return nil, err
		}

		if len(tt) == 0 {
			return nil, nil
		}

		if p.StepID == 0 && len(tt) > 0 {
			return tt[0], nil
		} else {
			for _, tMatch := range tt {
				if tMatch.StepID == p.StepID {
					return tMatch, nil
				}
			}
		}

		if !p.Trace {
			// when not doing a trace (designing the workflow)
			// we need to be more strict and disallow execution of
			// the misconfigured workflows and use of disabled triggers
			if t == nil {
				return nil, WorkflowErrUnknownWorkflowStep()
			} else if !t.Enabled {
				return nil, WorkflowErrDisabled()
			}
		}

		if t != nil {
			aProps.setTrigger(t)
			p.StepID = t.StepID
			p.EventType = t.EventType
			p.ResourceType = t.ResourceType

			// merge with input from trigger
			// with trigger input vars are overwritten by input vars
			p.Input = t.Input.MustMerge(p.Input)
		} else {
			p.EventType = "onTrace"
		}

		return nil, nil
	}()

	if err != nil {
		return
	}

	wait, sessionID, err = svc.exec(ctx, wf, p)

	if err != nil {
		return
	}

	if p.Async && !p.Wait && wf.CheckDeferred() {
		// deferred workflow, return right away and keep the workflow session
		// running without waiting for the execution
		return
	}

	// wait for the workflow to complete
	// reuse scope for results
	// this will be decoded back to event properties
	results, sessionID, _, stacktrace, err = wait(ctx)
	return
}

// validates workflow by trying to convert it to graph and checking assigned triggers
func (svc *workflow) validateWorkflow(ctx context.Context, wf *types.Workflow) (g *wfexec.Graph, runAs intAuth.Identifiable, err error) {
	var (
		tt []*types.Trigger
	)

	g, wf.Issues = Convert(svc, wf)

	tt, _, err = store.SearchAutomationTriggers(ctx, svc.store, types.TriggerFilter{
		WorkflowID: id.Strings(types.WorkflowSet{wf}.IDs()...),
		Deleted:    filter.StateExcluded,
		Disabled:   filter.StateExcluded,
	})

	if err != nil {
		return
	}

	wf.Issues = append(wf.Issues, validateWorkflowTriggers(wf, tt...)...)

	// Returns context with identity set to service user
	//
	// Current user (identity in the context) might not have
	// sufficient privileges to load info about invoker and runner
	sysUserCtx := func() context.Context {
		return intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
	}

	// @todo this might not be the smartest thing, users might get invalidated after
	//       we add cache them as workflow runners
	if wf.RunAs > 0 {
		if runAs, err = DefaultUser.FindByAny(sysUserCtx(), wf.RunAs); err != nil {
			wf.Issues = wf.Issues.Append(fmt.Errorf("failed to load run-as user %d: %w", wf.RunAs, err), nil)
		} else if !runAs.Valid() {
			wf.Issues = wf.Issues.Append(fmt.Errorf("invalid user %d used for workflow run-as", wf.RunAs), nil)
		}
	}

	return
}

func (svc *workflow) exec(ctx context.Context, wf *types.Workflow, p types.WorkflowExecParams) (WaitFn, uint64, error) {
	if wf.Issues != nil {
		return nil, 0, wf.Issues
	}

	defer svc.services.muxCache.Unlock()
	svc.services.muxCache.Lock()

	if svc.services.cache[wf.ID] == nil {
		return nil, 0, WorkflowErrInvalidID()
	}

	var (
		g     = svc.services.cache[wf.ID].g
		runAs = svc.services.cache[wf.ID].runAs

		scope *expr.Vars
	)

	// merge workflow scope with the input
	scope = wf.Scope.MustMerge(p.Input)

	return svc.services.session.Start(ctx, g, types.SessionStartParams{
		Invoker: intAuth.GetIdentityFromContext(ctx),
		Runner:  runAs,

		WorkflowID:   wf.ID,
		KeepFor:      wf.KeepSessions,
		Trace:        wf.Trace || p.Trace,
		Input:        scope,
		StepID:       p.StepID,
		EventType:    p.EventType,
		ResourceType: p.ResourceType,

		CallStack: wfexec.GetContextCallStack(ctx),
	})
}

func makeWorkflowHandler(svc *workflow, wf *types.Workflow, t *types.Trigger) eventbus.HandlerFn {
	return func(ctx context.Context, ev eventbus.Event) (err error) {
		var (
			scope *expr.Vars
		)

		if dec, is := ev.(varsEncoder); is {
			scope, err = dec.EncodeVars()
			if err != nil {
				return
			}
		}

		wait, _, err := svc.exec(ctx, wf, types.WorkflowExecParams{
			StepID:       t.StepID,
			EventType:    t.EventType,
			ResourceType: t.ResourceType,
			Input:        t.Input.MustMerge(scope),

			Trace: wf.Trace,
			Async: false,
		})

		if err != nil {
			return
		}

		if wf.CheckDeferred() {
			// deferred workflow, return right away and keep the workflow session
			// running without waiting for the execution
			return
		}

		// wait for the workflow to complete
		// reuse scope for results
		// this will be decoded back to event properties
		scope, _, _, _, err = wait(ctx)
		if err != nil {
			return
		}

		if dec, is := ev.(varsDecoder); is {
			return dec.DecodeVars(scope)
		}

		return
	}
}

func (svc *workflow) handleToID(h string) uint64 {
	svc.services.muxWIndex.RLock()
	defer svc.services.muxWIndex.RUnlock()

	return svc.services.wIndex[h]
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
