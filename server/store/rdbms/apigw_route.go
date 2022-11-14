package rdbms

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertApigwRouteFilter(f types.ApigwRouteFilter) (query squirrel.SelectBuilder, err error) {
	query = s.apigwRoutesSelectBuilder()
	query = filter.StateCondition(query, "ar.deleted_at", f.Deleted)
	query = filter.StateConditionNegBool(query, "ar.enabled", f.Disabled)

	if f.Query != "" {
		query = query.Where(squirrel.Like{"LOWER(ar.endpoint)": "%" + strings.ToLower(f.Query) + "%"})
	}

	return
}
