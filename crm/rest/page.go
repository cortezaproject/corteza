package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

type (
	Page struct {
		page service.PageService
	}
)

func (Page) New(pageSvc service.PageService) *Page {
	return &Page{pageSvc}
}

func (ctrl *Page) List(ctx context.Context, r *request.PageList) (interface{}, error) {
	return ctrl.page.With(ctx).Find()
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
	return ctrl.page.With(ctx).FindByID(r.ID)
}

func (ctrl *Page) Edit(ctx context.Context, r *request.PageEdit) (interface{}, error) {
	p := &types.Page{
		ID:          r.ID,
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
	return resputil.OK(), ctrl.page.With(ctx).DeleteByID(r.ID)
}
