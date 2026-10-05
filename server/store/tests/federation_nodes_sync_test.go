package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
)

func testFederationNodeSyncs(t *testing.T, s store.FederationNodeSyncs) {
	var (
		ctx = context.Background()
		req = require.New(t)

		nodeID  = id.Next()
		moduleA = id.Next()
		moduleB = id.Next()
		now     = time.Now().UTC().Truncate(time.Second)

		sync = func(moduleID uint64, status string, at time.Time) *types.NodeSync {
			return &types.NodeSync{
				NodeID:       nodeID,
				ModuleID:     moduleID,
				SyncType:     types.NodeSyncTypeStructure,
				SyncStatus:   status,
				TimeOfAction: at,
			}
		}
	)

	req.NoError(s.TruncateFederationNodeSyncs(ctx))
	req.NoError(s.CreateFederationNodeSync(ctx,
		sync(moduleA, types.NodeSyncStatusError, now.Add(-2*time.Minute)),
		sync(moduleA, types.NodeSyncStatusSuccess, now.Add(-time.Minute)),
		sync(moduleB, types.NodeSyncStatusError, now),
	))

	t.Run("search the last sync of one module", func(t *testing.T) {
		req := require.New(t)

		set, _, err := s.SearchFederationNodeSyncs(ctx, types.NodeSyncFilter{
			NodeID:   nodeID,
			ModuleID: moduleA,
			SyncType: types.NodeSyncTypeStructure,
			Sorting: filter.Sorting{
				Sort: filter.SortExprSet{&filter.SortExpr{Column: "time_of_action", Descending: true}},
			},
			Paging: filter.Paging{Limit: 1},
		})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(moduleA, set[0].ModuleID)
		req.Equal(types.NodeSyncStatusSuccess, set[0].SyncStatus)
		req.True(now.Add(-time.Minute).Equal(set[0].TimeOfAction))
	})
}
