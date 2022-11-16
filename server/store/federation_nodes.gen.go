package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/federation_nodes.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza/server/federation/types"
)

type (
	FederationNodes interface {
		SearchFederationNodes(ctx context.Context, f types.NodeFilter) (types.NodeSet, types.NodeFilter, error)
		LookupFederationNodeByID(ctx context.Context, id uint64) (*types.Node, error)
		LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, base_url string, shared_node_id uint64) (*types.Node, error)
		LookupFederationNodeBySharedNodeID(ctx context.Context, shared_node_id uint64) (*types.Node, error)

		CreateFederationNode(ctx context.Context, rr ...*types.Node) error

		UpdateFederationNode(ctx context.Context, rr ...*types.Node) error

		UpsertFederationNode(ctx context.Context, rr ...*types.Node) error

		DeleteFederationNode(ctx context.Context, rr ...*types.Node) error
		DeleteFederationNodeByID(ctx context.Context, ID uint64) error

		TruncateFederationNodes(ctx context.Context) error
	}
)

var _ *types.Node
var _ context.Context

// SearchFederationNodes returns all matching FederationNodes from store
func SearchFederationNodes(ctx context.Context, s FederationNodes, f types.NodeFilter) (types.NodeSet, types.NodeFilter, error) {
	return s.SearchFederationNodes(ctx, f)
}

// LookupFederationNodeByID searches for federation node by ID
//
// It returns federation node
func LookupFederationNodeByID(ctx context.Context, s FederationNodes, id uint64) (*types.Node, error) {
	return s.LookupFederationNodeByID(ctx, id)
}

// LookupFederationNodeByBaseURLSharedNodeID searches for node by shared-node-id and base-url
func LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, s FederationNodes, base_url string, shared_node_id uint64) (*types.Node, error) {
	return s.LookupFederationNodeByBaseURLSharedNodeID(ctx, base_url, shared_node_id)
}

// LookupFederationNodeBySharedNodeID searches for node by shared-node-id
func LookupFederationNodeBySharedNodeID(ctx context.Context, s FederationNodes, shared_node_id uint64) (*types.Node, error) {
	return s.LookupFederationNodeBySharedNodeID(ctx, shared_node_id)
}

// CreateFederationNode creates one or more FederationNodes in store
func CreateFederationNode(ctx context.Context, s FederationNodes, rr ...*types.Node) error {
	return s.CreateFederationNode(ctx, rr...)
}

// UpdateFederationNode updates one or more (existing) FederationNodes in store
func UpdateFederationNode(ctx context.Context, s FederationNodes, rr ...*types.Node) error {
	return s.UpdateFederationNode(ctx, rr...)
}

// UpsertFederationNode creates new or updates existing one or more FederationNodes in store
func UpsertFederationNode(ctx context.Context, s FederationNodes, rr ...*types.Node) error {
	return s.UpsertFederationNode(ctx, rr...)
}

// DeleteFederationNode Deletes one or more FederationNodes from store
func DeleteFederationNode(ctx context.Context, s FederationNodes, rr ...*types.Node) error {
	return s.DeleteFederationNode(ctx, rr...)
}

// DeleteFederationNodeByID Deletes FederationNode from store
func DeleteFederationNodeByID(ctx context.Context, s FederationNodes, ID uint64) error {
	return s.DeleteFederationNodeByID(ctx, ID)
}

// TruncateFederationNodes Deletes all FederationNodes from store
func TruncateFederationNodes(ctx context.Context, s FederationNodes) error {
	return s.TruncateFederationNodes(ctx)
}
