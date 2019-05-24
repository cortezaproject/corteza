package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/payload"
)

type (
	pagePayload struct {
		*types.Page

		CanGrant      bool `json:"canGrant"`
		CanUpdatePage bool `json:"canUpdatePage"`
		CanDeletePage bool `json:"canDeletePage"`
	}

	pageSetPayload struct {
		Filter types.PageFilter `json:"filter"`
		Set    []*pagePayload   `json:"set"`
	}

	Page struct {
		page       service.PageService
		attachment service.AttachmentService
		ac         pageAccessController
	}

	pageAccessController interface {
		CanGrant(context.Context) bool

		CanUpdatePage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool
	}
)

func (Page) New() *Page {
	return &Page{
		page:       service.DefaultPage,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl *Page) List(ctx context.Context, r *request.PageList) (interface{}, error) {
	f := types.PageFilter{
		NamespaceID: r.NamespaceID,
		ParentID:    r.SelfID,

		Query:   r.Query,
		PerPage: r.PerPage,
		Page:    r.Page,
	}

	set, filter, err := ctrl.page.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Page) Tree(ctx context.Context, r *request.PageTree) (interface{}, error) {
	return ctrl.page.With(ctx).Tree(r.NamespaceID)
}

func (ctrl *Page) Create(ctx context.Context, r *request.PageCreate) (interface{}, error) {
	var (
		err error
		mod = &types.Page{
			NamespaceID: r.NamespaceID,
			SelfID:      r.SelfID,
			ModuleID:    r.ModuleID,
			Title:       r.Title,
			Description: r.Description,
			Blocks:      r.Blocks,
			Visible:     r.Visible,
		}
	)

	mod, err = ctrl.page.With(ctx).Create(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Page) Read(ctx context.Context, r *request.PageRead) (interface{}, error) {
	mod, err := ctrl.page.With(ctx).FindByID(r.NamespaceID, r.PageID)
	return ctrl.makePayload(ctx, mod, err)

}

func (ctrl *Page) Reorder(ctx context.Context, r *request.PageReorder) (interface{}, error) {
	return resputil.OK(), ctrl.page.With(ctx).Reorder(r.NamespaceID, r.SelfID, payload.ParseUInt64s(r.PageIDs))
}

func (ctrl *Page) Update(ctx context.Context, r *request.PageUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.Page{
			NamespaceID: r.NamespaceID,
			ID:          r.PageID,
			SelfID:      r.SelfID,
			ModuleID:    r.ModuleID,
			Title:       r.Title,
			Description: r.Description,
			Blocks:      r.Blocks,
			Visible:     r.Visible,
		}
	)

	mod, err = ctrl.page.With(ctx).Update(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Page) Delete(ctx context.Context, r *request.PageDelete) (interface{}, error) {
	return resputil.OK(), ctrl.page.With(ctx).DeleteByID(r.NamespaceID, r.PageID)
}

func (ctrl *Page) Upload(ctx context.Context, r *request.PageUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.With(ctx).CreatePageAttachment(
		r.NamespaceID,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.PageID,
	)

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Page) makePayload(ctx context.Context, c *types.Page, err error) (*pagePayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &pagePayload{
		Page: c,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdatePage: ctrl.ac.CanUpdatePage(ctx, c),
		CanDeletePage: ctrl.ac.CanDeletePage(ctx, c),
	}, nil
}

func (ctrl Page) makeFilterPayload(ctx context.Context, nn types.PageSet, f types.PageFilter, err error) (*pageSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &pageSetPayload{Filter: f, Set: make([]*pagePayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}
