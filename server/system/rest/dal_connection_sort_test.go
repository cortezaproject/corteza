package rest

import (
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func ts(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func conn(handle string, created string, updated, deleted *time.Time) *types.DalConnection {
	return &types.DalConnection{
		Handle:    handle,
		CreatedAt: *ts(created),
		UpdatedAt: updated,
		DeletedAt: deleted,
	}
}

func handles(set types.DalConnectionSet) (out []string) {
	for _, c := range set {
		out = append(out, c.Handle)
	}
	return
}

// The column shows the most recent of deleted/updated/created, so the order has
// to come from the same value — sorting on updatedAt alone leaves every
// never-updated row with no key at all.
func TestSortDalConnectionSet_changedAt(t *testing.T) {
	set := types.DalConnectionSet{
		conn("never-updated", "2026-03-01T10:00:00Z", nil, nil),
		conn("updated-latest", "2026-01-01T10:00:00Z", ts("2026-05-01T10:00:00Z"), nil),
		conn("deleted-oldest", "2026-02-01T10:00:00Z", ts("2026-02-02T10:00:00Z"), ts("2026-02-03T10:00:00Z")),
	}

	ss, err := filter.NewSorting("coalesce(deletedAt,updatedAt,createdAt) DESC")
	require.NoError(t, err)

	sortDalConnectionSet(set, ss.Sort)
	require.Equal(t, []string{"updated-latest", "never-updated", "deleted-oldest"}, handles(set))

	ss, err = filter.NewSorting("coalesce(deletedAt,updatedAt,createdAt)")
	require.NoError(t, err)

	sortDalConnectionSet(set, ss.Sort)
	require.Equal(t, []string{"deleted-oldest", "never-updated", "updated-latest"}, handles(set))
}

func TestSortDalConnectionSet_plainColumn(t *testing.T) {
	set := types.DalConnectionSet{
		conn("charlie", "2026-01-01T10:00:00Z", nil, nil),
		conn("alpha", "2026-01-01T10:00:00Z", nil, nil),
		conn("bravo", "2026-01-01T10:00:00Z", nil, nil),
	}

	ss, err := filter.NewSorting("handle")
	require.NoError(t, err)

	sortDalConnectionSet(set, ss.Sort)
	require.Equal(t, []string{"alpha", "bravo", "charlie"}, handles(set))
}

// No expression means the caller sent no sort — the store's own order stands.
func TestSortDalConnectionSet_noSortKeepsOrder(t *testing.T) {
	set := types.DalConnectionSet{
		conn("charlie", "2026-01-01T10:00:00Z", nil, nil),
		conn("alpha", "2026-01-01T10:00:00Z", nil, nil),
	}

	sortDalConnectionSet(set, nil)
	require.Equal(t, []string{"charlie", "alpha"}, handles(set))
}

// A federation node arrives with no timestamps at all; it must not compare as
// though it changed at the zero time in the middle of the set.
func TestTimeSortKey(t *testing.T) {
	require.Equal(t, "", timeSortKey(nil))
	require.Equal(t, "", timeSortKey(&time.Time{}))
	require.Equal(t, "2026-05-01T10:00:00.000000000", timeSortKey(ts("2026-05-01T10:00:00Z")))

	// Comparable as strings across offsets: both are 12:00 UTC.
	require.Equal(t, timeSortKey(ts("2026-05-01T14:00:00+02:00")), timeSortKey(ts("2026-05-01T12:00:00Z")))
}
