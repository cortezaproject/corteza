package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"strings"
)

func (s Store) convertComposeModuleFilter(f types.ModuleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.composeModulesSelectBuilder()

	query = rh.FilterNullByState(query, "cmd.deleted_at", f.Deleted)

	if f.NamespaceID > 0 {
		query = query.Where("cmd.rel_namespace = ?", f.NamespaceID)
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cmd.name)": q},
			squirrel.Like{"LOWER(cmd.handle)": q},
		})
	}

	if f.Name != "" {
		query = query.Where(squirrel.Eq{"LOWER(cmd.name)": strings.ToLower(f.Name)})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"LOWER(cmd.handle)": strings.ToLower(f.Handle)})
	}

	return
}
