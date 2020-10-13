package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
)

func (s Store) convertFederationExposedModuleFilter(f types.ExposedModuleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationExposedModulesSelectBuilder()

	// query = filter.StateCondition(query, "cmd.deleted_at", f.Deleted)

	if f.NodeID > 0 {
		query = query.Where("cmd.rel_node = ?", f.NodeID)
	}

	if f.ComposeModuleID > 0 {
		query = query.Where("cmd.rel_compose_module = ?", f.ComposeModuleID)
	}

	if f.ComposeNamespaceID > 0 {
		query = query.Where("cmd.rel_compose_namespace = ?", f.ComposeNamespaceID)
	}

	// if f.Query != "" {
	// 	q := "%" + strings.ToLower(f.Query) + "%"
	// 	query = query.Where(squirrel.Or{
	// 		squirrel.Like{"LOWER(cmd.name)": q},
	// 		squirrel.Like{"LOWER(cmd.handle)": q},
	// 	})
	// }

	return
}
