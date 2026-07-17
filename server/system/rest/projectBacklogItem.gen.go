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

func (ctrl *ProjectBacklogItem) List(ctx context.Context, r *request.ProjectBacklogItemList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.projectBacklogItem.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ProjectBacklogItem) Create(ctx context.Context, r *request.ProjectBacklogItemCreate) (interface{}, error) {
	res := &types.ProjectBacklogItem{
		ProjectID:   r.ProjectID,
		Title:       r.Title,
		Description: r.Description,
		Category:    r.Category,
		EventID:     r.EventID,
		Assignee:    r.Assignee,
		Priority:    r.Priority,
		Status:      r.Status,
		DateDue:     r.DateDue,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectBacklogItem.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectBacklogItem) Read(ctx context.Context, r *request.ProjectBacklogItemRead) (interface{}, error) {
	res, err := ctrl.projectBacklogItem.FindByID(ctx, r.BacklogItemID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectBacklogItem) Update(ctx context.Context, r *request.ProjectBacklogItemUpdate) (interface{}, error) {
	res := &types.ProjectBacklogItem{
		ID:          r.BacklogItemID,
		Title:       r.Title,
		Description: r.Description,
		Category:    r.Category,
		EventID:     r.EventID,
		Assignee:    r.Assignee,
		Priority:    r.Priority,
		Status:      r.Status,
		DateDue:     r.DateDue,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.projectBacklogItem.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectBacklogItem) Delete(ctx context.Context, r *request.ProjectBacklogItemDelete) (interface{}, error) {
	return api.OK(), ctrl.projectBacklogItem.DeleteByID(ctx, r.BacklogItemID)
}
