package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
)

type (
	NodePairRequest struct{}
)

func (NodePairRequest) New() *NodePairRequest {
	return &NodePairRequest{}
}

func (ctrl NodePairRequest) RequestPairing(ctx context.Context, r *request.PairRequestRequestPairing) (interface{}, error) {
	fmt.Println("RequestPairing")
	return nil, nil
}
