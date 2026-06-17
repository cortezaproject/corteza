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

func (ctrl *Notification) List(ctx context.Context, r *request.NotificationList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.notification.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Notification) Create(ctx context.Context, r *request.NotificationCreate) (interface{}, error) {
	res := &types.Notification{
		Recipient: r.Recipient,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.notification.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Notification) Update(ctx context.Context, r *request.NotificationUpdate) (interface{}, error) {
	res := &types.Notification{
		ID:        r.NotificationID,
		Recipient: r.Recipient,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.notification.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Notification) Read(ctx context.Context, r *request.NotificationRead) (interface{}, error) {
	res, err := ctrl.notification.FindByID(ctx, r.NotificationID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Notification) Delete(ctx context.Context, r *request.NotificationDelete) (interface{}, error) {
	return api.OK(), ctrl.notification.DeleteByID(ctx, r.NotificationID)
}
