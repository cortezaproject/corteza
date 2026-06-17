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
	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/api"
)

func (ctrl *PageLayout) List(ctx context.Context, r *request.PageLayoutList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.pageLayout.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *PageLayout) Create(ctx context.Context, r *request.PageLayoutCreate) (interface{}, error) {
	res := &types.PageLayout{
		NamespaceID: r.NamespaceID,
		ParentID:    r.ParentID,
		Weight:      r.Weight,
		Handle:      r.Handle,
		Labels:      r.Labels,
		OwnedBy:     r.OwnedBy,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.pageLayout.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *PageLayout) Read(ctx context.Context, r *request.PageLayoutRead) (interface{}, error) {
	res, err := ctrl.pageLayout.FindByID(ctx, r.NamespaceID, r.PageLayoutID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *PageLayout) Update(ctx context.Context, r *request.PageLayoutUpdate) (interface{}, error) {
	res := &types.PageLayout{
		ID:          r.PageLayoutID,
		NamespaceID: r.NamespaceID,
		ParentID:    r.ParentID,
		Weight:      r.Weight,
		Handle:      r.Handle,
		Labels:      r.Labels,
		OwnedBy:     r.OwnedBy,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.pageLayout.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *PageLayout) Delete(ctx context.Context, r *request.PageLayoutDelete) (interface{}, error) {
	return api.OK(), ctrl.pageLayout.DeleteByID(ctx, r.NamespaceID, r.PageID, r.PageLayoutID)
}

func (ctrl *PageLayout) Undelete(ctx context.Context, r *request.PageLayoutUndelete) (interface{}, error) {
	return api.OK(), ctrl.pageLayout.UndeleteByID(ctx, r.NamespaceID, r.PageID, r.PageLayoutID)
}
