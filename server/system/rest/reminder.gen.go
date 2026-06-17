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

func (ctrl *Reminder) List(ctx context.Context, r *request.ReminderList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.reminder.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Reminder) Create(ctx context.Context, r *request.ReminderCreate) (interface{}, error) {
	res := &types.Reminder{
		Resource:   r.Resource,
		AssignedTo: r.AssignedTo,
		RemindAt:   r.RemindAt,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.reminder.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Reminder) Update(ctx context.Context, r *request.ReminderUpdate) (interface{}, error) {
	res := &types.Reminder{
		ID:         r.ReminderID,
		Resource:   r.Resource,
		AssignedTo: r.AssignedTo,
		RemindAt:   r.RemindAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.reminder.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Reminder) Read(ctx context.Context, r *request.ReminderRead) (interface{}, error) {
	res, err := ctrl.reminder.FindByID(ctx, r.ReminderID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Reminder) Delete(ctx context.Context, r *request.ReminderDelete) (interface{}, error) {
	return api.OK(), ctrl.reminder.DeleteByID(ctx, r.ReminderID)
}
