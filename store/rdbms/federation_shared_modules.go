package rdbms

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
)

func (s Store) convertFederationSharedModuleFilter(f types.SharedModuleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationSharedModulesSelectBuilder()

	// query = filter.StateCondition(query, "cmd.deleted_at", f.Deleted)

	if f.NodeID > 0 {
		query = query.Where("cmd.rel_node = ?", f.NodeID)
	}

	if f.Handle != "" {
		query = query.Where("cmd.handle = ?", f.Handle)
	}

	if f.Name != "" {
		query = query.Where("cmd.name = ?", f.Name)
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(cmd.name)": q},
			squirrel.Like{"LOWER(cmd.handle)": q},
		})
	}

	return
}
