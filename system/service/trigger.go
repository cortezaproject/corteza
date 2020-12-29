package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
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

		sessions map[uint64]*wfexec.Session

		// maps resolved trigger graphs to trigger ID (key, uint64)
		wfgs map[uint64]*wfexec.Graph

		mux *sync.RWMutex
	}

	triggerAccessController interface {
		CanAccess(context.Context) bool

		CanCreateTrigger(context.Context) bool
		CanReadTrigger(context.Context, *types.Trigger) bool
		CanUpdateTrigger(context.Context, *types.Trigger) bool
		CanDeleteTrigger(context.Context, *types.Trigger) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error
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

func Trigger(log *zap.Logger, s store.Storer, ar actionlog.Recorder, eb triggerEventTriggerHandler) *trigger {
	return &trigger{
		eventbus:  eb,
		actionlog: ar,
		store:     s,
		ac:        DefaultAccessControl,
		log:       log.Named("trigger"),
		wfgs:      make(map[uint64]*wfexec.Graph),
		triggers:  make(map[uint64]uintptr),
		mux:       &sync.RWMutex{},
	}
}

func (svc *trigger) Find(ctx context.Context, filter types.TriggerFilter) (rr types.TriggerSet, f types.TriggerFilter, err error) {
	var (
		wap = &triggerActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Trigger) (bool, error) {
		if !svc.ac.CanReadTrigger(ctx, res) {
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
			// some kind of control over who can see deleted or archived triggers
			if !svc.ac.CanAccess(ctx) {
				return TriggerErrNotAllowedToList()
			}
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

		if rr, f, err = store.SearchTriggers(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledTriggers(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, filter, svc.recordAction(ctx, wap, TriggerActionSearch, err)
}

func (svc *trigger) FindByID(ctx context.Context, triggerID uint64) (wf *types.Trigger, err error) {
	var (
		wap = &triggerActionProps{trigger: &types.Trigger{ID: triggerID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if wf, err = loadTrigger(ctx, s, triggerID); err != nil {
			return err
		}

		if !svc.ac.CanReadTrigger(ctx, wf) {
			return TriggerErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, wf); err != nil {
			return err
		}

		return nil
	})

	return wf, svc.recordAction(ctx, wap, TriggerActionLookup, err)
}

// Create adds new trigger resource and saves it into store
// It updates service's cache
func (svc *trigger) Create(ctx context.Context, new *types.Trigger) (wf *types.Trigger, err error) {
	var (
		wap   = &triggerActionProps{new: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanCreateTrigger(ctx) {
			return TriggerErrNotAllowedToCreate()
		}

		wf = &types.Trigger{
			ID:        nextID(),
			Labels:    new.Labels,
			Enabled:   new.Enabled,
			OwnedBy:   cUser,
			CreatedAt: *now(),
			CreatedBy: cUser,
		}

		if err = store.CreateTrigger(ctx, s, wf); err != nil {
			return
		}

		if err = label.Create(ctx, s, wf); err != nil {
			return
		}

		return
	})

	return wf, svc.recordAction(ctx, wap, TriggerActionCreate, err)
}

// Update modifies existing trigger resource in the store
func (svc *trigger) Update(ctx context.Context, upd *types.Trigger) (wf *types.Trigger, err error) {
	return svc.updater(ctx, upd.ID, TriggerActionUpdate, svc.handleUpdate(upd))
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
		wf      *types.Trigger
		aProps  = &triggerActionProps{trigger: &types.Trigger{ID: triggerID}}
		err     error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		wf, err = loadTrigger(ctx, s, triggerID)
		if err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, wf); err != nil {
			return err
		}

		aProps.setTrigger(wf)
		aProps.setUpdate(wf)

		if changes, err = fn(ctx, wf); err != nil {
			return err
		}

		if changes&triggerChanged > 0 {
			if err = store.UpdateTrigger(ctx, svc.store, wf); err != nil {
				return err
			}
		}

		if changes&triggerLabelsChanged > 0 {
			if err = label.Update(ctx, s, wf); err != nil {
				return
			}
		}

		return err
	})

	return wf, svc.recordAction(ctx, aProps, action, err)
}

func (svc trigger) handleUpdate(upd *types.Trigger) triggerUpdateHandler {
	return func(ctx context.Context, res *types.Trigger) (changes triggerChanges, err error) {
		if !svc.ac.CanUpdateTrigger(ctx, res) {
			return triggerUnchanged, TriggerErrNotAllowedToUpdate()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return triggerUnchanged, TriggerErrStaleData()
		}

		if !svc.ac.CanUpdateTrigger(ctx, res) {
			return triggerUnchanged, TriggerErrNotAllowedToUpdate()
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

		if res.WorkflowID != upd.WorkflowID {
			changes |= triggerChanged
			res.WorkflowID = upd.WorkflowID
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

func (svc trigger) handleDelete(ctx context.Context, wf *types.Trigger) (triggerChanges, error) {
	if !svc.ac.CanDeleteTrigger(ctx, wf) {
		return triggerUnchanged, TriggerErrNotAllowedToDelete()
	}

	if wf.DeletedAt != nil {
		// trigger already deleted
		return triggerUnchanged, nil
	}

	wf.DeletedAt = now()
	return triggerChanged, nil
}

func (svc trigger) handleUndelete(ctx context.Context, wf *types.Trigger) (triggerChanges, error) {
	if !svc.ac.CanDeleteTrigger(ctx, wf) {
		return triggerUnchanged, TriggerErrNotAllowedToUndelete()
	}

	if wf.DeletedAt == nil {
		// trigger not deleted
		return triggerUnchanged, nil
	}

	wf.DeletedAt = nil
	return triggerChanged, nil
}

// registerTriggers registeres workflows triggers to eventbus
//
// It preloads run-as identity and finds a starting step for each trigger
func (svc *trigger) registerTriggers(ctx context.Context, wf *types.Workflow, tt ...*types.Trigger) {
	svc.unregisterTriggers(tt...)

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

	for _, t := range tt {
		log := wfLog.With(
			zap.Uint64("triggerID", t.ID),
			zap.Uint64("workflowID", wf.ID),
		)

		if !t.Enabled {
			log.Debug("skipping disabled trigger")
			continue
		}

		var (
			start wfexec.Step
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

				session := wfexec.NewSession(ctx, g)
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

func (svc *trigger) unregisterTriggers(tt ...*types.Trigger) {
	defer svc.mux.Unlock()
	svc.mux.Lock()

	for _, t := range tt {
		if ptr, has := svc.triggers[t.ID]; has {
			svc.eventbus.Unregister(ptr)
			svc.log.Debug("trigger unregistered", zap.Uint64("triggerID", t.ID), zap.Uint64("workflowID", t.WorkflowID))
		}
	}
}

func loadTrigger(ctx context.Context, s store.Storer, workflowID uint64) (wf *types.Trigger, err error) {
	if workflowID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if wf, err = store.LookupTriggerByID(ctx, s, workflowID); errors.IsNotFound(err) {
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
