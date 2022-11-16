package service

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/cortezaproject/corteza/server/automation/types"
	cmpEvent "github.com/cortezaproject/corteza/server/compose/service/event"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/store"
	sysEvent "github.com/cortezaproject/corteza/server/system/service/event"
	"go.uber.org/zap"
)

type (
	trigger struct {
		eventbus  triggerEventTriggerHandler
		store     store.Storer
		actionlog actionlog.Recorder
		ac        triggerAccessController

		opt options.WorkflowOpt

		log *zap.Logger

		// maps registered triggers (value, uintptr) to trigger ID (key, uint64)
		// this will keep track of all our trigger registrations and help us do a cleanup on
		// trigger update.
		triggers map[uint64]uintptr

		reg map[uint64]map[uint64]uintptr

		workflow *workflow
		session  *session

		mux *sync.RWMutex
	}

	triggerAccessController interface {
		CanSearchTriggers(context.Context) bool
		CanManageTriggersOnWorkflow(context.Context, *types.Workflow) bool
		CanExecuteWorkflow(context.Context, *types.Workflow) bool
	}

	triggerEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	triggerUpdateHandler func(ctx context.Context, ns *types.Trigger) (triggerChanges, error)
	triggerChanges       uint8

	varsEncoder interface {
		EncodeVars() (*expr.Vars, error)
	}

	varsDecoder interface {
		DecodeVars(*expr.Vars) error
	}
)

const (
	triggerUnchanged     triggerChanges = 0
	triggerChanged       triggerChanges = 1
	triggerLabelsChanged triggerChanges = 2
)

func Trigger(log *zap.Logger, opt options.WorkflowOpt) *trigger {
	return &trigger{
		log:       log,
		opt:       opt,
		eventbus:  eventbus.Service(),
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		session:   DefaultSession,
		workflow:  DefaultWorkflow,
		triggers:  make(map[uint64]uintptr),
		reg:       make(map[uint64]map[uint64]uintptr),
		mux:       &sync.RWMutex{},
	}
}

