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
	"github.com/crusttech/human/server/store"
)

func (svc *ngAutomation) Search(ctx context.Context, filter types.NgAutomationFilter) (set types.NgAutomationSet, f types.NgAutomationFilter, err error) {
	var (
		aProps = &ngAutomationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.NgAutomation) (bool, error) {
		if !svc.ac.CanReadNgAutomation(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchNgAutomations(ctx) {
			return NgAutomationErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.NgAutomation{}.LabelResourceKind(),
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

		if set, f, err = store.SearchAutomationNgAutomations(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledNgAutomations(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, NgAutomationActionSearch, err)
}

// toLabeledNgAutomations converts to []label.LabeledResource
func toLabeledNgAutomations(set []*types.NgAutomation) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
