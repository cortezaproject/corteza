package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/federation/types"
)

func (s Store) convertFederationNodesSyncFilter(f types.NodeSyncFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationNodesSyncsSelectBuilder()

	if f.NodeID > 0 {
		query = query.Where("fdns.rel_node = ?", f.NodeID)
	}

	if f.ModuleID > 0 {
		query = query.Where("fdns.rel_module = ?", f.ModuleID)
	}

	if f.SyncStatus != "" {
		query = query.Where("fdns.sync_status = ?", f.SyncStatus)
	}

	if f.SyncType != "" {
		query = query.Where("fdns.sync_type = ?", f.SyncType)
	}

	return
}
