package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/federation_nodes_sync.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	FederationNodesSyncs interface {
		SearchFederationNodesSyncs(ctx context.Context, f types.NodeSyncFilter) (types.NodeSyncSet, types.NodeSyncFilter, error)
		LookupFederationNodesSyncByNodeID(ctx context.Context, node_id uint64) (*types.NodeSync, error)
		LookupFederationNodesSyncByNodeIDSyncTypeSyncStatus(ctx context.Context, node_id uint64, sync_type string, sync_status string) (*types.NodeSync, error)

		CreateFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) error

		UpdateFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) error

		UpsertFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) error

		DeleteFederationNodesSync(ctx context.Context, rr ...*types.NodeSync) error
		DeleteFederationNodesSyncByNodeID(ctx context.Context, nodeID uint64) error

		TruncateFederationNodesSyncs(ctx context.Context) error
	}
)

var _ *types.NodeSync
var _ context.Context

// SearchFederationNodesSyncs returns all matching FederationNodesSyncs from store
func SearchFederationNodesSyncs(ctx context.Context, s FederationNodesSyncs, f types.NodeSyncFilter) (types.NodeSyncSet, types.NodeSyncFilter, error) {
	return s.SearchFederationNodesSyncs(ctx, f)
}

// LookupFederationNodesSyncByNodeID searches for sync activity by node ID
//
// It returns sync activity
func LookupFederationNodesSyncByNodeID(ctx context.Context, s FederationNodesSyncs, node_id uint64) (*types.NodeSync, error) {
	return s.LookupFederationNodesSyncByNodeID(ctx, node_id)
}

// LookupFederationNodesSyncByNodeIDSyncTypeSyncStatus searches for activity by node, type and status
//
// It returns sync activity
func LookupFederationNodesSyncByNodeIDSyncTypeSyncStatus(ctx context.Context, s FederationNodesSyncs, node_id uint64, sync_type string, sync_status string) (*types.NodeSync, error) {
	return s.LookupFederationNodesSyncByNodeIDSyncTypeSyncStatus(ctx, node_id, sync_type, sync_status)
}

// CreateFederationNodesSync creates one or more FederationNodesSyncs in store
func CreateFederationNodesSync(ctx context.Context, s FederationNodesSyncs, rr ...*types.NodeSync) error {
	return s.CreateFederationNodesSync(ctx, rr...)
}

// UpdateFederationNodesSync updates one or more (existing) FederationNodesSyncs in store
func UpdateFederationNodesSync(ctx context.Context, s FederationNodesSyncs, rr ...*types.NodeSync) error {
	return s.UpdateFederationNodesSync(ctx, rr...)
}

// UpsertFederationNodesSync creates new or updates existing one or more FederationNodesSyncs in store
func UpsertFederationNodesSync(ctx context.Context, s FederationNodesSyncs, rr ...*types.NodeSync) error {
	return s.UpsertFederationNodesSync(ctx, rr...)
}

// DeleteFederationNodesSync Deletes one or more FederationNodesSyncs from store
func DeleteFederationNodesSync(ctx context.Context, s FederationNodesSyncs, rr ...*types.NodeSync) error {
	return s.DeleteFederationNodesSync(ctx, rr...)
}

// DeleteFederationNodesSyncByNodeID Deletes FederationNodesSync from store
func DeleteFederationNodesSyncByNodeID(ctx context.Context, s FederationNodesSyncs, nodeID uint64) error {
	return s.DeleteFederationNodesSyncByNodeID(ctx, nodeID)
}

// TruncateFederationNodesSyncs Deletes all FederationNodesSyncs from store
func TruncateFederationNodesSyncs(ctx context.Context, s FederationNodesSyncs) error {
	return s.TruncateFederationNodesSyncs(ctx)
}
