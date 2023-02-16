package rest

import (
	"context"
	"mime/multipart"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
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

func (ctrl *Icon) List(ctx context.Context, r *request.IconList) (interface{}, error) {
	var (
		err error
		f   = types.AttachmentFilter{
			Kind: types.IconAttachment,
		}
		set types.AttachmentSet
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err = ctrl.attachment.Find(ctx, f)
	return ctrl.makeIconFilterPayload(ctx, set, f, err)
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

func (ctrl *Icon) makeIconFilterPayload(ctx context.Context, nn types.AttachmentSet, f types.AttachmentFilter, err error) (*iconSetPayload, error) {
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
