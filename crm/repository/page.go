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
	mod := &types.Page{}
	return mod, r.db().Get(mod, "SELECT * FROM crm_page WHERE id=?", id)
}

func (r *page) FindByModuleID(id uint64) (*types.Page, error) {
	mod := &types.Page{}
	return mod, r.db().Get(mod, "SELECT * FROM crm_page WHERE module_id=?", id)
}

func (r *page) Find() ([]*types.Page, error) {
	mod := make([]*types.Page, 0)
	return mod, r.db().Select(&mod, "SELECT * FROM crm_page ORDER BY id ASC")
}

func (r *page) Create(mod *types.Page) (*types.Page, error) {
	mod.ID = factory.Sonyflake.NextID()
	return mod, r.db().Insert("crm_page", mod)
}

func (r *page) Update(mod *types.Page) (*types.Page, error) {
	return mod, r.db().Replace("crm_page", mod)
}

func (r *page) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_page WHERE id=?", id)
	return err
}
