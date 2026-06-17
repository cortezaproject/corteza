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
	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/pkg/api"
)

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return api.OK(), ctrl.record.DeleteByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
}

func (ctrl *Record) Undelete(ctx context.Context, r *request.RecordUndelete) (interface{}, error) {
	return api.OK(), ctrl.record.UndeleteByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
}
