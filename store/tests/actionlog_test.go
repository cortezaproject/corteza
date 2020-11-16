package tests

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testActionlog(t *testing.T, s store.Actionlogs) {
	var (
		ctx = context.Background()

		makeNew = func(id uint64) *actionlog.Action {
			// minimum data set for new user
			return &actionlog.Action{
				ID:        id,
				Timestamp: *now(),
				Action:    "test-action",
				Resource:  "test-resource",
			}
		}

		stringifySetRange = func(set actionlog.ActionSet) string {
			if len(set) == 0 {
				return ""
			}

			var out = strconv.FormatUint(set[0].ID, 10)

			if len(set) > 1 {
				out += ".." + strconv.FormatUint(set[len(set)-1].ID, 10)
			}

			return out
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, actionlog.ActionSet) {
			req := require.New(t)
			req.NoError(s.TruncateActionlogs(ctx))

			set := make([]*actionlog.Action, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(uint64(i))
			}

			req.NoError(s.CreateActionlog(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateActionlogs(ctx))
		action := &actionlog.Action{
			ID:        42,
			Timestamp: time.Now(),
			Resource:  "resource",
			Action:    "action",
		}
		req.NoError(s.CreateActionlog(ctx, action))
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by resource", func(t *testing.T) {
			req, set := truncAndFill(t, 5)

			set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Resource: "test-resource"})
			req.NoError(err)
			req.Len(set, 5)
		})

		t.Run("by action", func(t *testing.T) {
			req, set := truncAndFill(t, 5)

			set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Action: "test-action"})
			req.NoError(err)
			req.Len(set, 5)
		})

		t.Run("with paging", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.TruncateActionlogs(ctx))

			set := []*actionlog.Action{
				makeNew(1),
				makeNew(2),
				makeNew(3),
				makeNew(4),
				makeNew(5),
				makeNew(6),
				makeNew(7),
				makeNew(8),
				makeNew(9),
			}

			req.NoError(s.CreateActionlog(ctx, set...))
			f := actionlog.Filter{}

			// Fetch first page
			f.Limit = 1
			set, f, err := store.SearchActionlogs(ctx, s, f)
			req.NoError(err)
			req.Len(set, int(f.Limit))
			req.Equal("9", stringifySetRange(set))

			f.Limit = 3
			f.BeforeActionID = set[0].ID
			set, f, err = store.SearchActionlogs(ctx, s, f)
			req.NoError(err)
			req.Len(set, int(f.Limit))
			req.Equal("8..6", stringifySetRange(set))

			f.Limit = 9
			f.BeforeActionID = set[2].ID
			set, f, err = store.SearchActionlogs(ctx, s, f)
			req.NoError(err)
			req.Len(set, 5)
			req.Equal("5..1", stringifySetRange(set))
		})
	})

}
