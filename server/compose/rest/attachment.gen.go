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
)

func (ctrl *Attachment) List(ctx context.Context, r *request.AttachmentList) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.attachment.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}
