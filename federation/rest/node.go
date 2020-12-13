package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/titpetric/factory/resputil"
)

type (
	nodeServicer interface {
		Search(ctx context.Context, f types.NodeFilter) (types.NodeSet, types.NodeFilter, error)
		Create(ctx context.Context, n *types.Node) (*types.Node, error)
		CreateFromPairingURI(ctx context.Context, uri string) (*types.Node, error)

		Pair(ctx context.Context, nodeID uint64) error
		HandshakeConfirm(ctx context.Context, nodeID uint64) error
		HandshakeComplete(ctx context.Context, nodeID uint64, token string) error

		RegenerateNodeURI(ctx context.Context, nodeID uint64) (string, error)
	}

	Node struct {
		svcNode nodeServicer
	}

	moduleSetPayload struct {
		Filter types.NodeFilter `json:"filter"`
		Set    []*nodePayload   `json:"set"`
	}

	nodePayload struct {
		*types.Node
	}
)

func (Node) New() *Node {
	return &Node{
		svcNode: service.DefaultNode,
	}
}

func (ctrl Node) Search(ctx context.Context, r *request.NodeSearch) (interface{}, error) {
	set, f, err := ctrl.svcNode.Search(ctx, types.NodeFilter{
		Query:  r.Query,
		Status: r.Status,
	})

	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl Node) Create(ctx context.Context, r *request.NodeCreate) (interface{}, error) {
	if r.PairingURI != "" {
		return ctrl.svcNode.CreateFromPairingURI(ctx, r.PairingURI)
	} else {
		n := &types.Node{
			BaseURL: r.BaseURL,
			Name:    r.Name,
		}

		return ctrl.svcNode.Create(ctx, n)
	}
}

func (ctrl Node) GenerateURI(ctx context.Context, r *request.NodeGenerateURI) (interface{}, error) {
	return ctrl.svcNode.RegenerateNodeURI(ctx, r.NodeID)
}

func (ctrl Node) Pair(ctx context.Context, r *request.NodePair) (interface{}, error) {
	return resputil.OK(), ctrl.svcNode.Pair(ctx, r.NodeID)
}

func (ctrl Node) HandshakeConfirm(ctx context.Context, r *request.NodeHandshakeConfirm) (interface{}, error) {
	return resputil.OK(), ctrl.svcNode.HandshakeConfirm(ctx, r.NodeID)
}

func (ctrl Node) HandshakeComplete(ctx context.Context, r *request.NodeHandshakeComplete) (interface{}, error) {
	return resputil.OK(), ctrl.svcNode.HandshakeComplete(ctx, r.NodeID, r.TokenA)
}

func (ctrl Node) makePayload(ctx context.Context, m *types.Node, err error) (*nodePayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &nodePayload{
		Node: m,
	}, nil
}

func (ctrl Node) makeFilterPayload(ctx context.Context, nn types.NodeSet, f types.NodeFilter, err error) (*moduleSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &moduleSetPayload{Filter: f, Set: make([]*nodePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
