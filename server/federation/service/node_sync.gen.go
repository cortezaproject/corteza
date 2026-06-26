package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/federation/types"
)

func (svc *nodeSync) Search(ctx context.Context, filter types.NodeSyncFilter) (set types.NodeSyncSet, f types.NodeSyncFilter, err error) {
	var (
		aProps = &nodeSyncActionProps{nodeSyncFilter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, NodeSyncActionSearch, err)
}

func (svc *nodeSync) Create(ctx context.Context, new *types.NodeSync) (res *types.NodeSync, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &nodeSyncActionProps{nodeSync: new}
	)

	err = func() (err error) {
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, NodeSyncActionCreate, err)
}
