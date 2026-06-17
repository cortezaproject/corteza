package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *configuredConnection) Create(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &configuredConnectionActionProps{connection: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateConfiguredConnection(ctx) {
			return ConfiguredConnectionErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateConfiguredConnection(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionCreate, err)
}

// toLabeledConfiguredConnections converts to []label.LabeledResource
func toLabeledConfiguredConnections(set []*types.ConfiguredConnection) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
