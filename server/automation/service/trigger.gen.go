package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/label"
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, TriggerActionCreate, err)
}

func (svc *trigger) Update(ctx context.Context, upd *types.Trigger) (res *types.Trigger, err error) {
	var (
		aProps = &triggerActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, TriggerActionUpdate, err)
}

func (svc *trigger) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, TriggerActionDelete, err)
}

func (svc *trigger) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &triggerActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, TriggerActionUndelete, err)
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
