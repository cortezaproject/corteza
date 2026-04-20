package rest

import (
	"context"
	"github.com/crusttech/human/server/federation/rest/request"
	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/pkg/api"
)

type (
	handshakeInitializer interface {
		HandshakeInit(context.Context, uint64, string, uint64, string) error
	}

	NodeHandshake struct {
		svcNode handshakeInitializer
	}
)

func (NodeHandshake) New() *NodeHandshake {
	return &NodeHandshake{
		svcNode: service.DefaultNode,
	}
}

func (ctrl NodeHandshake) Initialize(ctx context.Context, r *request.NodeHandshakeInitialize) (interface{}, error) {
	return api.OK(), ctrl.svcNode.HandshakeInit(ctx, r.NodeID, r.PairToken, r.SharedNodeID, r.AuthToken)
}
