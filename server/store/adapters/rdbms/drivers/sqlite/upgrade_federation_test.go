package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// oldNodeSyncStore is a SQLite store holding federation_nodes_sync the way
// releases before 2025 left it: a primary key and the module in rel_module.
func oldNodeSyncStore(t *testing.T, name string) *rdbms.Store {
	ctx := context.Background()
	s, err := sqlite.Connect(ctx, "sqlite3://file:"+name+"?mode=memory&cache=shared")
	require.NoError(t, err)

	rs := s.(*rdbms.Store)
	_, err = rs.DB.ExecContext(ctx, `CREATE TABLE federation_nodes_sync (
		rel_node INTEGER NOT NULL, rel_module INTEGER NOT NULL, sync_type TEXT, sync_status TEXT,
		time_of_action TIMESTAMP NOT NULL, PRIMARY KEY (rel_node))`)
	require.NoError(t, err)
	_, err = rs.DB.ExecContext(ctx, `INSERT INTO federation_nodes_sync VALUES (1, 516687706984677377, 'sync_structure', 'success', ?)`, time.Unix(1700000000, 0).UTC())
	require.NoError(t, err)
	return rs
}

// The store upgrade, run twice the way every boot runs it, leaves the table
// without a primary key, its timestamp still a time, and the module in the
// column the model reads.
func TestUpgradeFederationNodeSync(t *testing.T) {
	ctx := context.Background()
	s := oldNodeSyncStore(t, "upgrade_federation_node_sync")

	for i := 0; i < 2; i++ {
		require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s), "run %d", i+1)
	}

	var pk int
	require.NoError(t, s.DB.QueryRowxContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('federation_nodes_sync') WHERE pk > 0`).Scan(&pk))
	require.Zero(t, pk)

	var (
		moduleID uint64
		at       time.Time
	)
	require.NoError(t, s.DB.QueryRowxContext(ctx, `SELECT rel_compose_module, time_of_action FROM federation_nodes_sync`).Scan(&moduleID, &at))
	require.Equal(t, uint64(516687706984677377), moduleID)
	require.Equal(t, int64(1700000000), at.Unix())
}
