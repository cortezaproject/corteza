package rdbms

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/federation/types"
)

func (s Store) convertFederationExposedModuleFilter(f types.ExposedModuleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationExposedModulesSelectBuilder()

	if f.NodeID > 0 {
		query = query.Where("cmd.rel_node = ?", f.NodeID)
	}

	if f.ComposeModuleID > 0 {
		query = query.Where("cmd.rel_compose_module = ?", f.ComposeModuleID)
	}

	if f.ComposeNamespaceID > 0 {
		query = query.Where("cmd.rel_compose_namespace = ?", f.ComposeNamespaceID)
	}

	if f.LastSync > 0 {
		t := time.Unix(int64(f.LastSync), 0)

		if !t.IsZero() {
			query = query.Where("(cmd.updated_at >= ? OR cmd.created_at >= ?)", t.UTC().Format(time.RFC3339), t.UTC().Format(time.RFC3339))
		}
	}

	return
}
