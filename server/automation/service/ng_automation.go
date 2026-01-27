package service

import (
	"context"
	"reflect"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/rbac"

	"github.com/cortezaproject/corteza/server/store"
	"go.uber.org/zap"
)

type (
	ngAutomation struct {
		eventbus  ngAutomationEventTriggerHandler
		store     store.Storer
		actionlog actionlog.Recorder
		ac        ngAutomationAccessController

		log *zap.Logger
	}

	ngAutomationAccessController interface {
		CanCreateNgAutomation(context.Context) bool
		CanSearchNgAutomations(context.Context) bool
		CanReadNgAutomation(context.Context, *types.NgAutomation) bool
		CanUpdateNgAutomation(context.Context, *types.NgAutomation) bool
		CanDeleteNgAutomation(context.Context, *types.NgAutomation) bool
		CanUndeleteNgAutomation(context.Context, *types.NgAutomation) bool

		Grant(ctx context.Context, rr ...*rbac.Rule) error

		ngAutomationExecController
	}

	ngAutomationExecController interface {
		// CanExecuteNgAutomation(context.Context, *types.NgAutomation) bool
	}

	ngAutomationEventTriggerHandler interface {
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	ngAutomationUpdateHandler func(ctx context.Context, ns *types.NgAutomation) (ngAutomationChanges, error)
	ngAutomationChanges       uint8

	ngAutomationInvokerCtxKey struct{}
)

const (
	ngAutomationUnchanged     ngAutomationChanges = 0
	ngAutomationChanged       ngAutomationChanges = 1
	ngAutomationLabelsChanged ngAutomationChanges = 2
	ngAutomationDefChanged    ngAutomationChanges = 4
)

func NgAutomation(log *zap.Logger, corredorOpt options.CorredorOpt) *ngAutomation {
	return &ngAutomation{
		log: log,
		// opt:         opt,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		eventbus:  eventbus.Service(),
	}
}

func (svc *ngAutomation) Search(ctx context.Context, filter types.NgAutomationFilter) (rr types.NgAutomationSet, f types.NgAutomationFilter, err error) {
	var (
		wap = &ngAutomationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.NgAutomation) (bool, error) {
		if !svc.ac.CanReadNgAutomation(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() (err error) {
		if !svc.ac.CanSearchNgAutomations(ctx) {
			return NgAutomationErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				"ng-automation",
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

		if rr, f, err = store.SearchAutomationNgAutomations(ctx, svc.store, filter); err != nil {
			return err
		}

		// if err = label.Load(ctx, svc.store, toLabeledNgAutomations(rr)...); err != nil {
		// 	return err
		// }

		return nil
	}()

	return rr, f, svc.recordAction(ctx, wap, NgAutomationActionSearch, err)
}

func (svc *ngAutomation) LookupByID(ctx context.Context, ngAutomationID uint64) (ngAtuomation *types.NgAutomation, err error) {
	var (
		wap = &ngAutomationActionProps{ngAutomation: &types.NgAutomation{ID: ngAutomationID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if ngAtuomation, err = loadNgAutomation(ctx, s, ngAutomationID); err != nil {
			return err
		}

		if !svc.ac.CanReadNgAutomation(ctx, ngAtuomation) {
			return NgAutomationErrNotAllowedToRead()
		}

		// if err = label.Load(ctx, svc.store, ngAtuomation); err != nil {
		// 	return err
		// }

		return nil
	})

	return ngAtuomation, svc.recordAction(ctx, wap, NgAutomationActionLookup, err)
}

// Create adds new ngAutomation resource and saves it into store
// It updates service's cache
func (svc *ngAutomation) Create(ctx context.Context, new *types.NgAutomation) (ngAtuomation *types.NgAutomation, err error) {
	var (
		wap   = &ngAutomationActionProps{ngAutomation: new}
		cUser = intAuth.GetIdentityFromContext(ctx).Identity()
		// g     *ngAtuomationexec.Graph
		// runAs intAuth.Identifiable
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanCreateNgAutomation(ctx) {
			return NgAutomationErrNotAllowedToCreate()
		}

		if new.Meta.Name == "" {
			return NgAutomationErrMissingName()
		}

		if !handle.IsValid(new.Handle) {
			return NgAutomationErrInvalidHandle()
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		ngAtuomation = &types.NgAutomation{
			ID:      nextID(),
			Handle:  new.Handle,
			Labels:  new.Labels,
			Meta:    new.Meta,
			Enabled: new.Enabled,

			Scope:    new.Scope,
			Triggers: new.Triggers,
			Steps:    new.Steps,
			Paths:    new.Paths,

			// @todo need to check against access control if current user can modify security descriptor
			RunAs:     new.RunAs,
			OwnedBy:   cUser,
			CreatedAt: *now(),
			CreatedBy: cUser,
		}

		wap.ngAutomation = ngAtuomation

		// if g, runAs, err = svc.validateNgAutomation(ctx, ngAtuomation); err != nil {
		// 	return
		// }

		// svc.updateCache(ngAtuomation, runAs, g)

		// if len(ngAtuomation.Issues) == 0 {
		// 	if err = svc.triggers.registerNgAutomations(ctx, ngAtuomation); err != nil {
		// 		return err
		// 	}
		// }

		if err = store.CreateAutomationNgAutomation(ctx, s, ngAtuomation); err != nil {
			return
		}

		// if err = label.Create(ctx, s, ngAtuomation); err != nil {
		// 	return
		// }

		wap.setNew(ngAtuomation)

		return
	})

	return ngAtuomation, svc.recordAction(ctx, wap, NgAutomationActionCreate, err)
}

// Update modifies existing ngAutomation resource in the store
func (svc *ngAutomation) Update(ctx context.Context, upd *types.NgAutomation) (*types.NgAutomation, error) {
	return svc.updater(ctx, upd.ID, NgAutomationActionUpdate, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		if upd.Meta.Name == "" {
			return ngAutomationUnchanged, NgAutomationErrMissingName()
		}

		if !svc.ac.CanUpdateNgAutomation(ctx, res) {
			return ngAutomationUnchanged, NgAutomationErrNotAllowedToUpdate()
		}

		handler := svc.handleUpdate(upd)
		return handler(ctx, res)
	})
}

func (svc *ngAutomation) DeleteByID(ctx context.Context, ngAutomationID uint64) error {
	return trim1st(svc.updater(ctx, ngAutomationID, NgAutomationActionDelete, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		changes, err := svc.handleDelete(ctx, res)
		if err != nil {
			return ngAutomationUnchanged, err
		}

		return changes, err
	}))
}

func (svc *ngAutomation) UndeleteByID(ctx context.Context, ngAutomationID uint64) error {
	return trim1st(svc.updater(ctx, ngAutomationID, NgAutomationActionUndelete, func(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
		var (
			changes, err = svc.handleUndelete(ctx, res)
		)

		if err != nil {
			return ngAutomationUnchanged, err
		}

		return changes, err

	}))
}

func (svc ngAutomation) uniqueCheck(ctx context.Context, res *types.NgAutomation) (err error) {
	if res.Handle != "" {
		if e, _ := store.LookupAutomationNgAutomationByHandle(ctx, svc.store, res.Handle); e != nil && e.ID != res.ID {
			return NgAutomationErrHandleNotUnique()
		}
	}

	return nil
}

func (svc *ngAutomation) updater(ctx context.Context, ngAutomationID uint64, action func(...*ngAutomationActionProps) *ngAutomationAction, fn ngAutomationUpdateHandler) (*types.NgAutomation, error) {
	var (
		changes ngAutomationChanges
		res     *types.NgAutomation
		aProps  = &ngAutomationActionProps{ngAutomation: &types.NgAutomation{ID: ngAutomationID}}
		err     error
		// g       *ngAtuomationexec.Graph
		// runAs   intAuth.Identifiable
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {

		res, err = loadNgAutomation(ctx, s, ngAutomationID)
		if err != nil {
			return
		}

		// if err = label.Load(ctx, svc.store, res); err != nil {
		// 	return err
		// }

		aProps.setNgAutomation(res)
		aProps.setUpdate(res)

		if changes, err = fn(ctx, res); err != nil {
			return err
		}

		// if g, runAs, err = svc.validateNgAutomation(ctx, res); err != nil {
		// 	return
		// }

		// svc.updateCache(res, runAs, g)

		// if len(res.Issues) == 0 {
		// 	if err = svc.triggers.registerNgAutomations(ctx, res); err != nil {
		// 		return err
		// 	}
		// }

		if changes&ngAutomationChanged > 0 {
			if err = store.UpdateAutomationNgAutomation(ctx, svc.store, res); err != nil {
				return err
			}
		}

		// if changes&ngAutomationLabelsChanged > 0 {
		// 	if err = label.Update(ctx, s, res); err != nil {
		// 		return
		// 	}
		// }

		return
	})

	return res, svc.recordAction(ctx, aProps, action, err)
}

func (svc ngAutomation) handleUpdate(upd *types.NgAutomation) ngAutomationUpdateHandler {
	return func(ctx context.Context, res *types.NgAutomation) (changes ngAutomationChanges, err error) {
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ngAutomationUnchanged, NgAutomationErrStaleData()
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ngAutomationUnchanged, NgAutomationErrInvalidHandle()
		}

		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return ngAutomationUnchanged, err
		}

		if !svc.ac.CanUpdateNgAutomation(ctx, res) {
			return ngAutomationUnchanged, NgAutomationErrNotAllowedToUpdate()
		}

		if res.Handle != upd.Handle {
			changes |= ngAutomationChanged
			res.Handle = upd.Handle
		}

		if res.Enabled != upd.Enabled {
			changes |= ngAutomationChanged | ngAutomationDefChanged
			res.Enabled = upd.Enabled
		}

		if upd.Labels != nil {
			if label.Changed(res.Labels, upd.Labels) {
				changes |= ngAutomationLabelsChanged
				res.Labels = upd.Labels
			}
		}

		if upd.Meta != nil {
			if !reflect.DeepEqual(upd.Meta, res.Meta) {
				changes |= ngAutomationChanged
				res.Meta = upd.Meta
			}
		}

		if upd.Scope != nil {
			if !reflect.DeepEqual(upd.Scope, res.Scope) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Scope = upd.Scope
			}
		}

		if upd.Triggers != nil {
			if !reflect.DeepEqual(upd.Triggers, res.Triggers) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Triggers = upd.Triggers
			}
		}

		if upd.Steps != nil {
			if !reflect.DeepEqual(upd.Steps, res.Steps) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Steps = upd.Steps
			}
		}

		if upd.Paths != nil {
			if !reflect.DeepEqual(upd.Paths, res.Paths) {
				changes |= ngAutomationChanged | ngAutomationDefChanged
				res.Paths = upd.Paths
			}
		}

		if res.RunAs != upd.RunAs {
			// @todo need to check against access control if current user can modify security descriptor
			changes |= ngAutomationChanged | ngAutomationDefChanged
			res.RunAs = upd.RunAs
		}

		if res.OwnedBy != upd.OwnedBy {
			// @todo need to check against access control if current user can modify owner
			changes |= ngAutomationChanged
			res.OwnedBy = upd.OwnedBy
		}

		if changes&ngAutomationChanged > 0 {
			res.UpdatedAt = now()
		}

		return
	}
}

