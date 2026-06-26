package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
)

func (svc *trigger) Search(ctx context.Context, filter types.TriggerFilter) (set types.TriggerSet, f types.TriggerFilter, err error) {
	var (
		aProps = &triggerActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchTriggers(ctx) {
			return TriggerErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, TriggerActionSearch, err)
}

func (svc *trigger) Create(ctx context.Context, new *types.Trigger) (res *types.Trigger, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &triggerActionProps{trigger: new, new: new}
	)

	err = func() (err error) {
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, TriggerActionCreate, err)
}

func (svc *trigger) Update(ctx context.Context, upd *types.Trigger) (res *types.Trigger, err error) {
	var (
		aProps = &triggerActionProps{update: upd}
		old    *types.Trigger
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadTrigger(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setTrigger(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return TriggerErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.WorkflowID = upd.WorkflowID
		res.Enabled = upd.Enabled
		res.ResourceType = upd.ResourceType
		res.EventType = upd.EventType
		res.OwnedBy = upd.OwnedBy
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateAutomationTrigger(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, TriggerActionUpdate, err, old, res)
}

func (svc *trigger) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerActionProps{}
		res    *types.Trigger
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadTrigger(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setTrigger(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, TriggerActionDelete, err)
}

func (svc *trigger) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerActionProps{}
		res    *types.Trigger
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadTrigger(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setTrigger(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, TriggerActionUndelete, err)
}

func loadTrigger(ctx context.Context, s store.AutomationTriggers, ID uint64) (res *types.Trigger, err error) {
	if ID == 0 {
		return nil, TriggerErrInvalidID()
	}

	if res, err = store.LookupAutomationTriggerByID(ctx, s, ID); errors.IsNotFound(err) {
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
