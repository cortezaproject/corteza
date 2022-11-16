package rdbms

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
)

func (s Store) convertComposePageFilter(f types.PageFilter) (query squirrel.SelectBuilder, err error) {
	query = s.composePagesSelectBuilder()

	query = filter.StateCondition(query, "cpg.deleted_at", f.Deleted)

	if f.NamespaceID > 0 {
		query = query.Where("cpg.rel_namespace = ?", f.NamespaceID)
	}

	if f.ModuleID > 0 {
		query = query.Where("cpg.rel_module = ?", f.ModuleID)
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"cpg.id": f.LabeledIDs})
	}

	if f.ParentID > 0 {
		query = query.Where("self_id = ?", f.ParentID)
	} else if f.Root {
		query = query.Where("self_id = 0")
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cpg.handle)": q},
			squirrel.Like{"LOWER(cpg.title)": q},
			squirrel.Like{"LOWER(cpg.description)": q},
		})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"LOWER(cpg.handle)": strings.ToLower(f.Handle)})
	}

	if f.Title != "" {
		query = query.Where(squirrel.Eq{"LOWER(cpg.title)": strings.ToLower(f.Title)})
	}

	return
}

func (s Store) ReorderComposePages(ctx context.Context, namespaceID, parentID uint64, pageIDs []uint64) (err error) {
	var (
		pages   types.PageSet
		pageMap = map[uint64]bool{}
		weight  = 1

		f = types.PageFilter{ParentID: parentID, NamespaceID: namespaceID}
	)

	if pages, _, err = s.SearchComposePages(ctx, f); err != nil {
		return
	}

	for _, page := range pages {
		pageMap[page.ID] = true
	}

	// honor parameter first
	for _, pageID := range pageIDs {
		if pageMap[pageID] {
			pageMap[pageID] = false
			err = s.execUpdateComposePages(ctx,
				squirrel.Eq{"cpg.id": pageID, "cpg.self_id": parentID},
				store.Payload{"weight": weight})

			if err != nil {
				return
			}

			weight++
		}
	}

	for pageID, update := range pageMap {
		if !update {
			continue
		}

		err = s.execUpdateComposePages(ctx,
			squirrel.Eq{"cpg.id": pageID, "cpg.self_id": parentID},
			store.Payload{"weight": weight})

		if err != nil {
			return
		}

		weight++

	}

	return
}
