package service

import (
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func connSortTs(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func connSortShorts(set types.ConnectionSet) (out []string) {
	for _, c := range set {
		out = append(out, c.Meta.Short)
	}
	return
}

// A catalog entry is appended after the store query, so the merged set is
// re-sorted here; a COALESCE expression carries its columns on the expression
// rather than in Column, and reading Column alone left the set untouched.
func TestSortConnectionSet_changedAt(t *testing.T) {
	mk := func(short, created string, updated *time.Time) *types.Connection {
		return &types.Connection{
			Meta:      types.ConnectionMeta{Short: short},
			CreatedAt: *connSortTs(created),
			UpdatedAt: updated,
		}
	}

	set := types.ConnectionSet{
		mk("never-updated", "2026-03-01T10:00:00Z", nil),
		mk("catalog-only", "2026-04-01T10:00:00Z", nil),
		mk("updated-latest", "2026-01-01T10:00:00Z", connSortTs("2026-05-01T10:00:00Z")),
	}

	ss, err := filter.NewSorting("coalesce(deletedAt,updatedAt,createdAt) DESC")
	require.NoError(t, err)

	sortConnectionSet(set, ss.Sort)
	require.Equal(t, []string{"updated-latest", "catalog-only", "never-updated"}, connSortShorts(set))
}

func TestSortConnectionSet_defaultsToNameAscending(t *testing.T) {
	set := types.ConnectionSet{
		{Meta: types.ConnectionMeta{Short: "Zeta"}},
		{Meta: types.ConnectionMeta{Short: "Alpha"}},
	}

	sortConnectionSet(set, nil)
	require.Equal(t, []string{"Alpha", "Zeta"}, connSortShorts(set))
}
