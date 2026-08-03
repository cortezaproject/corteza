package rdbms

import (
	"context"
	"sort"

	composeType "github.com/crusttech/human/server/compose/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) ReorderComposePages(ctx context.Context, namespaceID uint64, parentID uint64, pageIDs []uint64) (err error) {
	var (
		pages   composeType.PageSet
		pageMap = map[uint64]bool{}
		weight  = 1

		// Root is what constrains the search to top-level pages. The filter only
		// applies self_id = ParentID when ParentID > 0, so without this a reorder
		// of root pages (parentID 0) matched every page in the namespace and the
		// leftover pass below re-weighted all of them in map-iteration order.
		f = composeType.PageFilter{ParentID: parentID, NamespaceID: namespaceID, Root: parentID == 0}

		query = func(id uint64, weight int) *goqu.UpdateDataset {
			return s.Dialect.GOQU().
				Update(composePageTable).
				Set(goqu.Record{"weight": weight}).
				Where(goqu.C("id").Eq(id))
		}
	)

	if pages, _, err = s.SearchComposePages(ctx, f); err != nil {
		return
	}

	for _, app := range pages {
		pageMap[app.ID] = true
	}

	// honor parameter first
	for _, id := range pageIDs {
		if pageMap[id] {
			pageMap[id] = false

			if err = s.Exec(ctx, query(id, weight)); err != nil {
				return
			}

			weight++
		}
	}

	// Pages the caller did not name keep their relative order. Ranging over
	// pageMap directly would not: Go randomises map iteration, so reordering a
	// subset silently shuffled every other page, differently on each call.
	leftovers := make(composeType.PageSet, 0, len(pageMap))
	for _, p := range pages {
		if pageMap[p.ID] {
			leftovers = append(leftovers, p)
		}
	}
	sort.SliceStable(leftovers, func(i, j int) bool { return leftovers[i].Weight < leftovers[j].Weight })

	for _, p := range leftovers {
		if err = s.Exec(ctx, query(p.ID, weight)); err != nil {
			return
		}

		weight++

	}

	return
}
