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

func (ctrl *UserGroup) List(ctx context.Context, r *request.UserGroupList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.userGroup.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *UserGroup) Create(ctx context.Context, r *request.UserGroupCreate) (interface{}, error) {
	res := &types.UserGroup{
		Handle: r.Handle,
		Labels: r.Labels,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.userGroup.Create(ctx, res)
	if err == nil {
		err = ctrl.afterCreate(ctx, res, r)
	}
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *UserGroup) Update(ctx context.Context, r *request.UserGroupUpdate) (interface{}, error) {
	res := &types.UserGroup{
		ID:        r.UserGroupID.Num(),
		Handle:    r.Handle,
		Labels:    r.Labels,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.userGroup.Update(ctx, res)
	if err == nil {
		err = ctrl.afterUpdate(ctx, res, r)
	}
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *UserGroup) Read(ctx context.Context, r *request.UserGroupRead) (interface{}, error) {
	res, err := ctrl.userGroup.FindByID(ctx, r.UserGroupID.Num())
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *UserGroup) Delete(ctx context.Context, r *request.UserGroupDelete) (interface{}, error) {
	return api.OK(), ctrl.userGroup.DeleteByID(ctx, r.UserGroupID.Num())
}

func (ctrl *UserGroup) Undelete(ctx context.Context, r *request.UserGroupUndelete) (interface{}, error) {
	return api.OK(), ctrl.userGroup.UndeleteByID(ctx, r.UserGroupID.Num())
}
