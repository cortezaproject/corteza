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

func (ctrl *User) List(ctx context.Context, r *request.UserList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.user.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *User) Read(ctx context.Context, r *request.UserRead) (interface{}, error) {
	res, err := ctrl.user.FindByID(ctx, r.UserID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *User) Delete(ctx context.Context, r *request.UserDelete) (interface{}, error) {
	return api.OK(), ctrl.user.DeleteByID(ctx, r.UserID)
}

func (ctrl *User) Undelete(ctx context.Context, r *request.UserUndelete) (interface{}, error) {
	return api.OK(), ctrl.user.UndeleteByID(ctx, r.UserID)
}
