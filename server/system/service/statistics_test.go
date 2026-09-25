package service

import (
	"testing"
	"time"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestStatsResolveRange(t *testing.T) {
	now := time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC)
	ptr := func(t time.Time) *time.Time { return &t }

	t.Run("defaults to the last 30 days by day", func(t *testing.T) {
		r, bucket, err := statsResolveRange(StatisticsRequest{}, now)
		require.NoError(t, err)
		require.Equal(t, now, r.To)
		require.Equal(t, now.AddDate(0, 0, -30), r.From)
		require.Equal(t, types.SystemStatsBucketDay, bucket)
	})

	t.Run("picks week and month buckets by span", func(t *testing.T) {
		_, bucket, err := statsResolveRange(StatisticsRequest{From: ptr(now.AddDate(0, 0, -90))}, now)
		require.NoError(t, err)
		require.Equal(t, types.SystemStatsBucketWeek, bucket)

		_, bucket, err = statsResolveRange(StatisticsRequest{From: ptr(now.AddDate(-1, 0, 0))}, now)
		require.NoError(t, err)
		require.Equal(t, types.SystemStatsBucketMonth, bucket)
	})

	t.Run("keeps an explicit bucket", func(t *testing.T) {
		_, bucket, err := statsResolveRange(StatisticsRequest{From: ptr(now.AddDate(-1, 0, 0)), Bucket: "week"}, now)
		require.NoError(t, err)
		require.Equal(t, types.SystemStatsBucketWeek, bucket)
	})

	t.Run("rejects an empty, inverted, too long or unknown request", func(t *testing.T) {
		_, _, err := statsResolveRange(StatisticsRequest{From: ptr(now)}, now)
		require.Error(t, err)

		_, _, err = statsResolveRange(StatisticsRequest{From: ptr(now.AddDate(0, 0, 1))}, now)
		require.Error(t, err)

		_, _, err = statsResolveRange(StatisticsRequest{From: ptr(now.AddDate(-4, 0, 0))}, now)
		require.Error(t, err)

		_, _, err = statsResolveRange(StatisticsRequest{Bucket: "hour"}, now)
		require.Error(t, err)
	})
}

func TestStatsBuckets(t *testing.T) {
	// Wed 2026-09-23 → Sat 2026-10-03, local midnight boundaries
	from := time.Date(2026, 9, 23, 10, 0, 0, 0, time.Local)
	to := time.Date(2026, 10, 3, 9, 0, 0, 0, time.Local)
	r := types.SystemStatsRange{From: from, To: to}

	t.Run("day", func(t *testing.T) {
		b := statsBuckets(r, types.SystemStatsBucketDay)
		require.Len(t, b.labels, 11)
		require.Equal(t, "2026-09-23", b.labels[0])
		require.Equal(t, "2026-10-03", b.labels[10])
	})

	t.Run("week starts on the Monday containing from", func(t *testing.T) {
		b := statsBuckets(r, types.SystemStatsBucketWeek)
		require.Equal(t, []string{"2026-09-21", "2026-09-28"}, b.labels)
	})

	t.Run("month", func(t *testing.T) {
		b := statsBuckets(r, types.SystemStatsBucketMonth)
		require.Equal(t, []string{"2026-09-01", "2026-10-01"}, b.labels)
	})

	t.Run("roll sums days into their bucket per key and drops strays", func(t *testing.T) {
		b := statsBuckets(r, types.SystemStatsBucketWeek)
		rolled := b.roll([]types.SystemStatsDaily{
			{Day: "2026-09-23", Key: "ok", Count: 2},
			{Day: "2026-09-27", Key: "ok", Count: 3},
			{Day: "2026-09-29", Key: "ok", Count: 1},
			{Day: "2026-09-29", Key: "error", Count: 4},
			{Day: "2026-09-01", Key: "ok", Count: 99},
			{Day: "garbage", Key: "ok", Count: 99},
		})

		require.Equal(t, []uint{5, 1}, rolled["ok"])
		require.Equal(t, []uint{0, 4}, rolled["error"])
		require.Equal(t, []uint{5, 5}, b.sum(rolled["ok"], rolled["error"]))
		require.Equal(t, []uint{0, 0}, b.orEmpty(rolled["missing"]))
	})
}
