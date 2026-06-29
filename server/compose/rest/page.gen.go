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
)

func (ctrl *Page) List(ctx context.Context, r *request.PageList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.page.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Page) Create(ctx context.Context, r *request.PageCreate) (interface{}, error) {
	res := &types.Page{
		NamespaceID: r.NamespaceID,
		SelfID:      r.SelfID,
		ModuleID:    r.ModuleID,
		ProjectID:   r.ProjectID,
		Title:       r.Title,
		Handle:      r.Handle,
		Description: r.Description,
		Weight:      r.Weight,
		Labels:      r.Labels,
		Visible:     r.Visible,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.page.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Page) Read(ctx context.Context, r *request.PageRead) (interface{}, error) {
	res, err := ctrl.page.FindByID(ctx, r.NamespaceID, r.PageID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Page) Update(ctx context.Context, r *request.PageUpdate) (interface{}, error) {
	res := &types.Page{
		ID:          r.PageID,
		NamespaceID: r.NamespaceID,
		SelfID:      r.SelfID,
		ModuleID:    r.ModuleID,
		Title:       r.Title,
		Handle:      r.Handle,
		Description: r.Description,
		Weight:      r.Weight,
		Labels:      r.Labels,
		Visible:     r.Visible,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.page.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}
