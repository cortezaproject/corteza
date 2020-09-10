package rest

import (
	"context"
	"errors"
	"net/url"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/federation/util"
)

type (
	Node struct {
		svcNode service.NodeService
	}
)

var (
	ErrInvalidNodeCreateParams = errors.New("create node: missing or invalid parameters")
	ErrorIdentityMissingToken  = errors.New("identity: token missing")
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

	n := &types.Node{}
	if r.NodeURI != "" {
		uri, err := url.QueryUnescape(r.NodeURI)
		if err != nil {
			return nil, err
		}
		pr, err := util.DecodeURI(uri)
		if err != nil {
			return nil, err
		}
		n.Domain = pr.Domain
		n.Name = pr.Params.Name
		n.Status = types.NodeStatusPending
		n.NodeURI = r.NodeURI
	} else {
		n.Domain = r.Domain
		n.Name = r.Name
		n.Status = types.NodeStatusPending
	}

	return ctrl.svcNode.Create(ctx, n, r.MyDomain)
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
