package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
)

func (s Store) convertFederationNodesSyncFilter(f types.NodeSyncFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationNodesSyncsSelectBuilder()

	if f.NodeID > 0 {
		query = query.Where("fdns.rel_node = ?", f.NodeID)
	}

	if f.SyncStatus != "" {
		query = query.Where("fdns.rel_compose_module = ?", f.SyncStatus)
	}

	if f.SyncType != "" {
		query = query.Where("fdns.rel_compose_namespace = ?", f.SyncType)
	}

	return
}
