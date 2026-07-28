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

func (ctrl *ProjectTask) List(ctx context.Context, r *request.ProjectTaskList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectTask.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectTask) Create(ctx context.Context, r *request.ProjectTaskCreate) (interface{}, error) {
	res := &types.ProjectTask{
		ProjectID:     r.ProjectID,
		RevisionID:    r.RevisionID,
		Title:         r.Title,
		Description:   r.Description,
		TaskName:      r.TaskName,
		TaskType:      r.TaskType,
		Status:        r.Status,
		Severity:      r.Severity,
		Risk:          r.Risk,
		Owner:         r.Owner,
		ChangeOwner:   r.ChangeOwner,
		DateDue:       r.DateDue,
		CompletedDate: r.CompletedDate,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectTask.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectTask) Read(ctx context.Context, r *request.ProjectTaskRead) (interface{}, error) {
	res, err := ctrl.projectTask.FindByID(ctx, r.TaskID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectTask) Update(ctx context.Context, r *request.ProjectTaskUpdate) (interface{}, error) {
	res := &types.ProjectTask{
		ID:            r.TaskID,
		RevisionID:    r.RevisionID,
		Title:         r.Title,
		Description:   r.Description,
		TaskName:      r.TaskName,
		TaskType:      r.TaskType,
		Status:        r.Status,
		Severity:      r.Severity,
		Risk:          r.Risk,
		Owner:         r.Owner,
		ChangeOwner:   r.ChangeOwner,
		DateDue:       r.DateDue,
		CompletedDate: r.CompletedDate,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectTask.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectTask) Delete(ctx context.Context, r *request.ProjectTaskDelete) (interface{}, error) {
	return api.OK(), ctrl.projectTask.DeleteByID(ctx, r.TaskID)
}
