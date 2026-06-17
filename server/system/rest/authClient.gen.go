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

func (ctrl *AuthClient) List(ctx context.Context, r *request.AuthClientList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.authClient.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *AuthClient) Create(ctx context.Context, r *request.AuthClientCreate) (interface{}, error) {
	res := &types.AuthClient{
		Handle:      r.Handle,
		ValidGrant:  r.ValidGrant,
		RedirectURI: r.RedirectURI,
		Scope:       r.Scope,
		Trusted:     r.Trusted,
		Enabled:     r.Enabled,
		ValidFrom:   r.ValidFrom,
		ExpiresAt:   r.ExpiresAt,
		Labels:      r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.authClient.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *AuthClient) Update(ctx context.Context, r *request.AuthClientUpdate) (interface{}, error) {
	res := &types.AuthClient{
		ID:          r.ClientID,
		Handle:      r.Handle,
		ValidGrant:  r.ValidGrant,
		RedirectURI: r.RedirectURI,
		Scope:       r.Scope,
		Trusted:     r.Trusted,
		Enabled:     r.Enabled,
		ValidFrom:   r.ValidFrom,
		ExpiresAt:   r.ExpiresAt,
		Labels:      r.Labels,
		UpdatedAt:   r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.authClient.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *AuthClient) Read(ctx context.Context, r *request.AuthClientRead) (interface{}, error) {
	res, err := ctrl.authClient.FindByID(ctx, r.ClientID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *AuthClient) Delete(ctx context.Context, r *request.AuthClientDelete) (interface{}, error) {
	return api.OK(), ctrl.authClient.DeleteByID(ctx, r.ClientID)
}

func (ctrl *AuthClient) Undelete(ctx context.Context, r *request.AuthClientUndelete) (interface{}, error) {
	return api.OK(), ctrl.authClient.UndeleteByID(ctx, r.ClientID)
}
