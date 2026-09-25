package rdbms

import (
	"context"

	systemType "github.com/crusttech/human/server/system/types"
	"github.com/doug-martin/goqu/v9"
)

func (s *Store) CountUsers(ctx context.Context, f systemType.UserFilter) (c uint, _ error) {
	var (
		aux = struct {
			Count uint `db:"count"`
		}{}

		expr, _, err = s.Filters.User(s, f)

		query = s.Dialect.GOQU().
			From(userTable).
			Select(goqu.COUNT(goqu.Star()).As("count"))
	)

	if err != nil {
		return
	}

	if err = s.QueryOne(ctx, query.Where(expr...).Limit(1), &aux); err != nil {
		return
	}

	return aux.Count, nil
}
