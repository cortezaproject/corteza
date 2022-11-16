package rdbms

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

func (s Store) convertComposeChartFilter(f types.ChartFilter) (query squirrel.SelectBuilder, err error) {
	query = s.composeChartsSelectBuilder()

	query = filter.StateCondition(query, "cch.deleted_at", f.Deleted)

	if len(f.ChartID) > 0 {
		query = query.Where(squirrel.Eq{"cch.id": f.ChartID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"cch.id": f.LabeledIDs})
	}

	if f.NamespaceID > 0 {
		query = query.Where("cch.rel_namespace = ?", f.NamespaceID)
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cch.handle)": q},
			squirrel.Like{"LOWER(cch.name)": q},
		})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"LOWER(cch.handle)": strings.ToLower(f.Handle)})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"LOWER(cch.name)": strings.ToLower(f.Name)})
	}

	return
}