func (svc *trigger) Search(ctx context.Context, filter types.TriggerFilter) (rr types.TriggerSet, f types.TriggerFilter, err error) {
	var (
		wap = &triggerActionProps{filter: &filter}
	)

	err = func() (err error) {
		if !svc.ac.CanSearchTriggers(ctx) {
			return TriggerErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Trigger{}.LabelResourceKind(),
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

		if rr, f, err = store.SearchAutomationTriggers(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledTriggers(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, filter, svc.recordAction(ctx, wap, TriggerActionSearch, err)
}

// SearchOnManual finds first matching onManual trigger and returns it
//
// In case stepID is 0, first trigger is returned
func (svc *trigger) SearchOnManual(ctx context.Context, workflowID, stepID uint64) (*types.Trigger, error) {
	tt, _, err := svc.Search(ctx, types.TriggerFilter{
		WorkflowID: []uint64{workflowID},
		EventType:  "onManual",
	})

	if err != nil {
		return nil, err
	}

	if stepID == 0 && len(tt) > 0 {
		return tt[0], nil
	}

	for _, t := range tt {
		if t.StepID == stepID {
			return t, nil
		}
	}

	return nil, nil
}

func (svc *trigger) LookupByID(ctx context.Context, triggerID uint64) (res *types.Trigger, err error) {
	var (
		wap = &triggerActionProps{trigger: &types.Trigger{ID: triggerID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if !svc.ac.CanSearchTriggers(ctx) {
			return TriggerErrNotAllowedToRead()
		}

		if res, err = loadTrigger(ctx, s, triggerID); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, wap, TriggerActionLookup, err)
}

// Create adds new trigger resource and saves it into store
// It updates service's cache
func (svc *trigger) Create(ctx context.Context, new *types.Trigger) (res *types.Trigger, err error) {
	var (
		wap   = &triggerActionProps{new: new}
		cUser = auth.GetIdentityFromContext(ctx).Identity()
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			wf *types.Workflow
		)

		if wf, err = loadWorkflow(ctx, svc.store, new.WorkflowID); err != nil {
			return err
		}

		if !svc.ac.CanManageTriggersOnWorkflow(ctx, wf) {
			return TriggerErrNotAllowedToCreate()
		}

		res = &types.Trigger{
			ID:           nextID(),
			Enabled:      new.Enabled,
			WorkflowID:   new.WorkflowID,
			StepID:       new.StepID,
			ResourceType: new.ResourceType,
			EventType:    new.EventType,
			Constraints:  new.Constraints,
			Input:        new.Input,
			Labels:       new.Labels,
			Meta:         new.Meta,
			OwnedBy:      cUser,
			CreatedAt:    *now(),
			CreatedBy:    cUser,
		}

		if err = store.CreateAutomationTrigger(ctx, s, res); err != nil {
			return
		}

		if err = label.Create(ctx, s, res); err != nil {
			return
		}

		// Ignore workflow issues as those are defined by the workflow itself.
		// Internal errors should still be reported.
		if err = svc.registerWorkflow(ctx, wf, res); err != nil {
			if _, ok := err.(types.WorkflowIssueSet); ok {
				err = nil
			}
		}

		return
	})

	return res, svc.recordAction(ctx, wap, TriggerActionCreate, err)
}

// Update modifies existing trigger resource in the store
func (svc *trigger) Update(ctx context.Context, upd *types.Trigger) (*types.Trigger, error) {
	return svc.updater(ctx, upd.ID, TriggerActionUpdate, func(ctx context.Context, res *types.Trigger) (triggerChanges, error) {
		if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToUpdate()); err != nil {
			return triggerUnchanged, err
		}

		handler := svc.handleUpdate(upd)
		return handler(ctx, res)
	})
}

func (svc *trigger) DeleteByID(ctx context.Context, triggerID uint64) error {
	return trim1st(svc.updater(ctx, triggerID, TriggerActionDelete, svc.handleDelete))
}

func (svc *trigger) UndeleteByID(ctx context.Context, triggerID uint64) error {
	return trim1st(svc.updater(ctx, triggerID, TriggerActionUndelete, svc.handleUndelete))
}

func (svc trigger) updater(ctx context.Context, triggerID uint64, action func(...*triggerActionProps) *triggerAction, fn triggerUpdateHandler) (*types.Trigger, error) {
	var (
		changes triggerChanges
		res     *types.Trigger
		aProps  = &triggerActionProps{trigger: &types.Trigger{ID: triggerID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		res, err = loadTrigger(ctx, s, triggerID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setTrigger(res)
		aProps.setUpdate(res)

		if changes, err = fn(ctx, res); err != nil {
			return err
		}

		if changes&triggerChanged > 0 {
			if err = store.UpdateAutomationTrigger(ctx, svc.store, res); err != nil {
				return err
			}
		}

		if changes&triggerLabelsChanged > 0 {
			if err = label.Update(ctx, s, res); err != nil {
				return
			}
		}

		return err
	})

	return res, svc.recordAction(ctx, aProps, action, err)
}

func (svc trigger) handleUpdate(upd *types.Trigger) triggerUpdateHandler {
	return func(ctx context.Context, res *types.Trigger) (changes triggerChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return triggerUnchanged, TriggerErrStaleData()
		}

		if res.Enabled != upd.Enabled {
			changes |= triggerChanged
			res.Enabled = upd.Enabled
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= triggerLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if res.StepID != upd.StepID {
			changes |= triggerChanged
			res.StepID = upd.StepID
		}

		if res.EventType != upd.EventType {
			changes |= triggerChanged
			res.EventType = upd.EventType
		}

		if res.ResourceType != upd.ResourceType {
			changes |= triggerChanged
			res.ResourceType = upd.ResourceType
		}

		if upd.Meta != nil {
			if !reflect.DeepEqual(upd.Meta, res.Meta) {
				changes |= triggerChanged
				res.Meta = upd.Meta
			}
		}

		if upd.Input != nil {
			if !reflect.DeepEqual(upd.Input, res.Input) {
				changes |= triggerChanged
				res.Input = upd.Input
			}
		}

		if upd.Constraints != nil {
			if !reflect.DeepEqual(upd.Constraints, res.Constraints) {
				changes |= triggerChanged
				res.Constraints = upd.Constraints
			}
		}

		if res.OwnedBy != upd.OwnedBy {
			// @todo need to check against access control if current user can modify owner
			changes |= triggerChanged
			res.OwnedBy = upd.OwnedBy
		}

		if changes&triggerChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc trigger) handleDelete(ctx context.Context, res *types.Trigger) (triggerChanges, error) {
	if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToDelete()); err != nil {
		return triggerUnchanged, err
	}

	if res.DeletedAt != nil {
		// trigger already deleted
		return triggerUnchanged, nil
	}

	res.DeletedAt = now()
	return triggerChanged, nil
}

func (svc trigger) handleUndelete(ctx context.Context, res *types.Trigger) (triggerChanges, error) {
	if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToUndelete()); err != nil {
		return triggerUnchanged, err
	}

	if res.DeletedAt == nil {
		// trigger not deleted
		return triggerUnchanged, nil
	}

	res.DeletedAt = nil
	return triggerChanged, nil
}

func (svc trigger) canManageTrigger(ctx context.Context, res *types.Trigger, permErr error) error {
	if wf, err := loadWorkflow(ctx, svc.store, res.WorkflowID); err != nil {
		return err
	} else if !svc.ac.CanManageTriggersOnWorkflow(ctx, wf) {
		return permErr
	} else {
		return nil
	}
}

// registers all triggers on all given workflows
// before registering triggers on a workflow, all workflow triggers are unregistered
func (svc *trigger) registerWorkflows(ctx context.Context, workflows ...*types.Workflow) error {
	// load ALL triggers directly from store
	tt, _, err := store.SearchAutomationTriggers(ctx, svc.store, types.TriggerFilter{
		WorkflowID: types.WorkflowSet(workflows).IDs(),
		Deleted:    filter.StateInclusive,
		Disabled:   filter.StateExcluded,
	})

	if err != nil {
		return err
	}

	for _, wf := range workflows {
		svc.unregisterWorkflows(wf)

		if !wf.Enabled {
			continue
		}

		if wf.DeletedAt != nil {
			continue
		}

		if len(wf.Issues) > 0 {
			// workflow was processed before and issues were detected
			// and stored on the workflow; no need to continue
			continue
		}

		if err = svc.registerWorkflow(ctx, wf, tt.FilterByWorkflowID(wf.ID)...); err != nil {
			return err
		}
	}

	return nil
}

// updates trigger handler registration
//
// Loads associated workflow and registers specific trigger
func (svc *trigger) updateTriggerRegistration(ctx context.Context, t *types.Trigger) error {
	wf, err := loadWorkflow(ctx, svc.store, t.WorkflowID)
	if err != nil {
		return err
	}

	return svc.registerWorkflow(ctx, wf, t)
}

// registers one workflow and a set of triggers
func (svc *trigger) registerWorkflow(ctx context.Context, wf *types.Workflow, tt ...*types.Trigger) (err error) {
	var (
		runAs auth.Identifiable
	)

	// Returns context with identity set to service user
	//
	// Current user (identity in the context) might not have
	// sufficient privileges to load info about invoker and runner
	sysUserCtx := func() context.Context {
		return auth.SetIdentityToContext(ctx, auth.ServiceUser())
	}

	if !svc.opt.Register {
		return nil
	}

	if !wf.Enabled || len(wf.Issues) > 0 {
		// do not even try to register disabled
		// workflows or workflows with issues
		return nil
	}

	if len(types.TriggerSet(tt).FilterByWorkflowID(wf.ID)) < len(tt) {
		return fmt.Errorf("all triggers must reference the given workflow")
	}

	if wis := validateWorkflowTriggers(wf, tt...); len(wis) > 0 {
		// skip trigger registration of there is a trigger related issue(s)
		// on a specific workflow.
		//
		// this really happens since we run all validation on save,
		// but there might be workflow from the time when these checks were
		// not in place
		return nil
	}

	if wf.RunAs > 0 {
		if runAs, err = DefaultUser.FindByAny(sysUserCtx(), wf.RunAs); err != nil {
			return fmt.Errorf("failed to load run-as user %d: %w", wf.RunAs, err)
		} else if !runAs.Valid() {
			return fmt.Errorf("invalid user %d used for workflow run-as", wf.RunAs)
		}
	}

	svc.registerTriggers(wf, runAs, tt...)
	return nil
}

// registerTriggers registers workflows triggers to eventbus
//
// It preloads run-as identity and finds a starting step for each trigger
func (svc *trigger) registerTriggers(wf *types.Workflow, runAs auth.Identifiable, tt ...*types.Trigger) {
	var (
		handlerFn eventbus.HandlerFn
		err       error
		g         *wfexec.Graph
		issues    types.WorkflowIssueSet
		wfLog     = svc.log.
				With(zap.Uint64("workflowID", wf.ID))

		// register only enabled, undeleted workflows
		registerWorkflow = wf.Enabled && wf.DeletedAt == nil
	)

	// convert only registrable and workflows without issues
	if registerWorkflow && len(wf.Issues) == 0 {
		// Convert workflow only when valid (no issues, enable, not delete)
		if g, issues = Convert(svc.workflow, wf); len(issues) > 0 {
			wfLog.Error("failed to convert workflow to graph", zap.Error(issues))
			_ = issues.Walk(func(i *types.WorkflowIssue) error {
				wfLog.Debug("workflow issue found: "+i.Description, zap.Any("culprit", i.Culprit))
				return nil
			})
			g = nil
		}
	}

	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, t := range tt {
		log := wfLog.With(zap.Uint64("triggerID", t.ID))

		// always unregister
		if svc.reg[wf.ID] == nil {
			svc.reg[wf.ID] = make(map[uint64]uintptr)
		} else if ptr := svc.reg[wf.ID][t.ID]; ptr != 0 {
			// unregister handlers for this trigger if they exist
			svc.eventbus.Unregister(ptr)
		}

		// do not register disabled or deleted triggers
		if !registerWorkflow || !t.Enabled || t.DeletedAt != nil {
			continue
		}

		var (
			cnstr eventbus.ConstraintMatcher
			ops   = make([]eventbus.HandlerRegOp, 0, len(t.Constraints)+2)
		)

		if g == nil {
			handlerFn = func(_ context.Context, ev eventbus.Event) error {
				return errors.Internal(
					"trigger %s on %s failed due to invalid workflow %d: %s",
					ev.EventType(),
					ev.ResourceType(),
					wf.ID,
					wf.Issues,
				).Wrap(wf.Issues)
			}
		} else {
			handlerFn = makeWorkflowHandler(svc.workflow, wf, t)
		}

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

		svc.reg[wf.ID][t.ID] = svc.eventbus.Register(handlerFn, ops...)

		log.Debug("trigger registered",
			zap.String("eventType", t.EventType),
			zap.String("resourceType", t.ResourceType),
			zap.Any("constraints", t.Constraints),
		)
	}
}

func (svc *trigger) unregisterWorkflows(wwf ...*types.Workflow) {
	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, wf := range wwf {
		for triggerID, ptr := range svc.reg[wf.ID] {
			svc.eventbus.Unregister(ptr)
			svc.log.Debug("trigger unregistered", zap.Uint64("triggerID", triggerID), zap.Uint64("workflowID", wf.ID))
			delete(svc.triggers, wf.ID)
		}

		delete(svc.reg, wf.ID)
	}
}

func (svc *trigger) unregisterTriggers(tt ...*types.Trigger) {
	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, t := range tt {
		if svc.reg[t.WorkflowID] == nil {
			return
		}

		if ptr, has := svc.reg[t.WorkflowID][t.ID]; has {
			svc.eventbus.Unregister(ptr)
			svc.log.Debug("trigger unregistered", zap.Uint64("triggerID", t.ID), zap.Uint64("workflowID", t.WorkflowID))
			delete(svc.triggers, t.ID)
		}
	}
}

func loadTrigger(ctx context.Context, s store.Storer, triggerID uint64) (res *types.Trigger, err error) {
	if triggerID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if res, err = store.LookupAutomationTriggerByID(ctx, s, triggerID); errors.IsNotFound(err) {
		return nil, TriggerErrNotFound()
	}

	return
}

func loadWorkflowTriggers(ctx context.Context, s store.Storer, workflowID uint64) (tt types.TriggerSet, err error) {
	if workflowID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{WorkflowID: []uint64{workflowID}}); errors.IsNotFound(err) {
		return nil, TriggerErrNotFound()
	}

	return
}

// toLabeledTriggers converts to []label.LabeledResource
func toLabeledTriggers(set []*types.Trigger) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

// Checks if triggers are compatible with the workflow
//
// It ignores disabled triggers and does not care if triggers are in fact bond to the
// given workflow
func validateWorkflowTriggers(wf *types.Workflow, tt ...*types.Trigger) (wis types.WorkflowIssueSet) {
	var (
		// @todo find a better way how to flag events type that require
		//       run-as param to be set
		//       Possible solution: flag in definition that generates static
		//       list w/o the need  of cross-component imports
		requireRunAs = []eventbus.Event{
			sysEvent.SinkOnRequest(nil, nil),
			sysEvent.QueueOnMessage(nil),
			sysEvent.SystemOnInterval(),
			sysEvent.SystemOnTimestamp(),
			cmpEvent.ComposeOnInterval(),
			cmpEvent.ComposeOnTimestamp(),
		}
	)

	for i, t := range tt {
		if !t.Enabled || t.DeletedAt != nil {
			continue
		}

		for _, ev := range requireRunAs {
			if t.ResourceType == ev.ResourceType() && t.EventType == ev.EventType() {
				if wf.RunAs == 0 {
					wis = wis.Append(
						errors.InvalidData("%s for %s requires run-as to be set", ev.ResourceType(), ev.EventType()),
						map[string]int{"trigger": i},
					)
				}
			}
		}
	}

	return
}
