package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/label/types"
)

func (s Store) convertLabelFilter(f types.LabelFilter) (query squirrel.SelectBuilder, err error) {
	query = s.labelsSelectBuilder()

	query = query.Where(squirrel.Eq{"lbl.kind": f.Kind})

	if len(f.ResourceID) > 0 {
		query = query.Where(squirrel.Eq{"lbl.rel_resource": f.ResourceID})
	}

	if len(f.Filter) > 0 {
		kvOr := squirrel.Or{}
		for k, v := range f.Filter {
			kvOr = append(kvOr, squirrel.Eq{"name": k, "value": v})
		}
		query = query.Where(kvOr)

	}

	return
}

// DeleteExtraLabels removes all unlisted labels
func (s Store) DeleteExtraLabels(ctx context.Context, kind string, resourceID uint64, names ...string) (err error) {
	return s.execDeleteLabels(ctx, squirrel.And{
		squirrel.Eq{
			"kind":         kind,
			"rel_resource": resourceID,
		},
		squirrel.NotEq{
			"name": names,
		},
	})
}
