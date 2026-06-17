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

func (ctrl *Template) List(ctx context.Context, r *request.TemplateList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.renderer.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Template) Create(ctx context.Context, r *request.TemplateCreate) (interface{}, error) {
	res := &types.Template{
		Handle:   r.Handle,
		Language: r.Language,
		Partial:  r.Partial,
		Template: r.Template,
		OwnerID:  r.OwnerID,
		Labels:   r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.renderer.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Template) Read(ctx context.Context, r *request.TemplateRead) (interface{}, error) {
	res, err := ctrl.renderer.FindByID(ctx, r.TemplateID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Template) Update(ctx context.Context, r *request.TemplateUpdate) (interface{}, error) {
	res := &types.Template{
		ID:        r.TemplateID,
		Handle:    r.Handle,
		Language:  r.Language,
		Partial:   r.Partial,
		Template:  r.Template,
		OwnerID:   r.OwnerID,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.renderer.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Template) Delete(ctx context.Context, r *request.TemplateDelete) (interface{}, error) {
	return api.OK(), ctrl.renderer.DeleteByID(ctx, r.TemplateID)
}

func (ctrl *Template) Undelete(ctx context.Context, r *request.TemplateUndelete) (interface{}, error) {
	return api.OK(), ctrl.renderer.UndeleteByID(ctx, r.TemplateID)
}
