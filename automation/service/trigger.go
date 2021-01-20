package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type (
	trigger struct {
		eventbus  triggerEventTriggerHandler
		store     store.Storer
		actionlog actionlog.Recorder
		ac        triggerAccessController

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
		CanManageWorkflowTriggers(context.Context, *types.Workflow) bool
	}

	triggerEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	triggerUpdateHandler func(ctx context.Context, ns *types.Trigger) (triggerChanges, error)
	triggerChanges       uint8
)

const (
	triggerUnchanged     triggerChanges = 0
	triggerChanged       triggerChanges = 1
	triggerLabelsChanged triggerChanges = 2
)

func Trigger(log *zap.Logger) *trigger {
	return &trigger{
		log:       log,
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

		if !svc.ac.CanManageWorkflowTriggers(ctx, wf) {
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

		if err = svc.registerWorkflow(ctx, wf, res); err != nil {
			return
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
	} else if !svc.ac.CanManageWorkflowTriggers(ctx, wf) {
		return permErr
	} else {
		return nil
	}
}

// registers all triggers on all given workflows
// before registering triggers on a workflow, all workflow triggers are unregistered
func (svc *trigger) registerWorkflows(ctx context.Context, workflows ...*types.Workflow) error {
	// load ALL workflows directly from store
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

	if len(types.TriggerSet(tt).FilterByWorkflowID(wf.ID)) < len(tt) {
		return fmt.Errorf("all triggers must reference the given workflow")
	}

	if wf.RunAs > 0 {
		if runAs, err = DefaultUser.FindByID(ctx, wf.RunAs); err != nil {
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
		err   error
		g     *wfexec.Graph
		wfLog = svc.log.
			WithOptions(zap.AddStacktrace(zap.DPanicLevel)).
			With(zap.Uint64("workflowID", wf.ID))
	)

	if !wf.Enabled {
		wfLog.Debug("skipping disabled workflow")
		return
	}

	if wf.DeletedAt != nil {
		wfLog.Debug("skipping deleted workflow")
		return
	}

	if g, err = svc.workflow.toGraph(wf); err != nil {
		wfLog.Error("failed to convert workflow to graph", zap.Error(err))
		return
	}

	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, t := range tt {
		log := wfLog.With(zap.Uint64("triggerID", t.ID))

		if !t.Enabled {
			log.Debug("skipping disabled trigger")
			continue
		}

		if t.DeletedAt != nil {
			log.Debug("skipping deleted trigger")
			continue
		}

		var (
			handler = func(ctx context.Context, ev eventbus.Event) error {

				var (
					// create session scope from predefined workflow scope and trigger input
					scope = wf.Scope.Merge(t.Input)
					wait  WaitFn
				)

				// scope["event"] = ev

				if runAs == nil {
					// @todo can/should we get alternative identity from Event?
					//       for example:
					//         - use http auth header and get username
					//         - use from/to/replyTo and use that as an identifier
					runAs = auth.GetIdentityFromContext(ctx)
				}

				log.Debug("handling triggered workflow",
					zap.Any("event", ev),
					zap.Uint64("runAs", runAs.Identity()),
				)

				wait, err = svc.session.Start(g, runAs, types.SessionStartParams{
					WorkflowID:   wf.ID,
					KeepFor:      wf.KeepSessions,
					Trace:        wf.Trace,
					Input:        scope,
					StepID:       t.StepID,
					EventType:    t.EventType,
					ResourceType: t.ResourceType,
				})

				if err != nil {
					log.Error("workflow error", zap.Error(err))
					return err
				}

				return wait(ctx)
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

		if svc.reg[wf.ID] == nil {
			svc.reg[wf.ID] = make(map[uint64]uintptr)
		} else if ptr := svc.reg[wf.ID][t.ID]; ptr != 0 {
			// unregister handlers for this trigger if they exist
			svc.eventbus.Unregister(ptr)
		}

		svc.reg[wf.ID][t.ID] = svc.eventbus.Register(handler, ops...)

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

func loadTrigger(ctx context.Context, s store.Storer, workflowID uint64) (res *types.Trigger, err error) {
	if workflowID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if res, err = store.LookupAutomationTriggerByID(ctx, s, workflowID); errors.IsNotFound(err) {
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
