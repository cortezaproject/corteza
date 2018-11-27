package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	PageRepository interface {
		With(ctx context.Context, db *factory.DB) PageRepository

		FindByID(id uint64) (*types.Page, error)
		FindByModuleID(id uint64) (*types.Page, error)
		FindBySelfID(selfID uint64) (types.PageSet, error)
		FindAll() (types.PageSet, error)

		Create(mod *types.Page) (*types.Page, error)
		Update(mod *types.Page) (*types.Page, error)
		DeleteByID(id uint64) error

		Reorder(selfID uint64, pageIDs []uint64) error
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
	return page, nil
}

func (r *page) FindByModuleID(id uint64) (*types.Page, error) {
	page := &types.Page{}
	if err := r.db().Get(page, "SELECT * FROM crm_page WHERE module_id=?", id); err != nil {
		return nil, err
	}
	return page, nil
}

func (r *page) FindBySelfID(selfID uint64) (types.PageSet, error) {
	pages := types.PageSet{}
	if err := r.db().Select(&pages, "SELECT * FROM crm_page WHERE self_id = ? ORDER BY weight ASC", selfID); err != nil {
		return pages, err
	}
	return pages, nil
}

func (r *page) FindAll() (types.PageSet, error) {
	pages := types.PageSet{}
	return pages, r.db().Select(&pages, "SELECT * FROM crm_page ORDER BY self_id, weight ASC")
}

func (r *page) Reorder(selfID uint64, pageIDs []uint64) error {
	pageMap := map[uint64]bool{}
	if pages, err := r.FindBySelfID(selfID); err != nil {
		return nil
	} else {
		for _, page := range pages {
			pageMap[page.ID] = true
		}
	}
	weight := 1
	db := r.db()
	// honor parameter first
	for _, pageID := range pageIDs {
		if pageMap[pageID] {
			pageMap[pageID] = false
			if _, err := db.Exec("UPDATE crm_page set weight=? where id=? and self_id=?", weight, pageID, selfID); err != nil {
				return err
			}
			weight++
		}
	}
	for pageID, update := range pageMap {
		if update {
			if _, err := db.Exec("UPDATE crm_page set weight=? where id=? and self_id=?", weight, pageID, selfID); err != nil {
				return err
			}
			weight++
		}
	}
	return nil
}

func (r *page) Create(item *types.Page) (*types.Page, error) {
	page := &types.Page{}
	*page = *item

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
