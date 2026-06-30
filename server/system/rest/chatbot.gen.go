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

func (ctrl *Chatbot) List(ctx context.Context, r *request.ChatbotList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Chatbot) Create(ctx context.Context, r *request.ChatbotCreate) (interface{}, error) {
	res := &types.Chatbot{
		Handle:     r.Handle,
		ProjectID:  r.ProjectID,
		Name:       r.Name,
		Enabled:    r.Enabled,
		SessionTTL: r.SessionTTL,
		Labels:     r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chatbot) Read(ctx context.Context, r *request.ChatbotRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ChatbotID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chatbot) Update(ctx context.Context, r *request.ChatbotUpdate) (interface{}, error) {
	res := &types.Chatbot{
		ID:         r.ChatbotID,
		Handle:     r.Handle,
		Name:       r.Name,
		Enabled:    r.Enabled,
		SessionTTL: r.SessionTTL,
		Labels:     r.Labels,
		UpdatedAt:  r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.svc.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Chatbot) Delete(ctx context.Context, r *request.ChatbotDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ChatbotID)
}

func (ctrl *Chatbot) Undelete(ctx context.Context, r *request.ChatbotUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ChatbotID)
}
