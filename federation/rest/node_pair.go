package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
)

type (
	NodePair struct{}
)

func (NodePair) New() *NodePair {
	return &NodePair{}
}

func (ctrl NodePair) ApprovePairing(ctx context.Context, r *request.PairApprovePairing) (interface{}, error) {
	fmt.Println("ApprovePairing")
	return nil, nil
}
func (ctrl NodePair) CompletePairing(ctx context.Context, r *request.PairCompletePairing) (interface{}, error) {
	fmt.Println("CompletePairing")
	return nil, nil
}