func (svc ngAutomation) handleDelete(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
	if !svc.ac.CanDeleteNgAutomation(ctx, res) {
		return ngAutomationUnchanged, NgAutomationErrNotAllowedToDelete()
	}

	if res.DeletedAt != nil {
		// ngAutomation already deleted
		return ngAutomationUnchanged, nil
	}

	res.DeletedAt = now()
	return ngAutomationChanged, nil
}

func (svc ngAutomation) handleUndelete(ctx context.Context, res *types.NgAutomation) (ngAutomationChanges, error) {
	if !svc.ac.CanUndeleteNgAutomation(ctx, res) {
		return ngAutomationUnchanged, NgAutomationErrNotAllowedToUndelete()
	}

	if res.DeletedAt == nil {
		// ngAutomation not deleted
		return ngAutomationUnchanged, nil
	}

	res.DeletedAt = nil
	return ngAutomationChanged, nil
}

func loadNgAutomation(ctx context.Context, s store.Storer, ngAutomationID uint64) (res *types.NgAutomation, err error) {
	if ngAutomationID == 0 {
		return nil, NgAutomationErrInvalidID()
	}

	// @todo :)
	if res, err = store.LookupAutomationNgAutomationByID(ctx, s, ngAutomationID); errors.IsNotFound(err) {
		return nil, NgAutomationErrNotFound()
	}

	return
}
