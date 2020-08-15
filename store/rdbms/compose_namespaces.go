package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"strings"
)

func (s Store) convertComposeNamespaceFilter(f types.NamespaceFilter) (query squirrel.SelectBuilder, err error) {
	if f.Sort == "" {
		f.Sort = "id ASC"
	}

	query = s.QueryComposeNamespaces()

	query = rh.FilterNullByState(query, "cns.deleted_at", f.Deleted)

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cns.name)": q},
			squirrel.Like{"LOWER(cns.slug)": q},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"LOWER(cns.name)": strings.ToLower(f.Name)})
	}

	if f.Slug != "" {
		query = query.Where(squirrel.Eq{"LOWER(cns.slug)": strings.ToLower(f.Slug)})
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, s.ComposeModuleColumns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	return
}
