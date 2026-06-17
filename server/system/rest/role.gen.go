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
)

func (ctrl *Role) List(ctx context.Context, r *request.RoleList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.role.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Role) Read(ctx context.Context, r *request.RoleRead) (interface{}, error) {
	res, err := ctrl.role.FindByID(ctx, r.RoleID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Role) Delete(ctx context.Context, r *request.RoleDelete) (interface{}, error) {
	return api.OK(), ctrl.role.DeleteByID(ctx, r.RoleID)
}

func (ctrl *Role) Undelete(ctx context.Context, r *request.RoleUndelete) (interface{}, error) {
	return api.OK(), ctrl.role.UndeleteByID(ctx, r.RoleID)
}
