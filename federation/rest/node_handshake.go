package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
)

type (
	NodeHandshake struct {
		svcNode service.NodeService
	}
)

func (NodeHandshake) New() *NodeHandshake {
	return &NodeHandshake{
		svcNode: service.DefaultNode,
	}
}

func (ctrl NodeHandshake) Initialize(ctx context.Context, r *request.NodeHandshakeInitialize) (interface{}, error) {
	err := ctrl.svcNode.HandshakeInit(ctx, r.NodeID, r.NodeIDB, r.NodeURI, r.TokenB)
	return nil, err
}
