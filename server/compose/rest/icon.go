package rest

import (
	"context"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"mime/multipart"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	iconPayload struct {
		*attachmentPayload
	}

	iconSetPayload struct {
		Filter types.IconFilter `json:"filter"`
		Set    []*iconPayload   `json:"set"`
	}

	Icon struct {
		locale     service.ResourceTranslationsManagerService
		attachment service.AttachmentService
		ac         iconAccessController
	}

	iconAccessController interface {
		CanGrant(context.Context) bool
	}
)

func (Icon) New() *Icon {
	return &Icon{
		locale:     service.DefaultResourceTranslation,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

// makeFilter builds the icon search filter for the generated List
// controller. Icons are attachment-backed, so the filter is an
// AttachmentFilter scoped to IconAttachment and excluding deleted icons.
func (ctrl *Icon) makeFilter(ctx context.Context, r *request.IconList) (f types.AttachmentFilter, err error) {
	f = types.AttachmentFilter{
		Kind: types.IconAttachment,
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return
	}

	//Get only the undeleted icons
	f.Deleted = filter.StateExcluded

	return
}

func (ctrl *Icon) Upload(ctx context.Context, r *request.IconUpload) (interface{}, error) {
	file, err := r.Icon.Open()
	if err != nil {
		return nil, err
	}

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	a, err := ctrl.attachment.CreateIconAttachment(
		ctx,
		r.Icon.Filename,
		r.Icon.Size,
		file,
	)

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl *Icon) makeFilterPayload(ctx context.Context, nn types.AttachmentSet, f types.AttachmentFilter, err error) (*iconSetPayload, error) {
	if err != nil {
		return nil, err
	}

	var (
		a  *attachmentPayload
		ff types.IconFilter
	)

	ff.Paging = f.Paging
	ff.Sorting = f.Sorting

	res := &iconSetPayload{Filter: ff, Set: make([]*iconPayload, len(nn))}

	for i := range nn {
		a, _ = makeAttachmentPayload(ctx, nn[i], nil)
		res.Set[i] = &iconPayload{a}
	}

	return res, nil
}

func (ctrl *Icon) Delete(ctx context.Context, r *request.IconDelete) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot delete icon")
	}

	_, err := ctrl.attachment.FindByID(ctx, 0, r.IconID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.attachment.DeleteByID(ctx, 0, r.IconID)
}
