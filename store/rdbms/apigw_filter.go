package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertApigwFilterFilter(f types.ApigwFilterFilter) (query squirrel.SelectBuilder, err error) {
	query = s.apigwFiltersSelectBuilder()

	query = filter.StateCondition(query, "af.deleted_at", f.Deleted)
	query = filter.StateConditionNegBool(query, "af.enabled", f.Disabled)

	if f.RouteID > 0 {
		query = query.Where(squirrel.Eq{"af.rel_route": f.RouteID})
	}

	return
}
