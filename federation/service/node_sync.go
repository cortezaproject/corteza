package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	nodeSync struct {
		store     store.Storer
		actionlog actionlog.Recorder
	}

	NodeSyncService interface {
		Create(ctx context.Context, new *types.NodeSync) (*types.NodeSync, error)
		Search(ctx context.Context, f types.NodeSyncFilter) (types.NodeSyncSet, types.NodeSyncFilter, error)
		LookupLastSuccessfulSync(ctx context.Context, nodeID uint64, syncType string) (*types.NodeSync, error)
	}
)

func NodeSync() NodeSyncService {
	return &nodeSync{
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

func (svc nodeSync) Create(ctx context.Context, new *types.NodeSync) (*types.NodeSync, error) {
	var (
		aProps = &nodeSyncActionProps{nodeSync: new}
	)

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if _, err := DefaultNode.FindByID(ctx, new.NodeID); err != nil {
			return NodeSyncErrNodeNotFound()
		}

		return store.CreateFederationNodesSync(ctx, s, new)
	})

	return new, svc.recordAction(ctx, aProps, NodeSyncActionCreate, err)
}

func (svc nodeSync) Search(ctx context.Context, f types.NodeSyncFilter) (types.NodeSyncSet, types.NodeSyncFilter, error) {
	return store.SearchFederationNodesSyncs(ctx, svc.store, f)
}

func (svc nodeSync) LookupLastSuccessfulSync(ctx context.Context, nodeID uint64, syncType string) (ns *types.NodeSync, err error) {
	// todo - filter by synctype does not work
	s, _, err := store.SearchFederationNodesSyncs(ctx, svc.store, types.NodeSyncFilter{
		NodeID:     nodeID,
		SyncType:   syncType,
		SyncStatus: types.NodeSyncStatusSuccess,
		Sorting: filter.Sorting{
			Sort: filter.SortExprSet{
				&filter.SortExpr{Column: "time_action", Descending: true},
			},
		},
		Paging: filter.Paging{Limit: 1},
	})

	if err != nil || len(s) == 0 {
		return nil, err
	}

	return s[0], nil
}
