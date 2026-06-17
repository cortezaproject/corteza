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

func (ctrl *DalSensitivityLevel) List(ctx context.Context, r *request.DalSensitivityLevelList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.dalSensitivityLevel.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *DalSensitivityLevel) Create(ctx context.Context, r *request.DalSensitivityLevelCreate) (interface{}, error) {
	res := &types.DalSensitivityLevel{
		Handle: r.Handle,
		Level:  r.Level,
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.dalSensitivityLevel.Create(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalSensitivityLevel) Update(ctx context.Context, r *request.DalSensitivityLevelUpdate) (interface{}, error) {
	res := &types.DalSensitivityLevel{
		ID:        r.SensitivityLevelID,
		Handle:    r.Handle,
		Level:     r.Level,
		UpdatedAt: r.UpdatedAt,
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.dalSensitivityLevel.Update(ctx, res)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalSensitivityLevel) Read(ctx context.Context, r *request.DalSensitivityLevelRead) (interface{}, error) {
	res, err := ctrl.dalSensitivityLevel.FindByID(ctx, r.SensitivityLevelID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *DalSensitivityLevel) Delete(ctx context.Context, r *request.DalSensitivityLevelDelete) (interface{}, error) {
	return api.OK(), ctrl.dalSensitivityLevel.DeleteByID(ctx, r.SensitivityLevelID)
}

func (ctrl *DalSensitivityLevel) Undelete(ctx context.Context, r *request.DalSensitivityLevelUndelete) (interface{}, error) {
	return api.OK(), ctrl.dalSensitivityLevel.UndeleteByID(ctx, r.SensitivityLevelID)
}
