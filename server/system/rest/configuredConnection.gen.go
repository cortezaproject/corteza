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

func (ctrl *ConfiguredConnection) List(ctx context.Context, r *request.ConfiguredConnectionList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.configuredConnection.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *ConfiguredConnection) Read(ctx context.Context, r *request.ConfiguredConnectionRead) (interface{}, error) {
	res, err := ctrl.configuredConnection.FindByID(ctx, r.ConnectionID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ConfiguredConnection) Delete(ctx context.Context, r *request.ConfiguredConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.configuredConnection.DeleteByID(ctx, r.ConnectionID)
}
