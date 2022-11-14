package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/lann/builder"
)

// Count counts all rows that match conditions from given query builder
func Count(ctx context.Context, db dbLayer, q squirrel.SelectBuilder) (count uint, err error) {
	// Delete order-bys for counting
	q = builder.Delete(q, "OrderByParts").(squirrel.SelectBuilder)

	// Replace columns with count(*)
	q = builder.Delete(q, "Columns").(squirrel.SelectBuilder).Column("COUNT(*)")

	if sqlSelect, argsSelect, err := q.ToSql(); err != nil {
		return 0, err
	} else {
		if err := db.GetContext(ctx, &count, sqlSelect, argsSelect...); err != nil {
			return 0, err
		}
	}

	return count, nil
}
