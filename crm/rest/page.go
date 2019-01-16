package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/payload"
)

type (
	Page struct {
		page service.PageService
	}
)

func (Page) New() *Page {
	return &Page{
		page: service.DefaultPage,
	}
}

func (ctrl *Page) List(ctx context.Context, r *request.PageList) (interface{}, error) {
	if r.SelfID > 0 {
		return ctrl.page.With(ctx).FindBySelfID(r.SelfID)
	} else {
		return ctrl.page.With(ctx).Find()
	}
}

func (ctrl *Page) Tree(ctx context.Context, r *request.PageTree) (interface{}, error) {
	return ctrl.page.With(ctx).Tree()
}

func (ctrl *Page) Create(ctx context.Context, r *request.PageCreate) (interface{}, error) {
	p := &types.Page{
		SelfID:      r.SelfID,
		ModuleID:    r.ModuleID,
		Title:       r.Title,
		Description: r.Description,
		Blocks:      r.Blocks,
		Visible:     r.Visible,
	}
	return ctrl.page.With(ctx).Create(p)
}

func (ctrl *Page) Read(ctx context.Context, r *request.PageRead) (interface{}, error) {
	return ctrl.page.With(ctx).FindByID(r.PageID)
}

func (ctrl *Page) Reorder(ctx context.Context, r *request.PageReorder) (interface{}, error) {
	return resputil.OK(), ctrl.page.With(ctx).Reorder(r.SelfID, payload.ParseUInt64s(r.PageIDs))
}

func (ctrl *Page) Update(ctx context.Context, r *request.PageUpdate) (interface{}, error) {
	p := &types.Page{
		ID:          r.PageID,
		SelfID:      r.SelfID,
		ModuleID:    r.ModuleID,
		Title:       r.Title,
		Description: r.Description,
		Blocks:      r.Blocks,
		Visible:     r.Visible,
	}
	return ctrl.page.With(ctx).Update(p)
}

func (ctrl *Page) Delete(ctx context.Context, r *request.PageDelete) (interface{}, error) {
	return resputil.OK(), ctrl.page.With(ctx).DeleteByID(r.PageID)
}
