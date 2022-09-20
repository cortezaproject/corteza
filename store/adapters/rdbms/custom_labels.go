package rdbms

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

func (s Store) DeleteExtraLabels(ctx context.Context, kind string, resourceId uint64, name ...string) error {
	var (
		expr = []goqu.Expression{
			goqu.C("kind").Eq(kind),
			goqu.C("rel_resource").Eq(resourceId),
		}
	)

	if len(name) > 0 {
		expr = append(expr, goqu.C("name").NotIn(name))
	}

	return s.Exec(ctx, labelDeleteQuery(s.Dialect.GOQU(), expr...))
}
