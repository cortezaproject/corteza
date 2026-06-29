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

func (ctrl *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.module.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Module) Create(ctx context.Context, r *request.ModuleCreate) (interface{}, error) {
	res := &types.Module{
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		ProjectID:   r.ProjectID,
		Handle:      r.Handle,
		Labels:      r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.module.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	res, err := ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Module) Update(ctx context.Context, r *request.ModuleUpdate) (interface{}, error) {
	res := &types.Module{
		ID:          r.ModuleID,
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		Handle:      r.Handle,
		Labels:      r.Labels,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.module.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	return api.OK(), ctrl.module.DeleteByID(ctx, r.NamespaceID, r.ModuleID)
}
