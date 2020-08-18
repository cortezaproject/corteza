package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"strings"
)

func (s Store) convertComposePageFilter(f types.PageFilter) (query squirrel.SelectBuilder, err error) {
	query = s.QueryComposePages()

	query = rh.FilterNullByState(query, "cpg.deleted_at", f.Deleted)

	if f.NamespaceID > 0 {
		query = query.Where("cpg.rel_namespace = ?", f.NamespaceID)
	}

	if f.ParentID > 0 {
		query = query.Where("self_id = ?", f.ParentID)
	} else if f.Root {
		query = query.Where("self_id = 0")
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cpg.handle)": q},
			squirrel.Like{"LOWER(cpg.title)": q},
			squirrel.Like{"LOWER(cpg.description)": q},
		})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"LOWER(cpg.handle)": strings.ToLower(f.Handle)})
	}

	return
}
