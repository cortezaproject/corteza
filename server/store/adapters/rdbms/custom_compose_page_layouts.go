package rdbms

import (
	"context"

	composeType "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/doug-martin/goqu/v9"
)

func (s Store) ReorderComposePageLayouts(ctx context.Context, namespaceID uint64, pageID uint64, pageLayoutLayoutIDs []uint64) (err error) {
	var (
		pageLayouts   composeType.PageLayoutSet
		pageLayoutMap = map[uint64]bool{}
		weight        = 1

		f = composeType.PageLayoutFilter{PageID: pageID, NamespaceID: namespaceID}

		query = func(id uint64, weight int) *goqu.UpdateDataset {
			return s.Dialect.GOQU().
				Update(composePageLayoutTable).
				Set(goqu.Record{"weight": weight}).
				Where(goqu.C("id").Eq(id))
		}
	)

	if pageLayouts, _, err = s.SearchComposePageLayouts(ctx, f); err != nil {
		return
	}

	for _, app := range pageLayouts {
		pageLayoutMap[app.ID] = true
	}

	// honor parameter first
	for _, id := range pageLayoutLayoutIDs {
		if pageLayoutMap[id] {
			pageLayoutMap[id] = false

			if err = s.Exec(ctx, query(id, weight)); err != nil {
				return
			}

			weight++
		}
	}

	for id, update := range pageLayoutMap {
		if !update {
			continue
		}

		if err = s.Exec(ctx, query(id, weight)); err != nil {
			return
		}

		weight++

	}

	return
}
