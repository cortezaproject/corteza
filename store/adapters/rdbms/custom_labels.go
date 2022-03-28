package rdbms

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

func (s Store) DeleteExtraLabels(ctx context.Context, kind string, resourceId uint64, name ...string) error {
	return s.Exec(ctx, labelDeleteQuery(
		s.config.Dialect,
		goqu.C("kind").Eq(kind),
		goqu.C("rel_resource").Eq(resourceId),
		goqu.C("name").NotIn(name),
	))
}
