package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/crusttech/human/server/federation/rest/request"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/api"
)

func (ctrl *Node) Read(ctx context.Context, r *request.NodeRead) (interface{}, error) {
	res, err := ctrl.node.FindByID(ctx, r.NodeID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Node) Update(ctx context.Context, r *request.NodeUpdate) (interface{}, error) {
	res := &types.Node{
		ID:      r.NodeID,
		Name:    r.Name,
		Contact: r.Contact,
		BaseURL: r.BaseURL,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.node.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Node) Delete(ctx context.Context, r *request.NodeDelete) (interface{}, error) {
	return api.OK(), ctrl.node.DeleteByID(ctx, r.NodeID)
}

func (ctrl *Node) Undelete(ctx context.Context, r *request.NodeUndelete) (interface{}, error) {
	return api.OK(), ctrl.node.UndeleteByID(ctx, r.NodeID)
}
