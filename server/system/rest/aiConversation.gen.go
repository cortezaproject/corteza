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

func (ctrl *AiConversation) List(ctx context.Context, r *request.AiConversationList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *AiConversation) Read(ctx context.Context, r *request.AiConversationRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.AiConversationID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *AiConversation) Delete(ctx context.Context, r *request.AiConversationDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.AiConversationID)
}

func (ctrl *AiConversation) Undelete(ctx context.Context, r *request.AiConversationUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.AiConversationID)
}
