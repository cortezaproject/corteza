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

func (ctrl *Namespace) List(ctx context.Context, r *request.NamespaceList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.namespace.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Namespace) Create(ctx context.Context, r *request.NamespaceCreate) (interface{}, error) {
	res := &types.Namespace{
		Name:    r.Name,
		Labels:  r.Labels,
		Slug:    r.Slug,
		Enabled: r.Enabled,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.namespace.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Namespace) Read(ctx context.Context, r *request.NamespaceRead) (interface{}, error) {
	res, err := ctrl.namespace.FindByID(ctx, r.NamespaceID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Namespace) Update(ctx context.Context, r *request.NamespaceUpdate) (interface{}, error) {
	res := &types.Namespace{
		ID:        r.NamespaceID,
		Name:      r.Name,
		Slug:      r.Slug,
		Enabled:   r.Enabled,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.namespace.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Namespace) Delete(ctx context.Context, r *request.NamespaceDelete) (interface{}, error) {
	return api.OK(), ctrl.namespace.DeleteByID(ctx, r.NamespaceID)
}
