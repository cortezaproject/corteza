package rest

import (
	"context"
	"errors"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	Node struct {
		svcNode service.NodeService
	}
)

var (
	ErrInvalidNodeCreateParams = errors.New("create node: missing or invalid parameters")
)

func (Node) New() *Node {
	return &Node{
		svcNode: service.DefaultNode,
	}
}

func (ctrl Node) Create(ctx context.Context, r *request.NodeCreate) (interface{}, error) {
	if r.NodeURI == "" && (r.Domain == "" || r.Name == "") {
		return nil, ErrInvalidNodeCreateParams
	}

	if r.NodeURI != "" {
		return ctrl.svcNode.CreateFromURI(ctx, r.NodeURI, r.MyDomain)
	} else if r.Domain != "" && r.Name != "" {
		n := &types.Node{}
		n.Domain = r.Domain
		n.Name = r.Name
		n.Status = types.NodeStatusPending

		return ctrl.svcNode.Create(ctx, n, r.MyDomain)
	}

	return nil, ErrInvalidNodeCreateParams
}

func (ctrl Node) Pair(ctx context.Context, r *request.NodePair) (interface{}, error) {
	err := ctrl.svcNode.Pair(ctx, r.NodeID)
	return nil, err
}

func (ctrl Node) HandshakeConfirm(ctx context.Context, r *request.NodeHandshakeConfirm) (interface{}, error) {
	err := ctrl.svcNode.HandshakeConfirm(ctx, r.NodeID)
	return nil, err
}

func (ctrl Node) HandshakeComplete(ctx context.Context, r *request.NodeHandshakeComplete) (interface{}, error) {
	err := ctrl.svcNode.HandshakeComplete(ctx, r.NodeID, r.TokenA)
	return nil, err
}
