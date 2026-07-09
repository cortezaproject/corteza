package actionlog_test

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestActionlogReport(t *testing.T) {
	req := require.New(t)
	s := setupStore(t)
	ctx := context.Background()
	svc := actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDebugPolicy())
	req.NoError(store.TruncateActionlogs(ctx, s))

	var (
		day1 = time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
		day2 = time.Date(2026, 6, 2, 10, 0, 0, 0, time.UTC)

		mkAction = func(id uint64, ts time.Time, resource, action, errMsg string, actorID uint64) *actionlog.Action {
			return &actionlog.Action{
				ID:            id,
				Timestamp:     ts,
				Resource:      resource,
				Action:        action,
				Error:         errMsg,
				ActorID:       actorID,
				RequestOrigin: actionlog.RequestOrigin_API_REST,
				Severity:      actionlog.Notice,
				Meta:          actionlog.Meta{},
				Delta:         actionlog.Delta{},
				OldState:      actionlog.OldState("{}"),
			}
		}
	)

	req.NoError(store.CreateActionlog(ctx, s,
		mkAction(1, day1, testResource1, "create", "", 101),
		mkAction(2, day1.Add(5*time.Hour+30*time.Minute), testResource1, "create", "boom", 102),
		mkAction(3, day2, testResource1, "update", "", 101),
		mkAction(4, day2.Add(2*time.Hour), testResource2, "create", "", 103),
	))

	report := func(rr actionlog.ReportRequest) actionlog.ReportRowSet {
		set, err := svc.Report(ctx, rr)
		req.NoError(err)
		return set
	}

	t.Run("single dimension, all metrics", func(t *testing.T) {
		set := report(actionlog.ReportRequest{
			Dimensions: []string{"resource"},
			Metrics:    []string{"count", "actors", "errors"},
		})

		require.Len(t, set, 2)

		require.Equal(t, testResource1, set[0].Dimensions["resource"])
		require.Equal(t, float64(3), set[0].Metrics["count"])
		require.Equal(t, float64(2), set[0].Metrics["actors"])
		require.Equal(t, float64(1), set[0].Metrics["errors"])

		require.Equal(t, testResource2, set[1].Dimensions["resource"])
		require.Equal(t, float64(1), set[1].Metrics["count"])
		require.Equal(t, float64(1), set[1].Metrics["actors"])
		require.Equal(t, float64(0), set[1].Metrics["errors"])
	})

	t.Run("multiple dimensions, deterministic order", func(t *testing.T) {
		set := report(actionlog.ReportRequest{
			Dimensions: []string{"resource", "action"},
		})

		require.Len(t, set, 3)
		require.Equal(t, "create", set[0].Dimensions["action"])
		require.Equal(t, float64(2), set[0].Metrics["count"])
		require.Equal(t, "update", set[1].Dimensions["action"])
		require.Equal(t, testResource2, set[2].Dimensions["resource"])
	})

	t.Run("date dimension truncates distinct times to one day group", func(t *testing.T) {
		set := report(actionlog.ReportRequest{Dimensions: []string{"day"}})

		require.Len(t, set, 2)
		require.Equal(t, "2026-06-01", set[0].Dimensions["day"])
		require.Equal(t, float64(2), set[0].Metrics["count"])
		require.Equal(t, "2026-06-02", set[1].Dimensions["day"])
		require.Equal(t, float64(2), set[1].Metrics["count"])
	})

	t.Run("origin filter", func(t *testing.T) {
		set := report(actionlog.ReportRequest{
			Dimensions: []string{"resource"},
			Filter:     actionlog.Filter{Origin: actionlog.RequestOrigin_API_REST},
		})

		require.Len(t, set, 2)

		set = report(actionlog.ReportRequest{
			Dimensions: []string{"resource"},
			Filter:     actionlog.Filter{Origin: actionlog.RequestOrigin_API_GRPC},
		})

		require.Empty(t, set)
	})

	t.Run("actor dimension serialized as string", func(t *testing.T) {
		set := report(actionlog.ReportRequest{Dimensions: []string{"actor"}})

		require.Len(t, set, 3)
		require.Equal(t, "101", set[0].Dimensions["actor"])
		require.Equal(t, float64(2), set[0].Metrics["count"])
	})

	t.Run("no dimensions yields global totals", func(t *testing.T) {
		set := report(actionlog.ReportRequest{Metrics: []string{"count", "actors"}})

		require.Len(t, set, 1)
		require.Equal(t, float64(4), set[0].Metrics["count"])
		require.Equal(t, float64(3), set[0].Metrics["actors"])
	})

	t.Run("filter constrains rows", func(t *testing.T) {
		set := report(actionlog.ReportRequest{
			Dimensions: []string{"resource"},
			Filter:     actionlog.Filter{Resource: testResource1},
		})

		require.Len(t, set, 1)
		require.Equal(t, float64(3), set[0].Metrics["count"])
	})

	t.Run("time range filter", func(t *testing.T) {
		set := report(actionlog.ReportRequest{
			Dimensions: []string{"resource"},
			Filter:     actionlog.Filter{FromTimestamp: &day2},
		})

		require.Len(t, set, 2)
		for _, row := range set {
			require.Equal(t, float64(1), row.Metrics["count"])
		}
	})

	t.Run("unknown dimension errors", func(t *testing.T) {
		_, err := svc.Report(ctx, actionlog.ReportRequest{Dimensions: []string{"bogus"}})
		require.ErrorContains(t, err, `unknown report dimension "bogus"`)
	})
}
