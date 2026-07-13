package service

import (
	"context"

	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
)

type (
	nodeSyncAccessController interface{}

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

func (svc *nodeSync) onSearch(ctx context.Context, f types.NodeSyncFilter, aProps *nodeSyncActionProps) (types.NodeSyncSet, types.NodeSyncFilter, error) {
	return store.SearchFederationNodeSyncs(ctx, svc.store, f)
}

func (svc *nodeSync) onCreate(ctx context.Context, new *types.NodeSync) error {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if _, err := DefaultNode.FindByID(ctx, new.NodeID); err != nil {
			return NodeSyncErrNodeNotFound()
		}
		return store.CreateFederationNodeSync(ctx, s, new)
	})
}

func (svc nodeSync) LookupLastSuccessfulSync(ctx context.Context, nodeID uint64, syncType string) (ns *types.NodeSync, err error) {
	s, _, err := store.SearchFederationNodeSyncs(ctx, svc.store, types.NodeSyncFilter{
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
