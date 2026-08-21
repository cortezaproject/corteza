package rest

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/stretchr/testify/require"
)

// The list endpoint carried no sort and no paging for a long time: the request
// params existed on the wire but makeFilter never read them, so the webapp's
// limit and sort were accepted and discarded.
func TestDalConnectionMakeFilter(t *testing.T) {
	ctrl := DalConnection{}

	f, err := ctrl.makeFilter(context.Background(), &request.DalConnectionList{
		ConnectionID: []string{"1", "2"},
		Handle:       "primary-database",
		Type:         "corteza::system:dal-connection",
		Deleted:      uint(filter.StateInclusive),
		IncTotal:     true,
		Limit:        50,
		PageCursor:   "",
		Sort:         "coalesce(deletedAt,updatedAt,createdAt) DESC",
	})

	require.NoError(t, err)
	require.Equal(t, []string{"1", "2"}, f.DalConnectionID)
	require.Equal(t, "primary-database", f.Handle)
	require.Equal(t, "corteza::system:dal-connection", f.Type)
	require.Equal(t, filter.StateInclusive, f.Deleted)
	require.True(t, f.IncTotal)
	require.Equal(t, uint(50), f.Limit)
	require.Equal(t, "coalesce(deletedAt,updatedAt,createdAt) DESC", f.Sort.String())
}

// Unset deleted means exclude, not "any state".
func TestDalConnectionMakeFilterDefaultsDeletedToExcluded(t *testing.T) {
	f, err := DalConnection{}.makeFilter(context.Background(), &request.DalConnectionList{})

	require.NoError(t, err)
	require.Equal(t, filter.StateExcluded, f.Deleted)
}
