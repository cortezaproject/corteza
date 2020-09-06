package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
)

type (
	NodeIdentity struct{}
)

func (NodeIdentity) New() *NodeIdentity {
	return &NodeIdentity{}
}

func (ctrl NodeIdentity) GenerateNodeIdentity(ctx context.Context, r *request.IdentityGenerateNodeIdentity) (interface{}, error) {
	fmt.Println("GenerateNOdeIdentity")
	return nil, nil
}

func (ctrl NodeIdentity) RegisterOriginNode(ctx context.Context, r *request.IdentityRegisterOriginNode) (interface{}, error) {
	fmt.Println("RegisterORiginNode")
	return nil, nil
}
