package service

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/crusttech/human/server/automation/types"
	cmpEvent "github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/wfexec"
	"github.com/crusttech/human/server/store"
	sysEvent "github.com/crusttech/human/server/system/service/event"
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

	varsEncoder interface {
		EncodeVars() (*expr.Vars, error)
	}

	varsDecoder interface {
		DecodeVars(*expr.Vars) error
	}
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

// onSearch is the generated Search body handler.
//
// The recordAction wrapper, aProps and the standard CanSearchTriggers access
// check are owned by the generated trigger.gen.go. This returns the (possibly
// label-augmented) input filter as f, preserving the original behaviour.
func (svc *trigger) onSearch(ctx context.Context, filter types.TriggerFilter, wap *triggerActionProps) (rr types.TriggerSet, f types.TriggerFilter, err error) {
	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Trigger{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return rr, filter, err
		}

		// labels specified but no labeled resources found
		if len(filter.LabeledIDs) == 0 {
			return rr, filter, nil
		}
	}

	if rr, _, err = store.SearchAutomationTriggers(ctx, svc.store, filter); err != nil {
		return rr, filter, err
	}

	if err = label.Load(ctx, svc.store, toLabeledTriggers(rr)...); err != nil {
		return rr, filter, err
	}

	return rr, filter, nil
}

// onSearchOnManual finds first matching onManual trigger and returns it
//
// In case stepID is 0, first trigger is returned
func (svc *trigger) onSearchOnManual(ctx context.Context, _ *triggerActionProps, workflowID, stepID uint64) (res *types.Trigger, err error) {
	tt, _, err := svc.Search(ctx, types.TriggerFilter{
		WorkflowID: id.Strings(workflowID),
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

func (svc *trigger) FindByID(ctx context.Context, triggerID uint64) (res *types.Trigger, err error) {
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

// onCreate is the generated Create body handler.
//
// The recordAction wrapper, aProps and the res=new assignment are owned by the
// generated trigger.gen.go. Access runs against the loaded workflow
// (CanManageTriggersOnWorkflow), so the generated scaffold emits no standard
// access check (create is in customAccessOps).
//
// new is mutated in place so the generated res=new carries the created trigger.
// It adds a new trigger resource, saves it into the store and updates the
// service's registrations.
func (svc *trigger) onCreate(ctx context.Context, new *types.Trigger) (err error) {
	var (
		cUser = auth.GetIdentityFromContext(ctx).Identity()
	)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		var (
			wf *types.Workflow
		)

		if wf, err = loadWorkflow(ctx, svc.store, new.WorkflowID); err != nil {
			return err
		}

		if !svc.ac.CanManageTriggersOnWorkflow(ctx, wf) {
			return TriggerErrNotAllowedToCreate()
		}

		new.ID = nextID()
		new.OwnedBy = cUser
		new.CreatedAt = *now()
		new.CreatedBy = cUser
		new.UpdatedAt = nil
		new.UpdatedBy = 0
		new.DeletedAt = nil
		new.DeletedBy = 0

		if err = store.CreateAutomationTrigger(ctx, s, new); err != nil {
			return
		}

		if err = label.Create(ctx, s, new); err != nil {
			return
		}

		// Ignore workflow issues as those are defined by the workflow itself.
		// Internal errors should still be reported.
		if err = svc.registerWorkflow(ctx, wf, new); err != nil {
			if _, ok := err.(types.WorkflowIssueSet); ok {
				err = nil
			}
		}

		return
	})
}

// onUpdate receives a pre-loaded res and tx store from the generated Update
// scaffold. It verifies access, applies field changes, and persists to store.
func (svc *trigger) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Trigger, aProps *triggerActionProps, _, _ func() error) error {
	if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToUpdate()); err != nil {
		return err
	}

	changed := false

	if res.Enabled != upd.Enabled {
		changed = true
		res.Enabled = upd.Enabled
	}
	if res.StepID != upd.StepID {
		changed = true
		res.StepID = upd.StepID
	}
	if res.EventType != upd.EventType {
		changed = true
		res.EventType = upd.EventType
	}
	if res.ResourceType != upd.ResourceType {
		changed = true
		res.ResourceType = upd.ResourceType
	}
	if upd.Meta != nil && !reflect.DeepEqual(upd.Meta, res.Meta) {
		changed = true
		res.Meta = upd.Meta
	}
	if upd.Input != nil && !reflect.DeepEqual(upd.Input, res.Input) {
		changed = true
		res.Input = upd.Input
	}
	if upd.Constraints != nil && !reflect.DeepEqual(upd.Constraints, res.Constraints) {
		changed = true
		res.Constraints = upd.Constraints
	}
	if res.OwnedBy != upd.OwnedBy {
		changed = true
		res.OwnedBy = upd.OwnedBy
	}

	if changed {
		res.UpdatedAt = now()
		if err := store.UpdateAutomationTrigger(ctx, s, res); err != nil {
			return err
		}
	}

	if upd.Labels != nil && label.Changed(res.Labels, upd.Labels) {
		res.Labels = upd.Labels
		if err := label.Update(ctx, s, res); err != nil {
			return err
		}
	}

	return nil
}

// onDelete receives a pre-loaded res and tx store from the generated DeleteByID scaffold.
func (svc *trigger) onDelete(ctx context.Context, s store.Storer, res *types.Trigger, aProps *triggerActionProps) error {
	if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToDelete()); err != nil {
		return err
	}

	if res.DeletedAt != nil {
		return nil
	}

	res.DeletedAt = now()
	return store.UpdateAutomationTrigger(ctx, s, res)
}

// onUndelete receives a pre-loaded res and tx store from the generated UndeleteByID scaffold.
func (svc *trigger) onUndelete(ctx context.Context, s store.Storer, res *types.Trigger, aProps *triggerActionProps) error {
	if err := svc.canManageTrigger(ctx, res, TriggerErrNotAllowedToUndelete()); err != nil {
		return err
	}

	if res.DeletedAt == nil {
		return nil
	}

	res.DeletedAt = nil
	return store.UpdateAutomationTrigger(ctx, s, res)
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
		WorkflowID: id.Strings(types.WorkflowSet(workflows).IDs()...),
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
				With(logger.Uint64("workflowID", wf.ID))

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
		log := wfLog.With(logger.Uint64("triggerID", t.ID))

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
			svc.log.Debug("trigger unregistered", logger.Uint64("triggerID", triggerID), logger.Uint64("workflowID", wf.ID))
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
			svc.log.Debug("trigger unregistered", logger.Uint64("triggerID", t.ID), logger.Uint64("workflowID", t.WorkflowID))
			delete(svc.triggers, t.ID)
		}
	}
}

func loadWorkflowTriggers(ctx context.Context, s store.Storer, workflowID uint64) (tt types.TriggerSet, err error) {
	if workflowID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if tt, _, err = store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{WorkflowID: id.Strings(workflowID)}); errors.IsNotFound(err) {
		return nil, TriggerErrNotFound()
	}

	return
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

		if wf.Meta != nil && wf.Meta.SubWorkflow {
			wis = wis.Append(
				errors.InvalidData("workflow marked as sub-workflow cannot have enabled triggers"),
				map[string]int{"trigger": i},
			)
		}
	}

	return
}
