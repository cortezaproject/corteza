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
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/types"
)

func (ctrl *KnowledgeBase) List(ctx context.Context, r *request.KnowledgeBaseList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *KnowledgeBase) Create(ctx context.Context, r *request.KnowledgeBaseCreate) (interface{}, error) {
	res := &types.KnowledgeBase{
		Handle:      r.Handle,
		Title:       r.Title,
		Description: r.Description,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *KnowledgeBase) Read(ctx context.Context, r *request.KnowledgeBaseRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.KnowledgeBaseID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *KnowledgeBase) Update(ctx context.Context, r *request.KnowledgeBaseUpdate) (interface{}, error) {
	res := &types.KnowledgeBase{
		ID:          r.KnowledgeBaseID,
		Handle:      r.Handle,
		Title:       r.Title,
		Description: r.Description,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *KnowledgeBase) Delete(ctx context.Context, r *request.KnowledgeBaseDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.KnowledgeBaseID)
}

func (ctrl *KnowledgeBase) Undelete(ctx context.Context, r *request.KnowledgeBaseUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.KnowledgeBaseID)
}
