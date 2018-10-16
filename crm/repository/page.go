package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	PageRepository interface {
		With(ctx context.Context, db *factory.DB) PageRepository

		Find() ([]*types.Page, error)
		FindByID(id uint64) (*types.Page, error)
		FindByModuleID(id uint64) (*types.Page, error)

		Create(mod *types.Page) (*types.Page, error)
		Update(mod *types.Page) (*types.Page, error)
		DeleteByID(id uint64) error
	}

	page struct {
		*repository
	}
)

func Page(ctx context.Context, db *factory.DB) PageRepository {
	return (&page{}).With(ctx, db)
}

func (r *page) With(ctx context.Context, db *factory.DB) PageRepository {
	return &page{
		repository: r.repository.With(ctx, db),
	}
}

func (r *page) FindByID(id uint64) (*types.Page, error) {
	page := &types.Page{}
	if err := r.db().Get(page, "SELECT * FROM crm_page WHERE id=?", id); err != nil {
		return page, err
	}
	if err := r.fillModule(page); err != nil {
		return page, err
	}
	return page, nil
}

func (r *page) FindByModuleID(id uint64) (*types.Page, error) {
	page := &types.Page{}
	if err := r.db().Get(page, "SELECT * FROM crm_page WHERE module_id=?", id); err != nil {
		return page, err
	}
	if err := r.fillModule(page); err != nil {
		return page, err
	}
	return page, nil
}

func (r *page) Find() ([]*types.Page, error) {
	pages := make([]*types.Page, 0)
	if err := r.db().Select(&pages, "SELECT * FROM crm_page ORDER BY id ASC"); err != nil {
		return pages, err
	}
	for _, page := range pages {
		if err := r.fillModule(page); err != nil {
			return pages, err
		}
	}
	return pages, nil
}

func (r *page) Create(page *types.Page) (*types.Page, error) {
	page.ID = factory.Sonyflake.NextID()
	return page, r.db().Insert("crm_page", page)
}

func (r *page) Update(page *types.Page) (*types.Page, error) {
	return page, r.db().Replace("crm_page", page)
}

func (r *page) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_page WHERE id=?", id)
	return err
}

func (r *page) fillModule(page *types.Page) error {
	if page.ModuleID > 0 {
		api := Module(r.Context(), r.db())
		module, err := api.FindByID(page.ModuleID)
		page.Module = module
		return err
	}
	return nil
}
