package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	PageRepository interface {
		With(ctx context.Context, db *factory.DB) PageRepository

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
		FindByModuleID(namespaceID, moduleID uint64) (*types.Page, error)
		Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)

		Create(mod *types.Page) (*types.Page, error)
		Update(mod *types.Page) (*types.Page, error)
		DeleteByID(namespaceID, pageID uint64) error

		Reorder(namespaceID, selfID uint64, pageIDs []uint64) error
	}

	page struct {
		*repository
	}
)

const (
	ErrPageNotFound = repositoryError("PageNotFound")
)

func Page(ctx context.Context, db *factory.DB) PageRepository {
	return (&page{}).With(ctx, db)
}

func (r page) With(ctx context.Context, db *factory.DB) PageRepository {
	return &page{
		repository: r.repository.With(ctx, db),
	}
}

func (r page) table() string {
	return "compose_page"
}

func (r page) columns() []string {
	return []string{
		"id", "rel_namespace", "self_id", "rel_module", "title",
		"blocks", "description", "visible", "weight",
		"created_at", "updated_at", "deleted_at",
	}
}

func (r page) query() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table()).
		Where("deleted_at IS NULL")
}

func (r page) FindByID(namespaceID, pageID uint64) (*types.Page, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("id = ?", pageID)

		c = &types.Page{}
	)

	if namespaceID > 0 {
		query = query.Where("rel_namespace = ?", namespaceID)
	}

	return c, isFound(r.fetchOne(c, query), c.ID > 0, ErrPageNotFound)
}

func (r page) FindByModuleID(namespaceID, moduleID uint64) (*types.Page, error) {
	var (
		query = r.query().
			Columns(r.columns()...).
			Where("rel_module = ?", moduleID)

		c = &types.Page{}
	)

	if namespaceID > 0 {
		query = query.Where("rel_namespace = ?", namespaceID)
	}

	return c, isFound(r.fetchOne(c, query), c.ID > 0, ErrPageNotFound)
}

func (r page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	f = filter

	query := r.query()

	if filter.NamespaceID > 0 {
		query = query.Where("rel_namespace = ?", filter.NamespaceID)
	}

	if filter.ParentID > 0 {
		query = query.Where("self_id = ?", filter.ParentID)
	} else if filter.Root {
		query = query.Where("self_id = 0")
	}

	if f.Query != "" {
		q := "%" + f.Query + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", q, q)
	}

	if f.Count, err = r.count(query); err != nil || f.Count == 0 {
		return
	}

	query = query.
		Columns(r.columns()...).
		OrderBy("weight ASC")

	return set, f, r.fetchPaged(&set, query, f.Page, f.PerPage)
}

func (r page) Reorder(namespaceID, parentID uint64, pageIDs []uint64) error {
	var (
		pageMap = map[uint64]bool{}
		filter  = types.PageFilter{NamespaceID: namespaceID, ParentID: parentID}
	)

	if pages, _, err := r.Find(filter); err != nil {
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
			if _, err := db.Exec("UPDATE compose_page SET weight = ? WHERE id = ? AND self_id = ?", weight, pageID, parentID); err != nil {
				return err
			}
			weight++
		}
	}
	for pageID, update := range pageMap {
		if update {
			if _, err := db.Exec("UPDATE compose_page SET weight = ? WHERE id = ? AND self_id = ?", weight, pageID, parentID); err != nil {
				return err
			}
			weight++
		}
	}
	return nil
}

func (r page) Create(mod *types.Page) (*types.Page, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now().Truncate(time.Second)

	return mod, r.db().Insert(r.table(), mod)
}

func (r page) Update(mod *types.Page) (*types.Page, error) {
	now := time.Now().Truncate(time.Second)
	mod.UpdatedAt = &now

	return mod, r.db().Replace(r.table(), mod)
}

func (r page) DeleteByID(namespaceID, pageID uint64) error {
	_, err := r.db().Exec(
		"UPDATE "+r.table()+" SET deleted_at = NOW() WHERE rel_namespace = ? AND id = ?",
		namespaceID,
		pageID,
	)

	return err
}
